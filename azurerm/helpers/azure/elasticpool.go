package azure

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type errorType int
type skuType int

const (
	capacityError errorType = 0
	allOtherError errorType = 1
)

const (
	DTU   skuType = 0
	VCore skuType = 1
)

type sku struct {
	Name, Tier, Family                                string
	Capacity                                          int
	MaxAllowedGB, MaxSizeGb, MinCapacity, MaxCapacity float64
	ErrorType                                         errorType
	SkuType                                           skuType
}

// getDTUMaxGB: this map holds all of the DTU to 'max_size_gb' mappings based on a DTU lookup
//              note that the value can be below the returned value, except for 'basic' it's
//              value must match exactly what is returned else it will be rejected by the API
//              which will return a 'Internal Server Error'

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

// supportedDTUMaxGBValues: this map holds all of the valid 'max_size_gb' values
//                          for a DTU SKU type. If the 'max_size_gb' is anything
//                          other than the values in the map the API with throw
//                          an 'Internal Server Error'

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

// getvCoreMaxGB: this map holds all of the vCore to 'max_size_gb' mappings based on a vCore lookup
//                note that the value can be below the returned value

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

// getTierFromName: this map contains all of the valid mappings between 'name' and 'tier'
//                  the reason for this map is that the user may pass in an invalid mapping
//                  (e.g. name: "Basicpool" tier:"BusinessCritical") this map allows me
//                  to lookup the correct values in other maps even if the config file
//                  contains an invalid 'tier' attribute.

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

	s := sku{
		Name:        name.(string),
		Tier:        tier.(string),
		Family:      family.(string),
		Capacity:    capacity.(int),
		MaxSizeGb:   maxSizeGb.(float64),
		MinCapacity: minCapacity.(float64),
		MaxCapacity: maxCapacity.(float64),
		ErrorType:   capacityError,
		SkuType:     DTU,
	}

	if strings.HasPrefix(strings.ToLower(name.(string)), "gp_") || strings.HasPrefix(strings.ToLower(name.(string)), "bc_") {
		s.SkuType = VCore
	}

	// Convert Bytes to Gigabytes only if max_size_gb
	// has not changed else use max_size_gb
	if maxSizeBytes.(int) != 0 && !diff.HasChange("max_size_gb") {
		s.MaxSizeGb = float64(maxSizeBytes.(int) / 1024 / 1024 / 1024)
	}

	if s.SkuType == DTU {
		// DTU Checks
		s.MaxAllowedGB = getDTUMaxGB[strings.ToLower(getTierFromName[strings.ToLower(s.Name)])][s.Capacity]

		// Check to see if this is a valid SKU capacity combo
		if errMsg := doSKUValidation(s); errMsg != "" {
			return fmt.Errorf(errMsg)
		}

		// Check to see if this is a valid Max_Size_GB value
		s.ErrorType = allOtherError
		if errMsg := doSKUValidation(s); errMsg != "" {
			return fmt.Errorf(errMsg)
		}

	} else {
		// vCore Checks
		s.MaxAllowedGB = getvCoreMaxGB[strings.ToLower(getTierFromName[strings.ToLower(s.Name)])][strings.ToLower(getFamilyFromName(strings.ToLower(s.Name)))][s.Capacity]

		if errMsg := doSKUValidation(s); errMsg != "" {
			return fmt.Errorf(errMsg)
		}

		// Check to see if this is a valid Max_Size_GB value
		s.ErrorType = allOtherError
		if errMsg := doSKUValidation(s); errMsg != "" {
			return fmt.Errorf(errMsg)
		}
	}

	// Universal check for all SKUs
	if !nameTierIsValid(s.Name, s.Tier) {
		return fmt.Errorf("Mismatch between SKU name '%s' and tier '%s', expected 'tier' to be '%s'", s.Name, s.Tier, getTierFromName[strings.ToLower(s.Name)])
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
	return buildcapacityErrorString(name, capacity, m)
}

func getVCoreCapacityErrorMsg(name string, capacity int) string {
	m := getvCoreMaxGB[strings.ToLower(getTierFromName[strings.ToLower(name)])][strings.ToLower(getFamilyFromName(name))]
	return buildcapacityErrorString(name, capacity, m)
}

func buildcapacityErrorString(name string, capacity int, m map[int]float64) string {
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
	p := make([]int, 0, len(m))
	for k := range m {
		p = append(p, k)
	}

	// copy the values of the map of keys into a slice of ints
	for v := range p {
		a = append(a, p[v])
	}

	// sort the slice to get them in order
	sort.Ints(a)

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
	p := make([]int, 0, len(m))
	for k := range m {
		p = append(p, int(k))
	}

	// copy the values of the map of keys into a slice of ints
	for v := range p {
		a = append(a, p[v])
	}

	// sort the slice to get them in order
	sort.Ints(a)

	// build the error message
	for i := range a {

		if i < len(a)-1 {
			str += fmt.Sprintf("%d, ", a[i])
		} else {
			str += fmt.Sprintf("or %d GB", a[i])
		}
	}

	return str
}

