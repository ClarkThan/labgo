package gormshit

import (
	"errors"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	// sql "github.com/go-sql-driver/mysql"
	sql "github.com/go-sql-driver/mysql"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func initRWDB() {
	dsn := "test:12345687@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"
	//连接MYSQL
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "", SingularTable: true},
		Logger:         logger.Default.LogMode(logger.Silent),
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

func init() {
	dsn := "test:12345687@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"
	//连接MYSQL
	gormDB, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	db = gormDB
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
		log.Printf("demo7: %v\n", err)
		return
	}

	log.Printf("success query robot_ids: %v\n", ents)
}

func Main() {
	demo8()
}
