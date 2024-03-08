package middleware

import (
	"a-project-backend/svc/pkg/domain/model/exception"
	"a-project-backend/svc/pkg/domain/model/user"
	"a-project-backend/svc/pkg/uc"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizedUserIDField = "AuthorizedUserID"
)

type auth struct {
	loginUC uc.LoginUseCase
}

func NewAuth() auth {
	return auth{
		loginUC: uc.NewLoginUseCase(),
	}
}

func (a auth) VerifyUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt, err := getJWTFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}
		result, err := a.loginUC.Do(c, uc.LoginInput{JWT: user.JWT(jwt)})
		if result == nil || err != nil {
			c.AbortWithError(401, err)
			return
		}
		c.Set(AuthorizedUserIDField, result.UserID)
		c.Next()
	}
}

func getJWTFromHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")

	if len(header) < 8 || header[:7] != "Bearer " {
		return "", exception.ErrorInvalidHeader
	}
	return header[7:], nil
}
