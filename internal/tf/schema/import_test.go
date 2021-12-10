package schema

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestValidateResourceIDPriorToImport(t *testing.T) {
	testData := []struct {
		name             string
		id               string
		shouldBeImported bool
		validator        ResourceIDValidator
	}{
		{
			name:             "returns an error",
			shouldBeImported: false,
			validator: func(input string) error {
				return fmt.Errorf("Returns an error")
			},
		},
		{
			name:             "valid",
			shouldBeImported: true,
			validator: func(input string) error {
				return nil
			},
		},
	}

	for _, v := range testData {
		log.Printf("[DEBUG] Testing %q", v.name)

		f := ValidateResourceIDPriorToImport(v.validator)
		resourceData := &schema.ResourceData{}
		resourceData.SetId("hello")
		_, err := f.StateContext(context.TODO(), resourceData, nil)
		wasImported := err == nil

		if v.shouldBeImported != wasImported {
			t.Fatalf("Expected %t but got %t. Errors: %+v", v.shouldBeImported, wasImported, err)
		}
	}
}
