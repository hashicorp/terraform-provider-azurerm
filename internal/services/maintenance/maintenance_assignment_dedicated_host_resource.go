package maintenance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/dedicatedhosts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2021-05-01/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2021-05-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/parse"
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
			_, err := parse.MaintenanceAssignmentDedicatedHostID(id)
			return err
		}),

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
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     dedicatedhosts.ValidateHostID,
				DiffSuppressFunc: suppress.CaseDifference,
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

	dedicatedHostIdRaw := d.Get("dedicated_host_id").(string)
	dedicatedHostId, _ := dedicatedhosts.ParseHostID(dedicatedHostIdRaw)

	existingList, err := getMaintenanceAssignmentDedicatedHost(ctx, client, *dedicatedHostId, dedicatedHostIdRaw)

	if err != nil {
		return err
	}
	if existingList != nil && len(*existingList) > 0 {
		existing := (*existingList)[0]
		if existing.Id != nil && *existing.Id != "" {
			return tf.ImportAsExistsError("azurerm_maintenance_assignment_dedicated_host", configurationId.ID())
		}
	}

	// set assignment name to configuration name
	assignmentName := configurationId.ResourceName
	configurationAssignment := configurationassignments.ConfigurationAssignment{
		Name:     utils.String(assignmentName),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &configurationassignments.ConfigurationAssignmentProperties{
			MaintenanceConfigurationId: utils.String(configurationId.ID()),
			ResourceId:                 utils.String(dedicatedHostId.ID()),
		},
	}

	// It may take a few minutes after starting a VM for it to become available to assign to a configuration

	id := configurationassignments.NewProviders2ConfigurationAssignmentID(dedicatedHostId.SubscriptionId, dedicatedHostId.ResourceGroupName, "Microsoft.Compute", "hostGroups", dedicatedHostId.HostGroupName, "hosts", dedicatedHostId.HostName, assignmentName)
	err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if _, err := client.CreateOrUpdateParent(ctx, id, configurationAssignment); err != nil {
			if strings.Contains(err.Error(), "It may take a few minutes after starting a VM for it to become available to assign to a configuration") {
				return pluginsdk.RetryableError(fmt.Errorf("expected VM is available to assign to a configuration but was in pending state, retrying"))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("issuing creating request for Maintenance Assignment (Dedicated Host ID %q): %+v", dedicatedHostId.ID(), err))
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

	id, err := parse.MaintenanceAssignmentDedicatedHostID(d.Id())
	if err != nil {
		return err
	}

	resp, err := getMaintenanceAssignmentDedicatedHost(ctx, client, id.DedicatedHostId, id.DedicatedHostIdRaw)
	if err != nil {
		return err
	}
	if resp == nil || len(*resp) == 0 {
		d.SetId("")
		return nil
	}
	assignment := (*resp)[0]
	if assignment.Id == nil || *assignment.Id == "" {
		return fmt.Errorf("empty or nil ID of Maintenance Assignment (Dedicated Host ID: %q", id.DedicatedHostIdRaw)
	}

	d.Set("dedicated_host_id", id.DedicatedHostId.ID())

	if props := assignment.Properties; props != nil {
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
	return nil
}

func resourceArmMaintenanceAssignmentDedicatedHostDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MaintenanceAssignmentDedicatedHostID(d.Id())
	if err != nil {
		return err
	}

	providerConfigAssignmentId := configurationassignments.NewProviders2ConfigurationAssignmentID(id.DedicatedHostId.SubscriptionId, id.DedicatedHostId.ResourceGroupName, "Microsoft.Compute", "hostGroups", id.DedicatedHostId.HostGroupName, "hosts", id.DedicatedHostId.HostName, id.Name)

	if _, err := client.DeleteParent(ctx, providerConfigAssignmentId); err != nil {
		return fmt.Errorf("deleting Maintenance Assignment to resource %q: %+v", id.DedicatedHostIdRaw, err)
	}

	return nil
}

func getMaintenanceAssignmentDedicatedHost(ctx context.Context, client *configurationassignments.ConfigurationAssignmentsClient, hostId dedicatedhosts.HostId, dedicatedHostId string) (result *[]configurationassignments.ConfigurationAssignment, err error) {
	id := configurationassignments.NewResourceGroupProviderID(hostId.SubscriptionId, hostId.ResourceGroupName, "Microsoft.Compute", "hostGroups", hostId.HostGroupName, "hosts", hostId.HostName)

	resp, err := client.ListParent(ctx, id)

	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			err = fmt.Errorf("checking for presence of existing Maintenance assignment (Dedicated Host ID %q): %+v", dedicatedHostId, err)
			return nil, err
		}
	}
	return resp.Model.Value, nil
}
