package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("Password:")
	var pass string
	fmt.Scanln(&pass)
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), 14)

	fmt.Println(string(hash))
}
