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
