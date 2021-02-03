package sql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSqlFailoverGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlFailoverGroupCreateUpdate,
		Read:   resourceSqlFailoverGroupRead,
		Update: resourceSqlFailoverGroupCreateUpdate,
		Delete: resourceSqlFailoverGroupDelete,

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
				ValidateFunc: helper.ValidateMsSqlFailoverGroupName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: helper.ValidateMsSqlServerName,
			},

			"databases": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"partner_servers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"location": azure.SchemaLocationForDataSource(),

						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"readonly_endpoint_failover_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sql.ReadOnlyEndpointFailoverPolicyDisabled),
								string(sql.ReadOnlyEndpointFailoverPolicyEnabled),
							}, false),
						},
					},
				},
			},

			"read_write_endpoint_failover_policy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sql.Automatic),
								string(sql.Manual),
							}, false),
						},
						"grace_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
			},

			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSqlFailoverGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.FailoverGroupsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serverName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing SQL Failover Group %q (Resource Group %q, Server %q): %+v", name, resourceGroup, serverName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_sql_failover_group", *existing.ID)
		}
	}

	t := d.Get("tags").(map[string]interface{})

	properties := sql.FailoverGroup{
		FailoverGroupProperties: &sql.FailoverGroupProperties{
			ReadOnlyEndpoint:  expandSqlFailoverGroupReadOnlyPolicy(d),
			ReadWriteEndpoint: expandSqlFailoverGroupReadWritePolicy(d),
			PartnerServers:    expandSqlFailoverGroupPartnerServers(d),
		},
		Tags: tags.Expand(t),
	}

	if r, ok := d.Get("databases").(*schema.Set); ok && r.Len() > 0 {
		var databases []string
		for _, v := range r.List() {
			s := v.(string)
			databases = append(databases, s)
		}

		properties.Databases = &databases
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, serverName, name, properties)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for SQL Failover Group %q (Resource Group %q, Server %q): %+v", name, resourceGroup, serverName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for SQL Failover Group %q (Resource Group %q, Server %q): %+v", name, resourceGroup, serverName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serverName, name)
	if err != nil {
		return fmt.Errorf("Error issuing get request for SQL Failover Group %q (Resource Group %q, Server %q): %+v", name, resourceGroup, serverName, err)
	}

	d.SetId(*resp.ID)

	return resourceSqlFailoverGroupRead(d, meta)
}

func resourceSqlFailoverGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.FailoverGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FailoverGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Failover Group %q (Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.FailoverGroupProperties; props != nil {
		if err := d.Set("read_write_endpoint_failover_policy", flattenSqlFailoverGroupReadWritePolicy(props.ReadWriteEndpoint)); err != nil {
			return fmt.Errorf("setting `read_write_endpoint_failover_policy`: %+v", err)
		}

		if err := d.Set("readonly_endpoint_failover_policy", flattenSqlFailoverGroupReadOnlyPolicy(props.ReadOnlyEndpoint)); err != nil {
			return fmt.Errorf("setting `read_only_endpoint_failover_policy`: %+v", err)
		}

		if props.Databases != nil {
			d.Set("databases", set.FromStringSlice(*props.Databases))
		}
		d.Set("role", string(props.ReplicationRole))

		if err := d.Set("partner_servers", flattenSqlFailoverGroupPartnerServers(props.PartnerServers)); err != nil {
			return fmt.Errorf("Error setting `partner_servers`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSqlFailoverGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.FailoverGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FailoverGroupID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting SQL Failover Group %q (Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of SQL Failover Group %q (Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	return err
}

func expandSqlFailoverGroupReadWritePolicy(d *schema.ResourceData) *sql.FailoverGroupReadWriteEndpoint {
	vs := d.Get("read_write_endpoint_failover_policy").([]interface{})
	v := vs[0].(map[string]interface{})

	mode := sql.ReadWriteEndpointFailoverPolicy(v["mode"].(string))
	graceMins := int32(v["grace_minutes"].(int))

	policy := &sql.FailoverGroupReadWriteEndpoint{
		FailoverPolicy: mode,
	}

	if mode != sql.Manual {
		policy.FailoverWithDataLossGracePeriodMinutes = utils.Int32(graceMins)
	}

	return policy
}

func flattenSqlFailoverGroupReadWritePolicy(input *sql.FailoverGroupReadWriteEndpoint) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	policy := make(map[string]interface{})

	policy["mode"] = string(input.FailoverPolicy)

	if input.FailoverWithDataLossGracePeriodMinutes != nil {
		policy["grace_minutes"] = *input.FailoverWithDataLossGracePeriodMinutes
	}
	return []interface{}{policy}
}

func expandSqlFailoverGroupReadOnlyPolicy(d *schema.ResourceData) *sql.FailoverGroupReadOnlyEndpoint {
	vs := d.Get("readonly_endpoint_failover_policy").([]interface{})
	if len(vs) == 0 {
		return nil
	}

	v := vs[0].(map[string]interface{})
	mode := sql.ReadOnlyEndpointFailoverPolicy(v["mode"].(string))

	return &sql.FailoverGroupReadOnlyEndpoint{
		FailoverPolicy: mode,
	}
}

func flattenSqlFailoverGroupReadOnlyPolicy(input *sql.FailoverGroupReadOnlyEndpoint) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	policy := make(map[string]interface{})
	policy["mode"] = string(input.FailoverPolicy)

	return []interface{}{policy}
}

func expandSqlFailoverGroupPartnerServers(d *schema.ResourceData) *[]sql.PartnerInfo {
	servers := d.Get("partner_servers").([]interface{})
	partners := make([]sql.PartnerInfo, 0)

	for _, server := range servers {
		info := server.(map[string]interface{})

		id := info["id"].(string)
		partners = append(partners, sql.PartnerInfo{
			ID: &id,
		})
	}

	return &partners
}

func flattenSqlFailoverGroupPartnerServers(input *[]sql.PartnerInfo) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	if input != nil {
		for _, server := range *input {
			info := make(map[string]interface{})

			if v := server.ID; v != nil {
				info["id"] = *v
			}
			if v := server.Location; v != nil {
				info["location"] = *v
			}
			info["role"] = string(server.ReplicationRole)

			result = append(result, info)
		}
	}
	return result
}
