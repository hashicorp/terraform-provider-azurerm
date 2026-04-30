package sdk

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func setIDCallback(client any, id resourceids.ResourceId, d *pluginsdk.ResourceData, supportsIdentity bool) func() error {
	c, ok := client.(*clients.Client)
	if !ok {
		panic(fmt.Sprintf("internal-error: expected `*clients.Client` but got %T", client))
	}

	if c.Features.TrackPollingFailuresInState {
		return func() error {
			d.SetId(id.ID())
			if supportsIdentity {
				return pluginsdk.SetResourceIdentityData(d, id)
			}
			return nil
		}
	}

	return nil
}

// Untyped resources

func SetIDCallback(client any, id resourceids.ResourceId, d *pluginsdk.ResourceData) func() error {
	return setIDCallback(client, id, d, false)
}

func SetIDAndIdentityCallback(client any, id resourceids.ResourceId, d *pluginsdk.ResourceData) func() error {
	return setIDCallback(client, id, d, true)
}

// Typed resources

func (rmd ResourceMetaData) SetIDCallBack(id resourceids.ResourceId) func() error {
	return setIDCallback(rmd.Client, id, rmd.ResourceData, false)
}

func (rmd ResourceMetaData) SetIDAndIdentityCallback(id resourceids.ResourceId) func() error {
	return setIDCallback(rmd.Client, id, rmd.ResourceData, true)
}
