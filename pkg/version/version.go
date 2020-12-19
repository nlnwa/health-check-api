package version

import (
	"io/ioutil"
	"os"
)

func GetVersions(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
