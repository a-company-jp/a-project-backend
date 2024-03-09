package command

import "a-project-backend/svc/pkg/domain/model/user"

type User interface {
	CreateUser(*user.User) error
	UpdateUser(user.User) error
}
