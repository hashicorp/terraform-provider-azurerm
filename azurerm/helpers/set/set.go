package set

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
)

func HashInt(v interface{}) int {
	return hashcode.String(strconv.Itoa(v.(int)))
}

func HashStringIgnoreCase(v interface{}) int {
	return hashcode.String(strings.ToLower(v.(string)))
}
