package note

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
	"unicode/utf8"
)

// 随机数
func RandomNum() int {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return r.Intn(10) + 1
}

// 字符串类型转换
func StringConv() {
	i1 := 123
	s1 := "haha"
	s2 := fmt.Sprintf("%d/%s", i1, s1)
	fmt.Println(s2)

	var (
		i2 int
		s3 string
	)
	n, err := fmt.Sscanf(s2, "%d/%s", &i2, &s3)
	if err != nil {
		panic(err)
	}
	fmt.Println(n, i2, s3)

	s4 := strconv.FormatInt(123, 4)
	fmt.Println(s4)
	println(math.MaxUint8)
	i5, err := strconv.ParseUint("256", 10, 8)
	if err != nil {
		panic(err)
	}
	fmt.Println(i5)
}

// 中文字符常见操作
func UTF8Test() {
	s := "hello world, 你好世界"
	fmt.Println(utf8.RuneCountInString(s))
	fmt.Println(utf8.FullRuneInString(s[:len(s)-2]))

	str := "世界"
	fmt.Println(utf8.FullRuneInString(str))
	fmt.Println(utf8.FullRuneInString(str[:4]))

}

// time包测试
func TimePackageTest() {
	s, err := time.ParseDuration("1000s10h1m")
	if err != nil {
		panic(err)
	}
	fmt.Printf("s: %d\n", s)

	s1, err := time.Parse("2006.01.02 15-04:05", "2016.01.02 16-04:05")
	if err != nil {
		panic(err)
	}
	fmt.Printf("s1: %v\n", s1)
	fmt.Printf("time.Now(): %v\n", time.Now())
	time.After(10 * time.Second)
	fmt.Printf("time.Now(): %v\n", time.Now())
	// fmt.Println("start")
	// <-time.After(2 * time.Second)
	// fmt.Println("end")
	// var intChan = make(chan int, 1)
	// select {
	// case <-intChan:
	// 	fmt.Println("intChan")
	// case <-time.After(2 * time.Second):
	// 	fmt.Println("超时")
	// }
	fmt.Printf("time.Monday.String(): %v\n", time.Monday.String())
	fmt.Printf("time.Now(): %v\n", time.Now().Format("2006.01.02"))
	fmt.Printf("time.Now(): %v\n", time.Now().Format("2006.01.02"))
}

var (
	INFO  *log.Logger
	WARN  *log.Logger
	ERROR *log.Logger
)

func LogTest() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	INFO = log.New(logFile, "INFO:", log.LstdFlags|log.Llongfile)
	WARN = log.New(logFile, "WARN:", log.LstdFlags|log.Llongfile)
	ERROR = log.New(logFile, "ERROR:", log.LstdFlags|log.Llongfile)
	INFO.Println("INFO")
	// WARN.Panicln("WARN")
	ERROR.Fatalln("ERROR")
}
