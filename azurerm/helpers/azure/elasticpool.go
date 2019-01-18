package azure

import (
	"math"
	"strings"
)

func BasicGetMaxSizeGB(DTUs int) float64 {
	switch DTUs {
	case 50:
		return 4.8828125
	case 100:
		return 9.765625
	case 200:
		return 19.53125
	case 300:
		return 29.296875
	case 400:
		return 39.0625
	case 800:
		return 78.125
	case 1200:
		return 117.1875
	case 1600:
		return 156.25
	}
	// Invalid DTU
	return -1
}

func BasicIsCapacityValid(capacity int) bool {
	switch capacity {
	case 50:
		fallthrough
	case 100:
		fallthrough
	case 200:
		fallthrough
	case 300:
		fallthrough
	case 400:
		fallthrough
	case 800:
		fallthrough
	case 1200:
		fallthrough
	case 1600:
		return true
	}

	return false
}

func StandardGetMaxSizeGB(DTUs int) float64 {
	switch DTUs {
	case 50:
		return 500
	case 100:
		return 750
	case 200:
		return 1024
	case 300:
		return 1280
	case 400:
		return 1536
	case 800:
		return 2048
	case 1200:
		return 2560
	case 1600:
		return 3072
	case 2000:
		return 3584
	case 2500:
		fallthrough
	case 3000:
		return 4096
	}
	// Invalid DTU
	return -1
}

func StandardCapacityValid(capacity int) bool {
	switch capacity {
	case 50:
		fallthrough
	case 100:
		fallthrough
	case 200:
		fallthrough
	case 300:
		fallthrough
	case 400:
		fallthrough
	case 800:
		fallthrough
	case 1200:
		fallthrough
	case 1600:
		fallthrough
	case 2000:
		fallthrough
	case 2500:
		fallthrough
	case 3000:
		return true
	}

	return false
}

func PremiumGetMaxSizeGB(DTUs int) float64 {
	switch DTUs {
	case 125:
		fallthrough
	case 250:
		fallthrough
	case 500:
		fallthrough
	case 1000:
		return 1024
	case 1500:
		return 1536
	case 2000:
		return 2048
	case 2500:
		return 2560
	case 3000:
		return 3072
	case 3500:
		return 3584
	case 4000:
		return 4096
	}
	// Invalid DTU
	return -1
}

func PremiumCapacityValid(capacity int) bool {
	switch capacity {
	case 125:
		fallthrough
	case 250:
		fallthrough
	case 500:
		fallthrough
	case 1000:
		fallthrough
	case 1500:
		fallthrough
	case 2000:
		fallthrough
	case 2500:
		fallthrough
	case 3000:
		fallthrough
	case 3500:
		fallthrough
	case 4000:
		return true
	}

	return false
}

func StantardPremiumMaxGBValid(gb float64) bool {
	switch gb {
	case 50:
		fallthrough
	case 100:
		fallthrough
	case 150:
		fallthrough
	case 200:
		fallthrough
	case 250:
		fallthrough
	case 300:
		fallthrough
	case 400:
		fallthrough
	case 500:
		fallthrough
	case 750:
		fallthrough
	case 800:
		fallthrough
	case 1024:
		fallthrough
	case 1200:
		fallthrough
	case 1280:
		fallthrough
	case 1536:
		fallthrough
	case 1600:
		fallthrough
	case 1792:
		fallthrough
	case 2000:
		fallthrough
	case 2048:
		fallthrough
	case 2304:
		fallthrough
	case 2500:
		fallthrough
	case 2560:
		fallthrough
	case 2816:
		fallthrough
	case 3000:
		fallthrough
	case 3072:
		fallthrough
	case 3328:
		fallthrough
	case 3584:
		fallthrough
	case 3840:
		fallthrough
	case 4096:
		return true
	}

	return false
}

