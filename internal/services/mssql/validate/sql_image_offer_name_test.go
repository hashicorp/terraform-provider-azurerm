// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestSqlImageOfferName(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		// Valid cases
		{
			Value:    "SQL2019-WS2019",
			ErrCount: 0,
		},
		{
			Value:    "SQL2022-WS2012R2",
			ErrCount: 0,
		},
		// Invalid Cases
		{
			Value:    "",
			ErrCount: 1,
		},
		{
			Value:    "2019-WS2019",
			ErrCount: 1,
		},
		{
			Value:    "SQL2019-2019",
			ErrCount: 1,
		},
		{
			Value:    "SQL2019-WS20.19",
			ErrCount: 1,
		},
		{
			Value:    "SQL20.19-WS2019",
			ErrCount: 1,
		},
		{
			Value:    "SQL2019.WS2019",
			ErrCount: 1,
		},
	}

	for i, tc := range cases {
		_, errors := SqlImageOfferName(tc.Value, "sql_image_offer")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Case %d: Encountered %d error(s), expected %d", i, len(errors), tc.ErrCount)
		}
	}
}
