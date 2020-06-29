package tags

import (
	"fmt"
	"testing"
)

func TestExpand(t *testing.T) {
	testData := make(map[string]interface{})
	testData["key1"] = "value1"
	testData["key2"] = 21
	testData["key3"] = "value3"

	expanded := Expand(testData)

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
