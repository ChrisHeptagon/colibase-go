package utils

import (
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func VaildEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

func HandleEnvVars() error {
	MyEnv, err := godotenv.Read()
	if err != nil {
		panic(err)
	}
	for key, value := range MyEnv {
		os.Setenv(key, value)
	}
	return nil
}
