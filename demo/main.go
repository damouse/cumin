package main

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/damouse/cumin"
)

type asdf struct{}

func (a asdf) receiver() {}
func foo()               {}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func testnames() {
	// main.foo
	fmt.Println("name:", GetFunctionName(foo))

	// main.main.func1
	anonymous := func() {}
	fmt.Println("name:", GetFunctionName(anonymous))

	// main.(asdf).(main.receiver)
	a := asdf{}
	fmt.Println("name:", GetFunctionName(a.receiver))
}

func receive(a int, b string, c float64, d []string, e map[string]interface{}) (string, error) {
	fmt.Println(a, b, c, d, e)
	return "done", nil
}

func main() {
	res, err := cumin.Cumin(receive, []interface{}{1, "2", float64(3), []string{"4", "5"}, map[string]interface{}{"6": 7}})

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
