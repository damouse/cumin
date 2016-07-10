package main

import (
	"fmt"
	"reflect"
	"runtime"
)

func foo() {
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type asdf struct{}

func (a asdf) receiver() {

}

func main() {
	// main.foo
	fmt.Println("name:", GetFunctionName(foo))

	// main.main.func1
	anonymous := func() {}
	fmt.Println("name:", GetFunctionName(anonymous))

	// main.(asdf).(main.receiver)
	a := asdf{}
	fmt.Println("name:", GetFunctionName(a.receiver))
}
