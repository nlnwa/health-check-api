package version

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func GetVersions(filename string) (map[string]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var v map[string]string
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, err
	}
	return v, err
}
