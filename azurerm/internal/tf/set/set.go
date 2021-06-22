package set

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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

func HashIPv4AddressOrCIDR(ipv4 interface{}) int {
	warnings, errors := validate.IPv4Address(ipv4, "")

	// maybe cidr, just hash it
	if len(warnings) > 0 || len(errors) > 0 {
		return schema.HashString(ipv4)
	}

	// convert to cidr hash
	cidr := fmt.Sprintf("%s/32", ipv4.(string))
	return schema.HashString(cidr)
}
