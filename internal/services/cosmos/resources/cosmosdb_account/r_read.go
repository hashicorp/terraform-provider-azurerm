// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_account

import (
	"bytes"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbAccountRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseDatabaseAccountID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.DatabaseAccountsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.DatabaseAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := existing.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("kind", pointer.From(model.Kind))
		d.Set("tags", model.Tags)

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("offer_type", pointer.From(props.DatabaseAccountOfferType))
			d.Set("ip_range_filter", common.CosmosDBIpRulesToIpRangeFilter(props.IPRules))
			d.Set("endpoint", props.DocumentEndpoint)
			d.Set("free_tier_enabled", props.EnableFreeTier)
			d.Set("analytical_storage_enabled", props.EnableAnalyticalStorage)
			d.Set("public_network_access_enabled", pointer.From(props.PublicNetworkAccess) == cosmosdb.PublicNetworkAccessEnabled)

			d.Set("default_identity_type", "FirstPartyIdentity")
			if pointer.From(props.DefaultIdentity) != "" {
				d.Set("default_identity_type", props.DefaultIdentity)
			}

			d.Set("minimal_tls_version", pointer.From(props.MinimalTlsVersion))
			d.Set("create_mode", pointer.From(props.CreateMode))
			d.Set("partition_merge_enabled", pointer.From(props.EnablePartitionMerge))
			d.Set("burst_capacity_enabled", pointer.From(props.EnableBurstCapacity))
			d.Set("virtual_network_filter_enabled", pointer.From(props.IsVirtualNetworkFilterEnabled))
			d.Set("automatic_failover_enabled", pointer.From(props.EnableAutomaticFailover))
			d.Set("multiple_write_locations_enabled", pointer.From(props.EnableMultipleWriteLocations))

			if !features.SixPointOh() {
				d.Set("is_virtual_network_filter_enabled", pointer.From(props.IsVirtualNetworkFilterEnabled))
			}

			if v := props.KeyVaultKeyUri; v != nil {
				key, err := keyvault.ParseNestedItemID(*v, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
				if err != nil {
					return err
				}

				d.Set("key_vault_key_id", key.VersionlessID())
			}

			if err := d.Set("analytical_storage", flattenCosmosDBAccountAnalyticalStorageConfiguration(props.AnalyticalStorageConfiguration)); err != nil {
				return fmt.Errorf("setting `analytical_storage`: %+v", err)
			}

			if err := d.Set("capacity", flattenCosmosDBAccountCapacity(props.Capacity)); err != nil {
				return fmt.Errorf("setting `capacity`: %+v", err)
			}

			if err := d.Set("restore", flattenCosmosdbAccountRestoreParameters(props.RestoreParameters)); err != nil {
				return fmt.Errorf("setting `restore`: %+v", err)
			}

			if err = d.Set("consistency_policy", flattenAzureRmCosmosDBAccountConsistencyPolicy(props.ConsistencyPolicy)); err != nil {
				return fmt.Errorf("setting `consistency_policy`: %+v", err)
			}

			if err = d.Set("geo_location", flattenAzureRmCosmosDBAccountGeoLocations(props)); err != nil {
				return fmt.Errorf("setting `geo_location`: %+v", err)
			}

			if err = d.Set("capabilities", flattenAzureRmCosmosDBAccountCapabilities(props.Capabilities)); err != nil {
				return fmt.Errorf("setting `capabilities`: %+v", err)
			}

			if err = d.Set("virtual_network_rule", flattenAzureRmCosmosDBAccountVirtualNetworkRules(props.VirtualNetworkRules)); err != nil {
				return fmt.Errorf("setting `virtual_network_rule`: %+v", err)
			}

			accessKeyMetadataWritesEnabled := true
			if v := props.DisableKeyBasedMetadataWriteAccess; v != nil {
				accessKeyMetadataWritesEnabled = !*v
			}
			d.Set("access_key_metadata_writes_enabled", accessKeyMetadataWritesEnabled)

			if apiProps := props.ApiProperties; apiProps != nil {
				d.Set("mongo_server_version", pointer.From(apiProps.ServerVersion))
			}

			d.Set("network_acl_bypass_for_azure_services", pointer.From(props.NetworkAclBypass) == cosmosdb.NetworkAclBypassAzureServices)
			d.Set("network_acl_bypass_ids", pluginsdk.FlattenSlice(props.NetworkAclBypassResourceIds))
			d.Set("local_authentication_enabled", !pointer.From(props.DisableLocalAuth))

			policy, err := flattenCosmosdbAccountBackup(props.BackupPolicy)
			if err != nil {
				return fmt.Errorf("flattening `backup`: %w", err)
			}

			if err = d.Set("backup", policy); err != nil {
				return fmt.Errorf("setting `backup`: %+v", err)
			}

			d.Set("cors_rule", common.FlattenCosmosCorsRule(props.Cors))

			if err := d.Set("read_endpoints", flattenCosmosdbAccountReadWriteEndpoints(props.Locations)); err != nil {
				return fmt.Errorf("setting `read_endpoints`: %s", err)
			}

			if err := d.Set("write_endpoints", flattenCosmosdbAccountReadWriteEndpoints(props.WriteLocations)); err != nil {
				return fmt.Errorf("setting `write_endpoints`: %s", err)
			}

			// ListKeys returns a data structure containing a DatabaseAccountListReadOnlyKeysResult pointer
			// implying that it also returns the read only keys, however this appears to not be the case
			keys, err := client.DatabaseAccountsListKeys(ctx, *id)
			if err != nil {
				if response.WasNotFound(keys.HttpResponse) {
					log.Printf("[DEBUG] Keys were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.DatabaseAccountName, id.ResourceGroupName)
					d.SetId("")
					return nil
				}

				return fmt.Errorf("listing keys for %s: %w", id, err)
			}

			if keys.Model != nil {
				d.Set("primary_key", keys.Model.PrimaryMasterKey)
				d.Set("secondary_key", keys.Model.SecondaryMasterKey)
			}

			readonlyKeys, err := client.DatabaseAccountsListReadOnlyKeys(ctx, *id)
			if err != nil {
				if response.WasNotFound(keys.HttpResponse) {
					log.Printf("[DEBUG] Read Only Keys were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.DatabaseAccountName, id.ResourceGroupName)
					d.SetId("")
					return nil
				}

				return fmt.Errorf("listing read-only keys for %s: %w", id, err)
			}
			if readonlyKeys.Model != nil {
				d.Set("primary_readonly_key", readonlyKeys.Model.PrimaryReadonlyMasterKey)
				d.Set("secondary_readonly_key", readonlyKeys.Model.SecondaryReadonlyMasterKey)
			}

			connStringResp, err := client.DatabaseAccountsListConnectionStrings(ctx, *id)
			if err != nil {
				if response.WasNotFound(keys.HttpResponse) {
					log.Printf("[DEBUG] Connection Strings were not found for CosmosDB Account %q (Resource Group %q) - removing from state!", id.ResourceGroupName, id.ResourceGroupName)
					d.SetId("")
					return nil
				}

				return fmt.Errorf("listing connection strings for %s: %w", id, err)
			}

			var connStrings []string
			if connStringResp.Model.ConnectionStrings != nil {
				connStrings = make([]string, len(*connStringResp.Model.ConnectionStrings))
				for i, v := range *connStringResp.Model.ConnectionStrings {
					connStrings[i] = *v.ConnectionString
					if propertyName, propertyExists := connStringPropertyMap[*v.Description]; propertyExists {
						d.Set(propertyName, v.ConnectionString) // lintignore:R001
					}
				}
			}
		}
	}

	return nil
}

var connStringPropertyMap = map[string]string{
	"Primary SQL Connection String":                 "primary_sql_connection_string",
	"Secondary SQL Connection String":               "secondary_sql_connection_string",
	"Primary Read-Only SQL Connection String":       "primary_readonly_sql_connection_string",
	"Secondary Read-Only SQL Connection String":     "secondary_readonly_sql_connection_string",
	"Primary MongoDB Connection String":             "primary_mongodb_connection_string",
	"Secondary MongoDB Connection String":           "secondary_mongodb_connection_string",
	"Primary Read-Only MongoDB Connection String":   "primary_readonly_mongodb_connection_string",
	"Secondary Read-Only MongoDB Connection String": "secondary_readonly_mongodb_connection_string",
}

type databaseAccountCapabilities string

const (
	databaseAccountCapabilitiesEnableAggregationPipeline         databaseAccountCapabilities = "EnableAggregationPipeline"
	databaseAccountCapabilitiesEnableCassandra                   databaseAccountCapabilities = "EnableCassandra"
	databaseAccountCapabilitiesEnableGremlin                     databaseAccountCapabilities = "EnableGremlin"
	databaseAccountCapabilitiesEnableTable                       databaseAccountCapabilities = "EnableTable"
	databaseAccountCapabilitiesEnableServerless                  databaseAccountCapabilities = "EnableServerless"
	databaseAccountCapabilitiesEnableMongo                       databaseAccountCapabilities = "EnableMongo"
	databaseAccountCapabilitiesEnableMongo16MBDocumentSupport    databaseAccountCapabilities = "EnableMongo16MBDocumentSupport"
	databaseAccountCapabilitiesMongoDBv34                        databaseAccountCapabilities = "MongoDBv3.4"
	databaseAccountCapabilitiesMongoEnableDocLevelTTL            databaseAccountCapabilities = "mongoEnableDocLevelTTL"
	databaseAccountCapabilitiesDeleteAllItemsByPartitionKey      databaseAccountCapabilities = "DeleteAllItemsByPartitionKey"
	databaseAccountCapabilitiesDisableRateLimitingResponses      databaseAccountCapabilities = "DisableRateLimitingResponses"
	databaseAccountCapabilitiesAllowSelfServeUpgradeToMongo36    databaseAccountCapabilities = "AllowSelfServeUpgradeToMongo36"
	databaseAccountCapabilitiesEnableMongoRetryableWrites        databaseAccountCapabilities = "EnableMongoRetryableWrites"
	databaseAccountCapabilitiesEnableMongoRoleBasedAccessControl databaseAccountCapabilities = "EnableMongoRoleBasedAccessControl"
	databaseAccountCapabilitiesEnableUniqueCompoundNestedDocs    databaseAccountCapabilities = "EnableUniqueCompoundNestedDocs"
	databaseAccountCapabilitiesEnableNoSqlVectorSearch           databaseAccountCapabilities = "EnableNoSQLVectorSearch"
	databaseAccountCapabilitiesEnableNoSqlFullTextSearch         databaseAccountCapabilities = "EnableNoSQLFullTextSearch"
	databaseAccountCapabilitiesEnableTtlOnCustomPath             databaseAccountCapabilities = "EnableTtlOnCustomPath"
	databaseAccountCapabilitiesEnablePartialUniqueIndex          databaseAccountCapabilities = "EnablePartialUniqueIndex"
	databaseAccountCapabilitiesEnableFabricNetworkAclBypass      databaseAccountCapabilities = "EnableFabricNetworkAclBypass"
)

/*
	The mapping of capabilities and kinds of cosmosdb account confirmed by service team is as follows:

EnableMongo :                    	MongoDB
EnableCassandra :                	GlobalDocumentDB, Parse
EnableGremlin :                  	GlobalDocumentDB, Parse
EnableTable :                    	GlobalDocumentDB, Parse
EnableAggregationPipeline :      	GlobalDocumentDB, MongoDB, Parse
EnableServerless :               	GlobalDocumentDB, MongoDB, Parse
MongoDBv3.4 :                    	GlobalDocumentDB, MongoDB, Parse
mongoEnableDocLevelTTL :         	GlobalDocumentDB, MongoDB, Parse
DeleteAllItemsByPartitionKey :   	GlobalDocumentDB, MongoDB, Parse
DisableRateLimitingResponses :   	GlobalDocumentDB, MongoDB, Parse
AllowSelfServeUpgradeToMongo36 : 	GlobalDocumentDB, MongoDB, Parse
EnableMongoRetryableWrites :		MongoDB
EnableMongoRoleBasedAccessControl : MongoDB
EnableUniqueCompoundNestedDocs : 	MongoDB
EnableTtlOnCustomPath:              MongoDB
EnablePartialUniqueIndex:           MongoDB
EnableFabricNetworkAclBypass:    	GlobalDocumentDB
*/
var capabilitiesToKindMap = map[string]any{
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongo)):                       []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongo16MBDocumentSupport)):    []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRoleBasedAccessControl)): []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRetryableWrites)):        []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableUniqueCompoundNestedDocs)):    []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableTtlOnCustomPath)):             []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnablePartialUniqueIndex)):          []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableNoSqlVectorSearch)):           []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableNoSqlFullTextSearch)):         []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableCassandra)):                   []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableGremlin)):                     []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableTable)):                       []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableServerless)):                  []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableAggregationPipeline)):         []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesMongoDBv34)):                        []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesMongoEnableDocLevelTTL)):            []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesDeleteAllItemsByPartitionKey)):      []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesDisableRateLimitingResponses)):      []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesAllowSelfServeUpgradeToMongo36)):    []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindMongoDB)), strings.ToLower(string(cosmosdb.DatabaseAccountKindParse))},
	strings.ToLower(string(databaseAccountCapabilitiesEnableFabricNetworkAclBypass)):      []string{strings.ToLower(string(cosmosdb.DatabaseAccountKindGlobalDocumentDB))},
}

