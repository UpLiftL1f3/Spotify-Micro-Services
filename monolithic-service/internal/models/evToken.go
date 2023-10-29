package models

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/monolithic-service/internal/env"
	functions "github.com/UpLiftL1f3/Spotify-Micro-Services/shared/functions"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type EmailModel struct {
	*DBModel // Embed DBModel to provide access to the database
	EmailVD  EmailVerificationDocument
	EmailVR  EmailVerificationRequest //-> EmailVerificationRequest
}

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

type EmailVerificationRequest struct {
	UserID uuid.UUID `json:"userID"`
	Token  string    `json:"token"`
}

func (m *EmailModel) ValidateParams() error {
	if m.EmailVR.UserID == uuid.Nil {
		return errors.New("userID is required")
	}
	if m.EmailVR.Token == "" {
		return errors.New("a token is required")
	}
	return nil
}

func (m *EmailModel) VER_DeleteByUserID() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	deleteStmt := functions.BuildDeleteQuery(env.EmailVerificationTableName, "Owner")

	_, err := m.DB.ExecContext(ctx, deleteStmt, m.EmailVR.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (m *EmailModel) FindEmailVerToken() (*EmailVerificationDocument, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, owner, token, created_at, expires_at from spotifyClone_schema.email_verification_tokens where owner = $1`

	var emailVerificationBody EmailVerificationDocument
	fmt.Println("FindEmailVerToken 1: ", m.EmailVR.UserID)
	row := m.DB.QueryRowContext(ctx, query, m.EmailVR.UserID)

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

func (m *EmailModel) InsertEmailVerificationToken(token string, userID uuid.UUID) error {
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

	stmt, values := functions.BuildInsertQuery(env.EmailVerificationTableName, fields)

	var newID string
	queryErr := m.DB.QueryRowContext(ctx, stmt, values...).Scan(&newID)

	if queryErr != nil {
		fmt.Println("Query Row Error: ", err)
		return err
	}

	fmt.Println("Email Verification Document Made")

	return nil
}

func (m *EmailModel) CompareHashedToken(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(m.EmailVD.Token), []byte(plainText))
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

func (m *EmailModel) DeleteByID() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if m.EmailVD.ID == uuid.Nil {
		return errors.New("missing an ID")
	}

	deleteStmt := functions.BuildDeleteQuery(env.EmailVerificationTableName, "id")

	_, err := m.DB.ExecContext(ctx, deleteStmt, m.EmailVD.ID)
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

	var email EmailModel

	// -> generate Email Verification Document
	if err = email.InsertEmailVerificationToken(token, userID); err != nil {
		return "", err
	}

	return token, nil
}
