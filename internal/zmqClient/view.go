package zmqClient

import (
	"encoding/json"
	"fmt"
)

func ReqTest(data string) {
	jsonData := make(map[string]interface{})
	json.Unmarshal([]byte(data), jsonData)
	fmt.Printf("jsonData: %v\n", jsonData)
}
