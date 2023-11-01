package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/monolithic-service/internal/env"
	"github.com/UpLiftL1f3/Spotify-Micro-Services/shared/functions"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Avatar struct {
	Url      string
	PublicID string
}

// User is the structure which holds one user from the database.
type User struct {
	ID           uuid.UUID   `json:"id"`
	Email        string      `json:"email"`
	FirstName    string      `json:"firstName,"`
	LastName     string      `json:"lastName"`
	Password     string      `json:"password"`
	Verified     bool        `json:"verified"`
	Avatar       *Avatar     `json:"avatar,omitempty"`
	ActiveEvents []uuid.UUID `json:"activeEvents,omitempty"`
	Followers    []uuid.UUID `json:"followers,omitempty"`
	Following    []uuid.UUID `json:"following,omitempty"`
	Token        []string    `json:"token"`
	Active       int         `json:"active"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

//  __    __
// |  \  |  \
// | $$  | $$  _______   ______    ______
// | $$  | $$ /       \ /      \  /      \
// | $$  | $$|  $$$$$$$|  $$$$$$\|  $$$$$$\
// | $$  | $$ \$$    \ | $$    $$| $$   \$$
// | $$__/ $$ _\$$$$$$\| $$$$$$$$| $$
//  \$$    $$|       $$ \$$     \| $$
//   \$$$$$$  \$$$$$$$   \$$$$$$$ \$$
//  _______              __                __
// |       \            |  \              |  \
// | $$$$$$$\  ______  _| $$_     ______  | $$____    ______    _______   ______
// | $$  | $$ |      \|   $$ \   |      \ | $$    \  |      \  /       \ /      \
// | $$  | $$  \$$$$$$\\$$$$$$    \$$$$$$\| $$$$$$$\  \$$$$$$\|  $$$$$$$|  $$$$$$\
// | $$  | $$ /      $$ | $$ __  /      $$| $$  | $$ /      $$ \$$    \ | $$    $$
// | $$__/ $$|  $$$$$$$ | $$|  \|  $$$$$$$| $$__/ $$|  $$$$$$$ _\$$$$$$\| $$$$$$$$
// | $$    $$ \$$    $$  \$$  $$ \$$    $$| $$    $$ \$$    $$|       $$ \$$     \
//  \$$$$$$$   \$$$$$$$   \$$$$   \$$$$$$$ \$$$$$$$   \$$$$$$$ \$$$$$$$   \$$$$$$$
//  __       __             __      __                        __
// |  \     /  \           |  \    |  \                      |  \
// | $$\   /  $$  ______  _| $$_   | $$____    ______    ____| $$  _______
// | $$$\ /  $$$ /      \|   $$ \  | $$    \  /      \  /      $$ /       \
// | $$$$\  $$$$|  $$$$$$\\$$$$$$  | $$$$$$$\|  $$$$$$\|  $$$$$$$|  $$$$$$$
// | $$\$$ $$ $$| $$    $$ | $$ __ | $$  | $$| $$  | $$| $$  | $$ \$$    \
// | $$ \$$$| $$| $$$$$$$$ | $$|  \| $$  | $$| $$__/ $$| $$__| $$ _\$$$$$$\
// | $$  \$ | $$ \$$     \  \$$  $$| $$  | $$ \$$    $$ \$$    $$|       $$
//  \$$      \$$  \$$$$$$$   \$$$$  \$$   \$$  \$$$$$$   \$$$$$$$ \$$$$$$$

type UserModel struct {
	*DBModel // Embed DBModel to provide access to the database
	// User     User
}

// // GetUserByEmail gets a user by email address
// func (m *DBModel) GetUserByEmail(email string) (User, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	email = strings.ToLower(email)
// 	var u User

// 	row := m.DB.QueryRowContext(ctx, `
// 		select
// 			id, first_name, last_name, email, password, created_at, updated_at
// 		from
// 			users
// 		where email = ?`, email)

// 	err := row.Scan(
// 		&u.ID,
// 		&u.FirstName,
// 		&u.LastName,
// 		&u.Email,
// 		&u.Password,
// 		&u.CreatedAt,
// 		&u.UpdatedAt,
// 	)

// 	if err != nil {
// 		return u, err
// 	}

// 	return u, nil
// }

// GetAll returns a slice of all users, sorted by last name
func (u *UserModel) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name`

	rows, err := u.DB.QueryContext(ctx, query)
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
			pq.Array(&user.ActiveEvents),
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
func (u *UserModel) GetByID(userID uuid.UUID) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := functions.BuildFindOneQuery(env.UsersTableName, "id")

	var user User
	row := u.DB.QueryRowContext(ctx, stmt, userID)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Verified,
		&user.Avatar,
		pq.Array(&user.ActiveEvents),
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
func (u *UserModel) GetUserByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := functions.BuildFindOneQuery(env.UsersTableName, "email")

	var user User
	row := u.DB.QueryRowContext(ctx, query, email)

	var avatarJSON []byte // Create a []byte variable to store the JSON data

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Verified,
		&avatarJSON,
		pq.Array(&user.ActiveEvents),
		pq.Array(&user.Followers),
		pq.Array(&user.Following),
		pq.Array(&user.Token),
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		fmt.Println("error: ", err.Error())
		return nil, err

	}

	// Unmarshal the JSON data into the *models.Avatar struct
	if len(avatarJSON) > 0 {
		err = json.Unmarshal(avatarJSON, &user.Avatar)
		if err != nil {
			fmt.Println("error 2: ", err.Error())
			return nil, err
		}
	}

	return &user, nil
}

// GetByEmail returns one user by email
func (u *UserModel) FindOne(conditionFields map[string]interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	fmt.Println("findOne 1")
	query, values := functions.BuildFindOneQueryDynamic(env.UsersTableName, conditionFields)

	var user User
	row := u.DB.QueryRowContext(ctx, query, values...)
	fmt.Println("findOne 12")
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Verified,
		&user.Avatar,
		pq.Array(&user.ActiveEvents),
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
func (u *UserModel) Update(tableName string, updateFields map[string]interface{}, conditionField string, conditionValue any) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt, values := functions.BuildUpdateQuery(tableName, updateFields, conditionField, conditionValue)

	_, err := u.DB.ExecContext(ctx, stmt, values...)

	if err != nil {
		fmt.Println("updated error: ", err)
		return err
	}

	return nil
}

// Delete deletes one user from the database, by User.ID
func (u *UserModel) Delete(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := u.DB.ExecContext(ctx, stmt, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *UserModel) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := u.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new user into the database.
func (u *UserModel) InsertUser(user User) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	fmt.Println("user data", user)

	if err := CreateUserValidator(&user); err != nil {
		return uuid.Nil, err
	}

	hashedPassword, err := functions.HashString(strings.TrimSpace(user.Password))
	if err != nil {
		return uuid.Nil, fmt.Errorf("password error")
	}

	user.Password = hashedPassword

	var newID uuid.UUID
	stmt, values := functions.BuildInsertQuery(env.UsersTableName, user.UserToMap())
	fmt.Printf("INSPECT IT?: %s, %v", stmt, values)
	// stmt := `
	// INSERT INTO spotifyClone_schema.users (
	// 		email, first_name, last_name, password, verified, avatar,
	// 		active_events, followers, following, token, active, created_at, updated_at
	// 	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	// 	RETURNING id
	// `

	err = u.DB.QueryRowContext(ctx, stmt, values...).Scan(&newID)

	if err != nil {
		fmt.Println("Query Row Error: ", err)
		return uuid.Nil, err
	}

	fmt.Println("INSERT HIT THE END NIL AS ERROR")
	return newID, nil
}

