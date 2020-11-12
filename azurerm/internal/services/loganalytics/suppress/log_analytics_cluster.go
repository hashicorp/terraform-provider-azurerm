package suppress

import (
	"fmt"
	"net"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func LogAnalyticsClusterUrl(_, old, new string, _ *schema.ResourceData) bool {
	u, err := url.ParseRequestURI(old)
	if err != nil || u.Host == "" {
		return false
	}

	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}

	if new == fmt.Sprintf("%s://%s/", u.Scheme, host) {
		return true
	}

	return false
}
