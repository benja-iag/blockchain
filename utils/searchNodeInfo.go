package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

/*
	type nodeInfo struct {
		Port int `json:"port"`
		PID  int `json:"pid"`
	}

	func GetNodeInfo() *nodeInfo {
		filename := "port.pid"
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			fmt.Println("The file 'port.pid' does not exist.")
			return nil
		} else if err != nil {
			fmt.Println("Error checking the file 'port.pid':", err)
			return nil
		}

		fileData, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading the file 'port.pid':", err)
			return nil
		}

		var data nodeInfo
		err = json.Unmarshal(fileData, &data)
		if err != nil {
			fmt.Println("Error decoding the file 'port.pid' as JSON:", err)
			return nil
		}

		return &data
	}
*/
type nodeInfo struct {
	Port      int  `json:"port"`
	PID       int  `json:"pid"`
	Publisher bool `json:"publisher"`
}

func GetNodeInfo() *nodeInfo {
	filename := "port.pid"
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println("The file 'port.pid' does not exist.")
		return nil
	} else if err != nil {
		fmt.Println("Error checking the file 'port.pid':", err)
		return nil
	}

	fileData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading the file 'port.pid':", err)
		return nil
	}

	var data nodeInfo
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		fmt.Println("Error decoding the file 'port.pid' as JSON:", err)
		return nil
	}

	return &data
}
