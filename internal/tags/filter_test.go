package tags

import (
	"testing"
)

func TestFilter(t *testing.T) {
	testData := make(map[string]*string)
	valueData := [3]string{"value1", "value2", "value3"}

	testData["key1"] = &valueData[0]
	testData["key2"] = &valueData[1]
	testData["key3"] = &valueData[2]

	filtered := Filter(testData, "key1", "key3", "")

	if len(filtered) != 1 {
		t.Fatalf("Expected 1 result in filtered tag map, got %d", len(filtered))
	}

	if filtered["key2"] != &valueData[1] {
		t.Fatalf("Expected %v in filtered tag map, got %v", valueData[1], *filtered["key2"])
	}
}
