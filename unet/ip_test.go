package unet

import (
	"encoding/json"
	"fmt"
	"log"
)

func ExampleGetIpAddress() {
	result, err := GetIpAddress("114.115.116.112")
	if err != nil {
		fmt.Println(err)
	} else {
		data, _ := json.MarshalIndent(result, "", "  ")
		log.Println(string(data))
	}

	// output:
	//
}
