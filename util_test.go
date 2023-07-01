package subcmd

import "os"

// testSetenv sets the environment variable key to val and returns a function that restores the environment.
// It's like "testing.T".Setenv but for older versions of Go.
func testSetenv(key, val string) func() {
	oldval, ok := os.LookupEnv(key)
	os.Setenv(key, val)
	if ok {
		return func() { os.Setenv(key, oldval) }
	}
	return func() { os.Unsetenv(key) }
}
