package cumin

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCuminUnpacking(t *testing.T) {
	Convey("Functions that return arguments", t, func() {
		f := func() []interface{} {
			return []interface{}{1, 2, 3}
		}

		Convey("Should output just those arguments", func() {
			q := []interface{}{1, 2, 3}

			r, _ := Cumin(f, []interface{}{})
			fmt.Printf("%v\n", r)
			So(len(r[0].([]interface{})), ShouldEqual, len(q))
		})
	})

	Convey("Functions that return an error", t, func() {
		// A function that conditionally returns an error
		f := func(t bool) ([]interface{}, error) {
			if t {
				return nil, fmt.Errorf("An error")
			} else {
				return []interface{}{1, 2, 3}, nil
			}
		}

		// The arguments for that function
		q := []interface{}{1, 2, 3}

		Convey("Should not return that error if it was nil ", func() {
			r, _ := Cumin(f, []interface{}{false})
			fmt.Printf("%v\n", r)
			So(len(r[0].([]interface{})), ShouldEqual, len(q))
			So(len(r), ShouldEqual, 1)
		})

		Convey("Should only return the error if one occured", func() {
			r, e := Cumin(f, []interface{}{true})
			So(e.Error(), ShouldEqual, "An error")
			So(len(r), ShouldEqual, 0)
		})
	})
}

func TestCuminXNone(t *testing.T) {
	Convey("Functions that return nothing", t, func() {
		Convey("Should accept no args", func() {
			_, e := Cumin(noneNone, []interface{}{})
			So(e, ShouldBeNil)
		})

		Convey("Should accept one arg", func() {
			_, e := Cumin(oneNone, []interface{}{1})

			So(e, ShouldBeNil)
		})
	})
}

// Functions for cuminication
func noneNone()     {}
func oneNone(a int) {}

// Run test arguments through a round of JSON
func jsonicate(args ...interface{}) []interface{} {
	var dat []interface{}
	j, _ := json.Marshal(args)
	if err := json.Unmarshal(j, &dat); err != nil {
		panic(err)
	}

	return dat
}

func unmarshal(literal string) []interface{} {
	var dat []interface{}
	if err := json.Unmarshal([]byte(literal), &dat); err != nil {
		panic(err)
	}

	return dat
}
