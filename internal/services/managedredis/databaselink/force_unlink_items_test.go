package databaselink_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
)

func TestForceUnlinkItems(t *testing.T) {
	testCases := []struct {
		name          string
		oldItemList   []interface{}
		newItemList   []interface{}
		expectedNeed  bool
		expectedItems []string
		description   string
	}{
		{
			name:          "empty lists",
			oldItemList:   []interface{}{},
			newItemList:   []interface{}{},
			expectedNeed:  false,
			expectedItems: nil,
			description:   "both lists empty should return false with no items",
		},
		{
			name:          "empty old list with new items",
			oldItemList:   []interface{}{},
			newItemList:   []interface{}{"item1", "item2"},
			expectedNeed:  false,
			expectedItems: nil,
			description:   "no old items to unlink should return false",
		},
		{
			name:          "empty new list with old items",
			oldItemList:   []interface{}{"item1", "item2"},
			newItemList:   []interface{}{},
			expectedNeed:  true,
			expectedItems: []string{"item1", "item2"},
			description:   "all old items should be unlinked when new list is empty",
		},
		{
			name:          "identical lists",
			oldItemList:   []interface{}{"item1", "item2", "item3"},
			newItemList:   []interface{}{"item1", "item2", "item3"},
			expectedNeed:  false,
			expectedItems: nil,
			description:   "identical lists should return false with no items to unlink",
		},
		{
			name:          "identical lists different order",
			oldItemList:   []interface{}{"item1", "item2", "item3"},
			newItemList:   []interface{}{"item3", "item1", "item2"},
			expectedNeed:  false,
			expectedItems: nil,
			description:   "same items in different order should return false",
		},
		{
			name:          "new list has additional items",
			oldItemList:   []interface{}{"item1", "item2"},
			newItemList:   []interface{}{"item1", "item2", "item3"},
			expectedNeed:  false,
			expectedItems: nil,
			description:   "adding new items should not require unlinking",
		},
		{
			name:          "new list has some different items",
			oldItemList:   []interface{}{"item1", "item2"},
			newItemList:   []interface{}{"item1", "item3"},
			expectedNeed:  true,
			expectedItems: []string{"item2"},
			description:   "item2 should be unlinked since it's not in new list",
		},
		{
			name:          "new list is subset of old list",
			oldItemList:   []interface{}{"item1", "item2", "item3"},
			newItemList:   []interface{}{"item1", "item2"},
			expectedNeed:  true,
			expectedItems: []string{"item3"},
			description:   "item3 should be unlinked",
		},
		{
			name:          "completely different lists",
			oldItemList:   []interface{}{"item1", "item2"},
			newItemList:   []interface{}{"item3", "item4"},
			expectedNeed:  true,
			expectedItems: []string{"item1", "item2"},
			description:   "all old items should be unlinked",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			needUnlink, itemsToUnlink := databaselink.ForceUnlinkItems(tc.oldItemList, tc.newItemList)

			if needUnlink != tc.expectedNeed {
				t.Errorf("ForceUnlinkItems() needUnlink = %v, expected %v. %s", needUnlink, tc.expectedNeed, tc.description)
			}

			if !reflect.DeepEqual(itemsToUnlink, tc.expectedItems) {
				t.Errorf("ForceUnlinkItems() itemsToUnlink = %v, expected %v. %s", itemsToUnlink, tc.expectedItems, tc.description)
			}
		})
	}
}
