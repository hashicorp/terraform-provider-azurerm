package helper

import (
	"fmt"
	"regexp"
	"strconv"
)

func FormatSuppressionTTL(ttl int) string {
	//API will return error if send ttl = "-1", and API will return ttl = "-1" if send ttl =""
	if ttl == -1 {
		return ""
	}
	//convert TTL to API TTL form, which is "dd.hh:mm:ss"
	m, h, d := "0", "0", "0"
	s := strconv.Itoa(ttl % 60)
	ttl = ttl / 60
	if ttl != 0 {
		m = strconv.Itoa(ttl % 60)
		ttl = ttl / 60
		if ttl != 0 {
			h = strconv.Itoa(ttl % 24)
			ttl = ttl / 24
			if ttl != 0 {
				d = strconv.Itoa(ttl)
			}
		}
	}
	return d + "." + h + ":" + m + ":" + s
}

//Possible api values include: dd.hh:mm:ss, hh:mm:ss, return 0 if ttl string is in invalid format
func ParseSuppresionTTL(ttl string) (int, error) {
	if ttl == "-1" {
		return -1, nil
	}
	if re := regexp.MustCompile(`^(\d+\.)?(\d+):(\d+):(\d+)$`); re.Match([]byte(ttl)) {
		ttlList := re.FindStringSubmatch(ttl)
		days := 0
		if ttlList[1] != "" {
			days, _ = strconv.Atoi(ttlList[1][:len(ttlList[1])-1])
		}
		hours, _ := strconv.Atoi(ttlList[2])
		mins, _ := strconv.Atoi(ttlList[3])
		secs, _ := strconv.Atoi(ttlList[4])
		if hours >= 24 || mins >= 60 || secs >= 60 {
			return 0, fmt.Errorf("time in advisor suppression ttl string is not valid")
		}
		return days*24*60*60 + hours*60*60 + mins*60 + secs, nil
	}
	return 0, fmt.Errorf("advisor suppression ttl string is not valid")
}
