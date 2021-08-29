package socket

// Py_builtins_str is used by the pickle decoder to parse the server response into a format Go can understand
type Py_builtins_str struct{}

func (c Py_builtins_str) Call(args ...interface{}) (interface{}, error) {
	return args[0], nil
}
