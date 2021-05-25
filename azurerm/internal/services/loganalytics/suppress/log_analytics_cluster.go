package suppress

import (
	"fmt"
	"net"
	"net/url"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func LogAnalyticsClusterUrl(_, old, new string, _ *pluginsdk.ResourceData) bool {
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
