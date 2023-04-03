package https

import (
	"crypto/sha512"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm4"
	"github.com/tjfoc/gmsm/x509"
	"io/ioutil"
	"ivs-net-server/auth/google"
	"ivs-net-server/auth/models"
	"ivs-net-server/auth/sm"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strconv"
	"time"
)

//user_controller

// LoginBody 管理端用户登录
type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

// Login 登录
func Login(c *gin.Context) {
	var login LoginBody
	var user models.User
	err := c.ShouldBindJSON(&login)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "error",
		})
	} else {
		login.Email = "testdemob@sina.com"
		user = FindByEmail(login.Email)

		login.Code, _ = google.NewGoogleAuth().GetCode(user.GoogleCode)
		login.Password = user.UserPass

		code, _ := google.NewGoogleAuth().GetCode(user.GoogleCode)
		if login.Code == code && login.Password == user.UserPass {
			c.JSON(200, gin.H{
				"message": "success",
			})
		} else {
			c.JSON(400, gin.H{
				"message": "error",
			})
		}
	}
}

// PostReqType 第一次请求结构体
type PostReqType struct {
	UserEmail string `json:"UserEmail"`
	RandNum   []byte `json:"RandNum"`
	TimeStamp int64  `json:"TimeStamp"`
}

// RequestBody 请求字节流
type RequestBody struct {
	Res []byte `json:"res"`
}

// RespondBody 第一次响应体
type RespondBody struct {
	Mes       string
	Uuid      uuid.UUID
	UserName  string
	UserEmail string
	RandNum   []byte
	TimeStamp int64
}

// 第一次得到的客户端512字节
var firstClient []byte

// 第一次服务端512字节
var firstAuth []byte
var userBody models.User

// 第一次请求的时间戳
var firstTime int64

// Login1 登录
func Login1(c *gin.Context) {
	var req RequestBody
	var reqBody PostReqType
	var respondBody RespondBody
	_ = c.ShouldBindJSON(&req) // 解析req参数
	//私钥解密
	readFile, _ := ioutil.ReadFile("././conf/privateKey.txt")
	pem, _ := x509.ReadPrivateKeyFromPem(readFile, nil)
	firstClient, _ = pem.DecryptAsn1(req.Res)
	_ = json.Unmarshal(firstClient, &reqBody)
	firstTime = reqBody.TimeStamp

	//返回参数
	userBody = FindByEmail(reqBody.UserEmail)
	nowTime := time.Now().Unix()
	if userBody.ID == 0 || firstTime+30 < nowTime || firstTime-30 > nowTime {
		respondBody.Mes = "error"
		c.JSON(http.StatusBadRequest, respondBody)
		return
	} else {
		respondBody.Mes = ""
		respondBody.Uuid = userBody.UUid
		respondBody.UserName = userBody.UserName
		respondBody.UserEmail = userBody.UserEmail
		respondBody.TimeStamp = time.Now().Unix()
		respondBody.RandNum = RandUp(488 - len(userBody.UserName) - len(userBody.UserEmail))
		//客户端公钥加密
		firstAuth, _ = json.Marshal(respondBody)
		publicKey, _ := x509.ReadPublicKeyFromPem(userBody.PublicKey)
		encryptAsn1, _ := publicKey.EncryptAsn1(firstAuth, nil)
		req.Res = encryptAsn1
		c.JSON(http.StatusOK, req)
	}
}

// PostReqType2 第二次请求结构体
type PostReqType2 struct {
	Rand504   []byte   `json:"Rand504"`
	TimeStamp int64    `json:"TimeStamp"`
	Sha1024   [64]byte `json:"Sha1024"`
}

// Sha sha512计算的结构体
type Sha struct {
	Rand504   []byte
	TimeStamp int64
	Password  string
}

// RespondBody2 第二次响应体
type RespondBody2 struct {
	Mes     string
	RandNum []byte
}

// 第二次得到的客户端512字节
var secondClient []byte

// 第二次服务端512字节
var secondAuth []byte

