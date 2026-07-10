package testdata

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/networkprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkProfileCreateUpdate,
		Read:   resourceNetworkProfileRead,
		Update: resourceNetworkProfileCreateUpdate,
		Delete: resourceNetworkProfileDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := networkprofiles.ParseNetworkProfileID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceNetworkProfileCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourceNetworkProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkProfiles
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networkprofiles.ParseNetworkProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, networkprofiles.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NetworkProfileName)
	d.Set("resource_group_name", id.ResourceGroupName)

	return nil
}

func resourceNetworkProfileDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}