// Insert inserts a new user into the database.
// func (c *CreateUserRequest) Insert(tableName string) (uuid.UUID, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	if err := data.CreateNewUserValidator(c); err != nil {
// 		return uuid.Nil, err
// 	}

// 	hashedPassword, err := functions.HashString(strings.TrimSpace(c.Password))
// 	if err != nil {
// 		return uuid.Nil, fmt.Errorf("password error")
// 	}

// 	var newID uuid.UUID

// 	fields := map[string]interface{}{
// 		"email":      c.Email,
// 		"first_name": c.FirstName,
// 		"last_name":  c.LastName,
// 		"password":   hashedPassword,
// 		// Add more fields as needed
// 	}

// 	stmt, values := functions.BuildInsertQuery(tableName, fields)

// 	err = m.DB.QueryRowContext(ctx, stmt, values...).Scan(&newID)
// 	if err != nil {
// 		fmt.Println("Query Row Error: ", err)
// 		return uuid.Nil, err
// 	}

// 	fmt.Println("INSERT HIT THE END NIL AS ERROR", fields)
// 	return newID, nil
// }

// ResetPassword is the method we will use to change a user's password.
func (u *UserModel) GetPassword(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1`

	var user User
	row := u.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Verified,
		&user.Avatar,
		pq.Array(&user.ActiveEvents),
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
func (u *UserModel) ResetPassword(user User, password string) error {
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
	stmt, values := functions.BuildUpdateQuery(env.UsersTableName, fields, "id", user.ID)

	_, err = u.DB.ExecContext(ctx, stmt, values...)
	if err != nil {
		return err
	}

	return nil
}

//  __    __
// |  \  |  \
// | $$  | $$  _______   ______    ______
// | $$  | $$ /       \ /      \  /      \
// | $$  | $$|  $$$$$$$|  $$$$$$\|  $$$$$$\
// | $$  | $$ \$$    \ | $$    $$| $$   \$$
// | $$__/ $$ _\$$$$$$\| $$$$$$$$| $$
//  \$$    $$|       $$ \$$     \| $$
//   \$$$$$$  \$$$$$$$   \$$$$$$$ \$$
//  __       __             __      __                        __
// |  \     /  \           |  \    |  \                      |  \
// | $$\   /  $$  ______  _| $$_   | $$____    ______    ____| $$  _______
// | $$$\ /  $$$ /      \|   $$ \  | $$    \  /      \  /      $$ /       \
// | $$$$\  $$$$|  $$$$$$\\$$$$$$  | $$$$$$$\|  $$$$$$\|  $$$$$$$|  $$$$$$$
// | $$\$$ $$ $$| $$    $$ | $$ __ | $$  | $$| $$  | $$| $$  | $$ \$$    \
// | $$ \$$$| $$| $$$$$$$$ | $$|  \| $$  | $$| $$__/ $$| $$__| $$ _\$$$$$$\
// | $$  \$ | $$ \$$     \  \$$  $$| $$  | $$ \$$    $$ \$$    $$|       $$
//  \$$      \$$  \$$$$$$$   \$$$$  \$$   \$$  \$$$$$$   \$$$$$$$ \$$$$$$$

func (u *User) GenerateResetLink(token string) string {
	tokenStr := fmt.Sprintf("?token=%s", token)
	userStr := fmt.Sprintf("?userID=%s", u.ID)
	resetPassLink := env.ResetPasswordLink + tokenStr + userStr

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

func (u *User) String() string {
	return fmt.Sprintf("ID: %s, Email: %s, FirstName: %s, LastName: %s, Password: %s, Verified: %v, Avatar: %+v, ActiveEvents: %+v, Followers: %+v, Following: %+v, Token: %+v, Active: %d, CreatedAt: %s, UpdatedAt: %s",
		u.ID, u.Email, u.FirstName, u.LastName, u.Password, u.Verified, u.Avatar, u.ActiveEvents, u.Followers, u.Following, u.Token, u.Active, u.CreatedAt, u.UpdatedAt)
}

func (user *User) UserToMap() map[string]interface{} {
	if len(user.Token) == 0 {
		user.Token = nil
	}
	return map[string]interface{}{
		"ID":           user.ID,
		"Email":        user.Email,
		"FirstName":    user.FirstName,
		"LastName":     user.LastName,
		"Password":     user.Password,
		"Verified":     user.Verified,
		"Avatar":       user.Avatar,
		"ActiveEvents": user.ActiveEvents,
		"Followers":    user.Followers,
		"Following":    user.Following,
		"Token":        user.Token,
		"Active":       user.Active,
		"CreatedAt":    user.CreatedAt,
		"UpdatedAt":    user.UpdatedAt,
	}
}

//  __    __
// |  \  |  \
// | $$  | $$  _______   ______    ______
// | $$  | $$ /       \ /      \  /      \
// | $$  | $$|  $$$$$$$|  $$$$$$\|  $$$$$$\
// | $$  | $$ \$$    \ | $$    $$| $$   \$$
// | $$__/ $$ _\$$$$$$\| $$$$$$$$| $$
//  \$$    $$|       $$ \$$     \| $$
//   \$$$$$$  \$$$$$$$   \$$$$$$$ \$$
//  __     __           __  __        __             __      __
// |  \   |  \         |  \|  \      |  \           |  \    |  \
// | $$   | $$ ______  | $$ \$$  ____| $$  ______  _| $$_    \$$  ______   _______    _______
// | $$   | $$|      \ | $$|  \ /      $$ |      \|   $$ \  |  \ /      \ |       \  /       \
//  \$$\ /  $$ \$$$$$$\| $$| $$|  $$$$$$$  \$$$$$$\\$$$$$$  | $$|  $$$$$$\| $$$$$$$\|  $$$$$$$
//   \$$\  $$ /      $$| $$| $$| $$  | $$ /      $$ | $$ __ | $$| $$  | $$| $$  | $$ \$$    \
//    \$$ $$ |  $$$$$$$| $$| $$| $$__| $$|  $$$$$$$ | $$|  \| $$| $$__/ $$| $$  | $$ _\$$$$$$\
//     \$$$   \$$    $$| $$| $$ \$$    $$ \$$    $$  \$$  $$| $$ \$$    $$| $$  | $$|       $$
//      \$     \$$$$$$$ \$$ \$$  \$$$$$$$  \$$$$$$$   \$$$$  \$$  \$$$$$$  \$$   \$$ \$$$$$$$

// Assuming you have a function to convert a string to a UUID
// func parseUUID(str string) (uuid.UUID, error) {
// 	return uuid.Parse(str)
// }

// -> VALIDATE AN EXISTING USER
func (c *User) UserValidator() error {
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

// Validate checks whether all required fields are filled.
// -> CREATING A USER
func CreateNewUserValidator(u *User) error {
	if err := FirstNameValidator(strings.TrimSpace(u.FirstName)); err != nil {
		return err
	}
	if err := LastNameValidator(strings.TrimSpace(u.LastName)); err != nil {
		return err
	}

	if err := EmailValidator(u.Email); err != nil {
		return err
	}

	if err := PasswordValidator(strings.TrimSpace(u.Password)); err != nil {
		return err
	}

	return nil
}

func CreateUserValidator(u *User) error {
	err := FirstNameValidator(strings.TrimSpace(u.FirstName))
	if err != nil {
		return err
	}

	err = EmailValidator(u.Email)
	if err != nil {
		return err
	}

	err = PasswordValidator(strings.TrimSpace(u.Password))
	if err != nil {
		return err
	}

	return nil
}

func FirstNameValidator(firstName string) error {
	if firstName == "" {
		return fmt.Errorf("name is missing")
	}

	if len(firstName) < 3 {
		return fmt.Errorf("name has to be at least 3 characters")
	}

	if len(firstName) > 20 {
		return fmt.Errorf("name is too long! 20 characters is the max")
	}

	return nil
}
func LastNameValidator(lastName string) error {
	if lastName == "" {
		return fmt.Errorf("name is missing")
	}

	if len(lastName) < 2 {
		return fmt.Errorf("name has to be at least 2 characters")
	}

	if len(lastName) > 100 {
		return fmt.Errorf("name is too long! 20 characters is the max")
	}

	return nil
}

func EmailValidator(email string) error {
	if email == "" {
		return fmt.Errorf("email is missing")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	isValid := emailRegex.MatchString(email)
	if !isValid {
		return fmt.Errorf("invalid email")
	}

	return nil
}

func PasswordValidator(password string) error {
	if password == "" {
		return fmt.Errorf("password is missing")
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	// At least one uppercase letter
	if ok, _ := regexp.MatchString(`[A-Z]`, password); !ok {
		return fmt.Errorf("password must have at least one uppercase letter")
	}

	// At least one lowercase letter
	if ok, _ := regexp.MatchString(`[a-z]`, password); !ok {
		return fmt.Errorf("password must have at least one lowercase letter")
	}

	// At least one digit
	if ok, _ := regexp.MatchString(`\d`, password); !ok {
		return fmt.Errorf("password must have at least one digit")
	}

	// At least one special character
	if ok, _ := regexp.MatchString(`[@$!%*?&]`, password); !ok {
		return fmt.Errorf("password must have at least one special character")
	}

	return nil
}
