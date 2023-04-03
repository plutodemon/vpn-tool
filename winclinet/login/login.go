package login

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tjfoc/gmsm/sm4"
	"github.com/tjfoc/gmsm/x509"
	"io/ioutil"
	"ivs-net-winclinet/configure"
	"ivs-net-winclinet/models"
	"math/rand"
	"net/http"
	"time"
)

// 第一次得到的服务端512字节
var firstAuth []byte

// PostRequest 第一次POST请求
func PostRequest(encrypt []byte) (rsp RespondBody) {
	req := map[string]interface{}{"res": encrypt}
	body, _ := json.Marshal(req)
	url := urlPrefix + "/user/login1"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		logrus.Println("第一次请求发生错误：", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second} // 设置请求超时时长5s
	//返回体
	resp, err := client.Do(request)
	if err != nil {
		logrus.Println("响应发生错误: ", err)
		return
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 400 {
		_ = json.Unmarshal(respBody, &rsp)
		return
	}
	//处理返回值：私钥解密
	var byteBody ByteBody
	_ = json.Unmarshal(respBody, &byteBody)
	readFile, _ := ioutil.ReadFile("./conf/privateKey.txt")
	pem, _ := x509.ReadPrivateKeyFromPem(readFile, nil)
	firstAuth, _ = pem.DecryptAsn1(byteBody.Res)
	_ = json.Unmarshal(firstAuth, &rsp)
	return
}

// Mes 第一次请求发送的结构体
type Mes struct {
	UserEmail string
	RandNum   []byte
	TimeStamp int64
}

// RespondBody 第一次请求服务端响应字段
type RespondBody struct {
	Mes       string    `json:"Mes"`
	Uuid      uuid.UUID `json:"Uuid"`
	UserName  string    `json:"UserName"`
	UserEmail string    `json:"UserEmail"`
	RandNum   []byte    `json:"RandNum"`
	TimeStamp int64     `json:"TimeStamp"`
}

// ByteBody 响应字节流
type ByteBody struct {
	Res []byte `json:"Res"`
}

// 客户端第一次生成的512字节
var firstClient []byte

// 第一次请求的时间戳
var firstTime int64

// Login1 登录
func Login1(email string) string {
	email = "test1@sina.com"
	//随机字符串
	firstTime = time.Now().Unix()
	mes := Mes{
		UserEmail: email,
		RandNum:   RandUp(504 - len(email)),
		TimeStamp: firstTime,
	}
	firstClient, _ = json.Marshal(mes)

	//公钥加密
	readFile, _ := ioutil.ReadFile("./conf/publicKey.txt")
	pem, _ := x509.ReadPublicKeyFromPem(readFile)
	asn1, _ := pem.EncryptAsn1(firstClient, nil)

	rsp := PostRequest(asn1)
	if rsp.Mes == "error" {
		logrus.Println("登录发生错误（没有此用户）！！！")
		return "没有此用户!"
	} else {
		nowTime := time.Now().Unix()
		if rsp.TimeStamp+30 > nowTime && rsp.TimeStamp-30 < nowTime {
			//fmt.Println("得到的响应体:", rsp)
			logrus.Println("邮箱登录验证成功！")
			return ""
		} else {
			logrus.Println("登录发生错误（超时）！！！")
			return "请再试一次~"
		}
	}

}

// Mes2 第二次请求发送的结构体
type Mes2 struct {
	Rand504   []byte
	TimeStamp int64
	Sha1024   [64]byte
}

// RespondBody2 第二次响应体
type RespondBody2 struct {
	Mes     string `json:"Mes"`
	RandNum []byte `json:"RandNum"`
}

var respondBody2 RespondBody2

// Sha sha512计算的结构体
type Sha struct {
	Rand504   []byte
	TimeStamp int64
	Password  string
}

// 第二次客户端生成的512字节
var secondClient []byte

// 第二次服务端512字节
var secondAuth []byte

// Login2 登录
func Login2(pass string) string {

	pass = "test1"

	var mes2 Mes2
	mes2.Rand504 = RandUp(504)
	mes2.TimeStamp = time.Now().Unix()
	sha := Sha{
		Rand504:   mes2.Rand504,
		TimeStamp: mes2.TimeStamp,
		Password:  pass,
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
	mes2.Sha1024 = sum512
	secondClient, _ = json.Marshal(mes2)
	//公钥加密
	readFile, _ := ioutil.ReadFile("./conf/publicKey.txt")
	pem, _ := x509.ReadPublicKeyFromPem(readFile)
	asn1, _ := pem.EncryptAsn1(secondClient, nil)

	//发送第二次请求
	req := map[string]interface{}{"res": asn1}
	body, _ := json.Marshal(req)
	url := urlPrefix + "/user/login2"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		logrus.Println("第二次请求发生错误：", err)
		return "请再试一次~"
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second} // 设置请求超时时长5s
	//返回体
	resp, err := client.Do(request)
	if err != nil {
		logrus.Println("第二次响应发生错误: ", err)
		return "请再试一次~"
	}
	var byteBody ByteBody

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 400 {
		logrus.Println("登录发生错误：密码错误！！！")
		return "密码错误！"
	}
	//处理返回值
	_ = json.Unmarshal(respBody, &byteBody)
	_ = json.Unmarshal(byteBody.Res, &respondBody2)
	return ""
}

// AllByte 2048字节
type AllByte struct {
	FirstClient  []byte
	SecondClient []byte
	FirstAuth    []byte
	SecondAuth   []byte
}

// Request 第三次发送字节流
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
func Login3(code string) string {
	// sm4解密
	code = "123456"
	code += "0000000000"
	key := []byte(code)
	iv := make([]byte, sm4.BlockSize)
	secondAuth, _ = models.Sm4Decrypt(key, iv, respondBody2.RandNum)
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
	//发送第三次请求
	req := Request{Res: sum8192}
	body, _ := json.Marshal(req)
	url := urlPrefix + "/user/login3"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		logrus.Println("第三次请求发生错误：", err)
		return "请再试一次~"
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second} // 设置请求超时时长10s
	//返回体
	resp, err := client.Do(request)
	if err != nil {
		logrus.Println("第三次响应发生错误: ", err)
		return "请再试一次~"
	}
	if resp.StatusCode == 400 {
		logrus.Println("登录发生错误（验证码错误）！！！")
		return "验证码错误！"
	} else {
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
		logrus.Println("登录成功！！！")
		return "登录成功！"
	}

}

