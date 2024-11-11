// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = AssessmentMetadataId{}

func TestAssessmentMetadataIDFormatter(t *testing.T) {
	actual := NewAssessmentMetadataID("12345678-1234-9876-4563-123456789012", "metadata1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/assessmentMetadata/metadata1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestAssessmentMetadataID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AssessmentMetadataId
	}{
		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing AssessmentMetadataName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/",
			Error: true,
		},

		{
			// missing value for AssessmentMetadataName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/assessmentMetadata/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/assessmentMetadata/metadata1",
			Expected: &AssessmentMetadataId{
				SubscriptionId:         "12345678-1234-9876-4563-123456789012",
				AssessmentMetadataName: "metadata1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.SECURITY/ASSESSMENTMETADATA/METADATA1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := AssessmentMetadataID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.AssessmentMetadataName != v.Expected.AssessmentMetadataName {
			t.Fatalf("Expected %q but got %q for AssessmentMetadataName", v.Expected.AssessmentMetadataName, actual.AssessmentMetadataName)
		}
	}
}
