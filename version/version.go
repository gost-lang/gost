package version

import "fmt"

var (
	Version = "0.11.0"
)

func String() string {
	return fmt.Sprintf("%s", Version)
}
