package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type PortPID struct {
	Port      int  `json:"port"`
	PID       int  `json:"pid"`
	Publisher bool `json:"publisher"`
}

func CreatePortPIDFile(port, pid int, publisher bool) error {
	if _, err := os.Stat("port.pid"); err == nil {
		fmt.Println("The file 'port.pid' already exists.")
		return nil
	}

	file, err := os.Create("port.pid")
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	defer file.Close()

	data := PortPID{
		Port:      port,
		PID:       pid,
		Publisher: publisher,
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
