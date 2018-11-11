package commander

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
)

const (
	errTypeMismatch = `ParameterTypeError: %v; Got: %v; Expected: %v;`
	errNotExist     = `Command %v Not Found.`
	errNotFunction  = `The command %v is not a function.`
)

func errTypes(name, got, expect interface{}) error {
	return fmt.Errorf(errTypeMismatch, name, got, expect)
}

// Response is the structure of output from the called function.  If
// the call was successful, Error == nil, and the Result is an array
// of ouput. If there is an error, then Result == nil and Error ==
// "some error string".  Note: when calling Void functions:
// len(Result) = 0.
type Response struct {
	Result interface{}
	Error  interface{}
}

// Command is the structure of a callable command.  The Arguments are
// typed-checked at runtime.
type Command struct {
	Name string
	Args []interface{}
}

// Center contains a Map[string]interface{}, which maps a
// (Function Name) to its (Function).  This enables a valid 'Command' to
// lookup the function by name.
//
// When Creating a CommandCenter, make sure to create the map.
// Otherwise, using the method Call() will always return the error
// "Command Doesn't Exist".
type Center struct {
	FuncMap map[string]interface{}
}

// CallWithCommand is the same as Call, but using the predefined
// Command data structure.
func (c *Center) CallWithCommand(cmd *Command) (interface{}, error) {
	return c.Call(cmd.Name, cmd.Args...)
}

// Call attempts to call the function <name>(<args>...) and does type
// checks to confirm that it can be done.
func (c *Center) Call(name string, args ...interface{}) (interface{}, error) {

	// Retrieve the func:<name> from the map.
	f, ok := c.FuncMap[name]

	// check if func:<name> exists.
	if !ok {
		return nil, fmt.Errorf(errNotExist, name)
	}

	// retrieve the reflection of the func.
	t := reflect.TypeOf(f)

	// confirm that it is a callable func.
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf(errNotFunction, name)
	}

	// retrieve parameter types from target function.
	var paramTypes []reflect.Type
	for i := 0; i < t.NumIn(); i++ {
		paramTypes = append(paramTypes, t.In(i))
	}

	// retieve argument types (and values) from the user's call.
	var argTypes []reflect.Type
	var argVals []reflect.Value
	for _, v := range args {
		argTypes = append(argTypes, reflect.TypeOf(v))
		argVals = append(argVals, reflect.ValueOf(v))
	}

	// compare number of args to params.
	if len(argTypes) != len(paramTypes) {
		return nil, errTypes(name, argTypes, paramTypes)
	}

	// compare argument and parameter types.
	for i := 0; i < len(paramTypes); i++ {
		// case:  floating point --> integer conversion.
		if checkFloatToInt(argTypes[i], paramTypes[i], argVals[i]) {
			argVals[i] = reflect.ValueOf(int(argVals[i].Float()))
			continue
		}
		if paramTypes[i] != argTypes[i] {
			return nil, errTypes(name, argTypes, paramTypes)
		}
	}

	// call the function.
	result := reflect.ValueOf(f).Call(argVals)

	// NOTE:
	//
	// This method returns multiple outputs, type []interface{}
	// output := make([]interface{}, len(result))

	// for i, v := range result {
	// 	output[i] = v.Interface()
	// }

	// This method returns a single result, even if there are more.
	output := result[0].Interface()

	return output, nil
}

// converts the reflection of a float into a integer if that number is
// the same.  For example, float(9) converts to int(9), but float(9.1)
// must remain a float.
func checkFloatToInt(argT, paramT reflect.Type, argV reflect.Value) bool {

	// Confirm that (arg, param) are in (Float, Int)
	// If they aren't, then return immediately.
	if !(argT.Kind() == reflect.Float64 && paramT.Kind() == reflect.Int) {
		return false
	}

	// If a floating point equals its floor, then it's an integer.
	if math.Floor(argV.Float()) == argV.Float() {
		return true
	}

	return false
}

func (c *Command) String() string {
	s := c.Name
	s += "("
	lastIndex := len(c.Args) - 1
	for i, v := range c.Args {
		s += fmt.Sprintf("%#v", v)
		if i != lastIndex {
			s += ", "
		}
	}
	s += ")"
	return s
}

// ------------------------------------------------------------------
// Additional Ways to Call.  Here for Convenience.
// ------------------------------------------------------------------

// CallWithJson attempts to call the function using a JSON object that
// contains Name and Args.  The Structure will match the Command
// structure.
func (c *Center) CallWithJson(b []byte) (interface{}, error) {
	cmd := &Command{}
	err := json.Unmarshal(b, cmd)
	if err != nil {
		return nil, fmt.Errorf("JSON syntax error.")
	}
	return c.Call(cmd.Name, cmd.Args...)
}

// CallFromStrings calls a function in the Command Center using an
// array of strings as arguments for the function.
func (c *Center) CallWithStrings(name string, args []string) (interface{}, error) {
	a := ""
	last := len(args) - 1
	for i, v := range args {
		if i == last {
			a += v
			break
		}
		a += v + ","
	}
	v := fmt.Sprintf(`{"Name":"%s","Args":[%s]}`, name, a)
	return c.CallWithJson([]byte(v))
}

// CallWithFunctionParser parses a string based on the function
// syntax: functionName(arg1, arg2, ...) and uses the result to
// call the function.  Whitespace is removed entirely from the input
// string.  Trailing commas within the list of arguments are optional
// because they are removed.
//
// Arguments are treated as JSON values, so they follow the JSON
// encoding definitions.  If an error comes back as
// "Caller: JSON syntax error", it is most likely because one of the
// arguments is improperly formatted.
func (c *Center) CallWithFunctionString(s string) (interface{}, error) {
	name, args, err := parseFunctionSyntax(s)
	if err != nil {
		return "", fmt.Errorf("Parser: %v", err)
	}
	result, err := c.CallWithStrings(name, args)
	if err != nil {
		return "", fmt.Errorf("Caller: %v", err)
	}
	return result, nil
}
