package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/UpLiftL1f3/Spotify-Micro-Services/monolithic-service/internal/env"
	functions "github.com/UpLiftL1f3/Spotify-Micro-Services/shared/functions"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	ScopeAuthentication = "authentication"
)

type TokenModel struct {
	*DBModel // Embed DBModel to provide access to the database
	// User     User
}

// -> Token is the type for authentication tokens
type Token struct {
	PlainText string    `json:"token"`
	UserID    uuid.UUID `json:"_"`
	Hash      []byte    `json:"_"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// -> Generates a token that lasts for ttl and returns it
func GenerateToken(userID uuid.UUID, ttl time.Duration, scope string) (*Token, error) {

	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256(([]byte(token.PlainText)))
	token.Hash = hash[:]
	return token, nil
}

//  ________          __
// |        \        |  \
//  \$$$$$$$$______  | $$   __   ______   _______
//    | $$  /      \ | $$  /  \ /      \ |       \
//    | $$ |  $$$$$$\| $$_/  $$|  $$$$$$\| $$$$$$$\
//    | $$ | $$  | $$| $$   $$ | $$    $$| $$  | $$
//    | $$ | $$__/ $$| $$$$$$\ | $$$$$$$$| $$  | $$
//    | $$  \$$    $$| $$  \$$\ \$$     \| $$  | $$
//     \$$   \$$$$$$  \$$   \$$  \$$$$$$$ \$$   \$$
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

func (m *TokenModel) InsertToken(t *Token, u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	//* DELETE EXISTING TOKENS
	stmt := functions.BuildDeleteQuery(env.TokensTableName, "user_id")
	// fmt.Printf("INSPECT IT 1?: %s", stmt)
	_, err := m.DB.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		fmt.Println("Query Row Error: ", err)
		return err
	}

	var newID uuid.UUID
	stmt, values := functions.BuildInsertQuery(env.TokensTableName, t.TokenToMap(u))
	// fmt.Printf("INSPECT IT 2?: %s, %v", stmt, values)

	if err = m.DB.QueryRowContext(ctx, stmt, values...).Scan(&newID); err != nil {
		fmt.Println("Query Row Error: ", err)
		return err
	}

	fmt.Println("INSPECT IT 3")
	return nil
}

func (m *TokenModel) GetUserForToken(token string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	fmt.Println("Token BEFORE Hash:", token)

	tokenHash := sha256.Sum256([]byte(token))
	fmt.Println("TokenHash: ", tokenHash)

	//* Find the User's Email
	stmt := functions.BuildFindOneFieldQuery(env.TokensTableName, "token_hash", "email")
	fmt.Printf("INSPECT IT 1?: %s", stmt)
	userEmail := ""
	if err := m.DB.QueryRowContext(ctx, stmt, pq.Array(tokenHash)).Scan(&userEmail); err != nil {
		fmt.Println("Query Row Error: ", err)
		return nil, err
	}

	//* Find the User
	fmt.Println("USER EMAIL: ", userEmail)
	query := functions.BuildFindOneQuery(env.UsersTableName, "email")

	var user User
	row := m.DB.QueryRowContext(ctx, query, userEmail)

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

//	________          __
//
//	|        \        |  \
//	 \$$$$$$$$______  | $$   __   ______   _______
//	   | $$  /      \ | $$  /  \ /      \ |       \
//	   | $$ |  $$$$$$\| $$_/  $$|  $$$$$$\| $$$$$$$\
//	   | $$ | $$  | $$| $$   $$ | $$    $$| $$  | $$
//	   | $$ | $$__/ $$| $$$$$$\ | $$$$$$$$| $$  | $$
//	   | $$  \$$    $$| $$  \$$\ \$$     \| $$  | $$
//	    \$$   \$$$$$$  \$$   \$$  \$$$$$$$ \$$   \$$
//	 __       __             __      __                        __
//
// |  \     /  \           |  \    |  \                      |  \
// | $$\   /  $$  ______  _| $$_   | $$____    ______    ____| $$  _______
// | $$$\ /  $$$ /      \|   $$ \  | $$    \  /      \  /      $$ /       \
// | $$$$\  $$$$|  $$$$$$\\$$$$$$  | $$$$$$$\|  $$$$$$\|  $$$$$$$|  $$$$$$$
// | $$\$$ $$ $$| $$    $$ | $$ __ | $$  | $$| $$  | $$| $$  | $$ \$$    \
// | $$ \$$$| $$| $$$$$$$$ | $$|  \| $$  | $$| $$__/ $$| $$__| $$ _\$$$$$$\
// | $$  \$ | $$ \$$     \  \$$  $$| $$  | $$ \$$    $$ \$$    $$|       $$
//	\$$      \$$  \$$$$$$$   \$$$$  \$$   \$$  \$$$$$$   \$$$$$$$ \$$$$$$$

func (token *Token) TokenToMap(u *User) map[string]interface{} {
	return map[string]interface{}{
		"user_id":    token.UserID,
		"name":       u.LastName,
		"email":      u.Email,
		"token_hash": token.Hash,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
}

func (token *Token) TokenHashToString() string {
	tokenString := base64.StdEncoding.EncodeToString(token.Hash)
	return tokenString
}
