package gormshit

import (
	"context"
	gosql "database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	"github.com/ClarkThan/labgo/gormshit/model"
	"github.com/bwmarrin/snowflake"
	sql "github.com/go-sql-driver/mysql"
)

const (
	ABC = "hello"
)

const (
	ServingTypeConv   int8 = 1 << iota // 对话
	ServingTypeCall                    // 呼叫
	ServingTypeTicket                  // 工单
)

var (
	db            *gorm.DB
	snowflakeNode *snowflake.Node
)

func initRWDB() {
	dsn := "test:12345687@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"
	//连接MYSQL
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "", SingularTable: true},
		// Logger:         logger.Default.LogMode(logger.Silent),
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Silent,
				Colorful:      false,
			},
		),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	slaveDSN := "meiqia:JzpaqsFKtIacA!V@tcp(pc-8vbvpi114t895m715.mysql.polardb.zhangbei.rds.aliyuncs.com:3306)/meiqia?charset=utf8&parseTime=True&loc=Local"
	slaveDialector := mysql.New(mysql.Config{DSN: slaveDSN})

	resolver := dbresolver.Register(dbresolver.Config{
		// Sources:  []gorm.Dialector{},
		Replicas: []gorm.Dialector{slaveDialector},
		// Policy:   dbresolver.RandomPolicy{},
	}).SetMaxOpenConns(8).SetMaxIdleConns(10).SetConnMaxLifetime(time.Hour)

	err = gormDB.Use(resolver)
	if err != nil {
		log.Fatalf("配置读写分离失败: %v", err)
	}

	db = gormDB
}

func initSnowflake() {
	hostname, _ := os.Hostname() // hikari-new-89754945-62lv9

	var nodeID int64

	digitStr := strings.Join(regexp.MustCompile(`\d`).FindAllString(hostname, 8), "")

	if num, err := strconv.ParseInt(digitStr, 10, 64); err != nil && num < 100 {
		nodeID = rand.Int63n(1021)
	} else {
		nodeID = num % 1021
	}

	fmt.Println("--- snowflake id ", nodeID)
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		panic(err)
	}

	snowflakeNode = node
}

func init() {
	// dsn := "test:12345687@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"
	dsn := "meiqia:JzpaqsFKtIacA!V@tcp(pc-8vbvpi114t895m715.mysql.polardb.zhangbei.rds.aliyuncs.com:3306)/meiqia?charset=utf8&parseTime=True&loc=Local"
	//连接MYSQL
	gormDB, err := gorm.Open(mysql.Open(dsn)) // , &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	db = gormDB
	// initRWDB()

	initSnowflake()
}

