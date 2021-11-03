package firewallrules

import "fmt"

const defaultApiVersion = "2016-11-01"

func userAgent() string {
	return fmt.Sprintf("pandora/firewallrules/%s", defaultApiVersion)
}
