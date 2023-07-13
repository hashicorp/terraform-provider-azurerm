// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2020-06-25/guestconfigurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

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

	id := guestconfigurationassignments.NewProviders2GuestConfigurationAssignmentID(subscriptionId, resourceGroup, vmName, name)
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
			d.Set("compliance_status", string(pointer.From(props.ComplianceStatus)))
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
