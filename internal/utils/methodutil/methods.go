package methodutil

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/abufarhad/golang-starter-rest-api/internal/errors"
	"github.com/abufarhad/golang-starter-rest-api/internal/logger"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func LoadEnv() {
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Println("Error loading .env file", envErr)
	}
}

func IsEmpty(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func StructToStruct(input interface{}, output interface{}) error {
	if b, err := json.Marshal(input); err == nil {
		return json.Unmarshal(b, &output)
	} else {
		return err
	}
}

func MapToStruct(input map[string]interface{}, output interface{}) error {
	if b, err := json.Marshal(input); err == nil {
		return json.Unmarshal(b, &output)
	} else {
		return err
	}
}

func ParseJwtToken(token, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidJwtSigningMethod
		}
		return []byte(secret), nil
	})
}

func StringToIntArray(stringArray []string) []int {
	var res []int

	for _, v := range stringArray {
		if i, err := strconv.Atoi(v); err == nil {
			res = append(res, i)
		}
	}

	return res
}

func GenerateRandomStringOfLength(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	if length == 0 {
		length = 8
	}

	var b strings.Builder

	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}

func RecoverPanic() {
	if r := recover(); r != nil {
		logger.ErrorAsJson("error on panic recover: ", r)
	}
}

func ContainsUint(s []uint, item uint) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

func ContainsInt(s []int, item int) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

func ContainsString(s []string, item string) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

func GenerateHash(password *string) (*string, string) {

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(*password), 8)
	hashPassword := string(hashedPass)
	return &hashPassword, *password
}
