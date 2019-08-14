package envtest

import "os"

func With(vars map[string]string, callback func() error) error {
	defaults := make(map[string]string)

	for k, v := range vars {
		if val := os.Getenv(k); val != "" {
			defaults[k] = val
		}
		os.Setenv(k, v)
	}

	err := callback()

	for k := range vars {
		os.Unsetenv(k)
	}

	for k, v := range defaults {
		os.Setenv(k, v)
	}

	return err
}
