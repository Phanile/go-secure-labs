package users_postgres_repository

type UserModel struct {
	Id           int64
	Email        string
	PasswordHash string
	Role         string
}
