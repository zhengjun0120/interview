package middleware

import (
	"ai_interview/biz/entity"
	"ai_interview/pkg/response"
	"ai_interview/util"
	"github.com/gin-gonic/gin"
	"strings"
)

func JWTAuth() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		ctx := gCtx.Request.Context()

		authHeader := gCtx.GetHeader("Authorization")
		if authHeader == "" {
			r := response.NewResponse(gCtx)
			r.ErrorWithStatus(response.TOKEN_IS_EXPIRED, 401)
			gCtx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			r := response.NewResponse(gCtx)
			r.ErrorWithStatus(response.TOKEN_IS_EXPIRED, 401)
			gCtx.Abort()
			return
		}

		tokenString := parts[1]

		userID, err := util.ParseToken(tokenString)
		if userID == "" || err != nil {
			r := response.NewResponse(gCtx)
			r.ErrorWithStatus(response.TOKEN_IS_EXPIRED, 401)
			gCtx.Abort()
			return
		}

		// 将userID保存到上下文
		ctx = entity.WithUserID(ctx, userID)
		// 将上下文保存到gin.Context
		gCtx.Request = gCtx.Request.WithContext(ctx)

		gCtx.Next()
	}
}
