package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ivs-net-server/auth/configure"
)

var DB *gorm.DB

// 数据库初始化连接
// 读取conf/app.toml配置文件
func init() {
	DB, _ = gorm.Open(sqlite.Open(configure.Config.Get("db.sqlite.path").(string)), &gorm.Config{})
	//更新数据库的字段，只增不减
	_ = DB.AutoMigrate(&User{}, &NetWork{})
}
