package handler

import (
	"a-project-backend/gen/gModel"
	"a-project-backend/gen/gQuery"
	"a-project-backend/svc/pkg/domain/model/pkg_time"
	"a-project-backend/svc/pkg/middleware"
	"io"

	"github.com/google/uuid"
	"gorm.io/gorm"

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
			MilestoneID: uuid.New().String(),
			UserID:      c.GetString(middleware.AuthorizedUserIDField),
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
		c.ProtoBuf(201, respData)
	}
}

// UpdateMileStone マイルストーンの更新
func (m MileStone) UpdateMileStone() gin.HandlerFunc {
	return func(c *gin.Context) {
		mileStoneId := c.Param("milestone_id")
		if mileStoneId == "" {
			c.AbortWithStatusJSON(500, gin.H{"error": "id should not be null"})
		}

		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		var milestoneUpdateRequest pb_out.MilestoneUpdateRequest
		err = proto.Unmarshal(data, &milestoneUpdateRequest)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		_, err = m.q.Milestone.WithContext(c).
			Where(m.q.Milestone.MilestoneID.Eq(milestoneUpdateRequest.Milestone.MilestoneId),
				m.q.Milestone.UserID.Eq(c.GetString(middleware.AuthorizedUserIDField)),
			).Updates(
			map[string]interface{}{
				"title":       milestoneUpdateRequest.Milestone.Title,
				"content":     milestoneUpdateRequest.Milestone.Content,
				"begin_date":  milestoneUpdateRequest.Milestone.BeginDate,
				"finish_date": milestoneUpdateRequest.Milestone.FinishDate,
			},
		)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		c.AbortWithStatus(200)
	}
}

// DeleteMileStone マイルストーンの削除
func (m MileStone) DeleteMileStone() gin.HandlerFunc {
	return func(c *gin.Context) {
		mileStoneId := c.Param("milestone_id")
		if mileStoneId == "" {
			c.AbortWithStatusJSON(500, gin.H{"error": "id should not be null"})
		}
		_, err := m.q.Milestone.WithContext(c).Where(m.q.Milestone.MilestoneID.Eq(mileStoneId)).Delete()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		c.AbortWithStatus(200)
	}
}
