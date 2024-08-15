package util

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	client = resty.New()
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var client *resty.Client

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func PasswordHash(password string, salt string) string {
	sha := sha256.New()
	sha.Write([]byte(password))
	firstSha := sha.Sum(nil)
	sha = sha256.New()
	sha.Write([]byte(salt))
	sha.Write(firstSha)
	sha.Write([]byte(salt))
	passwordHash := sha.Sum(nil)
	return hex.EncodeToString(passwordHash)
}
func ForwardOrPanic(err error) {
	if err != nil {
		panic(err)
	}
}
func ForwardOrRollback(err error, tx *gorm.DB) {
	if err != nil {
		tx.Rollback()
		panic(err)
	}
}
func IntInSlice(haystack []int, needle int) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