func GeneralPurposeGetMaxSizeGB(vCores int, family string) float64 {

	if family == "gen4" {
		switch vCores {
		case 1:
			return 512
		case 2:
			return 756
		case 3:
			fallthrough
		case 4:
			fallthrough
		case 5:
			return 1536
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			fallthrough
		case 9:
			fallthrough
		case 10:
			return 2048
		case 16:
			return 3584
		case 24:
			return 4096
		}
		// Invalid vCore
		return -1
	}

	if family == "gen5" {
		switch vCores {
		case 2:
			return 512
		case 4:
			return 756
		case 6:
			fallthrough
		case 8:
			fallthrough
		case 10:
			return 1536
		case 12:
			fallthrough
		case 14:
			fallthrough
		case 16:
			return 2048
		case 18:
			fallthrough
		case 20:
			fallthrough
		case 24:
			return 3072
		case 32:
			fallthrough
		case 40:
			fallthrough
		case 80:
			return 4096
		}
		// Invalid vCore
		return -1
	}

	// Invalid family
	return -2
}

func GeneralPurposeCapacityValid(capacity int, family string) bool {
	if family == "gen4" {
		switch capacity {
		case 1:
			fallthrough
		case 2:
			fallthrough
		case 3:
			fallthrough
		case 4:
			fallthrough
		case 5:
			fallthrough
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			fallthrough
		case 9:
			fallthrough
		case 10:
			fallthrough
		case 16:
			fallthrough
		case 24:
			return true
		}

		return false
	}

	if family == "gen5" {
		switch capacity {
		case 2:
			fallthrough
		case 4:
			fallthrough
		case 6:
			fallthrough
		case 8:
			fallthrough
		case 10:
			fallthrough
		case 12:
			fallthrough
		case 14:
			fallthrough
		case 16:
			fallthrough
		case 18:
			fallthrough
		case 20:
			fallthrough
		case 24:
			fallthrough
		case 32:
			fallthrough
		case 40:
			fallthrough
		case 80:
			return true
		}

		return false
	}

	return false
}

func BusinessCriticalGetMaxSizeGB(vCores int, family string) float64 {

	if family == "gen4" {
		switch vCores {
		case 2:
			fallthrough
		case 3:
			fallthrough
		case 4:
			fallthrough
		case 5:
			fallthrough
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			fallthrough
		case 9:
			fallthrough
		case 10:
			fallthrough
		case 16:
			fallthrough
		case 24:
			return 1024
		}
		// Invalid vCore
		return -1
	}

	if family == "gen5" {
		switch vCores {
		case 4:
			return 1024
		case 6:
			fallthrough
		case 8:
			fallthrough
		case 10:
			return 1536
		case 12:
			fallthrough
		case 14:
			fallthrough
		case 16:
			fallthrough
		case 18:
			fallthrough
		case 20:
			return 3072
		case 24:
			fallthrough
		case 32:
			fallthrough
		case 40:
			fallthrough
		case 80:
			return 4096
		}
		// Invalid vCore
		return -1
	}

	// Invalid family
	return -2
}

func BusinessCriticalCapacityValid(capacity int, family string) bool {
	if family == "gen4" {
		switch capacity {
		case 2:
			fallthrough
		case 3:
			fallthrough
		case 4:
			fallthrough
		case 5:
			fallthrough
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			fallthrough
		case 9:
			fallthrough
		case 10:
			fallthrough
		case 16:
			fallthrough
		case 24:
			return true
		}

		return false
	}

	if family == "gen5" {
		switch capacity {
		case 4:
			fallthrough
		case 6:
			fallthrough
		case 8:
			fallthrough
		case 10:
			fallthrough
		case 12:
			fallthrough
		case 14:
			fallthrough
		case 16:
			fallthrough
		case 18:
			fallthrough
		case 20:
			fallthrough
		case 24:
			fallthrough
		case 32:
			fallthrough
		case 40:
			fallthrough
		case 80:
			return true
		}

		return false
	}

	return false
}

func IsInt(maxSizeGB float64) bool {
	// Get the maxSizeGB for the value in bytes
	max := 1073741824 * maxSizeGB

	// Check to see if the resulting max_size_bytes value is an integral value
	if max != math.Trunc(max) {
		return false
	}

	return true
}

func NameFamilyValid(name string, family string) bool {
	return strings.Contains(strings.ToLower(name), strings.ToLower(family))
}
