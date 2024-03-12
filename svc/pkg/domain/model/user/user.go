package user

type User struct {
	UserId        ID
	Username      string
	Firstname     string
	Lastname      string
	FirstnameKana string
	LastnameKana  string
	StatusMessage string
	Tags          []Tag
	IconImageHash *IconID
}

type ID string

type IconID string

type TagID string

type Tag struct {
	ID   TagID
	Name string
}

type JWT string
