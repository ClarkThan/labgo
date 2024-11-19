package clickshit

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"sync"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/singleflight"
)

const (
	ADDR     = "cc-8vbwcevpifh6o038y.clickhouse.ads.aliyuncs.com:3306"
	USERNAME = "meiqia"
	PASSWORD = "cmKInlS9W6YwjEiz"
	DSN      = "tcp://cc-8vbwcevpifh6o038y.clickhouse.ads.aliyuncs.com:3306?compress=true&username=meiqia&password=cmKInlS9W6YwjEiz"
	// PRD      = "tcp://cc-8vb5j516ra1re2uh7.clickhouse.ads.aliyuncs.com:3306?compress=true&username=salesadmin&password=q9K3gHbBGYkliEpD"
)

var (
	click *sqlx.DB
	// clickDB *sqlx.DB
	ctx = context.Background()
)

func init() {
	sqlxDB, err := sqlx.Open("clickhouse", DSN)
	if err != nil {
		panic(err)
	}
	if err := sqlxDB.Ping(); err != nil {
		log.Fatalf("ping failed: %v\n", err)
	}

	sqlxDB.SetMaxIdleConns(30)
	sqlxDB.SetMaxOpenConns(72)
	sqlxDB.SetConnMaxLifetime(290 * time.Second)
	click = sqlxDB

	// clickDB, err = sqlx.Open("clickhouse", PRD)
	// if err != nil {
	// 	panic(err)
	// }

	// if err := clickDB.Ping(); err != nil {
	// 	log.Fatalf("ping prod failed: %v\n", err)
	// }
}

func clickhouseTime(t time.Time) string {
	return t.Format(`2006-01-02T15:04:05.999999`)
}

func SqlxRowsToMap(rows *sqlx.Rows) (ds []map[string]any, err error) {
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

func Fetch(ctx context.Context, sql string, args ...any) (ds []map[string]any, err error) {
	rows, err := click.QueryxContext(ctx, sql, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("click report queryx: %w", err)
	}

	return SqlxRowsToMap(rows)
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
func demo2() {
	conn, err := click.Conn(ctx)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		log.Fatalf("conn err: %v\n", err)
	}

	var entID int64 = 10
	trafficType := 1
	visitID := "2dG4t5KoIL9tMFkL9zT8eLger53"
	createdTS := `2024-01-05T04:00:00`

	sql := `
		SELECT * FROM report.visit_conv_distributed
		WHERE ent_id = ? AND visit_id = ? AND traffic_type = ? AND sign = 1 and created_on >= toDateTime64(?, 6)
		ORDER BY version DESC
		LIMIT 1
	`

	rows, err := conn.QueryContext(ctx, sql, entID, visitID, trafficType, createdTS)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		log.Fatalf("query rows err: %v\n", err)
	}

	ds, err := SQLRowsToMap(rows)
	if err != nil {
		log.Fatalf("query data fail: %v\n", err)
	}

	if len(ds) > 0 {
		log.Printf("get %+v\n", ds[0])
	} else {
		log.Println("empty data")
	}
}

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

type RepoOverview4evaluateResp struct {
	GoodEvalNum   int64 `db:"good_eval_num"`   // 好评数
	MedEvalNum    int64 `db:"med_eval_num"`    // 中评数
	BadEvalNum    int64 `db:"bad_eval_num"`    // 差评数
	InviteEvalNum int64 `db:"invite_eval_num"` // 邀请评价数
	ConvNum       int64 `db:"conv_num"`        // 对话数
}

