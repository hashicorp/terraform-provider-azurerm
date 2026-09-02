// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestPipelineGroupDataCollectionRuleImmutableId(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{Input: "", Valid: false},
		{Input: "dcr-02f5b9b754b945d4b6bcdc6688d7e04f", Valid: true},
		{Input: "dcr-02F5B9B754B945D4B6BCDC6688D7E04F", Valid: true},
		{Input: "02f5b9b754b945d4b6bcdc6688d7e04f", Valid: false},
		{Input: "dcr-02f5b9b754b945d4b6bcdc6688d7e04", Valid: false},
		{Input: "dcr-02f5b9b754b945d4b6bcdc6688d7e04ff", Valid: false},
		{Input: "dcr-02f5b9b7-54b9-45d4-b6bc-dc6688d7e04f", Valid: false},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			_, errors := DataCollectionRuleImmutableId(tc.Input, "data_collection_rule_immutable_id")
			valid := len(errors) == 0
			if valid != tc.Valid {
				t.Fatalf("expected %q to be valid=%t, got valid=%t", tc.Input, tc.Valid, valid)
			}
		})
	}
}
