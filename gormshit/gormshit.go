package gormshit

import (
	"context"
	gosql "database/sql"
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

	// sql "github.com/go-sql-driver/mysql"
	"github.com/ClarkThan/labgo/gormshit/model"
	"github.com/bwmarrin/snowflake"
	sql "github.com/go-sql-driver/mysql"
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

	slaveDSN := "meiqia:f_xByc=9Dy+ZCbH1@tcp(pc-8vbvpi114t895m715.mysql.polardb.zhangbei.rds.aliyuncs.com:3306)/meiqia?charset=utf8&parseTime=True&loc=Local"
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
	dsn := "meiqia:f_xByc=9Dy+ZCbH1@tcp(pc-8vbvpi114t895m715.mysql.polardb.zhangbei.rds.aliyuncs.com:3306)/meiqia?charset=utf8&parseTime=True&loc=Local"
	//连接MYSQL
	gormDB, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	db = gormDB

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

func Main() {
	demo18()
}
