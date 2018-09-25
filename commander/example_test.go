package commander_test

import (
	"encoding/json"
	"fmt"
	"github.com/fractalbach/ninjaServer/commander"
)

// ==================================================
// Start with some functions.
// __________________________________________________

func a(f float64) bool {
	return true
}

func b(s string) (string, float64) {
	return "hello", 3.14
}

// ==================================================
//  Create the Command Center
// __________________________________________________

// You an rename the functions to whatever you like.  The string on
// the left side is the publicly accessible name.  The function, on
// the right side, is the actual function that is called.

var functionMap = map[string]interface{}{
	"alwaysTrue": a,
	"sayHello":   b,
}

var myCommandCenter = commander.CommandCenter{
	FuncMap: functionMap,
}

// ==================================================
//  Call a Command From JSON
// __________________________________________________

// In this example, a JSON message (probably originating from a
// network conection), is converted into the type Command.  This
// command is called by the CommandCenter, which produces the result.

var message = []byte(`{"Name":"sayHello", "Params":["ohai"]}`)

func Example() {
	
	cmd := new(commander.Command)
	json.Unmarshal(message, cmd)
	result, err := myCommandCenter.Call(cmd)
	resp := commander.Response{
		Result: result,
		Error: err,
	}
	output, _ := json.Marshal(resp)
	fmt.Println(string(output))
	// OUTPUT: {"Result":["hello",3.14],"Error":null}
}
