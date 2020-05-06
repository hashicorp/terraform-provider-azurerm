package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementUserCreateUpdate,
		Read:   resourceArmApiManagementUserRead,
		Update: resourceArmApiManagementUserCreateUpdate,
		Delete: resourceArmApiManagementUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"user_id": azure.SchemaApiManagementUserName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"first_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"last_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"confirmation": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.Invite),
					string(apimanagement.Signup),
				}, false),
			},

			"note": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"state": {
				Type:     schema.TypeString,
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

func resourceArmApiManagementUserCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.UsersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for API Management User creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	userId := d.Get("user_id").(string)

	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)
	state := d.Get("state").(string)
	note := d.Get("note").(string)
	password := d.Get("password").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, userId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing User %q (API Management Service %q / Resource Group %q): %s", userId, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_user", *existing.ID)
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

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, userId, properties, ""); err != nil {
		return fmt.Errorf("creating/updating User %q (API Management Service %q / Resource Group %q): %+v", userId, serviceName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, serviceName, userId)
	if err != nil {
		return fmt.Errorf("retrieving User %q (API Management Service %q / Resource Group %q): %+v", userId, serviceName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for User %q (API Management Service %q / Resource Group %q)", userId, serviceName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmApiManagementUserRead(d, meta)
}

func resourceArmApiManagementUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.UsersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	userId := id.Path["users"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, userId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("User %q was not found in API Management Service %q / Resource Group %q - removing from state!", userId, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on User %q (API Management Service %q / Resource Group %q): %+v", userId, serviceName, resourceGroup, err)
	}

	d.Set("user_id", userId)
	d.Set("api_management_name", serviceName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.UserContractProperties; props != nil {
		d.Set("first_name", props.FirstName)
		d.Set("last_name", props.LastName)
		d.Set("email", props.Email)
		d.Set("note", props.Note)
		d.Set("state", string(props.State))
	}

	return nil
}

func resourceArmApiManagementUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.UsersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	userId := id.Path["users"]

	log.Printf("[DEBUG] Deleting User %q (API Management Service %q / Resource Grouo %q)", userId, serviceName, resourceGroup)
	deleteSubscriptions := utils.Bool(true)
	notify := utils.Bool(false)
	resp, err := client.Delete(ctx, resourceGroup, serviceName, userId, "", deleteSubscriptions, notify)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting User %q (API Management Service %q / Resource Group %q): %+v", userId, serviceName, resourceGroup, err)
		}
	}

	return nil
}
