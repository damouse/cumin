package cumin

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Used to check for the presence of an error
var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// Convert and apply args to arbitrary function fn
func Cumin(fn interface{}, args []interface{}) ([]interface{}, error) {
	reciever := reflect.TypeOf(fn)
	var ret []interface{}

	if reciever.Kind() != reflect.Func {
		return ret, fmt.Errorf("Handler is not a function!")
	}

	if reciever.NumIn() != len(args) {
		return ret, fmt.Errorf("Cumin Type Error: expected %d args for function %s, got %d", reciever.NumIn(), reciever, len(args))
	}

	// Iterate over the params listed in the method and try their casts
	values := make([]reflect.Value, len(args))

	for i := 0; i < reciever.NumIn(); i++ {
		param := reciever.In(i)
		arg := GetValueOf(args[i])

		if param == arg.Type() {
			values[i] = arg
		} else if arg.Type().ConvertibleTo(param) {
			values[i] = arg.Convert(param)
		} else {
			return ret, fmt.Errorf("Cumin Type Error: expected %s for arg[%d] in (%s), got %s.", param, i, reciever, arg.Type())
		}
	}

	// Perform the call, collect the results, and return them
	result := reflect.ValueOf(fn).Call(values)

	// If the last value is an error type check its value and finish early
	if len(result) > 0 && result[len(result)-1].Type().Implements(errorInterface) {
		if e := result[len(result)-1]; !e.IsNil() {
			m := e.MethodByName("Error").Call([]reflect.Value{})
			return nil, fmt.Errorf("%v", m[0])
		} else {
			result = result[:len(result)-1]
		}
	}

	// Else return the actual results
	for _, x := range result {
		ret = append(ret, x.Interface())
	}

	return ret, nil
}

// Wraps the function fn in a Curry struct and returns it. Returns an error if fn
// is not a function
func NewCurry(fn interface{}) (*Curry, error) {
	typ := reflect.TypeOf(fn)

	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("Handler is not a function!")
	}

	c := &Curry{
		name:   GetFunctionName(fn),
		fnType: typ,
		fn:     fn,
	}

	return c, nil
}

// Wraps a function in a struct and extracts its name infromation.
type Curry struct {
	fn     interface{}
	fnType reflect.Type
	name   string
}

// Invokes a curried function with the passed arguments. If the function returned an error
// that error is returned, else the results of the function are returned as a slice
func (c *Curry) Invoke(args []interface{}) ([]interface{}, error) {
	if c.fnType.NumIn() != len(args) {
		return nil, fmt.Errorf("Cumin Type Error: expected %d args for function %s, got %d", c.fnType.NumIn(), c.fnType, len(args))
	}

	// Iterate over the params listed in the method and try their casts
	values := make([]reflect.Value, len(args))

	for i := 0; i < c.fnType.NumIn(); i++ {
		param := c.fnType.In(i)
		arg := GetValueOf(args[i])

		if param == arg.Type() {
			values[i] = arg

		} else if arg.Type().ConvertibleTo(param) {
			values[i] = arg.Convert(param)

		} else {
			return nil, fmt.Errorf("Cumin Type Error: expected %s for arg[%d] in (%s), got %s.", param, i, c.fnType, arg.Type())
		}
	}

	// Perform the call, collect the results, and return them
	result := reflect.ValueOf(c.fn).Call(values)

	// If the last value is an error type check its value and finish early
	if len(result) > 0 && result[len(result)-1].Type().Implements(errorInterface) {
		if e := result[len(result)-1]; !e.IsNil() {
			m := e.MethodByName("Error").Call([]reflect.Value{})
			return nil, fmt.Errorf("%v", m[0])
		} else {
			result = result[:len(result)-1]
		}
	}

	// Else return the actual results
	var ret []interface{}
	for _, x := range result {
		ret = append(ret, x.Interface())
	}

	return ret, nil
}

// Return the name of the curried function
func (c *Curry) Name() string {
	return c.name
}

// Reflects the function name and slices off the package. Panics if not given a
// function
func GetFunctionName(fn interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()

	// Method above returns functions in the form :  main.foo
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}

// Wraps reflect.ValueOf to handle the case where an integer value is stored as
// a float64, as JSON unmarshal does.
func GetValueOf(x interface{}) reflect.Value {
	value := reflect.ValueOf(x)

	if value.Kind() == reflect.Float64 {
		asfloat := value.Float()
		asint := int(asfloat)
		if float64(asint) == asfloat {
			value = reflect.ValueOf(asint)
		}
	}

	return value
}
