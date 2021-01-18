package set

import (
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func HashInt(v interface{}) int {
	return schema.HashString(strconv.Itoa(v.(int)))
}

func HashStringIgnoreCase(v interface{}) int {
	return schema.HashString(strings.ToLower(v.(string)))
}

func FromStringSlice(slice []string) *schema.Set {
	set := &schema.Set{F: schema.HashString}
	for _, v := range slice {
		set.Add(v)
	}
	return set
}

// HashIPv6Address normalizes an IPv6 address and returns a hash for it
func HashIPv6Address(ipv6 interface{}) int {
	return schema.HashString(normalizeIPv6Address(ipv6))
}

// NormalizeIPv6Address returns the normalized notation of an IPv6
func normalizeIPv6Address(ipv6 interface{}) string {
	if ipv6 == nil || ipv6.(string) == "" {
		return ""
	}
	r := net.ParseIP(ipv6.(string))
	if r == nil {
		return ""
	}
	return r.String()
}
