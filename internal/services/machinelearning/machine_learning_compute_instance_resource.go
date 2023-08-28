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
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceComputeInstance() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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

			"location": commonschema.Location(),

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
}

func resourceComputeInstanceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
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

	computeInstance := &machinelearningcomputes.ComputeInstance{
		Properties: &machinelearningcomputes.ComputeInstanceProperties{
			VMSize:                          utils.String(d.Get("virtual_machine_size").(string)),
			Subnet:                          subnet,
			SshSettings:                     expandComputeSSHSetting(d.Get("ssh").([]interface{})),
			PersonalComputeInstanceSettings: expandComputePersonalComputeInstanceSetting(d.Get("assign_to_user").([]interface{})),
			EnableNodePublicIP:              pointer.To(d.Get("node_public_ip_enabled").(bool)),
		},
		ComputeLocation:  utils.String(d.Get("location").(string)),
		Description:      utils.String(d.Get("description").(string)),
		DisableLocalAuth: utils.Bool(!d.Get("local_auth_enabled").(bool)),
	}
	authType := d.Get("authorization_type").(string)
	if authType != "" {
		computeInstance.Properties.ComputeInstanceAuthorizationType = utils.ToPtr(machinelearningcomputes.ComputeInstanceAuthorizationType(authType))
	}

	parameters := machinelearningcomputes.ComputeResource{
		Properties: computeInstance,
		Identity:   identity,
		Location:   utils.String(location.Normalize(d.Get("location").(string))),
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.ComputeCreateOrUpdate(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating Machine Learning Compute (%q): %+v", id, err)
	}
	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for creation of Machine Learning Compute (%q): %+v", id, err)
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
		return fmt.Errorf("retrieving Machine Learning Compute (%q): %+v", id, err)
	}

	d.Set("name", id.ComputeName)
	workspaceId := workspaces.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	if location := resp.Model.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity, err := flattenIdentity(resp.Model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	props := resp.Model.Properties.(machinelearningcomputes.ComputeInstance)

	if props.DisableLocalAuth != nil {
		d.Set("local_auth_enabled", !*props.DisableLocalAuth)
	}
	d.Set("description", props.Description)
	if props.Properties != nil {
		d.Set("virtual_machine_size", props.Properties.VMSize)
		if props.Properties.Subnet != nil {
			d.Set("subnet_resource_id", props.Properties.Subnet.Id)
		}
		d.Set("authorization_type", string(pointer.From(props.Properties.ComputeInstanceAuthorizationType)))
		d.Set("ssh", flattenComputeSSHSetting(props.Properties.SshSettings))
		d.Set("assign_to_user", flattenComputePersonalComputeInstanceSetting(props.Properties.PersonalComputeInstanceSettings))
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
		UnderlyingResourceAction: utils.ToPtr(machinelearningcomputes.UnderlyingResourceActionDelete),
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
			SshPublicAccess: utils.ToPtr(machinelearningcomputes.SshPublicAccessDisabled),
		}
	}
	value := input[0].(map[string]interface{})
	return &machinelearningcomputes.ComputeInstanceSshSettings{
		SshPublicAccess: utils.ToPtr(machinelearningcomputes.SshPublicAccessEnabled),
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
