package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"io"
	"math/rand"
	"os"
	"time"
)

// 将session信息进行保存（测试使用）
var cookieStore = sessions.NewCookieStore([]byte("studyEcho"))

// 初始化操作
func init() {
	rand.Seed(time.Now().UnixNano())
	os.Mkdir("./log", 0755)
}

// 实现一个登陆的案例
func main() {
	// 创建一个echo实例
	e := echo.New()

	// 配置日志信息
	configureLogger(e)

	// 设置静态路由
	e.Static("img", "./01-demo/img")
	e.File("/favicon.ico", "./01-demo/img/favicon.ico")
	e.File("/logo.png", "./01-demo/img/logo.png")

	// 设置中间件
	setMiddleWare(e)

	// 注册路由
	RegisterRouter(e)

	// 启动服务
	e.Logger.Fatal(e.Start(":2020"))

}

// configureLogger 设置当前服务器的Logger信息
func configureLogger(e *echo.Echo) {
	// 设置日志级别为info
	e.Logger.SetLevel(log.INFO)
	// 记录业务日志到文件中
	logFile, err := os.OpenFile("./log/echo.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	// 设置日志输出位置(文件以及终端)
	e.Logger.SetOutput(io.MultiWriter(logFile, os.Stdout))
}

// setMiddleWare 设置中间件
func setMiddleWare(e *echo.Echo) {
	// access log输出到文件中
	accessLog, err := os.OpenFile("./log/access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	// 设置日志文件的输出路径（文件以及终端）
	middleware.DefaultLoggerConfig.Output = accessLog
	middleware.DefaultLoggerConfig.Output = os.Stdout

	// 设置对应的中间件
	e.Use(middleware.Logger())  // 使用日志中间件记录http请求信息
	e.Use(AutoLogin)            // 使用自定义中间件验证用户的登陆信息
	e.Use(middleware.Recover()) // 使用恢复中间件恢复panic恐慌状态
}

