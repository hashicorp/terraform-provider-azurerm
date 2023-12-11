// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"net"
	"strings"
)

// Verifies that there and no duplicate CIDRs in the passed input based on CIDR type (e.g. IPv4 or IPv6)
func FrontDoorRuleCidrOverlap(input []interface{}, key string) (warnings []string, errors []error) {
	// verify there are no duplicates in the CIDRs
	if len(input) > 1 {
		tmp := make(map[string]bool)
		for _, CIDR := range input {
			v, ok := CIDR.(string)
			if !ok {
				errors = append(errors, fmt.Errorf("expected %q to be a string", key))
				return warnings, errors
			}

			if _, value := tmp[v]; !value {
				tmp[v] = true
			} else {
				errors = append(errors, fmt.Errorf("%q CIDRs must be unique, there is a duplicate entry for CIDR %q in the %q field. Please remove the duplicate entry and re-apply", key, CIDR, key))
				return warnings, errors
			}
		}
	} else {
		_, ok := input[0].(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return warnings, errors
		}

		return warnings, errors
	}

	// separate the CIDRs into IPv6 and IPv4 variants
	IPv4CIDRs := make([]string, 0)
	IPv6CIDRs := make([]string, 0)

	for _, matchValue := range input {
		if matchValue != nil {
			CIDR := matchValue.(string)

			if strings.Count(CIDR, ":") < 2 {
				IPv4CIDRs = append(IPv4CIDRs, CIDR)
			}
			if strings.Count(CIDR, ":") >= 2 {
				IPv6CIDRs = append(IPv6CIDRs, CIDR)
			}
		}
	}

	// check to see if the IPv4 CIDRs overlap
	if len(IPv4CIDRs) > 1 {
		for _, sourceCIDR := range IPv4CIDRs {
			for _, checkCIDR := range IPv4CIDRs {
				if sourceCIDR == checkCIDR {
					continue
				}

				cidrOverlaps, err := validateCIDROverlap(sourceCIDR, checkCIDR)
				if err != nil {
					errors = append(errors, err)
					return warnings, errors
				}

				if cidrOverlaps {
					errors = append(errors, fmt.Errorf("the IPv4 %q CIDR %q address range overlaps with %q IPv4 CIDR address range", key, sourceCIDR, checkCIDR))
					return warnings, errors
				}
			}
		}
	}

	// check to see if the IPv6 CIDRs overlap
	if len(IPv6CIDRs) > 1 {
		for _, sourceCIDR := range IPv6CIDRs {
			for _, checkCIDR := range IPv6CIDRs {
				if sourceCIDR == checkCIDR {
					continue
				}

				cidrOverlaps, err := validateCIDROverlap(sourceCIDR, checkCIDR)
				if err != nil {
					errors = append(errors, fmt.Errorf("unable to validate IPv6 CIDR address ranges overlap: %+v", err))
					return warnings, errors
				}

				if cidrOverlaps {
					errors = append(errors, fmt.Errorf("the %q IPv6 CIDR %q address range overlaps with %q IPv6 CIDR address range", key, sourceCIDR, checkCIDR))
					return warnings, errors
				}
			}
		}
	}

	return warnings, errors
}

func validateCIDROverlap(sourceCIDR string, checkCIDR string) (bool, error) {
	_, sourceNetwork, err := net.ParseCIDR(sourceCIDR)
	if err != nil {
		return false, err
	}

	sourceOnes, sourceBits := sourceNetwork.Mask.Size()
	if sourceOnes == 0 && sourceBits == 0 {
		return false, fmt.Errorf("%q CIDR must be in its canonical form", sourceCIDR)
	}

	_, checkNetwork, err := net.ParseCIDR(checkCIDR)
	if err != nil {
		return false, err
	}

	checkOnes, checkBits := checkNetwork.Mask.Size()
	if checkOnes == 0 && checkBits == 0 {
		return false, fmt.Errorf("%q CIDR must be in its canonical form", checkCIDR)
	}

	ipStr := checkNetwork.IP.String()
	checkIp := net.ParseIP(ipStr)
	if checkIp == nil {
		return false, fmt.Errorf("unable to parse %q, invalid IP address", ipStr)
	}

	ipStr = sourceNetwork.IP.String()
	sourceIp := net.ParseIP(ipStr)
	if sourceIp == nil {
		return false, fmt.Errorf("unable to parse %q, invalid IP address", ipStr)
	}

	// swap the check values depending on which CIDR is more specific
	if sourceOnes > checkOnes {
		sourceNetwork = checkNetwork
		checkIp = sourceIp
	}

	// CIDRs overlap was detected
	if sourceNetwork.Contains(checkIp) {
		return true, nil
	}

	// CIDR overlap was not detected
	return false, nil
}