// Login2 登录
func Login2(c *gin.Context) {
	var req RequestBody
	var reqBody PostReqType2
	var respondBody RespondBody2
	_ = c.ShouldBindJSON(&req) // 解析req参数
	//私钥解密
	readFile, _ := ioutil.ReadFile("././conf/privateKey.txt")
	pem, _ := x509.ReadPrivateKeyFromPem(readFile, nil)
	secondClient, _ = pem.DecryptAsn1(req.Res)
	_ = json.Unmarshal(secondClient, &reqBody)
	sha := Sha{
		Rand504:   reqBody.Rand504,
		TimeStamp: reqBody.TimeStamp,
		Password:  userBody.UserPass,
	}
	marshal, _ := json.Marshal(sha)
	sum512 := sha512.Sum512(marshal)
	for i := 0; i < 1023; i++ {
		tempByte := make([]byte, 65)
		for k := 0; k < 64; k++ {
			tempByte[k] = sum512[k]
		}
		sum512 = sha512.Sum512(tempByte)
	}
	if reqBody.Sha1024 != sum512 {
		respondBody.Mes = "error"
		c.JSON(http.StatusBadRequest, respondBody)
		return
	} else {
		secondAuth = RandUp(512)
		// sm4加密
		code := "123456" //google.NewGoogleAuth().GetCode(userBody.GoogleCode)
		code += "0000000000"
		key := []byte(code)
		iv := make([]byte, sm4.BlockSize)
		ciphertext, _ := sm.Sm4Encrypt(key, iv, secondAuth)
		respondBody.RandNum = ciphertext
		req.Res, _ = json.Marshal(respondBody)
		c.JSON(http.StatusOK, req)
	}
}

// AllByte 2048字节
type AllByte struct {
	FirstClient  []byte
	SecondClient []byte
	FirstAuth    []byte
	SecondAuth   []byte
}

// Request 第三次请求字节流
type Request struct {
	Res [64]byte `json:"Res"`
}

// Rc4SecretKey rc4密钥结构
type Rc4SecretKey struct {
	Rc4Key    [64]byte
	TimeStamp int64
}

// Rc4Key 密钥
var Rc4Key [64]byte

// Login3 登录
func Login3(c *gin.Context) {
	allByte := AllByte{
		FirstClient:  firstClient,
		SecondClient: secondClient,
		FirstAuth:    firstAuth,
		SecondAuth:   secondAuth,
	}
	marshal, _ := json.Marshal(allByte)
	sum8192 := sha512.Sum512(marshal)
	for i := 0; i < 8191; i++ {
		tempByte := make([]byte, 65)
		for k := 0; k < 64; k++ {
			tempByte[k] = sum8192[k]
		}
		sum8192 = sha512.Sum512(tempByte)
	}
	var req Request
	_ = c.ShouldBindJSON(&req) // 解析req参数
	if sum8192 == req.Res {
		sum4096 := sha512.Sum512(marshal)
		for i := 0; i < 4095; i++ {
			tempByte := make([]byte, 65)
			for k := 0; k < 64; k++ {
				tempByte[k] = sum4096[k]
			}
			sum4096 = sha512.Sum512(tempByte)
		}
		rc4SecretKey := Rc4SecretKey{
			Rc4Key:    sum4096,
			TimeStamp: firstTime,
		}
		secretKey, _ := json.Marshal(rc4SecretKey)
		Rc4Key = sha512.Sum512(secretKey)
		for i := 0; i < 4095; i++ {
			tempByte := make([]byte, 65)
			for k := 0; k < 64; k++ {
				tempByte[k] = Rc4Key[k]
			}
			Rc4Key = sha512.Sum512(tempByte)
		}
		c.JSON(http.StatusOK, "OK")
	} else {
		c.JSON(http.StatusBadRequest, "ERROR")
	}
}

type DhcpBody struct {
	Mes     string
	Id      uint
	Ip      uint32
	Mask    uint32
	Gateway uint32
	Dns     uint32
}

var DhcpGet DhcpBody

type dhcpGetBody struct {
	Email string `json:"email"`
	Ip    uint32 `json:"ip"`
}