type Membership struct {
	RobotID   int64     `gorm:"column:robot_id"`
	EntID     int64     `gorm:"column:ent_id"`
	SubType   string    `gorm:"column:sub_type"`
	Provider  string    `gorm:"column:provider"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (b *Membership) TableName() string {
	return "robot_membership"
}

func demo1() {
	b := new(Membership)
	fields := []string{"robot_id", "sub_type", "provider", "updated_at"}
	err := db.Clauses(dbresolver.Write).Table(b.TableName()).Select(fields).
		Where("robot_id = ? AND ent_id = ?", 63, 10).First(b).Error
	if err != nil {
		log.Println("damn!")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("X Not Found!")
			return
		}
		log.Printf("error: %v\n", err)
	}

	// log.Printf("%#v\n", b)
	log.Println(b.SubType, b.Provider, b.UpdatedAt)

	m := Membership{RobotID: 63, EntID: 10, SubType: "helplook"}
	// result := db.Create(&m)
	result := db.Select("robot_id", "ent_id", "sub_type").Create(&m)
	if err := result.Error; err != nil {
		var e *sql.MySQLError
		if errors.As(err, &e) && e.Number == 1062 {
			if err := db.Model(&Membership{}).Where("robot_id = ?", 12).Update("provider", "fucking").Error; err != nil {
				log.Fatalf("update error: %v\n", err)
			}
			log.Println("update ok")
			return
		}
		log.Printf("create error: %v\n", result.Error)
	} else {
		log.Println("create ok:", result.RowsAffected)
	}
}

func demo2() {
	var botIDs []int64
	err := db.Raw("select robot_id from robot_membership").Scan(&botIDs).Error
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	log.Printf("success query robot_ids: %v\n", botIDs)
}

func demo3() {
	var entCnts struct {
		EntID int64 `gorm:"column:ent_id"`
		Cnt   int64 `gorm:"column:cnt"`
	}

	err := db.Raw("select ent_id, count(1) cnt from robot_membership where ent_id in ? group by ent_id order by cnt limit 1", []int64{10, 22, 80}).Scan(&entCnts).Error
	if err != nil {
		log.Printf("demo3: %v\n", err)
		return
	}

	log.Printf("success query robot_ids: %v\n", entCnts)
}

func demo4() {
	var robotID int64
	err := db.Model(&Membership{}).Where("ent_id in ? limit 1", []int64{10, 22, 80}).Pluck("robot_id", &robotID).Error
	if err != nil {
		log.Printf("demo4: %v\n", err)
		return
	}

	log.Printf("success query robot_ids: %v\n", robotID)
}

func demo5() {
	var robotID int64
	err := db.Raw("select robot_id from robot_membership where ent_id in ? limit 1", []int64{10, 22, 80}).First(&robotID).Error
	if err != nil {
		log.Printf("demo5: %v\n", err)
		return
	}

	log.Printf("success query robot_ids: %v\n", robotID)
}

func demo6() {
	var exists bool
	// err := db.Raw("select exists (select 1 from robot_membership where ent_id in ?)", []int64{100, 220, 800}).First(&exists).Error
	err := db.Model(&Membership{}).Select("count(*) > 0").Where("ent_id in ?", []int64{100, 220, 80}).First(&exists).Error
	if err != nil {
		log.Printf("demo6: %v\n", err)
		return
	}

	log.Printf("success query robot_ids: %v\n", exists)
}

func demo7() {

	var bots []Membership
	err := db.Select("ent_id", "robot_id").Where("ent_id > ?", 10).Find(&bots).Order("ent_id").Error
	// err := db.Raw("select ent_id, robot_id from robot_membership where ent_id in ?", []int64{100, 220, 800}).Order("ent_id").Scan(&bots).Error
	if err != nil {
		log.Printf("demo7: %v\n", err)
		return
	}

	log.Printf("success query robot_ids: %v\n", bots)
}

func demo8() {

	var ents []struct {
		EntID   int64 `gorm:"column:ent_id"`
		RobotID int64 `gorm:"column:robot_id"`
	}
	err := db.Raw("select a.ent_id, a.robot_id from robot_membership a inner join gptbot_membership b on b.robot_id = a.ent_id").Order("ent_id").Scan(&ents).Error
	if err != nil {
		log.Printf("demo8: %v\n", err)
		return
	}

	log.Printf("success query robot_ids: %v\n", ents)
}

func demo9() {
	var grp string
	err := db.Raw("SELECT group_id FROM agent WHERE id = ?", 1).First(&grp).Error
	if err != nil {
		log.Printf("demo9: %v\n", err)
		return
	}

	log.Printf("got group_id str :%s\n", grp)
}

type Agent struct {
	EntID   int64  `gorm:"column:ent_id"`
	GroupID string `gorm:"column:group_id"`
	Name    string `gorm:"column:name"`
	Email   string `gorm:"column:email"`
}

func (*Agent) TableName() string {
	return "agent"
}

func demo10() {
	// var a Agent
	// err := db.Raw("SELECT ent_id, group_id FROM agent WHERE id = ?", 5).First(&a).Error
	// if err != nil {
	// 	log.Printf("demo10: %v\n", err)
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		log.Println("got it ....")
	// 	}
	// 	return
	// }

	a := new(Agent)
	// err := db.Table(a.TableName()).Select("ent_id", "group_id").Where("name = ?", "ranyax").First(a).Error
	err := db.Model(new(Agent)).Select("ent_id", "group_id").Where("name = ?", "ranya").First(a).Error
	// err := db.Model(new(Agent)).Where("name = ?", "ranyax").First(&a).Error
	if err != nil {
		log.Printf("demo10: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("got it ....")
		}
		return
	}

	log.Printf("got ent_id: %d   group_id str :%s   %t\n", a.EntID, a.GroupID, a.GroupID == "23")
}

func demo11() {
	var id int64
	err := db.Raw("select id from agent where ent_id = ? order by name limit 1", 10).Scan(&id).Error
	if err != nil {
		log.Printf("demo11: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("got it ....")
		}
		return
	}

	log.Println("got ", id)
}

type ConvActive struct {
	ID                    int64      `gorm:"column:id"`
	EntID                 int64      `gorm:"column:enterprise_id"`
	TrackID               string     `gorm:"column:track_id"`
	VisitID               string     `gorm:"column:visit_id"`
	AgentID               int64      `gorm:"column:agent_id"`
	AgentType             string     `gorm:"column:agent_type"`
	Assignee              int64      `gorm:"column:assignee"`
	CreatedOn             *time.Time `gorm:"column:created_on"`
	EndedOn               *time.Time `gorm:"column:ended_on"`
	EndedBy               string     `gorm:"column:ended_by"`
	MsgNum                int        `gorm:"column:msg_num"`
	ClientMsgNum          int        `gorm:"column:client_msg_num"`
	AgentMsgNum           int        `gorm:"column:agent_msg_num"`
	AgentEffectiveMsgNum  int        `gorm:"column:agent_effective_msg_num"`
	FirstResponseWaitTime int        `gorm:"column:first_response_wait_time"`
	ClientFirstSendTime   *time.Time `gorm:"column:client_first_send_time"`
	FirstMsgCreatedOn     *time.Time `gorm:"column:first_msg_created_on"`
	LastMsgCreatedOn      *time.Time `gorm:"column:last_msg_created_on"`
	LastMsgContent        string     `gorm:"column:last_msg_content"`
	ConvDuration          int        `gorm:"column:conversation_duration"`
	QualityGrade          string     `gorm:"column:quality_grade"`
	Clues                 string     `gorm:"column:clues"`
	URL                   string     `gorm:"column:url"`
	Title                 string     `gorm:"column:title"`
	LastUpdated           *time.Time `gorm:"column:last_updated"`
	Source                string     `gorm:"column:source"`
	SubSource             string     `gorm:"column:sub_source"`
}

func (*ConvActive) TableName() string {
	return "conversation_active"
}

type AgentConvCnt struct {
	AgentID int64 `gorm:"column:agent_id"`
	ConvCnt int64 `gorm:"column:conv_cnt"`
}

func demo12() {
	entID := 10
	agentIDs := []int64{65}
	var convCnts []*AgentConvCnt
	// err := r.meiqiaDB.Raw("SELECT agent_id, count(1) conv_cnt FROM conversation_active WHERE enterprise_id = ? AND agent_id in ? AND ended_on is NULL GROUP BY agent_id", entID, agentIDs).Scan(&convCnts).Error
	err := db.Select("agent_id", "count(1) conv_cnt").Model(&ConvActive{}).
		Where("enterprise_id = ? AND agent_id in ? AND ended_on is NULL", entID, agentIDs).
		Group("agent_id").Scan(&convCnts).Error

	if err != nil {
		log.Printf("get agent conv cnt %v\n", err)
		return
	}

	log.Printf("ret: %d -> %d\n", convCnts[0].AgentID, convCnts[0].ConvCnt)
}

type AgentInfo struct {
	ID          int64  `gorm:"column:id;primary_key;"`
	EntID       int64  `gorm:"column:enterprise_id"`
	Email       string `gorm:"column:email"`
	Nickname    string `gorm:"column:nickname"`
	Realname    string `gorm:"column:realname"`
	Avatar      string `gorm:"column:avatar"`
	GroupID     int64  `gorm:"column:group_id"`
	ServingType int64  `gorm:"column:serving_type"`
	Privilege   string `gorm:"column:privilege"`
}

func (*AgentInfo) TableName() string {
	return "agent"
}

func demo13() {
	entID := 12
	agentIDs := []int64{392, 391, 378}

	var agents []*AgentInfo
	err := db.Model(&AgentInfo{}).Select("id", "serving_limit", "rank", "group_id").
		Where("enterprise_id = ? AND id in ? AND status = ?", entID, agentIDs, "on_duty").Scan(&agents).Error
	// builder := r.meiqiaDB.Raw("SELECT id, serving_limit, `rank`, group_id FROM agent WHERE enterprise_id = ? AND id in ? AND status = 'on_duty'", entID, agentIDs)

	if err != nil {
		log.Printf("get on_duty agents %v", err)
		return
	}

	log.Printf("ret %v\n", agents[0])
}

type AdvanceSelectingRule struct {
	ID              int64  `gorm:"column:id"`
	Inverted        int8   `gorm:"column:inverted"`
	MatchType       string `gorm:"column:match_type"`
	AllocOption     string `gorm:"column:alloc_option"`
	MatchRules      []byte `gorm:"column:match_rules"`
	Targets         []byte `gorm:"column:targets"`
	OverflowTargets []byte `gorm:"column:overflow_targets"`
}

func demo14() {
	var rules []*AdvanceSelectingRule
	entID := 12
	// agentIDs := []int64{392, 391, 378}
	err := db.Raw("SELECT id, match_type, targets, match_rules, inverted, overflow_targets, alloc_option FROM selecting_rule WHERE enterprise_id = ? AND enabled = 1 ORDER BY `rank`", entID).Scan(&rules).Error
	log.Println("<------")
	if err != nil {
		log.Printf("get on_duty agents %v", err)
		return
	}

	log.Println("------>")
	if len(rules) > 0 {
		log.Printf("ret[0] = %v\n", rules[0])
	}
}

func demo15() {
	var exists bool
	err := db.Model(&AgentInfo{}).Select("count(*) > 0").Where("enterprise_id = ?", 9999999).First(&exists).Error
	if err != nil {
		log.Printf("xxxx  %v", err)
		return
	}
	log.Println(exists)
}

func stringify(val any) string {
	switch value := val.(type) {
	case string:
		return value
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", val)
	}
}

func demo16Aux(ctx context.Context, body, params, ext map[string]any, saveAsActive bool) map[string]any {
	evtID := snowflakeNode.Generate().Int64()
	agentID, _ := params["agent"].(int64)
	var actionCode int64 = 17
	var entID int64 = 63790

	targetIDStr := stringify(params["target_id"]) // target_id 可能是 agent_id, track_id, 空
	targetKind, _ := params["target_kind"].(string)
	trackID, _ := params["track_id"].(string) // 可能为空
	createdOn := time.Now().UTC()
	bodyDat, _ := json.Marshal(body)

	dbEvt := model.Event{
		ID:            evtID,
		Action:        actionCode,
		EntID:         entID,
		TargetID:      targetIDStr,
		TargetKind:    targetKind,
		Body:          bodyDat,
		TrackID:       trackID,
		AgentID:       agentID,
		AgentNickname: "",
		CreatedOn:     createdOn,
	}

	if saveAsActive {
		dbActiveEvt := model.EventActive(dbEvt)
		if result := db.Create(&dbActiveEvt); result.RowsAffected != 1 {
			log.Fatalf("insert event_active fail: %v", result.Error)
		}
	} else {
		if result := db.Create(&dbEvt); result.RowsAffected != 1 {
			log.Fatalf("insert event fail: %v", result.Error)
		}
	}

	evtIDStr := fmt.Sprintf("%020d", evtID)

	evt := map[string]any{
		"id":             evtIDStr,
		"enterprise_id":  entID,
		"action":         "init_conv",
		"body":           body,
		"agent_id":       agentID,
		"agent_nickname": "",
		"realname":       "",
		"avatar":         "",
		"track_id":       trackID,
		"target_id":      targetIDStr,
		"target_kind":    targetKind,
		"created_on":     createdOn,
	}

	for k, v := range ext {
		evt[k] = v
	}

	return evt
}

func demo16() {
	body := map[string]any{"msg_id": 110, "client_id": "fuckyoushit", "conv_id": 110110, "agent_id": 119, "msg_created_time": time.Now().UTC()}
	params := map[string]any{"agent": 119, "target_kind": "conv", "target_id": 110110, "track_id": "fuckyoushit"}
	ext := map[string]any{"source": "baidu_bcp"}
	evt := demo16Aux(context.Background(), body, params, ext, false)
	dat, _ := json.Marshal(evt)
	log.Println("event: ", string(dat))
}

func demo17() {
	fields := []string{"id"}
	conv := new(model.Conversation)
	err := db.Table(conv.TableName()).Select(fields).Where("id = ? AND enterprise_id = ? AND ended_on is NULL", 14162, 10).First(conv).Error
	if err != nil {
		log.Printf("no active conv: %v\n", err)
		return
	}

	log.Printf("got conv %d\n", conv.ID)
}

type WxFansInfoModel struct {
	ID            int64            `gorm:"column:id" json:"id"`
	EntID         int64            `gorm:"column:ent_id" json:"ent_id"` //  enterprise id
	Fansopenid    string           `gorm:"column:fansopenid" json:"fansopenid"`
	Gzopenid      gosql.NullString `gorm:"column:gzopenid" json:"gzopenid"`
	Createdtime   gosql.NullInt64  `gorm:"column:createdtime" json:"createdtime"`
	TrackID       gosql.NullString `gorm:"column:track_id" json:"track_id"`
	VisitID       gosql.NullString `gorm:"column:visit_id" json:"visit_id"`
	Subscribe     gosql.NullString `gorm:"column:subscribe" json:"subscribe"`
	Nickname      gosql.NullString `gorm:"column:nickname" json:"nickname"`
	Sex           gosql.NullString `gorm:"column:sex" json:"sex"`
	City          gosql.NullString `gorm:"column:city" json:"city"`
	Country       gosql.NullString `gorm:"column:country" json:"country"`
	Language      gosql.NullString `gorm:"column:language" json:"language"`
	Wxfansinfocol gosql.NullString `gorm:"column:wxfansinfocol" json:"wxfansinfocol"`
	SubscribeTime gosql.NullInt64  `gorm:"column:subscribe_time" json:"subscribe_time"`
	Unionid       gosql.NullString `gorm:"column:unionid" json:"unionid"`
	Remark        gosql.NullString `gorm:"column:remark" json:"remark"`
	Groupid       gosql.NullString `gorm:"column:groupid" json:"groupid"`
	Nextonlinetmp gosql.NullInt64  `gorm:"column:nextonlinetmp" json:"nextonlinetmp"`
	Province      gosql.NullString `gorm:"column:province" json:"province"`
	FollowSource  string           `gorm:"column:follow_source" json:"follow_source"` //  visitor srouce
}

func (WxFansInfoModel) TableName() string {
	return "wx_fans_info"
}

func demo18() {
	var res WxFansInfoModel
	var entID int64 = 7
	var trackID string = "2M5Kv8STIDDtRlus2YZ9nbLvfN4"
	err := db.First(&res, `ent_id = ? and track_id = ?`, entID, trackID).Error
	if err != nil {
		log.Fatalf("fail: %v\n", err)
	}
	log.Printf("got: %+v\n", res)
}

type ConvEvaluation struct {
	ID        int64           `gorm:"column:id;primary_key"`
	ConvID    int64           `gorm:"column:conversation_id"`
	Level     int8            `gorm:"column:level"`
	Resolved  gosql.NullInt16 `gorm:"column:resolved"`
	Content   string          `gorm:"column:content"`
	CreatedOn *time.Time      `gorm:"column:created_on"`
	UpdatedOn *time.Time      `gorm:"column:updated_on"`
}

func (*ConvEvaluation) TableName() string {
	return "evaluation"
}

func demo19() {
	evaluation := ConvEvaluation{
		Level: -1, // 这个默认值不能为0
		// Resolved: -1, // 这个默认值不能为0
	}
	err := db.Table(evaluation.TableName()).Select([]string{"level", "resolved"}).
		Where(`conversation_id = ?`, 16839).First(&evaluation).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("got err: %v\n", err)
	}

	log.Printf("evaluation: %+v\n", evaluation)
	if evaluation.Resolved.Valid {
		log.Println("kkkkk", evaluation.Resolved.Int16)
	}
	// log.Println(evaluation.Level, evaluation.Level == 2, evaluation.Resolved == -1)
	// log.Println(binaryMatch(evaluation.Level, 2))
	// log.Println(binaryMatch(evaluation.Level, 1))
	// log.Println(binaryMatch(evaluation.Level, 0))
	// log.Println(binaryMatch(evaluation.Resolved, 1))
	// log.Println(binaryMatch(evaluation.Resolved, 0))
	var m map[string]any
	log.Println(parseSignedInt[int8](m["resolved"]))
}

func parseSignedInt[T ~int | ~int8 | ~int64](val any) (t T) {
	return parseInteger[T](val, -1)
}

func parseInteger[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](val any, defaultVal ...T) (t T) {
	var defVal T = 0
	if len(defaultVal) > 0 {
		defVal = defaultVal[0]
	}
	if val == nil {
		return defVal
	}

	switch v := val.(type) {
	case int:
		return T(v)
	case int8:
		return T(v)
	case int16:
		return T(v)
	case int32:
		return T(v)
	case int64:
		return T(v)
	case uint:
		return T(v)
	case uint8:
		return T(v)
	case uint16:
		return T(v)
	case uint32:
		return T(v)
	case uint64:
		return T(v)
	case float64:
		return T(v)
	default:
		return defVal
	}
}

func binaryMatch(level, lvlNum int8) int8 {
	if level == lvlNum {
		return 1
	}

	return 0
}

func GetHumanAgentIDs(entID int64, serveType ...int8) ([]int64, error) {
	var servingType int8 = ServingTypeConv
	if len(serveType) > 0 {
		servingType = serveType[0]
	}

	var servingTypes []int8
	// combinationCnt := int(math.Pow(2, 3)) - 1
	for i := 1; i <= 7; i++ {
		it := int8(i)
		if it&servingType == servingType {
			servingTypes = append(servingTypes, it)
		}
	}

	log.Println("serving_type = ", servingTypes)
	var agentIDs []int64
	err := db.Raw("SELECT id FROM agent WHERE enterprise_id = ? AND deleted_at IS NULL AND serving_type in (?) AND privilege not like 'bot%'", entID, servingTypes).Scan(&agentIDs).Error
	if err != nil {
		log.Printf("query human agent ids err: %v\n", err)
		return nil, err
	}

	return agentIDs, nil
}

func demo20() {
	aids, _ := GetHumanAgentIDs(10)
	log.Println("agent_ids: ", len(aids), aids)
}

func demo21() {
	log.Println(ServingTypeConv, ServingTypeCall, ServingTypeTicket)
}

// type SmartGuideTemplate struct {
// 	TemplateID   string            `gorm:"column:template_id"`
// 	Target       string            `gorm:"column:target"`
// 	Rank         float64           `gorm:"column:rank"`
// 	ManualStatus string            `gorm:"column:manual_status"`
// 	Settings     SmartGuideSetting `gorm:"column:settings"`
// 	Details      SmartGuideDetail  `gorm:"column:details"`
// 	CreatedAt    *time.Time        `gorm:"column:created_at"`
// }

// func (t *SmartGuideTemplate) TableName() string {
// 	return "enterprise_mpush_template"
// }

// type VisitorRule struct {
// 	Name string `json:"name"`
// 	Val  string `json:"val"`
// 	Op   int8   `json:"op"`
// }

// // {"auto_status":"open","delay":4,"push_begin_time":"","push_end_time":"","push_time":"offline","push_visitor_type":1,"visitor_rules":[],"work_rule":1}
// type SmartGuideSetting struct {
// 	AutoStatus      string        `json:"auto_status"`
// 	Delay           int64         `json:"delay"`
// 	PushBeginTime   string        `json:"push_begin_time"`
// 	PushEndTime     string        `json:"push_end_time"`
// 	PushTime        string        `json:"push_time"`
// 	VisitorRules    []VisitorRule `json:"visitor_rules"`
// 	PushVisitorType int8          `json:"push_visitor_type"`
// 	WorkRule        int8          `json:"work_rule"`
// }

// func (s SmartGuideSetting) Value() (driver.Value, error) {
// 	return json.Marshal(s)
// }

// func (s *SmartGuideSetting) Scan(v any) error {
// 	switch data := v.(type) {
// 	case []byte:
// 		return json.Unmarshal(data, &s)
// 	case string:
// 		return json.Unmarshal([]byte(data), &s)
// 	default:
// 		return fmt.Errorf("invalid setting field data: %v", v)
// 	}
// }

// type innerType struct {
// 	Type      string `json:"type"`
// 	Countdown int64  `json:"countdown,omitempty"`
// }

// type SmartGuideDetail struct {
// 	AutoAutoBounde   innerType `json:"auto_auto_bounce"`
// 	AutoStatus       string    `gorm:"column:auto_status"`
// 	Delay            int64     `gorm:"column:delay"`
// 	GapSecond        int64     `json:"gap_second"`
// 	GreetingsContent string    `json:"greetings_content"`
// 	GreetingsType    int8      `json:"greetings_type"`
// 	ManualAutoBounce innerType `json:"manual_auto_bounce"`
// 	ManualGapSecond  int64     `json:"manual_gap_second"`
// 	ManualPushRate   int8      `json:"manual_push_rate"`
// 	PushRate         int8      `json:"push_rate"`
// 	PushTime         string    `json:"push_time"`
// 	PushVisitorType  int8      `json:"push_visitor_type"`
// 	SenderID         int64     `json:"sender_id"`
// }

// func (s SmartGuideDetail) Value() (driver.Value, error) {
// 	return json.Marshal(s)
// }

// func (s *SmartGuideDetail) Scan(v any) error {
// 	switch data := v.(type) {
// 	case []byte:
// 		return json.Unmarshal(data, &s)
// 	case string:
// 		return json.Unmarshal([]byte(data), &s)
// 	default:
// 		return fmt.Errorf("invalid detail field data: %v", v)
// 	}
// }

/*
{
    "auto_auto_bounce": {
        "type": "stop"
    },
    "auto_status": "open",
    "delay": 4,
    "gap_second": 0,
    "greetings_content": "\u003cp\u003e\u003c/p\u003e",
    "greetings_content_draft": "{\"blocks\":[{\"key\":\"43km3\",\"text\":\"\",\"type\":\"unstyled\",\"depth\":0,\"inlineStyleRanges\":[],\"entityRanges\":[],\"data\":{}}],\"entityMap\":{}}",
    "greetings_style": {
        "actions": [
            {
                "height": 60,
                "id": "lCI31F",
                "position": {
                    "bottom": "auto",
                    "left": 52,
                    "right": "auto",
                    "top": 78
                },
                "type": 1,
                "width": 80
            }
        ],
        "bgi": {
            "height": 200,
            "src": "https://meiqia-upload-qa.meiqiausercontent.com/widget/10/ysfc/FZbxH8QsNsMVlj6hhHNK.png",
            "width": 200
        }
    },

}

*/

type SmartGuideTemplate struct {
	TemplateID string             `gorm:"column:template_id" json:"template_id"`
	Target     string             `gorm:"column:target" json:"target"`
	Settings   *SmartGuideSetting `gorm:"column:settings" json:"settings"`
	Details    *SmartGuideDetail  `gorm:"column:details" json:"details"`
	// Rank         float64            `gorm:"column:rank" json:"rank"`
	// ManualStatus string             `gorm:"column:manual_status" json:"manual_status"`
}

func (t *SmartGuideTemplate) TableName() string {
	return "enterprise_mpush_template"
}

// 这些字段只有在匹配的时候才有用，后续智能引导处理的过程中不需要，这里目的是保存redis是可以节省点空间
func (t *SmartGuideTemplate) Trim() *SmartGuideTemplate {
	t.Settings.PushBeginTime = ""
	t.Settings.PushEndTime = ""
	t.Settings.PushTime = ""
	t.Settings.VisitorRules = nil
	t.Settings.PushVisitorType = 0
	t.Settings.WorkRule = 0

	return t
}

type VisitorRule struct {
	Name string `json:"name"`
	Val  string `json:"val"`
	Op   int8   `json:"op"`
}

// {"auto_status":"open","delay":4,"push_begin_time":"","push_end_time":"","push_time":"offline","push_visitor_type":1,"visitor_rules":[],"work_rule":1}
type SmartGuideSetting struct {
	AutoStatus      string        `json:"auto_status"` // 自动推送开启状态
	Delay           int64         `json:"delay"`       // 访客访问网站delay秒后自动推送
	PushBeginTime   string        `json:"push_begin_time,omitempty"`
	PushEndTime     string        `json:"push_end_time,omitempty"`
	PushTime        string        `json:"push_time,omitempty"` // 推送时间范围: all | online | offline | custom
	VisitorRules    []VisitorRule `json:"visitor_rules,omitempty"`
	PushVisitorType int8          `json:"push_visitor_type,omitempty"` // 推送访客类型: 1:全部 | 2:无线索访客 | 3:自定义规则
	WorkRule        int8          `json:"work_rule,omitempty"`
}

func (s SmartGuideSetting) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SmartGuideSetting) Scan(v any) error {
	switch data := v.(type) {
	case []byte:
		return json.Unmarshal(data, &s)
	case string:
		return json.Unmarshal([]byte(data), &s)
	default:
		return fmt.Errorf("invalid setting field data: %v", v)
	}
}

type innerType struct {
	Type      string `json:"type"`
	Countdown int64  `json:"countdown,omitempty"`
}

type SmartGuideDetail struct {
	AutoAutoBounde   innerType `json:"auto_auto_bounce"`
	ManualAutoBounce innerType `json:"manual_auto_bounce"`
	Type             string    `json:"type"`
	GreetingsContent string    `json:"greetings_content"`
	SenderID         int64     `json:"sender_id"`
	GapSecond        int64     `json:"gap_second"`
	ManualGapSecond  int64     `json:"manual_gap_second"`
	GreetingsType    int8      `json:"greetings_type"`
	ManualPushRate   int8      `json:"manual_push_rate"`
	PushRate         int8      `json:"push_rate"`
	GreetingsStyle   any       `json:"greetings_style"` // TODO(any OK)
	// 下面4个字段以 settings 中的为准
	// AutoStatus       string    `gorm:"column:auto_status"`
	// Delay            int64     `gorm:"column:delay"`
	// PushTime         string    `json:"push_time"`
	// PushVisitorType int8  `json:"push_visitor_type"`
}

func (s SmartGuideDetail) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SmartGuideDetail) Scan(v any) error {
	switch data := v.(type) {
	case []byte:
		return json.Unmarshal(data, &s)
	case string:
		return json.Unmarshal([]byte(data), &s)
	default:
		return fmt.Errorf("invalid detail field data: %v", v)
	}
}

func demo22() {
	entID := 10

	var templates []SmartGuideTemplate
	// err := db.Raw("SELECT template_id, target, `rank`, manual_status, created_at, settings, details FROM enterprise_mpush_template WHERE ent_id = ? AND template_status = 0 ORDER BY `rank`, `created_at`", entID).Scan(&templates).Error
	// err := db.Model(&SmartGuideTemplate{}).Select("template_id", "target", "`rank`", "manual_status", "created_at", "settings", "details").
	err := db.Model(&SmartGuideTemplate{}).Select("template_id", "target", "settings", "details").
		Where("ent_id = ? AND template_status = 1", entID).
		Find(&templates).Error
	if err != nil {
		log.Printf("err: %v\n", err)
		return
	}

	for _, tpl := range templates {
		bs, _ := json.Marshal(tpl)
		fmt.Println(tpl.TemplateID, string(bs), "\n")
	}

	// log.Println(templates[0].Details.GreetingsStyle)
	// dat := map[string]any{"name": "foo", "data": templates[0].Details.GreetingsStyle}
	// bs, _ := json.Marshal(dat)
	// log.Println(string(bs))
	// log.Println(templates[1].TemplateID, templates[1].Settings.VisitorRules)
}

func demo33() {
	var entIDs []int64
	err := db.Raw("SELECT ent_id FROM feishu_account WHERE enabled = 2").Scan(&entIDs).Error
	if err != nil {
		panic(err)
	}

	set := make(map[int64]struct{}, len(entIDs))
	for _, entID := range entIDs {
		set[entID] = struct{}{}
	}

	log.Printf("set: %v\n", set)
}

type DevClientTrack struct {
	DevClientID string `gorm:"column:dev_client_id"`
	TrackID     string `gorm:"column:track_id"`
}

func demo34() {
	var exists bool
	// var exists int
	if err := db.Model(&Membership{}).Select("count(*) > 0").Where("robot_id = ?", 61).First(&exists).Error; err != nil {
		log.Fatalf("query robot_membership error: %v\n", err)
	}
	fmt.Println(exists)

	var entID int64 = 2
	trackIDs := []string{"2MJhQbmOVxn411vDL7ODq2d4pgQ", "2NH1waNRUcXQJxaN0ggColpRSxP", "2NPHs9vG1XFxPbd33gw1jhcRsDv"}
	var data []*DevClientTrack
	err := db.Table(`dev_client_track`).Select(`dev_client_id, track_id`).Where(`enterprise_id = ? AND track_id in ?`, entID, trackIDs).Scan(&data).Error
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	ret := make(map[string]string, len(data))
	for _, v := range data {
		ret[v.TrackID] = v.DevClientID
	}

	for k, v := range ret {
		fmt.Println(k, v)
	}
}

func Main() {
	demo34()
}
