package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-09-01/locks"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmManagementLock() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmManagementLockCreateUpdate,
		Read:   resourceArmManagementLockRead,
		Delete: resourceArmManagementLockDelete,
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmManagementLockName,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"lock_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(locks.CanNotDelete),
					string(locks.ReadOnly),
				}, false),
			},

			"notes": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 512),
			},
		},
	}
}

func resourceArmManagementLockCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Management Lock creation.")

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetByScope(ctx, scope, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Management Lock %q (Scope %q): %s", name, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_management_lock", *existing.ID)
		}
	}

	lockLevel := d.Get("lock_level").(string)
	notes := d.Get("notes").(string)

	lock := locks.ManagementLockObject{
		ManagementLockProperties: &locks.ManagementLockProperties{
			Level: locks.LockLevel(lockLevel),
			Notes: utils.String(notes),
		},
	}

	if _, err := client.CreateOrUpdateByScope(ctx, scope, name, lock); err != nil {
		return err
	}

	read, err := client.GetByScope(ctx, scope, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID of AzureRM Management Lock %q (Scope %q)", name, scope)
	}

	d.SetId(*read.ID)
	return resourceArmManagementLockRead(d, meta)
}

func resourceArmManagementLockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureRMLockId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetByScope(ctx, id.Scope, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Management Lock %q (Scope %q): %+v", id.Name, id.Scope, err)
	}

	d.Set("name", resp.Name)
	d.Set("scope", id.Scope)

	if props := resp.ManagementLockProperties; props != nil {
		d.Set("lock_level", string(props.Level))
		d.Set("notes", props.Notes)
	}

	return nil
}

func resourceArmManagementLockDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureRMLockId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DeleteByScope(ctx, id.Scope, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Management Lock %q (Scope %q): %+v", id.Name, id.Scope, err)
	}

	return nil
}

type AzureManagementLockId struct {
	Scope string
	Name  string
}

func parseAzureRMLockId(id string) (*AzureManagementLockId, error) {
	segments := strings.Split(id, "/providers/Microsoft.Authorization/locks/")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected ID to be in the format `{scope}/providers/Microsoft.Authorization/locks/{name} - got %d segments", len(segments))
	}

	scope := segments[0]
	name := segments[1]
	lockId := AzureManagementLockId{
		Scope: scope,
		Name:  name,
	}
	return &lockId, nil
}

func validateArmManagementLockName(v interface{}, k string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile(`[A-Za-z0-9-_]`).MatchString(input) {
		errors = append(errors, fmt.Errorf("%s can only consist of alphanumeric characters, dashes and underscores", k))
	}

	if len(input) >= 260 {
		errors = append(errors, fmt.Errorf("%s can only be a maximum of 260 characters", k))
	}

	return warnings, errors
}
