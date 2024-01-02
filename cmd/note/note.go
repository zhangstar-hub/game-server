package note

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

// 常量与变量
func VariableAndConstant(n int) {
	if n < 3 {
		fmt.Println("Variable and constant")
	}
}

func SwitchCase(n int) {
	for i := 0; i < 10; i++ {
		switch n {
		case 1:
			fmt.Println("one")
		case 2:
			fmt.Println("two")
		case 3:
			fmt.Println("three")
		default:
			fmt.Println("default")
		}
	}
}

// for 无限循环测试
func For() {
	index := 1
	for {
		fmt.Println(index, time.Second)
		index += 1
		time.Sleep(1 * time.Second)
	}
}

// if test
func If() {
	i := 2
	fmt.Printf("%p \n", &i)
	if i := 2; i <= 2 {
		fmt.Printf("i的地址为：%p\n", &i)
	} else if i := 3; i < 3 {
		fmt.Printf("i的地址为：%p\n", &i) // 这里的重新声明会导致编译错误
	}
	{
		i := 3
		fmt.Printf("%p \n", &i)
	}
}

// 闭包测试
func Closure() {
	var i int
	f := func() {
		i++
		fmt.Println(i)
	}
	f()
}

// defer 错误处理
func Defer() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
}

func init() {
	fmt.Println("note init")
}

func ForTest(n int) {
	b := [...][4]int{
		{1, 2, 3},
		{1, 2, 3},
		{1, 2, 3},
	}
	b[1] = [4]int{1, 2, 3, 4}
	// fmt.Println(b)
}

// 切片测试
func Slice() {
	var a = [...]int{1, 2, 3}
	// var b []int = a[1:]
	var b = a[1:]
	b[0] = 100
	fmt.Println(a, b)
	c := make([]int, 3, 4)
	fmt.Printf("%v, %v\n", c, cap(c))
	d := []int{1, 2}
	copy(d, b)
	fmt.Println(b, d)

	s := "你好世界"
	s2 := []byte(s)
	fmt.Printf("%v, %v\n", s2, cap(s2))
	fmt.Printf("%s\n", s2)
	for i, b := range s {
		fmt.Printf("%d, %c\n", i, b)
	}
}

func SelectByKey(text ...string) {
	fmt.Printf("text: %v\n", text)
	fmt.Printf("text: %T\n", text)
	for i, v := range text {
		fmt.Printf("%d, %s\n", i, v)
	}
}

func MapTest() {
	m := make(map[int]string, 2)
	m[1] = "one"
	m[2] = "tow"
	println(len(m))
	m[3] = "three"
	println(len(m))
	for k, v := range m {
		fmt.Printf("%d, %s\n", k, v)
	}
}

type textMes struct {
	text string
	name string
}

type imeMes struct {
	text string
	name string
}

func (t *textMes) SetText() {
	t.text = "textMes"
	fmt.Printf("%s\n", t.text)
}

func (t *textMes) getName() {
	t.name = "textMes getName"
	fmt.Printf("%s\n", t.name)
}

func (i *imeMes) SetText() {
	i.text = "imeMes"
	fmt.Printf("%s\n", i.text)
}

func (i *imeMes) getName() {
	i.name = "imeMes getName"
	fmt.Printf("%s\n", i.name)
}

type Mes interface {
	SetText()
}

func SendMessage(mes Mes) {
	mes.SetText()
	switch t := mes.(type) {
	case *textMes:
		t.getName()
	case *imeMes:
		t.getName()
	}
}

func Interface() {
	t1 := textMes{}
	SendMessage(&t1)
	t2 := imeMes{}
	SendMessage(&t2)
}

func ChannelTest() {
	c := make(chan int, 10)
	for i := 1; i < 10000; i++ {
		go func(j int) {
			c <- j
		}(i)
	}

For:
	for {

		select {
		case v := <-c:
			fmt.Printf("%d\n", v)
		default:
			break For
		}
	}

}

// 计算字符串长度
func StringLength(str string) int {
	return len([]rune(str))
}

func ArgsTest() {
	fmt.Printf("os.Args: %v\n", os.Args)
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("os.Args[%d]: %s %T\n", i, os.Args[i], os.Args[i])
	}
	x := flag.String("x", "", "input x")
	b := flag.Int("b", 12, "input b")
	flag.Parse()
	fmt.Printf("x: %v\n", *x)
	fmt.Printf("b: %v\n", *b)
}

func LockTest() {
	var lock sync.Mutex
	var wg sync.WaitGroup
	var once sync.Once

	cond := sync.NewCond(&lock)
	ret := 1
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func() {
			println("before lock")
			lock.Lock()
			time.Sleep(100 * time.Millisecond)
			ret += 1
			fmt.Printf("ret: %v\n", ret)
			wg.Done()
			lock.Unlock()
			println("after unlock")
			once.Do(func() {
				fmt.Println("once")
			})
		}()
	}

	ret = 1
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func() {
			println("before lock")
			cond.L.Lock()
			ret += 1
			cond.Wait()
			cond.L.Unlock()
			wg.Done()
			println("after unlock")
		}()
	}
	time.Sleep(2 * time.Second)
	cond.Signal()
	time.Sleep(2 * time.Second)
	cond.Signal()
	time.Sleep(2 * time.Second)
	cond.Broadcast()

	wg.Wait()

}
