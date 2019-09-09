package version

import "os"

const (
	undefinedVersion = "dev-undefined"
)

var Version = undefinedVersion

func init() {
	if Version == undefinedVersion {
		override := os.Getenv("VERSION_OVERRIDE")
		if override != "" {
			Version = override
		}
	}
}
