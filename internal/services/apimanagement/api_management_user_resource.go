package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementUser() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementUserCreateUpdate,
		Read:   resourceApiManagementUserRead,
		Update: resourceApiManagementUserCreateUpdate,
		Delete: resourceApiManagementUserDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.UserID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(45 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(45 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"user_id": schemaz.SchemaApiManagementUserName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"first_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"email": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"last_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"confirmation": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.ConfirmationInvite),
					string(apimanagement.ConfirmationSignup),
				}, false),
			},

			"note": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"password": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.UserStateActive),
					string(apimanagement.UserStateBlocked),
					string(apimanagement.UserStatePending),
				}, false),
			},
		},
	}
}

func resourceApiManagementUserCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.UsersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for API Management User creation.")

	id := parse.NewUserID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("user_id").(string))

	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)
	state := d.Get("state").(string)
	note := d.Get("note").(string)
	password := d.Get("password").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_user", id.ID())
		}
	}

	properties := apimanagement.UserCreateParameters{
		UserCreateParameterProperties: &apimanagement.UserCreateParameterProperties{
			FirstName: utils.String(firstName),
			LastName:  utils.String(lastName),
			Email:     utils.String(email),
		},
	}

	confirmation := d.Get("confirmation").(string)
	if confirmation != "" {
		properties.UserCreateParameterProperties.Confirmation = apimanagement.Confirmation(confirmation)
	}
	if note != "" {
		properties.UserCreateParameterProperties.Note = utils.String(note)
	}
	if password != "" {
		properties.UserCreateParameterProperties.Password = utils.String(password)
	}
	if state != "" {
		properties.UserCreateParameterProperties.State = apimanagement.UserState(state)
	}

	notify := utils.Bool(false)
	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.Name, properties, notify, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementUserRead(d, meta)
}

func resourceApiManagementUserRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.UsersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.UserID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("user_id", id.Name)
	d.Set("api_management_name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.UserContractProperties; props != nil {
		d.Set("first_name", props.FirstName)
		d.Set("last_name", props.LastName)
		d.Set("email", props.Email)
		d.Set("note", props.Note)
		d.Set("state", string(props.State))
	}

	return nil
}

func resourceApiManagementUserDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.UsersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.UserID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", *id)
	deleteSubscriptions := utils.Bool(true)
	notify := utils.Bool(false)
	resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.Name, "", deleteSubscriptions, notify, apimanagement.AppTypeDeveloperPortal)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
