// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_user

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/user"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementUserCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.UsersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := user.NewUserID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("user_id").(string))

	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)
	state := d.Get("state").(string)
	note := d.Get("note").(string)
	password := d.Get("password").(string)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
		properties.Properties.Confirmation = pointer.ToEnum[user.Confirmation](confirmation)
	}
	if note != "" {
		properties.Properties.Note = pointer.To(note)
	}
	if password != "" {
		properties.Properties.Password = pointer.To(password)
	}
	if state != "" {
		properties.Properties.State = pointer.ToEnum[user.UserState](state)
	}

	if _, err := client.CreateOrUpdate(ctx, id, properties, user.CreateOrUpdateOperationOptions{Notify: pointer.To(false)}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementUserRead(d, meta)
}
