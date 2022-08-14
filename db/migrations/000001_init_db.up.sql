CREATE TABLE IF NOT EXISTS users (
               id SERIAL PRIMARY KEY,
         username CHARACTER VARYING(50),
    password_hash CHARACTER VARYING(64)
 );

CREATE TABLE IF NOT EXISTS orders (
             id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users (id),
         number CHARACTER VARYING(50),
    uploaded_at TIMESTAMPTZ,
         status CHARACTER VARYING(50),
        accrual INTEGER DEFAULT 0,
     withdrawal INTEGER DEFAULT 0
 )