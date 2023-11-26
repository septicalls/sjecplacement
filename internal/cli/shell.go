package cli

import (
	"bufio"
	"fmt"
	"os"

	"sjecplacement.in/internal/models"
	"sjecplacement.in/internal/validator"
)

func UserShell(m *models.UserModel) {
	scanner := bufio.NewScanner(os.Stdin)

	var name, email, password string

	for {
		fmt.Print("Enter Name: ")
		scanner.Scan()
		name = scanner.Text()

		if !validator.NotBlank(name) || !validator.MaxChar(name, 100) {
			continue
		}

		break
	}

	for {
		fmt.Print("Enter Email: ")
		scanner.Scan()
		email = scanner.Text()

		if !validator.Matches(email, validator.EmailRX) {
			continue
		}

		break
	}

	for {
		fmt.Print("Enter Password: ")
		scanner.Scan()
		password = scanner.Text()

		if !validator.MinChar(password, 8) || !validator.MaxChar(password, 64) {
			continue
		}

		break
	}

	err := m.Insert(name, email, password)
	if err != nil {
		fmt.Printf("Error inserting user: %v\n", err)
		return
	}

	fmt.Printf("%s inserted successfully!", name)
}
