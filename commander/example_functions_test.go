package commander_test

import (
	"fmt"
	"github.com/tilegame/gameserver/commander"
)

func add(a float64, b float64) float64 {
	return a + b
}

func mul(a float64, b float64) float64 {
	return a * b
}

func Example_callWithFunctionString() {

	// ==================================================
	//  Create the Command Center
	// __________________________________________________

	functionMap := map[string]interface{}{
		"add": add,
		"mul": mul,
	}

	center := commander.Center{
		FuncMap: functionMap,
	}

	// ==================================================
	//  Call Commands with String
	// __________________________________________________

	ex1 := "add(31, 11)"
	ex2 := "mul(6, 7)"

	result1, err1 := center.CallWithFunctionString(ex1)
	fmt.Println(result1, err1)

	result2, err2 := center.CallWithFunctionString(ex2)
	fmt.Println(result2, err2)

	// Output:
	// 42 <nil>
	// 42 <nil>
}
