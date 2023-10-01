package data

import "os"

var (
	EmailVerificationTableName string
	UsersTableName             string
	ResetPasswordTableName     string
	ResetPasswordLink          string
)

func LoadEnvVariables() {
	// Load environment variables from the .env file
	// err := godotenv.Load("")
	// if err != nil {
	// 	log.Fatal("Error loading .env file", err)
	// }

	// Access environment variables
	UsersTableName = os.Getenv("Users_Table_Name")
	EmailVerificationTableName = os.Getenv("Email_Verification_Table_Name")
	ResetPasswordTableName = os.Getenv("Reset_Password_Tokens")
	ResetPasswordLink = os.Getenv("ResetPasswordLink")

}
