package mssql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmMSSQLManagedInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMSSQLManagedInstanceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"administrator_login": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"identity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
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

			"collation": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"license_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"maintenance_configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"minimal_tls_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"proxy_override": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_data_endpoint_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"storage_size_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"timezone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vcores": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"fully_qualified_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dns_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmMSSQLManagedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ManagedInstancesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("while reading managed SQL instance %s: %v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", azure.NormalizeLocation(*resp.Location))
	d.Set("type", (resp.Type))

	if id := resp.ID; id != nil {
		d.SetId(*id)
	}

	if err := d.Set("identity", flattenMIIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("while setting `identity`: %+v", err)
	}

	if props := resp.ManagedInstanceProperties; props != nil {
		d.Set("fully_qualified_domain_name", props.FullyQualifiedDomainName)
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("subnet_id", props.SubnetID)
		d.Set("state", props.State)
		d.Set("license_type", props.LicenseType)
		d.Set("vcores", props.VCores)
		d.Set("storage_size_gb", props.StorageSizeInGB)
		d.Set("collation", props.Collation)
		d.Set("dns_zone", props.DNSZone)
		d.Set("public_data_endpoint_enabled", props.PublicDataEndpointEnabled)
		d.Set("proxy_override", props.ProxyOverride)
		d.Set("timezone_id", props.TimezoneID)
		d.Set("instance_pool_id", props.InstancePoolID)
		d.Set("maintenance_configuration_id", props.MaintenanceConfigurationID)
		d.Set("minimal_tls_version", props.MinimalTLSVersion)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenMIIdentity(identity *sql.ResourceIdentity) []interface{} {
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
