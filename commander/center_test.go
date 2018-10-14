package commander

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	for i, example := range cases {

		f := &Command{}
		r := &Response{}

		t.Logf("======[ Example %v ]======\n", i)

		err := json.Unmarshal([]byte(example.s), f)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Command: %s \n", f)

		result, err := center.CallWithCommand(f)
		r.Result = result

		if err != nil {
			if example.shouldBeValid {
				t.Error("unexpected error!", err)
			}
			r.Error = err.Error()
		}

		b, err := json.MarshalIndent(r, "", "\t")
		if err != nil {
			t.Error(err)
		}

		t.Logf("Result:\n%s\n", b)
	}
}
