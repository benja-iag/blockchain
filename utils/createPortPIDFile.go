package utils

/*import (
	"fmt"
	"os"
)*/

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
/*func CreatePortPIDFile(port, pid int) error {
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
}*/

import (
	"encoding/json"
	"fmt"
	"os"
)

type PortPID struct {
	Port int `json:"port"`
	PID  int `json:"pid"`
}

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

	data := PortPID{
		Port: port,
		PID:  pid,
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("Error encoding data to JSON: %v", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

//func main() {
//	err := createPortPIDFile(3001, 999)
//	if err != nil {
//		fmt.Println("Error:", err)
//	} else {
//		fmt.Println("File 'port.pid' created successfully.")
//	}
//}
