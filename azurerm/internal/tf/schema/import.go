package schema

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceIDValidator takes a Resource ID and confirms that it's Valid
type ResourceIDValidator func(resourceId string) error

// ValidateResourceIDPriorToImport parses the Resource ID to confirm it's
// valid for this Resource prior to performing an import - allowing for incorrect
// Resource ID's to be caught prior to Import and subsequent crashes
func ValidateResourceIDPriorToImport(idParser ResourceIDValidator) *schema.ResourceImporter {
	return ValidateResourceIDPriorToImportThen(idParser, schema.ImportStatePassthroughContext)
}

// ValidateResourceIDPriorToImportThen parses the Resource ID to confirm it's
// valid for this Resource prior to calling the importer - allowing for incorrect
// Resource ID's to be caught prior to Import and subsequent crashes
func ValidateResourceIDPriorToImportThen(idParser ResourceIDValidator, importer schema.StateContextFunc) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			log.Printf("[DEBUG] Importing Resource - parsing %q", d.Id())

			if err := idParser(d.Id()); err != nil {
				return []*schema.ResourceData{d}, fmt.Errorf("Error parsing Resource ID %q: %+v", d.Id(), err)
			}

			return importer(ctx, d, meta)
		},
	}
}
