package command

type Milestone interface {
	Create(*Milestone) error
	Update(Milestone) error
	Delete([]string) error
}