func doSKUValidation(s sku) string {
	errMsg := ""

	if s.ErrorType == capacityError {
		if s.SkuType == DTU {
			// DTU Based Capacity Errors
			if s.MaxAllowedGB == 0 {
				return getDTUCapacityErrorMsg(s.Name, s.Capacity)
			}
		} else {
			// vCore Based Capacity Errors
			if s.MaxAllowedGB == 0 {
				return getVCoreCapacityErrorMsg(s.Name, s.Capacity)
			}
		}
	} else {
		// AllOther Errors
		if s.SkuType == DTU {
			if strings.EqualFold(s.Name, "BasicPool") {
				// Basic SKU does not let you pick your max_size_GB they are fixed values
				if s.MaxSizeGb != s.MaxAllowedGB {
					return fmt.Sprintf("service tier 'Basic' with a 'capacity' of %d must have a 'max_size_gb' of %.7f GB, got %.7f GB", s.Capacity, s.MaxAllowedGB, s.MaxSizeGb)
				}
			} else {
				// All other DTU based SKUs
				if s.MaxSizeGb > s.MaxAllowedGB {
					return fmt.Sprintf("service tier '%s' with a 'capacity' of %d must have a 'max_size_gb' no greater than %d GB, got %d GB", getTierFromName[strings.ToLower(s.Name)], s.Capacity, int(s.MaxAllowedGB), int(s.MaxSizeGb))
				}

				if int(s.MaxSizeGb) < 50 {
					return fmt.Sprintf("service tier '%s', must have a 'max_size_gb' value equal to or greater than 50 GB, got %d GB", getTierFromName[strings.ToLower(s.Name)], int(s.MaxSizeGb))
				}

				// Check to see if the max_size_gb value is valid for this SKU type and capacity
				if !supportedDTUMaxGBValues[s.MaxSizeGb] {
					return getDTUNotValidSizeErrorMsg(s.Name, s.MaxSizeGb)
				}
			}

			// All Other DTU based SKU Checks
			if s.Family != "" {
				return fmt.Sprintf("Invalid attribute 'family' (%s) for service tiers 'Basic', 'Standard', and 'Premium', remove the 'family' attribute from the configuration file", s.Family)
			}
			if s.MinCapacity != math.Trunc(s.MinCapacity) {
				return fmt.Sprintf("service tiers 'Basic', 'Standard', and 'Premium' must have whole numbers as their 'minCapacity'")
			}

			if s.MaxCapacity != math.Trunc(s.MaxCapacity) {
				return fmt.Sprintf("service tiers 'Basic', 'Standard', and 'Premium' must have whole numbers as their 'maxCapacity'")
			}

			if s.MinCapacity < 0.0 {
				return fmt.Sprintf("service tiers 'Basic', 'Standard', and 'Premium' per_database_settings 'min_capacity' must be equal to or greater than zero")
			}
		} else {
			// vCore Based Errors
			if s.MaxSizeGb > s.MaxAllowedGB {
				return fmt.Sprintf("service tier '%s' %s with a 'capacity' of %d vCores must have a 'max_size_gb' between 5 GB and %d GB, got %d GB", getTierFromName[strings.ToLower(s.Name)], s.Family, s.Capacity, int(s.MaxAllowedGB), int(s.MaxSizeGb))
			}

			if int(s.MaxSizeGb) < 5 {
				return fmt.Sprintf("service tier '%s' must have a 'max_size_gb' value equal to or greater than 5 GB, got %d GB", getTierFromName[strings.ToLower(s.Name)], int(s.MaxSizeGb))
			}

			if s.MaxSizeGb != math.Trunc(s.MaxSizeGb) {
				return fmt.Sprintf("'max_size_gb' must be a whole number, got %f GB", s.MaxSizeGb)
			}

			if !nameContainsFamily(s.Name, s.Family) {
				return fmt.Sprintf("Mismatch between SKU name '%s' and family '%s', expected '%s'", s.Name, s.Family, getFamilyFromName(s.Name))
			}

			if s.MaxCapacity > float64(s.Capacity) {
				return fmt.Sprintf("service tier '%s' perDatabaseSettings 'maxCapacity' must not be higher than the SKUs 'capacity'(%d) value", getTierFromName[strings.ToLower(s.Name)], s.Capacity)
			}

			if s.MinCapacity > s.MaxCapacity {
				return fmt.Sprintf("perDatabaseSettings 'maxCapacity' must be greater than or equal to the perDatabaseSettings 'minCapacity' value")
			}
		}
	}

	return errMsg
}
