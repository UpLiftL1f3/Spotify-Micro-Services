package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	functions "github.com/UpLiftL1f3/Spotify-Micro-Services/shared/functions"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//  _______                                  __
// |       \                                |  \
// | $$$$$$$\  ______    _______   ______  _| $$_
// | $$__| $$ /      \  /       \ /      \|   $$ \
// | $$    $$|  $$$$$$\|  $$$$$$$|  $$$$$$\\$$$$$$
// | $$$$$$$\| $$    $$ \$$    \ | $$    $$ | $$ __
// | $$  | $$| $$$$$$$$ _\$$$$$$\| $$$$$$$$ | $$|  \
// | $$  | $$ \$$     \|       $$ \$$     \  \$$  $$
//  \$$   \$$  \$$$$$$$ \$$$$$$$   \$$$$$$$   \$$$$
//  _______                                                                        __
// |       \                                                                      |  \
// | $$$$$$$\ ______    _______   _______  __   __   __   ______    ______    ____| $$
// | $$__/ $$|      \  /       \ /       \|  \ |  \ |  \ /      \  /      \  /      $$
// | $$    $$ \$$$$$$\|  $$$$$$$|  $$$$$$$| $$ | $$ | $$|  $$$$$$\|  $$$$$$\|  $$$$$$$
// | $$$$$$$ /      $$ \$$    \  \$$    \ | $$ | $$ | $$| $$  | $$| $$   \$$| $$  | $$
// | $$     |  $$$$$$$ _\$$$$$$\ _\$$$$$$\| $$_/ $$_/ $$| $$__/ $$| $$      | $$__| $$
// | $$      \$$    $$|       $$|       $$ \$$   $$   $$ \$$    $$| $$       \$$    $$
//  \$$       \$$$$$$$ \$$$$$$$  \$$$$$$$   \$$$$$\$$$$   \$$$$$$  \$$        \$$$$$$$
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
//                            \$$                                                                                \$$

