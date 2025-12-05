package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// Tool to hash a password from command line argument
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run tools/hash_password.go <password>")
		return
	}

	password := os.Args[1]
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		return
	}

	fmt.Println(string(hashed))
}

/* Usage:
   1. go run tools/hash_password.go your_password_here
   2. Insert the output hash into your database or configuration as needed.
*/
