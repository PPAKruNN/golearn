package indisk

import (
	"os"
)

func openOrCreateFile(pathToFile string) *os.File {

	handle, err := os.Open(pathToFile)
	if err != nil {
		handle, err = os.Create(pathToFile)

		if err != nil {
			panic("Could not open database file. Exiting.")
		}
	}

	return handle
}

func resetFile(pathToFile string) error {

	handle, err := os.OpenFile(pathToFile, os.O_RDWR|os.O_CREATE, 0755)
	defer handle.Close()

	_, err = handle.Write([]byte("[]"))

	if err != nil {
		panic("Failed reseting accounts file" + " " + err.Error())
	}

	return nil

}

func saveInFile(pathToFile string, data []byte) {

	handle, err := os.OpenFile(pathToFile, os.O_RDWR|os.O_CREATE, 0755)
	_, err = handle.Write(data)

	if err != nil {
		panic("Failed saving accounts to disk" + " " + err.Error())
	}

}
