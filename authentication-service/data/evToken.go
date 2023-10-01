package data

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	functions "github.com/UpLiftL1f3/Spotify-Micro-Services/shared/functions"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//  __     __                      __   ______
// |  \   |  \                    |  \ /      \
// | $$   | $$  ______    ______   \$$|  $$$$$$\ __    __
// | $$   | $$ /      \  /      \ |  \| $$_  \$$|  \  |  \
//  \$$\ /  $$|  $$$$$$\|  $$$$$$\| $$| $$ \    | $$  | $$
//   \$$\  $$ | $$    $$| $$   \$$| $$| $$$$    | $$  | $$
//    \$$ $$  | $$$$$$$$| $$      | $$| $$      | $$__/ $$
//     \$$$    \$$     \| $$      | $$| $$       \$$    $$
//      \$      \$$$$$$$ \$$       \$$ \$$       _\$$$$$$$
//                                              |  \__| $$
//                                               \$$    $$
//                                                \$$$$$$
//  ________                          __  __
// |        \                        |  \|  \
// | $$$$$$$$ ______ ____    ______   \$$| $$
// | $$__    |      \    \  |      \ |  \| $$
// | $$  \   | $$$$$$\$$$$\  \$$$$$$\| $$| $$
// | $$$$$   | $$ | $$ | $$ /      $$| $$| $$
// | $$_____ | $$ | $$ | $$|  $$$$$$$| $$| $$
// | $$     \| $$ | $$ | $$ \$$    $$| $$| $$
//  \$$$$$$$$ \$$  \$$  \$$  \$$$$$$$ \$$ \$$
//  _______                                                       __
// |       \                                                     |  \
// | $$$$$$$\  ______    ______   __    __   ______    _______  _| $$_
// | $$__| $$ /      \  /      \ |  \  |  \ /      \  /       \|   $$ \
// | $$    $$|  $$$$$$\|  $$$$$$\| $$  | $$|  $$$$$$\|  $$$$$$$ \$$$$$$
// | $$$$$$$\| $$    $$| $$  | $$| $$  | $$| $$    $$ \$$    \   | $$ __
// | $$  | $$| $$$$$$$$| $$__| $$| $$__/ $$| $$$$$$$$ _\$$$$$$\  | $$|  \
// | $$  | $$ \$$     \ \$$    $$ \$$    $$ \$$     \|       $$   \$$  $$
//  \$$   \$$  \$$$$$$$  \$$$$$$$  \$$$$$$   \$$$$$$$ \$$$$$$$     \$$$$
//                           | $$
//                           | $$
//                            \$$

type VerifyEmailRequest struct {
	UserID uuid.UUID `json:"userID"`
	Token  string    `json:"token"`
}

// Validate checks whether all required fields are filled.
func (v *VerifyEmailRequest) ValidateID() error {
	if v.UserID == uuid.Nil {
		return errors.New("userID is required")
	}

	return nil
}

func (v *VerifyEmailRequest) Validate() error {
	if v.UserID == uuid.Nil {
		return errors.New("userID is required")
	}
	if v.Token == "" {
		return errors.New("a token is required")
	}
	return nil
}

