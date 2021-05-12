package machinelearning

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2020-04-01/machinelearningservices"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/machinelearning/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceComputeCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeClusterCreate,
		Read:   resourceComputeClusterRead,
		Update: resourceComputeClusterUpdate,
		Delete: resourceComputeClusterDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ComputeClusterID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"machine_learning_workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"vm_size": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vm_priority": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"identity": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(machinelearningservices.SystemAssigned),
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"scale_settings": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_node_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"min_node_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"node_idle_time_before_scale_down": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"subnet_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceComputeClusterCreate(d *schema.ResourceData, meta interface{}) error {
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

	// AML Compute Cluster configuration (not needed here presumably)
	vmSize := d.Get("vm_size").(string)
	vmPriority := d.Get("vm_priority").(string)

	scaleSettings := expandScaleSettings(d.Get("scale_settings").([]interface{}))

	subnetResourceId := d.Get("subnet_resource_id").(string)

	description := d.Get("description").(string)

	location := d.Get("location").(string)

	identity := d.Get("identity").([]interface{})

	// Get SKU from Workspace
	workspace, err := mlWorkspacesClient.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name)
	if err != nil {
		return err
	}
	sku := workspace.Sku

	t := d.Get("tags").(map[string]interface{})

	subnetId := machinelearningservices.ResourceID{ID: utils.String(subnetResourceId)}

	computeClusterProperties := machinelearningservices.AmlCompute{
		// Properties - AML Compute properties
		Properties: &machinelearningservices.AmlComputeProperties{
			// VMSize - Virtual Machine Size
			VMSize: utils.String(vmSize),
			// VMPriority - Virtual Machine priority. Possible values include: 'Dedicated', 'LowPriority'
			VMPriority: machinelearningservices.VMPriority(vmPriority),
			// ScaleSettings - Scale settings for AML Compute
			ScaleSettings: scaleSettings,
			// UserAccountCredentials - Credentials for an administrator user account that will be created on each compute node.
			// Subnet - Virtual network subnet resource ID the compute nodes belong to.
			Subnet: &subnetId,
			// RemoteLoginPortPublicAccess - State of the public SSH port. Possible values are: Disabled - Indicates that the public ssh port is closed on all nodes of the cluster. Enabled - Indicates that the public ssh port is open on all nodes of the cluster. NotSpecified - Indicates that the public ssh port is closed on all nodes of the cluster if VNet is defined, else is open all public nodes. It can be default only during cluster creation time, after creation it will be either enabled or disabled. Possible values include: 'RemoteLoginPortPublicAccessEnabled', 'RemoteLoginPortPublicAccessDisabled', 'RemoteLoginPortPublicAccessNotSpecified'
			// AllocationState - READ-ONLY; Allocation state of the compute. Possible values are: steady - Indicates that the compute is not resizing. There are no changes to the number of compute nodes in the compute in progress. A compute enters this state when it is created and when no operations are being performed on the compute to change the number of compute nodes. resizing - Indicates that the compute is resizing; that is, compute nodes are being added to or removed from the compute. Possible values include: 'Steady', 'Resizing'
			// AllocationStateTransitionTime - READ-ONLY; The time at which the compute entered its current allocation state.
			// Errors - READ-ONLY; Collection of errors encountered by various compute nodes during node setup.
			// CurrentNodeCount - READ-ONLY; The number of compute nodes currently assigned to the compute.
			// TargetNodeCount - READ-ONLY; The target number of compute nodes for the compute. If the allocationState is resizing, this property denotes the target node count for the ongoing resize operation. If the allocationState is steady, this property denotes the target node count for the previous resize operation.
			// NodeStateCounts - READ-ONLY; Counts of various node states on the compute.
		},
		// ComputeLocation - Location for the underlying compute
		ComputeLocation: &location,
		// ProvisioningState - READ-ONLY; The provision state of the cluster. Valid values are Unknown, Updating, Provisioning, Succeeded, and Failed. Possible values include: 'ProvisioningStateUnknown', 'ProvisioningStateUpdating', 'ProvisioningStateCreating', 'ProvisioningStateDeleting', 'ProvisioningStateSucceeded', 'ProvisioningStateFailed', 'ProvisioningStateCanceled'
		// Description - The description of the Machine Learning compute.
		Description: &description,
		// CreatedOn - READ-ONLY; The date and time when the compute was created.
		// ModifiedOn - READ-ONLY; The date and time when the compute was last modified.
		// ResourceID - ARM resource id of the underlying compute
		// ProvisioningErrors - READ-ONLY; Errors during provisioning
		// IsAttachedCompute - READ-ONLY; Indicating whether the compute was provisioned by user and brought from outside if true, or machine learning service provisioned it if false.
		// ComputeType - Possible values include: 'ComputeTypeCompute', 'ComputeTypeAKS1', 'ComputeTypeAmlCompute1', 'ComputeTypeVirtualMachine1', 'ComputeTypeHDInsight1', 'ComputeTypeDataFactory1', 'ComputeTypeDatabricks1', 'ComputeTypeDataLakeAnalytics1'
		ComputeType: "ComputeTypeAmlCompute1",
	}

	amlComputeProperties, isAmlCompute := (machinelearningservices.BasicCompute).AsAmlCompute(computeClusterProperties)
	if !isAmlCompute {
		return fmt.Errorf("error: No Compute cluster")
	}

	computeClusterParameters := machinelearningservices.ComputeResource{
		// Properties - Compute properties
		Properties: amlComputeProperties,
		// ID - READ-ONLY; Specifies the resource ID.
		// Name - READ-ONLY; Specifies the name of the resource.
		// Identity - The identity of the resource.
		Identity: expandComputeClusterIdentity(identity),
		// Location - Specifies the location of the resource.
		Location: &location,
		// Type - READ-ONLY; Specifies the type of the resource.
		// Tags - Contains resource tags defined as key/value pairs.
		Tags: tags.Expand(t),
		// Sku - The sku of the workspace.
		Sku: sku,
	}

	future, err := mlComputeClient.CreateOrUpdate(ctx, workspaceID.ResourceGroup, workspaceID.Name, name, computeClusterParameters)
	if err != nil {
		return fmt.Errorf("error creating Compute Cluster %q in workspace %q (Resource Group %q): %+v",
			name, workspaceID.Name, workspaceID.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, mlComputeClient.Client); err != nil {
		return fmt.Errorf("error waiting for creation of Compute Cluster %q in workspace %q (Resource Group %q): %+v",
			name, workspaceID.Name, workspaceID.ResourceGroup, err)
	}

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := parse.NewComputeClusterID(subscriptionId, workspaceID.ResourceGroup, workspaceID.Name, name)
	d.SetId(id.ID())

	return resourceComputeClusterRead(d, meta)
}