func demo7() {
	sqlStmt := `SELECT sum(visit_conv_distributed.is_good_conv * visit_conv_distributed.sign) AS good_eval_num,
sum(visit_conv_distributed.is_medium_conv * visit_conv_distributed.sign) AS med_eval_num,
sum(visit_conv_distributed.is_bad_conv * visit_conv_distributed.sign) AS bad_eval_num,
sum(visit_conv_distributed.is_invite_evaluate * visit_conv_distributed.sign) AS invite_eval_num,
sum(visit_conv_distributed.is_new_conv * visit_conv_distributed.sign) AS conv_num
FROM report.visit_conv_distributed WHERE visit_conv_distributed.ent_id = 10 
AND visit_conv_distributed.conv_created_on >= toDateTime64('2024-09-01T02:00:00',6) 
AND visit_conv_distributed.conv_created_on < toDateTime64('2024-09-04T03:05:00',6) 
AND visit_conv_distributed.is_effective = 1  
HAVING sum(sign) > 0`

	// AND visit_conv_distributed.agent_id in (2090324,1953238,1964086,1969204,1980508,1980510,1980978,1997510,2005766)

	var resp RepoOverview4evaluateResp

	err := click.GetContext(ctx, &resp, sqlStmt)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Fatalf("overview evaluate sql: %s -> %v", sqlStmt, err)
	}

	fmt.Println(resp)
}

func demo8() {
	g := new(singleflight.Group)
	for i := 0; i < 2; i++ {
		_, _, shared := g.Do("demo7", func() (any, error) {
			demo7()
			return nil, nil
		})
		fmt.Println(shared)
	}
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 2; j++ {
				label := "demo7"
				_, _, shared := g.Do(label, func() (any, error) {
					demo7()
					return nil, nil
				})
				fmt.Println(label, shared)
			}
		}()
	}
	wg.Wait()
	fmt.Println("concurrent process end")
}

func demo9() {
	// var m *MM
	// log.Println(m.IsOk())
	// demo1_2()
	var wg sync.WaitGroup
	for i := 0; i < 120; i++ {
		wg.Add(1)
		go func() {
			random := rand.New(rand.NewSource(time.Now().UnixNano()))
			defer wg.Done()
			for i := 0; i < 200; i++ {
				demo7()
				time.Sleep(time.Duration(random.Int31n(300)) * time.Millisecond)
			}
		}()
	}
	wg.Wait()
	fmt.Println("concurrent process end")
}

func demo10() {
	conn, err := click.Conn(ctx)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		log.Fatalf("conn err: %v\n", err)
	}

	entID := 10
	startTS := "2024-11-15 16:00:00"
	endTS := "2024-11-19T16:00:00"
	agentIDs := []int64{9112, 1456, 8997, 9106, 9124, 9125, 2107440, 2107442, 2107456, 2107465}

	// clickhouse/v2 不支持 hasAny([?], redirects.from_agent_id)
	sql := `
		SELECT 
  			toStartOfHour(created_on) as hour, sum(mbot_match_cnt * sign * is_new_visit) as match_cnt, 
  			sumIf(sign, agent_type == 3 OR has(redirects.from_agent_type, 3)) as ini_conv_cnt, 
  			sumIf(is_effective * sign, agent_type == 3 OR has(redirects.from_agent_type, 3)) as effective_conv_cnt, 
  			sum(is_mbot_conv_complete * sign) as fin_conv_cnt, 
  			sumIf(accurate_clue_cnt * sign, agent_type == 3 OR has(redirects.from_agent_type, 3)) as clue_cnt, 
  			sumIf(sign, has(redirects.from_agent_type, 3)) as direct_cnt 
		FROM report.visit_conv_distributed
		WHERE ent_id = ? and conv_ended_on >= toDateTime64(?, 6) 
		and conv_ended_on < toDateTime64(?, 6) 
		and (agent_id in (?) or hasAny(?, redirects.from_agent_id))
		GROUP BY hour HAVING sum(sign) > 0;
	`

	rows, err := conn.QueryContext(ctx, sql, entID, startTS, endTS, agentIDs, agentIDs)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		log.Fatalf("query rows err: %v\n", err)
	}

	ds, err := SQLRowsToMap(rows)
	if err != nil {
		log.Fatalf("query data fail: %v\n", err)
	}

	if len(ds) > 0 {
		log.Printf("get %+v\n", ds[0])
	} else {
		log.Println("empty data")
	}
}

