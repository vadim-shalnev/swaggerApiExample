package main

import (
	"database/sql"
	"log"
)

func createTablesIfNotExist(db *sql.DB) {
	createUsersTable := `
CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        "role" VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
    );
    `

	createTableQuery := `
CREATE TABLE IF NOT EXISTS qwery_history (
    id SERIAL PRIMARY KEY,
    user_key INTEGER NOT NULL,
    query VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (user_key) REFERENCES users(id)
);
`

	createTableResponse := `
CREATE TABLE IF NOT EXISTS response_history (
    id SERIAL PRIMARY KEY,
    qwery_key INTEGER NOT NULL,
    address VARCHAR(255) NOT NULL,
    lng VARCHAR(255) NOT NULL,
    lat VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    FOREIGN KEY (qwery_key) REFERENCES qwery_history(id)
);
`
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to start transaction:", err)
	}

	_, err = tx.Exec(createUsersTable)
	if err != nil {
		tx.Rollback()
		log.Fatal("Failed to create users table:", err)
	}

	_, err = tx.Exec(createTableQuery)
	if err != nil {
		tx.Rollback()
		log.Fatal("Failed to create geocodes table:", err)
	}

	_, err = tx.Exec(createTableResponse)
	if err != nil {
		tx.Rollback()
		log.Fatal("Failed to create response_history table:", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}

	log.Println("Tables created successfully")
}