// AddressIp 可用地址列表
func AddressIp(c *gin.Context) {
	var dhcpBody dhcpGetBody
	_ = c.ShouldBindJSON(&dhcpBody)
	var mUser models.User
	user, _ := mUser.FindByEmail(dhcpBody.Email)
	var netWork models.NetWork
	netWork.FindNetById(user.NetId)
	var ipool IpPool
	for _, pool := range pools {
		if pool.ipNet == netWork.Ip && pool.mask == netWork.Mask {
			ipool = pool
			break
		} else {
			return
		}
	}
	var ip uint32 = 0
	if user.NetIp != 0 {
		if ipool.pool[user.NetIp] {
			DhcpGet.Mes = "ip已被占用！！！"
			c.JSON(http.StatusBadRequest, DhcpGet)
			return
		} else {
			ip = user.NetIp
		}
	} else {
		for i, b := range ipool.pool {
			if !b {
				ip = i
				break
			}
		}
	}
	if ip == 0 || ipool.max == 0 {
		DhcpGet.Mes = "ip已满，请稍后重试！！！"
		c.JSON(http.StatusBadRequest, DhcpGet)
		return
	}
	DhcpGet.Ip = ip
	DhcpGet.Mes = "获取ip成功"
	c.JSON(http.StatusOK, DhcpGet)
	return
}

// Dhcp 获取地址
func Dhcp(c *gin.Context) {
	var dhcpBody dhcpGetBody
	_ = c.ShouldBindJSON(&dhcpBody)
	var mUser models.User
	user, _ := mUser.FindByEmail(dhcpBody.Email)

	var netWork models.NetWork
	netWork.FindNetById(user.NetId)
	var ipool IpPool
	for _, pool := range pools {
		if pool.ipNet == netWork.Ip && pool.mask == netWork.Mask {
			ipool = pool
			break
		} else {
			return
		}
	}

	if ipool.max == 0 {
		DhcpGet.Mes = "ip已满，请稍后重试！！！"
		c.JSON(http.StatusBadRequest, DhcpGet)
		return
	}

	if ipool.pool[dhcpBody.Ip] {
		DhcpGet.Mes = "ip已被占用，刷新重试！！！"
		c.JSON(http.StatusBadRequest, DhcpGet)
		return
	} else {
		ipool.max = ipool.max - 1
		ipool.pool[dhcpBody.Ip] = true
	}

	DhcpGet.Ip = dhcpBody.Ip
	DhcpGet.Mask = netWork.Mask
	DhcpGet.Id = netWork.ID
	DhcpGet.Mes = "获取ip成功"
	DhcpGet.Gateway = netWork.NetGateway
	DhcpGet.Dns = netWork.DnsAddress
	c.JSON(http.StatusOK, DhcpGet)

	var client ClientSpeed
	client.Lose = 0
	client.Email = user.UserEmail
	client.Id = netWork.ID
	client.Ip = dhcpBody.Ip
	client.UserLevel = user.UserLevel

	Clients.Add(&client)
	UpSort = append(UpSort, dhcpBody.Ip)
	sort.Slice(UpSort, func(i, j int) bool {
		return UpSort[i] < UpSort[j]
	})
	return
}

// DisConn 断开连接
func DisConn(c *gin.Context) {
	var disBody struct {
		Id   uint   `json:"id"`
		Ip   uint32 `json:"ip"`
		Mask uint32 `json:"mask"`
	}
	_ = c.ShouldBindJSON(&disBody)
	var ipool IpPool
	for _, pool := range pools {
		if pool.ipNet == (disBody.Ip&disBody.Mask) && pool.mask == disBody.Mask {
			ipool = pool
			break
		}
	}
	ipool.pool[disBody.Ip] = false
	ipool.max = ipool.max + 1
	c.JSON(http.StatusOK, "")

	Clients.Del(disBody.Ip)

	i := sort.Search(len(UpSort), func(i int) bool {
		return UpSort[i] >= disBody.Ip
	})
	if i < len(UpSort) && UpSort[i] == disBody.Ip {
		UpSort = append(UpSort[:i], UpSort[i+1:]...)
	}
	sort.Slice(UpSort, func(i, j int) bool {
		return UpSort[i] < UpSort[j]
	})
	return
}

type IpPool struct {
	ipNet uint32
	mask  uint32
	max   uint64
	pool  map[uint32]bool
}

var pools []IpPool

// initIpPool 初始化地址池
func initIpPool() {
	var netWork models.NetWork
	nets := netWork.FindNet()
	pools = make([]IpPool, len(nets))
	for i, work := range nets {
		pools[i].ipNet = work.Ip
		pools[i].mask = work.Mask
		pools[i].max = work.MaxAllow
		num := work.Mask & (^work.Mask + 1)
		pools[i].pool = make(map[uint32]bool, num)
		var t uint32
		for t = 2; t < num; t++ {
			pools[i].pool[work.Ip+t] = false
		}
	}
}

func init() {
	rand.Seed(time.Now().Unix())
	initIpPool()
}

