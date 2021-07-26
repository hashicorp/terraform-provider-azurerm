package azuread

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/graph"
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

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validate.NoEmptyStrings,
				ConflictsWith: []string{"object_id"},
			},

			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"owners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
		g, err := graph.GroupGetByDisplayName(&client, ctx, name)
		if err != nil {
			return fmt.Errorf("Error finding Azure AD Group with display name %q: %+v", name, err)
		}
		group = *g
	} else {
		return fmt.Errorf("one of `object_id` or `name` must be supplied")
	}

	if group.ObjectID == nil {
		return fmt.Errorf("Group objectId is nil")
	}
	d.SetId(*group.ObjectID)

	d.Set("object_id", group.ObjectID)
	d.Set("name", group.DisplayName)

	if v, ok := group.AdditionalProperties["description"]; ok {
		d.Set("description", v.(string))
	}

	members, err := graph.GroupAllMembers(client, ctx, d.Id())
	if err != nil {
		return err
	}
	d.Set("members", members)

	owners, err := graph.GroupAllOwners(client, ctx, d.Id())
	if err != nil {
		return err
	}
	d.Set("owners", owners)

	return nil
}
