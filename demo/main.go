package main

import (
	"fmt"

	"github.com/damouse/cumin"
)

type asdf struct {
	i int
}

func (a asdf) ConvertArgument(arg interface{}) (interface{}, error) {
	return nil, nil
}

func bar(a asdf) {

}

func main() {
	c, _ := cumin.NewCurry(bar)
	r, err := c.Invoke([]interface{}{1})

	fmt.Println(r, err)
}
