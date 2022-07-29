package repository

type Authorization interface {
}

type Registration interface {
}

type Repository struct {
	Authorization
	Registration
}

func NewRepository() *Repository {

	return &Repository{}
}
