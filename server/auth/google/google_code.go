package google

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type GoogleAuth struct {
}

func NewGoogleAuth() *GoogleAuth {
	return &GoogleAuth{}
}

func (google *GoogleAuth) un() int64 {
	return time.Now().UnixNano() / 1000 / 30
}

func (google *GoogleAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (google *GoogleAuth) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (google *GoogleAuth) base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func (google *GoogleAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (google *GoogleAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

func (google *GoogleAuth) oneTimePassword(key []byte, data []byte) uint32 {
	hash := google.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := google.toUint32(hashParts)
	return number % 1000000
}

// GetSecret 获取秘钥
func (google *GoogleAuth) GetSecret() string {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, google.un())

	return strings.ToUpper(google.base32encode(google.hmacSha1(buf.Bytes(), nil)))
}

// GetCode 获取动态码
func (google *GoogleAuth) GetCode(secret string) (string, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := google.base32decode(secretUpper)
	if err != nil {
		return "", err
	}
	number := google.oneTimePassword(secretKey, google.toBytes(time.Now().Unix()/30))
	return fmt.Sprintf("%06d", number), nil
}

// GetQrcode 获取动态码二维码内容
func (google *GoogleAuth) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}

// GetQrcodeUrl 获取动态码二维码图片地址,这里是第三方二维码api
func (google *GoogleAuth) GetQrcodeUrl(user, secret string) string {
	qrcode := google.GetQrcode(user, secret)
	width := "200"
	height := "200"
	data := url.Values{}
	data.Set("data", qrcode)
	return "https://api.qrserver.com/v1/create-qr-code/?" + data.Encode() + "&size=" + width + "x" + height + "&ecc=M"
}

// VerifyCode 验证动态码
func (google *GoogleAuth) VerifyCode(secret, code string) (bool, error) {
	_code, err := google.GetCode(secret)
	fmt.Println(_code, code, err)
	if err != nil {
		return false, err
	}
	return _code == code, nil
}

// InitAuth 开启二次认证
func (google *GoogleAuth) InitAuth(user string) (secret, code string) {
	// 秘钥
	secret = NewGoogleAuth().GetSecret()
	fmt.Println("Secret:", secret)

	// 动态码(每隔30s会动态生成一个6位数的数字)
	code, err := NewGoogleAuth().GetCode(secret)
	fmt.Println("Code:", code, err)

	// 用户名
	qrCode := NewGoogleAuth().GetQrcode(user, code)
	fmt.Println("Qrcode", qrCode)

	// 打印二维码地址
	qrCodeUrl := NewGoogleAuth().GetQrcodeUrl(user, secret)
	fmt.Println("QrcodeUrl", qrCodeUrl)

	return
}
