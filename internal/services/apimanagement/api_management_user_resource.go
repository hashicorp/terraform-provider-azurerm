// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/user"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementUser() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementUserCreateUpdate,
		Read:   resourceApiManagementUserRead,
		Update: resourceApiManagementUserCreateUpdate,
		Delete: resourceApiManagementUserDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := user.ParseUserID(id)
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

			"resource_group_name": commonschema.ResourceGroupName(),

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
					string(user.ConfirmationInvite),
					string(user.ConfirmationSignup),
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
					string(user.UserStateActive),
					string(user.UserStateBlocked),
					string(user.UserStatePending),
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

	id := user.NewUserID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("user_id").(string))

	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)
	state := d.Get("state").(string)
	note := d.Get("note").(string)
	password := d.Get("password").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_user", id.ID())
		}
	}

	properties := user.UserCreateParameters{
		Properties: &user.UserCreateParameterProperties{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		},
	}

	confirmation := d.Get("confirmation").(string)
	if confirmation != "" {
		properties.Properties.Confirmation = pointer.To(user.Confirmation(confirmation))
	}
	if note != "" {
		properties.Properties.Note = pointer.To(note)
	}
	if password != "" {
		properties.Properties.Password = pointer.To(password)
	}
	if state != "" {
		properties.Properties.State = pointer.To(user.UserState(state))
	}

	if _, err := client.CreateOrUpdate(ctx, id, properties, user.CreateOrUpdateOperationOptions{Notify: pointer.To(false)}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementUserRead(d, meta)
}

func resourceApiManagementUserRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.UsersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := user.ParseUserID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("user_id", id.UserId)
	d.Set("api_management_name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("first_name", pointer.From(props.FirstName))
			d.Set("last_name", pointer.From(props.LastName))
			d.Set("email", pointer.From(props.Email))
			d.Set("note", pointer.From(props.Note))
			d.Set("state", string(pointer.From(props.State)))
		}
	}

	return nil
}

func resourceApiManagementUserDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.UsersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := user.ParseUserID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", *id)
	resp, err := client.Delete(ctx, *id, user.DeleteOperationOptions{AppType: pointer.To(user.AppTypeDeveloperPortal), DeleteSubscriptions: pointer.To(true), Notify: pointer.To(false)})
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
