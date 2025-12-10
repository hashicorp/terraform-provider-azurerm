// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"math/big"
)

func NumberOfIpAddresses(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string containing a positive integer", key))
		return warnings, errors
	}

	n, ok := new(big.Int).SetString(v, 10)
	if !ok || n.Sign() <= 0 || n.BitLen() > 128 {
		errors = append(errors, fmt.Errorf("expected %q to be a positive integer", key))
		return warnings, errors
	}

	// Check if n & (n - 1) == 0
	if new(big.Int).And(n, new(big.Int).Sub(n, big.NewInt(1))).Cmp(big.NewInt(0)) != 0 {
		// Further check if it fits within 128 bits
		errors = append(errors, fmt.Errorf("expected %q to be a power of 2", key))
	}

	return warnings, errors
}
