package main

import (
	//"blockchain1/commandLine"
	"blockchain1/utils"
	"encoding/json"
	"fmt"

	"os"
)

func main() {

	defer os.Exit(0)

	//commandLine.Execute()

	//Usage of the searchNodeInfo function	
	info := utils.GetNodeInfo()
	if info != nil {
		jsonData, err := json.MarshalIndent(info, "", "    ")
		if err != nil {
			fmt.Println("Error converting to JSON.", err)
			return
		}
		fmt.Println("Node Information (JSON):")
		fmt.Println(string(jsonData))
	} else {
		fmt.Println("Node information is nil.")
	}

	// Usage of the CreatePortPIDFile function
	err := utils.CreatePortPIDFile(3001, 999, true)
	if err != nil {
		fmt.Println("Error creating the 'port.pid' file.", err)
	} 

}
