package tools

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"os"
	"strconv"
	"time"
)

var (
	jwtAccessExpiredMinutes  int
	jwtRefreshExpiredMinutes int
	jwtSecret                []byte
)

type JWTAccess struct {
	AccessToken         string
	RefreshToken        string
	AccessTokenExpired  time.Time
	RefreshTokenExpired time.Time
}

func InitJWTTool() error {
	jwtAccessExpStr := os.Getenv("JWT_EXPIRED_ACCESS_MIN")
	parsedAccess, err := strconv.Atoi(jwtAccessExpStr)
	if err != nil {
		return err
	}
	jwtAccessExpiredMinutes = parsedAccess

	jwtRefreshExpStr := os.Getenv("JWT_EXPIRED_REFRESH_MIN")
	parsedRefresh, err := strconv.Atoi(jwtRefreshExpStr)
	if err != nil {
		return err
	}
	jwtRefreshExpiredMinutes = parsedRefresh

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	return nil
}

func GenerateJWT(pairUUID string) (*JWTAccess, error) {

	// access token
	token := jwt.New(jwt.SigningMethodHS256)
	expAccessToken := time.Now().Add(time.Duration(jwtAccessExpiredMinutes) * time.Minute)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expAccessToken.Unix()
	claims["pairUUID"] = pairUUID

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	expRefreshToken := time.Now().Add(time.Duration(jwtRefreshExpiredMinutes) * time.Minute)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["exp"] = expRefreshToken.Unix()
	rtClaims["pairUUID"] = pairUUID
	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	return &JWTAccess{
		AccessToken:         tokenString,
		RefreshToken:        refreshTokenString,
		AccessTokenExpired:  expAccessToken,
		RefreshTokenExpired: expRefreshToken,
	}, nil
}

// VerifyJWT return pair UUID of access+refresh
// anyToken = accessToken or refreshToken
func VerifyJWT(anyToken string) (ok bool, pairUUID string, err error) {
	token, err := jwt.Parse(anyToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("cannot parse token")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return false, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		pairUUID := claims["pairUUID"].(string)
		return true, pairUUID, nil
	}
	return false, "", nil
}
