package data

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

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

// User is the structure which holds one user from the database.
type EmailVerificationToken struct {
	Owner     uuid.UUID `json:"owner"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}

// Generate a NEW email Verification Token
func GenerateToken(length int, userId uuid.UUID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	verificationToken := ""

	for i := 0; i < length; i++ {
		randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNumber := randomGenerator.Intn(100) + 1

		verificationToken += fmt.Sprint(randomNumber)
	}

	stmt := `
		INSERT INTO spotifyClone_schema.email_verification_tokens (
			owner, token
		) VALUES ($1, $2)
		RETURNING id
	`

	var newID string
	err := db.QueryRowContext(ctx, stmt,
		userId,            // UUID of the owner
		verificationToken, // Token string
	).Scan(&newID)

	if err != nil {
		fmt.Println("Query Row Error: ", err)
		return "", err
	}

	return verificationToken, nil
}