func demo11() {
	entID := 10
	startTS := "2024-11-15 16:00:00"
	endTS := "2024-11-19T16:00:00"
	agentIDs := []int64{9112, 1456, 8997, 9106, 9124, 9125, 2107440, 2107442, 2107456, 2107465}
	values := map[string]any{"ent_id": entID, "start": startTS, "end": endTS, "agent_ids": agentIDs}

	// clickhouse/v2 不支持 hasAny([:agent_ids] 的写法
	sql := `
		SELECT 
  			toStartOfHour(created_on) as hour, sum(mbot_match_cnt * sign * is_new_visit) as match_cnt, 
  			sumIf(sign, agent_type == 3 OR has(redirects.from_agent_type, 3)) as ini_conv_cnt, 
  			sumIf(is_effective * sign, agent_type == 3 OR has(redirects.from_agent_type, 3)) as effective_conv_cnt, 
  			sum(is_mbot_conv_complete * sign) as fin_conv_cnt, 
  			sumIf(accurate_clue_cnt * sign, agent_type == 3 OR has(redirects.from_agent_type, 3)) as clue_cnt, 
  			sumIf(sign, has(redirects.from_agent_type, 3)) as direct_cnt 
		FROM report.visit_conv_distributed
		WHERE ent_id = :ent_id and conv_ended_on >= toDateTime64(:start, 6) 
		and conv_ended_on < toDateTime64(:end, 6) 
		and (agent_id in (:agent_ids) or hasAny(:agent_ids, redirects.from_agent_id))
		GROUP BY hour HAVING sum(sign) > 0;
	`

	rows, err := click.NamedQueryContext(ctx, sql, values)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		log.Fatalf("query rows err: %v\n", err)
	}

	ds, err := SqlxRowsToMap(rows)
	if err != nil {
		log.Fatalf("query data fail: %v\n", err)
	}

	if len(ds) > 0 {
		log.Printf("get %+v\n", ds[0])
	} else {
		log.Println("empty data")
	}
}

func demo12() {
	sql := `SELECT * FROM salesadmin.trigger_d_ limit 1;`
	ds, err := Fetch(ctx, sql)
	if err != nil {
		log.Fatalf("query data fail: %v\n", err)
	}

	log.Printf("got: %v\n", ds)
}

/*
CREATE TABLE privatization_cafe_activity.activity_lt
(
    `id` String COMMENT '主键ID',
    `app_id` String COMMENT '应用ID',
    `module_code` String COMMENT '模块编码',
    `module_name` String COMMENT '模块名称',
    `ent_id` String COMMENT '企业ID',
    `ent_code` String COMMENT '企业Code',
    `ent_name` String COMMENT '企业名称',
    `operator_id` String COMMENT '操作人ID',
    `operator_name` String COMMENT '操作人名称',
    `action_event` String COMMENT '行为事件',
    `action_event_name` String COMMENT '行为事件名称',
    `action_time` UInt64 COMMENT '行为发生的时间戳',
    `action_data` String COMMENT '行为内容（JSON格式）',
    `create_time` UInt64 COMMENT '保存数据的时间戳',
    `ds` UInt32 COMMENT '日期（分区使用）'
)
*/

var LocalTableName string = "privatization_cafe_activity.activity_lt"

func demo13() {
	tx, err := click.Begin()
	if err != nil {
		log.Fatalf("begin tx err: %v\n", err)
	}

	// clickhouse/v1 居然能容忍 下面的语句 （多了 values）
	stmt, err := tx.Prepare(`insert into ` + LocalTableName +
		` values (id, app_id, module_code, module_name, ent_id, ent_code,ent_name,operator_id,operator_name,action_event, action_event_name,action_time,action_data,create_time,ds)
		VALUES (?, ?,?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?, ? );`)
	if err != nil {
		log.Fatalf("prepare stmt err: %v\n", err)
	}
	if _, err := stmt.Exec(
		"3b60927592ff49ad8a2f74cbfdc87eee",
		"22222",
		"SALE_ENTERPRISE",
		"售卖-企业",
		"65b31d6724701833711a8f63",
		"77",
		"xlsx",
		"29g9pbgr8wlkz3do6tm61ly74s3ki2lw",
		"廖明双",
		"USER_LOGIN",
		"通行证-登录",
		1706178865799,
		`{"belong":"system"}`,
		1706178865842,
		20240125,
	); err != nil {
		log.Fatalf("exec stmt err: %v\n", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("commit tx err: %v\n", err)
	}

	log.Println("insert success")
}

func Main() {
	demo13()
}
