package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const SecrectKey = "HUYNHMANH22111999"

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, email string, expirationHours int)(string,error){
	expirationTime := time.Now().Add(time.Hour * time.Duration(expirationHours))
	claims := Claims{
		UserID: userID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	
	tokenString,err := token.SignedString([]byte(SecrectKey))
	if err != nil{
		return "",err
	}
	
	return tokenString,nil
}

func VerifyToken(tokenString string) (*Claims,error){
	claims := &Claims{}

	token,err := jwt.ParseWithClaims(tokenString,claims,func(token *jwt.Token)(interface{},error){
		if _,ok := token.Method.(*jwt.SigningMethodHMAC);!ok{
			return nil,errors.New("unexpected signing method")
		}
		return []byte(SecrectKey),nil
	})
	if err != nil {
		return nil,err
	}
	
	if !token.Valid{
		return nil, errors.New("invalid token")
	}
	return claims,nil
}