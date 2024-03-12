package handler

import (
	"a-project-backend/svc/pkg/domain/command"
	"a-project-backend/svc/pkg/domain/model/milestone"
	"a-project-backend/svc/pkg/domain/model/pkg_time"
	"a-project-backend/svc/pkg/domain/model/user"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/openhacku-team-a/a-project-frontend/proto/golang/pb_out"
)

type MileStone struct {
	mileStoneCommand command.Milestone
}

func NewMileStone(mileStoneCommand command.Milestone) MileStone {
	return MileStone{
		mileStoneCommand: mileStoneCommand,
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
		target := milestone.Milestone{
			ID:         milestone.ID(""),
			UserID:     user.ID(milestoneCreateRequest.Milestone.UserId),
			Title:      milestoneCreateRequest.Milestone.Title,
			Content:    milestoneCreateRequest.Milestone.Content,
			ImageID:    milestone.ImageID(milestoneCreateRequest.Milestone.ImageHash),
			BeginDate:  beginDate,
			FinishDate: finishDate,
		}
		err = m.mileStoneCommand.Create(&target)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		resp := pb_out.MilestoneCreateResponse{
			Milestone: &pb_out.Milestone{
				UserId:      string(target.UserID),
				MilestoneId: string(target.ID),
				Title:       target.Title,
				Content:     target.Content,
				ImageHash:   string(target.ImageID),
				BeginDate:   target.BeginDate.String(),
				FinishDate:  target.FinishDate.String(),
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

// UpdateMileStone マイルストーンの更新
func (m MileStone) UpdateMileStone() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		var milestoneUpdateRequest pb_out.MilestoneUpdateRequest
		err = proto.Unmarshal(data, &milestoneUpdateRequest)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		if milestoneUpdateRequest.Milestone.MilestoneId != "" {
			c.AbortWithStatusJSON(500, gin.H{"error": "milestone id is not empty"})
		}

		beginDate, err := pkg_time.FromString(milestoneUpdateRequest.Milestone.BeginDate)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		finishDate, err := pkg_time.FromString(milestoneUpdateRequest.Milestone.FinishDate)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		target := milestone.Milestone{
			ID:         milestone.ID(""),
			UserID:     user.ID(milestoneUpdateRequest.Milestone.UserId),
			Title:      milestoneUpdateRequest.Milestone.Title,
			Content:    milestoneUpdateRequest.Milestone.Content,
			ImageID:    milestone.ImageID(milestoneUpdateRequest.Milestone.ImageHash),
			BeginDate:  beginDate,
			FinishDate: finishDate,
		}
		err = m.mileStoneCommand.Update(target)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		c.Data(200, "application/octet-stream", nil)
	}
}

// DeleteMileStone マイルストーンの削除
func (m MileStone) DeleteMileStone() gin.HandlerFunc {
	return func(c *gin.Context) {
		mileStoneId := milestone.ID(c.Param("id"))
		err := m.mileStoneCommand.Delete([]milestone.ID{mileStoneId})
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		c.Data(200, "application/octet-stream", nil)
	}
}
