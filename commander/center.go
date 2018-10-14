package commander

import (
	"fmt"
	"reflect"
	"encoding/json"
)

const (
	errTypeMismatch = `ParameterTypeError: %v; Got: %v; Expected: %v;`
	errNotExist     = `Command %v Not Found.`
	errNotFunction = `The command %v is not a function.`
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
	Name   string
	Args  []interface{}
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

// CallWithJson attempts to call the function using a JSON object that
// contains Name and Args.  The Structure will match the Command
// structure.
func (c *Center) CallWithJson(b []byte) (interface{}, error) {
	cmd := &Command{}
	err := json.Unmarshal(b, cmd)
	if err != nil {
		return nil, err
	}
	return c.Call(cmd.Name, cmd.Args...)
}

// CallFromStrings calls a function in the Command Center using an
// array of strings as arguments for the function.
func (c *Center) CallWithStrings(name string, args []string) (interface{}, error) {
	a := ""
	last := len(args)-1
	for i, v := range(args) {
		if i == last {
			a += v
			break
		}
		a += v + ","
	}
	v := fmt.Sprintf(`{"Name":"%s","Args":[%s]}`, name, a)
	return c.CallWithJson([]byte(v))
}

// Call attempts to call the function <name>(<args>...) and does type
// checks to confirm that it can be done.
func (c *Center) Call(name string, args ...interface{}) (interface{}, error) {
	
	// check if function <name> exists.
	f, ok := c.FuncMap[name]
	if !ok {
		return nil, fmt.Errorf(errNotExist, name)
	}

	// retrieve the reflection type of the target function.
	t := reflect.TypeOf(f)

	// confirm that the target is a callable function.
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

	// compare arguments to parameters.
	if len(argTypes) != len(paramTypes) {
		return nil, errTypes(name, argTypes, paramTypes)
	}
	for i:=0; i<len(paramTypes); i++ {
		if paramTypes[i] != argTypes[i] {
			return nil, errTypes(name, argTypes, paramTypes)
		}
	}

	// call the function.
	result := reflect.ValueOf(f).Call(argVals)

	// convert the reflection types into interface{} types that
	// the user can make use of.
	output := make([]interface{}, len(result))
	for i, v := range result {
		output[i] = v.Interface()
	}
	return output, nil 
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
