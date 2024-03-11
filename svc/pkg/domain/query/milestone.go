package query

import (
	"errors"

	"a-project-backend/svc/pkg/domain/model/milestone"
)

type IMilestoneQuery interface {
	GetMilestoneByID(milestoneID string) (milestone.Milestone, error)
	GetMilestonesByIDs(milestoneIDs []string) ([]milestone.Milestone, error)
	GetMilestonesByUserID(userID string) ([]milestone.Milestone, error)
	GetMilestonesByUserIDs(userIDs []string) ([]milestone.Milestone, error)
	SearchMilestonesByUserName(userName string) ([]milestone.Milestone, error)
	SearchMilestonesByTitle(title string) ([]milestone.Milestone, error)
}

func NewMilestoneQuery() *MilestoneQuery {
	return &MilestoneQuery{}
}

type MilestoneQuery struct{}

func (q *MilestoneQuery) GetMilestoneByID(milestoneID milestone.ID) (milestone.Milestone, error) {
	return milestone.Milestone{}, errors.New("not implemented")
}

func (q *MilestoneQuery) GetMilestonesByIDs(milestoneIDs []milestone.ID) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (q *MilestoneQuery) GetMilestonesByUserID(userID milestone.UserID) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (q *MilestoneQuery) GetMilestonesByUserIDs(userIDs []milestone.UserID) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (q *MilestoneQuery) SearchMilestonesByUserName(userName string) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (q *MilestoneQuery) SearchMilestonesByTitle(title string) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}
