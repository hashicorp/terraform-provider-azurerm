package machinelearning

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2021-07-01/machinelearningservices"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
				ValidateFunc: validation.StringInSlice([]string{string(machinelearningservices.VMPriorityDedicated), string(machinelearningservices.VMPriorityLowPriority)}, false),
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptionalForceNew(),

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

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"ssh": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"admin_username": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
						},
						"admin_password": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: []string{"ssh.0.admin_password", "ssh.0.key_value"},
						},
						"key_value": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: []string{"ssh.0.admin_password", "ssh.0.key_value"},
						},
					},
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ssh_public_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true, // `ssh_public_access_enabled` sets to `true` by default even if unspecified
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
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceID, err := parse.WorkspaceID(d.Get("machine_learning_workspace_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewComputeClusterID(workspaceID.SubscriptionId, workspaceID.ResourceGroup, workspaceID.Name, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_machine_learning_compute_cluster", id.ID())
	}

	computeClusterAmlComputeProperties := machinelearningservices.AmlComputeProperties{
		VMSize:                 utils.String(d.Get("vm_size").(string)),
		VMPriority:             machinelearningservices.VMPriority(d.Get("vm_priority").(string)),
		ScaleSettings:          expandScaleSettings(d.Get("scale_settings").([]interface{})),
		UserAccountCredentials: expandUserAccountCredentials(d.Get("ssh").([]interface{})),
	}

	computeClusterAmlComputeProperties.RemoteLoginPortPublicAccess = machinelearningservices.RemoteLoginPortPublicAccessDisabled
	if d.Get("ssh_public_access_enabled").(bool) {
		computeClusterAmlComputeProperties.RemoteLoginPortPublicAccess = machinelearningservices.RemoteLoginPortPublicAccessEnabled
	}

	if subnetId, ok := d.GetOk("subnet_resource_id"); ok && subnetId.(string) != "" {
		computeClusterAmlComputeProperties.Subnet = &machinelearningservices.ResourceID{ID: utils.String(subnetId.(string))}
	}

	computeClusterProperties := machinelearningservices.AmlCompute{
		Properties:       &computeClusterAmlComputeProperties,
		ComputeLocation:  utils.String(d.Get("location").(string)),
		Description:      utils.String(d.Get("description").(string)),
		DisableLocalAuth: utils.Bool(!d.Get("local_auth_enabled").(bool)),
	}

	// Get SKU from Workspace
	workspace, err := mlWorkspacesClient.Get(ctx, workspaceID.ResourceGroup, workspaceID.Name)
	if err != nil {
		return err
	}

	identity, err := expandIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	computeClusterParameters := machinelearningservices.ComputeResource{
		Properties: computeClusterProperties,
		Identity:   identity,
		Location:   computeClusterProperties.ComputeLocation,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku:        workspace.Sku,
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName, computeClusterParameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceComputeClusterRead(d, meta)
}

func resourceComputeClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ComputeClusterID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Compute Cluster ID `%q`: %+v", d.Id(), err)
	}

	computeResource, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName)
	if err != nil {
		if utils.ResponseWasNotFound(computeResource.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ComputeName)

	workspaceId := parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	// use ComputeResource to get to AKS Cluster ID and other properties
	computeCluster, isComputeCluster := machinelearningservices.BasicCompute.AsAmlCompute(computeResource.Properties)
	if !isComputeCluster {
		return fmt.Errorf("compute resource %s is not an AML Compute cluster", id.ComputeName)
	}

	if computeCluster.DisableLocalAuth != nil {
		d.Set("local_auth_enabled", !*computeCluster.DisableLocalAuth)
	}
	d.Set("description", computeCluster.Description)
	if props := computeCluster.Properties; props != nil {
		d.Set("vm_size", props.VMSize)
		d.Set("vm_priority", props.VMPriority)
		d.Set("scale_settings", flattenScaleSettings(props.ScaleSettings))
		d.Set("ssh", flattenUserAccountCredentials(props.UserAccountCredentials))
		if props.Subnet != nil {
			d.Set("subnet_resource_id", props.Subnet.ID)
		}

		switch props.RemoteLoginPortPublicAccess {
		case machinelearningservices.RemoteLoginPortPublicAccessNotSpecified:
			d.Set("ssh_public_access_enabled", nil)
		case machinelearningservices.RemoteLoginPortPublicAccessEnabled:
			d.Set("ssh_public_access_enabled", true)
		case machinelearningservices.RemoteLoginPortPublicAccessDisabled:
			d.Set("ssh_public_access_enabled", false)
		}
	}

	if location := computeResource.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity, err := flattenIdentity(computeResource.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, computeResource.Tags)
}

func resourceComputeClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ComputeClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.ComputeName, machinelearningservices.UnderlyingResourceActionDelete)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}
	return nil
}

func expandScaleSettings(input []interface{}) *machinelearningservices.ScaleSettings {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	maxNodeCount := int32(v["max_node_count"].(int))
	minNodeCount := int32(v["min_node_count"].(int))
	scaleDownNodes := v["scale_down_nodes_after_idle_duration"].(string)

	return &machinelearningservices.ScaleSettings{
		MaxNodeCount:                &maxNodeCount,
		MinNodeCount:                &minNodeCount,
		NodeIdleTimeBeforeScaleDown: &scaleDownNodes,
	}
}

func expandUserAccountCredentials(input []interface{}) *machinelearningservices.UserAccountCredentials {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	return &machinelearningservices.UserAccountCredentials{
		AdminUserName:         utils.String(v["admin_username"].(string)),
		AdminUserPassword:     utils.String(v["admin_password"].(string)),
		AdminUserSSHPublicKey: utils.String(v["key_value"].(string)),
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

func flattenUserAccountCredentials(credentials *machinelearningservices.UserAccountCredentials) interface{} {
	if credentials == nil {
		return []interface{}{}
	}
	var username string
	if credentials.AdminUserName != nil {
		username = *credentials.AdminUserName
	}

	var adminPassword string
	if credentials.AdminUserPassword != nil {
		adminPassword = *credentials.AdminUserPassword
	}

	var sshPublicKey string
	if credentials.AdminUserSSHPublicKey != nil {
		sshPublicKey = *credentials.AdminUserSSHPublicKey
	}

	return []interface{}{
		map[string]interface{}{
			"admin_username": username,
			"admin_password": adminPassword,
			"key_value":      sshPublicKey,
		},
	}
}
