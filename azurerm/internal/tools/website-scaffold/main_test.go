package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

const (
	RESOURCE_NAME    = "azurerm_foobar"
	BRAND_NAME       = "Foobar"
	RESOURCE_ID      = "12345"
	WEBSITE_CATEGORY = "Foobar Category"
)

func setupDocGen(isDataSource bool, resource *schema.Resource) documentationGenerator {
	return documentationGenerator{
		resourceName: RESOURCE_NAME,
		brandName:    BRAND_NAME,
		resourceId: func(s string) *string {
			return &s
		}(RESOURCE_ID),
		isDataSource:      isDataSource,
		websiteCategories: []string{WEBSITE_CATEGORY},
		resource:          resource,
	}
}

func TestResourceArgumentBlock(t *testing.T) {
	expectedOut := fmt.Sprintf(`## Arguments Reference

The following arguments are supported:

* %[1]sblock%[1]s - (Required) A %[1]sblock%[1]s block as defined below.

* %[1]sfoo_enabled%[1]s - (Required) Should the TODO be enabled?

* %[1]sfoo_id%[1]s - (Required) The ID of the TODO.

* %[1]slist%[1]s - (Required) Specifies an array of TODO.

* %[1]slocation%[1]s - (Required) The Azure Region where the Foobar should exist. Changing this forces a new Foobar to be created.

* %[1]smap%[1]s - (Required) Specifies a map of TODO.

* %[1]sname%[1]s - (Required) The Name which should be used for this Foobar. Changing this forces a new Foobar to be created.

* %[1]sresource_group_name%[1]s - (Required) The name of the Resource Group where the Foobar should exist. Changing this forces a new Foobar to be created.

* %[1]sset%[1]s - (Required) Specifies an array of TODO.

---

* %[1]stags%[1]s - (Optional) A mapping of tags which should be assigned to the Foobar.

---

A %[1]sblock%[1]s block supports the following:`, "`")

	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),
			"location":            azure.SchemaLocation(),
			"foo_enabled": {
				Type:     schema.TypeString,
				Required: true,
			},
			"foo_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"block": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
				},
			},
			"list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"set": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"map": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tags.Schema(),
		},
	}
	gen := setupDocGen(false, resource)

	actualOut := gen.argumentsBlock()

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(actualOut, expectedOut, true)
	hasDiff := false
	for _, diff := range diffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			hasDiff = true
			break
		}
	}
	if hasDiff {
		t.Fatal(dmp.DiffText1(diffs))
	}
}
