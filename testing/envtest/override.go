package envtest

import "os"

func Override(vars map[string]string, callback func() error) error {
	current := make(map[string]string)
	for k, v := range vars {
		current[k] = os.Getenv(k)
		os.Setenv(k, v)
	}
	err := callback()
	for k, v := range current {
		os.Setenv(k, v)
	}
	return err
}
