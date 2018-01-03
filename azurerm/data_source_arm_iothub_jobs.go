package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmIotHubJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmIotHubJobsRead,

		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameForDataSourceSchema(),
			"iot_hub_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceArmIotHubJobsRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient
	log.Printf("[INFO] Acquiring IotHub Job ID")

	iothubName := d.Get("iot_hub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	jobID := d.Get("job_id").(string)

	jobResp, err := iothubClient.GetJob(resourceGroup, iothubName, jobID)
	if err != nil {
		return fmt.Errorf("Error retrieving job with id: %s", jobID)
	}

	d.Set("job_id", jobResp.JobID)
	d.Set("status", jobResp.Status)
	d.Set("start", jobResp.StartTimeUtc)
	d.Set("end", jobResp.EndTimeUtc)
	d.Set("type", jobResp.Type)
	if jobResp.Status == "failed" {
		d.Set("status_fail", jobResp.FailureReason)
	}
	return nil
}
