package azure

import (
	"strings"
)

var basicDTUMaxGB = map[int]float64{
	50:   4.8828125,
	100:  9.765625,
	200:  19.53125,
	300:  29.296875,
	400:  39.0625,
	800:  78.125,
	1200: 117.1875,
	1600: 156.25,
}

var standardDTUMaxGB = map[int]float64{
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
}

var premiumDTUMaxGB = map[int]float64{
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

var generalPurposeGen4vCoreMaxGB = map[int]float64{
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
}

var generalPurposeGen5vCoreMaxGB = map[int]float64{
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
}

var businessCriticalGen4vCoreMaxGB = map[int]float64{
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
}

var businessCriticalGen5vCoreMaxGB = map[int]float64{
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
}

func BasicGetMaxSizeGB(DTUs int) float64 {
	return basicDTUMaxGB[DTUs]
}

func StandardGetMaxSizeGB(DTUs int) float64 {
	return standardDTUMaxGB[DTUs]
}

func PremiumGetMaxSizeGB(DTUs int) float64 {
	return premiumDTUMaxGB[DTUs]
}

func StandardPremiumMaxGBValid(gb float64) bool {
	return supportedMaxGBValues[gb]
}

func GeneralPurposeGetMaxSizeGB(vCores int, family string) float64 {
	maxGB := 0.0

	if strings.EqualFold(family, "Gen4") {
		maxGB = generalPurposeGen4vCoreMaxGB[vCores]
	}

	if strings.EqualFold(family, "Gen5") {
		maxGB = generalPurposeGen5vCoreMaxGB[vCores]
	}

	return maxGB
}

func BusinessCriticalGetMaxSizeGB(vCores int, family string) float64 {
	maxGB := 0.0

	if strings.EqualFold(family, "Gen4") {
		maxGB = businessCriticalGen4vCoreMaxGB[vCores]
	}

	if strings.EqualFold(family, "Gen5") {
		maxGB = businessCriticalGen5vCoreMaxGB[vCores]
	}

	return maxGB
}

func NameFamilyValid(name string, family string) bool {
	return strings.Contains(strings.ToLower(name), strings.ToLower(family))
}

func GetTier(name string) string {
	validTier := ""

	if strings.EqualFold(name, "BasicPool") {
		validTier = "Basic"
	}

	if strings.EqualFold(name, "StandardPool") {
		validTier = "Standard"
	}

	if strings.EqualFold(name, "PremiumPool") {
		validTier = "Premium"
	}

	if strings.HasPrefix(strings.ToLower(name), "gp_") {
		validTier = "GeneralPurpose"
	}

	if strings.HasPrefix(strings.ToLower(name), "bc_") {
		validTier = "BusinessCritical"
	}

	return validTier
}
