package command

type CreateCmd struct {
	userID     string
	title      string
	content    string
	imageID    string
	beginDate  string
	finishDate string
}

type UpdateCmd struct {
	milestoneID string
	title       string
	content     string
	imageID     string
	beginDate   string
	finishDate  string
}
