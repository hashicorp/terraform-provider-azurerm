package pluginsdk

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IDValidationFunc func(id string) error

type ImporterFunc = func(ctx context.Context, d *ResourceData, meta interface{}) ([]*ResourceData, error)

// DefaultImporter is a wrapper around the default importer within the Plugin SDK
// at this point resources should be using ImporterValidatingResourceId, but this
// is providing a compatibility shim for the moment
func DefaultImporter() *schema.ResourceImporter {
	// NOTE: we should do a secondary sweep and move things _off_ of this, since all resources
	// should be validating the Resource ID at import time at this point forwards
	return &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
}

// ImporterValidatingResourceId validates the ID provided at import time is valid
// using the validateFunc.
func ImporterValidatingResourceId(validateFunc IDValidationFunc) *schema.ResourceImporter {
	thenFunc := func(ctx context.Context, d *ResourceData, meta interface{}) ([]*ResourceData, error) {
		return []*ResourceData{d}, nil
	}
	return ImporterValidatingResourceIdThen(validateFunc, thenFunc)
}

// ImporterValidatingResourceIdThen validates the ID provided at import time is valid
// using the validateFunc then runs the 'thenFunc', allowing the import to be customised.
func ImporterValidatingResourceIdThen(validateFunc IDValidationFunc, thenFunc ImporterFunc) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(ctx context.Context, d *ResourceData, meta interface{}) ([]*ResourceData, error) {
			log.Printf("[DEBUG] Importing Resource - parsing %q", d.Id())

			if err := validateFunc(d.Id()); err != nil {
				return []*ResourceData{d}, fmt.Errorf("parsing Resource ID %q: %+v", d.Id(), err)
			}

			return thenFunc(ctx, d, meta)
		},
	}
}
