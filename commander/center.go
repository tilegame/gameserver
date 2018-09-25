/*Package commander provides a data structure that can call functions.

The motivation for writing this package was to allow user-accessible
functions to be written more easily.  The main use-case is to call a
function with a JSON message.  The JSON provides the Function's Name,
and the Arguments for that function.  It's very similar to JSON-RPC.

Example Usage:
   // some example functions that you want accessible.
   f1 := func() {}
   f2 := func() {}

*/package commander

import (
	"fmt"
	"reflect"
)

const (
	errTypeMismatch = `ParameterTypeError: %v; Got: %v; Expected: %v;`
	errNotExist     = `Command %v Not Found.`
)

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
	Params []interface{}
}

// CommandCenter contains a Map[string]interface{}, which maps a
// (Function Name) to its (Function).  This enables a valid 'Command' to
// lookup the function by name.
//
// When Creating a CommandCenter, make sure to create the map.
// Otherwise, using the method Call() will always return the error
// "Command Doesn't Exist".
type CommandCenter struct {
	FuncMap map[string]interface{}
}

// Call uses the given 'Command' to lookup a function in the
// 'CommandCenter.FuncMap'.  If it exists, and all of the argument
// types match the parameters, then the function will be called.
func (center *CommandCenter) Call(c *Command) (interface{}, error) {

	// see if the command functions exists.
	f, ok := center.FuncMap[c.Name]
	if !ok {
		return nil, fmt.Errorf(errNotExist, c.Name)
	}

	// convert to the reflection type
	fr := reflect.ValueOf(f)

	// gather up the parameters types.
	var list []reflect.Type
	var args []reflect.Value

	for _, p := range c.Params {
		list = append(list, reflect.TypeOf(p))
		args = append(args, reflect.ValueOf(p))
	}

	//  We only want to compare the parameters, not its output.
	//  From the real command, create a list of parameter types.
	realtype := reflect.TypeOf(f)
	n := realtype.NumIn()
	expectedList := make([]reflect.Type, n)

	for i := 0; i < n; i++ {
		expectedList[i] = realtype.In(i)
	}

	// Create "function types" that can be compared to each other.
	// Doing the comparison like this allows for more helpful
	// error messages that can be sent back to the user.
	ftype := reflect.FuncOf(list, nil, false)
	rtype := reflect.FuncOf(expectedList, nil, false)

	if rtype != ftype {
		return nil, fmt.Errorf(errTypeMismatch, c.Name, ftype, realtype)
	}

	// If we get this far, then the arguments and parameters
	// match, and the function can be safely called.
	result := fr.Call(args)
	output := make([]interface{}, len(result))
	for i, v := range result {
		output[i] = v.Interface()
	}
	return output, nil
}

func (c *Command) String() string {
	s := c.Name
	s += "("
	lastIndex := len(c.Params) - 1
	for i, v := range c.Params {
		s += fmt.Sprintf("%#v", v)
		if i != lastIndex {
			s += ", "
		}
	}
	s += ")"
	return s
}
