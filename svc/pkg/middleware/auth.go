package middleware

import (
	"a-project-backend/gen/gModel"
	"a-project-backend/gen/gQuery"
	"a-project-backend/svc/pkg/domain/model/exception"
	"a-project-backend/svc/pkg/uc"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	AuthorizedUserIDField = "AuthorizedUserID"
	AuthorizedUserField   = "AuthorizedUser"
)

type auth struct {
	loginUC uc.LoginUseCase
	db      *gorm.DB
	q       *gQuery.Query
}

func NewAuth(db *gorm.DB) auth {
	return auth{
		loginUC: uc.NewLoginUseCase(),
		db:      db,
		q:       gQuery.Use(db),
	}
}

func (a auth) VerifyUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt, err := getJWTFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}
		firebaseAuth, err := a.loginUC.Do(c, uc.LoginInput{JWT: jwt})
		if firebaseAuth == nil || err != nil {
			c.AbortWithError(401, err)
			return
		}

		// Userの存在チェック
		if firebaseAuth.UserID == "" {
			c.AbortWithStatusJSON(500, "user_id is null")
			return
		}
		u, err := a.q.User.WithContext(c).Where(a.q.User.FirebaseUID.Eq(firebaseAuth.UserID)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 存在しなければ、作成する
				u = &gModel.User{
					UserID:      uuid.New().String(),
					FirebaseUID: firebaseAuth.UserID,
				}
				err := a.q.User.WithContext(c).Create(u)
				if err != nil {
					c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
					return
				}
			}
		}

		c.Set(AuthorizedUserIDField, u.UserID)
		c.Set(AuthorizedUserField, *u)
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
