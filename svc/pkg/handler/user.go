package handler

import (
	"a-project-backend/gen/gQuery"
	"a-project-backend/pkg/config"
	"a-project-backend/svc/pkg/domain/model/exception"
	"a-project-backend/svc/pkg/domain/model/pkg_time"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/openhacku-team-a/a-project-frontend/proto/golang/pb_out"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
	q  *gQuery.Query
}

const (
	GCSUserIconFolder       = "user_icon"
	GCSMilestoneImageFolder = "milestone_img"
)

func NewUser(db *gorm.DB) User {
	return User{
		db: db,
		q:  gQuery.Use(db),
	}
}

func (h User) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		u, err := h.q.User.Where(gQuery.User.UserID.Eq(userID)).First()
		if err != nil {
			if errors.Is(err, exception.ErrNotFound) {
				c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
			} else {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			}
		}

		ms, err := h.q.Milestone.Where(gQuery.Milestone.UserID.Eq(userID)).Find()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		tags := make([]*pb_out.Tag, len(u.Tags))
		for i, t := range u.Tags {
			tags[i] = &pb_out.Tag{
				TagId:   t.TagID,
				TagName: t.TagName,
			}
		}

		conf := config.Get()

		var milestones []*pb_out.Milestone
		for _, m := range ms {
			imgUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", conf.Application.GCS.BucketName, GCSMilestoneImageFolder, m.ImageHash)
			milestones = append(milestones, &pb_out.Milestone{
				MilestoneId: m.MilestoneID,
				UserId:      m.UserID,
				Title:       m.Title,
				Content:     m.Content,
				ImageUrl:    &imgUrl,
				BeginDate:   pkg_time.DateOnly(m.BeginDate).String(),
				FinishDate:  pkg_time.DateOnly(m.FinishDate).String(),
			})
		}

		iconUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", conf.Application.GCS.BucketName, GCSUserIconFolder, u.IconImageHash)

		resp := pb_out.UserInfoResponse{
			UserData: &pb_out.UserData{
				UserId:        u.UserID,
				Username:      u.Username,
				Firstname:     u.Firstname,
				Lastname:      u.Lastname,
				FirstnameKana: u.FirstnameKana,
				LastnameKana:  u.LastnameKana,
				StatusMessage: u.StatusMessage,
				Tag:           tags,
				IconUrl:       &iconUrl,
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
