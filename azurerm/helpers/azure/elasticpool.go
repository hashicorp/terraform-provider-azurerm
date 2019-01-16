package azure

import (
	"fmt"
	"math"
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
	switch {
	case capacity == 50:
	case capacity == 100:
	case capacity == 200:
	case capacity == 300:
	case capacity == 400:
	case capacity == 800:
	case capacity == 1200:
	case capacity == 1600:
	default:
		return false
	}

	return true
}

func StandardGetMaxSizeGB(DTUs int) float64 {
	switch DTUs {
	case 50:
		return 500
	case 100:
		return 750
	case 200:
		return 1000
	case 300:
		return 1250
	case 400:
		return 1500
	case 800:
		return 2000
	case 1200:
		return 2500
	case 1600:
		return 3000
	case 2000:
		return 3500
	case 2500:
		fallthrough
	case 3000:
		return 4000
	}
	// Invalid DTU
	return -1
}

func StandardCapacityValid(capacity int) bool {
	switch {
	case capacity == 50:
	case capacity == 100:
	case capacity == 200:
	case capacity == 300:
	case capacity == 400:
	case capacity == 800:
	case capacity == 1200:
	case capacity == 1600:
	case capacity == 2000:
	case capacity == 2500:
	case capacity == 3000:
	default:
		return false
	}

	return true
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
		return 1000
	case 1500:
		return 1500
	case 2000:
		return 2000
	case 2500:
		return 2500
	case 3000:
		return 3000
	case 3500:
		return 3500
	case 4000:
		return 4000
	}
	// Invalid DTU
	return -1
}

func PremiumCapacityValid(capacity int) bool {
	switch {
	case capacity == 125:
	case capacity == 250:
	case capacity == 500:
	case capacity == 1000:
	case capacity == 1500:
	case capacity == 2000:
	case capacity == 2500:
	case capacity == 3000:
	case capacity == 3500:
	case capacity == 4000:
	default:
		return false
	}

	return true
}

func GeneralPurposeGetMaxSizeGB(vCores int, SKU string) float64 {

	if SKU == "gen4" {
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
			return 1500
		case 6:
			fallthrough
		case 7:
			fallthrough
		case 8:
			fallthrough
		case 9:
			fallthrough
		case 10:
			return 2000
		case 16:
			return 3500
		case 24:
			return 4000
		}
		// Invalid vCore
		return -1
	} else if SKU == "gen5" {
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
			return 1500
		case 12:
			fallthrough
		case 14:
			fallthrough
		case 16:
			return 2000
		case 18:
			fallthrough
		case 20:
			fallthrough
		case 24:
			return 3000
		case 32:
			fallthrough
		case 40:
			fallthrough
		case 80:
			return 4000
		}
		// Invalid vCore
		return -1
	}

	// Invalid SKU
	return -2
}

func GeneralPurposeCapacityValid(capacity int, family string) bool {
	if family == "gen4" {
		switch {
		case capacity == 1:
		case capacity == 2:
		case capacity == 3:
		case capacity == 4:
		case capacity == 5:
		case capacity == 6:
		case capacity == 7:
		case capacity == 8:
		case capacity == 9:
		case capacity == 10:
		case capacity == 16:
		case capacity == 24:
		default:
			return false
		}

		return true
	}

	if family == "gen5" {
		switch {
		case capacity == 2:
		case capacity == 4:
		case capacity == 6:
		case capacity == 8:
		case capacity == 10:
		case capacity == 12:
		case capacity == 14:
		case capacity == 16:
		case capacity == 18:
		case capacity == 20:
		case capacity == 24:
		case capacity == 32:
		case capacity == 40:
		case capacity == 80:
		default:
			return false
		}

		return true
	}

	return false
}

func BusinessCriticalGetMaxSizeGB(vCores int, SKU string) float64 {

	if SKU == "gen4" {
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
			return 1000
		}
		// Invalid vCore
		return -1
	} else if SKU == "gen5" {
		switch vCores {
		case 4:
			return 1000
		case 6:
			fallthrough
		case 8:
			fallthrough
		case 10:
			return 1500
		case 12:
			fallthrough
		case 14:
			fallthrough
		case 16:
			fallthrough
		case 18:
			fallthrough
		case 20:
			return 3000
		case 24:
			fallthrough
		case 32:
			fallthrough
		case 40:
			fallthrough
		case 80:
			return 4000
		}
		// Invalid vCore
		return -1
	}

	// Invalid SKU
	return -2
}

func BusinessCriticalCapacityValid(capacity int, family string) bool {
	if family == "gen4" {
		switch {
		case capacity == 2:
		case capacity == 3:
		case capacity == 4:
		case capacity == 5:
		case capacity == 6:
		case capacity == 7:
		case capacity == 8:
		case capacity == 9:
		case capacity == 10:
		case capacity == 16:
		case capacity == 24:
		default:
			return false
		}

		return true
	}

	if family == "gen5" {
		switch {
		case capacity == 4:
		case capacity == 6:
		case capacity == 8:
		case capacity == 10:
		case capacity == 12:
		case capacity == 14:
		case capacity == 16:
		case capacity == 18:
		case capacity == 20:
		case capacity == 24:
		case capacity == 32:
		case capacity == 40:
		case capacity == 80:
		default:
			return false
		}

		return true
	}

	return false
}

func IsMaxGBValid(gbIncrement int64, maxSizeGB float64) (msg string, ok bool) {
	// Get the increment for the value in bytes
	// and the maxSizeGB in bytes
	inc := 1073741824 * gbIncrement
	max := 1073741824 * maxSizeGB

	// Check to see if the resulting max_size_bytes value is an integral value
	if max != math.Trunc(float64(max)) {
		return "max_size_gb is not a valid value", false
	}

	// Check to see if the maxSizeGB follows the increment constraint
	if max/float64(inc) != math.Trunc(float64(max/float64(inc))) {
		return fmt.Sprintf("max_size_gb must be defined in increments of %d GB", gbIncrement), false
	}

	return "", true
}
