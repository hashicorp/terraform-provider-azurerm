package databaselink_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
)

func TestForceLinkNeeded(t *testing.T) {
	testCases := []struct {
		name        string
		oldItemList []interface{}
		newItemList []interface{}
		expected    bool
		description string
	}{
		{
			name:        "empty lists",
			oldItemList: []interface{}{},
			newItemList: []interface{}{},
			expected:    false,
			description: "both lists empty should return false",
		},
		{
			name:        "empty old list with new items",
			oldItemList: []interface{}{},
			newItemList: []interface{}{"item1", "item2"},
			expected:    true,
			description: "new items added to empty old list should return true",
		},
		{
			name:        "empty new list with old items",
			oldItemList: []interface{}{"item1", "item2"},
			newItemList: []interface{}{},
			expected:    false,
			description: "removing all items should return false (no new items to link)",
		},
		{
			name:        "identical lists",
			oldItemList: []interface{}{"item1", "item2", "item3"},
			newItemList: []interface{}{"item1", "item2", "item3"},
			expected:    false,
			description: "identical lists should return false",
		},
		{
			name:        "identical lists different order",
			oldItemList: []interface{}{"item1", "item2", "item3"},
			newItemList: []interface{}{"item3", "item1", "item2"},
			expected:    false,
			description: "same items in different order should return false",
		},
		{
			name:        "new list has additional items",
			oldItemList: []interface{}{"item1", "item2"},
			newItemList: []interface{}{"item1", "item2", "item3"},
			expected:    true,
			description: "adding new items should return true",
		},
		{
			name:        "new list has some different items",
			oldItemList: []interface{}{"item1", "item2"},
			newItemList: []interface{}{"item1", "item3"},
			expected:    true,
			description: "replacing some items should return true (item3 is new)",
		},
		{
			name:        "new list is subset of old list",
			oldItemList: []interface{}{"item1", "item2", "item3"},
			newItemList: []interface{}{"item1", "item2"},
			expected:    false,
			description: "removing items should return false (no new items to link)",
		},
		{
			name:        "completely different lists",
			oldItemList: []interface{}{"item1", "item2"},
			newItemList: []interface{}{"item3", "item4"},
			expected:    true,
			description: "completely different items should return true",
		},
		{
			name:        "single item added",
			oldItemList: []interface{}{"database1"},
			newItemList: []interface{}{"database1", "database2"},
			expected:    true,
			description: "adding a single new database should return true",
		},
		{
			name:        "single item removed",
			oldItemList: []interface{}{"database1", "database2"},
			newItemList: []interface{}{"database1"},
			expected:    false,
			description: "removing a database should return false",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := databaselink.ForceLinkNeeded(tc.oldItemList, tc.newItemList)
			if result != tc.expected {
				t.Errorf("ForceLinkNeeded() = %v, expected %v. %s", result, tc.expected, tc.description)
			}
		})
	}
}
