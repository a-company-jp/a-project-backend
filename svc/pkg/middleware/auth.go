package middleware

import (
	"a-project-backend/gen/gModel"
	"a-project-backend/gen/gQuery"
	"a-project-backend/svc/pkg/domain/model/exception"
	"a-project-backend/svc/pkg/uc"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	AuthorizedUserIDField = "AuthorizedUserID"
)

type auth struct {
	certString map[string]string
	loginUC    uc.LoginUseCase
	db         *gorm.DB
	q          *gQuery.Query
}

func NewAuth(db *gorm.DB) auth {
	resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")
	if err != nil {
		log.Fatalf("Failed to make a request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read the response body: %v", err)
	}

	var result map[string]string

	if err = json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Failed to json unmarshal: %v", err)
	}
	return auth{
		certString: result,
		loginUC:    uc.NewLoginUseCase(),
		db:         db,
		q:          gQuery.Use(db),
	}
}

func (a auth) VerifyUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt, err := getJWTFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}
		result, err := a.loginUC.Do(c, uc.LoginInput{JWT: jwt})
		if result == nil || err != nil {
			c.AbortWithError(401, err)
			return
		}

		// Userの存在チェック
		if result.UserID == "" {
			c.AbortWithStatusJSON(500, "user_id is null")
			return
		}
		_, err = a.q.User.WithContext(c).Where(gQuery.User.FirebaseUID.Eq(result.UserID)).First()
		if err != nil {
			if errors.Is(err, exception.ErrNotFound) {
				// 存在しなければ、作成する
				err := a.q.User.WithContext(c).Create(&gModel.User{
					UserID:      uuid.New().String(),
					FirebaseUID: result.UserID,
				})
				if err != nil {
					c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
					return
				}
			}
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
