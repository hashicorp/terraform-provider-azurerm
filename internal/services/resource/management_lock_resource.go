// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/managementlocks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceManagementLock() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagementLockCreate,
		Read:   resourceManagementLockRead,
		Delete: resourceManagementLockDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := managementlocks.ParseScopedLockID(id)
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
				ValidateFunc: validate.ManagementLockName,
			},

			"scope": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"lock_level": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(managementlocks.LockLevelCanNotDelete),
					string(managementlocks.LockLevelReadOnly),
				}, false),
			},

			"notes": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 512),
			},
		},
	}
}

func resourceManagementLockCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := managementlocks.NewScopedLockID(d.Get("scope").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.GetByScope(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_management_lock", id.ID())
		}
	}

	payload := managementlocks.ManagementLockObject{
		Properties: managementlocks.ManagementLockProperties{
			Level: managementlocks.LockLevel(d.Get("lock_level").(string)),
			Notes: utils.String(d.Get("notes").(string)),
		},
	}

	if _, err := client.CreateOrUpdateByScope(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceManagementLockRead(d, meta)
}

func resourceManagementLockRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managementlocks.ParseScopedLockID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetByScope(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.LockName)
	d.Set("scope", id.Scope)

	if model := resp.Model; model != nil {
		d.Set("lock_level", string(model.Properties.Level))
		d.Set("notes", model.Properties.Notes)
	}

	return nil
}

func resourceManagementLockDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managementlocks.ParseScopedLockID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.DeleteByScope(ctx, *id); err != nil {
		// @tombuildsstuff: this is intentionally here in case the parent is gone, since we're under a scope
		// which isn't ideal (as this logic shouldn't be present for most resources) but should for this one
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
