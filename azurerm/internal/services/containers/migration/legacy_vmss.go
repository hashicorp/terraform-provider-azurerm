package migration

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func LegacyVMSSV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Upgrade: legacyVMSSUpgradeV0ToV1,
	}
}

func legacyVMSSUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	// @tombuildsstuff: NOTE, this state migration is essentially pointless
	// however it existed in the legacy migration so even though this is
	//  essentially a noop there's no reason this shouldn't be the same I guess

	client := meta.(*clients.Client).Compute.VMScaleSetClient
	ctx, cancel := context.WithTimeout(meta.(*clients.Client).StopContext, 5*time.Minute)
	defer cancel()

	resGroup := rawState["resource_group_name"].(string)
	name := rawState["name"].(string)

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return rawState, err
	}

	rawState["id"] = *read.ID
	return rawState, nil
}
