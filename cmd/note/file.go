package note

import (
	"fmt"
	"os"
)

// 文件是否存在
func FileExists(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		return !fi.IsDir()
	}
	return os.IsExist(err)
}

// 更具文件目录创建文件夹
func MkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// file test
func FileTest() {
	// a, err := os.ReadDir("./")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	return
	// }
	// for _, v := range a {
	// 	fmt.Printf("i: %v %v %v\n", v.Name(), v.Type(), v.IsDir())
	// }
	// fmt.Printf("a: %v\n", a)
	// fmt.Println(os.UserHomeDir())
	file, err := os.OpenFile("f1.json", os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer file.Close()
	data, err := os.ReadFile("f1.json")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("data: %v\n", data)
	data = append(data, []byte(`
	{
		"name": "zhangsan",
		"age": 18
	}
	`)...)
	err = os.WriteFile("f1.json", data, os.ModeDir)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
}
