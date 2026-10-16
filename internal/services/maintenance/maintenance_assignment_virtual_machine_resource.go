// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance

//go:generate go run ../../tools/generator-tests resourceidentity -properties "scope:virtual_machine_id" -compare-values "name:maintenance_configuration_id" -test-expect-non-empty -has-id-casing-bug

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmMaintenanceAssignmentVirtualMachine() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmMaintenanceAssignmentVirtualMachineCreate,
		Read:   resourceArmMaintenanceAssignmentVirtualMachineRead,
		Delete: resourceArmMaintenanceAssignmentVirtualMachineDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentityThen(&configurationassignments.ScopedConfigurationAssignmentId{}, func(ctx context.Context, d *pluginsdk.ResourceData, m any) ([]*pluginsdk.ResourceData, error) {
			id, err := configurationassignments.ParseScopedConfigurationAssignmentID(d.Id())
			if err != nil {
				return nil, err
			}
			if _, err := virtualmachines.ParseVirtualMachineID(id.Scope); err != nil {
				return nil, fmt.Errorf("parsing %q as Virtual Machine ID: %+v", id.Scope, err)
			}
			return []*pluginsdk.ResourceData{d}, nil
		}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&configurationassignments.ScopedConfigurationAssignmentId{}),
		},

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AssignmentVirtualMachineV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"location": commonschema.Location(),

			"maintenance_configuration_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     maintenanceconfigurations.ValidateMaintenanceConfigurationID,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"virtual_machine_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualmachines.ValidateVirtualMachineID,
			},
		},
	}
}

func resourceArmMaintenanceAssignmentVirtualMachineCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualMachineId, err := virtualmachines.ParseVirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}

	configurationId, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(d.Get("maintenance_configuration_id").(string))
	if err != nil {
		return err
	}

	id := configurationassignments.NewScopedConfigurationAssignmentID(virtualMachineId.ID(), configurationId.MaintenanceConfigurationName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_maintenance_assignment_virtual_machine", id.ID())
		}
	}

	// set assignment name to configuration name
	assignmentName := configurationId.MaintenanceConfigurationName
	configurationAssignment := configurationassignments.ConfigurationAssignment{
		Name:     new(assignmentName),
		Location: new(location.Normalize(d.Get("location").(string))),
		Properties: &configurationassignments.ConfigurationAssignmentProperties{
			MaintenanceConfigurationId: new(configurationId.ID()),
			ResourceId:                 new(virtualMachineId.ID()),
		},
	}

	// It may take a few minutes after starting a VM for it to become available to assign to a configuration
	if err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if _, err := client.CreateOrUpdate(ctx, id, configurationAssignment); err != nil {
			if strings.Contains(err.Error(), "It may take a few minutes after starting a VM for it to become available to assign to a configuration") {
				return pluginsdk.RetryableError(errors.New("expected VM is available to assign to a configuration but was in pending state, retrying"))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("issuing creating request for %s: %+v", id, err))
		}

		return nil
	}); err != nil {
		return err
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return fmt.Errorf("setting resource identity: %w", err)
	}
	return resourceArmMaintenanceAssignmentVirtualMachineRead(d, meta)
}

func resourceArmMaintenanceAssignmentVirtualMachineRead(d *pluginsdk.ResourceData, meta any) error {
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

	return resourceArmMaintenanceAssignmentVirtualMachineFlatten(d, id, resp.Model)
}

func resourceArmMaintenanceAssignmentVirtualMachineFlatten(d *pluginsdk.ResourceData, id *configurationassignments.ScopedConfigurationAssignmentId, model *configurationassignments.ConfigurationAssignment) error {
	vmId, err := virtualmachines.ParseVirtualMachineID(id.Scope)
	if err != nil {
		return err
	}

	d.Set("virtual_machine_id", vmId.ID())

	if model != nil {
		loc := location.NormalizeNilable(model.Location)
		// location isn't returned by the API
		if loc == "" {
			if existingLoc, ok := d.GetOk("location"); ok {
				loc = existingLoc.(string)
			}
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
	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceArmMaintenanceAssignmentVirtualMachineDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := configurationassignments.ParseScopedConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
