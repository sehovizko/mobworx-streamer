package errspec

import "fmt"

func InvalidParameter(name, value string) error {
	return fmt.Errorf("invalid parameter: name=%s; value=%s", name, value)
}

func ParameterIsUndefined(name string) error {
	return fmt.Errorf("parameter is undefined: %s", name)
}

func ParameterShouldBeNull(name, value string) error {
	return fmt.Errorf("parameter should be null: name=%s; value=%s", name, value)
}
