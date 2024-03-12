package command

import "a-project-backend/svc/pkg/domain/model/milestone"

type Milestone interface {
	Create(*milestone.Milestone) error
	Update(milestone.Milestone) error
	Delete([]milestone.ID) error
}