type ResetPasswordRequest struct {
	UserID   uuid.UUID `json:"userID"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Token    string    `json:"token"`
}

// Validate checks whether all required fields are filled.
func (r *ResetPasswordRequest) ValidateID() error {
	if r.UserID == uuid.Nil {
		return errors.New("userID is required")
	}
	return nil
}
func (r *ResetPasswordRequest) ValidateEmail() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	return nil
}
func (r *ResetPasswordRequest) ValidateWithoutToken() error {
	if r.UserID == uuid.Nil {
		return errors.New("userID is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
func (r *ResetPasswordRequest) ValidateWithoutPassword() error {
	if r.UserID == uuid.Nil {
		return errors.New("userID is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Token == "" {
		return errors.New("Token is required")
	}
	return nil
}

func (v *ResetPasswordRequest) DeleteByUserID() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	deleteStmt := functions.BuildDeleteQuery(EmailVerificationTableName, "Owner")

	_, err := db.ExecContext(ctx, deleteStmt, v.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (v *ResetPasswordRequest) FindResetPassToken() (*ResetPasswordDocument, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := functions.BuildFindOneQuery(ResetPasswordTableName, "owner")

	var resetPasswordDocument ResetPasswordDocument
	row := db.QueryRowContext(ctx, query, v.UserID)

	err := row.Scan(
		&resetPasswordDocument.ID,
		&resetPasswordDocument.Owner,
		&resetPasswordDocument.Token,
		&resetPasswordDocument.CreatedAt,
		&resetPasswordDocument.ExpiresAt,
	)
	fmt.Println("FindEmailVerToken 2")
	if err != nil {
		return nil, err
	}

	return &resetPasswordDocument, nil
}

//  _______                                  __
// |       \                                |  \
// | $$$$$$$\  ______    _______   ______  _| $$_
// | $$__| $$ /      \  /       \ /      \|   $$ \
// | $$    $$|  $$$$$$\|  $$$$$$$|  $$$$$$\\$$$$$$
// | $$$$$$$\| $$    $$ \$$    \ | $$    $$ | $$ __
// | $$  | $$| $$$$$$$$ _\$$$$$$\| $$$$$$$$ | $$|  \
// | $$  | $$ \$$     \|       $$ \$$     \  \$$  $$
//  \$$   \$$  \$$$$$$$ \$$$$$$$   \$$$$$$$   \$$$$
//  _______                                                                        __
// |       \                                                                      |  \
// | $$$$$$$\ ______    _______   _______  __   __   __   ______    ______    ____| $$
// | $$__/ $$|      \  /       \ /       \|  \ |  \ |  \ /      \  /      \  /      $$
// | $$    $$ \$$$$$$\|  $$$$$$$|  $$$$$$$| $$ | $$ | $$|  $$$$$$\|  $$$$$$\|  $$$$$$$
// | $$$$$$$ /      $$ \$$    \  \$$    \ | $$ | $$ | $$| $$  | $$| $$   \$$| $$  | $$
// | $$     |  $$$$$$$ _\$$$$$$\ _\$$$$$$\| $$_/ $$_/ $$| $$__/ $$| $$      | $$__| $$
// | $$      \$$    $$|       $$|       $$ \$$   $$   $$ \$$    $$| $$       \$$    $$
//  \$$       \$$$$$$$ \$$$$$$$  \$$$$$$$   \$$$$$\$$$$   \$$$$$$  \$$        \$$$$$$$
//  _______                                                                    __
// |       \                                                                  |  \
// | $$$$$$$\  ______    _______  __    __  ______ ____    ______   _______  _| $$_
// | $$  | $$ /      \  /       \|  \  |  \|      \    \  /      \ |       \|   $$ \
// | $$  | $$|  $$$$$$\|  $$$$$$$| $$  | $$| $$$$$$\$$$$\|  $$$$$$\| $$$$$$$\\$$$$$$
// | $$  | $$| $$  | $$| $$      | $$  | $$| $$ | $$ | $$| $$    $$| $$  | $$ | $$ __
// | $$__/ $$| $$__/ $$| $$_____ | $$__/ $$| $$ | $$ | $$| $$$$$$$$| $$  | $$ | $$|  \
// | $$    $$ \$$    $$ \$$     \ \$$    $$| $$ | $$ | $$ \$$     \| $$  | $$  \$$  $$
//  \$$$$$$$   \$$$$$$   \$$$$$$$  \$$$$$$  \$$  \$$  \$$  \$$$$$$$ \$$   \$$   \$$$$

// ResetPasswordDocument is the structure which holds one ResetPasswordDocument from the database.
type ResetPasswordDocument struct {
	ID        uuid.UUID `json:"id"`
	Owner     uuid.UUID `json:"owner"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func InsertResetPasswordToken(token string, userID uuid.UUID) error {
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

	query, values := functions.BuildInsertQuery(ResetPasswordTableName, fields)

	var newID string
	queryErr := db.QueryRowContext(ctx, query, values...).Scan(&newID)

	if queryErr != nil {
		fmt.Println("Query Row Error: ", err)
		return err
	}

	fmt.Println("Email Verification Document Made")

	return nil
}

func (r *ResetPasswordDocument) CompareHashedToken(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(r.Token), []byte(plainText))
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

func FindResetPassToken(userID uuid.UUID) (*ResetPasswordDocument, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := functions.BuildFindOneQuery(ResetPasswordTableName, "owner")

	var resetPasswordDocument ResetPasswordDocument
	row := db.QueryRowContext(ctx, query, userID)

	err := row.Scan(
		&resetPasswordDocument.ID,
		&resetPasswordDocument.Owner,
		&resetPasswordDocument.Token,
		&resetPasswordDocument.CreatedAt,
		&resetPasswordDocument.ExpiresAt,
	)
	fmt.Println("FindEmailVerToken 2")
	if err != nil {
		return nil, err
	}

	// If err is ErrNoRows, it means the row was not found, and you can handle it accordingly
	if err == sql.ErrNoRows {
		return &ResetPasswordDocument{}, nil
	}

	return &resetPasswordDocument, nil
}

func (r *ResetPasswordDocument) DeleteByID() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if r.ID == uuid.Nil {
		return errors.New("missing an ID")
	}

	deleteStmt := functions.BuildDeleteQuery(ResetPasswordTableName, "id")

	_, err := db.ExecContext(ctx, deleteStmt, r.ID)
	if err != nil {
		return err
	}

	return nil
}

// RESET PASSWORD
func GenerateTokenAndCreateRPDocument(userID uuid.UUID) (string, error) {
	// Generate the Email Token
	token, err := functions.GenerateHexToken(36)
	if err != nil {
		return "", err
	}

	// -> generate Email Verification Document
	if err = InsertResetPasswordToken(token, userID); err != nil {
		return "", err
	}

	return token, nil
}

func FindAndDeleteByID(userID uuid.UUID) error {
	resetPassDocument, err := FindResetPassToken(userID)
	if err != nil {
		return err
	}

	if resetPassDocument == (&ResetPasswordDocument{}) {
		// It's empty
		return nil
	}

	if err := resetPassDocument.DeleteByID(); err != nil {
		return err
	}

	return nil
}
