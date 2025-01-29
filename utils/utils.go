package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// type Claims struct {
// 	UserID   int64    `json:"user_id"`
// 	Roles    []string `json:"roles"`
// 	Username string   `json:"username"`
// 	jwt.RegisteredClaims
// }

type Claims struct {
	UserInfo UserInfo
	jwt.RegisteredClaims
}

type UserInfo struct {
	UserID   int64
	Username string
	Roles    []string
}

var jwtSecret = []byte(os.Getenv("jwtsecret"))

func CreateToken(userID int64, username string, roles []string, duration time.Duration) (string, error) {
	// claims := &Claims{
	// 	UserID:   userID,
	//  Username: username,
	// 	Roles:    roles,
	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
	// 	},
	// }

	// UserInfo := map[string]interface{}{
	// 	"UserID":   userID,
	// 	"Username": username,
	// 	"Roles":    roles,
	// }

	userinfo := UserInfo{
		UserID:   userID,
		Username: username,
		Roles:    roles,
	}

	claims := &Claims{
		UserInfo: userinfo,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtSecret)
}

func VerifyToken(tokenStr string) (Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return *claims, err
	}

	return *claims, nil
}

func ConvertToStringSlice(input interface{}) []string {
	interfaceSlice, ok := input.([]interface{})
	if !ok {
		fmt.Println(fmt.Errorf("input is not a slice of interface{}"))
		return nil
	}

	stringSlice := make([]string, len(interfaceSlice))
	for i, v := range interfaceSlice {
		str, ok := v.(string)
		if !ok {
			fmt.Println(fmt.Errorf("value at index %d is not a string", i))
			return nil
		}
		stringSlice[i] = str
	}
	return stringSlice

	// aInterface := input.([]interface{})
	// aString := make([]string, len(aInterface))
	// for i, v := range aInterface {
	// 	aString[i] = v.(string)
	// }
	// return aString
}
