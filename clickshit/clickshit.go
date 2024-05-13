package clickshit

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"

	// "github.com/ClickHouse/clickhouse-go/v2"
	// "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
)

const (
	ADDR     = "cc-8vbwcevpifh6o038y.clickhouse.ads.aliyuncs.com:3306"
	USERNAME = "meiqia"
	PASSWORD = "cmKInlS9W6YwjEiz"
	DSN      = "tcp://cc-8vbwcevpifh6o038y.clickhouse.ads.aliyuncs.com:3306?compress=true&username=meiqia&password=cmKInlS9W6YwjEiz"
)

var (
	click *sqlx.DB
	// sqlDB *sql.DB
	ctx = context.Background()
	// conn  driver.Conn

	dailer = &net.Dialer{
		Timeout: 3 * time.Second,
	}
)

func init() {
	sqlxDB, err := sqlx.Open("clickhouse", DSN)
	if err != nil {
		panic(err)
	}
	if err := sqlxDB.Ping(); err != nil {
		log.Fatalf("ping failed: %v\n", err)
	}

	click = sqlxDB
}

// func init() {
// 	c, err := clickhouse.Open(&clickhouse.Options{
// 		Addr: []string{ADDR},
// 		Auth: clickhouse.Auth{
// 			Database: "report",
// 			Username: USERNAME,
// 			Password: PASSWORD,
// 		},
// 		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
// 			return dailer.DialContext(ctx, "tcp", addr)
// 		},
// 		Debug: false,
// 		// Debugf: func(format string, v ...any) {
// 		// 	fmt.Printf(format+"\n", v...)
// 		// },
// 		Settings: clickhouse.Settings{
// 			"max_execution_time": 60,
// 		},
// 		Compression: &clickhouse.Compression{
// 			Method: clickhouse.CompressionLZ4,
// 		},
// 		DialTimeout:          time.Second * 30,
// 		MaxOpenConns:         5,
// 		MaxIdleConns:         5,
// 		ConnMaxLifetime:      time.Duration(10) * time.Minute,
// 		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
// 		BlockBufferSize:      10,
// 		MaxCompressionBuffer: 10240,
// 		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
// 			Products: []struct {
// 				Name    string
// 				Version string
// 			}{
// 				{Name: "my-app", Version: "0.1"},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		log.Fatalf("damn1: %v\n", err)
// 	}

// 	if err := c.Ping(context.Background()); err != nil {
// 		log.Fatalf("damn2: %v\n", err)
// 	}

// 	conn = c
// }

// func init() {
// 	op, err := clickhouse.ParseDSN(DSN)
// 	if err != nil {
// 		log.Fatalf("parse dsn fail: %v\n", err)
// 	}

// 	sqlDB = clickhouse.OpenDB(&clickhouse.Options{
// 		Addr: op.Addr,
// 		Auth: clickhouse.Auth{
// 			Database: "report",
// 			Username: op.Auth.Username,
// 			Password: op.Auth.Password,
// 		},
// 		Settings: clickhouse.Settings{
// 			"max_execution_time": 60,
// 			"send_logs_level":    "trace",
// 		},
// 		DialTimeout: time.Second * 30,
// 		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
// 			var d net.Dialer
// 			return d.DialContext(ctx, "tcp", addr)
// 			// return dailer.DialContext(ctx, "tcp", addr)
// 		},
// 		Compression: &clickhouse.Compression{
// 			Method: clickhouse.CompressionLZ4,
// 		},
// 		Debug:                false,
// 		Protocol:             clickhouse.Native,
// 		BlockBufferSize:      10,
// 		MaxCompressionBuffer: 10240,
// 		// ConnMaxLifetime:      time.Duration(10) * time.Minute,
// 		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
// 		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
// 			Products: []struct {
// 				Name    string
// 				Version string
// 			}{
// 				{Name: "hikari-report", Version: "0.1"},
// 			},
// 		},
// 	})
// 	sqlDB.SetMaxIdleConns(5)
// 	sqlDB.SetMaxOpenConns(10)
// 	sqlDB.SetConnMaxLifetime(time.Hour)

