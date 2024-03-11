package repository

import (
	"errors"

	"a-project-backend/svc/pkg/domain/model/milestone"
)

type IMilestoneRepo interface {
	SelectByID(milestoneID string) (milestone.Milestone, error)
	SelectsByIDs(milestoneIDs []string) ([]milestone.Milestone, error)
	SelectsByUserID(userID string) ([]milestone.Milestone, error)
	SelectsByUserIDs(userIDs []string) ([]milestone.Milestone, error)
	SelectByUserName(userName string) ([]milestone.Milestone, error)
	SelectByTitle(title string) ([]milestone.Milestone, error)
	INSERT(milestone milestone.Milestone) error
	Update(milestone milestone.Milestone) error
	DELETE(milestoneIDs []string) error
}

type MilestoneRepo struct{}

func NewMilestoneRepo() *MilestoneRepo {
	return &MilestoneRepo{}
}

func (r *MilestoneRepo) SelectByID(milestoneID string) (milestone.Milestone, error) {
	return milestone.Milestone{}, errors.New("not implemented")
}

func (r *MilestoneRepo) SelectsByIDs(milestoneIDs []string) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (r *MilestoneRepo) SelectsByUserID(userID string) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (r *MilestoneRepo) SelectsByUserIDs(userIDs []string) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (r *MilestoneRepo) SelectByUserName(userName string) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (r *MilestoneRepo) SelectByTitle(title string) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (r *MilestoneRepo) INSERT(milestone milestone.Milestone) error {
	return errors.New("not implemented")
}

func (r *MilestoneRepo) Update(milestone milestone.Milestone) error {
	return errors.New("not implemented")
}

func (r *MilestoneRepo) DELETE(milestoneIDs []string) error {
	return errors.New("not implemented")
}

