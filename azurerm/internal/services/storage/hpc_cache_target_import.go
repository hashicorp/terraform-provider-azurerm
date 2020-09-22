package storage

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2020-03-01/storagecache"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func importHpcCache(kind storagecache.TargetType) func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	return func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
		id, err := parsers.HPCCacheTargetID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).Storage.StorageTargetsClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := client.Get(ctx, id.ResourceGroup, id.Cache, id.Name)
		if err != nil {
			return nil, fmt.Errorf("retrieving HPC Cache Target %q (Resource Group %q, Cahe %q): %+v", id.Name, id.ResourceGroup, id.Cache, err)
		}

		if resp.Type == nil {
			return nil, fmt.Errorf(`HPC Cache Target %q (Resource Group %q, Cahe %q) nil "type"`, id.Name, id.ResourceGroup, id.Cache)
		}

		if *resp.Type != string(kind) {
			return nil, fmt.Errorf(`HPC Cache Target %q (Resource Group %q, Cahe %q) "type" mismatch, expected "%s", got "%s"`, id.Name, id.ResourceGroup, id.Cache, kind, *resp.Type)
		}
		return []*schema.ResourceData{d}, nil
	}
}
