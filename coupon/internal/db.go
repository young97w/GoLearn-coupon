package internal

import (
	"coupon/coupon_srv/model"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"dbName"`
	UserName string `mapstructure:"userName"`
	Password string `mapstructure:"password"`
}

func InitDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	var err error
	//dsn := "username:password@tcp(127.0.0.1:3306)/orm_test?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConf.DBConfig.UserName,
		AppConf.DBConfig.Password,
		AppConf.DBConfig.Host,
		AppConf.DBConfig.Port,
		AppConf.DBConfig.DBName,
	)
	zap.S().Infof(dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("数据库连接失败" + err.Error())
	}

	err = DB.AutoMigrate(&model.Account{}, &model.Coffee{}, &model.Coupon{})
	if err != nil {
		panic(err)

	}
	fmt.Println("连接成功！")
}

func Paginate(pageSize, pageNo int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case pageSize >= 100:
			pageSize = 10
		case pageSize < 1:
			pageSize = 5
		}
		if pageNo < 1 {
			pageNo = 1
		}
		return db.Offset((pageNo - 1) * pageSize).Limit(pageSize)
	}
}
