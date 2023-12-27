package dbrshit

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/timonwong/dbr"
	"github.com/timonwong/dbr/dialect"
)

const (
	batchInsertGPTSQL = "INSERT IGNORE INTO `gptbot_corpus` (`ent_id`, `robot_id`, `title`, `content`) VALUES "
)

var (
	sess *dbr.Session
	ctx  context.Context = context.Background()
)

func Main() {
	demo9()
}

func demo1() {
	entID, robotID, docs := genData(9, 5)
	var builder strings.Builder
	holder := "(?,?,?,?),"

	builder.WriteString(batchInsertGPTSQL)
	builder.Grow(len(holder) * 4 * len(docs))

	values := make([]any, 0, 4*len(docs))
	for _, d := range docs {
		values = append(values, entID, robotID, d.Title, d.Content)
		builder.WriteString(holder)
	}
	sql := strings.TrimRight(builder.String(), " ,")

	log.Printf("++++++++ batch insert gptdocs info, sql: %s   value cnt: %d   first:%v", sql, len(values), values[:4])

	result, err := sess.InsertBySql(sql, values...).ExecContext(ctx)
	if err != nil {
		log.Printf("batch insert gptdocs fail: %v", err)
		return
	}

	cnt, _ := result.RowsAffected()
	log.Println("成功数", cnt)
}

func demo2() {
	result, err := sess.UpdateBySql("UPDATE `gptbot_membership` SET `status` = `status` - ?", 10).ExecContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	cnt, _ := result.RowsAffected()
	log.Println("affected ", cnt)
}

func demo3() {
	var flowData []byte
	err := sess.Select("chat_flow").From("settings").Where(dbr.Eq("id", 1)).LoadOneContext(ctx, &flowData)
	if err != nil {
		log.Fatalf("damn it: %v\n", err)
	}
	log.Println(string(flowData))

	dat := []byte(`{"foo":["bar","baz"]}`)
	ret, err := sess.InsertInto("settings").Columns("ent_id", "chat_flow").Values(67930, dat).ExecContext(ctx)
	if err != nil {
		log.Fatalf("insert error: %v\n", err)
	}
	id, err := ret.LastInsertId()
	if err != nil {
		log.Println("last insert id fetching error", err)
	}

	log.Println("new id", id)
}

type TMap map[string]any

func (t TMap) Value() (driver.Value, error) {
	if len(t) == 0 {
		return []byte(`{}`), nil
	}

	return json.Marshal(t)
}

func (t *TMap) Scan(v any) error {
	switch data := v.(type) {
	case []byte:
		return json.Unmarshal(data, &t)
	case string:
		return json.Unmarshal([]byte(data), &t)
	default:
		return fmt.Errorf("invalid field data: %v", v)
	}
}

type Setting struct {
	ID       int64 `json:"id" db:"id"`
	EntID    int64 `json:"ent_id" db:"ent_id"`
	Status   int64 `json:"status" db:"status"`
	ChatFlow TMap  `json:"chat_flow,omitempty" db:"chat_flow"`
}

func demo4() {
	var ss []Setting
	cnt, err := sess.Select("id", "ent_id", "chat_flow", "status").
		From("settings").Where(dbr.Eq("id", []int{11, 5})).Load(&ss)

	if err != nil {
		log.Fatalf("fetching settings failed: %v\n", err)
	}

	log.Println("got ", cnt)
	if dat, err := json.Marshal(ss); err == nil {
		log.Println("stringify: ", string(dat))
	}

	sss := []Setting{
		{EntID: 23, Status: 1},
		{EntID: 45, Status: 0, ChatFlow: TMap(map[string]any{})},
	}
	if dat, err := json.Marshal(sss); err == nil {
		log.Println("sss: ", string(dat))
	}
}

