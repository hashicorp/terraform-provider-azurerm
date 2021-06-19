package machinelearning

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2020-04-01/machinelearningservices"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/machinelearning/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceComputeCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceComputeClusterCreate,
		Read:   resourceComputeClusterRead,
		Delete: resourceComputeClusterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ComputeClusterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"machine_learning_workspace_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"vm_size": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vm_priority": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(machinelearningservices.Dedicated), string(machinelearningservices.LowPriority)}, false),
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(machinelearningservices.SystemAssigned),
							}, false),
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"scale_settings": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"max_node_count": {
							Type:     pluginsdk.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"min_node_count": {
							Type:     pluginsdk.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"scale_down_nodes_after_idle_duration": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"subnet_resource_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceComputeClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	mlWorkspacesClient := meta.(*clients.Client).MachineLearning.WorkspacesClient
	mlComputeClient := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	// Get Machine Learning Workspace Name and Resource Group from ID
	workspaceID, err := parse.WorkspaceID(d.Get("machine_learning_workspace_id").(string))
	if err != nil {
		return err
	}

	existing, err := mlComputeClient.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("error checking for existing Compute Cluster %q in Workspace %q (Resource Group %q): %s",
				name, workspaceID.Name, workspaceID.ResourceGroup, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_machine_learning_compute_cluster", *existing.ID)
	}

	computeClusterProperties := machinelearningservices.AmlCompute{
		Properties: &machinelearningservices.AmlComputeProperties{
			VMSize:        utils.String(d.Get("vm_size").(string)),
			VMPriority:    machinelearningservices.VMPriority(d.Get("vm_priority").(string)),
			ScaleSettings: expandScaleSettings(d.Get("scale_settings").([]interface{})),
			Subnet:        &machinelearningservices.ResourceID{ID: utils.String(d.Get("subnet_resource_id").(string))},
		},
		ComputeLocation: utils.String(d.Get("location").(string)),
		Description:     utils.String(d.Get("description").(string)),
	}

	amlComputeProperties, isAmlCompute := (machinelearningservices.BasicCompute).AsAmlCompute(computeClusterProperties)
	if !isAmlCompute {
		return fmt.Errorf("no compute cluster")
	}

	// Get SKU from Workspace
	workspace, err := mlWorkspacesClient.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name)
	if err != nil {
		return err
	}

	computeClusterParameters := machinelearningservices.ComputeResource{
		Properties: amlComputeProperties,
		Identity:   expandComputeClusterIdentity(d.Get("identity").([]interface{})),
		Location:   computeClusterProperties.ComputeLocation,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku:        workspace.Sku,
	}

	future, err := mlComputeClient.CreateOrUpdate(ctx, workspaceID.ResourceGroup, workspaceID.Name, name, computeClusterParameters)
	if err != nil {
		return fmt.Errorf("creating Compute Cluster %q in workspace %q (Resource Group %q): %+v",
			name, workspaceID.Name, workspaceID.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, mlComputeClient.Client); err != nil {
		return fmt.Errorf("waiting for creation of Compute Cluster %q in workspace %q (Resource Group %q): %+v",
			name, workspaceID.Name, workspaceID.ResourceGroup, err)
	}

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := parse.NewComputeClusterID(subscriptionId, workspaceID.ResourceGroup, workspaceID.Name, name)
	d.SetId(id.ID())

	return resourceComputeClusterRead(d, meta)
}

func resourceComputeClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	mlComputeClient := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ComputeClusterID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Compute Cluster ID `%q`: %+v", d.Id(), err)
	}

	computeResource, err := mlComputeClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName)
	if err != nil {
		if utils.ResponseWasNotFound(computeResource.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Compute Cluster %q in Workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroup, err)
	}

	d.Set("name", id.ComputeName)

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	workspaceId := parse.NewWorkspaceID(subscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	// use ComputeResource to get to AKS Cluster ID and other properties
	computeCluster, isComputeCluster := (machinelearningservices.BasicCompute).AsAmlCompute(computeResource.Properties)
	if !isComputeCluster {
		return fmt.Errorf("compute resource %s is not an Aml Compute cluster", id.ComputeName)
	}

	d.Set("vm_size", computeCluster.Properties.VMSize)
	d.Set("vm_priority", computeCluster.Properties.VMPriority)
	d.Set("scale_settings", flattenScaleSettings(computeCluster.Properties.ScaleSettings))
	d.Set("subnet_resource_id", computeCluster.Properties.Subnet.ID)

	if location := computeResource.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("identity", flattenComputeClusterIdentity(computeResource.Identity)); err != nil {
		return fmt.Errorf("flattening identity on Workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.ResourceGroup, err)
	}

	return tags.FlattenAndSet(d, computeResource.Tags)
}

func resourceComputeClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	mlComputeClient := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.ComputeClusterID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Compute Cluster ID `%q`: %+v", d.Id(), err)
	}
	future, err := mlComputeClient.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName, machinelearningservices.Detach)
	if err != nil {
		return fmt.Errorf("deleting Compute Cluster %q in workspace %q (Resource Group %q): %+v", id.ComputeName, id.WorkspaceName, id.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, mlComputeClient.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Compute Cluster %q in workspace %q (Resource Group %q): %+v", id.ComputeName, id.WorkspaceName, id.ResourceGroup, err)
	}
	return nil
}

func expandScaleSettings(input []interface{}) *machinelearningservices.ScaleSettings {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	max_node_count := int32(v["max_node_count"].(int))
	min_node_count := int32(v["min_node_count"].(int))
	scale_down_nodes_after_idle_duration := v["scale_down_nodes_after_idle_duration"].(string)

	return &machinelearningservices.ScaleSettings{
		MaxNodeCount:                &max_node_count,
		MinNodeCount:                &min_node_count,
		NodeIdleTimeBeforeScaleDown: &scale_down_nodes_after_idle_duration,
	}
}

func flattenScaleSettings(scaleSettings *machinelearningservices.ScaleSettings) []interface{} {
	if scaleSettings == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"max_node_count":                       scaleSettings.MaxNodeCount,
			"min_node_count":                       scaleSettings.MinNodeCount,
			"scale_down_nodes_after_idle_duration": scaleSettings.NodeIdleTimeBeforeScaleDown,
		},
	}
}

func expandComputeClusterIdentity(input []interface{}) *machinelearningservices.Identity {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &machinelearningservices.Identity{
		Type: machinelearningservices.ResourceIdentityType(v["type"].(string)),
	}
}

func flattenComputeClusterIdentity(identity *machinelearningservices.Identity) []interface{} {
	if identity == nil {
		return []interface{}{}
	}

	principalID := ""
	if identity.PrincipalID != nil {
		principalID = *identity.PrincipalID
	}

	tenantID := ""
	if identity.TenantID != nil {
		tenantID = *identity.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(identity.Type),
			"principal_id": principalID,
			"tenant_id":    tenantID,
		},
	}
}
