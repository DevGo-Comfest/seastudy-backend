package util

import (
	"log"
	"os"
)

func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("JWT_SECRET environment variable not set, continue with default.")
		secret = "default"
	}
	return []byte(secret)
}
