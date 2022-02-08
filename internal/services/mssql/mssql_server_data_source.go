package mssql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceMsSqlServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMsSqlServerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"administrator_login": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"fully_qualified_domain_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": func() *schema.Schema {
				// TODO: Document the breaking change (user_assigned_identity_ids -> identity_ids) coming in 3.0
				if !features.ThreePointOhBeta() {
					return &schema.Schema{
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"user_assigned_identity_ids": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"principal_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"tenant_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					}
				}

				return commonschema.SystemOrUserAssignedIdentityComputed()
			}(),

			"restorable_dropped_database_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceMsSqlServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	restorableDroppedDatabasesClient := meta.(*clients.Client).MSSQL.RestorableDroppedDatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	d.SetId(id.ID())
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.ServerProperties; props != nil {
		d.Set("version", props.Version)
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("fully_qualified_domain_name", props.FullyQualifiedDomainName)
	}

	identity, err := flattenSqlServerIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	restorableListPage, err := restorableDroppedDatabasesClient.ListByServerComplete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("listing %s Restorable Dropped Databases: %v", id, err)
	}
	if err := d.Set("restorable_dropped_database_ids", flattenSqlServerRestorableDatabases(restorableListPage.Response())); err != nil {
		return fmt.Errorf("setting `restorable_dropped_database_ids`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
