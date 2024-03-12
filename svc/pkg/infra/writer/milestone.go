package writer

import (
	"a-project-backend/db/model"
	"a-project-backend/db/query"
	"errors"
)

type Milestone struct {
	db *query.Query
}

func NewMilestone(db *query.Query) *Milestone {
	return &Milestone{db: db}
}

func (m Milestone) Create(target *model.Milestone) error {
	if target == nil {
		return errors.New("milestone is nil")
	}
	if target.MilestoneID != "" {
		return errors.New("milestone id is not empty")
	}
	dbModel := model.Milestone{
		MilestoneID: target.MilestoneID,
		UserID:      target.UserID,
		Title:       target.Title,
		Content:     target.Content,
		ImageHash:   target.ImageHash,
		BeginDate:   target.BeginDate,
		FinishDate:  target.FinishDate,
	}
	if err := m.db.Milestone.Create(&dbModel); err != nil {
		return err
	}
	return nil
}

func (m Milestone) Update(target model.Milestone) error {
	if target.MilestoneID == "" {
		return errors.New("milestone id is empty")
	}
	dbModel := model.Milestone{
		MilestoneID: target.MilestoneID,
		UserID:      target.UserID,
		Title:       target.Title,
		Content:     target.Content,
		ImageHash:   target.ImageHash,
		BeginDate:   target.BeginDate,
		FinishDate:  target.FinishDate,
	}
	if _, err := m.db.Milestone.Updates(&dbModel); err != nil {
		return err
	}
	return nil
}
