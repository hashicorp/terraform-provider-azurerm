// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimedisableinteractivequery"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimeenableinteractivequery"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryIntegrationRuntimeAzure() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceDataFactoryIntegrationRuntimeAzureCreate,
		Read:   resourceDataFactoryIntegrationRuntimeAzureRead,
		Update: resourceDataFactoryIntegrationRuntimeAzureUpdate,
		Delete: resourceDataFactoryIntegrationRuntimeAzureDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DataFactoryIntegrationRuntimeAzureV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := integrationruntimes.ParseIntegrationRuntimeID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^([a-zA-Z0-9](-|-?[a-zA-Z0-9]+)+[a-zA-Z0-9])$`),
					`Invalid name for Azure Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
				),
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: factories.ValidateFactoryID,
			},

			"location": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					location.EnhancedValidate,
					validation.StringInSlice([]string{"AutoResolve"}, false),
				),
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},

			"cleanup_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"compute_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(integrationruntimes.DataFlowComputeTypeGeneral),
				ValidateFunc: validation.StringInSlice([]string{
					string(integrationruntimes.DataFlowComputeTypeGeneral),
					string(integrationruntimes.DataFlowComputeTypeComputeOptimized),
					string(integrationruntimes.DataFlowComputeTypeMemoryOptimized),
				}, false),
			},

			"core_count": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  8,
				ValidateFunc: validation.IntInSlice([]int{
					8, 16, 32, 48, 80, 144, 272,
				}),
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"interactive_authoring_time_to_live_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{10, 30, 60, 120}),
			},

			"time_to_live_min": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  0,
			},

			"virtual_network_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}

	return resource
}

