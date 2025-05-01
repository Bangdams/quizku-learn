package util

import (
	"os"
	"strconv"
	"time"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenJWT(request *entity.User) (string, error) {
	var token model.TokenPyload
	duration := os.Getenv("DURATION_JWT_TOKEN")
	lifeTime, _ := strconv.Atoi(duration)

	now := time.Now()
	token.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "QuizKu",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * time.Duration(lifeTime))),
	}

	token.UserID = request.ID
	token.Email = request.Email
	token.Name = request.Name
	token.Role = request.Role

	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	return _token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}
