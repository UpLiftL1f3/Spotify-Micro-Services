package helperFunc

import "github.com/UpLiftL1f3/Spotify-Micro-Services/monolithic-service/internal/models"

func UserAvatarDereference(user *models.User) models.Avatar {
	if user.Avatar != nil {
		return *user.Avatar
	} else {
		return models.Avatar{
			Url:      "",
			PublicID: "",
		}
	}
}
