package configmanger

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	CFGPath = "config"
)

var CONFIG map[string]interface{}

func init() {
	CONFIG = make(map[string]interface{})

	files, err := os.ReadDir(CFGPath)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		file, err := os.OpenFile(filepath.Join(CFGPath, f.Name()), os.O_RDWR, 0444)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		config_data := make(map[string]interface{})

		if err := json.NewDecoder(file).Decode(&config_data); err != nil {
			panic(err)
		}
		CONFIG[file.Name()] = config_data
	}
}
