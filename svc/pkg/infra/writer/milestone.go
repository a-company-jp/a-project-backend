package writer

import (
	"a-project-backend/db/model"
	"a-project-backend/db/query"
	"a-project-backend/svc/pkg/domain/model/milestone"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Milestone struct {
	gormDB *gorm.DB
	db     *query.Query
}

func NewMilestone(db *gorm.DB) *Milestone {
	return &Milestone{db: query.Use(db), gormDB: db}
}

func (m Milestone) Create(target *milestone.Milestone) error {
	if target == nil {
		return errors.New("milestone is nil")
	}
	if target.ID != "" {
		return errors.New("milestone id is not empty")
	}
	dbModel := model.Milestone{
		MilestoneID: string(target.ID),
		UserID:      string(target.UserID),
		Title:       target.Title,
		Content:     target.Content,
		ImageHash:   string(target.ImageID),
		BeginDate:   int32(target.BeginDate),
		FinishDate:  int32(target.FinishDate),
	}
	if err := m.db.Milestone.Create(&dbModel); err != nil {
		return err
	}
	return nil
}

func (m Milestone) Update(target milestone.Milestone) error {
	if target.ID == "" {
		return errors.New("milestone id is empty")
	}
	dbModel := model.Milestone{
		MilestoneID: string(target.ID),
		UserID:      string(target.UserID),
		Title:       target.Title,
		Content:     target.Content,
		ImageHash:   string(target.ImageID),
		BeginDate:   int32(target.BeginDate),
		FinishDate:  int32(target.FinishDate),
	}
	if _, err := m.db.Milestone.Updates(&dbModel); err != nil {
		return err
	}
	return nil
}

func (m Milestone) Delete(milestoneID []milestone.ID) error {
	if len(milestoneID) == 0 {
		return errors.New("milestone id is empty")
	}
	ids := make([]string, len(milestoneID))
	for i, id := range milestoneID {
		ids[i] = string(id)
	}
	if err := m.gormDB.Where("milestone_id IN ?", ids).Delete(&model.Milestone{}); err != nil {
		return fmt.Errorf("failed to delete milestone: %w", err)
	}
	return nil
}
