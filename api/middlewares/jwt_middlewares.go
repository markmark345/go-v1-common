package middlewares

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markmark345/go-v1-common/api/contexts"
	"github.com/markmark345/go-v1-common/api/responses"
)

type JwtClamis struct {
	UuId  string `json:"uuid"`
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
	jwt.RegisteredClaims
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			responses.Failure(c, "missing_authroization_header", errors.New("missing authorization header"))
			c.Abort()
			return
		}

		accessToken := authorization
		if strings.HasPrefix(authorization, "Bearer ") {
			accessToken = strings.TrimPrefix(authorization, "Bearer ")
		}

		ctx := c.Request.Context()

		var p jwt.Parser
		jwtToken, _, err := p.ParseUnverified(accessToken, &JwtClamis{})
		if err != nil {
			responses.Failure(c, "invalid_access_token", err)
			c.Abort()
			return
		}

		if jwtToken.Valid {
			jwtClaims := jwtToken.Claims.(*JwtClamis)
			ctx = contexts.SetUseeId(ctx, jwtClaims.UuId)

			req := c.Request.WithContext(ctx)
			c.Request = req

			c.Next()
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			responses.Failure(c, "expired_token", err)
			c.Abort()
			return
		}

	}
}
