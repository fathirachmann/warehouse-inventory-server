package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run tools/hash_password.go <password>")
		return
	}

	password := os.Args[1]
	// Using cost 10 to match the application's default
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		return
	}

	fmt.Println(string(hashed))
}