func init() {
	log.Println("mysql orm of dbr run...")
	dbConf := &DBConf{
		Ip:           "127.0.0.1",
		Port:         3306,
		User:         "test",
		Password:     "12345687",
		Database:     "test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ParseTime:    true,
	}

	dsn, err := dbConf.DSN()
	if err != nil {
		log.Fatalln("err: ", err)
	}

	// opens a database
	conn, _ := dbr.Open("mysql", dsn, &dbr.NullEventReceiver{})

	// 设置连接池相关参数
	conn.SetMaxOpenConns(dbConf.MaxOpenConns)
	conn.SetMaxIdleConns(dbConf.MaxIdleConns)
	conn.SetConnMaxLifetime(dbConf.MaxLifetime)
	conn.SetConnMaxIdleTime(dbConf.MaxIdleTime)

	sess = conn.NewSession(nil) // 每次查询都是一个session会话操作
}

// User
type User struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Age  int    `json:"age" db:"age"`
}

// TableName table
func (User) TableName() string {
	return "user"
}

// DBConf DB config
type DBConf struct {
	Ip        string
	Port      int // 默认3306
	User      string
	Password  string
	Database  string
	Charset   string // 字符集 utf8mb4 支持表情符号
	Collation string // 整理字符集 utf8mb4_unicode_ci

	MaxIdleConns int // 空闲pool个数
	MaxOpenConns int // 最大open connection个数

	// sets the maximum amount of time a connection may be reused.
	// 设置连接可以重用的最大时间
	// 给db设置一个超时时间，时间小于数据库的超时时间
	MaxLifetime time.Duration // 数据库超时时间
	MaxIdleTime time.Duration // 最大空闲时间

	// 连接超时/读取超时/写入超时设置
	Timeout      time.Duration // Dial timeout
	ReadTimeout  time.Duration // I/O read timeout
	WriteTimeout time.Duration // I/O write timeout

	ParseTime bool   // 格式化时间类型
	Loc       string // 时区字符串 Local,PRC
}

func (conf *DBConf) DSN() (string, error) {
	if conf.Ip == "" {
		conf.Ip = "127.0.0.1"
	}

	if conf.Port == 0 {
		conf.Port = 3306
	}

	if conf.Charset == "" {
		conf.Charset = "utf8mb4"
	}

	// 默认字符序，定义了字符的比较规则
	if conf.Collation == "" {
		conf.Collation = "utf8mb4_general_ci"
	}

	if conf.Loc == "" {
		conf.Loc = "Local"
	}

	if conf.Timeout == 0 {
		conf.Timeout = 10 * time.Second
	}

	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 5 * time.Second
	}

	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 5 * time.Second
	}

	if conf.MaxLifetime == 0 {
		conf.MaxLifetime = 20 * time.Minute
	}

	if conf.MaxIdleTime == 0 {
		conf.MaxIdleTime = 10 * time.Minute
	}

	// mysql connection time loc.
	loc, err := time.LoadLocation(conf.Loc)
	if err != nil {
		return "", err
	}

	// mysql config
	mysqlConf := mysql.Config{
		User:   conf.User,
		Passwd: conf.Password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", conf.Ip, conf.Port),
		DBName: conf.Database,
		// Connection parameters
		Params: map[string]string{
			"charset": conf.Charset,
		},
		Collation:            conf.Collation,
		Loc:                  loc,               // Location for time.Time values
		Timeout:              conf.Timeout,      // Dial timeout
		ReadTimeout:          conf.ReadTimeout,  // I/O read timeout
		WriteTimeout:         conf.WriteTimeout, // I/O write timeout
		AllowNativePasswords: true,              // Allows the native password authentication method
		ParseTime:            conf.ParseTime,    // Parse time values to time.Time
	}

	return mysqlConf.FormatDSN(), nil
}

type GPTDocItem struct {
	Title   string
	Content string
}

func genData(s, cnt int) (entID, robotID int64, docs []GPTDocItem) {
	entID = 1
	robotID = 23
	for i := s; i < s+cnt; i++ {
		docs = append(docs, GPTDocItem{
			Title:   fmt.Sprintf("%d-title", i),
			Content: fmt.Sprintf("%d-content", i),
		})
	}

	return
}

