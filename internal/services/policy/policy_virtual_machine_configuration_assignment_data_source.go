package policy

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/guestconfiguration/mgmt/2020-06-25/guestconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourcePolicyVirtualMachineConfigurationAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePolicyVirtualMachineConfigurationAssignmentRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"virtual_machine_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"content_hash": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"content_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"assignment_hash": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"compliance_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"latest_report_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"last_compliance_status_checked": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePolicyVirtualMachineConfigurationAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.GuestConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	vmName := d.Get("virtual_machine_name").(string)
	name := d.Get("name").(string)

	id := parse.NewVirtualMachineConfigurationAssignmentID(subscriptionId, resourceGroup, vmName, name)

	resp, err := client.Get(ctx, id.ResourceGroup, id.GuestConfigurationAssignmentName, id.VirtualMachineName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] guestConfiguration %q was not found", id.GuestConfigurationAssignmentName)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.GuestConfigurationAssignmentName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("virtual_machine_name", vmName)

	if props := resp.Properties; props != nil {
		if v := props.AssignmentHash; v != nil {
			d.Set("assignment_hash", v)
		}

		if v := string(props.ComplianceStatus); v != "" {
			d.Set("compliance_status", v)
		}

		if v := props.LatestReportID; v != nil {
			d.Set("latest_report_id", v)
		}

		if v := props.LastComplianceStatusChecked; v != nil {
			d.Set("last_compliance_status_checked", v.Format(time.RFC3339))
		}

		contentHash, contentUri := dataSourceFlattenGuestConfigurationAssignment(props.GuestConfiguration)

		if contentHash != nil {
			d.Set("content_hash", contentHash)
		}

		if contentUri != nil {
			d.Set("content_uri", contentUri)
		}
	}
	return nil
}

func dataSourceFlattenGuestConfigurationAssignment(input *guestconfiguration.Navigation) (*string, *string) {
	if input == nil {
		return nil, nil
	}

	var contentHash *string
	if input.ContentHash != nil {
		contentHash = input.ContentHash
	}
	var contentUri *string
	if input.ContentURI != nil {
		contentUri = input.ContentURI
	}

	return contentHash, contentUri
}
