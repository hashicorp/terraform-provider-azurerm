package importexport

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/importexport/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func importAzureImportExportJob(jobType string, resourceType string) func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	return func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
		id, err := parse.StorageImportExportJobID(d.Id())
		if err != nil {
			return []*schema.ResourceData{}, err
		}

		client := meta.(*clients.Client).ImportExport.JobClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		job, err := client.Get(ctx, id.Name, id.ResourceGroup)
		if err != nil {
			return []*schema.ResourceData{}, fmt.Errorf("failure retrieving Azure Import Export Job %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if job.Properties == nil {
			return []*schema.ResourceData{}, fmt.Errorf("failure retrieving Azure Import Export Job %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
		}

		if job.Properties.JobType == nil {
			return []*schema.ResourceData{}, fmt.Errorf("failure retrieving Azure Import Export Job %q (Resource Group %q): `properties.JobType` was nil", id.Name, id.ResourceGroup)
		}
		if *job.Properties.JobType != jobType {
			return []*schema.ResourceData{}, fmt.Errorf("the %q resource only supports Azure %s Job", resourceType, jobType)
		}

		return []*schema.ResourceData{d}, nil
	}
}
