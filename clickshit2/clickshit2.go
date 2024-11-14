package clickshit2

// var test int = 11

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const (
	ADDR     = "cc-8vbwcevpifh6o038y.clickhouse.ads.aliyuncs.com:3306"
	USERNAME = "meiqia"
	PASSWORD = "cmKInlS9W6YwjEiz"
	DSN      = "tcp://cc-8vbwcevpifh6o038y.clickhouse.ads.aliyuncs.com:3306?compress=true&username=meiqia&password=cmKInlS9W6YwjEiz"
	// PRD      = "tcp://cc-8vb5j516ra1re2uh7.clickhouse.ads.aliyuncs.com:3306?compress=true&username=salesadmin&password=q9K3gHbBGYkliEpD"
)

var (
	sqlDB *sql.DB
	conn  driver.Conn

	dailer = &net.Dialer{
		Timeout: 3 * time.Second,
	}

	// click *sqlx.DB
	// clickDB *sqlx.DB
	ctx = context.Background()
)

func init() {
	init1()
	init2()
}

func init1() {
	c, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{ADDR},
		Auth: clickhouse.Auth{
			Database: "report",
			Username: USERNAME,
			Password: PASSWORD,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			return dailer.DialContext(ctx, "tcp", addr)
		},
		Debug: false,
		// Debugf: func(format string, v ...any) {
		// 	fmt.Printf(format+"\n", v...)
		// },
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "my-app", Version: "0.1"},
			},
		},
	})
	if err != nil {
		log.Fatalf("damn1: %v\n", err)
	}

	if err := c.Ping(context.Background()); err != nil {
		log.Fatalf("damn2: %v\n", err)
	}

	conn = c
}

func init2() {
	op, err := clickhouse.ParseDSN(DSN)
	if err != nil {
		log.Fatalf("parse dsn fail: %v\n", err)
	}

	sqlDB = clickhouse.OpenDB(&clickhouse.Options{
		Addr: op.Addr,
		Auth: clickhouse.Auth{
			Database: "report",
			Username: op.Auth.Username,
			Password: op.Auth.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
			"send_logs_level":    "trace",
		},
		DialTimeout: time.Second * 30,
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
			// return dailer.DialContext(ctx, "tcp", addr)
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                false,
		Protocol:             clickhouse.Native,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		// ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "hikari-report", Version: "0.1"},
			},
		},
	})
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 10,
	}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
		fmt.Println("progress: ", p)
	}))
	if err := sqlDB.PingContext(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		log.Fatalf("ping fail: %v\n", err)
	}
}

func demo3() {
	var entID int64 = 10
	trafficType := 1
	visitID := "2dG4t5KoIL9tMFkL9zT8eLger53"
	createdTS := `2024-03-05T04:00:00`

	sql := `
		SELECT conv_id, agent_type FROM report.visit_conv_distributed
		WHERE ent_id = ? AND traffic_type = ? AND visit_id = ? AND sign = 1 and created_on >= toDateTime64(?, 6)
		ORDER BY version DESC
		LIMIT 1
	`
	ctx := clickhouse.Context(
		context.Background(),
		clickhouse.WithSettings(clickhouse.Settings{"max_block_size": 2}),
		clickhouse.WithProgress(func(p *clickhouse.Progress) { fmt.Println("----> ", p) }),
	)
	rows, err := conn.Query(ctx, sql, entID, trafficType, visitID, createdTS)
	defer conn.Close()
	if err != nil {
		log.Fatalf("query rows err: %v\n", err)
	}
	defer rows.Close()

	data, err := DriverRows2Map(rows)
	if err != nil {
		log.Fatalf("translate driver rows fail: %v\n", err)
	}
	if len(data) > 0 {
		log.Printf("get: %v\n", len(data[0]))
	}

	var sts []struct {
		ConvID    uint64 `ch:"conv_id"`
		AgentType uint8  `ch:"agent_type"`
	}

	if err := conn.Select(ctx, &sts, sql, entID, trafficType, visitID, createdTS); err != nil {
		log.Fatalf("select failed: %v\n", err)
	}

	log.Printf("st: %v\n", sts)
}

