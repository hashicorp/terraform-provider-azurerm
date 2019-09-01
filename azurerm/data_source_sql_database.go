package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
)

func dataSourceSqlDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSqlDatabaseRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"create_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_database_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"restore_point_in_time": {
				Type:         schema.TypeString,
				Computed:     true,
			},

			"edition": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"collation": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"max_size_bytes": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"requested_service_objective_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"requested_service_objective_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"source_database_deletion_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"elastic_pool_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"encryption": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_secondary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"read_scale": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sql.DatabasesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, serverName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Sql Database %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Sql Database %q (Resource Group %q): %s", name, resourceGroup, err)
	}

	if id := resp.ID; id != nil {
		d.SetId(*resp.ID)
	}

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.DatabaseProperties; props != nil {
		d.Set("collation", props.Collation)
		d.Set("collation", props.Collation)
		d.Set("default_secondary_location", props.DefaultSecondaryLocation)
		d.Set("edition", string(props.Edition))
		d.Set("elastic_pool_name", props.ElasticPoolName)
		d.Set("max_size_bytes", props.MaxSizeBytes)
		d.Set("requested_service_objective_name", string(props.RequestedServiceObjectiveName))

		if cd := props.CreationDate; cd != nil {
			d.Set("creation_date", cd.String())
		}

		if rsoid := props.RequestedServiceObjectiveID; rsoid != nil {
			d.Set("requested_service_objective_id", rsoid.String())
		}

		if rpit := props.RestorePointInTime; rpit != nil {
			d.Set("restore_point_in_time", rpit.String())
		}

		if sddd := props.SourceDatabaseDeletionDate; sddd != nil {
			d.Set("source_database_deletion_date", sddd.String())
		}

		d.Set("encryption", flattenEncryptionStatus(props.TransparentDataEncryption))

		readScale := props.ReadScale
		if readScale == sql.ReadScaleEnabled {
			d.Set("read_scale", true)
		} else {
			d.Set("read_scale", false)
		}
	}

	//
	//
	//d.Set("name", resp.Name)
	//d.Set("resource_group_name", resourceGroup)
	//if location := resp.Location; location != nil {
	//	d.Set("location", azure.NormalizeLocation(*location))
	//}
	//
	//d.Set("server_name", serverName)
	//
	//if props := resp.DatabaseProperties; props != nil {
	//	// TODO: set `create_mode` & `source_database_id` once this issue is fixed:
	//	// https://github.com/Azure/azure-rest-api-specs/issues/1604
	//

	//}

	return tags.FlattenAndSet(d, resp.Tags)
}
