package models

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	//uuid
	UUid uuid.UUID `json:"uUid"`
	//用户名
	UserName string `json:"userName" binding:"required"`
	//密码
	UserPass string `json:"userPass"`
	//是否允许密码登录
	IsAllowed bool `json:"isAllowed" gorm:"default:false"`
	//公钥
	PublicKey []byte `json:"publicKey"`
	//邮箱
	UserEmail string `json:"userEmail"`
	//google code
	GoogleCode string `json:"googleCode"`
	//域
	NetId uint `json:"netId" binding:"required"`
	//子网
	NetIp uint32 `json:"netIp" gorm:"default:0"`
	//用户等级
	UserLevel uint8 `json:"userLevel" gorm:"default:0"`

	gorm.Model
}

// Login 登录验证
func (u *User) Login(user *User) (flag uint) {
	var userList User
	DB.Where("user_name = ?", user.UserName).Find(&userList)
	if userList.UserPass == user.UserPass {
		logrus.Println("登录成功:", user.UserName)
		flag = 1
		return
	} else {
		logrus.Println("登录失败")
		flag = 0
		return
	}
}

// CreateUser 创建用户
func (u *User) CreateUser() (flag uint) {
	DB.Create(u) // 通过数据的指针来创建
	flag = u.ID
	logrus.Println("新增成功:", flag)
	return
}

// FindById 通过id查找
func (u *User) FindById(id uint) (err error) {
	tx := DB.First(u, id)
	err = tx.Error
	return err
}

// FindByName 通过name查找
func (u *User) FindByName(name string) (id uint, err error) {
	tx := DB.Where("user_name = ?", name).Find(u)
	err = tx.Error
	id = u.ID
	return id, err
}

// FindByEmail 通过email查找
func (u *User) FindByEmail(email string) (user User, err error) {
	tx := DB.Where("user_email = ?", email).Find(u)
	err = tx.Error
	return *u, err
}

// UserList ui界面结构体
type UserList struct {
	Id        uint   `json:"id"`
	UserName  string `json:"userName"`
	UserPass  string `json:"userPass"`
	UserEmail string `json:"userEmail"`
	NetId     uint   `json:"netId" `
	NetIp     uint32 `json:"netIp" `
	IsAllowed bool   `json:"isAllowed"`
	UserLevel uint8  `json:"userLevel"`
}

// FindUi ui界面展示的信息
func (u *User) FindUi() (users []UserList) {
	DB.Model(u).Find(&users)
	return
}

// UpdateById 更新
func (u *User) UpdateById() {
	DB.Model(u).Updates(map[string]interface{}{"UserName": u.UserName,
		"IsAllowed": u.IsAllowed, "UserEmail": u.UserEmail,
		"NetId": u.NetId, "NetIp": u.NetIp, "UserLevel": u.UserLevel})
}

// DeleteById 删除
func (u *User) DeleteById(id string) {
	DB.Delete(u, id)
}
