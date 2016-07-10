package cumin

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

// Functions for cuminication
func noneNone()     {}
func oneNone(a int) {}

func TestCurryNonFunction(t *testing.T) {
	_, err := NewCurry(1)
	assert.NotNil(t, err)
}

func TestCurryFunctionNames(t *testing.T) {
	c, err := NewCurry(noneNone)

	assert.Nil(t, err)
	assert.Equal(t, "noneNone", c.Name())
}

// No args, no return
func TestCurriedInvocationNone(t *testing.T) {
	c, _ := NewCurry(noneNone)
	r, err := c.Invoke([]interface{}{})

	assert.Nil(t, err)
	assert.Len(t, r, 0)
}

// Takes one arg, no return
func TestCurriedInvocationArgs(t *testing.T) {
	c, _ := NewCurry(oneNone)
	r, err := c.Invoke([]interface{}{1})

	assert.Nil(t, err)
	assert.Len(t, r, 0)
}

// Curried functions that take no args and have no return
func TestArgsNoReturn(t *testing.T) {
	fn := func(a int) {}
	c, _ := NewCurry(fn)

	Convey("Succeeds with good arguments", t, func() {
		r, err := c.Invoke([]interface{}{1})

		So(err, ShouldBeNil)
		So(len(r), ShouldEqual, 0)
	})

	Convey("Fails with too many arguments", t, func() {
		r, err := c.Invoke([]interface{}{1, 2})

		So(r, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestNoArgsReturn(t *testing.T) {
	fn := func() int {
		return 1
	}

	c, _ := NewCurry(fn)

	Convey("Succeeds with good arguments", t, func() {
		r, err := c.Invoke([]interface{}{})

		So(err, ShouldBeNil)
		So(len(r), ShouldEqual, 1)
		So(r[0], ShouldEqual, 1)
	})

	Convey("Fails with too many arguments", t, func() {
		r, err := c.Invoke([]interface{}{1, 2})

		So(r, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestErrReturn(t *testing.T) {
	fn := func(a int, b string) (int, error) {
		if a == 1 {
			return 0, fmt.Errorf("some error")
		} else {
			return 1, nil
		}
	}

	c, _ := NewCurry(fn)

	Convey("Forwards internal errors", t, func() {
		r, err := c.Invoke([]interface{}{1, ""})

		So(err.Error(), ShouldEqual, "some error")
		So(len(r), ShouldEqual, 0)
	})

	Convey("Returns nothing on bad arguments", t, func() {
		r, err := c.Invoke([]interface{}{})

		So(err, ShouldNotBeNil)
		So(len(r), ShouldEqual, 0)
	})
}
