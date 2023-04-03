package models

// User 用户表
type User struct {
	//用户名
	UserName string `json:"userName"`
	//密码
	UserPass string `json:"userPass"`
	//公钥
	PublicKey []byte `json:"publicKey"`
	//私钥
	PrivateKey []byte
	//邮箱
	UserEmail string `json:"userEmail"`
	//google code
	GoogleCode string `json:"googleCode"`
}
