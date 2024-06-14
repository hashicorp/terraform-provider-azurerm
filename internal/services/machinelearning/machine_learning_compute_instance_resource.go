// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceComputeInstance() *pluginsdk.Resource {
	resource := pluginsdk.Resource{
		Create: resourceComputeInstanceCreate,
		Read:   resourceComputeInstanceRead,
		Delete: resourceComputeInstanceDelete,

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
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]{3,24}$`),
					"It can include letters, digits and dashes. It must start with a letter, end with a letter or digit, and be between 3 and 24 characters in length."),
			},

			"machine_learning_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"virtual_machine_size": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"authorization_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(machinelearningcomputes.ComputeInstanceAuthorizationTypePersonal),
				}, false),
			},

			"assign_to_user": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"tenant_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},

						"object_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptionalForceNew(),

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
						"public_key": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"username": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"port": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"subnet_resource_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"node_public_ip_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"tags": commonschema.TagsForceNew(),
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["location"] = &pluginsdk.Schema{
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			Deprecated:       "The `azurerm_machine_learning_compute_instance` must be deployed to the same location as the associated `azurerm_machine_learning_workspace` resource, as the `location` fields must be the same the `location` field no longer has any effect and will be removed in version 4.0 of the AzureRM Provider",
			ValidateFunc:     location.EnhancedValidate,
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
		}
	}

	return &resource
}

func resourceComputeInstanceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	mlWorkspacesClient := meta.(*clients.Client).MachineLearning.Workspaces
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceID, _ := workspaces.ParseWorkspaceID(d.Get("machine_learning_workspace_id").(string))
	id := machinelearningcomputes.NewComputeID(subscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.ComputeGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing Machine Learning Compute (%q): %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_machine_learning_compute_instance", id.ID())
		}
	}

	identity, err := expandIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	var subnet *machinelearningcomputes.ResourceId
	if subnetId, ok := d.GetOk("subnet_resource_id"); ok {
		subnet = &machinelearningcomputes.ResourceId{
			Id: subnetId.(string),
		}
	}

	if !d.Get("node_public_ip_enabled").(bool) && d.Get("subnet_resource_id").(string) == "" {
		return fmt.Errorf("`subnet_resource_id` must be set if `node_public_ip_enabled` is set to `false`")
	}

	// NOTE: The 'ComputeResource' struct contains the information
	// which is related to the parent resource of the instance that is
	// to be deployed (e.g., the workspace), which is why we need to
	// GET the workspace to discover the location it has been deployed to.
	// If we do not set the correct location, the identity will be created
	// and then orphaned in the incorrect region.
	workspace, err := mlWorkspacesClient.Get(ctx, *workspaceID)
	if err != nil {
		return err
	}

	model := workspace.Model
	if model == nil {
		return fmt.Errorf("machine learning %s Workspace: model is nil", id)
	}

	if model.Location == nil {
		return fmt.Errorf("machine learning %s Workspace: model `Location` is nil", id)
	}

	parameters := machinelearningcomputes.ComputeResource{
		Identity: identity,
		Location: pointer.To(azure.NormalizeLocation(*model.Location)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	// NOTE: In 4.0 the 'location' field will be deprecated...
	props := machinelearningcomputes.ComputeInstance{
		Properties: &machinelearningcomputes.ComputeInstanceProperties{
			VMSize:                          utils.String(d.Get("virtual_machine_size").(string)),
			Subnet:                          subnet,
			SshSettings:                     expandComputeSSHSetting(d.Get("ssh").([]interface{})),
			PersonalComputeInstanceSettings: expandComputePersonalComputeInstanceSetting(d.Get("assign_to_user").([]interface{})),
			EnableNodePublicIP:              pointer.To(d.Get("node_public_ip_enabled").(bool)),
		},
		Description:      utils.String(d.Get("description").(string)),
		DisableLocalAuth: utils.Bool(!d.Get("local_auth_enabled").(bool)),
	}

	// NOTE: The 'location' field is not supported for instances, "Compute clusters can be created in
	// a different region than your workspace. This functionality is only available for compute
	// clusters, not compute instances"
	//
	// https://learn.microsoft.com/azure/machine-learning/how-to-create-attach-compute-cluster?view=azureml-api-2&tabs=python#limitations

	if v, ok := d.GetOk("authorization_type"); ok {
		props.Properties.ComputeInstanceAuthorizationType = pointer.To(machinelearningcomputes.ComputeInstanceAuthorizationType(v.(string)))
	}

	parameters.Properties = props

	future, err := client.ComputeCreateOrUpdate(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating Machine Learning %s: %+v", id, err)
	}
	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for creation of Machine Learning %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceComputeInstanceRead(d, meta)
}

func resourceComputeInstanceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Machine Learning Compute %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Machine Learning %s: %+v", id, err)
	}

	workspaceId := workspaces.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)

	model := resp.Model
	if model == nil {
		return fmt.Errorf("machine learning %s: model is nil", id)
	}

	props := model.Properties.(machinelearningcomputes.ComputeInstance)

	d.Set("name", id.ComputeName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	identity, err := flattenIdentity(model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}

	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props.DisableLocalAuth != nil {
		d.Set("local_auth_enabled", !*props.DisableLocalAuth)
	}

	d.Set("description", props.Description)

	if props.Properties != nil {
		d.Set("virtual_machine_size", props.Properties.VMSize)
		d.Set("authorization_type", string(pointer.From(props.Properties.ComputeInstanceAuthorizationType)))
		d.Set("ssh", flattenComputeSSHSetting(props.Properties.SshSettings))
		d.Set("assign_to_user", flattenComputePersonalComputeInstanceSetting(props.Properties.PersonalComputeInstanceSettings))

		if props.Properties.Subnet != nil {
			d.Set("subnet_resource_id", props.Properties.Subnet.Id)
		}

		enableNodePublicIP := true
		if props.Properties.ConnectivityEndpoints.PublicIPAddress == nil {
			enableNodePublicIP = false
		}

		d.Set("node_public_ip_enabled", enableNodePublicIP)
	}

	return tags.FlattenAndSet(d, resp.Model.Tags)
}

func resourceComputeInstanceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.ComputeDelete(ctx, *id, machinelearningcomputes.ComputeDeleteOperationOptions{
		UnderlyingResourceAction: pointer.To(machinelearningcomputes.UnderlyingResourceActionDelete),
	})
	if err != nil {
		return fmt.Errorf("deleting Machine Learning Compute (%q): %+v", id, err)
	}

	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of the Machine Learning Compute (%q): %+v", id, err)
	}
	return nil
}

func expandComputePersonalComputeInstanceSetting(input []interface{}) *machinelearningcomputes.PersonalComputeInstanceSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	value := input[0].(map[string]interface{})
	return &machinelearningcomputes.PersonalComputeInstanceSettings{
		AssignedUser: &machinelearningcomputes.AssignedUser{
			ObjectId: value["object_id"].(string),
			TenantId: value["tenant_id"].(string),
		},
	}
}

func expandComputeSSHSetting(input []interface{}) *machinelearningcomputes.ComputeInstanceSshSettings {
	if len(input) == 0 {
		return &machinelearningcomputes.ComputeInstanceSshSettings{
			SshPublicAccess: pointer.To(machinelearningcomputes.SshPublicAccessDisabled),
		}
	}
	value := input[0].(map[string]interface{})
	return &machinelearningcomputes.ComputeInstanceSshSettings{
		SshPublicAccess: pointer.To(machinelearningcomputes.SshPublicAccessEnabled),
		AdminPublicKey:  utils.String(value["public_key"].(string)),
	}
}

func flattenComputePersonalComputeInstanceSetting(settings *machinelearningcomputes.PersonalComputeInstanceSettings) interface{} {
	if settings == nil || settings.AssignedUser == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"tenant_id": settings.AssignedUser.TenantId,
			"object_id": settings.AssignedUser.ObjectId,
		},
	}
}

func flattenComputeSSHSetting(settings *machinelearningcomputes.ComputeInstanceSshSettings) interface{} {
	if settings == nil || strings.EqualFold(string(*settings.SshPublicAccess), string(machinelearningcomputes.SshPublicAccessDisabled)) {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"public_key": settings.AdminPublicKey,
			"username":   settings.AdminUserName,
			"port":       settings.SshPort,
		},
	}
}
