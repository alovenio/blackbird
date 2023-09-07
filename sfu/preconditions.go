package sfu

import (
	"fmt"
	"regexp"
	"strings"
)

func isNotBlank(n string, v string) error {
	if len(strings.TrimSpace(v)) == 0 {
		return fmt.Errorf("%s must not be blank", n)
	}
	return nil
}

func isId(n string, v string) error {
	valid := false
	if len(v) == IdLen {
		if match, _ := regexp.MatchString("[A-Za-z0-9=+\\-]", v); match {
			valid = true
		}
	}
	if !valid {
		return fmt.Errorf("%s must be a valid id", v)
	}
	return nil
}
