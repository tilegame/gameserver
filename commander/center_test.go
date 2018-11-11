package commander

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// ==================================================
// Some Example Commands
// __________________________________________________

func command1(s string, f float64) string {
	out := fmt.Sprint("command 1 was executed, params passed:", s, f)
	return out
}

func command2(s string, i int) string {
	out := fmt.Sprint("woot", s, i)
	return out
}

func add(a, b int) int {
	return a + b
}

func multInt(a, b int) int {
	return a * b
}

func multFloat(a, b float64) float64 {
	return a * b
}

func gimmeTrue() bool {
	return true
}

// ==================================================
// Creating the Actual Command Center
// __________________________________________________

var center = Center{map[string]interface{}{
	"Command1":  command1,
	"Command2":  command2,
	"GimmeTrue": gimmeTrue,
	"Add":       add,
	"multInt":   multInt,
	"multFloat": multFloat,
}}

// ==================================================
// Example Cases
// __________________________________________________

type ex struct {
	s             string
	shouldBeValid bool
}

var cases = []ex{
	ex{`{"Name": "Command1", "Args": ["hello world!", 54]}`, true},
	ex{`{"Name": "Command1", "Args": [123, 123]}`, false},
	ex{`{"Name": "GimmeTrue", "Args": []}`, false},
}

// ==================================================
// The Actual Tests
// __________________________________________________

// Initial test for the programmer's sake.  Makes sure that all of the
// interface{} values in the function map are actual functions.
func TestFuncMap(t *testing.T) {
	for _, val := range center.FuncMap {
		if reflect.ValueOf(val).Kind() != reflect.Func {
			t.Errorf(`Command "%v" must be a function.`, val)
			t.FailNow()
		}
	}
}

// func TestSingleExample(t *testing.T) {
// 	message := []byte{`{"Name":"Command1", "Args":["hello", 123]}`}
// 	cmd := new(Command)
// 	json.Unmarshal(message, cmd)
// }

func TestCases(t *testing.T) {
	for _, example := range cases {
		// t.Logf("======[ Example %v ]======\n", i)
		f := &Command{}
		r := &Response{}

		err := json.Unmarshal([]byte(example.s), f)
		if err != nil {
			t.Error(err)
		}
		// t.Logf("Command: %s \n", f)

		result, err := center.CallWithCommand(f)
		r.Result = result

		if err != nil {
			if example.shouldBeValid {
				t.Error("unexpected error!", err)
			}
			r.Error = err.Error()
		}

		_, err = json.MarshalIndent(r, "", "\t")
		if err != nil {
			t.Error(err)
		}
		// t.Logf("Result:\n%s\n", b)
	}
}

func TestStrings(t *testing.T) {

	cases := []struct {
		in  string
		out string
		ok  bool
	}{
		{"multInt 10 12", "120", true},
		{"multInt 10 1", "10", true},
		{"multInt 10 1.1", "", false},
		{"multInt 10", "", false},
		{"multInt 1.1 1.1", "", false},
		{"multInt ", "", false},
		{"multFloat 10 10", "100", true},
		{"multFloat 10.0 10.1", "101", true},
	}

	for _, c := range cases {
		arr := strings.Split(c.in, " ")
		name, args := arr[0], arr[1:]
		val, err := center.CallWithStrings(name, args)

		if (err != nil) && (c.ok) {
			t.Errorf("Command String:(%v) Expected Error, but got (%v)",
				c.in, err)
			continue
		}
		if (err == nil) && (!c.ok) {
			t.Errorf("Command String Got an Unexpected Error: %v", c.in)
			continue
		}

		if (err != nil) && (!c.ok) {
			// All good.  We expected an error for this case.
			continue
		}

		result := fmt.Sprint(val)
		if result != c.out {
			t.Errorf("Command String Failed:(%s) Expected:(%v), Got:(%v)",
				c.in, result, c.out)
			continue
		}

	}
}
