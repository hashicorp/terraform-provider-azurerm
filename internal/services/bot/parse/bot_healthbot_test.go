// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = BotHealthbotId{}

func TestBotHealthbotIDFormatter(t *testing.T) {
	actual := NewBotHealthbotID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "bot1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.HealthBot/healthBots/bot1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestBotHealthbotID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *BotHealthbotId
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
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing HealthBotName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.HealthBot/",
			Error: true,
		},

		{
			// missing value for HealthBotName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.HealthBot/healthBots/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.HealthBot/healthBots/bot1",
			Expected: &BotHealthbotId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resourceGroup1",
				HealthBotName:  "bot1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.HEALTHBOT/HEALTHBOTS/BOT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := BotHealthbotID(v.Input)
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
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.HealthBotName != v.Expected.HealthBotName {
			t.Fatalf("Expected %q but got %q for HealthBotName", v.Expected.HealthBotName, actual.HealthBotName)
		}
	}
}
