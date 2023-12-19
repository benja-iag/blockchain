// we need to delete all the files inside ./tmp/blocks folder
package utils

import (
	"fmt"
	"log"
	"os"
)

func DeleteBlockchain() {
	fmt.Println("Deleting blockchain...")
	files, err := os.ReadDir("./tmp/blocks")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		err := os.Remove("./tmp/blocks/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Deleting blockchain... NICE")
}
