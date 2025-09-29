/*
Copyright © 2025 lixw
*/
package middleware

import (
	"net/http"
	"strings"

	"github.com/ethanli-dev/go-app-layout/pkg/i18n"
	"github.com/gin-gonic/gin"
)

func I18n() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Next()
			return
		}
		locale := ctx.GetHeader("Accept-Language")
		if locale != "" {
			languages := strings.Split(locale, ",")
			if len(languages) > 0 {
				locale = languages[0]
			}
		}
		ctx.Request = ctx.Request.WithContext(i18n.SetLocale(ctx, locale))
		ctx.Next()
	}
}
