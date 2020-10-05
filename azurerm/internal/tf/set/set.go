package set

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func HashInt(v interface{}) int {
	return hashcode.String(strconv.Itoa(v.(int)))
}

func HashStringIgnoreCase(v interface{}) int {
	return hashcode.String(strings.ToLower(v.(string)))
}

func FromStringSlice(slice []string) *schema.Set {
	set := &schema.Set{F: schema.HashString}
	for _, v := range slice {
		set.Add(v)
	}
	return set
}
