package reader

import (
	"a-project-backend/db/query"
	"a-project-backend/svc/pkg/domain/model/milestone"
	"a-project-backend/svc/pkg/domain/model/pkg_time"
	"a-project-backend/svc/pkg/domain/model/user"
	"fmt"
	"gorm.io/gorm"
)

type Milestone struct {
	db *query.Query
}

func NewMilestone(db *gorm.DB) *Milestone {
	return &Milestone{db: query.Use(db)}
}

func (m Milestone) GetByID(milestoneID milestone.ID) (milestone.Milestone, error) {
	found, err := m.db.Milestone.Select(query.Milestone.MilestoneID.Eq(string(milestoneID))).First()
	if err != nil {
		return milestone.Milestone{}, fmt.Errorf("failed to get milestone by id: %w", err)
	}
	return milestone.Milestone{
		ID:         milestone.ID(found.MilestoneID),
		UserID:     user.ID(found.UserID),
		Title:      found.Title,
		Content:    found.Content,
		ImageID:    milestone.ImageID(found.ImageHash),
		BeginDate:  pkg_time.DateOnly(found.BeginDate),
		FinishDate: pkg_time.DateOnly(found.FinishDate),
	}, nil
}

func (m Milestone) GetByUserID(userID user.ID) ([]milestone.Milestone, error) {
	found, err := m.db.Milestone.Select(query.Milestone.UserID.Eq(string(userID))).Find()
	if err != nil {
		return nil, fmt.Errorf("failed to get milestone by user id: %w", err)
	}
	milestones := make([]milestone.Milestone, len(found))
	for i, m := range found {
		milestones[i] = milestone.Milestone{
			ID:         milestone.ID(m.MilestoneID),
			UserID:     user.ID(m.UserID),
			Title:      m.Title,
			Content:    m.Content,
			ImageID:    milestone.ImageID(m.ImageHash),
			BeginDate:  pkg_time.DateOnly(m.BeginDate),
			FinishDate: pkg_time.DateOnly(m.FinishDate),
		}
	}
	return milestones, nil
}