// If the consistency policy of the Cosmos DB Database Account is not bounded staleness,
// any changes to the configuration for bounded staleness should be suppressed.
func suppressConsistencyPolicyStalenessConfiguration(_, _, _ string, d *pluginsdk.ResourceData) bool {
	consistencyPolicyList := d.Get("consistency_policy").([]any)
	if len(consistencyPolicyList) == 0 || consistencyPolicyList[0] == nil {
		return false
	}

	consistencyPolicy := consistencyPolicyList[0].(map[string]any)

	return consistencyPolicy["consistency_level"].(string) != string(cosmosdb.DefaultConsistencyLevelBoundedStaleness)
}

func findZoneRedundant(locations *[]cosmosdb.Location, id string) bool {
	if locations == nil {
		return false
	}

	for _, location := range *locations {
		if location.Id != nil && *location.Id == id {
			if location.IsZoneRedundant != nil {
				return *location.IsZoneRedundant
			}
		}
	}

	return false
}

func resourceAzureRMCosmosDBAccountGeoLocationHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		location := location.Normalize(m["location"].(string))
		priority := int32(m["failover_priority"].(int))
		zone_redundant := m["zone_redundant"].(bool)

		fmt.Fprintf(&buf, "%s-%d-%t", location, priority, zone_redundant)
	}

	return pluginsdk.HashString(buf.String())
}

