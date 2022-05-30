package alertsmanagement

import "fmt"

const defaultApiVersion = "2021-08-08"

func userAgent() string {
	return fmt.Sprintf("pandora/alertsmanagement/%s", defaultApiVersion)
}
