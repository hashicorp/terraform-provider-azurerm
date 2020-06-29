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

func FromInt32Slice(slice []int32) *schema.Set {
	set := &schema.Set{F: HashInt}
	for _, v := range slice {
		set.Add(int(v))
	}

	return set
}

func ToSliceInt32P(set *schema.Set) *[]int32 {
	var slice []int32
	for _, m := range set.List() {
		slice = append(slice, int32(m.(int)))
	}

	return &slice
}

func FromStringSlice(slice []string) *schema.Set {
	set := &schema.Set{F: schema.HashString}
	for _, v := range slice {
		set.Add(v)
	}
	return set
}
