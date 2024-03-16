package debug

import (
	"encoding/json"
	"fmt"
)

func PrintAsJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ") // Indent with spaces
	if err != nil {
		fmt.Println("error marshalling data into to JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}