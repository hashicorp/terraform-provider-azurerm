package suppress

import (
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// IPv6Compression compares two IPv6 addresses if they are equal (even if compressed)
func IPv6Compression(_, old, new string, _ *schema.ResourceData) bool {
	oldIP := net.ParseIP(old)
	if oldIP == nil {
		return false
	}

	newIP := net.ParseIP(new)
	if newIP == nil {
		return false
	}

	return oldIP.Equal(newIP)
}
