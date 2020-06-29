package azure

import (
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
)

// NormalizeIPv6Address returns the normalized notation of an IPv6
func NormalizeIPv6Address(ipv6 interface{}) string {
	if ipv6 == nil || ipv6.(string) == "" {
		return ""
	}
	r := net.ParseIP(ipv6.(string))
	if r == nil {
		return ""
	}
	return r.String()
}

// HashIPv6Address normalizes an IPv6 address and returns a hash for it
func HashIPv6Address(ipv6 interface{}) int {
	return hashcode.String(NormalizeIPv6Address(ipv6))
}
