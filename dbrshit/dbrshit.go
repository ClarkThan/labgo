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
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/go-sql-driver/mysql"
	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/mapstructure"
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

func init() {
	// initQA("meiqia")
	initLocal()
}

func initQA(database string) {
	// dsn := "test:12345687@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("meiqia:f_xByc=9Dy+ZCbH1@tcp(pc-8vbvpi114t895m715.mysql.polardb.zhangbei.rds.aliyuncs.com:3306)/%s?charset=utf8&parseTime=True&loc=Local", database)

	// opens a database
	conn, err := dbr.Open("mysql", dsn, &dbr.NullEventReceiver{})
	if err != nil {
		log.Fatalln("err: ", err)
	}

	log.Println("mysql orm of dbr run...")

	// 设置连接池相关参数
	conn.SetMaxOpenConns(3)
	conn.SetMaxIdleConns(2)
	// conn.SetConnMaxLifetime(dbConf.MaxLifetime)
	// conn.SetConnMaxIdleTime(dbConf.MaxIdleTime)

	sess = conn.NewSession(nil) // 每次查询都是一个session会话操作
}

func initLocal() {
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
	// var provider string

	type Info struct {
		RobotID  int64  `db:"robot_id"`
		Provider string `db:"provider"`
	}

	var info Info
	_, err := sess.Select("provider,robot_id").From("robot_membership").Where(dbr.Gt("robot_id", 35)).LoadContext(ctx, &info)
	if err != nil {
		// dead code, LoadContext 不会返回 ErrNotFound
		if errors.Is(err, dbr.ErrNotFound) {
			log.Println("fuck not found")
			return
		}
		log.Fatalf("query provider from gptbot_membership failed: %v\n", err)
	}
	log.Println("provider = ", info)
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

func demo10() {
	trackIDs := []string{"1fLJKWKKSuQA4hXE9LeB91hVhBU", "1xFSvRLBJMzwVwrF0ZZpYutSZZY", "1zOjmrd12vyU5lTW5qsXvg6faHA"}
	builder := sess.Select("MAX(id) as last_visit_id").From("visit").
		Where("enterprise_id = ?", 10).
		Where("track_id in ?", trackIDs).
		GroupBy("enterprise_id", "track_id").
		OrderBy("NULL")

	log.Println(ToSQL(builder))

	lastVisitIDs := make([]string, 0, len(trackIDs))
	_, err := builder.Load(&lastVisitIDs)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	log.Printf("last_visit_ids: %v\n", lastVisitIDs)
}

func demo11() {
	ids := []int64{1510, 2160, 1880, 999}
	existIDs := make([]int64, 0, len(ids))
	cnt, err := sess.Select("id").From("enterprise_sub_source").Where(dbr.Eq("id", ids)).LoadContext(ctx, &existIDs)
	if err != nil {
		log.Fatalf("query id fail: %v\n", err)
	}

	log.Printf("exist ids(%d): %v\n", cnt, existIDs)
}

type M struct {
	Name                string         `json:"name"`
	No                  int64          `json:"no"`
	TagIDs              []int          `json:"tag_ids"`
	ClientFirstSendTime *time.Time     `json:"client_first_send_time"`
	Attrs               map[string]any `json:"attrs"`
}

func (m *M) ToMap() map[string]any {
	st := structs.New(m)
	st.TagName = `json`
	ret := st.Map()
	// if m.ClientFirstSendTime == nil {
	// 	delete(ret, "client_first_send_time")
	// }
	return ret
}

func orEmptyList(val any) any {
	// if !reflect.ValueOf(val).IsNil() {
	// 	return val
	// }
	if lst, ok := val.([]any); ok && len(lst) > 0 {
		return lst
	}

	return []any{}
}

func demo12() {
	// demo11()
	// m := map[string]any{
	// 	"name":    "air",
	// 	"no":      23,
	// 	"tag_ids": nil,
	// }

	// now := time.Now().UTC()
	x := M{Name: "hello", No: 0}
	m := x.ToMap()
	t, ok := m["client_first_send_time"].(*time.Time)
	log.Printf("m: %v %t %t\n", m, ok, t == nil)
	mStr, _ := json.Marshal(m)
	log.Println("old json:", string(mStr))

	ret := map[string]any{
		"name":    m["name"],
		"tag_ids": orEmptyList(m["tag_ids"]),
		"no":      m["no"],
		"fuck":    nil,
	}

	for _, k := range []string{"first_msg_created_on", "attrs", "tag_ids", "client_first_send_time", "first_response_wait_time", "first_response_agent_id"} {
		if v, exists := m[k]; exists && v != nil {
			log.Println("add some non nil field", k)
			ret[k] = v
		}
	}

	log.Printf("ret: %v\n", ret)

	for k, v := range ret {
		if v == nil {
			log.Println("del nil kv", k)
			delete(ret, k)
			continue
		}
		// empty := false
		// vv := reflect.ValueOf(v)
		// switch vv.Kind() {
		// case reflect.Slice, reflect.Map, reflect.Pointer, reflect.Interface:
		// 	empty = vv.IsNil()
		// 	// default:
		// 	// 	log.Println(k, "-->", v)
		// }
		// if empty {

		if reflect.ValueOf(v).IsZero() {
			log.Println("del empty kv", k)
			delete(ret, k)
		}
	}
	log.Printf("ret: %v\n", ret)
	retStr, _ := json.Marshal(ret)
	log.Println("json:", string(retStr))
}

func demo13() {
	var maxRank int64
	err := sess.Select("MAX(`rank`)").
		From("selecting_rule").
		Where("enterprise_id = ?", 10).
		LoadOneContext(ctx, &maxRank)

	if err != nil {
		log.Printf("err: %v\n", err)
	}

	log.Println("maxRank", maxRank)
}

type Event struct {
	ID           int64 `json:"id" mapstructure:"id"`
	EnterpriseID int64 `json:"enterprise_id" mapstructure:"enterprise_id"`
	TicketID     int64 `json:"ticket_id" mapstructure:"ticket_id"`
	AgentID      int64 `json:"agent_id" mapstructure:"agent_id"`

	Content     string `json:"content" mapstructure:"content"`
	ContentType string `json:"content_type" mapstructure:"content_type"`
	Type        string `json:"type" mapstructure:"type"`
	Action      string `json:"action" mapstructure:"action"`
	FromType    string `json:"from_type" mapstructure:"from_type"`

	CreatedAt time.Time `json:"created_at" mapstructure:"created_at,asis"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at,asis"`
}

type AgentAvatar struct {
	Avatar string `json:"avatar"`
}

type Reply struct {
	*Event

	MediaURL string       `json:"media_url,omitempty" mapstructure:"media_url"`
	Agent    *AgentAvatar `json:"agent,omitempty" mapstructure:"agent"`
}

func demo14() {
	evt := Reply{
		Event: &Event{
			ID:          134,
			AgentID:     345,
			ContentType: "hello",
		},
		MediaURL: "http://meiqia.com",
	}

	var e any = evt
	data := make(map[string]any)
	if err := mapstructure.Decode(e, &data); err != nil {
		log.Printf("err: %v\n", err)
	}

	bs, _ := json.Marshal(data)

	log.Println(string(bs))
}

type TicketConf struct {
	SDKEnabled bool `json:"sdk_enabled"`
	WebEnabled bool `json:"web_enabled"`
}

func demo15() {
	cfg := TicketConf{true, true}
	// dat := []byte(`{"addresser_name":"","captcha":"close","category":"close","contactRule":"multi","content_fill_type":"placeholder","content_placeholder":"留言遇到的问题和图片","content_title":"留言内容（请具体描述您所遇到的问题，并上传产品和快递面单照片，以及留言您希望的处理方式）","custom_fields":[{"metainfo":[],"name":"tel","placeholder":"请确保您的联系方式准确无误。","required":true,"type":"string"}],"defaultTemplate":"open","email":"close","intro":"亲爱的达漫会员：您好，客服当前不在线，请留言您遇到的具体问题+订单编号以及您的联系方式。\n      人工客服在线时间：08:30-17:30，其余时间人工客服不在线，如有问题可致电4001580008，电话客服服务时间：08:30-21:00。非服务时间（当日17:30至次日8:30）客服将会在次日上班后依次处理，如您着急，请在服务时间内拨打400电话向客服反馈。\n\n自助售后流程： 进入“达漫电商”微信公众号-达漫商城-我的-我的订单-找到对应订单确认收货后即可申请售后（上传问题图片处理更快哦~）","name":"close","permission":"close","qq":"close","tel":"open","upload_image":"open","wechat":"close"}	`)
	var settings any = nil
	err := mapstructure.Decode(settings, &cfg)
	fmt.Printf("%#v\n,  %v", cfg, err)
}

type agentModel struct {
	ID              int64          `db:"id"`
	EntID           int64          `db:"enterprise_id"`
	Avatar          dbr.NullString `db:"avatar"`
	Email           dbr.NullString `db:"email"`
	Nickname        dbr.NullString `db:"nickname"`
	Realname        dbr.NullString `db:"realname"`
	Password        dbr.NullString `db:"password"`
	Privilege       dbr.NullString `db:"privilege"`
	Cellphone       dbr.NullString `db:"cellphone"`
	Status          dbr.NullString `db:"status"`
	PublicCellphone dbr.NullString `db:"public_cellphone"`
	PublicEmail     dbr.NullString `db:"public_email"`
	QQ              dbr.NullString `db:"qq"`
	Signature       dbr.NullString `db:"signature"`
	Telephone       dbr.NullString `db:"telephone"`
	Weixin          dbr.NullString `db:"weixin"`
	Audience        dbr.NullString `db:"audience"`
	Token           dbr.NullString `db:"token"`
	WorkNum         dbr.NullString `db:"work_num"`
	EmailActivated  dbr.NullInt64  `db:"email_activated"`
	GroupID         int64          `db:"group_id"`
	Rank            int64          `db:"rank"`
	ServingLimit    int64          `db:"serving_limit"`
	ReadFeatureID   dbr.NullInt64  `db:"read_feature_id"`
	CreatedOn       *time.Time     `db:"created_on"`
	LastUpdated     *time.Time     `db:"last_updated"`
	LastActiveTime  *time.Time     `db:"last_active_time"`
	LastLoginTime   *time.Time     `db:"last_login_time"`
	DeletedAt       *time.Time     `db:"deleted_at"`
	ServingType     dbr.NullInt64  `db:"serving_type"`
	Phone           dbr.NullString `db:"phone"`
}

type Agent struct {
	ID              int64      `json:"id" db:"id"`
	EntID           int64      `json:"enterprise_id" mapstructure:"ent_id" db:"enterprise_id" kun:"descr=企业ID"`
	Avatar          string     `json:"avatar" mapstructure:"avatar" db:"avatar" kun:"descr=头像"`
	Email           string     `json:"email" mapstructure:"email" db:"email" kun:"descr=客服邮箱"`
	Nickname        string     `json:"nickname" mapstructure:"nickname" db:"nickname" kun:"descr=昵称"`
	Realname        string     `json:"realname" mapstructure:"realname" db:"realname" kun:"descr=真实姓名"`
	Password        string     `json:"-" db:"password"`
	Privilege       string     `json:"privilege" db:"privilege" kun:"descr=权限"`
	Cellphone       string     `json:"cellphone" db:"cellphone"`
	Status          string     `json:"status" db:"status" kun:"descr=在线状态"`
	PublicCellphone string     `json:"public_cellphone" db:"public_cellphone"`
	PublicEmail     string     `json:"public_email" db:"public_email"`
	QQ              string     `json:"qq" db:"qq"`
	Signature       string     `json:"signature" db:"signature"`
	Telephone       string     `json:"telephone" db:"telephone"`
	Weixin          string     `json:"weixin" db:"weixin"`
	Audience        string     `json:"-" db:"audience"`
	Token           string     `json:"token" db:"token"`
	WorkNum         string     `json:"work_num" db:"work_num" kun:"descr=工号"`
	EmailActivated  int64      `json:"email_activated" db:"email_activated" kun:"descr=邮件是否激活"`
	GroupID         int64      `json:"group_id" mapstructure:"group_id" db:"group_id" kun:"descr=分组ID"`
	Rank            int64      `json:"rank" db:"rank"`
	ServingLimit    int64      `json:"serving_limit" db:"serving_limit" kun:"descr=服务上限"`
	ReadFeatureID   int64      `json:"read_feature_id" db:"read_feature_id"`
	CreatedOn       *time.Time `json:"created_on" db:"created_on" kun:"descr=创建时间"`
	LastUpdated     *time.Time `json:"-" db:"last_updated"`
	LastActiveTime  *time.Time `json:"-" db:"last_active_time"`
	LastLoginTime   *time.Time `json:"-" db:"last_login_time"`
	DeletedAt       *time.Time `json:"deleted_at" db:"deleted_at" kun:"descr=删除时间"`
	ServingType     int64      `json:"serving_type" db:"serving_type" mapstructure:"serving_type"`
	ServingTypes    []string   `json:"serving_types" db:"-" kun:"descr=账号类型"`
	IsOnline        bool       `json:"is_online" db:"-"`
}

func (agentM *agentModel) toAgent() *Agent {
	return &Agent{
		ID:              agentM.ID,
		EntID:           agentM.EntID,
		GroupID:         agentM.GroupID,
		Avatar:          agentM.Avatar.String,
		Email:           agentM.Email.String,
		Nickname:        agentM.Nickname.String,
		Realname:        agentM.Realname.String,
		Password:        agentM.Password.String,
		Cellphone:       agentM.Cellphone.String,
		Status:          agentM.Status.String,
		Privilege:       agentM.Privilege.String,
		PublicCellphone: agentM.PublicCellphone.String,
		PublicEmail:     agentM.PublicEmail.String,
		QQ:              agentM.QQ.String,
		Signature:       agentM.Signature.String,
		Telephone:       agentM.Telephone.String,
		Weixin:          agentM.Weixin.String,
		Audience:        agentM.Audience.String,
		Token:           agentM.Token.String,
		WorkNum:         agentM.WorkNum.String,
		EmailActivated:  agentM.EmailActivated.Int64,
		Rank:            agentM.Rank,
		ServingLimit:    agentM.ServingLimit,
		ReadFeatureID:   agentM.ReadFeatureID.Int64,
		CreatedOn:       agentM.CreatedOn,
		LastUpdated:     agentM.LastUpdated,
		LastActiveTime:  agentM.LastActiveTime,
		LastLoginTime:   agentM.LastLoginTime,
		DeletedAt:       agentM.DeletedAt,
		ServingType:     agentM.ServingType.Int64,
		// ServingTypes:    models.ConvertServingType(agentM.ServingType.Int64),
	}
}

func demo16() {
	// go-sql-mysql 1.7.1 和 1.8.1(有问题)
	columns := []string{"*"}
	var agentM agentModel
	err := sess.Select(columns...).
		From("agent_info").
		Where(dbr.Eq("id", 1)).
		LoadOneContext(ctx, &agentM)

	if err != nil {
		log.Fatalf("got err: %v\n", err)
	}

	agent := agentM.toAgent()
	log.Printf("agent: %v\n", agent)
}

func asString(src any) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func demo17() {
	var f1 float64 = 999999
	var f2 float64 = 1000000
	var f3 []uint8 = []uint8{51, 57, 52, 53, 48, 48, 48}
	fmt.Println(asString(f1))
	fmt.Println(asString(f2))
	fmt.Println(asString(f3))
}

func demo18() {
	// COALESCE(MAX(`rank`),0)
	// IFNULL(SUM(`rank`), 0)
	var maxRank float64
	err := sess.Select("IFNULL(SUM(`rank`), 0)").
		From("selecting_rule").
		Where("enterprise_id = ?", 99999).
		LoadOneContext(ctx, &maxRank)
	if err != nil {
		log.Fatalf("got err: %v\n", err)
	}
	log.Println("got", maxRank)
}

func Main() {
	demo18()
}
