package azuread

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/graph"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

func dataGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupsRead,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"object_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MinItems:      1,
				ConflictsWith: []string{"names"},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.UUID,
				},
			},

			"names": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				MinItems:      1,
				ConflictsWith: []string{"object_ids"},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
			},
		},
	}
}

func dataSourceGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	var groups []graphrbac.ADGroup
	expectedCount := 0

	if names, ok := d.Get("names").([]interface{}); ok && len(names) > 0 {
		expectedCount = len(names)
		for _, v := range names {
			g, err := graph.GroupGetByDisplayName(&client, ctx, v.(string))
			if err != nil {
				return fmt.Errorf("Error finding Azure AD Group with display name %q: %+v", v.(string), err)
			}
			groups = append(groups, *g)
		}
	} else if oids, ok := d.Get("object_ids").([]interface{}); ok && len(oids) > 0 {
		expectedCount = len(oids)
		for _, v := range oids {
			resp, err := client.Get(ctx, v.(string))
			if err != nil {
				return fmt.Errorf("Error making Read request on AzureAD Group with ID %q: %+v", v.(string), err)
			}

			groups = append(groups, resp)
		}
	} else {
		return fmt.Errorf("one of `object_ids` or `names` must be supplied")
	}

	if len(groups) != expectedCount {
		return fmt.Errorf("Unexpected number of groups returned (%d != %d)", len(groups), expectedCount)
	}

	names := make([]string, 0, len(groups))
	oids := make([]string, 0, len(groups))
	for _, u := range groups {
		if u.ObjectID == nil || u.DisplayName == nil {
			return fmt.Errorf("User with nil ObjectId or UPN was found: %v", u)
		}

		oids = append(oids, *u.ObjectID)
		names = append(names, *u.DisplayName)
	}

	h := sha1.New()
	if _, err := h.Write([]byte(strings.Join(names, "-"))); err != nil {
		return fmt.Errorf("Unable to compute hash for names: %v", err)
	}

	d.SetId("groups#" + base64.URLEncoding.EncodeToString(h.Sum(nil)))
	d.Set("object_ids", oids)
	d.Set("names", names)
	return nil
}
