package resource

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-09-01/locks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceManagementLock() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceManagementLockCreateUpdate,
		Read:   resourceManagementLockRead,
		Delete: resourceManagementLockDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ParseManagementLockID(id)
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
					string(locks.CanNotDelete),
					string(locks.ReadOnly),
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

func resourceManagementLockCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Management Lock creation.")

	id := parse.NewManagementLockID(d.Get("scope").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.GetByScope(ctx, id.Scope, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_management_lock", id.ID())
		}
	}

	lock := locks.ManagementLockObject{
		ManagementLockProperties: &locks.ManagementLockProperties{
			Level: locks.LockLevel(d.Get("lock_level").(string)),
			Notes: utils.String(d.Get("notes").(string)),
		},
	}

	if _, err := client.CreateOrUpdateByScope(ctx, id.Scope, id.Name, lock); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceManagementLockRead(d, meta)
}

func resourceManagementLockRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseManagementLockID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetByScope(ctx, id.Scope, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("scope", id.Scope)

	if props := resp.ManagementLockProperties; props != nil {
		d.Set("lock_level", string(props.Level))
		d.Set("notes", props.Notes)
	}

	return nil
}

func resourceManagementLockDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseManagementLockID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.DeleteByScope(ctx, id.Scope, id.Name); err != nil {
		// @tombuildsstuff: this is intentionally here in case the parent is gone, since we're under a scope
		// which isn't ideal (as this shouldn't be present for most resources) but should for this one
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
