package repository

const (
	queryGetUser   = "SELECT username, password_hash FROM users WHERE username = $1"
	queryCreateUse = "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
)
