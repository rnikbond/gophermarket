package postgres

const (
	queryCheckUsername = "SELECT id FROM users WHERE username = $1"
	queryCheckUserAuth = "SELECT id FROM users WHERE username = $1 AND password_hash = $2"
	queryCreateUse     = "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
)
