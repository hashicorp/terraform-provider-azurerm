package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

func TestExpandARMTags(t *testing.T) {
	testData := make(map[string]interface{})
	testData["key1"] = "value1"
	testData["key2"] = 21
	testData["key3"] = "value3"

	expanded := tags.Expand(testData)

	if len(expanded) != 3 {
		t.Fatalf("Expected 3 results in expanded tag map, got %d", len(expanded))
	}

	for k, v := range testData {
		var strVal string
		switch v := v.(type) {
		case string:
			strVal = v
		case int:
			strVal = fmt.Sprintf("%d", v)
		}

		if *expanded[k] != strVal {
			t.Fatalf("Expanded value %q incorrect: expected %q, got %q", k, strVal, *expanded[k])
		}
	}
}

func TestFilterARMTags(t *testing.T) {
	testData := make(map[string]*string)
	valueData := [3]string{"value1", "value2", "value3"}

	testData["key1"] = &valueData[0]
	testData["key2"] = &valueData[1]
	testData["key3"] = &valueData[2]

	filtered := tags.Filter(testData, "key1", "key3", "")

	if len(filtered) != 1 {
		t.Fatalf("Expected 1 result in filtered tag map, got %d", len(filtered))
	}

	if filtered["key2"] != &valueData[1] {
		t.Fatalf("Expected %v in filtered tag map, got %v", valueData[1], *filtered["key2"])
	}
}