// RandUp 随机字符串
func RandUp(n int) []byte {
	var longLetters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ=_")
	if n <= 0 {
		return []byte{}
	}
	b := make([]byte, n)
	arc := uint8(0)
	if _, err := rand.Read(b[:]); err != nil {
		return []byte{}
	}
	for i, x := range b {
		arc = x & 63
		b[i] = longLetters[arc]
	}
	return b
}

// ui界面用户信息
type uiUser struct {
	Id uint `json:"id"`
	//用户名
	UserName string `json:"userName"`
	//密码
	UserPass string `json:"userPass"`
	//邮箱
	UserEmail string `json:"userEmail"`
	//域
	NetId uint `json:"netId" `
	//子网
	NetIp net.IP `json:"netIp" `
	//是否允许密码登录
	IsAllowed bool `json:"isAllowed"`
	//用户等级
	UserLevel uint8 `json:"userLevel"`
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var ui uiUser
	var user models.User
	err := c.ShouldBindJSON(&ui)
	if err != nil {
		c.JSON(400, gin.H{
			"mes": "error",
		})
	} else {
		user.UserName = ui.UserName
		user.UserEmail = ui.UserEmail
		user.UserPass = ui.UserPass
		user.NetId = ui.NetId
		user.UserLevel = ui.UserLevel
		user.NetIp = IPToUInt32(ui.NetIp)
		user.UUid, _ = uuid.NewUUID()
		//生成密钥对
		privateKey, _ := sm2.GenerateKey(nil)
		pem, _ := x509.WritePrivateKeyToPem(privateKey, nil)
		err := ioutil.WriteFile("././conf/client_private_key.txt", pem, 0666)
		if err != nil {
			logrus.Println("写入密钥失败！！！")
			return
		}
		key := privateKey.Public().(*sm2.PublicKey)
		user.PublicKey, _ = x509.WritePublicKeyToPem(key)
		//生成google—code
		secret, _ := google.NewGoogleAuth().InitAuth(user.UserEmail)
		url := google.NewGoogleAuth().GetQrcodeUrl(user.UserEmail, secret)
		user.GoogleCode = secret
		user.CreateUser()
		c.JSON(200, gin.H{
			"url": url,
			"key": string(pem),
		})
	}
}

// FindByID 查找信息
func FindByID(c *gin.Context) {
	var user models.User
	intNum, _ := strconv.Atoi(c.Param("id"))
	id := uint(intNum)
	err := user.FindById(id)
	if err != nil {
		return
	}
}

// FindByName 查找信息
func FindByName(name string) (id uint) {
	var user models.User
	id, err := user.FindByName(name)
	if err != nil || id == 0 {
		logrus.Println("用户名查找错误：", err)
		return
	}
	return
}

// FindByEmail 查找信息
func FindByEmail(email string) (user models.User) {
	user, err := user.FindByEmail(email)
	if err != nil {
		logrus.Println("邮箱查找错误：", err)
		return
	}
	return
}

// FindUi ui界面展示的信息
func FindUi(c *gin.Context) {
	var user models.User
	users := user.FindUi()
	uUser := make([]uiUser, len(users))
	for i, user := range users {
		uUser[i].Id = user.Id
		uUser[i].UserName = user.UserName
		uUser[i].UserEmail = user.UserEmail
		uUser[i].UserPass = user.UserPass
		uUser[i].IsAllowed = user.IsAllowed
		uUser[i].NetId = user.NetId
		uUser[i].NetIp = UInt32ToIP(user.NetIp)
		uUser[i].UserLevel = user.UserLevel
	}
	c.JSON(200, uUser)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	var user models.User
	user.DeleteById(c.Query("id"))
	c.JSON(200, gin.H{"mes": "success"})
}

// Update 更新信息
func Update(c *gin.Context) {
	var user models.User
	var uUser uiUser
	err := c.ShouldBindJSON(&uUser)
	if err != nil {
		c.JSON(400, gin.H{
			"mes": err,
		})
	} else {
		user.ID = uUser.Id
		user.UserName = uUser.UserName
		user.UserEmail = uUser.UserEmail
		user.IsAllowed = uUser.IsAllowed
		user.NetId = uUser.NetId
		user.NetIp = IPToUInt32(uUser.NetIp)
		user.UserLevel = uUser.UserLevel
		user.UpdateById()
		c.JSON(200, gin.H{
			"mes": "success",
		})
	}
}
