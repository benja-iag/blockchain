package utils

import (
	"fmt"
	"os"
	"strings"
)

type nodeInfo struct {
	Port string
	PID  string
}


/*	func GetNodeInfo() *nodeInfo {
		filename := "port.pid"
		fileData, err := os.ReadFile(filename)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("El archivo 'port.pid' no existe.")
				return nil
			}
			fmt.Println("Error al leer el archivo 'port.pid':", err)
			return nil
		}

		lines := strings.Split(string(fileData), "\n")
		if len(lines) < 2 {
			fmt.Println("El archivo 'port.pid' no contiene la información requerida.")
			return nil
		}

		return &nodeInfo{
			Port: strings.TrimSpace(lines[0]),
			PID:  strings.TrimSpace(lines[1]),
		}
	}*/

func GetNodeInfo() *nodeInfo {
	filename := "port.pid"
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println("El archivo 'port.pid' no existe.")
		return nil
	} else if err != nil {
		fmt.Println("Error al verificar el archivo 'port.pid':", err)
		return nil
	}

	fileData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error al leer el archivo 'port.pid':", err)
		return nil
	}

	lines := strings.Split(string(fileData), "\n")
	if len(lines) < 2 {
		fmt.Println("El archivo 'port.pid' no contiene la información requerida.")
		return nil
	}

	return &nodeInfo{
		Port: strings.TrimSpace(lines[0]),
		PID:  strings.TrimSpace(lines[1]),
	}
}


