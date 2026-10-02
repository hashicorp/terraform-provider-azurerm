// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package billing_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/billing"
)

func TestIsReturnNotAllowed(t *testing.T) {
	testData := []struct {
		Name     string
		Err      error
		Expected bool
	}{
		{
			Name:     "ReturnFailed code",
			Err:      fmt.Errorf("ReturnFailed: the reservation could not be returned"),
			Expected: true,
		},
		{
			Name:     "ReturnNotAllowed code",
			Err:      fmt.Errorf("ReturnNotAllowed: this reservation type cannot be returned"),
			Expected: true,
		},
		{
			Name:     "cannot be returned phrase",
			Err:      fmt.Errorf("the reservation cannot be returned after the refund window"),
			Expected: true,
		},
		{
			Name:     "refund period phrase",
			Err:      fmt.Errorf("the refund period for this reservation has expired"),
			Expected: true,
		},
		{
			Name:     "is not eligible phrase",
			Err:      fmt.Errorf("the reservation is not eligible for return"),
			Expected: true,
		},
		{
			Name:     "keywords are matched case-insensitively",
			Err:      fmt.Errorf("RETURNFAILED due to policy"),
			Expected: true,
		},
		{
			Name:     "unrelated error returns false",
			Err:      fmt.Errorf("internal server error"),
			Expected: false,
		},
		{
			Name:     "empty error message returns false",
			Err:      fmt.Errorf(""),
			Expected: false,
		},
	}

	for _, v := range testData {
		t.Run(v.Name, func(t *testing.T) {
			actual := billing.IsReturnNotAllowed(v.Err)
			if actual != v.Expected {
				t.Fatalf("IsReturnNotAllowed(%q): expected %v but got %v", v.Err.Error(), v.Expected, actual)
			}
		})
	}
}
