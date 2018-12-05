package azure

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func SchemaLocation() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		StateFunc:        NormalizeLocation,
		DiffSuppressFunc: SuppressLocationDiff,
	}
}

func SchemaLocationForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
}

func SchemaLocationDeprecated() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		ForceNew:         true,
		Optional:         true,
		StateFunc:        NormalizeLocation,
		DiffSuppressFunc: SuppressLocationDiff,
		Deprecated:       "location is no longer used",
	}
}

func FlattenAndSetLocation(d *schema.ResourceData, location *string) error {
	if location != nil {
		if err := d.Set("location", NormalizeLocation(*location)); err != nil {
			return fmt.Errorf("Error setting `location`: %s", err)
		}
	}
	return nil
}

// azureRMNormalizeLocation is a function which normalises human-readable region/location
// names (e.g. "West US") to the values used and returned by the Azure API (e.g. "westus").
// In state we track the API internal version as it is easier to go from the human form
// to the canonical form than the other way around.
func NormalizeLocation(location interface{}) string {
	input := location.(string)
	return strings.Replace(strings.ToLower(input), " ", "", -1)
}

func SuppressLocationDiff(k, old, new string, _ *schema.ResourceData) bool {
	return NormalizeLocation(old) == NormalizeLocation(new)
}

func HashAzureLocation(location interface{}) int {
	return hashcode.String(NormalizeLocation(location.(string)))
}
