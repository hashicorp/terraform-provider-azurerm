package logic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic/validate"
	networkParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIntegrationServiceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIntegrationServiceEnvironmentCreateUpdate,
		Read:   resourceArmIntegrationServiceEnvironmentRead,
		Update: resourceArmIntegrationServiceEnvironmentCreateUpdate,
		Delete: resourceArmIntegrationServiceEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Hour),
			Delete: schema.DefaultTimeout(5 * time.Hour),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationServiceEnvironmentName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			// Maximum scale units that you can add	10 - https://docs.microsoft.com/en-US/azure/logic-apps/logic-apps-limits-and-config#integration-service-environment-ise
			// Developer Always 0 capacity
			"sku_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Developer_0",
				ValidateFunc: validation.StringInSlice([]string{
					"Developer_0",
					"Premium_0",
					"Premium_1",
					"Premium_2",
					"Premium_3",
					"Premium_4",
					"Premium_5",
					"Premium_6",
					"Premium_7",
					"Premium_8",
					"Premium_9",
					"Premium_10",
				}, false),
			},

			"access_endpoint_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // The access end point type cannot be changed once the integration service environment is provisioned.
				ValidateFunc: validation.StringInSlice([]string{
					string(logic.IntegrationServiceEnvironmentAccessEndpointTypeInternal),
					string(logic.IntegrationServiceEnvironmentAccessEndpointTypeExternal),
				}, false),
			},

			"virtual_network_subnet_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true, // The network configuration subnets cannot be updated after integration service environment is created.
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.ValidateSubnetID,
				},
				MinItems: 4,
				MaxItems: 4,
			},

			"connector_endpoint_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"connector_outbound_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"workflow_endpoint_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"workflow_outbound_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("sku_name", func(old, new, meta interface{}) bool {
				oldSku := strings.Split(old.(string), "_")
				newSku := strings.Split(new.(string), "_")
				// The SKU cannot be changed once integration service environment has been provisioned. -> we need ForceNew
				return oldSku[0] != newSku[0]
			}),
		),
	}
}

func resourceArmIntegrationServiceEnvironmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationServiceEnvironmentClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Integration Service Environment creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing Integration Service Environment %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_integration_service_environment", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	accessEndpointType := d.Get("access_endpoint_type").(string)
	virtualNetworkSubnetIds := d.Get("virtual_network_subnet_ids").(*schema.Set).List()
	t := d.Get("tags").(map[string]interface{})

	sku, err := expandIntegrationServiceEnvironmentSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name` for Integration Service Environment %q (Resource Group %q): %v", name, resourceGroup, err)
	}

	integrationServiceEnvironment := logic.IntegrationServiceEnvironment{
		Name:     &name,
		Location: &location,
		Properties: &logic.IntegrationServiceEnvironmentProperties{
			NetworkConfiguration: &logic.NetworkConfiguration{
				AccessEndpoint: &logic.IntegrationServiceEnvironmentAccessEndpoint{
					Type: logic.IntegrationServiceEnvironmentAccessEndpointType(accessEndpointType),
				},
				Subnets: expandSubnetResourceID(virtualNetworkSubnetIds),
			},
		},
		Sku:  sku,
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, integrationServiceEnvironment)
	if err != nil {
		return fmt.Errorf("creating/updating Integration Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of Integration Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Integration Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read Integration Service Environment %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmIntegrationServiceEnvironmentRead(d, meta)
}

func resourceArmIntegrationServiceEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationServiceEnvironmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationServiceEnvironmentID(d.Id())
	if err != nil {
		return err
	}

	name := id.Name
	resourceGroup := id.ResourceGroup

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Integration Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if err := d.Set("sku_name", flattenIntegrationServiceEnvironmentSkuName(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku_name`: %+v", err)
	}

	if props := resp.Properties; props != nil {
		if netCfg := props.NetworkConfiguration; netCfg != nil {
			if accessEndpoint := netCfg.AccessEndpoint; accessEndpoint != nil {
				d.Set("access_endpoint_type", accessEndpoint.Type)
			}

			d.Set("virtual_network_subnet_ids", flattenSubnetResourceID(netCfg.Subnets))
		}

		if props.EndpointsConfiguration == nil || props.EndpointsConfiguration.Connector == nil {
			d.Set("connector_endpoint_ip_addresses", []interface{}{})
			d.Set("connector_outbound_ip_addresses", []interface{}{})
		} else {
			d.Set("connector_endpoint_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Connector.AccessEndpointIPAddresses))
			d.Set("connector_outbound_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Connector.OutgoingIPAddresses))
		}

		if props.EndpointsConfiguration == nil || props.EndpointsConfiguration.Workflow == nil {
			d.Set("workflow_endpoint_ip_addresses", []interface{}{})
			d.Set("workflow_outbound_ip_addresses", []interface{}{})
		} else {
			d.Set("workflow_endpoint_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Workflow.AccessEndpointIPAddresses))
			d.Set("workflow_outbound_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Workflow.OutgoingIPAddresses))
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmIntegrationServiceEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationServiceEnvironmentClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationServiceEnvironmentID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Integration Service Environment ID `%q`: %+v", d.Id(), err)
	}

	name := id.Name
	resourceGroup := id.ResourceGroup

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("retrieving Integration Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// Get subnet IDs before delete
	subnetIDs := getSubnetIDs(&resp)

	// Not optimal behavior for now
	// It deletes synchronously and resource is not available anymore after return from delete operation
	// Next, after return - delete operation is still in progress in the background and is still occupying subnets.
	// As workaround we are checking on all involved subnets presence of serviceAssociationLink and resourceNavigationLink
	// If the operation fails we are lost. We do not have original resource and we cannot resume delete operation.
	// User has to wait for completion of delete operation in the background.
	// It would be great to have async call with future struct
	if resp, err := client.Delete(ctx, resourceGroup, name); err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting Integration Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:                   []string{string(logic.WorkflowProvisioningStateDeleting)},
		Target:                    []string{string(logic.WorkflowProvisioningStateDeleted)},
		MinTimeout:                5 * time.Minute,
		Refresh:                   integrationServiceEnvironmentDeleteStateRefreshFunc(ctx, meta.(*clients.Client), d.Id(), subnetIDs),
		Timeout:                   d.Timeout(schema.TimeoutDelete),
		ContinuousTargetOccurence: 1,
		NotFoundChecks:            1,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for deletion of Integration Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func flattenIntegrationServiceEnvironmentSkuName(input *logic.IntegrationServiceEnvironmentSku) string {
	if input == nil {
		return ""
	}

	return fmt.Sprintf("%s_%d", string(input.Name), *input.Capacity)
}

func expandIntegrationServiceEnvironmentSkuName(skuName string) (*logic.IntegrationServiceEnvironmentSku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("sku_name (%s) has the wrong number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var sku logic.IntegrationServiceEnvironmentSkuName
	switch parts[0] {
	case string(logic.IntegrationServiceEnvironmentSkuNameDeveloper):
		sku = logic.IntegrationServiceEnvironmentSkuNameDeveloper
	case string(logic.IntegrationServiceEnvironmentSkuNamePremium):
		sku = logic.IntegrationServiceEnvironmentSkuNamePremium
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku %s", skuName, parts[0])
	}

	capacity, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("cannot convert sku_name %s capacity %s to int", skuName, parts[1])
	}

	if sku != logic.IntegrationServiceEnvironmentSkuNamePremium && capacity > 0 {
		return nil, fmt.Errorf("`capacity` can only be greater than zero for `sku_name` `Premium`")
	}

	return &logic.IntegrationServiceEnvironmentSku{
		Name:     sku,
		Capacity: utils.Int32(int32(capacity)),
	}, nil
}

func expandSubnetResourceID(input []interface{}) *[]logic.ResourceReference {
	results := make([]logic.ResourceReference, 0)
	for _, item := range input {
		results = append(results, logic.ResourceReference{
			ID: utils.String(item.(string)),
		})
	}
	return &results
}

func flattenSubnetResourceID(input *[]logic.ResourceReference) []interface{} {
	subnetIDs := make([]interface{}, 0)
	if input == nil {
		return subnetIDs
	}

	for _, resourceRef := range *input {
		if resourceRef.ID == nil || *resourceRef.ID == "" {
			continue
		}

		subnetIDs = append(subnetIDs, resourceRef.ID)
	}

	return subnetIDs
}

func getSubnetIDs(input *logic.IntegrationServiceEnvironment) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	if props := input.Properties; props != nil {
		if netCfg := props.NetworkConfiguration; netCfg != nil {
			return flattenSubnetResourceID(netCfg.Subnets)
		}
	}

	return results
}

func integrationServiceEnvironmentDeleteStateRefreshFunc(ctx context.Context, client *clients.Client, iseID string, subnetIDs []interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		linkExists, err := linkExists(ctx, client, iseID, subnetIDs)
		if err != nil {
			return string(logic.WorkflowProvisioningStateDeleting), string(logic.WorkflowProvisioningStateDeleting), err
		}

		if linkExists {
			return string(logic.WorkflowProvisioningStateDeleting), string(logic.WorkflowProvisioningStateDeleting), nil
		}

		return string(logic.WorkflowProvisioningStateDeleted), string(logic.WorkflowProvisioningStateDeleted), nil
	}
}

func linkExists(ctx context.Context, client *clients.Client, iseID string, subnetIDs []interface{}) (bool, error) {
	for _, subnetID := range subnetIDs {
		if subnetID == nil {
			continue
		}

		id := *(subnetID.(*string))
		log.Printf("Checking links on subnetID: %q\n", id)

		hasLink, err := serviceAssociationLinkExists(ctx, client.Network.ServiceAssociationLinkClient, iseID, id)
		if err != nil {
			return false, err
		}

		if hasLink {
			return true, nil
		} else {
			hasLink, err := resourceNavigationLinkExists(ctx, client.Network.ResourceNavigationLinkClient, id)
			if err != nil {
				return false, err
			}

			if hasLink {
				return true, nil
			}
		}
	}

	return false, nil
}

func serviceAssociationLinkExists(ctx context.Context, client *network.ServiceAssociationLinksClient, iseID string, subnetID string) (bool, error) {
	id, err := networkParse.SubnetID(subnetID)
	if err != nil {
		return false, err
	}

	resp, err := client.List(ctx, id.ResourceGroup, id.VirtualNetworkName, id.Name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return false, nil
		}
		return false, fmt.Errorf("retrieving Service Association Links from Virtual Network %q, subnet %q (Resource Group %q): %+v", id.VirtualNetworkName, id.Name, id.ResourceGroup, err)
	}

	if resp.Value != nil {
		for _, link := range *resp.Value {
			if link.ServiceAssociationLinkPropertiesFormat != nil && link.ServiceAssociationLinkPropertiesFormat.Link != nil {
				if strings.EqualFold(iseID, *link.ServiceAssociationLinkPropertiesFormat.Link) {
					log.Printf("Has Service Association Link: %q\n", *link.ID)
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func resourceNavigationLinkExists(ctx context.Context, client *network.ResourceNavigationLinksClient, subnetID string) (bool, error) {
	id, err := networkParse.SubnetID(subnetID)
	if err != nil {
		return false, err
	}

	resp, err := client.List(ctx, id.ResourceGroup, id.VirtualNetworkName, id.Name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return false, nil
		}
		return false, fmt.Errorf("retrieving Resource Navigation Links from Virtual Network %q, subnet %q (Resource Group %q): %+v", id.VirtualNetworkName, id.Name, id.ResourceGroup, err)
	}

	if resp.Value != nil {
		for _, link := range *resp.Value {
			log.Printf("Has Resource Navigation Link: %q\n", *link.ID)
			return true, nil
		}
	}

	return false, nil
}
