package backend

import (
	"fmt"
	"sync"

	"github.com/WinnieVenice/packer/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB = nil
	once sync.Once
)

func init() {
	once.Do(func() {
		const dsn = "root:123456@tcp(127.0.0.1:3306)/cq?charset=utf8mb4&parseTime=True&loc=Local"
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Printf("db conn failed, dsn = (%+v), err = (%+v)\n", dsn, err)
			panic(err)
		}
		fmt.Println("conn mysql succ")
		err = db.AutoMigrate(&model.UserInfo{}, &model.PlatformUserInfo{})
		if err != nil {
			panic(err)
		}

	})
}

func MustInit() {}

func GetDB() *gorm.DB {
	return db
}
