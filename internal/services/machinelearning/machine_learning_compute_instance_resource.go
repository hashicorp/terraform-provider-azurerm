package machinelearning

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
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
			_, err := parse.ComputeID(id)
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
				ValidateFunc: validate.WorkspaceID,
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
				ValidateFunc: networkValidate.SubnetID,
			},

			"tags": commonschema.TagsForceNew(),
		},
	}
}

func resourceComputeInstanceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
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

	identity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	var subnet *machinelearningcomputes.ResourceId
	if subnetId, ok := d.GetOk("subnet_resource_id"); ok {
		subnet = &machinelearningcomputes.ResourceId{
			Id: subnetId.(string),
		}
	}

	computeInstance := machinelearningcomputes.ComputeInstance{
		Properties: &machinelearningcomputes.ComputeInstanceProperties{
			VMSize:                          utils.String(d.Get("virtual_machine_size").(string)),
			Subnet:                          subnet,
			SshSettings:                     expandComputeSSHSetting(d.Get("ssh").([]interface{})),
			PersonalComputeInstanceSettings: expandComputePersonalComputeInstanceSetting(d.Get("assign_to_user").([]interface{})),
		},
		ComputeLocation:  utils.String(d.Get("location").(string)),
		Description:      utils.String(d.Get("description").(string)),
		DisableLocalAuth: utils.Bool(!d.Get("local_auth_enabled").(bool)),
	}
	authType := d.Get("authorization_type").(string)
	if authType != "" {
		computeInstance.Properties.ComputeInstanceAuthorizationType = utils.ToPtr(machinelearningcomputes.ComputeInstanceAuthorizationType(authType))
	}

	var compute machinelearningcomputes.Compute = computeInstance
	parameters := machinelearningcomputes.ComputeResource{
		Properties: pointer.To(compute),
		Identity:   identity,
		Location:   utils.String(location.Normalize(d.Get("location").(string))),
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.ComputeCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceComputeInstanceRead(d, meta)
}

func resourceComputeInstanceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
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

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		identity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if computeInstance, ok := (*props).(machinelearningcomputes.ComputeInstance); ok {
				localAuthEnabled := true
				if computeInstance.DisableLocalAuth != nil {
					localAuthEnabled = !*computeInstance.DisableLocalAuth
				}
				d.Set("local_auth_enabled", localAuthEnabled)

				d.Set("description", computeInstance.Description)
				if instanceProps := computeInstance.Properties; instanceProps != nil {
					d.Set("virtual_machine_size", instanceProps.VMSize)
					subnetId := ""
					if computeInstance.Properties.Subnet != nil {
						subnetId = instanceProps.Subnet.Id
					}
					d.Set("subnet_resource_id", subnetId)
					d.Set("authorization_type", string(pointer.From(instanceProps.ComputeInstanceAuthorizationType)))
					d.Set("ssh", flattenComputeSSHSetting(instanceProps.SshSettings))
					d.Set("assign_to_user", flattenComputePersonalComputeInstanceSetting(instanceProps.PersonalComputeInstanceSettings))
				}
			}
		}
	}

	return nil
}

func resourceComputeInstanceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	opts := machinelearningcomputes.ComputeDeleteOperationOptions{
		UnderlyingResourceAction: utils.ToPtr(machinelearningcomputes.UnderlyingResourceActionDelete),
	}
	if err := client.ComputeDeleteThenPoll(ctx, *id, opts); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
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
