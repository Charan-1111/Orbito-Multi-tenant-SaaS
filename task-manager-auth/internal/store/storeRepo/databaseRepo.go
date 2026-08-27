package repository

type Repository interface {
	Close()
	RegisterUser(username, password string) error
}
