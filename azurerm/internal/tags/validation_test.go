package tags

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateMaximumNumberOfTags(t *testing.T) {
	tagsMap := make(map[string]interface{})
	for i := 0; i < 51; i++ {
		tagsMap[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}

	_, es := Validate(tagsMap, "tags")

	if len(es) != 1 {
		t.Fatal("Expected one validation error for too many tags")
	}

	if !strings.Contains(es[0].Error(), "a maximum of 50 tags") {
		t.Fatal("Wrong validation error message for too many tags")
	}
}

func TestValidateTagMaxKeyLength(t *testing.T) {
	tooLongKey := strings.Repeat("long", 128) + "a"
	tagsMap := make(map[string]interface{})
	tagsMap[tooLongKey] = "value"

	_, es := Validate(tagsMap, "tags")
	if len(es) != 1 {
		t.Fatal("Expected one validation error for a key which is > 512 chars")
	}

	if !strings.Contains(es[0].Error(), "maximum length for a tag key") {
		t.Fatal("Wrong validation error message maximum tag key length")
	}

	if !strings.Contains(es[0].Error(), tooLongKey) {
		t.Fatal("Expected validated error to contain the key name")
	}

	if !strings.Contains(es[0].Error(), "513") {
		t.Fatal("Expected the length in the validation error for tag key")
	}
}

func TestValidateTagMaxValueLength(t *testing.T) {
	tagsMap := make(map[string]interface{})
	tagsMap["toolong"] = strings.Repeat("long", 64) + "a"

	_, es := Validate(tagsMap, "tags")
	if len(es) != 1 {
		t.Fatal("Expected one validation error for a value which is > 256 chars")
	}

	if !strings.Contains(es[0].Error(), "maximum length for a tag value") {
		t.Fatal("Wrong validation error message for maximum tag value length")
	}

	if !strings.Contains(es[0].Error(), "toolong") {
		t.Fatal("Expected validated error to contain the key name")
	}

	if !strings.Contains(es[0].Error(), "257") {
		t.Fatal("Expected the length in the validation error for value")
	}
}
