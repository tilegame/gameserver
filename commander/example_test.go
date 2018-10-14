package commander_test

import (
	"fmt"
	"github.com/fractalbach/ninjaServer/commander"
)

// ==================================================
// Start with some functions.
// __________________________________________________

func sum(f1 float64, f2 float64) float64 {
	return f1 + f2
}

// ==================================================
//  Create the Command Center
// __________________________________________________

// You an rename the functions to whatever you like.  The string on
// the left side is the publicly accessible name.  The function, on
// the right side, is the actual function that is called.

var functionMap = map[string]interface{}{
	"sum": sum,
}

var myCommandCenter = commander.Center{
	FuncMap: functionMap,
}

// ==================================================
//  Call a Command
// __________________________________________________

func Example() {

	name := "sum"
	arg1 := 10.5
	arg2 := 20.2

	results, err := myCommandCenter.Call(name, arg1, arg2)
	fmt.Println(results, err)
	// OUTPUT: [30.7] <nil>
}
