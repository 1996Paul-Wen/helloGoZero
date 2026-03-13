package util

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

const JWTKeyUserID = "userId"
const JWTEXP = "exp"
const JWTIAT = "iat"

// GenerateToken 生成JWT
func GenerateToken(secret string, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ParseToken 解析JWT
func ParseToken(secret string, tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

func GetUserID(token *jwt.Token) (uint64, error) {
	//提取 claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		if userIdVal, ok := claims[JWTKeyUserID]; ok {
			// 尝试转换为 float64（标准 JSON 数字类型）
			if userIdFloat, ok := userIdVal.(float64); ok {
				userId := uint64(userIdFloat)
				return userId, nil
			}
			return 0, fmt.Errorf("%+v 无法转换为uint64", userIdVal)
		}
		return 0, fmt.Errorf("%+v not found in jwt claims", JWTKeyUserID)

	}
	return 0, fmt.Errorf("jwt claims not found in given token")
}
