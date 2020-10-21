package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMsSqlDatabaseExtendedAuditingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMsSqlDatabaseExtendedAuditingPolicyCreateUpdate,
		Read:   resourceArmMsSqlDatabaseExtendedAuditingPolicyRead,
		Update: resourceArmMsSqlDatabaseExtendedAuditingPolicyCreateUpdate,
		Delete: resourceArmMsSqlDatabaseExtendedAuditingPolicyDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.MssqlDatabaseExtendedAuditingPolicyID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"database_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MsSqlDatabaseID,
			},

			"storage_endpoint": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPS,
			},

			"storage_account_access_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_account_access_key_is_secondary": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"retention_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 3285),
			},
		},
	}
}

func resourceArmMsSqlDatabaseExtendedAuditingPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabaseExtendedBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MsSql Database Extended Auditing Policy creation.")

	dbId, err := parse.MsSqlDatabaseID(d.Get("database_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, dbId.ResourceGroup, dbId.MsSqlServer, dbId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Failed to check for presence of existing Database %q Sql Auditing (MsSql Server %q / Resource Group %q): %s", dbId.Name, dbId.MsSqlServer, dbId.ResourceGroup, err)
			}
		}

		// if state is not disabled, we should import it.
		if existing.ID != nil && *existing.ID != "" && existing.ExtendedDatabaseBlobAuditingPolicyProperties != nil && existing.ExtendedDatabaseBlobAuditingPolicyProperties.State != sql.BlobAuditingPolicyStateDisabled {
			return tf.ImportAsExistsError("azurerm_mssql_database_extended_auditing_policy", *existing.ID)
		}
	}

	params := sql.ExtendedDatabaseBlobAuditingPolicy{
		ExtendedDatabaseBlobAuditingPolicyProperties: &sql.ExtendedDatabaseBlobAuditingPolicyProperties{
			State:                      sql.BlobAuditingPolicyStateEnabled,
			StorageEndpoint:            utils.String(d.Get("storage_endpoint").(string)),
			IsStorageSecondaryKeyInUse: utils.Bool(d.Get("storage_account_access_key_is_secondary").(bool)),
			RetentionDays:              utils.Int32(int32(d.Get("retention_in_days").(int))),

			// NOTE: this works around a regression in the Azure API detailed here:
			// https://github.com/Azure/azure-rest-api-specs/issues/11271
			IsAzureMonitorTargetEnabled: utils.Bool(true),
		},
	}

	if v, ok := d.GetOk("storage_account_access_key"); ok {
		params.ExtendedDatabaseBlobAuditingPolicyProperties.StorageAccountAccessKey = utils.String(v.(string))
	}

	if _, err = client.CreateOrUpdate(ctx, dbId.ResourceGroup, dbId.MsSqlServer, dbId.Name, params); err != nil {
		return fmt.Errorf("creating MsSql Database %q Extended Auditing Policy (Sql Server %q / Resource Group %q): %+v", dbId.Name, dbId.MsSqlServer, dbId.ResourceGroup, err)
	}

	read, err := client.Get(ctx, dbId.ResourceGroup, dbId.MsSqlServer, dbId.Name)
	if err != nil {
		return fmt.Errorf("retrieving MsSql Database %q Extended Auditing Policy (MsSql Server Name %q / Resource Group %q): %+v", dbId.Name, dbId.MsSqlServer, dbId.ResourceGroup, err)
	}

	if read.ID == nil || *read.ID == "" {
		return fmt.Errorf("reading MsSql Database %q Extended Auditing Policy (MsSql Server Name %q / Resource Group %q) ID is empty or nil", dbId.Name, dbId.MsSqlServer, dbId.ResourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMsSqlDatabaseExtendedAuditingPolicyRead(d, meta)
}

func resourceArmMsSqlDatabaseExtendedAuditingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabaseExtendedBlobAuditingPoliciesClient
	dbClient := meta.(*clients.Client).MSSQL.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MssqlDatabaseExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MsSqlServer, id.MsDBName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading MsSql Database %s Extended Auditing Policy (MsSql Server Name %q / Resource Group %q): %s", id.MsDBName, id.MsSqlServer, id.ResourceGroup, err)
	}

	dbResp, err := dbClient.Get(ctx, id.ResourceGroup, id.MsSqlServer, id.MsDBName)
	if err != nil || *dbResp.ID == "" {
		return fmt.Errorf("reading MsSql Database %q ID is empty or nil(Resource Group %q): %s", id.MsSqlServer, id.ResourceGroup, err)
	}

	d.Set("database_id", dbResp.ID)

	if props := resp.ExtendedDatabaseBlobAuditingPolicyProperties; props != nil {
		d.Set("storage_endpoint", props.StorageEndpoint)
		d.Set("storage_account_access_key_is_secondary", props.IsStorageSecondaryKeyInUse)
		d.Set("retention_in_days", props.RetentionDays)
	}

	return nil
}

func resourceArmMsSqlDatabaseExtendedAuditingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabaseExtendedBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MssqlDatabaseExtendedAuditingPolicyID(d.Id())
	if err != nil {
		return err
	}

	params := sql.ExtendedDatabaseBlobAuditingPolicy{
		ExtendedDatabaseBlobAuditingPolicyProperties: &sql.ExtendedDatabaseBlobAuditingPolicyProperties{
			State: sql.BlobAuditingPolicyStateDisabled,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.MsSqlServer, id.MsDBName, params); err != nil {
		return fmt.Errorf("deleting MsSql Database %q Extended Auditing Policy( MsSql Server %q / Resource Group %q): %+v", id.MsDBName, id.MsSqlServer, id.ResourceGroup, err)
	}

	return nil
}
