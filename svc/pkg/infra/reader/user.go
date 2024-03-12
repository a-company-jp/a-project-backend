package reader

import (
	"a-project-backend/db/query"
	"a-project-backend/svc/pkg/domain/model/user"
	"gorm.io/gorm"
)

type User struct {
	db *query.Query
}

func NewUser(db *gorm.DB) *User {
	return &User{
		db: query.Use(db),
	}
}

func (r User) GetByID(id user.ID) (user.User, error) {
	u, err := r.db.User.Where(query.User.UserID.Eq(string(id))).First()
	tags := make([]user.Tag, len(u.Tags))
	result := user.User{
		UserId:        user.ID(u.UserID),
		Username:      u.Username,
		Firstname:     u.Firstname,
		Lastname:      u.Lastname,
		FirstnameKana: u.FirstnameKana,
		LastnameKana:  u.LastnameKana,
		StatusMessage: u.StatusMessage,
		Tags:          tags,
		IconImageHash: (*user.IconID)(u.IconImageHash),
	}
	return result, err
}
