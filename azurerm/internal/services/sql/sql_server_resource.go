package sql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSqlServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlServerCreateUpdate,
		Read:   resourceSqlServerRead,
		Update: resourceSqlServerCreateUpdate,
		Delete: resourceSqlServerDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlServerName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"2.0",
					"12.0",
				}, true),
			},

			"administrator_login": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"administrator_login_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"connection_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(sql.ServerConnectionTypeDefault),
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.ServerConnectionTypeDefault),
					string(sql.ServerConnectionTypeProxy),
					string(sql.ServerConnectionTypeRedirect),
				}, false),
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"extended_auditing_policy": helper.ExtendedAuditingSchema(),

			"fully_qualified_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSqlServerCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ServersClient
	auditingClient := meta.(*clients.Client).Sql.ServerExtendedBlobAuditingPoliciesClient
	connectionClient := meta.(*clients.Client).Sql.ServerConnectionPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	adminUsername := d.Get("administrator_login").(string)
	version := d.Get("version").(string)

	t := d.Get("tags").(map[string]interface{})
	metadata := tags.Expand(t)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing SQL Server %q (Resource Group %q): %+v", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_sql_server", *existing.ID)
		}
	}

	parameters := sql.Server{
		Location: utils.String(location),
		Tags:     metadata,
		ServerProperties: &sql.ServerProperties{
			Version:            utils.String(version),
			AdministratorLogin: utils.String(adminUsername),
		},
	}

	if _, ok := d.GetOk("identity"); ok {
		sqlServerIdentity := expandAzureRmSqlServerIdentity(d)
		parameters.Identity = sqlServerIdentity
	}

	if d.HasChange("administrator_login_password") {
		adminPassword := d.Get("administrator_login_password").(string)
		parameters.ServerProperties.AdministratorLoginPassword = utils.String(adminPassword)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for SQL Server %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasConflict(future.Response()) {
			return fmt.Errorf("SQL Server names need to be globally unique and %q is already in use.", name)
		}

		return fmt.Errorf("Error waiting on create/update future for SQL Server %q (Resource Group %q): %+v", name, resGroup, err)
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error issuing get request for SQL Server %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	connection := sql.ServerConnectionPolicy{
		ServerConnectionPolicyProperties: &sql.ServerConnectionPolicyProperties{
			ConnectionType: sql.ServerConnectionType(d.Get("connection_policy").(string)),
		},
	}
	if _, err = connectionClient.CreateOrUpdate(ctx, resGroup, name, connection); err != nil {
		return fmt.Errorf("Error issuing create/update request for SQL Server %q Connection Policy (Resource Group %q): %+v", name, resGroup, err)
	}

	auditingProps := sql.ExtendedServerBlobAuditingPolicy{
		ExtendedServerBlobAuditingPolicyProperties: helper.ExpandAzureRmSqlServerBlobAuditingPolicies(d.Get("extended_auditing_policy").([]interface{})),
	}
	if _, err = auditingClient.CreateOrUpdate(ctx, resGroup, name, auditingProps); err != nil {
		return fmt.Errorf("Error issuing create/update request for SQL Server %q Blob Auditing Policies(Resource Group %q): %+v", name, resGroup, err)
	}

	return resourceSqlServerRead(d, meta)
}

func resourceSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ServersClient
	auditingClient := meta.(*clients.Client).Sql.ServerExtendedBlobAuditingPoliciesClient
	connectionClient := meta.(*clients.Client).Sql.ServerConnectionPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] SQL Server %q (Resource Group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving SQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	connection, err := connectionClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Blob Connection Policy for SQL Server %q (Resource Group %q): %v ", id.Name, id.ResourceGroup, err)
	}

	auditingResp, err := auditingClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Blob Auditing Policies for SQL Server %q (Resource Group %q): %v ", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("identity", flattenAzureRmSqlServerIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %+v", err)
	}

	if serverProperties := resp.ServerProperties; serverProperties != nil {
		d.Set("version", serverProperties.Version)
		d.Set("administrator_login", serverProperties.AdministratorLogin)
		d.Set("fully_qualified_domain_name", serverProperties.FullyQualifiedDomainName)
	}

	if props := connection.ServerConnectionPolicyProperties; props != nil {
		d.Set("connection_policy", string(props.ConnectionType))
	}

	if err := d.Set("extended_auditing_policy", helper.FlattenAzureRmSqlServerBlobAuditingPolicies(&auditingResp, d)); err != nil {
		return fmt.Errorf("setting `extended_auditing_policy`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSqlServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting SQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of SQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandAzureRmSqlServerIdentity(d *schema.ResourceData) *sql.ResourceIdentity {
	identities := d.Get("identity").([]interface{})
	if len(identities) == 0 {
		return &sql.ResourceIdentity{}
	}
	identity := identities[0].(map[string]interface{})
	identityType := sql.IdentityType(identity["type"].(string))
	return &sql.ResourceIdentity{
		Type: identityType,
	}
}

func flattenAzureRmSqlServerIdentity(identity *sql.ResourceIdentity) []interface{} {
	if identity == nil {
		return []interface{}{}
	}
	result := make(map[string]interface{})
	result["type"] = identity.Type
	if identity.PrincipalID != nil {
		result["principal_id"] = identity.PrincipalID.String()
	}
	if identity.TenantID != nil {
		result["tenant_id"] = identity.TenantID.String()
	}

	return []interface{}{result}
}
