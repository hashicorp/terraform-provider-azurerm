package azure

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type ErrorType int
type SkuType int

const (
	CapacityError ErrorType = 0
	AllOtherError ErrorType = 1
)

const (
	DTU   SkuType = 0
	VCore SkuType = 1
)

var getDTUMaxGB = map[string]map[int]float64{
	"basic": {
		50:   4.8828125,
		100:  9.765625,
		200:  19.53125,
		300:  29.296875,
		400:  39.0625,
		800:  78.125,
		1200: 117.1875,
		1600: 156.25,
	},
	"standard": {
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
	"premium": {
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

func MSSQLElasticPoolValidateSKU(diff *schema.ResourceDiff) error {

	name, _ := diff.GetOk("sku.0.name")
	tier, _ := diff.GetOk("sku.0.tier")
	capacity, _ := diff.GetOk("sku.0.capacity")
	family, _ := diff.GetOk("sku.0.family")
	maxSizeBytes, _ := diff.GetOk("max_size_bytes")
	maxSizeGb, _ := diff.GetOk("max_size_gb")
	minCapacity, _ := diff.GetOk("per_database_settings.0.min_capacity")
	maxCapacity, _ := diff.GetOk("per_database_settings.0.max_capacity")

	skuToValidate := DTU

	if strings.HasPrefix(strings.ToLower(name.(string)), "gp_") || strings.HasPrefix(strings.ToLower(name.(string)), "bc_") {
		skuToValidate = VCore
	}

	// Convert Bytes to Gigabytes only if max_size_gb
	// has not changed else use max_size_gb
	if maxSizeBytes != 0 && !diff.HasChange("max_size_gb") {
		maxSizeGb = float64(maxSizeBytes.(int) / 1024 / 1024 / 1024)
	}

	if skuToValidate == DTU {
		// DTU Checks
		maxAllowedGB := getDTUMaxGB[strings.ToLower(getTierFromName[strings.ToLower(name.(string))])][capacity.(int)]

		// Check to see if this is a valid SKU capacity combo
		if errMsg := getErrorMsg(name.(string), family.(string), capacity.(int), maxAllowedGB, maxSizeGb.(float64), minCapacity.(float64), maxCapacity.(float64), CapacityError, skuToValidate); errMsg != "" {
			return fmt.Errorf(errMsg)
		}

		// Check to see if this is a valid Max_Size_GB value
		if errMsg := getErrorMsg(name.(string), family.(string), capacity.(int), maxAllowedGB, maxSizeGb.(float64), minCapacity.(float64), maxCapacity.(float64), AllOtherError, skuToValidate); errMsg != "" {
			return fmt.Errorf(errMsg)
		}

	} else {
		// vCore Checks
		maxAllowedGB := getvCoreMaxGB[strings.ToLower(getTierFromName[strings.ToLower(name.(string))])][strings.ToLower(getFamilyFromName(strings.ToLower(name.(string))))][capacity.(int)]

		if errMsg := getErrorMsg(name.(string), family.(string), capacity.(int), maxAllowedGB, maxSizeGb.(float64), minCapacity.(float64), maxCapacity.(float64), CapacityError, skuToValidate); errMsg != "" {
			return fmt.Errorf(errMsg)
		}

		// Check to see if this is a valid Max_Size_GB value
		if errMsg := getErrorMsg(name.(string), family.(string), capacity.(int), maxAllowedGB, maxSizeGb.(float64), minCapacity.(float64), maxCapacity.(float64), AllOtherError, skuToValidate); errMsg != "" {
			return fmt.Errorf(errMsg)
		}
	}

	// Universal check for all SKUs
	if !nameTierIsValid(name.(string), tier.(string)) {
		return fmt.Errorf("Mismatch between SKU name '%s' and tier '%s', expected 'tier' to be '%s'", name.(string), tier.(string), getTierFromName[strings.ToLower(name.(string))])
	}

	return nil
}

func nameContainsFamily(name string, family string) bool {
	return strings.Contains(strings.ToLower(name), strings.ToLower(family))
}

func nameTierIsValid(name string, tier string) bool {
	if strings.EqualFold(name, "BasicPool") && !strings.EqualFold(tier, "Basic") ||
		strings.EqualFold(name, "StandardPool") && !strings.EqualFold(tier, "Standard") ||
		strings.EqualFold(name, "PremiumPool") && !strings.EqualFold(tier, "Premium") ||
		strings.HasPrefix(strings.ToLower(name), "gp_") && !strings.EqualFold(tier, "GeneralPurpose") ||
		strings.HasPrefix(strings.ToLower(name), "bc_") && !strings.EqualFold(tier, "BusinessCritical") {
		return false
	}

	return true
}

func getFamilyFromName(name string) string {
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

func getDTUCapacityErrorMsg(name string, capacity int) string {
	m := getDTUMaxGB[strings.ToLower(getTierFromName[strings.ToLower(name)])]
	return buildCapacityErrorString(name, capacity, m)
}

func getVCoreCapacityErrorMsg(name string, capacity int) string {
	m := getvCoreMaxGB[strings.ToLower(getTierFromName[strings.ToLower(name)])][strings.ToLower(getFamilyFromName(name))]
	return buildCapacityErrorString(name, capacity, m)
}

func buildCapacityErrorString(name string, capacity int, m map[int]float64) string {
	var a []int
	str := ""
	family := getFamilyFromName(name)
	tier := getTierFromName[strings.ToLower(name)]

	if family == "" {
		str += fmt.Sprintf("service tier '%s' must have a 'capacity'(%d) of ", tier, capacity)
	} else {
		str += fmt.Sprintf("service tier '%s' %s must have a 'capacity'(%d) of ", tier, family, capacity)
	}

	// copy the keys into another map
	possible := make([]int, 0, len(m))
	for k := range m {
		possible = append(possible, k)
	}

	// copy the values of the map of keys into a slice of ints
	for v := range possible {
		a = append(a, possible[v])
	}

	// sort the slice to get them in order
	sort.Sort(sort.IntSlice(a))

	// build the error message
	for i := range a {
		if i < len(a)-1 {
			str += fmt.Sprintf("%d, ", a[i])
		} else {
			if family == "" {
				str += fmt.Sprintf("or %d DTUs", a[i])
			} else {
				str += fmt.Sprintf("or %d vCores", a[i])
			}
		}
	}

	return str
}

func getDTUNotValidSizeErrorMsg(name string, maxSizeGb float64) string {
	m := supportedDTUMaxGBValues
	var a []int
	str := fmt.Sprintf("'max_size_gb'(%d) is not a valid value for service tier '%s', 'max_size_gb' must have a value of ", int(maxSizeGb), getTierFromName[strings.ToLower(name)])

	// copy the keys into another map
	possible := make([]int, 0, len(m))
	for k := range m {
		possible = append(possible, int(k))
	}

	// copy the values of the map of keys into a slice of ints
	for v := range possible {
		a = append(a, possible[v])
	}

	// sort the slice to get them in order
	sort.Sort(sort.IntSlice(a))

	// build the error message
	for index := range a {

		if index < len(a)-1 {
			str += fmt.Sprintf("%d, ", a[index])
		} else {
			str += fmt.Sprintf("or %d GB", a[index])
		}
	}

	return str
}

func getErrorMsg(name string, family string, capacity int, maxAllowedGB float64, maxSizeGb float64, minCapacity float64, maxCapacity float64, errType ErrorType, skuType SkuType) string {
	errMsg := ""

	if errType == CapacityError {
		if skuType == DTU {
			// DTU Based Capacity Errors
			if maxAllowedGB == 0 {
				return getDTUCapacityErrorMsg(name, capacity)
			}
		} else {
			// vCore Based Capacity Errors
			if maxAllowedGB == 0 {
				return getVCoreCapacityErrorMsg(name, capacity)
			}
		}
	} else {
		// AllOther Errors
		if skuType == DTU {
			if strings.EqualFold(name, "BasicPool") {
				// Basic SKU does not let you pick your max_size_GB they are fixed values
				if maxSizeGb != maxAllowedGB {
					return fmt.Sprintf("service tier 'Basic' with a 'capacity' of %d must have a 'max_size_gb' of %.7f GB, got %.7f GB", capacity, maxAllowedGB, maxSizeGb)
				}
			} else {
				// All other DTU based SKUs
				if maxSizeGb > maxAllowedGB {
					return fmt.Sprintf("service tier '%s' with a 'capacity' of %d must have a 'max_size_gb' no greater than %d GB, got %d GB", getTierFromName[strings.ToLower(name)], capacity, int(maxAllowedGB), int(maxSizeGb))
				}

				if int(maxSizeGb) < 50 {
					return fmt.Sprintf("service tier '%s', must have a 'max_size_gb' value equal to or greater than 50 GB, got %d GB", getTierFromName[strings.ToLower(name)], int(maxSizeGb))
				}

				// Check to see if the max_size_gb value is valid for this SKU type and capacity
				if !supportedDTUMaxGBValues[maxSizeGb] {
					return getDTUNotValidSizeErrorMsg(name, maxSizeGb)
				}
			}

			// All Other DTU based SKU Checks
			if family != "" {
				return fmt.Sprintf("Invalid attribute 'family' (%s) for service tiers 'Basic', 'Standard', and 'Premium', remove the 'family' attribute from the configuration file", family)
			}
			if minCapacity != math.Trunc(minCapacity) {
				return fmt.Sprintf("service tiers 'Basic', 'Standard', and 'Premium' must have whole numbers as their 'minCapacity'")
			}

			if maxCapacity != math.Trunc(maxCapacity) {
				return fmt.Sprintf("service tiers 'Basic', 'Standard', and 'Premium' must have whole numbers as their 'maxCapacity'")
			}

			if minCapacity < 0.0 {
				return fmt.Sprintf("service tiers 'Basic', 'Standard', and 'Premium' per_database_settings 'min_capacity' must be equal to or greater than zero")
			}
		} else {
			// vCore Based AllOther Errors
			if maxSizeGb > maxAllowedGB {
				return fmt.Sprintf("service tier '%s' %s with a 'capacity' of %d vCores must have a 'max_size_gb' between 5 GB and %d GB, got %d GB", getTierFromName[strings.ToLower(name)], family, capacity, int(maxAllowedGB), int(maxSizeGb))
			}

			if int(maxSizeGb) < 5 {
				return fmt.Sprintf("service tier '%s' must have a 'max_size_gb' value equal to or greater than 5 GB, got %d GB", getTierFromName[strings.ToLower(name)], int(maxSizeGb))
			}

			if maxSizeGb != math.Trunc(maxSizeGb) {
				return fmt.Sprintf("'max_size_gb' must be a whole number, got %f GB", maxSizeGb)
			}

			if !nameContainsFamily(name, family) {
				return fmt.Sprintf("Mismatch between SKU name '%s' and family '%s', expected '%s'", name, family, getFamilyFromName(name))
			}

			if maxCapacity > float64(capacity) {
				return fmt.Sprintf("service tier '%s' perDatabaseSettings 'maxCapacity' must not be higher than the SKUs 'capacity'(%d) value", getTierFromName[strings.ToLower(name)], capacity)
			}

			if minCapacity > maxCapacity {
				return fmt.Sprintf("perDatabaseSettings 'maxCapacity' must be greater than or equal to the perDatabaseSettings 'minCapacity' value")
			}
		}
	}

	return errMsg
}