func (e *VerifyEmailRequest) DeleteByUserID() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	deleteStmt := functions.BuildDeleteQuery(EmailVerificationTableName, "Owner")

	_, err := db.ExecContext(ctx, deleteStmt, e.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (v *VerifyEmailRequest) FindEmailVerToken() (*EmailVerificationDocument, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, owner, token, created_at, expires_at from spotifyClone_schema.email_verification_tokens where owner = $1`

	var emailVerificationBody EmailVerificationDocument
	fmt.Println("FindEmailVerToken 1: ", v.UserID)
	row := db.QueryRowContext(ctx, query, v.UserID)

	err := row.Scan(
		&emailVerificationBody.ID,
		&emailVerificationBody.Owner,
		&emailVerificationBody.Token,
		&emailVerificationBody.CreatedAt,
		&emailVerificationBody.ExpiresAt,
	)
	fmt.Println("FindEmailVerToken 2")
	if err != nil {
		return nil, err
	}

	return &emailVerificationBody, nil
}

//  ________  __       __   ______   ______  __
// |        \|  \     /  \ /      \ |      \|  \
// | $$$$$$$$| $$\   /  $$|  $$$$$$\ \$$$$$$| $$
// | $$__    | $$$\ /  $$$| $$__| $$  | $$  | $$
// | $$  \   | $$$$\  $$$$| $$    $$  | $$  | $$
// | $$$$$   | $$\$$ $$ $$| $$$$$$$$  | $$  | $$
// | $$_____ | $$ \$$$| $$| $$  | $$ _| $$_ | $$_____
// | $$     \| $$  \$ | $$| $$  | $$|   $$ \| $$     \
//  \$$$$$$$$ \$$      \$$ \$$   \$$ \$$$$$$ \$$$$$$$$
//  __     __  ________  _______   ______  ________  ______   ______    ______  ________  ______   ______   __    __
// |  \   |  \|        \|       \ |      \|        \|      \ /      \  /      \|        \|      \ /      \ |  \  |  \
// | $$   | $$| $$$$$$$$| $$$$$$$\ \$$$$$$| $$$$$$$$ \$$$$$$|  $$$$$$\|  $$$$$$\\$$$$$$$$ \$$$$$$|  $$$$$$\| $$\ | $$
// | $$   | $$| $$__    | $$__| $$  | $$  | $$__      | $$  | $$   \$$| $$__| $$  | $$     | $$  | $$  | $$| $$$\| $$
//  \$$\ /  $$| $$  \   | $$    $$  | $$  | $$  \     | $$  | $$      | $$    $$  | $$     | $$  | $$  | $$| $$$$\ $$
//   \$$\  $$ | $$$$$   | $$$$$$$\  | $$  | $$$$$     | $$  | $$   __ | $$$$$$$$  | $$     | $$  | $$  | $$| $$\$$ $$
//    \$$ $$  | $$_____ | $$  | $$ _| $$_ | $$       _| $$_ | $$__/  \| $$  | $$  | $$    _| $$_ | $$__/ $$| $$ \$$$$
//     \$$$   | $$     \| $$  | $$|   $$ \| $$      |   $$ \ \$$    $$| $$  | $$  | $$   |   $$ \ \$$    $$| $$  \$$$
//      \$     \$$$$$$$$ \$$   \$$ \$$$$$$ \$$       \$$$$$$  \$$$$$$  \$$   \$$   \$$    \$$$$$$  \$$$$$$  \$$   \$$
//  ________   ______   __    __  ________  __    __
// |        \ /      \ |  \  /  \|        \|  \  |  \
//  \$$$$$$$$|  $$$$$$\| $$ /  $$| $$$$$$$$| $$\ | $$
//    | $$   | $$  | $$| $$/  $$ | $$__    | $$$\| $$
//    | $$   | $$  | $$| $$  $$  | $$  \   | $$$$\ $$
//    | $$   | $$  | $$| $$$$$\  | $$$$$   | $$\$$ $$
//    | $$   | $$__/ $$| $$ \$$\ | $$_____ | $$ \$$$$
//    | $$    \$$    $$| $$  \$$\| $$     \| $$  \$$$
//     \$$     \$$$$$$  \$$   \$$ \$$$$$$$$ \$$   \$$

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
// type evTokenModels struct {
// 	User User
// }

// EmailVerificationDocument is the structure which holds one EmailVerificationDocument from the database.
type EmailVerificationDocument struct {
	ID        uuid.UUID `json:"id"`
	Owner     uuid.UUID `json:"owner"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// User is the structure which holds one user from the database.
type CreateEmailVerificationRequest struct {
	Owner uuid.UUID `json:"owner"`
	Token string    `json:"token"`
}

func InsertEmailVerificationToken(token string, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedToken, err := functions.HashString(strings.TrimSpace(token))
	if err != nil {
		return fmt.Errorf("password error")
	}

	fields := map[string]interface{}{
		"owner": userID,
		"token": hashedToken,
		// Add more fields as needed
	}

	stmt, values := functions.BuildInsertQuery(EmailVerificationTableName, fields)

	var newID string
	queryErr := db.QueryRowContext(ctx, stmt, values...).Scan(&newID)

	if queryErr != nil {
		fmt.Println("Query Row Error: ", err)
		return err
	}

	fmt.Println("Email Verification Document Made")

	return nil
}

func (e *EmailVerificationDocument) CompareHashedToken(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e.Token), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (e *EmailVerificationDocument) DeleteByID() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if e.ID == uuid.Nil {
		return errors.New("missing an ID")
	}

	deleteStmt := functions.BuildDeleteQuery(EmailVerificationTableName, "id")

	_, err := db.ExecContext(ctx, deleteStmt, e.ID)
	if err != nil {
		return err
	}

	return nil
}

// EMAIL VERIFICATION
func GenerateTokenAndCreateEVDocument(userID uuid.UUID) (string, error) {
	// Generate the Email Token
	token, err := functions.GenerateToken(6)
	if err != nil {
		return "", err
	}

	// -> generate Email Verification Document
	if err = InsertEmailVerificationToken(token, userID); err != nil {
		return "", err
	}

	return token, nil
}