func resourceAzureRMCosmosDBAccountCapabilitiesHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		fmt.Fprintf(&buf, "%s-", m["name"].(string))
	}

	return pluginsdk.HashString(buf.String())
}

func resourceAzureRMCosmosDBAccountVirtualNetworkRuleHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		buf.WriteString(strings.ToLower(m["id"].(string)))
	}

	return pluginsdk.HashString(buf.String())
}

func checkCapabilitiesCanBeUpdated(kind string, oldCapabilities *[]cosmosdb.Capability, newCapabilities *[]cosmosdb.Capability) bool {
	// The feedback from service team : capabilities that can be added to an existing account
	canBeAddedCaps := []string{
		strings.ToLower(string(databaseAccountCapabilitiesDeleteAllItemsByPartitionKey)),
		strings.ToLower(string(databaseAccountCapabilitiesDisableRateLimitingResponses)),
		strings.ToLower(string(databaseAccountCapabilitiesAllowSelfServeUpgradeToMongo36)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableAggregationPipeline)),
		strings.ToLower(string(databaseAccountCapabilitiesMongoDBv34)),
		strings.ToLower(string(databaseAccountCapabilitiesMongoEnableDocLevelTTL)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableMongo16MBDocumentSupport)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRetryableWrites)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRoleBasedAccessControl)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableUniqueCompoundNestedDocs)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableTtlOnCustomPath)),
		strings.ToLower(string(databaseAccountCapabilitiesEnablePartialUniqueIndex)),
		strings.ToLower(string(databaseAccountCapabilitiesEnableFabricNetworkAclBypass)),
	}

	// The feedback from service team: capabilities that can be removed from an existing account
	canBeRemovedCaps := []string{
		strings.ToLower(string(databaseAccountCapabilitiesEnableMongoRetryableWrites)),
		strings.ToLower(string(databaseAccountCapabilitiesDisableRateLimitingResponses)),
	}

	// first check the new capabilities can be added
	for _, capability := range *newCapabilities {
		if capability.Name == nil {
			continue
		}
		existedPreviously := false
		for _, existing := range *oldCapabilities {
			if existing.Name != nil && strings.EqualFold(*existing.Name, *capability.Name) {
				existedPreviously = true
				break
			}
		}
		if existedPreviously {
			continue
		}

		// retrieve a list of capabilities for this DB type
		supportedKindsForCapability, ok := capabilitiesToKindMap[strings.ToLower(*capability.Name)]
		if !ok {
			return false
		}

		// first check if this is supported
		if isSupported := slices.Contains(supportedKindsForCapability.([]string), strings.ToLower(kind)); !isSupported {
			return false
		}

		// then check if it can be added via an update
		if !slices.Contains(canBeAddedCaps, strings.ToLower(*capability.Name)) {
			return false
		}
	}
	// then check if we're removing any that they can be removed
	for _, capability := range *oldCapabilities {
		existsNow := false
		for _, new := range *newCapabilities {
			if new.Name != nil && strings.EqualFold(*new.Name, *capability.Name) {
				existsNow = true
				break
			}
		}
		if existsNow {
			continue
		}

		if !slices.Contains(canBeRemovedCaps, strings.ToLower(*capability.Name)) {
			return false
		}
	}

	return true
}

func prepareCapabilities(capabilities any) *[]cosmosdb.Capability {
	output := make([]cosmosdb.Capability, 0)
	for _, v := range capabilities.(*pluginsdk.Set).List() {
		m := v.(map[string]any)
		if c, ok := m["name"].(string); ok {
			capability := cosmosdb.Capability{
				Name: pointer.To(c),
			}
			output = append(output, capability)
		}
	}
	return &output
}
