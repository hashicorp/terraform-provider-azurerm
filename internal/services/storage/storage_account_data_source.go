// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceStorageAccount() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourceStorageAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.StorageAccountName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"account_kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"account_tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"account_replication_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"access_tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"custom_domain": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"https_traffic_only_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"min_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"allow_nested_items_to_be_public": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"is_hns_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"nfsv3_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"primary_location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_blob_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_blob_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_blob_internet_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_blob_internet_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_blob_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_blob_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_blob_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_blob_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_blob_internet_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_blob_internet_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_blob_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_blob_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_queue_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_queue_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_queue_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_queue_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_queue_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_queue_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_queue_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_queue_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_table_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_table_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_table_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_table_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_table_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_table_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_table_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_table_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_web_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_web_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_web_internet_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_web_internet_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_web_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_web_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_web_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_web_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_web_internet_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_web_internet_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_web_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_web_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_dfs_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_dfs_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_dfs_internet_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_dfs_internet_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_dfs_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_dfs_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_dfs_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_dfs_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_dfs_internet_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_dfs_internet_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_dfs_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_dfs_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_file_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_file_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_file_internet_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_file_internet_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_file_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_file_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_file_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_file_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_file_internet_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_file_internet_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_file_microsoft_host": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_file_microsoft_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_blob_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_blob_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"queue_encryption_key_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"table_encryption_key_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"infrastructure_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"azure_files_authentication": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"directory_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"active_directory": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"domain_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"netbios_domain_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"forest_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"domain_guid": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"domain_sid": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"storage_sid": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
						"default_share_level_permission": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"dns_endpoint_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["enable_https_traffic_only"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Computed: true,
		}
	}

	return resource
}

func dataSourceStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageDomainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Storage domain suffix for environment %q", meta.(*clients.Client).Account.Environment.Name)
	}

	id := commonids.NewStorageAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.GetProperties(ctx, id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	listKeysOpts := storageaccounts.DefaultListKeysOperationOptions()
	listKeysOpts.Expand = pointer.To(storageaccounts.ListKeyExpandKerb)
	keys, err := client.ListKeys(ctx, id, listKeysOpts)
	if err != nil {
		hasWriteLock := response.WasConflict(keys.HttpResponse)
		doesntHavePermissions := response.WasForbidden(keys.HttpResponse) || response.WasStatusCode(keys.HttpResponse, http.StatusUnauthorized)
		if !hasWriteLock && !doesntHavePermissions {
			return fmt.Errorf("listing Keys for %s: %+v", id, err)
		}
	}

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("account_kind", string(pointer.From(model.Kind)))

		// NOTE: we should expose EdgeZone in the future

		if sku := model.Sku; sku != nil {
			d.Set("account_tier", pointer.From(sku.Tier))
			d.Set("account_replication_type", strings.Split(string(sku.Name), "_")[1])
		}

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}

		if props := model.Properties; props != nil {
			d.Set("access_tier", pointer.From(props.AccessTier))
			d.Set("allow_nested_items_to_be_public", pointer.From(props.AllowBlobPublicAccess))
			if err := d.Set("custom_domain", flattenAccountCustomDomain(props.CustomDomain)); err != nil {
				return fmt.Errorf("setting `custom_domain`: %+v", err)
			}
			d.Set("https_traffic_only_enabled", pointer.From(props.SupportsHTTPSTrafficOnly))
			d.Set("is_hns_enabled", pointer.From(props.IsHnsEnabled))
			d.Set("nfsv3_enabled", pointer.From(props.IsNfsV3Enabled))
			d.Set("primary_location", location.NormalizeNilable(props.PrimaryLocation))
			d.Set("secondary_location", location.NormalizeNilable(props.SecondaryLocation))

			if !features.FourPointOhBeta() {
				d.Set("enable_https_traffic_only", pointer.From(props.SupportsHTTPSTrafficOnly))
			}

			// Setting the encryption key type to "Service" in PUT. The following GET will not return the queue/table in the service list of its response.
			// So defaults to setting the encryption key type to "Service" if it is absent in the GET response. Also, define the default value as "Service" in the schema.
			infrastructureEncryption := false
			queueEncryptionKeyType := string(storageaccounts.KeyTypeService)
			tableEncryptionKeyType := string(storageaccounts.KeyTypeService)
			if encryption := props.Encryption; encryption != nil {
				infrastructureEncryption = pointer.From(encryption.RequireInfrastructureEncryption)

				if encryption.Services != nil {
					if encryption.Services.Queue != nil && encryption.Services.Queue.KeyType != nil {
						queueEncryptionKeyType = string(*encryption.Services.Queue.KeyType)
					}
					if encryption.Services.Table != nil && encryption.Services.Table.KeyType != nil {
						tableEncryptionKeyType = string(*encryption.Services.Table.KeyType)
					}
				}
			}
			d.Set("infrastructure_encryption_enabled", infrastructureEncryption)
			d.Set("queue_encryption_key_type", queueEncryptionKeyType)
			d.Set("table_encryption_key_type", tableEncryptionKeyType)

			// For storage account created using old API, the response of GET call will not return "min_tls_version"
			minTlsVersion := string(storageaccounts.MinimumTlsVersionTLSOneZero)
			if props.MinimumTlsVersion != nil {
				minTlsVersion = string(*props.MinimumTlsVersion)
			}
			d.Set("min_tls_version", minTlsVersion)

			// DNSEndpointType is null when unconfigured - therefore default this to Standard
			dnsEndpointType := storageaccounts.DnsEndpointTypeStandard
			if props.DnsEndpointType != nil {
				dnsEndpointType = *props.DnsEndpointType
			}
			d.Set("dns_endpoint_type", string(dnsEndpointType))

			if err := d.Set("azure_files_authentication", flattenAccountAzureFilesAuthentication(props.AzureFilesIdentityBasedAuthentication)); err != nil {
				return fmt.Errorf("setting `azure_files_authentication`: %+v", err)
			}
		}
	}

	var primaryEndpoints *storageaccounts.Endpoints
	var secondaryEndpoints *storageaccounts.Endpoints
	var routingPreference *storageaccounts.RoutingPreference
	if model := resp.Model; model != nil && model.Properties != nil {
		primaryEndpoints = model.Properties.PrimaryEndpoints
		routingPreference = model.Properties.RoutingPreference
		secondaryEndpoints = model.Properties.SecondaryEndpoints
	}
	endpoints := flattenAccountEndpoints(primaryEndpoints, secondaryEndpoints, routingPreference)
	if err := endpoints.set(d); err != nil {
		return err
	}

	storageAccountKeys := make([]storageaccounts.StorageAccountKey, 0)
	if keys.Model != nil && keys.Model.Keys != nil {
		storageAccountKeys = *keys.Model.Keys
	}
	keysAndConnectionStrings := flattenAccountAccessKeysAndConnectionStrings(id.StorageAccountName, *storageDomainSuffix, storageAccountKeys, endpoints)
	if err := keysAndConnectionStrings.set(d); err != nil {
		return err
	}

	return nil
}
