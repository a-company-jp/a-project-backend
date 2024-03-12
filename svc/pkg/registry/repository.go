package registry

import (
	"a-project-backend/svc/pkg/domain/command"
	"a-project-backend/svc/pkg/domain/query"
	"a-project-backend/svc/pkg/infra/reader"
	"a-project-backend/svc/pkg/infra/writer"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r Repository) NewUserQuery() query.User {
	return reader.NewUser(r.db)
}

func (r Repository) NewMilestoneQuery() query.Milestone {
	return reader.NewMilestone(r.db)
}

func (r Repository) NewUserCommand() command.User {
	return writer.NewUser(r.db)
}

func (r Repository) NewMilestoneCommand() command.Milestone {
	return writer.NewMilestone(r.db)
}
