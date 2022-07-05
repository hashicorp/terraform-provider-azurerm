package machinelearning

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2021-07-01/machinelearningservices"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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

			"location": azure.SchemaLocation(),

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
					string(machinelearningservices.ComputeInstanceAuthorizationTypePersonal),
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

			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceComputeInstanceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceID, _ := parse.WorkspaceID(d.Get("machine_learning_workspace_id").(string))
	id := parse.NewComputeID(subscriptionId, workspaceID.ResourceGroup, workspaceID.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing Machine Learning Compute (%q): %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_machine_learning_compute_instance", id.ID())
		}
	}

	identity, err := expandIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	var subnet *machinelearningservices.ResourceID
	if subnetId, ok := d.GetOk("subnet_resource_id"); ok {
		subnet = &machinelearningservices.ResourceID{
			ID: utils.String(subnetId.(string)),
		}
	}

	parameters := machinelearningservices.ComputeResource{
		Properties: &machinelearningservices.ComputeInstance{
			Properties: &machinelearningservices.ComputeInstanceProperties{
				VMSize:                           utils.String(d.Get("virtual_machine_size").(string)),
				Subnet:                           subnet,
				SSHSettings:                      expandComputeSSHSetting(d.Get("ssh").([]interface{})),
				ComputeInstanceAuthorizationType: machinelearningservices.ComputeInstanceAuthorizationType(d.Get("authorization_type").(string)),
				PersonalComputeInstanceSettings:  expandComputePersonalComputeInstanceSetting(d.Get("assign_to_user").([]interface{})),
			},
			ComputeLocation:  utils.String(d.Get("location").(string)),
			Description:      utils.String(d.Get("description").(string)),
			DisableLocalAuth: utils.Bool(!d.Get("local_auth_enabled").(bool)),
		},
		Identity: identity,
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating Machine Learning Compute (%q): %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Machine Learning Compute (%q): %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceComputeInstanceRead(d, meta)
}

func resourceComputeInstanceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ComputeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Machine Learning Compute %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Machine Learning Compute (%q): %+v", id, err)
	}

	d.Set("name", id.Name)
	workspaceId := parse.NewWorkspaceID(subscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity, err := flattenIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props, ok := resp.Properties.AsComputeInstance(); ok && props != nil {
		if props.DisableLocalAuth != nil {
			d.Set("local_auth_enabled", !*props.DisableLocalAuth)
		}
		d.Set("description", props.Description)
		if props.Properties != nil {
			d.Set("virtual_machine_size", props.Properties.VMSize)
			if props.Properties.Subnet != nil {
				d.Set("subnet_resource_id", props.Properties.Subnet.ID)
			}
			d.Set("authorization_type", props.Properties.ComputeInstanceAuthorizationType)
			d.Set("ssh", flattenComputeSSHSetting(props.Properties.SSHSettings))
			d.Set("assign_to_user", flattenComputePersonalComputeInstanceSetting(props.Properties.PersonalComputeInstanceSettings))
		}
	} else {
		return fmt.Errorf("compute resource %s is not a ComputeInstance Compute", id)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceComputeInstanceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.ComputeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.ComputeID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, machinelearningservices.UnderlyingResourceActionDelete)
	if err != nil {
		return fmt.Errorf("deleting Machine Learning Compute (%q): %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the Machine Learning Compute (%q): %+v", id, err)
	}
	return nil
}

func expandComputePersonalComputeInstanceSetting(input []interface{}) *machinelearningservices.PersonalComputeInstanceSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	value := input[0].(map[string]interface{})
	return &machinelearningservices.PersonalComputeInstanceSettings{
		AssignedUser: &machinelearningservices.AssignedUser{
			ObjectID: utils.String(value["object_id"].(string)),
			TenantID: utils.String(value["tenant_id"].(string)),
		},
	}
}

func expandComputeSSHSetting(input []interface{}) *machinelearningservices.ComputeInstanceSSHSettings {
	if len(input) == 0 {
		return &machinelearningservices.ComputeInstanceSSHSettings{
			SSHPublicAccess: machinelearningservices.SSHPublicAccessDisabled,
		}
	}
	value := input[0].(map[string]interface{})
	return &machinelearningservices.ComputeInstanceSSHSettings{
		SSHPublicAccess: machinelearningservices.SSHPublicAccessEnabled,
		AdminPublicKey:  utils.String(value["public_key"].(string)),
	}
}

func flattenComputePersonalComputeInstanceSetting(settings *machinelearningservices.PersonalComputeInstanceSettings) interface{} {
	if settings == nil || settings.AssignedUser == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"tenant_id": settings.AssignedUser.TenantID,
			"object_id": settings.AssignedUser.ObjectID,
		},
	}
}

func flattenComputeSSHSetting(settings *machinelearningservices.ComputeInstanceSSHSettings) interface{} {
	if settings == nil || strings.EqualFold(string(settings.SSHPublicAccess), string(machinelearningservices.SSHPublicAccessDisabled)) {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"public_key": settings.AdminPublicKey,
			"username":   settings.AdminUserName,
			"port":       settings.SSHPort,
		},
	}
}
