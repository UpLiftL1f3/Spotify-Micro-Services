package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	functions "github.com/UpLiftL1f3/Spotify-Micro-Services/shared/functions"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3 // three seconds (wow)

var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User:              User{},
		EmailVerification: EmailVerificationDocument{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User              User
	EmailVerification EmailVerificationDocument
}

type Avatar struct {
	Url      string
	PublicID string
}

type CreateUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	Password  string `json:"password"`
}

// User is the structure which holds one user from the database.
type User struct {
	ID             uuid.UUID   `json:"id"`
	Email          string      `json:"email"`
	FirstName      string      `json:"firstName,"`
	LastName       *string     `json:"lastName,omitempty"`
	Password       string      `json:"password"`
	Verified       bool        `json:"verified"`
	Avatar         *Avatar     `json:"avatar,omitempty"`
	FavoritesAudio []uuid.UUID `json:"favoritesAudio,omitempty"`
	Followers      []uuid.UUID `json:"followers,omitempty"`
	Following      []uuid.UUID `json:"following,omitempty"`
	Token          []string    `json:"token"`
	Active         int         `json:"active"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// Validate checks whether all required fields are filled.
func (c *CreateUserRequest) Validate() error {
	if c.Email == "" {
		return errors.New("email is required")
	}
	if c.FirstName == "" {
		return errors.New("first name is required")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
func (c *User) Validate() error {
	if c.Email == "" {
		return errors.New("email is required")
	}
	if c.ID == uuid.Nil {
		return errors.New("userID is required")
	}
	if c.FirstName == "" {
		return errors.New("first name is required")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// GetAll returns a slice of all users, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Verified,
			&user.Avatar,
			pq.Array(&user.FavoritesAudio),
			pq.Array(&user.Followers),
			pq.Array(&user.Following),
			pq.Array(&user.Token),
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// GetByID returns one user by userID
func (u *User) GetByID(userID uuid.UUID) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := functions.BuildFindOneQuery(UsersTableName, "id")

	var user User
	row := db.QueryRowContext(ctx, stmt, userID)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Verified,
		&user.Avatar,
		pq.Array(&user.FavoritesAudio),
		pq.Array(&user.Followers),
		pq.Array(&user.Following),
		pq.Array(&user.Token),
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail returns one user by email
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := functions.BuildFindOneQuery(UsersTableName, "email")

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Verified,
		&user.Avatar,
		pq.Array(&user.FavoritesAudio),
		pq.Array(&user.Followers),
		pq.Array(&user.Following),
		pq.Array(&user.Token),
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail returns one user by email
func (u *User) FindOne(conditionFields map[string]interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	fmt.Println("findOne 1")
	query, values := functions.BuildFindOneQueryDynamic(UsersTableName, conditionFields)

	var user User
	row := db.QueryRowContext(ctx, query, values...)
	fmt.Println("findOne 12")
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Verified,
		&user.Avatar,
		pq.Array(&user.FavoritesAudio),
		pq.Array(&user.Followers),
		pq.Array(&user.Following),
		pq.Array(&user.Token),
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	fmt.Println("findOne 13", query)
	fmt.Printf("findOne 13: %v", values...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *User) Update(tableName string, updateFields map[string]interface{}, conditionField string, conditionValue any) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt, values := functions.BuildUpdateQuery(tableName, updateFields, conditionField, conditionValue)

	_, err := db.ExecContext(ctx, stmt, values...)

	if err != nil {
		fmt.Println("updated error: ", err)
		return err
	}

	return nil
}

// Delete deletes one user from the database, by User.ID
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new user into the database.
func (u *User) Insert() (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if err := CreateUserValidator(u); err != nil {
		return uuid.Nil, err
	}

	hashedPassword, err := functions.HashString(strings.TrimSpace(u.Password))
	if err != nil {
		return uuid.Nil, fmt.Errorf("password error")
	}

	var newID uuid.UUID
	stmt := `
		INSERT INTO spotifyClone_schema.users (
			email, first_name, last_name, password, verified, avatar, 
			favorites_audio, followers, following, token, active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	err = db.QueryRowContext(ctx, stmt,
		u.Email,
		strings.TrimSpace(u.FirstName),
		strings.TrimSpace(*u.LastName),
		hashedPassword,
		u.Verified,
		u.Avatar,
		u.FavoritesAudio,
		u.Followers,
		u.Following,
		u.Token,
		u.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		fmt.Println("Query Row Error: ", err)
		return uuid.Nil, err
	}

	fmt.Println("INSERT HIT THE END NIL AS ERROR")
	return newID, nil
}

// Insert inserts a new user into the database.
func (c *CreateUserRequest) Insert(tableName string) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if err := CreateNewUserValidator(c); err != nil {
		return uuid.Nil, err
	}

	hashedPassword, err := functions.HashString(strings.TrimSpace(c.Password))
	if err != nil {
		return uuid.Nil, fmt.Errorf("password error")
	}

	var newID uuid.UUID

	fields := map[string]interface{}{
		"email":      c.Email,
		"first_name": c.FirstName,
		"password":   hashedPassword,
		// Add more fields as needed
	}

	stmt, values := functions.BuildInsertQuery(tableName, fields)

	err = db.QueryRowContext(ctx, stmt, values...).Scan(&newID)
	if err != nil {
		fmt.Println("Query Row Error: ", err)
		return uuid.Nil, err
	}

	fmt.Println("INSERT HIT THE END NIL AS ERROR")
	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) GetPassword(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Verified,
		&user.Avatar,
		pq.Array(&user.FavoritesAudio),
		pq.Array(&user.Followers),
		pq.Array(&user.Following),
		pq.Array(&user.Token),
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"password": hashedPassword,
		// Add more fields as needed
	}
	stmt, values := functions.BuildUpdateQuery(UsersTableName, fields, "id", u.ID)

	_, err = db.ExecContext(ctx, stmt, values...)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GenerateResetLink(token string) string {
	tokenStr := fmt.Sprintf("?token=%s", token)
	userStr := fmt.Sprintf("?userID=%s", u.ID)
	resetPassLink := ResetPasswordLink + tokenStr + userStr

	return resetPassLink
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
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

// Assuming you have a function to convert a string to a UUID
// func parseUUID(str string) (uuid.UUID, error) {
// 	return uuid.Parse(str)
// }
