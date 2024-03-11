package command

import "a-project-backend/svc/pkg/domain/model/user"

type User interface {
	Create(*user.User) error
	Update(user.User) error
}
