package azure

import (
	"fmt"
	"strings"
)

type ErrorType int

const (
	Capacity  ErrorType = 0
	MaxSizeGB ErrorType = 1
)

var getDTUErrorMsg = map[string]string{
	"basicpool_capacity":     "service tier 'Basic' must have a 'capacity'(%d) of 50, 100, 200, 300, 400, 800, 1200 or 1600 DTUs",
	"basicpool_maxSizeGB":    "service tier 'Basic' with a 'capacity' of %d must have a 'max_size_gb' of %.7f GB, got %.7f GB",
	"standardpool_capacity":  "service tier 'Standard' must have a 'capacity'(%d) of 50, 100, 200, 300, 400, 800, 1200, 1600, 2000, 2500 or 3000 DTUs",
	"standardpool_maxSizeGB": "service tier 'Standard' with a 'capacity' of %d must have a 'max_size_gb' no greater than %d GB, got %d GB",
	"premiumpool_capacity":   "service tier 'Premium' must have a 'capacity'(%d) of 125, 250, 500, 1000, 1500, 2000, 2500, 3000, 3500 or 4000 DTUs",
	"premiumpool_maxSizeGB":  "service tier 'Premium' with a 'capacity' of %d must have a 'max_size_gb' no greater than %d GB, got %d GB",
}

var getvCoreErrorMsg = map[string]string{
	"generalpurpose_gen4_capacity":   "service tier 'GeneralPurpose' Gen4 must have a 'capacity'(%d) of 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 16 or 24 vCores",
	"generalpurpose_gen5_capacity":   "service tier 'GeneralPurpose' Gen5 must have a 'capacity'(%d) of 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40 or 80 vCores",
	"businesscritical_gen4_capacity": "service tier 'BusinessCritical' Gen4 must have a 'capacity'(%d) of 2, 3, 4, 5, 6, 7, 8, 9, 10, 16 or 24 vCores",
	"businesscritical_gen5_capacity": "service tier 'BusinessCritical' Gen5 must have a 'capacity'(%d) of 4, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40 or 80 vCores",
}

var getDTUMaxGB = map[string]float64{
	"basicpool_50":      4.8828125,
	"basicpool_100":     9.765625,
	"basicpool_200":     19.53125,
	"basicpool_300":     29.296875,
	"basicpool_400":     39.0625,
	"basicpool_800":     78.125,
	"basicpool_1200":    117.1875,
	"basicpool_1600":    156.25,
	"standardpool_50":   500,
	"standardpool_100":  750,
	"standardpool_200":  1024,
	"standardpool_300":  1280,
	"standardpool_400":  1536,
	"standardpool_800":  2048,
	"standardpool_1200": 2560,
	"standardpool_1600": 3072,
	"standardpool_2000": 3584,
	"standardpool_2500": 4096,
	"standardpool_3000": 4096,
	"premiumpool_125":   1024,
	"premiumpool_250":   1024,
	"premiumpool_500":   1024,
	"premiumpool_1000":  1024,
	"premiumpool_1500":  1536,
	"premiumpool_2000":  2048,
	"premiumpool_2500":  2560,
	"premiumpool_3000":  3072,
	"premiumpool_3500":  3584,
	"premiumpool_4000":  4096,
}

var supportedMaxGBValues = map[float64]bool{
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

var getvCoreMaxGB = map[string]float64{
	"generalpurpose_gen4_1":    512,
	"generalpurpose_gen4_2":    756,
	"generalpurpose_gen4_3":    1536,
	"generalpurpose_gen4_4":    1536,
	"generalpurpose_gen4_5":    1536,
	"generalpurpose_gen4_6":    2048,
	"generalpurpose_gen4_7":    2048,
	"generalpurpose_gen4_8":    2048,
	"generalpurpose_gen4_9":    2048,
	"generalpurpose_gen4_10":   2048,
	"generalpurpose_gen4_16":   3584,
	"generalpurpose_gen4_24":   4096,
	"generalpurpose_gen5_2":    512,
	"generalpurpose_gen5_4":    756,
	"generalpurpose_gen5_6":    1536,
	"generalpurpose_gen5_8":    1536,
	"generalpurpose_gen5_10":   1536,
	"generalpurpose_gen5_12":   2048,
	"generalpurpose_gen5_14":   2048,
	"generalpurpose_gen5_16":   2048,
	"generalpurpose_gen5_18":   3072,
	"generalpurpose_gen5_20":   3072,
	"generalpurpose_gen5_24":   3072,
	"generalpurpose_gen5_32":   4096,
	"generalpurpose_gen5_40":   4096,
	"generalpurpose_gen5_80":   4096,
	"businesscritical_gen4_2":  1024,
	"businesscritical_gen4_3":  1024,
	"businesscritical_gen4_4":  1024,
	"businesscritical_gen4_5":  1024,
	"businesscritical_gen4_6":  1024,
	"businesscritical_gen4_7":  1024,
	"businesscritical_gen4_8":  1024,
	"businesscritical_gen4_9":  1024,
	"businesscritical_gen4_10": 1024,
	"businesscritical_gen4_16": 1024,
	"businesscritical_gen4_24": 1024,
	"businesscritical_gen5_4":  1024,
	"businesscritical_gen5_6":  1536,
	"businesscritical_gen5_8":  1536,
	"businesscritical_gen5_10": 1536,
	"businesscritical_gen5_12": 3072,
	"businesscritical_gen5_14": 3072,
	"businesscritical_gen5_16": 3072,
	"businesscritical_gen5_18": 3072,
	"businesscritical_gen5_20": 3072,
	"businesscritical_gen5_24": 4096,
	"businesscritical_gen5_32": 4096,
	"businesscritical_gen5_40": 4096,
	"businesscritical_gen5_80": 4096,
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
	errType := "capacity"

	if errorType == MaxSizeGB {
		errType = "maxSizeGB"
	}

	return getDTUErrorMsg[fmt.Sprintf("%s_%s", strings.ToLower(name), errType)]
}

func MSSQLElasticPoolGetvCoreBasedErrorMsg(tier string, family string) string {
	return getvCoreErrorMsg[fmt.Sprintf("%s_%s_capacity", strings.ToLower(tier), strings.ToLower(family))]
}

func MSSQLElasticPoolGetDTUMaxSizeGB(name string, capacity int) float64 {
	return getDTUMaxGB[fmt.Sprintf("%s_%d", strings.ToLower(name), capacity)]
}

func MSSQLElasticPoolIsValidMaxGBSizeForSKU(gb float64) bool {
	return supportedMaxGBValues[gb]
}

func MSSQLElasticPoolGetvCoreMaxSizeGB(tier string, family string, vCores int) float64 {
	return getvCoreMaxGB[fmt.Sprintf("%s_%s_%d", strings.ToLower(tier), strings.ToLower(family), vCores)]
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
