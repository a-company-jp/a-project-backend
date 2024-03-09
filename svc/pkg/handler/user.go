package handler

import (
	"a-project-backend/svc/pkg/domain/model/exception"
	"a-project-backend/svc/pkg/domain/model/user"
	"a-project-backend/svc/pkg/domain/query"
	"errors"
	"github.com/gin-gonic/gin"
)

type User struct {
	userQ query.User
}

func NewUser(userQ query.User) User {
	return User{
		userQ: userQ,
	}
}

func (u User) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := user.ID(c.Param("user_id"))
		_, err := u.userQ.GetUserByID(userID)
		if err != nil {
			if errors.Is(err, exception.ErrNotFound) {
				c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
			} else {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			}
		}
		// TODO: create response using protobuf struct
		//resp := pb_out.UserInfoResponse {}
		//c.Data()
	}
}
