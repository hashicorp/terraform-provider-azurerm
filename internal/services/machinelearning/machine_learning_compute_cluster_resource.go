// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
			_, err := machinelearningcomputes.ParseComputeID(id)
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

			"location": commonschema.Location(),

			"vm_size": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vm_priority": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(machinelearningcomputes.VMPriorityDedicated), string(machinelearningcomputes.VMPriorityLowPriority)}, false),
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

			"node_public_ip_enabled": {
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

			"tags": commonschema.TagsForceNew(),
		},
	}
}

func resourceComputeClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	mlWorkspacesClient := meta.(*clients.Client).MachineLearning.Workspaces
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceID, err := workspaces.ParseWorkspaceID(d.Get("machine_learning_workspace_id").(string))
	if err != nil {
		return err
	}

	id := machinelearningcomputes.NewComputeID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, d.Get("name").(string))

	existing, err := client.ComputeGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_machine_learning_compute_cluster", id.ID())
	}

	if !d.Get("node_public_ip_enabled").(bool) && d.Get("subnet_resource_id").(string) == "" {
		return fmt.Errorf("`subnet_resource_id` must be set if `node_public_ip_enabled` is set to `false`")
	}

	vmPriority := machinelearningcomputes.VMPriority(d.Get("vm_priority").(string))
	computeClusterAmlComputeProperties := machinelearningcomputes.AmlComputeProperties{
		VMSize:                 utils.String(d.Get("vm_size").(string)),
		VMPriority:             &vmPriority,
		ScaleSettings:          expandScaleSettings(d.Get("scale_settings").([]interface{})),
		UserAccountCredentials: expandUserAccountCredentials(d.Get("ssh").([]interface{})),
		EnableNodePublicIP:     pointer.To(d.Get("node_public_ip_enabled").(bool)),
	}

	computeClusterAmlComputeProperties.RemoteLoginPortPublicAccess = utils.ToPtr(machinelearningcomputes.RemoteLoginPortPublicAccessDisabled)
	if d.Get("ssh_public_access_enabled").(bool) {
		computeClusterAmlComputeProperties.RemoteLoginPortPublicAccess = utils.ToPtr(machinelearningcomputes.RemoteLoginPortPublicAccessEnabled)
	}

	if subnetId, ok := d.GetOk("subnet_resource_id"); ok && subnetId.(string) != "" {
		computeClusterAmlComputeProperties.Subnet = &machinelearningcomputes.ResourceId{Id: subnetId.(string)}
	}

	computeClusterProperties := machinelearningcomputes.AmlCompute{
		Properties:       &computeClusterAmlComputeProperties,
		ComputeLocation:  utils.String(d.Get("location").(string)),
		Description:      utils.String(d.Get("description").(string)),
		DisableLocalAuth: utils.Bool(!d.Get("local_auth_enabled").(bool)),
	}

	// Get SKU from Workspace
	workspace, err := mlWorkspacesClient.Get(ctx, *workspaceID)
	if err != nil {
		return err
	}

	identity, err := expandIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	computeClusterParameters := machinelearningcomputes.ComputeResource{
		Properties: computeClusterProperties,
		Identity:   identity,
		Location:   computeClusterProperties.ComputeLocation,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku: &machinelearningcomputes.Sku{
			Name: workspace.Model.Sku.Name,
			Tier: utils.ToPtr(machinelearningcomputes.SkuTier(*workspace.Model.Sku.Tier)),
		},
	}

	future, err := client.ComputeCreateOrUpdate(ctx, id, computeClusterParameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceComputeClusterRead(d, meta)
}

func resourceComputeClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Compute Cluster ID `%q`: %+v", d.Id(), err)
	}

	computeResource, err := client.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(computeResource.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ComputeName)

	workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	// use ComputeResource to get to AKS Cluster ID and other properties
	computeCluster := computeResource.Model.Properties.(machinelearningcomputes.AmlCompute)

	if computeCluster.DisableLocalAuth != nil {
		d.Set("local_auth_enabled", !*computeCluster.DisableLocalAuth)
	}
	d.Set("description", computeCluster.Description)
	if props := computeCluster.Properties; props != nil {
		d.Set("vm_size", props.VMSize)
		d.Set("vm_priority", string(pointer.From(props.VMPriority)))
		d.Set("scale_settings", flattenScaleSettings(props.ScaleSettings))
		d.Set("ssh", flattenUserAccountCredentials(props.UserAccountCredentials))
		enableNodePublicIP := true
		if props.EnableNodePublicIP != nil {
			enableNodePublicIP = *props.EnableNodePublicIP
		}
		d.Set("node_public_ip_enabled", enableNodePublicIP)
		if props.Subnet != nil {
			d.Set("subnet_resource_id", props.Subnet.Id)
		}

		switch *props.RemoteLoginPortPublicAccess {
		case machinelearningcomputes.RemoteLoginPortPublicAccessNotSpecified:
			d.Set("ssh_public_access_enabled", nil)
		case machinelearningcomputes.RemoteLoginPortPublicAccessEnabled:
			d.Set("ssh_public_access_enabled", true)
		case machinelearningcomputes.RemoteLoginPortPublicAccessDisabled:
			d.Set("ssh_public_access_enabled", false)
		}
	}

	if location := computeResource.Model.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity, err := flattenIdentity(computeResource.Model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, computeResource.Model.Tags)
}

func resourceComputeClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.ComputeDelete(ctx, *id, machinelearningcomputes.ComputeDeleteOperationOptions{
		UnderlyingResourceAction: utils.ToPtr(machinelearningcomputes.UnderlyingResourceActionDelete),
	})
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}
	return nil
}

func expandScaleSettings(input []interface{}) *machinelearningcomputes.ScaleSettings {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	maxNodeCount := int64(v["max_node_count"].(int))
	minNodeCount := int64(v["min_node_count"].(int))
	scaleDownNodes := v["scale_down_nodes_after_idle_duration"].(string)

	return &machinelearningcomputes.ScaleSettings{
		MaxNodeCount:                maxNodeCount,
		MinNodeCount:                &minNodeCount,
		NodeIdleTimeBeforeScaleDown: &scaleDownNodes,
	}
}

func expandUserAccountCredentials(input []interface{}) *machinelearningcomputes.UserAccountCredentials {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})

	return &machinelearningcomputes.UserAccountCredentials{
		AdminUserName:         v["admin_username"].(string),
		AdminUserPassword:     utils.String(v["admin_password"].(string)),
		AdminUserSshPublicKey: utils.String(v["key_value"].(string)),
	}
}

func flattenScaleSettings(scaleSettings *machinelearningcomputes.ScaleSettings) []interface{} {
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

func flattenUserAccountCredentials(credentials *machinelearningcomputes.UserAccountCredentials) interface{} {
	if credentials == nil {
		return []interface{}{}
	}
	var username string
	if credentials.AdminUserName != "" {
		username = credentials.AdminUserName
	}

	var adminPassword string
	if credentials.AdminUserPassword != nil {
		adminPassword = *credentials.AdminUserPassword
	}

	var sshPublicKey string
	if credentials.AdminUserSshPublicKey != nil {
		sshPublicKey = *credentials.AdminUserSshPublicKey
	}

	return []interface{}{
		map[string]interface{}{
			"admin_username": username,
			"admin_password": adminPassword,
			"key_value":      sshPublicKey,
		},
	}
}
