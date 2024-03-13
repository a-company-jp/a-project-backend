package handler

import (
	"a-project-backend/svc/pkg/domain/command"
	"a-project-backend/svc/pkg/domain/model/exception"
	"a-project-backend/svc/pkg/domain/model/user"
	"a-project-backend/svc/pkg/domain/query"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/openhacku-team-a/a-project-frontend/proto/golang/pb_out"
)

type User struct {
	userQ      query.User
	milestoneQ query.Milestone
	userC      command.User
}

func NewUser(userQ query.User, msQ query.Milestone) User {
	return User{
		userQ:      userQ,
		milestoneQ: msQ,
	}
}

func (h User) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := user.ID(c.Param("user_id"))
		u, err := h.userQ.GetByID(userID)
		if err != nil {
			if errors.Is(err, exception.ErrNotFound) {
				c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
			} else {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			}
		}
		ms, err := h.milestoneQ.GetByUserID(u.UserId)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": fmt.Errorf("failed to get milestone by user id: %v", err)})
		}

		tags := make([]*pb_out.Tag, len(u.Tags))
		for i, t := range u.Tags {
			tags[i] = &pb_out.Tag{
				TagId:   string(t.ID),
				TagName: t.Name,
			}
		}

		var milestones []*pb_out.Milestone
		for _, m := range ms {
			milestones = append(milestones, &pb_out.Milestone{
				MilestoneId: string(m.ID),
				UserId:      string(m.UserID),
				Title:       m.Title,
				Content:     m.Content,
				ImageHash:   string(m.ImageID),
				BeginDate:   m.BeginDate.String(),
				FinishDate:  m.FinishDate.String(),
			})
		}

		resp := pb_out.UserInfoResponse{
			UserData: &pb_out.UserData{
				UserId:        string(u.UserId),
				Username:      u.Username,
				Firstname:     u.Firstname,
				Lastname:      u.Lastname,
				FirstnameKana: u.FirstnameKana,
				LastnameKana:  u.LastnameKana,
				StatusMessage: u.StatusMessage,
				Tag:           tags,
				IconImageHash: string(*u.IconImageHash),
			},
			Milestones: milestones,
		}
		respData, err := proto.Marshal(&resp)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Data(200, "application/octet-stream", respData)
	}
}

//func (h User) UpdateUserInfo() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		data, err := io.ReadAll(c.Request.Body)
//		if err != nil {
//			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
//			return
//		}
//
//		var req pb_out.UserInfoUpdateRequest
//		if err := proto.Unmarshal(data, &req); err != nil {
//			c.AbortWithStatusJSON(500, gin.H{"error": fmt.Sprintf("failed to unmarshal, err: %v\n", err)})
//			return
//		}
//
//		target := user.User{
//			UserId:        user.ID(req.UserData.UserId),
//			FirebaseUID:   req.UserData.,
//			Username:      "",
//			Firstname:     "",
//			Lastname:      "",
//			FirstnameKana: "",
//			LastnameKana:  "",
//			StatusMessage: "",
//			Tags:          nil,
//			IconImageHash: nil,
//		}
//		h.userC.Update()
//	}
//}
