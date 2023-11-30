package utils

import (
	"fmt"
	"os"
)

/*
	func CreatePortPIDFile(port, pid int) error {
		file, err := os.Create("port.pid")
		if err != nil {
			return fmt.Errorf("Error creating file: %v", err)
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("PORT=%d\nPID=%d\n", port, pid))
		if err != nil {
			return fmt.Errorf("Error writing to file: %v", err)
		}

		return nil
	}
*/
func CreatePortPIDFile(port, pid int) error {
	// Verificar si el archivo 'port.pid' ya existe
	if _, err := os.Stat("port.pid"); err == nil {
		fmt.Println("El archivo 'port.pid' ya existe.")
		return nil
	}

	file, err := os.Create("port.pid")
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("PORT=%d\nPID=%d\n", port, pid))
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

//JSON

//func main() {
//	err := createPortPIDFile(3001, 999)
//	if err != nil {
//		fmt.Println("Error:", err)
//	} else {
//		fmt.Println("File 'port.pid' created successfully.")
//	}
//}
