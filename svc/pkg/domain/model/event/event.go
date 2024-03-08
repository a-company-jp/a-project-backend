package event

import (
	"a-project-backend/svc/pkg/domain/model/pkg_time"
	"a-project-backend/svc/pkg/domain/model/user"
)

type ID string

type Event struct {
	ID         ID
	UserID     user.ID
	Title      string
	Content    string
	ImageID    ImageID
	BeginDate  pkg_time.DateOnly
	FinishDate pkg_time.DateOnly
}

type ImageID string
