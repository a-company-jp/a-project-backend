package service

import (
	"errors"

	"a-project-backend/svc/pkg/domain/command"
	"a-project-backend/svc/pkg/domain/model/milestone"
)

type IMilestoneService interface {
	// basic CRUD
	Create(command.CreateCmd) error
	GetByID(milestoneID string) (milestone.Milestone, error)
	Update(command.UpdateCmd) error
	Delete(milestoneID string) error
	// query
	GetByUserID(userID string) ([]milestone.Milestone, error)
	GetByTitle(title string) ([]milestone.Milestone, error)
	GetByUserName(userName string) (milestone.Milestone, error)
	Pagination() ([10]milestone.Milestone, error)
}

type MilestoneService struct{}

func NewMilestoneService() *MilestoneService {
	return &MilestoneService{}
}

func (s *MilestoneService) Create(userID string, title string, content string, beginDate string, finishDate string) error {
	return errors.New("not implemented")
}

func (s *MilestoneService) GetByID(milestoneID string) (milestone.Milestone, error) {
	return milestone.Milestone{}, errors.New("not implemented")
}

func (s *MilestoneService) GetByUesrID(userID string) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (s *MilestoneService) Pagenation() ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (s *MilestoneService) GetByUserName(userName string) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}

func (s *MilestoneService) GetByTitle(title string) ([]milestone.Milestone, error) {
	return []milestone.Milestone{}, errors.New("not implemented")
}
