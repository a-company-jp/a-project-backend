package query

import "a-project-backend/svc/pkg/domain/model/user"

type User interface {
	GetUserByID(userID user.ID) (user.User, error)
}
