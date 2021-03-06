package cumin

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Used to check for the presence of an error
var errorInterface = reflect.TypeOf((*error)(nil)).Elem()
var spicyInterface = reflect.TypeOf((*Spicy)(nil)).Elem()

// A wrapped function
type Curry struct {
	fn     interface{}
	fnType reflect.Type
	name   string
}

// Curry.Invoke attempts to gracefully make types play nicely between incoming arrguments and expected types
// If the expected type conforms to Spicy then Spicy.Convert is invoked with the argument, which allows
// implementers to customize the conversion. This function *must* return the same type as its implementer
type Spicy interface {
	ConvertArgument(interface{}) (interface{}, error)
}

// Wraps the function fn in a Curry struct and returns it. Returns an error if fn
// is not a function
func NewCurry(fn interface{}) (*Curry, error) {
	typ := reflect.TypeOf(fn)

	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("Handler is not a function!")
	}

	c := &Curry{
		name:   getFunctionName(fn),
		fnType: typ,
		fn:     fn,
	}

	return c, nil
}

func (c *Curry) Name() string {
	return c.name
}

func (c *Curry) Type() reflect.Type {
	return c.fnType
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
		arg := getValueOf(args[i])

		// This type wants to override conversion
		if param.Implements(spicyInterface) {
			// fmt.Println("Is Spicy")
		}

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

// Reflects the function name and slices off the package. Panics if not given a function.
func getFunctionName(fn interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()

	// Method above returns functions in the form :  main.foo
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}

// Wraps reflect.ValueOf to handle the case where an integer value is stored as a float64, as JSON unmarshal does.
func getValueOf(x interface{}) reflect.Value {
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