func demo5() {
	conds := dbr.And(dbr.Eq("ent_id", 1), dbr.Eq("id", 2))
	q := "ba"
	if q != "" {
		conds = dbr.And(conds, dbr.Expr("`title` LIKE ?", fmt.Sprintf(`%%%s%%`, q)))
	}
	var ret []struct {
		ID    int64  `db:"id"`
		Title string `db:"title"`
	}

	sqlBuilder := sess.Select("id", "title").From("foo").Where(conds)
	log.Println("sql: ", ToSQL(sqlBuilder))
	cnt, err := sqlBuilder.Load(&ret)
	if err != nil {
		log.Fatalf("query foo failed: %v\n", err)
	}

	log.Printf("got %d columns\n", cnt)
	log.Printf("ret: %+v\n", ret)
}

func demo6() {
	var provider string
	conds := dbr.And(dbr.Eq("robot_id", 1849587), dbr.Eq("status", 1))
	_, err := sess.Select("provider").From("gptbot_membership").Where(conds).LoadContext(ctx, &provider)
	if err != nil {
		// dead code, LoadContext 不会返回 ErrNotFound
		if errors.Is(err, dbr.ErrNotFound) {
			log.Println("fuck not found")
			return
		}
		log.Fatalf("query provider from gptbot_membership failed: %v\n", err)
	}
	log.Println("provider = ", provider)
}

func demo7() {
	var cnt []int
	err := sess.Select("count(1)").From("gptbot_membership").
		Where(dbr.Eq("status", 1)).GroupBy("provider").LoadOneContext(ctx, &cnt)
	if err != nil {
		// dead code, LoadContext 不会返回 ErrNotFound
		if errors.Is(err, dbr.ErrNotFound) {
			log.Println("fuck not found")
			return
		}
		log.Fatalf("query provider from gptbot_membership failed: %v\n", err)
	}
	log.Printf("provider = %v\n", cnt)
}

func demo8() {
	var cnt int
	err := sess.Select("count(*)").From("gptbot_membership").Where(dbr.Eq("robot_id", 10)).LoadOneContext(ctx, &cnt)
	if err != nil {
		// dead code, LoadContext 不会返回 ErrNotFound
		if errors.Is(err, dbr.ErrNotFound) {
			log.Println("fuck not found")
			return
		}
		log.Fatalf("query provider from gptbot_membership failed: %v\n", err)
	}
	log.Printf("robot_id = %v\n", cnt)
}

func demo9() {
	var robotID int64
	if err := sess.SelectBySql("SELECT `robot_id` FROM `gptbot_membership` WHERE `provider` = ?", "wenxin").
		LoadOneContext(ctx, &robotID); err != nil {
		log.Printf("error = %v\n", err)
	}
	log.Println(robotID)

	var bot struct {
		RobotID  int64  `db:"robot_id"`
		EntID    int64  `db:"ent_id"`
		Provider string `db:"provider"`
	}
	err := sess.Select("robot_id", "ent_id", "provider").From("gptbot_membership").
		Where(dbr.Eq("status", 1)).OrderBy("robot_id DESC").Limit(1).LoadOneContext(ctx, &bot)
	if err != nil {
		log.Printf("damn: %v\n", err)
	}
	log.Printf("ret = %+v\n", bot)

	var dat io.Reader = bytes.NewBuffer([]byte(`fuckyou`))
	data1, _ := io.ReadAll(dat)
	data2, _ := io.ReadAll(dat) // 不能重复读
	if cmp.Equal(data1, data2) {
		log.Fatal("damn", data1, data2)
	}
}

func ToSQL(builder dbr.Builder) string {
	sql, _ := dbr.InterpolateForDialect("?", []any{builder}, dialect.MySQL)
	return sql
}
