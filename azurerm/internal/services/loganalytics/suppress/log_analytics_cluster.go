package suppress

import (
	"fmt"
	"log"
	"net"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func LogAnalyticsClusterUrl(_, old, new string, _ *schema.ResourceData) bool {
	// verify the uri is valid
	log.Printf("[INFO] Suppress Log Analytics Cluster URI: %s", old)
	u, err := url.ParseRequestURI(old)
	if err != nil || u.Host == "" {
		return false
	}

	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}

	log.Printf("[INFO] Suppress Log Analytics Cluster URI: %s == %s", new, fmt.Sprintf("%s://%s/", u.Scheme, host))
	if new == fmt.Sprintf("%s://%s/", u.Scheme, host) {
		return true
	}

	return false
}