func resourceComputeClusterRead(d *schema.ResourceData, meta interface{}) error {
	mlComputeClient := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	mlWorkspacesClient := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ComputeClusterID(d.Id())
	if err != nil {
		return fmt.Errorf("error parsing Compute Cluster ID `%q`: %+v", d.Id(), err)
	}

	computeResource, err := mlComputeClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName)
	if err != nil {
		if utils.ResponseWasNotFound(computeResource.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error making Read request on Compute Cluster %q in Workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.WorkspaceName, id.ResourceGroup, err)
	}

	d.Set("name", id.ComputeName)

	workspaceResp, err := mlWorkspacesClient.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		return err
	}
	d.Set("machine_learning_workspace_id", workspaceResp.ID)

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
		return fmt.Errorf("error flattening identity on Workspace %q (Resource Group %q): %+v",
			id.ComputeName, id.ResourceGroup, err)
	}

	return tags.FlattenAndSet(d, computeResource.Tags)
}

func resourceComputeClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ComputeClusterID(d.Id())
	if err != nil {
		return err
	}

	update := machinelearningservices.ClusterUpdateParameters{
		ClusterUpdateProperties: &machinelearningservices.ClusterUpdateProperties{},
	}

	if d.HasChange("scale_settings") {
		scaleSettings := d.Get("scale_settings").([]interface{})
		update.ScaleSettings = expandScaleSettings(scaleSettings)
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName, update); err != nil {
		return fmt.Errorf("error updating Machine Learning Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroup, err)
	}

	return resourceComputeClusterRead(d, meta)
}

func resourceComputeClusterDelete(d *schema.ResourceData, meta interface{}) error {
	mlComputeClient := meta.(*clients.Client).MachineLearning.MachineLearningComputeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.ComputeClusterID(d.Id())
	if err != nil {
		return fmt.Errorf("error parsing Compute Cluster ID `%q`: %+v", d.Id(), err)
	}
	underlying_resource_action := machinelearningservices.Detach
	future, err := mlComputeClient.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName, underlying_resource_action)
	if err != nil {
		return fmt.Errorf("error deleting Compute Cluster %q in workspace %q (Resource Group %q): %+v", id.ComputeName, id.WorkspaceName, id.ResourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, mlComputeClient.Client); err != nil {
		return fmt.Errorf("error waiting for deletion of Compute Cluster %q in workspace %q (Resource Group %q): %+v", id.ComputeName, id.WorkspaceName, id.ResourceGroup, err)
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
	node_idle_time_before_scale_down := string(v["node_idle_time_before_scale_down"].(string))

	return &machinelearningservices.ScaleSettings{
		MaxNodeCount:                &max_node_count,
		MinNodeCount:                &min_node_count,
		NodeIdleTimeBeforeScaleDown: &node_idle_time_before_scale_down,
	}
}

func flattenScaleSettings(scaleSettings *machinelearningservices.ScaleSettings) []interface{} {
	if scaleSettings == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"max_node_count":                   scaleSettings.MaxNodeCount,
			"min_node_count":                   scaleSettings.MinNodeCount,
			"node_idle_time_before_scale_down": scaleSettings.NodeIdleTimeBeforeScaleDown,
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

	t := string(identity.Type)

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
			"type":         t,
			"principal_id": principalID,
			"tenant_id":    tenantID,
		},
	}
}
