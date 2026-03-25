package middleware

import (
	"ai_interview/biz/entity"
	"ai_interview/infra/storage"
	"ai_interview/pkg/response"
	"ai_interview/pkg/zlog"
	"ai_interview/util"
	"github.com/gin-gonic/gin"
	"strings"
)

func JWTAuth() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		ctx := gCtx.Request.Context()

		var authHeader string
		var tokenString string
		authHeader = gCtx.Query("token")
		if authHeader == "" {
			authHeader = gCtx.GetHeader("Authorization")
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

			tokenString = parts[1]
		} else {
			tokenString = authHeader
		}

		userID, err := util.ParseToken(tokenString)
		if userID == "" || err != nil {
			r := response.NewResponse(gCtx)
			r.ErrorWithStatus(response.TOKEN_IS_EXPIRED, 401)
			gCtx.Abort()
			return
		}

		us := storage.GetUserStorage()
		ok, err := us.CheckUser(ctx, userID)
		// 如果用户存在，将userID保存到上下文
		if ok {
			// 将userID保存到上下文
			ctx = entity.WithUserID(ctx, userID)
			// 将上下文保存到gin.Context
			gCtx.Request = gCtx.Request.WithContext(ctx)

			gCtx.Next()
		} else {
			// 如果用户不存在，将talentID保存到上下文
			if err != nil {
				r := response.NewResponse(gCtx)
				r.ErrorWithStatus(response.TOKEN_IS_EXPIRED, 401)
				zlog.Errorf("中间件鉴权，数据库错误: %v", err)
				gCtx.Abort()
				return
			} else {
				ctx = entity.WithTalentID(ctx, userID)
				gCtx.Request = gCtx.Request.WithContext(ctx)
				gCtx.Next()
			}
		}

		//// 将userID保存到上下文
		//ctx = entity.WithUserID(ctx, userID)
		//// 将上下文保存到gin.Context
		//gCtx.Request = gCtx.Request.WithContext(ctx)
		//
		//gCtx.Next()
	}
}
