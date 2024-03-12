package writer

import (
	"a-project-backend/db/model"
	"a-project-backend/db/query"
	"a-project-backend/svc/pkg/domain/model/user"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	db *query.Query
}

func NewUser(db *gorm.DB) *User {
	return &User{db: query.Use(db)}
}

func (w User) Create(u *user.User) error {
	if u == nil {
		return errors.New("user is nil")
	}
	if u.UserId != "" {
		return errors.New("user id is not empty")
	}
	newID := uuid.New().String()
	tags := make([]model.Tag, len(u.Tags))
	for i, tag := range u.Tags {
		tags[i] = model.Tag{
			TagID:   string(tag.ID),
			TagName: tag.Name,
		}
	}
	dbModel := model.User{
		UserID:        newID,
		FirebaseUID:   u.FirebaseUID,
		Username:      u.Username,
		Firstname:     u.Firstname,
		Lastname:      u.Lastname,
		FirstnameKana: u.FirstnameKana,
		LastnameKana:  u.LastnameKana,
		StatusMessage: u.StatusMessage,
		IconImageHash: (*string)(u.IconImageHash),
		Tags:          tags,
	}
	if err := w.db.User.Create(&dbModel); err != nil {
		return err
	}
	u.UserId = user.ID(newID)
	return nil
}

func (w User) Update(u model.User) error {
	if u.UserID == "" {
		return errors.New("user id is empty")
	}
	tags := make([]model.Tag, len(u.Tags))
	for i, tag := range u.Tags {
		tags[i] = model.Tag{
			TagID:   tag.TagID,
			TagName: tag.TagName,
		}
	}
	dbModel := model.User{
		UserID:        u.UserID,
		FirebaseUID:   u.FirebaseUID,
		Username:      u.Username,
		Firstname:     u.Firstname,
		Lastname:      u.Lastname,
		FirstnameKana: u.FirstnameKana,
		LastnameKana:  u.LastnameKana,
		StatusMessage: u.StatusMessage,
		IconImageHash: u.IconImageHash,
		Tags:          tags,
	}
	_, err := w.db.User.Updates(&dbModel)
	if err != nil {
		return err
	}
	return nil
}
