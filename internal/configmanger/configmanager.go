package configmanger

import (
	"encoding/json"
	"my_app/internal/utils"
	"os"
	"path/filepath"
)

const (
	CFGPath = "config"
)

// 不加载的文件列表
var ExcludeFiles = [...]string{"config.json"}

var CONFIG map[string]interface{}

func LoadCofnig() {
	CONFIG = make(map[string]interface{})

	files, err := os.ReadDir(CFGPath)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if index := utils.ArrayIndexOfString(ExcludeFiles[:], f.Name()); index >= 0 {
			continue
		}
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
