package mssql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMSSQLManagedInstanceAdmin() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMSSQLManagedInstanceAdminCreateUpdate,
		Read:   resourceArmMSSQLManagedInstanceAdminRead,
		Update: resourceArmMSSQLManagedInstanceAdminCreateUpdate,
		Delete: resourceArmMSSQLManagedInstanceAdminDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"managed_instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: azure.ValidateResourceID,
			},

			"login_username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"object_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.IsUUID,
			},

			"tenant_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validation.IsUUID,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmMSSQLManagedInstanceAdminCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	adminClient := meta.(*clients.Client).MSSQL.ManagedInstanceAdministratorsClient
	managedInstanceClient := meta.(*clients.Client).MSSQL.ManagedInstancesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	managedInstanceId := d.Get("managed_instance_id").(string)

	id, err := azure.ParseAzureResourceID(managedInstanceId)
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["managedInstances"]

	if _, err := managedInstanceClient.Get(ctx, resGroup, name); err != nil {
		return fmt.Errorf("while reading managed SQL instance %s: %v", name, err)
	}

	if d.IsNewResource() {
		existing, err := adminClient.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("while checking for presence of existing managed sql instance aad admin details %q (Resource Group %q): %+v", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mssql_managed_instance_admin", *existing.ID)
		}
	}

	sid, _ := uuid.FromString(d.Get("object_id").(string))
	managedInstanceAdmin := sql.ManagedInstanceAdministrator{
		ManagedInstanceAdministratorProperties: &sql.ManagedInstanceAdministratorProperties{
			AdministratorType: utils.String("ActiveDirectory"),
			Sid:               &sid,
		},
	}

	if v, exists := d.GetOk("login_username"); exists {
		managedInstanceAdmin.ManagedInstanceAdministratorProperties.Login = utils.String(v.(string))
	}

	if v, exists := d.GetOk("tenant_id"); exists {
		tid, _ := uuid.FromString(v.(string))
		managedInstanceAdmin.ManagedInstanceAdministratorProperties.TenantID = &tid
	}

	adminFuture, err := adminClient.CreateOrUpdate(ctx, resGroup, name, managedInstanceAdmin)
	if err != nil {
		return fmt.Errorf("while creating Managed SQL Instance %q AAD admin details (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = adminFuture.WaitForCompletionRef(ctx, adminClient.Client); err != nil {
		return fmt.Errorf("while waiting for creation of Managed SQL Instance %q AAD admin details (Resource Group %q): %+v", name, resGroup, err)
	}

	result, err := adminClient.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("while making get request for managed SQL instance AAD Admin details %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if result.ID == nil {
		return fmt.Errorf("while getting ID from managed SQL instance %q AAD Admin details (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*result.ID)

	return resourceArmMSSQLManagedInstanceAdminRead(d, meta)

}

func resourceArmMSSQLManagedInstanceAdminRead(d *schema.ResourceData, meta interface{}) error {
	adminClient := meta.(*clients.Client).MSSQL.ManagedInstanceAdministratorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["managedInstances"]

	adminResp, err := adminClient.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("while reading managed instance %s AAD admin: %v", name, err)
	}

	managedInstanceId, _ := azure.GetSQLResourceParentId(d.Id())
	if err != nil {
		return err
	}
	d.Set("managed_instance_id", managedInstanceId)
	d.Set("name", adminResp.Name)
	d.Set("type", adminResp.Type)

	if props := adminResp.ManagedInstanceAdministratorProperties; props != nil {
		d.Set("admin_type", props.AdministratorType)
		d.Set("login_username", props.Login)
		d.Set("object_id", props.Sid.String())
		d.Set("tenant_id", props.TenantID.String())
	}
	return nil
}

func resourceArmMSSQLManagedInstanceAdminDelete(d *schema.ResourceData, meta interface{}) error {
	adminClient := meta.(*clients.Client).MSSQL.ManagedInstanceAdministratorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["managedInstances"]

	future, err := adminClient.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("while deleting managed SQL instance %s admin details: %+v", name, err)
	}

	return future.WaitForCompletionRef(ctx, adminClient.Client)
}
