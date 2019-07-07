package azuread

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

func dataGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceActiveDirectoryGroupRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validate.UUID,
				ConflictsWith: []string{"name"},
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validate.NoEmptyStrings,
				ConflictsWith: []string{"object_id"},
			},
		},
	}
}

func dataSourceActiveDirectoryGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	var group graphrbac.ADGroup

	if oId, ok := d.Get("object_id").(string); ok && oId != "" {
		// use the object_id to find the Azure AD application
		resp, err := client.Get(ctx, oId)
		if err != nil {
			if ar.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error: AzureAD Group with ID %q was not found", oId)
			}

			return fmt.Errorf("Error making Read request on AzureAD Group with ID %q: %+v", oId, err)
		}

		group = resp
	} else if name, ok := d.Get("name").(string); ok && name != "" {
		filter := fmt.Sprintf("displayName eq '%s'", name)

		resp, err := client.ListComplete(ctx, filter)
		if err != nil {
			return fmt.Errorf("Error listing Azure AD Groups for filter %q: %+v", filter, err)
		}

		values := resp.Response().Value
		if values == nil {
			return fmt.Errorf("nil values for AD Groups matching %q", filter)
		}
		if len(*values) == 0 {
			return fmt.Errorf("Found no AD Groups matching %q", filter)
		}
		if len(*values) > 2 {
			return fmt.Errorf("Found multiple AD Groups matching %q", filter)
		}

		group = (*values)[0]
		if group.DisplayName == nil {
			return fmt.Errorf("nil DisplayName for AD Groups matching %q", filter)
		}
		if *group.DisplayName != name {
			return fmt.Errorf("displayname for AD Groups matching %q does is does not match(%q!=%q)", filter, *group.DisplayName, name)
		}
	} else {
		return fmt.Errorf("one of `object_id` or `name` must be supplied")
	}

	if group.ObjectID == nil {
		return fmt.Errorf("Group objectId is nil")
	}
	d.SetId(*group.ObjectID)

	d.Set("object_id", group.ObjectID)
	d.Set("name", group.DisplayName)
	return nil
}
