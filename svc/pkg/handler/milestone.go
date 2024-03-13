package handler

import (
	"a-project-backend/gen/gModel"
	"a-project-backend/gen/gQuery"
	"a-project-backend/svc/pkg/domain/model/pkg_time"
	"gorm.io/gorm"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/openhacku-team-a/a-project-frontend/proto/golang/pb_out"
)

type MileStone struct {
	db *gorm.DB
	q  *gQuery.Query
}

func NewMileStone(db *gorm.DB) MileStone {
	return MileStone{
		db: db,
		q:  gQuery.Use(db),
	}
}

// PostMileStone マイルストーンの作成
func (m MileStone) PostMileStone() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		var milestoneCreateRequest pb_out.MilestoneCreateRequest
		err = proto.Unmarshal(data, &milestoneCreateRequest)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		if milestoneCreateRequest.Milestone.MilestoneId != "" {
			c.AbortWithStatusJSON(500, gin.H{"error": "milestone id is not empty"})
		}

		beginDate, err := pkg_time.FromString(milestoneCreateRequest.Milestone.BeginDate)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		finishDate, err := pkg_time.FromString(milestoneCreateRequest.Milestone.FinishDate)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		target := gModel.Milestone{
			MilestoneID: milestoneCreateRequest.Milestone.MilestoneId,
			UserID:      milestoneCreateRequest.Milestone.UserId,
			Title:       milestoneCreateRequest.Milestone.Title,
			Content:     milestoneCreateRequest.Milestone.Content,
			ImageHash:   "",
			BeginDate:   int32(beginDate),
			FinishDate:  int32(finishDate),
		}
		if err = m.q.Milestone.Create(&target); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		resp := pb_out.MilestoneCreateResponse{
			Milestone: &pb_out.Milestone{
				UserId:      target.UserID,
				MilestoneId: target.MilestoneID,
				Title:       target.Title,
				Content:     target.Content,
				ImageUrl:    nil,
				BeginDate:   pkg_time.DateOnly(target.BeginDate).String(),
				FinishDate:  pkg_time.DateOnly(target.FinishDate).String(),
			},
		}

		respData, err := proto.Marshal(&resp)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(201, "application/octet-stream", respData)
	}
}
