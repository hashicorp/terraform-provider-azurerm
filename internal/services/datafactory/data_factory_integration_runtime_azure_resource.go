// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

func resourceDataFactoryIntegrationRuntimeAzure() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceDataFactoryIntegrationRuntimeAzureCreateUpdate,
		Read:   resourceDataFactoryIntegrationRuntimeAzureRead,
		Update: resourceDataFactoryIntegrationRuntimeAzureCreateUpdate,
		Delete: resourceDataFactoryIntegrationRuntimeAzureDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationRuntimeID(id)
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

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
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
				Default:  string(datafactory.DataFlowComputeTypeGeneral),
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.DataFlowComputeTypeGeneral),
					string(datafactory.DataFlowComputeTypeComputeOptimized),
					string(datafactory.DataFlowComputeTypeMemoryOptimized),
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

	if !features.FourPointOhBeta() {
		resource.Schema["cleanup_enabled"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		}
	}

	return resource
}

func resourceDataFactoryIntegrationRuntimeAzureCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworks
	subscriptionId := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewIntegrationRuntimeID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_integration_runtime_azure", id.ID())
		}
	}

	description := d.Get("description").(string)

	managedIntegrationRuntime := datafactory.ManagedIntegrationRuntime{
		Description: &description,
		Type:        datafactory.TypeBasicIntegrationRuntimeTypeManaged,
		ManagedIntegrationRuntimeTypeProperties: &datafactory.ManagedIntegrationRuntimeTypeProperties{
			ComputeProperties: expandDataFactoryIntegrationRuntimeAzureComputeProperties(d),
		},
	}

	if d.Get("virtual_network_enabled").(bool) {
		virtualNetworkName, err := getManagedVirtualNetworkName(ctx, managedVirtualNetworksClient, id.SubscriptionId, id.ResourceGroup, id.FactoryName)
		if err != nil {
			return err
		}
		if virtualNetworkName == nil {
			return fmt.Errorf("virtual network feature for azure integration runtime is only available after managed virtual network for this data factory is enabled")
		}
		managedIntegrationRuntime.ManagedVirtualNetwork = &datafactory.ManagedVirtualNetworkReference{
			Type:          utils.String("ManagedVirtualNetworkReference"),
			ReferenceName: virtualNetworkName,
		}
	}

	basicIntegrationRuntime, _ := managedIntegrationRuntime.AsBasicIntegrationRuntime()

	integrationRuntime := datafactory.IntegrationRuntimeResource{
		Name:       &id.Name,
		Properties: basicIntegrationRuntime,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, integrationRuntime, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryIntegrationRuntimeAzureRead(d, meta)
}

func resourceDataFactoryIntegrationRuntimeAzureRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	managedIntegrationRuntime, convertSuccess := resp.Properties.AsManagedIntegrationRuntime()
	if !convertSuccess {
		return fmt.Errorf("converting managed integration runtime to Azure integration runtime %s", *id)
	}

	if managedIntegrationRuntime.Description != nil {
		d.Set("description", managedIntegrationRuntime.Description)
	}

	virtualNetworkEnabled := false
	if managedIntegrationRuntime.ManagedVirtualNetwork != nil && managedIntegrationRuntime.ManagedVirtualNetwork.ReferenceName != nil {
		virtualNetworkEnabled = true
	}
	d.Set("virtual_network_enabled", virtualNetworkEnabled)

	if computeProps := managedIntegrationRuntime.ComputeProperties; computeProps != nil {
		if location := computeProps.Location; location != nil {
			d.Set("location", location)
		}

		if dataFlowProps := computeProps.DataFlowProperties; dataFlowProps != nil {
			if computeType := &dataFlowProps.ComputeType; computeType != nil {
				d.Set("compute_type", string(*computeType))
			}

			if coreCount := dataFlowProps.CoreCount; coreCount != nil {
				d.Set("core_count", coreCount)
			}

			if timeToLive := dataFlowProps.TimeToLive; timeToLive != nil {
				d.Set("time_to_live_min", timeToLive)
			}

			d.Set("cleanup_enabled", dataFlowProps.Cleanup)
		}
	}

	return nil
}

func resourceDataFactoryIntegrationRuntimeAzureDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandDataFactoryIntegrationRuntimeAzureComputeProperties(d *pluginsdk.ResourceData) *datafactory.IntegrationRuntimeComputeProperties {
	location := azure.NormalizeLocation(d.Get("location").(string))
	coreCount := int32(d.Get("core_count").(int))
	timeToLiveMin := int32(d.Get("time_to_live_min").(int))

	cleanup := true
	if features.FourPointOhBeta() {
		cleanup = d.Get("cleanup_enabled").(bool)
	} else {
		// nolint staticcheck
		if v, ok := d.GetOkExists("cleanup_enabled"); ok {
			cleanup = v.(bool)
		}
	}

	return &datafactory.IntegrationRuntimeComputeProperties{
		Location: &location,
		DataFlowProperties: &datafactory.IntegrationRuntimeDataFlowProperties{
			ComputeType: datafactory.DataFlowComputeType(d.Get("compute_type").(string)),
			CoreCount:   &coreCount,
			TimeToLive:  &timeToLiveMin,
			Cleanup:     utils.Bool(cleanup),
		},
	}
}
