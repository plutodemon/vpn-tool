package https

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"ivs-net-server/auth/configure"
	"strings"
	"sync"
	"time"
)

var router *gin.Engine
var sessionCache *cache.Cache

func init() {

	// 创建一个默认过期时间为 1 分钟的缓存，每 10 分钟清除一次过期项目
	sessionCache = cache.New(1*time.Minute, 10*time.Minute)
	router = gin.Default()
	router.Use(MyMiddleware())
	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", Login)
		userGroup.POST("/login1", Login1)
		userGroup.POST("/login2", Login2)
		userGroup.POST("/login3", Login3)
		userGroup.POST("/dhcp", Dhcp)
		userGroup.POST("/address", AddressIp)
		userGroup.POST("/disconn", DisConn)
		userGroup.GET("/find/:id", FindByID)
		userGroup.POST("/update", Update)
	}
	uiGroup := router.Group("/ui")
	{
		uiGroup.POST("/createUser", CreateUser)
		uiGroup.GET("/findUser", FindUi)
		uiGroup.GET("/deleteUser/", DeleteUser)
		uiGroup.POST("/editUser", Update)

		uiGroup.POST("/createNet", CreateNet)
		uiGroup.GET("/findNet", FindNet)
		uiGroup.GET("/deleteNet/", DeleteNet)
		uiGroup.POST("/editNet", EditNet)

		uiGroup.GET("/getSpeed", GetSpeed)
		uiGroup.GET("/getClientSpeed", GetClientSpeed)
	}

}

// MyMiddleware 中间件
func MyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*	// 请求前
			//获取资源
			obj := c.Request.URL.RequestURI()
			//获取方法
			act := c.Request.Method
			//获取实体
			var user *models.User
			_ = c.ShouldBindJSON(&user)
			userId := int(FindByName(user.UserName))
			fmt.Println(obj, act, userId)

			//策略校验
			allow, _ := power.Enforcer.Enforce(userId, 1, "pass")
			logrus.Println("校验结果:", allow)*/

		path := c.FullPath()

		if strings.HasPrefix(path, "/ui") {
			if Regular != nil {
				Regular.Stop()
			}
			if ws != nil {
				ws.Close()

			}
			if wsC != nil {
				wsC.Close()
			}
		}

		/*if path == "/ui/*" {
			if Regular != nil {
				Regular.Stop()
			}
			flagCorn = true
		} else if flagCorn && path != "/ui/getSpeed" {
			ws.Close()
			Regular.Stop()
			flagCorn = false
		}*/

		c.Next()

		return
	}
}

func RunHttpsServer(WG *sync.WaitGroup) {
	defer WG.Done()
	err := router.Run(configure.Config.Get("servers.https.port").(string))
	if err != nil {
		return
	}
}
