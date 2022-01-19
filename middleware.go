package main

import "github.com/labstack/echo/v4"

// AutoLogin 自定义中间件，如果上次记住了则自动登陆
func AutoLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// 获取用户cookie信息
		cookie, err := ctx.Cookie("username")
		// 当前用户已经登陆
		if err == nil && cookie.Value != "" {
			// 将user信息放在context中，即记住用户信息
			user := &User{Username: cookie.Value}
			ctx.Set("user", user)
		}
		// 返回中间件
		return next(ctx)
	}
}
