package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var atoz map[string]int = map[string]int{
	"a":1,
	"b":2,
	"c":3,
	"d":4,
	"e":5,
	"f":6,
	"g":7,
	"h":8,
	"i":9,
	"j":10,
	"k":11,
	"l":12,
	"m":13,
	"n":14,
	"o":15,
	"p":16,
	"q":17,
	"r":18,
	"s":19,
	"t":20,
	"u":21,
	"v":22,
	"w":23,
	"x":24,
	"y":25,
	"z":26,
}

func replaceAtPosition(originaltext string, indexofcharacter int, replacement string) string {
	runes := []rune(originaltext )
	partOne := string(runes[0:indexofcharacter-1])
	partTwo := string(runes[indexofcharacter:len(runes)])
	return partOne + replacement + partTwo
}

// Replace the nth occurrence of old in s by new.
func replaceNth(s, old, new string, n int) string {
	i := 0
	for m := 1; m <= n; m++ {
		x := strings.Index(s[i:], old)
		if x < 0 {
			break
		}
		i += x
		if m == n {
			return s[:i] + new + s[i+len(old):]
		}
		i += len(old)
	}
	return s
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}


// RegisterRouter 注册请求路由处理
func RegisterRouter(e *echo.Echo) {
	e.Static("/static", "")
	// 首页路由
	e.GET("/", func(ctx echo.Context) error {
		// 获取当前的请求参数msg的值
		data := map[string]interface{}{
			"msg": ctx.QueryParam("msg"),
		}
		// 判断用户是否已经登陆，从context中获取user
		if user, ok := ctx.Get("user").(*User); ok {
			// 加载模版信息
			tpl := template.Must(template.ParseFiles("./board.html"))
			ctx.Logger().Info("this is login page")
			// 用户已经登陆则从context中获取用户的信息
			data["username"] = user.Username
			data["had_login"] = true
			// 将模版以及数据信息写入到缓冲区中
			var buf bytes.Buffer
			err := tpl.Execute(&buf, data)
			if err != nil {
				return err
			}
			fmt.Println(user.Username)
			// 将模版信息以html的方式返回
			return ctx.HTML(http.StatusOK, buf.String())
		} else {
			// 加载模版信息
			tpl := template.Must(template.ParseFiles("./login.html"))
			ctx.Logger().Info("this is login page")
			// 用户没有登陆则从session中获取用户的登陆信息
			sess := getCookieSession(ctx)
			if flashes := sess.Flashes("username"); len(flashes) > 0 {
				data["username"] = flashes[0]
			}
			sess.Save(ctx.Request(), ctx.Response())
			// 将模版以及数据信息写入到缓冲区中
			var buf bytes.Buffer
			err := tpl.Execute(&buf, data)
			if err != nil {
				return err
			}
			// 将模版信息以html的方式返回
			return ctx.HTML(http.StatusOK, buf.String())
		}

	})
	// 登陆路由
	e.POST("/login", func(ctx echo.Context) error {
		// 获取请求参数
		username := ctx.FormValue("username")
		passwd := ctx.FormValue("passwd")
		remember_me := ctx.FormValue("remember_me")
		if username == "admin" && passwd == "123456" || username == "mzj" && passwd == "123456"  || username == "gdy" && passwd == "123456" || username == "zjb" && passwd == "123456" || username == "hkh" && passwd == "123456" || username == "wxm" && passwd == "123456" {
			// 用户名密码正确则用标准库种cookie
			cookie := &http.Cookie{
				Name:     "username",
				Value:    username,
				HttpOnly: true,
			}
			// 查看用户是否选择记住用户名
			if remember_me == "1" {
				cookie.MaxAge = 7 * 24 * 3600 // 7天
			}
			ctx.SetCookie(cookie)
			// 重定向到首页
			return ctx.Redirect(http.StatusSeeOther, "/board")
		}
		// 用户名密码错误则返回相应信息
		// 首先使用session保存当前用户的用户名信息
		session := getCookieSession(ctx)
		session.AddFlash("username", username)
		err := session.Save(ctx.Request(), ctx.Response())
		if err != nil {
			return ctx.Redirect(http.StatusSeeOther, "/?msg="+err.Error())
		}
		return ctx.Redirect(http.StatusSeeOther, "/?msg=用户名或者密码错误")
	})
	// 注销路由
	e.GET("/logout", func(ctx echo.Context) error {
		// 覆盖当前cookie
		cookie := &http.Cookie{
			Name:    "username",
			Value:   "",
			Expires: time.Now().Add(-1e9),
			MaxAge:  -1,
		}
		ctx.SetCookie(cookie)
		return ctx.Redirect(http.StatusSeeOther, "/")
	})
	//公告板路由
	e.GET("/board", func(ctx echo.Context)error {
		htmlstr := `        <div class="scroller _h" style="overflow-x: scroll; box-sizing: border-box; margin: 0px; border: 0px;"><div class="wea-doc-detail-content-main">

            <div id="weaDocDetailHtmlContent" style=""><style type="text/css">.mouldDemoMain {
                margin: 5% auto;
            }
            .mouldDemoMainContent {
                margin-top:20px;
                padding: 0 12%;
                position: relative;
            }
            .mouldDemoMainTxt {
                text-align: left;
                text-justify: inter-ideograph;
            }
            </style>
                <div class="mouldDemoMain">
                    <div class="mouldDemoMainContent">
                        <div class="mouldDemoMainTxt">
                            <div style="min-height: 450px; margin-top: 21px;"><div class="Section0"><p align="center" style="text-align:center"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:10.5pt"><span style="mso-pagination:widow-orphan"><span style="line-height:30.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:Calibri"><b><span style="font-size:28.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#ff0000"><span style="letter-spacing:1.3000pt"><span style="font-weight:bold"><span style="mso-font-kerning:1.0000pt"><font face="宋体">温氏股份养禽事业部文件</font></span></span></span></span></span></span></span></span></b></span></span></span></span></span></span></span></p><p style="text-align:justify">&nbsp;</p><p style="text-align:justify">&nbsp;</p><p align="center" style="text-align:center"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:10.5pt"><span style="mso-pagination:widow-orphan"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:Calibri"><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:0.0000pt"><font face="仿宋_GB2312">养禽事业部办〔</font>2021<font face="仿宋_GB2312">〕</font></span></span></span></span></span><span style="font-family:Times New Roman,Times,serif"><span style="font-size:16.0000pt">24</span></span><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:0.0000pt"><font face="仿宋_GB2312">号</font></span></span></span></span></span></span></span></span></span></span></span></span></p><hr style="width:100%; height:2px; background:#ff0000; border:0px"><p align="center" class="p" style="margin-top:5px; margin-bottom:5px; text-align:center"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:12pt"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><b><span style="font-size:22.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="font-weight:bold"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="宋体">关于召开</font></span></span></span></span></span></span></span></span></span></span></b><b><span style="font-size:22.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="font-weight:bold"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">202</span></span></span></span></span></span></span></span></span></span></b><b><span style="font-size:22.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="font-weight:bold"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">2<font face="宋体">年温氏股份养禽事业部</font></span></span></span></span></span></span></span></span></span></span></b></span></span></span></span></span></span><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:12pt"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><b><span style="font-size:22.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="font-weight:bold"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="宋体">工作大会的通知</font></span></span></span></span></span></span></span></span></span></span></b></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-align:justify">&nbsp;</p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-align:justify"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:12pt"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">养禽事业部下属各单位：</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:12pt"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">经研究，兹定于</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">202</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">2</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">年</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">1</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">月</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">6</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">日（星期四）下午召开</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">202</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">2</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">年温氏股份养禽事业部工作大会，具体事宜通知如下：</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:12pt"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><b><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="font-weight:bold"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">一、会议安排</font></span></span></span></span></span></span></span></span></span></span></b></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:12pt"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-size:16pt"><span style="font-family:宋体"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">（一）时间：</font></span></span></span></span></span></span></span></span></span></span><span style="font-size:16pt"><span style="font-family:宋体"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">202</span></span></span></span></span></span></span></span></span></span><span style="font-size:16pt"><span style="font-family:宋体"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">2</span></span></span></span></span></span></span></span></span></span><span style="font-size:16pt"><span style="font-family:宋体"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">年</font></span></span></span></span></span></span></span></span></span></span><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">1</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">月</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">6</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">日</font></span></span></span></span></span></span></span></span></span></span><span style="font-size:16pt"><span style="font-family:宋体"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">（星期四），下午</font>14:30~17:45</span></span></span></span></span></span></span></span></span></span><span style="font-size:16pt"><span style="font-family:宋体"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">，报到时间</font></span></span></span></span></span></span></span></span></span></span><span style="font-size:16pt"><span style="font-family:宋体"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">14:00~14:25</span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:12pt"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">（二）地点：</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">集团总部大楼一楼影剧院</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:12pt"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">（三）会议形式：现场会</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">+</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">视频会</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:12pt"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">（四）会议主持人：简仿辉副总裁</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:12pt"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">（五）参会范围</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0500pt; text-align:justify"><span style="layout-grid:15.6000pt"><span style="page:Section0"><span style="font-size:12pt"><span style="mso-pagination:none"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">1.</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">现场会：事业部总裁级干部；事业部负责全面工作的总经理级干部；云浮市内单位副总经理、高级技术职务人员、总经理助理；云浮市内四级公司行政经理级干部；事业部职能部门科室负责人；集团</font>“<font face="仿宋_GB2312">降本增效</font><font face="Times New Roman">”</font><font face="仿宋_GB2312">标兵单位负责人</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:仿宋_GB2312"><span style="mso-ascii-font-family:'Times New Roman'"><span style="mso-hansi-font-family:'Times New Roman'"><span style="mso-bidi-font-family:'Times New Roman'"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">。</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p></div><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0500pt; text-align:justify"><span style="font-size:12pt"><span style="mso-pagination:none"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">2.</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">视频会：云浮市外三级公司副总经理、总经理助理，云浮市外四级公司行政经理级干部。</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">现场会具体名单后续将另发传阅通知到各参会人员。</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">（六）会议议程</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">1.</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">奏唱温氏之歌；</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">2.</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">养禽事业部财务信息部总经理梁冰飞作《养禽事业部</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">20</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">21</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">年生产经营情况通报暨</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">202</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">2</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">年经营预算》；</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">3.</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">养禽事业部副总裁覃健萍作《</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">20</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">22</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">年养禽事业部工作报告》；</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">4</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">.</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">集团首席运营官兼养禽事业部总裁秦开田作</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">202</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">2</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">年养禽事业部工作重点解读；</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">5</span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">.</span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">集团领导讲话；</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">6.202</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">1</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">年养禽事业部专项奖表彰；</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">7.</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">参会全体人员朗诵</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">“</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">温氏十大正气观</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">”</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">。</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><b><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="font-weight:bold"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">二、晚宴安排</font></span></span></span></span></span></span></span></span></span></span></b></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">事业部总裁级干部，各单位行政总经理级干部、高级技术职务人员、总经理助理，云浮市内四级公司行政经理级干部，统一参加集团在总部三楼餐厅组织的晚宴。</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:仿宋_GB2312"><span style="mso-hansi-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="text-transform:none"><span style="font-style:normal"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">事业部职能部门科室负责人用餐另行安排</font></span></span></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:仿宋_GB2312"><span style="mso-ascii-font-family:'Times New Roman'"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="text-transform:none"><span style="font-style:normal"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">。</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p class="p" style="margin-top:5px; margin-bottom:5px; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="mso-list:l0 level1 lfo1"><span style="font-family:宋体"><b><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="font-weight:bold"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">&nbsp; &nbsp; 三、现场会相关要求</font></span></span></span></span></span></span></span></span></span></span></b></span></span></span></span></span></span></span></p><p style="text-indent:32.0000pt; text-align:justify"><span style="font-size:10.5pt"><span style="mso-char-indent-count:2.0000"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:Calibri"><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt">1.</span></span></span></span></span></span></span></span></span><span style="font-family:Calibri"><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt"><font face="仿宋_GB2312">参会人员</font></span></span></span></span></span><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:仿宋_GB2312"><span style="mso-ascii-font-family:'Times New Roman'"><span style="mso-hansi-font-family:'Times New Roman'"><span style="mso-bidi-font-family:'Times New Roman'"><span style="mso-font-kerning:1.0000pt"><font face="仿宋_GB2312">须着</font></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt"><font face="仿宋_GB2312">正装（男：深色西装</font>+<font face="仿宋_GB2312">白色衬衫；女：深色套装；主席台领导打领带）</font></span></span></span></span></span><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:仿宋_GB2312"><span style="mso-ascii-font-family:'Times New Roman'"><span style="mso-hansi-font-family:'Times New Roman'"><span style="mso-bidi-font-family:'Times New Roman'"><span style="mso-font-kerning:1.0000pt"><font face="仿宋_GB2312">，</font></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:仿宋_GB2312"><span style="mso-hansi-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="text-transform:none"><span style="font-style:normal"><span style="mso-font-kerning:1.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">佩戴温氏股份胸徽和口罩</font></span></span></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:仿宋_GB2312"><span style="mso-hansi-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="text-transform:none"><span style="font-style:normal"><span style="mso-font-kerning:1.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">。</font></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p style="text-indent:32.0000pt; text-align:justify"><span style="font-size:10.5pt"><span style="mso-char-indent-count:2.0000"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:Calibri"><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt">2.<font face="仿宋_GB2312">根据疫情防控工作要求，</font></span></span></span></span></span><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt"><font face="仿宋_GB2312">省外参会人员到新兴后要求做一次核酸检测；</font></span></span></span></span></span><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt"><font face="仿宋_GB2312">参会人员会前</font>14<font face="仿宋_GB2312">天内如有新冠肺炎疑似症状、疫情严重地区人员接触史、疫情严重地区驻留史或其他任何疑似情况的，</font></span></span></span></span></span><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt"><font face="仿宋_GB2312">请立即与会务人员联系。</font></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:宋体"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">3</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">.</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">参会人员如需请假，请于</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">1</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">月</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">4</span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">日前报事业部行政综合部简诗桐处，汇总后报事业部领导审批。</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p style="text-indent:32.0000pt; text-align:justify"><span style="font-size:10.5pt"><span style="mso-char-indent-count:2.0000"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:Calibri"><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt">4.<font face="仿宋_GB2312">事业部下属各级单位要尽快组织本级年度工作大会，将更多的精力投入工作中。</font></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">会议联系人：简诗桐、周琼厚</font></span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify"><span style="font-size:12pt"><span style="background:#ffffff"><span style="text-justify:inter-ideograph"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:宋体"><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff"><font face="仿宋_GB2312">联系电话：</font></span></span></span></span></span></span></span></span></span><span style="font-size:16.0000pt"><span style="background:#ffffff"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:sans-serif"><span style="color:#000000"><span style="letter-spacing:-0.0500pt"><span style="mso-font-kerning:0.0000pt"><span style="mso-shading:#ffffff">0766-2929693</span></span></span></span></span></span></span></span></span></span></span></span></span></span></span></p><p align="justify" class="p" style="margin-top:5px; margin-bottom:5px; text-indent:32.0000pt; text-align:justify">&nbsp;</p><p style="text-align:justify">&nbsp;</p><p style="text-indent:32.0000pt; text-align:justify"><span style="font-size:10.5pt"><span style="mso-char-indent-count:2.0000"><span style="mso-pagination:widow-orphan"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:Calibri"><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt">&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;<font face="仿宋_GB2312">温氏股份养禽事业部</font></span></span></span></span></span></span></span></span></span></span></span></p><p style="text-indent:232.0000pt; text-align:justify"><span style="font-size:10.5pt"><span style="mso-char-indent-count:14.5000"><span style="mso-pagination:widow-orphan"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:Calibri"><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt">&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; 2021<font face="仿宋_GB2312">年</font><font face="Times New Roman">12</font><font face="仿宋_GB2312">月</font><font face="Times New Roman">29</font><font face="仿宋_GB2312">日</font></span></span></span></span></span></span></span></span></span></span></span></p><p style="text-indent:-0.0500pt; text-align:justify">&nbsp;</p><p style="text-indent:-0.0500pt; text-align:justify">&nbsp;</p><p style="text-indent:-0.0500pt; text-align:justify">&nbsp;</p><p style="text-align:justify">&nbsp;</p><p style="text-indent:-0.0500pt; text-align:justify">&nbsp;</p><hr><p style="text-indent:16.0000pt; text-align:justify"><span style="font-size:10.5pt"><span style="mso-char-indent-count:1.0000"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:Calibri"><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt"><font face="仿宋_GB2312">养禽事业部行政综合部&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;</font> 2021<font face="仿宋_GB2312">年</font><font face="Times New Roman">12</font><font face="仿宋_GB2312">月</font><font face="Times New Roman">29</font><font face="仿宋_GB2312">日印发</font></span></span></span></span></span><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="letter-spacing:-1.0000pt"><span style="mso-font-kerning:1.0000pt"><font face="仿宋_GB2312">（共印</font>1<font face="仿宋_GB2312">份）</font></span></span></span></span></span></span></span></span></span></span></span></p><hr><p style="margin-right:39px; mso-para-margin-right:-0.7000gd; text-align:justify"><span style="font-size:10.5pt"><span style="line-height:28.0000pt"><span style="mso-line-height-rule:exactly"><span style="font-family:Calibri"><span style="position:absolute; margin-left:39px; margin-top:39px; width:605.0000px"><span style="z-index:1"><span style="height:2.0000px"><img alt="" height="2px" src="file:///C:\Users\ADMINI~1\AppData\Local\Temp\ksohtml4600\wps2.png" width="605px"></span></span></span><span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt"><font face="仿宋_GB2312">拟稿：周琼厚</font>&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp;<font face="仿宋_GB2312">校对：陈永华</font>&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;<font face="仿宋_GB2312">签发：</font></span></span></span></span></span> <span style="font-size:16.0000pt"><span style="mso-spacerun:'yes'"><span style="font-family:'Times New Roman'"><span style="mso-fareast-font-family:仿宋_GB2312"><span style="mso-font-kerning:1.0000pt"><font face="仿宋_GB2312">秦开田</font></span></span></span></span></span></span></span></span></span></p><hr></div>
                        </div>
                    </div>
                </div></div>
        </div></div>
`
		rexp1 := `mai诗`
		fmt.Println(regexp.MatchString(rexp1, htmlstr))

		re3, _ := regexp.Compile("诗");
		rep2 := re3.ReplaceAllString(htmlstr, "无字天书");
		fmt.Println(rep2);

		re, _ := regexp.Compile("[，。：；、]");
		all_ix := re.FindAllIndex([]byte(htmlstr), -1);
		fmt.Println(all_ix);

		//FindAllSubmatch查找所有子匹配项
		fmt.Println("sub matches")
		all_sub := re.FindAllSubmatch([]byte(htmlstr), -1);
		fmt.Println(string(all_sub[0][0]));
		fmt.Println(string(all_sub[1][0]));
		for i := 0; i < len(all_sub); i++ {
			fmt.Println("for each::")
			fmt.Println(string(all_sub[i][0]))
		}

		n := int64(8)
		s := fmt.Sprintf("%b", n)
		fmt.Println(s)
		s2 := fmt.Sprintf("%07b",n)
		fmt.Println(s2)
		for i := 0; i < len(s2); i++ {
			ch := string(s2[i])
			fmt.Println(ch)
		}


		all_sub_ix := re.FindAllSubmatchIndex([]byte(htmlstr), -1);
		fmt.Println(all_sub_ix);

		ret, _ := regexp.Compile("[中。；：、]");
		all_ix = ret.FindAllIndex([]byte("<div>    中文"), -1);
		fmt.Println(all_ix);
		fmt.Println(replaceAtPosition("<div>    中文", 10, "a中" ));
		//这个函数不行，懒得追究原因
		//fmt.Println(replaceAtPosition(htmlstr, 7006, "麦"));

		//strings.replace只能替换前n个字符串
		//htmlstr = strings.Replace(htmlstr, "，", "%，", 2);

		//替换第n个字符
		//htmlstr = replaceNth(htmlstr, "，", " ，", 2);



		// 加载模版信息
		tpl := template.Must(template.ParseFiles("./board.html"))
		ctx.Logger().Info("this is login page")
		// 获取当前的请求参数msg的值
		//data := map[string]interface{}{
		//	"msg": ctx.QueryParam("msg"),
		//}
		data := map[string]interface{}{
			"msg": template.HTML(htmlstr),
		}
		// 判断用户是否已经登陆，从context中获取user
		//初始化二进制表达
		bstr := ""
		if user, ok := ctx.Get("user").(*User); ok {
			// 用户已经登陆则从context中获取用户的信息
			data["username"] = user.Username
			data["had_login"] = true
			fmt.Printf("登录用户：——> %s \n", user.Username)
			//映射用户名到数字组合
			for i := 0; i < len(user.Username); i++ {
				ch := string(user.Username[i])
				fmt.Println(atoz[ch])
				//转换为二进制
				sname := fmt.Sprintf("%07b",atoz[ch])
				//表达1-26只需要5个二进制数，故从2开始，反正前面两位一直是0
				for i := 2; i < len(sname); i++ {
					ch := string(sname[i])
					bstr = bstr + ch
				}
			}
			fmt.Println(bstr)


			//重复2次，增加抗干扰能力
			bstr = bstr + bstr + bstr + bstr + bstr + bstr

			l :=0;
			m := 0;
			o :=0;
			p := 0;
			q := 0;
			for i := 0; i < len(bstr) && i < len(all_sub); i++ {
				if string(all_sub[i][0]) == "，"{
					l = l + 1;
				}
				if string(all_sub[i][0]) == "。"{
					m = m + 1;
				}
				if string(all_sub[i][0]) == "："{
					o = o + 1;
				}
				if string(all_sub[i][0]) == "；"{
					p = p + 1;
				}
				if string(all_sub[i][0]) == "、"{
					q = q + 1;
				}
				fmt.Println("for each bstr::")
				fmt.Println(string(bstr[i]))
				if string(bstr[i]) == "1" && user.Username != "admin" {
					if string(all_sub[i][0]) == "，"{
						htmlstr = replaceNth(htmlstr, "，", " ，", l);
					}
					if string(all_sub[i][0]) == "。"{
						htmlstr = replaceNth(htmlstr, "。", " 。", m);
					}
					if string(all_sub[i][0]) == "："{
						htmlstr = replaceNth(htmlstr, "：", " ：", o);
					}
					if string(all_sub[i][0]) == "；"{
						htmlstr = replaceNth(htmlstr, "；", " ；", p);
					}
					if string(all_sub[i][0]) == "、"{
						htmlstr = replaceNth(htmlstr, "、", " 、", q);
					}
				}
				//admin用户对所有目标符号标记
				if user.Username == "admin" {
					if string(all_sub[i][0]) == "，"{
						htmlstr = replaceNth(htmlstr, "，", "<span style=\"background-color: #ff6600;\">，</span>", l);
					}
					if string(all_sub[i][0]) == "。"{
						htmlstr = replaceNth(htmlstr, "。", "<span style=\"background-color: #ff6600;\">。</span>", m);
					}
					if string(all_sub[i][0]) == "："{
						htmlstr = replaceNth(htmlstr, "：", "<span style=\"background-color: #ff6600;\">：</span>", o);
					}
					if string(all_sub[i][0]) == "；"{
						htmlstr = replaceNth(htmlstr, "；", "<span style=\"background-color: #ff6600;\">；</span>", p);
					}
					if string(all_sub[i][0]) == "、"{
						htmlstr = replaceNth(htmlstr, "、", "<span style=\"background-color: #ff6600;\">、</span>", q);
					}
				}

			}
			data["msg"] =  template.HTML(htmlstr)

		} else {
			// 用户没有登陆则从session中获取用户的登陆信息
			sess := getCookieSession(ctx)
			if flashes := sess.Flashes("username"); len(flashes) > 0 {
				data["username"] = flashes[0]
			}
			sess.Save(ctx.Request(), ctx.Response())
			return ctx.Redirect(http.StatusSeeOther, "/?msg=你又没有登录是吧？")
		}
		// 将模版以及数据信息写入到缓冲区中
		var buf bytes.Buffer
		err := tpl.Execute(&buf, data)
		if err != nil {
			return err
		}
		// 将模版信息以html的方式返回
		return ctx.HTML(http.StatusOK, buf.String())
	})
	//公告板路由
	e.GET("/download", func(ctx echo.Context)error {
		ctx.Logger().Info("downloading")
		// 获取当前的请求参数msg的值
		data := map[string]interface{}{
			"msg": ctx.QueryParam("msg"),
		}
		// 判断用户是否已经登陆，从context中获取user
		if user, ok := ctx.Get("user").(*User); ok {
			// 用户已经登陆则从context中获取用户的信息
			data["username"] = user.Username
			data["had_login"] = true
			// Read from docx file
			r, err := ReadDocxFile("./oa_file_download.docx")
			// Or read from memory
			// r, err := docx.ReadDocxFromMemory(data io.ReaderAt, size int64)
			if err != nil {
				panic(err)
			}
			docx1 := r.Editable()
			// Replace like https://golang.org/pkg/strings/#Replace
			//docx1.Replace("test", "225776", -1)
			//把用户名写入docx属性
			docx1.ReplaceCustom("KSOProductBuildVer",user.Username, -1)
			//docx1.Replace("old_1_2", "new_1_2", -1)
			//docx1.ReplaceLink("http://example.com/", "https://github.com/nguyenthenguyen/docx", 1)
			//docx1.ReplaceHeader("out with the old", "in with the new")
			//docx1.ReplaceFooter("Change This Footer", "new footer")
			docx1.WriteToFile("./new_result" + GetMD5Hash(user.Username)+ ".docx")
			fmt.Println("file generated")
			fmt.Println(fmt.Sprintf("%s",GetMD5Hash(user.Username)))

			//docx2 := r.Editable()
			//docx2.Replace("old_2_1", "new_2_1", -1)
			//docx2.Replace("old_2_2", "new_2_2", -1)
			//docx2.WriteToFile("./new_result_2.docx")

			// Or write to ioWriter
			// docx2.Write(ioWriter io.Writer)

			r.Close()
			file := "./new_result" + GetMD5Hash(user.Username)+ ".docx"
			// 讀取檔案
			downloadBytes, err := ioutil.ReadFile(file)

			if err != nil {
				fmt.Println(err)
			}
			// 取得檔案的 MIME type
			mime := http.DetectContentType(downloadBytes)
			fileSize := len(string(downloadBytes))
			ctx.Response().Header().Set("Content-Type", mime)
			ctx.Response().Header().Set("Content-Disposition", "attachment; filename="+"秘密文档.docx")
			ctx.Response().Header().Set("Content-Length", strconv.Itoa(fileSize))
			return ctx.Stream(http.StatusOK, echo.MIMEOctetStream, bytes.NewReader(downloadBytes))
		} else {
			// 用户没有登陆则从session中获取用户的登陆信息
			sess := getCookieSession(ctx)
			if flashes := sess.Flashes("username"); len(flashes) > 0 {
				data["username"] = flashes[0]
			}
			sess.Save(ctx.Request(), ctx.Response())
			return ctx.Redirect(http.StatusSeeOther, "/?msg=登陆了没。。。？不登陆不能下载啊，这是基本常识")
		}
	})

	// 首页路由
	e.GET("/check", func(ctx echo.Context) error {
		// 获取当前的请求参数msg的值
		data := map[string]interface{}{
			"msg": ctx.QueryParam("msg"),
		}
		// 判断用户是否已经登陆，从context中获取user
		if user, ok := ctx.Get("user").(*User); ok {
			// 加载模版信息
			tpl := template.Must(template.ParseFiles("./check.html"))
			ctx.Logger().Info("this is login page")
			// 用户已经登陆则从context中获取用户的信息
			data["username"] = user.Username
			data["had_login"] = true
			// 将模版以及数据信息写入到缓冲区中
			var buf bytes.Buffer
			err := tpl.Execute(&buf, data)
			if err != nil {
				return err
			}
			fmt.Println(user.Username)
			// 将模版信息以html的方式返回
			return ctx.HTML(http.StatusOK, buf.String())
		} else {
			// 加载模版信息
			tpl := template.Must(template.ParseFiles("./login.html"))
			ctx.Logger().Info("this is login page")
			// 用户没有登陆则从session中获取用户的登陆信息
			sess := getCookieSession(ctx)
			if flashes := sess.Flashes("username"); len(flashes) > 0 {
				data["username"] = flashes[0]
			}
			sess.Save(ctx.Request(), ctx.Response())
			// 将模版以及数据信息写入到缓冲区中
			var buf bytes.Buffer
			err := tpl.Execute(&buf, data)
			if err != nil {
				return err
			}
			// 将模版信息以html的方式返回
			return ctx.HTML(http.StatusOK, buf.String())
		}

	})

	e.POST("/upload",func(ctx echo.Context)error {
		status := ""
		if user, ok := ctx.Get("user").(*User); ok {
			// 用户已经登陆则从context中获取用户的信息
			if user.Username != ""{
				// 通过echo.Contxt实例的FormFile函数获取客户端上传的单个文件
				file,err:=ctx.FormFile("filename") //filename要与前端对应上
				if err!=nil{
					return err
				}
				// 先打开文件源
				src,err:=file.Open()
				if err!=nil{
					return err
				}
				defer src.Close()
				// 下面创建保存路径文件 file.Filename 即上传文件的名字 创建upload文件夹
				dst,err:=os.Create("upload/"+user.Username+".docx")
				if err !=nil {
					return err
				}
				defer dst.Close()
				// 下面将源拷贝到目标文件
				if _,err=io.Copy(dst,src);err !=nil{
					return err
				}
				// Read from docx file
				r, err := ReadDocxFile("upload/"+user.Username+".docx")
				if err != nil {
					status = "呵呵别玩了，这个文件不是在这里下载的，识别不到！"
				} else {
				docx2 := r.Editable()
				textreadout, err := docx2.readCustom2(r.zipReader.files())
				fmt.Printf("recustom2 %s", textreadout)
				//r.customs若为空，无err返回，不好处理
				//docx1 := r.customs

					//docx1: <?xml version="1.0" encoding="UTF-8" standalone="yes"?>
					//<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/custom-properties" xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes"><property fmtid="{D5CDD505-2E9C-101B-9397-08002B2CF9AE}" pid="2" name="mzj"><vt:lpwstr>2052-0.0.0.0</vt:lpwstr></property></Properties>
					rename, err := regexp.Compile("name=\"(.*?)\"");
					if err != nil{
						fmt.Println(err)
						status = "呵呵，别玩我了，不是在这里OA下载的文件，识别不到"
					} else {
						//one := rename.Find([]byte(docx1));
						sub := rename.FindSubmatch([]byte(textreadout));
						//fmt.Println(string(one));
						fmt.Printf("find sub %s", string(sub[1]))
						//status = string(one)
						status = fmt.Sprintf("<h1><img src=\"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAAABmJLR0QA/wD/AP+gvaeTAAAGsklEQVRoge2YW0xb9x3Hv+fmGzZgGwMJ4LtzM4FRmgttJ61pNbXdlqVSpHXapiqTquZl0V667aWTOm0v1bTuVvVhD9G6bmn60rvU2xoWRWlSJU2AOAnE9YVggs3BBttwjs/t3wcGBWwDcYz74s8TOofzO9/v+f7//9//b6BOnTp16tSpHGqz/9jXN/CISshThEDMzFC/mpj4TNhKYXf7PnazhRWV/OzggZ6nCwWJDA3fNE9M4Ni9y63e+9ZNxOkNPmRttp/M5eY8qqIw1ZVaOQTkxHj0+t9WXiubSJc3uM/cYPn0x8dPcN39B6DT6Tb9otP/+CuohWlQKz7THT6HE7/7cwWyV/PuqX/i3ddfta+9Tpd7wNps/+PRY8e5+wa+fVcmAKD/wYcxmcpiOiNgOiMglV5Ak721Atmbp1wilCjMH+zdP1BRUf+evfjlH/5SuaoKKJmI07nXozcYaUtTc03F3AulE2HUfr3BSL328ks1llOeRDyGX79YPuWSRnSc8eDe/gPMw48d2TJhd8vgB2+ve7+kEb1efzjD8zj78ftbImqJncEeAMBoaBgGoxH7HzqEoUvn0Xv/A/j83KcQha97YGI8sm6tkkZEcb4t2LcP9pb2Ksouz56e+wEAfDKJji4f+GQS3kBwtSZhAQCQz2YBaIW1NUpOdophXn3j5CvCRl+hlkiigGcOH8KFwU9mVJp+Y+39sp3d6d39qI7Vv3nkJ8+Yt3e6MXTxDOYyMxAEAUajEYqqQtNU6Djd8rW+Bw7B0mhDPpvFpfP/hbnRWrI2IRooqmwLK0noysX52Vn+2Xjk+r9L3V93i+JydT9pbm7617Ff/KZh8P3TIHIeNE2hxW7DrXAMfp8LM+lZAAT5eQlP/Og4YuGb5IO3/yMoksxRFPVyqWEAgscJhQQFDG/WCAGG45HQqcU/i1l30xiPX3szYOqfTMS/DPh392H4wkdwuTqLTJiMRjQ6XFAVCR++dUpUJelBAM9rwOXxaOi1og/kCT6mgfrt7ci1S5s1shEb56uSM3xyCo727aBptshEi92G0bEoPDu68eXYdTAsczYWC10Fwd8pkBc6Onat2he53W4DAL+eFkeqZWJTRjRgeTynZ+cQ8LuXTbS3OSAIAqw2K2wt2zA1EZNkSfwIAOLx0BkAJ1kdc97j6f7u1xVNvaCo0XA4XDzkttIIy9AD9tZW8MlJdHW2g5/JYCmJW+E4LBYzdAYTAGByPCZoCi4vPRuPhn6vAc9pIC86Pd0/BQBC0Y9QRDtXTRPABnPE6+1vkmV557YODy7+7z1IsgyWodHe5kAul4fH3YnJO0loKgdNI5hN80aWFa6srHE7GnoHwDsA4PP1tMqaepwG9YNqG1k3EU0Tj3S6vQWWY5EYj4Fl6FVJ8DMZiKIEl38n0nwSDMtOh8Ph7No6fr+/0ent/r6iqYM0RV6JRkND1TaybiI6o/Hne3r2WeazszAZOLTYbRgLxxDwuTCV5EEIAcdx6HAHMDkeR0EqWF2eYNHyKKvIArgKinouFgltyb6nrBG3+1vNiiLt9wR24cr5T9DWutg7Av9fsQghcLTYcHV4FE1WBy5/dq7A0NTzkfC1P22F0I0onwgtPdnl2iVzOr1hOplAQcgvL7tLJsbCMdgdDgDAndtRceVErzVl54hBb3o22Lffkk7dAZ9KLS+7K034PF0wNTQuT3S9XvqiluJXUtKI39/nkFW11+3bjRtDFxEMBsDPZJZNLA2xxGQSTVbH8kQfHR3N1drAEiWNKJp81OPbqbAci7lMCrlcfpWJpSEmCCK8u7oxPZUARdPfWBpAGSM6veGpQLDXPJuexsJ8DoQQbGt3IBIdR8DvBsex0Os4GE0NaLI6MJWIS7IkDtZa/EqKjHR2DhgVWT7g9OzA8OdnYTab4GixIZ9fwI6AB/n5xQOOxWKGrKgAgMTtmEBUpWobwEooMsIwc9+xO9pEvcEAYX4Orq7tEAQRBoMeU0keoihCUVSM3orC0mSDphHMzfBGhpGvfhMGligyQiiqZ1un0wQsnsoWu3jDqj3WWDgGh90Kp3cH0tNTZTt6LSmeIxQy+eycDAAFRQVN0wjdCCOVSkMQChgauQkdx0GjDfDtvg83Ri5Lqqa9XnPlayhqiBzNvBWL3HopMR7B0adPrPtwZOw6Rr64sCBq6r3/qHuPlDzqurzd3+NY9nSX20+1d7hMFL06OFEQlEQ8vJDmU3lRln5YzZNepZQ9s/t8Pa2Kph4G4KaoouRmCagRs5F8HAqFpK2VWKdOnTp16mw9XwEpfPgjEdDEoQAAAABJRU5ErkJggg==\"/>泄密者 %s</h1>", string(sub[1]))
					}
				}

			}
		} else {
			//未登录直接返回
			return ctx.Redirect(http.StatusSeeOther, "/?msg=没有登录，不让查")
		}
		return ctx.HTML(http.StatusOK,status)
	})
}

// getCookieSession 获取当前登陆context中的session信息
func getCookieSession(ctx echo.Context) *sessions.Session {
	session, _ := cookieStore.Get(ctx.Request(), "request-scope")
	return session
}