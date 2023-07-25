// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maintenance

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/dedicatedhosts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmMaintenanceAssignmentDedicatedHost() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmMaintenanceAssignmentDedicatedHostCreate,
		Read:   resourceArmMaintenanceAssignmentDedicatedHostRead,
		Delete: resourceArmMaintenanceAssignmentDedicatedHostDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			parsed, err := configurationassignments.ParseScopedConfigurationAssignmentID(id)
			if err != nil {
				return err
			}

			if _, err := dedicatedhosts.ParseHostID(parsed.Scope); err != nil {
				return fmt.Errorf("parsing %q as a Dedicated Host ID: %+v", parsed.Scope, err)
			}

			return nil
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AssignmentDedicatedHostV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"location": commonschema.Location(),

			"maintenance_configuration_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     maintenanceconfigurations.ValidateMaintenanceConfigurationID,
				DiffSuppressFunc: suppress.CaseDifference, // TODO remove in 4.0 with a work around or when https://github.com/Azure/azure-rest-api-specs/issues/8653 is fixed
			},

			"dedicated_host_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: dedicatedhosts.ValidateHostID,
			},
		},
	}
}

func resourceArmMaintenanceAssignmentDedicatedHostCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	configurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(d.Get("maintenance_configuration_id").(string))
	if err != nil {
		return err
	}

	dedicatedHostId, err := dedicatedhosts.ParseHostID(d.Get("dedicated_host_id").(string))
	if err != nil {
		return err
	}

	id := configurationassignments.NewScopedConfigurationAssignmentID(dedicatedHostId.ID(), configurationId.MaintenanceConfigurationName)
	resp, err := client.GetParent(ctx, id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(resp.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_maintenance_assignment_dedicated_host", id.ID())
	}

	// set assignment name to configuration name
	configurationAssignment := configurationassignments.ConfigurationAssignment{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &configurationassignments.ConfigurationAssignmentProperties{
			MaintenanceConfigurationId: utils.String(configurationId.ID()),
			ResourceId:                 utils.String(dedicatedHostId.ID()),
		},
	}

	// TODO: refactor to using a context-aware poller
	// It may take a few minutes after starting a VM for it to become available to assign to a configuration
	err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if _, err := client.CreateOrUpdateParent(ctx, id, configurationAssignment); err != nil {
			if strings.Contains(err.Error(), "It may take a few minutes after starting a VM for it to become available to assign to a configuration") {
				return pluginsdk.RetryableError(fmt.Errorf("expected VM is available to assign to a configuration but was in pending state, retrying"))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("issuing creating request for %s: %+v", id, err))
		}

		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceArmMaintenanceAssignmentDedicatedHostRead(d, meta)
}

func resourceArmMaintenanceAssignmentDedicatedHostRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurationassignments.ParseScopedConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("checking for presence of existing %s: %+v", *id, err)
	}

	dedicatedHostId, err := dedicatedhosts.ParseHostID(id.Scope)
	if err != nil {
		return fmt.Errorf("parsing %q as a dedicated host id: %+v", id.Scope, err)
	}
	d.Set("dedicated_host_id", dedicatedHostId.ID())

	if model := resp.Model; model != nil {
		loc := location.NormalizeNilable(model.Location)
		// location isn't returned by the API
		if loc == "" {
			loc = d.Get("location").(string)
		}
		d.Set("location", loc)

		if props := model.Properties; props != nil {
			maintenanceConfigurationId := ""
			if props.MaintenanceConfigurationId != nil {
				parsedId, err := maintenanceconfigurations.ParseMaintenanceConfigurationIDInsensitively(*props.MaintenanceConfigurationId)
				if err != nil {
					return fmt.Errorf("parsing %q: %+v", *props.MaintenanceConfigurationId, err)
				}
				maintenanceConfigurationId = parsedId.ID()
			}
			d.Set("maintenance_configuration_id", maintenanceConfigurationId)
		}
	}

	return nil
}

func resourceArmMaintenanceAssignmentDedicatedHostDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurationassignments.ParseScopedConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.DeleteParent(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