type DhcpBody struct {
	Mes     string `json:"Mes"`
	Id      uint   `json:"Id"`
	Ip      uint32 `json:"Ip"`
	Mask    uint32 `json:"Mask"`
	Gateway uint32 `json:"Gateway"`
	Dns     uint32 `json:"Dns"`
}

func dhcpIp(u string, userEmail string, ip uint32) (flag bool, dhcpBody DhcpBody) {
	var req map[string]interface{}
	if ip == 0 {
		req = map[string]interface{}{"email": userEmail}
	} else {
		req = map[string]interface{}{"email": userEmail, "ip": ip}
	}
	body, _ := json.Marshal(req)
	url := urlPrefix + u
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		logrus.Println("dhcp获取失败：", err)
		return false, dhcpBody
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second} // 设置请求超时时长5s
	//返回体
	resp, err := client.Do(request)
	if err != nil {
		logrus.Println("dhcp获取失败：", err)
		return false, dhcpBody
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 400 {
		_ = json.Unmarshal(respBody, &dhcpBody)
		return false, dhcpBody
	}
	_ = json.Unmarshal(respBody, &dhcpBody)
	return true, dhcpBody
}

// AddressIp 获得地址列表
func AddressIp(userEmail string, ip uint32) (flag bool, dhcpBody DhcpBody) {
	return dhcpIp("/user/address", userEmail, ip)
}

// Dhcp 获得ip地址
func Dhcp(userEmail string, ip uint32) (flag bool, dhcpBody DhcpBody) {
	return dhcpIp("/user/dhcp", userEmail, ip)
}

// DisConn 断开连接
func DisConn(dhcpBody DhcpBody) bool {
	if dhcpBody.Ip == 0 {
		return false
	}
	req := map[string]interface{}{"id": dhcpBody.Id, "ip": dhcpBody.Ip, "mask": dhcpBody.Mask}
	body, _ := json.Marshal(req)
	url := urlPrefix + "/user/disconn"
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second} // 设置请求超时时长5s
	//返回体
	resp, _ := client.Do(request)
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		logrus.Println("断开连接失败!")
		return false
	} else {
		return true
	}
}

var urlPrefix string

func init() {
	rand.Seed(time.Now().Unix())
	GetUrl()
}
func GetUrl() {
	urlPrefix = "http://" + configure.Config.Get("servers.auth_address").(string) + ":" + configure.Config.Get("servers.auth_port").(string)
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
