package cumin

import (
	"fmt"
	"reflect"
)

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