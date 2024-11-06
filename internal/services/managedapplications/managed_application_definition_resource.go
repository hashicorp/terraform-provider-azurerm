// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedapplications

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applicationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedapplications/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedApplicationDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagedApplicationDefinitionCreate,
		Read:   resourceManagedApplicationDefinitionRead,
		Update: resourceManagedApplicationDefinitionUpdate,
		Delete: resourceManagedApplicationDefinitionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := applicationdefinitions.ParseApplicationDefinitionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApplicationDefinitionName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ApplicationDefinitionDisplayName,
			},

			"lock_level": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(applicationdefinitions.ApplicationLockLevelCanNotDelete),
					string(applicationdefinitions.ApplicationLockLevelNone),
					string(applicationdefinitions.ApplicationLockLevelReadOnly),
				}, false),
			},

			"authorization": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"role_definition_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"service_principal_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"create_ui_definition": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				ConflictsWith:    []string{"package_file_uri"},
				RequiredWith:     []string{"main_template"},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ApplicationDefinitionDescription,
			},

			"main_template": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				ConflictsWith:    []string{"package_file_uri"},
				RequiredWith:     []string{"create_ui_definition"},
			},

			"package_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"package_file_uri": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceManagedApplicationDefinitionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := applicationdefinitions.NewApplicationDefinitionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("failed to check for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_managed_application_definition", id.ID())
	}

	parameters := applicationdefinitions.ApplicationDefinition{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: applicationdefinitions.ApplicationDefinitionProperties{
			Authorizations: expandManagedApplicationDefinitionAuthorization(d.Get("authorization").(*pluginsdk.Set).List()),
			Description:    pointer.To(d.Get("description").(string)),
			DisplayName:    pointer.To(d.Get("display_name").(string)),
			IsEnabled:      pointer.To(d.Get("package_enabled").(bool)),
			LockLevel:      applicationdefinitions.ApplicationLockLevel(d.Get("lock_level").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("create_ui_definition"); ok {
		parameters.Properties.CreateUiDefinition = &v
	}

	if v, ok := d.GetOk("main_template"); ok {
		parameters.Properties.MainTemplate = &v
	}

	if v, ok := d.GetOk("package_file_uri"); ok {
		parameters.Properties.PackageFileUri = pointer.To(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceManagedApplicationDefinitionRead(d, meta)
}

func resourceManagedApplicationDefinitionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationdefinitions.ParseApplicationDefinitionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := existing.Model

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("display_name") {
		payload.Properties.DisplayName = pointer.To(d.Get("display_name").(string))
	}

	if d.HasChange("package_enabled") {
		payload.Properties.IsEnabled = pointer.To(d.Get("package_enabled").(bool))
	}

	if d.HasChange("authorization") {
		payload.Properties.Authorizations = expandManagedApplicationDefinitionAuthorization(d.Get("authorization").(*pluginsdk.Set).List())
	}

	if d.HasChange("create_ui_definition") {
		// handle API error: The 'MainTemplate, CreateUiDefinition' properties should be empty if package zip file uri is provided.
		if v, ok := d.GetOk("create_ui_definition"); ok {
			payload.Properties.CreateUiDefinition = pointer.To(v)
		} else {
			payload.Properties.CreateUiDefinition = nil
		}
	}

	if d.HasChange("main_template") {
		// handle API error: The 'MainTemplate, CreateUiDefinition' properties should be empty if package zip file uri is provided.
		if v, ok := d.GetOk("main_template"); ok {
			payload.Properties.MainTemplate = pointer.To(v)
		} else {
			payload.Properties.MainTemplate = nil
		}
	}

	if d.HasChange("package_file_uri") {
		payload.Properties.PackageFileUri = pointer.To(d.Get("package_file_uri").(string))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	// update payload only supports tags, so we'll continue using CreateOrUpdate method here
	if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceManagedApplicationDefinitionRead(d, meta)
}

func resourceManagedApplicationDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationdefinitions.ParseApplicationDefinitionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ApplicationDefinitionName)
	d.Set("resource_group_name", id.ResourceGroupName) // missing from response?

	if model := resp.Model; model != nil {
		p := model.Properties

		d.Set("location", location.NormalizeNilable(model.Location))

		if err := d.Set("authorization", flattenManagedApplicationDefinitionAuthorization(p.Authorizations)); err != nil {
			return fmt.Errorf("setting `authorization`: %+v", err)
		}
		d.Set("description", p.Description)
		d.Set("display_name", p.DisplayName)
		d.Set("package_enabled", p.IsEnabled)
		d.Set("lock_level", string(p.LockLevel))

		// the following are not returned from the API so lets pull it from state
		if v, ok := d.GetOk("create_ui_definition"); ok {
			d.Set("create_ui_definition", v.(string))
		}
		if v, ok := d.GetOk("main_template"); ok {
			d.Set("main_template", v.(string))
		}
		if v, ok := d.GetOk("package_file_uri"); ok {
			d.Set("package_file_uri", v.(string))
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceManagedApplicationDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationdefinitions.ParseApplicationDefinitionID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("issuing AzureRM delete request for '%s': %+v", *id, err)
	}

	return nil
}

func expandManagedApplicationDefinitionAuthorization(input []interface{}) *[]applicationdefinitions.ApplicationAuthorization {
	results := make([]applicationdefinitions.ApplicationAuthorization, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		result := applicationdefinitions.ApplicationAuthorization{
			RoleDefinitionId: v["role_definition_id"].(string),
			PrincipalId:      v["service_principal_id"].(string),
		}

		results = append(results, result)
	}
	return &results
}

func flattenManagedApplicationDefinitionAuthorization(input *[]applicationdefinitions.ApplicationAuthorization) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"role_definition_id":   item.RoleDefinitionId,
			"service_principal_id": item.PrincipalId,
		})
	}

	return results
}
