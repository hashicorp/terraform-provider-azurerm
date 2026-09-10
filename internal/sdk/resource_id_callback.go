package sdk

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func setIDCallback(client any, id resourceids.ResourceId, d *pluginsdk.ResourceData, supportsIdentity bool, identityType pluginsdk.ResourceTypeForIdentity) func() error {
	c, ok := client.(*clients.Client)
	if !ok {
		panic(fmt.Sprintf("internal-error: expected `*clients.Client` but got %T", client))
	}

	if c.Features.PersistIDOnCreateBeforePollingForCompletion {
		return func() error {
			d.SetId(id.ID())
			if supportsIdentity {
				return pluginsdk.SetResourceIdentityData(d, id, identityType)
			}
			return nil
		}
	}

	return nil
}

// Untyped resources

func SetIDCallback(client any, id resourceids.ResourceId, d *pluginsdk.ResourceData) func() error {
	return setIDCallback(client, id, d, false, pluginsdk.ResourceTypeForIdentityDefault)
}

func SetIDAndIdentityCallback(client any, id resourceids.ResourceId, d *pluginsdk.ResourceData) func() error {
	return setIDCallback(client, id, d, true, pluginsdk.ResourceTypeForIdentityDefault)
}

func SetIDAndIdentityWithTypeCallback(client any, id resourceids.ResourceId, d *pluginsdk.ResourceData, identityType pluginsdk.ResourceTypeForIdentity) func() error {
	return setIDCallback(client, id, d, true, identityType)
}

// Typed resources

func (rmd ResourceMetaData) SetIDCallback(id resourceids.ResourceId) func() error {
	return setIDCallback(rmd.Client, id, rmd.ResourceData, false, pluginsdk.ResourceTypeForIdentityDefault)
}

func (rmd ResourceMetaData) SetIDAndIdentityCallback(id resourceids.ResourceId) func() error {
	return setIDCallback(rmd.Client, id, rmd.ResourceData, true, pluginsdk.ResourceTypeForIdentityDefault)
}

func (rmd ResourceMetaData) SetIDAndIdentityWithTypeCallback(id resourceids.ResourceId, identityType pluginsdk.ResourceTypeForIdentity) func() error {
	return setIDCallback(rmd.Client, id, rmd.ResourceData, true, identityType)
}