func resourceDataFactoryIntegrationRuntimeAzureCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	enableInteractiveQueryClient := meta.(*clients.Client).DataFactory.IntegrationRuntimeEnableInteractiveQueryClient
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworks

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := integrationruntimes.NewIntegrationRuntimeID(dataFactoryId.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	existing, err := client.Get(ctx, id, integrationruntimes.DefaultGetOperationOptions())
	if err != nil && !response.WasNotFound(existing.HttpResponse) {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_data_factory_integration_runtime_azure", id.ID())
	}

	managedIntegrationRuntime := integrationruntimes.ManagedIntegrationRuntime{
		Description: pointer.To(d.Get("description").(string)),
		Type:        integrationruntimes.IntegrationRuntimeTypeManaged,
		TypeProperties: integrationruntimes.ManagedIntegrationRuntimeTypeProperties{
			ComputeProperties: &integrationruntimes.IntegrationRuntimeComputeProperties{
				Location: pointer.To(location.Normalize(d.Get("location").(string))),
				DataFlowProperties: &integrationruntimes.IntegrationRuntimeDataFlowProperties{
					ComputeType: pointer.ToEnum[integrationruntimes.DataFlowComputeType](d.Get("compute_type").(string)),
					CoreCount:   pointer.To(int64(d.Get("core_count").(int))),
					TimeToLive:  pointer.To(int64(d.Get("time_to_live_min").(int))),
					Cleanup:     pointer.To(d.Get("cleanup_enabled").(bool)),
				},
			},
		},
	}

	if d.Get("virtual_network_enabled").(bool) {
		virtualNetworkName, err := getManagedVirtualNetworkName(ctx, managedVirtualNetworksClient, id.SubscriptionId, id.ResourceGroupName, id.FactoryName)
		if err != nil {
			return err
		}
		if virtualNetworkName == nil {
			return fmt.Errorf("virtual network feature for azure integration runtime is only available after managed virtual network for this data factory is enabled")
		}
		managedIntegrationRuntime.ManagedVirtualNetwork = &integrationruntimes.ManagedVirtualNetworkReference{
			Type:          integrationruntimes.ManagedVirtualNetworkReferenceTypeManagedVirtualNetworkReference,
			ReferenceName: *virtualNetworkName,
		}
	}

	integrationRuntime := integrationruntimes.IntegrationRuntimeResource{
		Name:       &id.IntegrationRuntimeName,
		Properties: managedIntegrationRuntime,
	}

	if _, err := client.CreateOrUpdate(ctx, id, integrationRuntime, integrationruntimes.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if ttl := d.Get("interactive_authoring_time_to_live_in_minutes").(int); ttl != 0 {
		// Interactive Authoring/Query can only be modified once the integration runtime is online
		poller := pollers.NewPoller(custompollers.NewDataFactoryIntegrationRuntimeStatusPoller(client, id), time.Second*5, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		err := poller.PollUntilDone(ctx)
		if err != nil {
			return fmt.Errorf("waiting for state change of %s: %+v", id, err)
		}

		payload := integrationruntimeenableinteractivequery.EnableInteractiveQueryRequest{
			AutoTerminationMinutes: pointer.To(int64(ttl)),
		}

		if err := enableInteractiveQueryClient.IntegrationRuntimeEnableInteractiveQueryThenPoll(ctx, integrationruntimeenableinteractivequery.IntegrationRuntimeId(id), payload); err != nil {
			return fmt.Errorf("enabling interactive authoring on %s: %+v", id, err)
		}
	}

	return resourceDataFactoryIntegrationRuntimeAzureRead(d, meta)
}

func resourceDataFactoryIntegrationRuntimeAzureUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	disableInteractiveQueryClient := meta.(*clients.Client).DataFactory.IntegrationRuntimeDisableInteractiveQueryClient
	enableInteractiveQueryClient := meta.(*clients.Client).DataFactory.IntegrationRuntimeEnableInteractiveQueryClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, integrationruntimes.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `Model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `Properties` was nil", id)
	}

	props, ok := existing.Model.Properties.(integrationruntimes.ManagedIntegrationRuntime)
	if !ok {
		return fmt.Errorf("retrieving %s: asserting `IntegrationRuntime` as `ManagedIntegrationRuntime`", id)
	}

	if props.TypeProperties.ComputeProperties == nil {
		props.TypeProperties.ComputeProperties = &integrationruntimes.IntegrationRuntimeComputeProperties{}
	}

	if props.TypeProperties.ComputeProperties.DataFlowProperties == nil {
		props.TypeProperties.ComputeProperties.DataFlowProperties = &integrationruntimes.IntegrationRuntimeDataFlowProperties{}
	}
	dfProps := props.TypeProperties.ComputeProperties.DataFlowProperties

	if d.HasChange("cleanup_enabled") {
		dfProps.Cleanup = pointer.To(d.Get("cleanup_enabled").(bool))
	}

	if d.HasChange("compute_type") {
		dfProps.ComputeType = pointer.ToEnum[integrationruntimes.DataFlowComputeType](d.Get("compute_type").(string))
	}

	if d.HasChange("core_count") {
		dfProps.CoreCount = pointer.To(int64(d.Get("core_count").(int)))
	}

	if d.HasChange("description") {
		props.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("time_to_live_min") {
		dfProps.TimeToLive = pointer.To(int64(d.Get("time_to_live_min").(int)))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model, integrationruntimes.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if d.HasChange("interactive_authoring_time_to_live_in_minutes") {
		poller := pollers.NewPoller(custompollers.NewDataFactoryIntegrationRuntimeStatusPoller(client, *id), time.Second*5, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		err := poller.PollUntilDone(ctx)
		if err != nil {
			return fmt.Errorf("waiting for state change of %s: %+v", id, err)
		}

		if ttl := d.Get("interactive_authoring_time_to_live_in_minutes").(int); ttl > 0 {
			payload := integrationruntimeenableinteractivequery.EnableInteractiveQueryRequest{
				AutoTerminationMinutes: pointer.To(int64(ttl)),
			}

			if err := enableInteractiveQueryClient.IntegrationRuntimeEnableInteractiveQueryThenPoll(ctx, integrationruntimeenableinteractivequery.IntegrationRuntimeId(*id), payload); err != nil {
				return fmt.Errorf("enabling interactive authoring on %s: %+v", id, err)
			}
		} else {
			if err := disableInteractiveQueryClient.IntegrationRuntimeDisableInteractiveQueryThenPoll(ctx, integrationruntimedisableinteractivequery.IntegrationRuntimeId(*id)); err != nil {
				return fmt.Errorf("disabling interactive authoring on %s: %+v", id, err)
			}
		}
	}

	return resourceDataFactoryIntegrationRuntimeAzureRead(d, meta)
}

func resourceDataFactoryIntegrationRuntimeAzureRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName)

	resp, err := client.Get(ctx, *id, integrationruntimes.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.IntegrationRuntimeName)
	d.Set("data_factory_id", dataFactoryId.ID())

	if model := resp.Model; model != nil {
		runTime, ok := model.Properties.(integrationruntimes.ManagedIntegrationRuntime)
		if !ok {
			return fmt.Errorf("asserting `IntegrationRuntime` as `ManagedIntegrationRuntime` for %s", *id)
		}

		d.Set("description", pointer.From(runTime.Description))
		d.Set("virtual_network_enabled", runTime.ManagedVirtualNetwork != nil && runTime.ManagedVirtualNetwork.ReferenceName != "")

		if computeProps := runTime.TypeProperties.ComputeProperties; computeProps != nil {
			d.Set("location", location.NormalizeNilable(computeProps.Location))

			if dataFlowProps := computeProps.DataFlowProperties; dataFlowProps != nil {
				d.Set("compute_type", string(pointer.From(dataFlowProps.ComputeType)))
				d.Set("core_count", dataFlowProps.CoreCount)
				d.Set("time_to_live_min", dataFlowProps.TimeToLive)
				d.Set("cleanup_enabled", dataFlowProps.Cleanup)
			}
		}

		// This is currently non-functional, the API doesn't return the InteractiveQuery properties
		// See: https://github.com/Azure/azure-rest-api-specs/issues/39594
		if iqProps := runTime.TypeProperties.InteractiveQuery; iqProps != nil {
			d.Set("interactive_authoring_time_to_live_in_minutes", iqProps.AutoTerminationMinutes)
		}
	}

	return nil
}

func resourceDataFactoryIntegrationRuntimeAzureDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