func DriverRows2Map(rows driver.Rows) (ds []map[string]any, err error) {
	columns := rows.Columns()
	if len(columns) == 0 {
		return
	}
	columnTypes := rows.ColumnTypes()

	for rows.Next() {
		values := make([]any, len(columns))
		for i := range values {
			values[i] = reflect.New(columnTypes[i].ScanType()).Interface()
		}
		if err := rows.Scan(values...); err != nil {
			return nil, fmt.Errorf("scan row data: %w", err)
		}

		row := make(map[string]any, len(columns))
		for i, column := range columns {
			if values[i] == nil {
				continue
			}
			row[column] = t(values[i])
		}
		ds = append(ds, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows last check: %w", err)
	}

	return
}

func t(v any) any {
	switch val := v.(type) {
	case *int8:
		return *val
	case *int16:
		return *val
	case *int32:
		return *val
	case *int64:
		return *val
	case *uint8:
		return *val
	case *uint16:
		return *val
	case *uint32:
		return *val
	case *uint64:
		return *val
	case *string:
		return *val
	case *time.Time:
		return *val
	case *[]int8:
		return *val
	case *[]int16:
		return *val
	case *[]int32:
		return *val
	case *[]int64:
		return *val
	case *[]uint8:
		return *val
	case *[]uint16:
		return *val
	case *[]uint32:
		return *val
	case *[]uint64:
		return *val
	case *[]string:
		return *val
	case *[]time.Time:
		return *val
	default:
		return v
	}
}

type ConvData struct {
	AgentEffectiveMsgNum int64  `db:"agent_effective_msg_num"`
	AgentMsgNum          int64  `db:"agent_msg_num"`
	VisitPageCnt         int64  `db:"visit_page_cnt"`
	VisitID              string `db:"visit_id"`
	Country              string `db:"country"`
}

func demo4() {
	var entID int64 = 10
	trafficType := 1
	visitID := "2ogwNyxQFPk0jBxYYVXbE78Z4gL"
	createdTS := `2024-03-05T04:00:00`
	// 占位符用 ？ 和 $1 都可以
	sql := `
		SELECT agent_effective_msg_num, agent_msg_num, visit_page_cnt, visit_id, country FROM report.visit_conv_distributed
		WHERE ent_id = $1 AND traffic_type = $2 AND visit_id = $3 AND sign = 1 and created_on >= toDateTime64($4, 6)
		ORDER BY version DESC
		LIMIT 1
	`

	ctx := clickhouse.Context(
		context.Background(),
		clickhouse.WithSettings(clickhouse.Settings{"max_block_size": 2}),
		clickhouse.WithProgress(func(p *clickhouse.Progress) { fmt.Println("----> ", p) }),
		clickhouse.WithProfileInfo(func(p *clickhouse.ProfileInfo) { fmt.Println("profile info: ", p) }),
		// clickhouse.WithLogs(func(log *clickhouse.Log) { fmt.Println("log info: ", log) }),
	)

	rows, err := sqlDB.QueryContext(ctx, sql, entID, trafficType, visitID, createdTS)
	if err != nil {
		log.Fatalf("query rows err: %v\n", err)
	}
	defer rows.Close()

	data, err := SQLRowsToMap(rows)
	if err != nil {
		log.Fatalf("translate fail: %v\n", err)
	}

	if len(data) > 0 {
		log.Printf("get: %d, %v\n", len(data[0]), data[0])
	} else {
		log.Println("shit.....")
	}
}

func demo5() {
	// db, err := sql.Open("clickhouse", "clickhouse://cc-8vbwcevpifh6o038y.clickhouse.ads.aliyuncs.com:3306?compress=true&username=meiqia&password=cmKInlS9W6YwjEiz")
	// if err != nil {
	// 	log.Fatalf("opening database failed: %v", err)
	// }
	// defer db.Close()
	// stmt, err := sqlDB.PrepareContext(ctx, "SELECT * FROM system.metrics")
	// if err != nil {
	// 	log.Fatalf("preparing statement failed: %v", err)
	// }

	// rows, err := stmt.QueryContext(ctx)
	rows, err := sqlDB.QueryContext(ctx, "SELECT * FROM system.metrics limit 2")
	if err != nil {
		log.Fatalf("executing preparing statement failed: %v", err)
	}

	// var ms []Metric
	// if err := rows.Scan(&ms); err != nil {
	// 	log.Fatalf("scanning rows failed: %v", err)
	// }
	// fmt.Printf("len: %d,  vals: %+v\n", len(ms), ms)

	columns, err := rows.Columns()
	fmt.Println("row columns:", columns)
	fmt.Println("err:", err)
	dataMap, err := SQLRowsToMap(rows)
	fmt.Println(dataMap)
	fmt.Println(err)
}

type Metric struct {
	Metric      string `ch:"metric"`
	Value       int64  `ch:"value"`
	Description string `ch:"description"`
}

// type Metric struct {
// 	Metric      string
// 	Value       int64
// 	Description string
// }

func demo7() {
	stmt, err := sqlDB.PrepareContext(ctx, "SELECT * FROM system.metrics WHERE metric = $1 LIMIT 1")
	if err != nil {
		log.Fatalf("preparing statement failed: %v", err)
	}
	rows1, err := stmt.QueryContext(ctx, "Query")
	if err != nil {
		log.Fatalf("executing preparing statement1 failed: %v", err)
	}
	var ms1 []Metric
	if err := rows1.Scan(&ms1); err != nil {
		log.Fatalf("scanning rows failed: %v", err)
	}
	fmt.Printf("len: %d,  vals: %+v\n", len(ms1), ms1)

	rows2, err := stmt.QueryContext(ctx, "Merge")
	if err != nil {
		log.Fatalf("executing preparing statement2 failed: %v", err)
	}
	var ms2 []Metric
	if err := rows2.Scan(&ms2); err != nil {
		log.Fatalf("scanning rows failed: %v", err)
	}
	fmt.Printf("len: %d,  vals: %+v\n", len(ms2), ms2)
}

func demo8() {
	rows, err := sqlDB.QueryContext(ctx, "SELECT * FROM system.metrics limit 2")
	if err != nil {
		log.Fatalf("executing preparing statement failed: %v", err)
	}
	data, err := SQLRowsToMap(rows)
	if err != nil {
		log.Fatalf("scanning rows failed: %v", err)
	}
	fmt.Printf("len: %d,  vals: %+v\n", len(data), data)
}

func demo6() {
	var metrics []*Metric
	err := conn.Select(ctx, &metrics, "SELECT * FROM system.metrics limit 5")
	if err != nil {
		log.Fatalf("executing preparing statement failed: %v", err)
		return
	}
	fmt.Printf("len: %d,  vals: %+v\n\n", len(metrics), metrics)

	var m Metric
	row := conn.QueryRow(ctx, "SELECT * FROM system.metrics Where metric != 'shit' limit 1")
	if err := row.ScanStruct(&m); err == nil {
		fmt.Printf("got one: %+v\n", m)
	} else if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("--- got err: %v\n", err)
	}
}

func Main() {
	demo6()
}
