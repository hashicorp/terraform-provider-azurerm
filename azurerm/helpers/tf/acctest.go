package tf

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
)

func AccRandTimeInt() int {
	// acctest.RantInt() returns a value of size:
	// 000000000000000000
	// YYMMddHHmmsshhRRRR

	//go format: 2006-01-02 15:04:05.00

	timeStr := strings.Replace(time.Now().Local().Format("060102150405.00"), ".", "", 1) //no way to not have a .?
	postfix := acctest.RandStringFromCharSet(4, "0123456789")

	i, err := strconv.Atoi(timeStr + postfix)
	if err != nil {
		panic(err)
	}

	return i
}

func AccCheckResourceAttributes(attributes map[string]string, check ...string) error {
	for _, a := range check {
		_, ok := attributes[a]
		if !ok {
			return fmt.Errorf("missing %s", a)
		}
	}

	return nil
}