//		ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
//			"max_block_size": 10,
//		}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
//			fmt.Println("progress: ", p)
//		}))
//		if err := sqlDB.PingContext(ctx); err != nil {
//			if exception, ok := err.(*clickhouse.Exception); ok {
//				fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
//			}
//			log.Fatalf("ping fail: %v\n", err)
//		}
//	}
func clickhouseTime(t time.Time) string {
	return t.Format(`2006-01-02T15:04:05.999999`)
}

func Fetch(ctx context.Context, sql string, args ...any) (ds []map[string]any, err error) {
	rows, err := click.QueryxContext(ctx, sql, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("click report queryx: %w", err)
	}

	for rows.Next() {
		rowMap := make(map[string]any)
		if err := rows.MapScan(rowMap); err != nil {
			return nil, fmt.Errorf("rows map scan: %w", err)
		}

		for k, v := range rowMap {
			switch val := v.(type) {
			case time.Time:
				rowMap[k] = clickhouseTime(val)
			case []uint8:
				if len(val) == 0 {
					delete(rowMap, k)
				} else {
					newSlice := make([]any, 0, len(val))
					for _, it := range val {
						newSlice = append(newSlice, it)
					}
					rowMap[k] = newSlice
				}
			}
		}

		ds = append(ds, rowMap)
	}

	return
}

func demo1() {
	var entID int64 = 10
	trafficType := 1
	visitID := "2dG4t5KoIL9tMFkL9zT8eLger53"
	createdTS := `2024-03-05T04:00:00`
	// trackID := "2bLIt1spHIXXoJ9QaUE6KuKHXFc"
	// createdTS := clickhouseTime(time.Now().UTC().Add(-30 * 24 * time.Hour))

	sql := `
		SELECT conv_id, agent_type FROM report.visit_conv_distributed
		WHERE ent_id = ? AND traffic_type = ? AND visit_id = ? AND sign = 1 and created_on >= toDateTime64(?, 6)
		ORDER BY version DESC
		LIMIT 1
	`

	var data []struct {
		ConvID    int64 `db:"conv_id"`
		AgentType int64 `db:"agent_type"`
	}

	if err := click.Select(&data, sql, entID, trafficType, visitID, createdTS); err != nil {
		log.Printf("click query fail: %v\n", err)
		return
	}

	if len(data) == 0 {
		log.Println("empty daty")
		return
	}

	log.Printf("got %+v\n", data[0])
}

func demo1_2() {
	// var entID int64 = 10
	// trafficType := 1
	// visitID := "2dG4t5KoIL9tMFkL9zT8eLger53"
	// createdTS := `2024-03-05T04:00:00`

	// sql := `
	// 	SELECT * FROM report.visit_conv_distributed
	// 	WHERE ent_id = ? AND traffic_type = ? AND visit_id = ? AND sign = 1 and created_on >= toDateTime64(?, 6)
	// 	ORDER BY version DESC
	// 	LIMIT 1
	// `

	// var data []struct {
	// 	ConvID    int64 `db:"conv_id"`
	// 	AgentType int64 `db:"agent_type"`
	// }

	// click.QueryContext()
	// click.Query()
	// click.QueryRow()
	// click.QueryRowContext()
	// click.QueryRowx()
	// click.QueryRowxContext()
	// click.Queryx()

	// data, err := Fetch(ctx, sql, entID, trafficType, visitID, createdTS)
	sql := `
	SELECT * FROM report.visit_conv_distributed
	WHERE ent_id = ? AND conv_id = ? AND sign = 1 AND version = 5
	ORDER BY version DESC
	LIMIT 1
`
	data, err := Fetch(ctx, sql, 10, 15989) // 12, 16022)
	if err != nil {
		log.Printf("click query fail: %v\n", err)
		return
	}

	if len(data) > 0 {
		log.Printf("got %+v\n", data[0])

		s := data[0]["redirects.to_agent_type"]
		m := map[string]any{"shit": s}
		dat, _ := json.Marshal(m)
		dd, ok := s.([]any)

		log.Println(string(dat), ok, dd, isEmptySlice(s), s == nil)

		x := data[0]["first_msg_created_on"]
		t, ok := x.(time.Time)
		log.Println("time", t, ok)
	}

	bs := []uint8{65, 66}
	var ns []any
	for _, v := range bs {
		ns = append(ns, v)
	}

	m := map[string]any{"fuck": ns}
	dat, _ := json.Marshal(m)
	log.Println("got", string(dat))

	first := ns[0]
	kk, ok := first.(uint8)
	log.Println(kk, ok)

	// var xx any = []int8{}
	// log.Println("fuck", reflect.ValueOf(xx).IsZero())
	// var kk []any
	// for i := 0; i < 0; i++ {
	// 	kk = append(kk, reflect.ValueOf(xx).Index(0).Interface())
	// }
	// yy, _ := xx.([]uint8)
	// var zz []int8
	// for _, y := range yy {
	// 	zz = append(zz, int8(y))
	// }
	// // q := zz[0]
	// // pp, ok := q.(uint8)
	// // log.Println(pp, ok)
	// m := map[string]any{"shit": kk}
	// dat, _ := json.Marshal(m)
	// log.Println(string(dat))

	// log.Println(reflect.ValueOf(xx).Index(0).Interface())
}

