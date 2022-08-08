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
	queryCreateOrder = `INSERT INTO orders (user_id, number,status, uploaded_at) 
                        VALUES($1,$2,$3,$4)`

	queryUpdateOrder = `UPDATE orders
                        SET status = $1
                        WHERE number = $2`

	queryUpdateAccrual = `UPDATE orders
                          SET accrual = $1
                          WHERE number = $2`

	queryOrderUserID = `SELECT user_id 
                        FROM orders 
                        WHERE number = $1`

	queryOrdersByStatuses = `SELECT number, status
                             FROM orders
                             WHERE status = ANY($1)`
)

// Accruals
const (
	queryUserAccruals = `SELECT SUM(accrual)
                         FROM orders
                         WHERE user_id = $1`
)
