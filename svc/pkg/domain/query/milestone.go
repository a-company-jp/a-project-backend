package query

import (
	"a-project-backend/svc/pkg/domain/model/milestone"
	"a-project-backend/svc/pkg/domain/model/user"
)

type Milestone interface {
	GetByID(milestoneID milestone.ID) (milestone.Milestone, error)
	GetByIDs(milestoneIDs []milestone.ID) ([]milestone.Milestone, error)
	GetByUserID(userID user.ID) ([]milestone.Milestone, error)
	GetByUserIDs(userIDs []user.ID) ([]milestone.Milestone, error)
	SearchByTitle(title string) ([]milestone.Milestone, error)
}