func isEmptySlice(a any) bool {
	if a == nil {
		return true
	}

	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Slice {
		return v.Len() == 0
	}

	return false
}

// 更推荐 demo1_2
// func demo2() {
// 	conn, err := click.Conn(ctx)
// 	if conn != nil {
// 		defer conn.Close()
// 	}
// 	if err != nil {
// 		log.Fatalf("conn err: %v\n", err)
// 	}

// 	var entID int64 = 10
// 	trafficType := 1
// 	visitID := "2dG4t5KoIL9tMFkL9zT8eLger53"
// 	createdTS := `2024-03-05T04:00:00`

// 	sql := `
// 		SELECT * FROM report.visit_conv_distributed
// 		WHERE ent_id = ? AND traffic_type = ? AND visit_id = ? AND sign = 1 and created_on >= toDateTime64(?, 6)
// 		ORDER BY version DESC
// 		LIMIT 1
// 	`

// 	rows, err := conn.QueryContext(ctx, sql, entID, trafficType, visitID, createdTS)
// 	if rows != nil {
// 		defer rows.Close()
// 	}
// 	if err != nil {
// 		log.Fatalf("query rows err: %v\n", err)
// 	}

// 	ds, err := SQLRowsToMap(rows)
// 	if err != nil {
// 		log.Fatalf("query data fail: %v\n", err)
// 	}

// 	if len(ds) > 0 {
// 		log.Printf("get %+v\n", ds[0])
// 	}
// }

// func demo3() {
// 	var entID int64 = 10
// 	trafficType := 1
// 	visitID := "2dG4t5KoIL9tMFkL9zT8eLger53"
// 	createdTS := `2024-03-05T04:00:00`

// 	sql := `
// 		SELECT conv_id, agent_type FROM report.visit_conv_distributed
// 		WHERE ent_id = ? AND traffic_type = ? AND visit_id = ? AND sign = 1 and created_on >= toDateTime64(?, 6)
// 		ORDER BY version DESC
// 		LIMIT 1
// 	`
// 	ctx := clickhouse.Context(
// 		context.Background(),
// 		clickhouse.WithSettings(clickhouse.Settings{"max_block_size": 2}),
// 		clickhouse.WithProgress(func(p *clickhouse.Progress) { fmt.Println("----> ", p) }),
// 	)
// 	rows, err := conn.Query(ctx, sql, entID, trafficType, visitID, createdTS)
// 	defer conn.Close()
// 	if err != nil {
// 		log.Fatalf("query rows err: %v\n", err)
// 	}
// 	defer rows.Close()

// 	data, err := DriverRows2Map(rows)
// 	if err != nil {
// 		log.Fatalf("translate driver rows fail: %v\n", err)
// 	}
// 	if len(data) > 0 {
// 		log.Printf("get: %v\n", len(data[0]))
// 	}

// 	var sts []struct {
// 		ConvID    uint64 `ch:"conv_id"`
// 		AgentType uint8  `ch:"agent_type"`
// 	}

// 	if err := conn.Select(ctx, &sts, sql, entID, trafficType, visitID, createdTS); err != nil {
// 		log.Fatalf("select failed: %v\n", err)
// 	}

// 	log.Printf("st: %v\n", sts)
// }

// func DriverRows2Map(rows driver.Rows) (ds []map[string]any, err error) {
// 	columns := rows.Columns()
// 	if len(columns) == 0 {
// 		return
// 	}
// 	columnTypes := rows.ColumnTypes()

