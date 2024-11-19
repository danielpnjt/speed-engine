package utils

import (
	"fmt"
	"regexp"

	"github.com/danielpnjt/speed-engine/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s%s", config.GetString("hashKey"), password)), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(hashedPassword, rawPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(fmt.Sprintf("%s%s", config.GetString("hashKey"), rawPassword))) == nil
}

// func GenerateInitialPasswordFromName(name string) string {
// 	str, _ := uuid.GenerateUUID()
// 	initial := ""
// 	for _, word := range strings.Split(name, " ") {
// 		if len(word) > 0 {
// 			initial += string(word[0])
// 		}
// 	}
// 	return fmt.Sprintf("%sb%s", strings.ToUpper(initial), strings.ReplaceAll(str, "-", ""))[:12]
// }

func IsValidPassword(password string) bool {
	secure := false
	regex := []string{".{8,}", "[a-z]+", "[A-Z]+", "[0-9]+", "\\S+$"}
	for _, r := range regex {
		t, _ := regexp.MatchString(r, password)
		if !t {
			secure = false
			break
		}
		secure = true
	}
	return secure
}
