package suppress

import (
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func JsonDiff(_, old, new string, _ *schema.ResourceData) bool {
	if old == new {
		return true
	}
	var oldJson interface{}
	var newJson interface{}

	err := json.Unmarshal([]byte(old), &oldJson)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(new), &newJson)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(oldJson, newJson)
}
