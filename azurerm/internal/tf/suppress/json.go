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

	json.Unmarshal([]byte(old), &oldJson)
	json.Unmarshal([]byte(new), &newJson)

	return reflect.DeepEqual(oldJson, newJson)
}
