package videoanalyzer

import "fmt"

const defaultApiVersion = "2021-05-01-preview"

func userAgent() string {
	return fmt.Sprintf("pandora/videoanalyzer/%s", defaultApiVersion)
}
