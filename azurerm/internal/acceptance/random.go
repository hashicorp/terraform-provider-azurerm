package acceptance

import (
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func RandTimeInt() int {
	// acctest.RantInt() returns a value of size:
	// 000000000000000000
	// YYMMddHHmmsshhRRRR

	// go format: 2006-01-02 15:04:05.00

	timeStr := strings.Replace(time.Now().Local().Format("060102150405.00"), ".", "", 1) // no way to not have a .?
	postfix := acctest.RandStringFromCharSet(4, "0123456789")

	i, err := strconv.Atoi(timeStr + postfix)
	if err != nil {
		panic(err)
	}

	return i
}

// RandString generates a random alphanumeric string of the length specified
func RandString(strlen int) string {
	return acctest.RandString(strlen)
}

func RandStringFromCharSet(strlen int, charSet string) string {
	return acctest.RandStringFromCharSet(strlen, charSet)
}
