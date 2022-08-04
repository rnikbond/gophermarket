package postgres

// User
const (
	queryCreateUser = `INSERT INTO users (username, password_hash) 
                       VALUES ($1, $2)`

	queryGetUserID = `SELECT id
                      FROM users 
                      WHERE username = $1 AND password_hash = $2`

	queryGetUserIDByName = `SELECT id
                            FROM users 
                             WHERE username = $1`
)

// Order
const (
	queryCreateOrder = `INSERT INTO orders (user_id, number,status,created_at) 
                        VALUES($1,$2,$3,$4)`

	queryOrderUserID = `SELECT user_id 
                        FROM orders 
                        WHERE number = $1`
)
