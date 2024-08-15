package util

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"sduonline-recruitment/pkg/conf"
)

type UserClaims struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
	jwt.StandardClaims
}

func ParseJWT(token string) (*UserClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Conf.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*UserClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
func GenerateJWT(userID int, roleID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		UserID:         userID,
		RoleID:         roleID,
		StandardClaims: jwt.StandardClaims{},
	})
	str, err := token.SignedString([]byte(conf.Conf.JWTSecret))
	if err != nil {
		panic(err)
	}
	return str
}
func ExtractUserClaims(c *gin.Context) *UserClaims {
	claims, ok := c.Get("userClaims")
	if !ok {
		panic("userClaims must exists")
	}
	userClaims, ok := claims.(*UserClaims)
	if !ok {
		panic("userClaims cannot be converted")
	}
	return userClaims
}
func ExtractSectionID(c *gin.Context) int {
	sectionID, ok := c.Get("sectionID")
	if !ok {
		panic("sectionID must exists")
	}
	sectionIDi, ok := sectionID.(int)
	if !ok {
		panic("sectionID cannot be converted")
	}
	return sectionIDi
}
func HasSectionID(c *gin.Context) bool {
	_, ok := c.Get("sectionID")
	return ok
}
