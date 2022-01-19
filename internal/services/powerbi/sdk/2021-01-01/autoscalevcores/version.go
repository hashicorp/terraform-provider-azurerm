package autoscalevcores

import "fmt"

const defaultApiVersion = "2021-01-01"

func userAgent() string {
	return fmt.Sprintf("pandora/autoscalevcores/%s", defaultApiVersion)
}
