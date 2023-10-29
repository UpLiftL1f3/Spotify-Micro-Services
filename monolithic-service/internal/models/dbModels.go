package models

import "database/sql"

// DBModel is the type for database connection values
type DBModel struct {
	DB                *sql.DB
	User              UserModel
	EmailVerification EmailModel
}
