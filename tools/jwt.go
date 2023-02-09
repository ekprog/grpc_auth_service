package tools

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"os"
	"strconv"
	"time"
)

var (
	jwtExpiredMinutes int
	jwtSecret         []byte
)

func InitJWTTool() error {
	jwtExpStr := os.Getenv("JWT_EXPIRED_MIN")
	parsed, err := strconv.Atoi(jwtExpStr)
	if err != nil {
		return err
	}
	jwtExpiredMinutes = parsed

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	return nil
}

func GenerateJWT(userId int64) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	exp := time.Now().Add(time.Duration(jwtExpiredMinutes) * time.Minute)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = exp.Unix()
	claims["authorized"] = true
	claims["user_id"] = userId

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, exp, nil
}

func VerifyJWT(tokenString string) (bool, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("cannot parse token")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return false, 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userIdFloat := claims["user_id"].(float64)
		userId := int64(userIdFloat)
		return true, userId, nil
	}
	return false, 0, nil
}
