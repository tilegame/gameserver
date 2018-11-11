package commander

import (
	"fmt"
	"strings"
)

func parseFunctionSyntax(s string) (string, []string, error) {
	var name string
	var args []string
	var err error
	var currentArgument string

	s = strings.Replace(s, " ", "", -1)

	for i, r := range s {
		switch r {
		case '(':
			s = s[i+1:]
			goto parseArgs
		default:
			name += string(r)
		}
	}
	err = fmt.Errorf("syntax error: '(' not found.")
	goto ret
parseArgs:
	currentArgument = ""
	for i, r := range s {
		switch r {
		case ')':
			s = s[i+1:]
			if len(currentArgument) > 0 {
				args = append(args, currentArgument)
			}
			goto finish
		case ',':
			args = append(args, currentArgument)
			currentArgument = ""
			continue
		default:
			currentArgument += string(r)
		}
	}
	err = fmt.Errorf("syntax error: ')' not found.")
	goto ret
finish:
	if len(s) != 0 {
		err = fmt.Errorf("syntax error: extra characters found after ')'.")
	}
ret:
	return name, args, err
}
