package main

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const (
	RESOURCE_NAME    = "azurerm_foobar"
	BRAND_NAME       = "Foobar"
	RESOURCE_ID      = "12345"
	WEBSITE_CATEGORY = "Foobar Category"
)

func TestResourceArgumentBlock(t *testing.T) {
	expectedOut := strings.ReplaceAll(`## Arguments Reference

The following arguments are supported:

* 'block2' - (Required) One or more 'block2' blocks as defined below.

* 'foo_enabled' - (Required) Should the TODO be enabled?

* 'foo_id' - (Required) The ID of the TODO.

* 'list' - (Required) Specifies a list of TODO.

* 'location' - (Required) The Azure Region where the Foobar should exist. Changing this forces a new Foobar to be created.

* 'map' - (Required) Specifies a list of TODO.

* 'name' - (Required) The name which should be used for this Foobar. Changing this forces a new Foobar to be created.

* 'resource_group_name' - (Required) The name of the Resource Group where the Foobar should exist. Changing this forces a new Foobar to be created.

* 'set' - (Required) Specifies a list of TODO.

---

* 'tags' - (Optional) A mapping of tags which should be assigned to the Foobar.

---

A 'block1' block supports the following:

* 'nest_attr1' - (Optional) TODO.

---

A 'block2' block supports the following:

* 'block1' - (Required) A 'block1' block as defined above.

* 'block3' - (Required) One or more 'block3' blocks as defined below.

* 'nest_attr2' - (Optional) TODO.

---

A 'block3' block supports the following:

* 'nest_attr3' - (Optional) TODO.`, "'", "`")

	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"foo_enabled": {
				Type:     schema.TypeString,
				Required: true,
			},
			"foo_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"block2": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nest_attr2": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// lintignore:XS003
						"block1": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nest_attr1": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						// lintignore:XS003
						"block3": {
							Type:     schema.TypeList,
							MinItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nest_attr3": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
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
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	gen := setupDocGen(false, resource)
	actualOut := gen.argumentsBlock()

	runTest(t, expectedOut, actualOut)
}

func runTest(t *testing.T, expected, actual string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(actual, expected, true)
	hasDiff := false
	for _, diff := range diffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			hasDiff = true
			break
		}
	}
	if hasDiff {
		t.Fatal(dmp.DiffPrettyText(diffs))
	}
}

func setupDocGen(isDataSource bool, resource *schema.Resource) documentationGenerator {
	toStrPtr := func(input string) *string {
		return &input
	}
	return documentationGenerator{
		resourceName:      RESOURCE_NAME,
		brandName:         BRAND_NAME,
		resourceId:        toStrPtr(RESOURCE_ID),
		isDataSource:      isDataSource,
		websiteCategories: []string{WEBSITE_CATEGORY},
		resource:          resource,
	}
}
