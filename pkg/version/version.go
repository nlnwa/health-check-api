package version

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var Version = "undefined"

func GetNotes(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer func() {
		_ = f.Close()
	}()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}
	var v map[string]string
	if err := json.Unmarshal(data, &v); err != nil {
		return nil
	}

	var notes []string
	for key, value := range v {
		notes = append(notes, key+": "+value)
	}
	return notes
}