// 	for rows.Next() {
// 		values := make([]any, len(columns))
// 		for i := range values {
// 			values[i] = reflect.New(columnTypes[i].ScanType()).Interface()
// 		}
// 		if err := rows.Scan(values...); err != nil {
// 			return nil, fmt.Errorf("scan row data: %w", err)
// 		}

// 		row := make(map[string]any, len(columns))
// 		for i, column := range columns {
// 			if values[i] == nil {
// 				continue
// 			}
// 			row[column] = t(values[i])
// 		}
// 		ds = append(ds, row)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("rows last check: %w", err)
// 	}

// 	return
// }

// func t(v any) any {
// 	switch val := v.(type) {
// 	case *int8:
// 		return *val
// 	case *int16:
// 		return *val
// 	case *int32:
// 		return *val
// 	case *int64:
// 		return *val
// 	case *uint8:
// 		return *val
// 	case *uint16:
// 		return *val
// 	case *uint32:
// 		return *val
// 	case *uint64:
// 		return *val
// 	case *string:
// 		return *val
// 	case *time.Time:
// 		return *val
// 	case *[]int8:
// 		return *val
// 	case *[]int16:
// 		return *val
// 	case *[]int32:
// 		return *val
// 	case *[]int64:
// 		return *val
// 	case *[]uint8:
// 		return *val
// 	case *[]uint16:
// 		return *val
// 	case *[]uint32:
// 		return *val
// 	case *[]uint64:
// 		return *val
// 	case *[]string:
// 		return *val
// 	case *[]time.Time:
// 		return *val
// 	default:
// 		return v
// 	}
// }

// func demo4() {
// 	var entID int64 = 10
// 	trafficType := 1
// 	visitID := "2dG4t5KoIL9tMFkL9zT8eLger53"
// 	createdTS := `2024-03-05T04:00:00`
// 	// 占位符用 ？ 和 $1 都可以
// 	sql := `
// 		SELECT * FROM report.visit_conv_distributed
// 		WHERE ent_id = $1 AND traffic_type = $2 AND visit_id = $3 AND sign = 1 and created_on >= toDateTime64($4, 6)
// 		ORDER BY version DESC
// 		LIMIT 1
// 	`

// 	ctx := clickhouse.Context(
// 		context.Background(),
// 		clickhouse.WithSettings(clickhouse.Settings{"max_block_size": 2}),
// 		clickhouse.WithProgress(func(p *clickhouse.Progress) { fmt.Println("----> ", p) }),
// 		clickhouse.WithProfileInfo(func(p *clickhouse.ProfileInfo) { fmt.Println("profile info: ", p) }),
// 		// clickhouse.WithLogs(func(log *clickhouse.Log) { fmt.Println("log info: ", log) }),
// 	)

// 	rows, err := sqlDB.QueryContext(ctx, sql, entID, trafficType, visitID, createdTS)
// 	if err != nil {
// 		log.Fatalf("query rows err: %v\n", err)
// 	}
// 	defer rows.Close()

// 	data, err := SQLRowsToMap(rows)
// 	if err != nil {
// 		log.Fatalf("translate fail: %v\n", err)
// 	}

// 	if len(data) > 0 {
// 		log.Printf("get: %v\n", len(data[0]))
// 	}
// }

func SQLRowsToMap(rows *sql.Rows) (ds []map[string]any, err error) {
	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		err = fmt.Errorf("get columns: %w", err)
		return
	}

	// 迭代每一行数据
	for rows.Next() {
		// 创建一个用于存储列值的切片
		values := make([]any, len(columns))
		for i := range values {
			values[i] = new(any)
		}
		// 将列值扫描到values切片中
		err := rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("scan values: %w", err)
		}

		// 创建一个用于存储结果的空map
		row := make(map[string]any)

		// 将列名和对应的值存储到结果map中
		for i, column := range columns {
			if values[i] == nil {
				continue
			}
			row[column] = *(values[i].(*any))
		}

		ds = append(ds, row)
	}

	// 检查是否有错误发生
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows last check: %w", err)
	}

	return
}

type MM struct {
	Name string
}

func (m *MM) IsOk() bool {
	return m != nil && len(m.Name) > 2
}

func Main() {
	var m *MM
	log.Println(m.IsOk())
	demo1_2()
}
