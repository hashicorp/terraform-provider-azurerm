package helper

import (
	"strconv"
	"strings"
)

func ConvertToAdvisorSuppresionTTL(ttl int)string{
	if ttl == -1{
		return ""
	}
	s,m,h,d := "0","0","0","0"
	s = strconv.Itoa(ttl%60)
	ttl = ttl/60
	if ttl!=0{
		m = strconv.Itoa(ttl%60)
		ttl = ttl/60
		if ttl!=0{
			h = strconv.Itoa(ttl%24)
			ttl = ttl/24
			if ttl!=0{
				d = strconv.Itoa(ttl)
			}
		}
	}
	return d+"."+h+":"+m+":"+s
}


//Possible api values include: dd.hh:mm:ss, hh:mm:ss
func ParseAdvisorSuppresionTTL(ttl string)int{
	if ttl == "-1"{
		return -1
	}
	if strings.Contains(ttl,"."){
		daysSplit := strings.Split(ttl,".")
		days,err := strconv.Atoi(daysSplit[0])
		if err!=nil{
			return 0
		}
		return days + ParseAdvisorSuppresionTTL(daysSplit[1])
	}
	if strings.Contains(ttl,":"){
		timesSplit := strings.Split(ttl,":")
		hours,err := strconv.Atoi(timesSplit[0])
		if err!=nil{
			return 0
		}
		mins,err := strconv.Atoi(timesSplit[1])
		if err!=nil{
			return 0
		}
		secs,err := strconv.Atoi(timesSplit[2])
		if err!=nil{
			return 0
		}
		return hours*60+ mins*60+secs
	}
	return 0
}
