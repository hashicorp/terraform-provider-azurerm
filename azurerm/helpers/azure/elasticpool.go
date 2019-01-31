package azure

import (
	"strings"
)

type ErrorType int

const (
	Capacity  ErrorType = 0
	MaxSizeGB ErrorType = 1
)

var getDTUErrorMsg = map[string]map[ErrorType]string{
	"basicpool": {
		Capacity:  "service tier 'Basic' must have a 'capacity'(%d) of 50, 100, 200, 300, 400, 800, 1200 or 1600 DTUs",
		MaxSizeGB: "service tier 'Basic' with a 'capacity' of %d must have a 'max_size_gb' of %.7f GB, got %.7f GB",
	},
	"standardpool": {
		Capacity:  "service tier 'Standard' must have a 'capacity'(%d) of 50, 100, 200, 300, 400, 800, 1200, 1600, 2000, 2500 or 3000 DTUs",
		MaxSizeGB: "service tier 'Standard' with a 'capacity' of %d must have a 'max_size_gb' no greater than %d GB, got %d GB",
	},
	"premiumpool": {
		Capacity:  "service tier 'Premium' must have a 'capacity'(%d) of 125, 250, 500, 1000, 1500, 2000, 2500, 3000, 3500 or 4000 DTUs",
		MaxSizeGB: "service tier 'Premium' with a 'capacity' of %d must have a 'max_size_gb' no greater than %d GB, got %d GB",
	},
}

var getvCoreErrorMsg = map[string]map[string]string{
	"generalpurpose": {
		"gen4": "service tier 'GeneralPurpose' Gen4 must have a 'capacity'(%d) of 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 16 or 24 vCores",
		"gen5": "service tier 'GeneralPurpose' Gen5 must have a 'capacity'(%d) of 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40 or 80 vCores",
	},
	"businesscritical": {
		"gen4": "service tier 'BusinessCritical' Gen4 must have a 'capacity'(%d) of 2, 3, 4, 5, 6, 7, 8, 9, 10, 16 or 24 vCores",
		"gen5": "service tier 'BusinessCritical' Gen5 must have a 'capacity'(%d) of 4, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40 or 80 vCores",
	},
}

var getDTUMaxGB = map[string]map[int]float64{
	"basicpool": {
		50:   4.8828125,
		100:  9.765625,
		200:  19.53125,
		300:  29.296875,
		400:  39.0625,
		800:  78.125,
		1200: 117.1875,
		1600: 156.25,
	},
	"standardpool": {
		50:   500,
		100:  750,
		200:  1024,
		300:  1280,
		400:  1536,
		800:  2048,
		1200: 2560,
		1600: 3072,
		2000: 3584,
		2500: 4096,
		3000: 4096,
	},
	"premiumpool": {
		125:  1024,
		250:  1024,
		500:  1024,
		1000: 1024,
		1500: 1536,
		2000: 2048,
		2500: 2560,
		3000: 3072,
		3500: 3584,
		4000: 4096,
	},
}

var supportedDTUMaxGBValues = map[float64]bool{
	50:   true,
	100:  true,
	150:  true,
	200:  true,
	250:  true,
	300:  true,
	400:  true,
	500:  true,
	750:  true,
	800:  true,
	1024: true,
	1200: true,
	1280: true,
	1536: true,
	1600: true,
	1792: true,
	2000: true,
	2048: true,
	2304: true,
	2500: true,
	2560: true,
	2816: true,
	3000: true,
	3072: true,
	3328: true,
	3584: true,
	3840: true,
	4096: true,
}

var getvCoreMaxGB = map[string]map[string]map[int]float64{
	"generalpurpose": {
		"gen4": {
			1:  512,
			2:  756,
			3:  1536,
			4:  1536,
			5:  1536,
			6:  2048,
			7:  2048,
			8:  2048,
			9:  2048,
			10: 2048,
			16: 3584,
			24: 4096,
		},
		"gen5": {
			2:  512,
			4:  756,
			6:  1536,
			8:  1536,
			10: 1536,
			12: 2048,
			14: 2048,
			16: 2048,
			18: 3072,
			20: 3072,
			24: 3072,
			32: 4096,
			40: 4096,
			80: 4096,
		},
	},
	"businesscritical": {
		"gen4": {
			2:  1024,
			3:  1024,
			4:  1024,
			5:  1024,
			6:  1024,
			7:  1024,
			8:  1024,
			9:  1024,
			10: 1024,
			16: 1024,
			24: 1024,
		},
		"gen5": {
			4:  1024,
			6:  1536,
			8:  1536,
			10: 1536,
			12: 3072,
			14: 3072,
			16: 3072,
			18: 3072,
			20: 3072,
			24: 4096,
			32: 4096,
			40: 4096,
			80: 4096,
		},
	},
}

var getTierFromName = map[string]string{
	"basicpool":    "Basic",
	"standardpool": "Standard",
	"premiumpool":  "Premium",
	"gp_gen4":      "GeneralPurpose",
	"gp_gen5":      "GeneralPurpose",
	"bc_gen4":      "BusinessCritical",
	"bc_gen5":      "BusinessCritical",
}

func MSSQLElasticPoolGetDTUBasedErrorMsg(name string, errorType ErrorType) string {
	return getDTUErrorMsg[strings.ToLower(name)][errorType]
}

func MSSQLElasticPoolGetvCoreBasedErrorMsg(tier string, family string) string {
	return getvCoreErrorMsg[strings.ToLower(tier)][strings.ToLower(family)]
}

func MSSQLElasticPoolGetDTUMaxSizeGB(name string, capacity int) float64 {
	return getDTUMaxGB[strings.ToLower(name)][capacity]
}

func MSSQLElasticPoolIsValidDTUMaxGBSize(gigabytes float64) bool {
	return supportedDTUMaxGBValues[gigabytes]
}

func MSSQLElasticPoolGetvCoreMaxSizeGB(tier string, family string, vCores int) float64 {
	return getvCoreMaxGB[strings.ToLower(tier)][strings.ToLower(family)][vCores]
}

func MSSQLElasticPoolNameContainsFamily(name string, family string) bool {
	return strings.Contains(strings.ToLower(name), strings.ToLower(family))
}

func MSSQLElasticPoolNameTierIsValid(name string, tier string) bool {
	if strings.EqualFold(name, "BasicPool") && !strings.EqualFold(tier, "Basic") ||
		strings.EqualFold(name, "StandardPool") && !strings.EqualFold(tier, "Standard") ||
		strings.EqualFold(name, "PremiumPool") && !strings.EqualFold(tier, "Premium") ||
		strings.HasPrefix(strings.ToLower(name), "gp_") && !strings.EqualFold(tier, "GeneralPurpose") ||
		strings.HasPrefix(strings.ToLower(name), "bc_") && !strings.EqualFold(tier, "BusinessCritical") {
		return false
	}

	return true
}

func MSSQLElasticPoolGetTierFromSKUName(name string) string {
	return getTierFromName[strings.ToLower(name)]
}

func MSSQLElasticPoolGetFamilyFromSKUName(name string) string {
	if !strings.HasPrefix(strings.ToLower(name), "gp_") && !strings.HasPrefix(strings.ToLower(name), "bc_") {
		return ""
	}

	nameFamily := name[3:]
	retFamily := "Gen4" // Default

	if strings.EqualFold(nameFamily, "Gen5") {
		retFamily = "Gen5"
	}

	return retFamily
}
