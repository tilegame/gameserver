package commander

import (
	"testing"
	"encoding/json"
	"reflect"
)

var cases = []ex{
	ex{`{"Name": "Command1", "Params": ["hello world!", 54]}`, true},
	ex{`{"Name": "Command1", "Params": [123, 123]}`, false},
	ex{`{"Name": "GimmeTrue", "Params": []}`, false},
}

type ex struct {
	s             string
	shouldBeValid bool
}

// Initial test for the programmer's sake.  Makes sure that all of the
// interface{} values in the function map are actual functions.
func TestFuncMap(t *testing.T) {
	for _, val := range cmdmap {
		if reflect.ValueOf(val).Kind() != reflect.Func {
			t.Errorf(`Command "%v" must be a function.`, val)
			t.FailNow()
		}
	}
}


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

		result, err := f.Call()
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
		
		t.Logf("Result:\n%s\n",b)
	}
}
