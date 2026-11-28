// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_virtual_machine_configuration_assignment

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2024-04-05/guestconfigurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Guest Configuration Policy.",
		Read:        dataSourcePolicyVirtualMachineConfigurationAssignmentRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of the Guest Configuration Assignment.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"virtual_machine_name": {
				Description: "Only retrieve Policy Set Definitions from this Management Group.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"content_hash": {
				Description: "The content hash for the Guest Configuration package.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"content_uri": {
				Description: "The content URI where the Guest Configuration package is stored.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"assignment_hash": {
				Description: "Combined hash of the configuration package and parameters.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"compliance_status": {
				Description: "A value indicating compliance status of the machine for the assigned guest configuration.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"latest_report_id": {
				Description: "The ID of the latest report for the guest configuration assignment.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"last_compliance_status_checked": {
				Description: "Date and time, in RFC3339 format, when the machine's compliance status was last checked.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcePolicyVirtualMachineConfigurationAssignmentRead(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.GuestConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	vmName := d.Get("virtual_machine_name").(string)
	name := d.Get("name").(string)

	id := guestconfigurationassignments.NewVirtualMachineProviders2GuestConfigurationAssignmentID(subscriptionId, resourceGroup, vmName, name)
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] guestConfiguration %q was not found", id.GuestConfigurationAssignmentName)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.GuestConfigurationAssignmentName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("virtual_machine_name", vmName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("assignment_hash", pointer.From(props.AssignmentHash))
			d.Set("compliance_status", pointer.FromEnum(props.ComplianceStatus))
			d.Set("latest_report_id", pointer.From(props.LatestReportId))
			d.Set("last_compliance_status_checked", pointer.From(props.LastComplianceStatusChecked))

			contentHash, contentUri := dataSourceFlattenGuestConfigurationAssignment(props.GuestConfiguration)

			if contentHash != nil {
				d.Set("content_hash", contentHash)
			}

			if contentUri != nil {
				d.Set("content_uri", contentUri)
			}
		}
	}
	return nil
}

func dataSourceFlattenGuestConfigurationAssignment(input *guestconfigurationassignments.GuestConfigurationNavigation) (*string, *string) {
	if input == nil {
		return nil, nil
	}

	var contentHash *string
	if input.ContentHash != nil {
		contentHash = input.ContentHash
	}
	var contentUri *string
	if input.ContentUri != nil {
		contentUri = input.ContentUri
	}

	return contentHash, contentUri
}
