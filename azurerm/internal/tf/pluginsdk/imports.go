package pluginsdk

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceImporter = schema.ResourceImporter

type IDValidationFunc func(id string) error
type ImporterFunc = func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error)

func DefaultImporter() *ResourceImporter {
	return &ResourceImporter{
		State: ImportStatePassthrough,
	}
}

// ImporterValidatingResourceId validates the ID provided at import time is valid
// using the validateFunc.
func ImporterValidatingResourceId(validateFunc IDValidationFunc) *ResourceImporter {
	var thenFunc = func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		return []*ResourceData{d}, nil
	}
	return ImporterValidatingResourceIdThen(validateFunc, thenFunc)
}

// ImporterValidatingResourceIdThen validates the ID provided at import time is valid
// using the validateFunc then runs the 'thenFunc', allowing the import to be customised.
func ImporterValidatingResourceIdThen(validateFunc IDValidationFunc, thenFunc ImporterFunc) *ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			log.Printf("[DEBUG] Importing Resource - parsing %q", d.Id())

			if err := validateFunc(d.Id()); err != nil {
				return []*schema.ResourceData{d}, fmt.Errorf("parsing Resource ID %q: %+v", d.Id(), err)
			}

			// TODO: replace with a real one with Plugin SDKv2
			ctx := context.TODO()
			return thenFunc(ctx, d, meta)
		},
	}
}
