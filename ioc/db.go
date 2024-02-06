package ioc

import (
	"github/yyfzy/mybook/config"
	"github/yyfzy/mybook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	//db, err := gorm.Open(mysql.Open("root:root@tcp(webook-mysql:13309)/webook"))

	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
