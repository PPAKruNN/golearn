package main

import (
	"fmt"
	"os"
)

func main() {

	handle, err := os.Open("accounts.json")

	if err != nil {
		fmt.Print(err)

		handle, err = os.Create("accounts.json")

		if err != nil {
			fmt.Print(err)
			panic("Could not create file")
		}

	}

	fmt.Printf("\nO nome do arquivo é: %s\n", handle.Name())
}
