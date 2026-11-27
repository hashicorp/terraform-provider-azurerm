// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_account

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbAccountCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewDatabaseAccountID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.DatabaseAccountsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cosmosdb_account", id.ID())
		}
	}

	databaseAccountNameID := cosmosdb.NewDatabaseAccountNameID(id.DatabaseAccountName)
	dbCheckNameResp, err := client.DatabaseAccountsCheckNameExists(ctx, databaseAccountNameID)
	if !response.WasNotFound(dbCheckNameResp.HttpResponse) {
		if err != nil {
			// TODO: remove when https://github.com/Azure/azure-sdk-for-go/issues/9891 is fixed
			if !response.WasStatusCode(dbCheckNameResp.HttpResponse, http.StatusInternalServerError) {
				return fmt.Errorf("checking whether the CosmosDB Account name (`%s`) is available: %+v", id.DatabaseAccountName, err)
			}
		}
		return fmt.Errorf("a CosmosDB Account with the chosen name (`%s`) already exists, please specify a different name", id.DatabaseAccountName)
	}
	geoLocations, err := expandAzureRmCosmosDBAccountGeoLocations(d)
	if err != nil {
		return fmt.Errorf("expanding `geo_location`: %+v", err)
	}

	expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	account := cosmosdb.DatabaseAccountCreateUpdateParameters{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Kind:     pointer.ToEnum[cosmosdb.DatabaseAccountKind](d.Get("kind").(string)),
		Identity: expandedIdentity,
		Properties: cosmosdb.DatabaseAccountCreateUpdateProperties{
			DatabaseAccountOfferType:           cosmosdb.DatabaseAccountOfferType(d.Get("offer_type").(string)),
			IPRules:                            common.CosmosDBIpRangeFilterToIpRules(*pluginsdk.ExpandStringSlice(d.Get("ip_range_filter").(*pluginsdk.Set).List())),
			IsVirtualNetworkFilterEnabled:      pointer.To(d.Get("virtual_network_filter_enabled").(bool)),
			EnableFreeTier:                     pointer.To(d.Get("free_tier_enabled").(bool)),
			EnableAutomaticFailover:            pointer.To(d.Get("automatic_failover_enabled").(bool)),
			ConsistencyPolicy:                  expandAzureRmCosmosDBAccountConsistencyPolicy(d),
			Locations:                          geoLocations,
			Capabilities:                       expandAzureRmCosmosDBAccountCapabilities(d),
			MinimalTlsVersion:                  pointer.ToEnum[cosmosdb.MinimalTlsVersion](d.Get("minimal_tls_version").(string)),
			VirtualNetworkRules:                expandAzureRmCosmosDBAccountVirtualNetworkRules(d),
			EnableMultipleWriteLocations:       pointer.To(d.Get("multiple_write_locations_enabled").(bool)),
			EnablePartitionMerge:               pointer.To(d.Get("partition_merge_enabled").(bool)),
			EnableBurstCapacity:                pointer.To(d.Get("burst_capacity_enabled").(bool)),
			PublicNetworkAccess:                expandCosmosdbAccountPublicNetworkAccess(d.Get("public_network_access_enabled").(bool)),
			EnableAnalyticalStorage:            pointer.To(d.Get("analytical_storage_enabled").(bool)),
			Cors:                               common.ExpandCosmosCorsRule(d.Get("cors_rule").([]any)),
			DisableKeyBasedMetadataWriteAccess: pointer.To(!d.Get("access_key_metadata_writes_enabled").(bool)),
			NetworkAclBypass:                   expandCosmosdbAccountNetworkBypass(d.Get("network_acl_bypass_for_azure_services").(bool)),
			NetworkAclBypassResourceIds:        pluginsdk.ExpandStringSlice(d.Get("network_acl_bypass_ids").([]any)),
			DisableLocalAuth:                   pointer.To(!d.Get("local_authentication_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if !features.SixPointOh() && !pluginsdk.IsExplicitlyNullInConfig(d, "is_virtual_network_filter_enabled") {
		account.Properties.IsVirtualNetworkFilterEnabled = pointer.To(d.Get("is_virtual_network_filter_enabled").(bool))
	}

	if v, ok := d.GetOk("default_identity_type"); ok {
		account.Properties.DefaultIdentity = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("analytical_storage"); ok {
		account.Properties.AnalyticalStorageConfiguration = expandCosmosDBAccountAnalyticalStorageConfiguration(v.([]any))
	}

	if v, ok := d.GetOk("capacity"); ok {
		account.Properties.Capacity = expandCosmosDBAccountCapacity(v.([]any))
	}

	var createMode string
	if v, ok := d.GetOk("create_mode"); ok {
		createMode = v.(string)
		account.Properties.CreateMode = pointer.ToEnum[cosmosdb.CreateMode](createMode)
	}

	if v, ok := d.GetOk("restore"); ok {
		account.Properties.RestoreParameters = expandCosmosdbAccountRestoreParameters(v.([]any))
	}

	if v, ok := d.GetOk("mongo_server_version"); ok {
		account.Properties.ApiProperties = &cosmosdb.ApiProperties{
			ServerVersion: pointer.ToEnum[cosmosdb.ServerVersion](v.(string)),
		}
	}

	if v, ok := d.GetOk("backup"); ok {
		policy, err := expandCosmosdbAccountBackup(v.([]any), false, createMode)
		if err != nil {
			return fmt.Errorf("expanding `backup`: %+v", err)
		}
		account.Properties.BackupPolicy = policy
	} else if createMode != "" {
		return errors.New("`create_mode` only works when `backup.type` is `Continuous`")
	}

	var key *keyvault.NestedItemID
	if v, ok := d.GetOk("key_vault_key_id"); ok {
		keyId, err := keyvault.ParseNestedItemID(v.(string), keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
		if err != nil {
			return err
		}
		key = keyId
	}

	if key != nil {
		account.Properties.KeyVaultKeyUri = pointer.To(key.ID())
	}

	// additional validation on MaxStalenessPrefix as it varies depending on if the DB is multi region or not
	consistencyPolicy := account.Properties.ConsistencyPolicy
	if len(geoLocations) > 1 && consistencyPolicy != nil && consistencyPolicy.DefaultConsistencyLevel == cosmosdb.DefaultConsistencyLevelBoundedStaleness {
		if msp := consistencyPolicy.MaxStalenessPrefix; msp != nil && pointer.From(msp) < 100000 {
			return fmt.Errorf("max_staleness_prefix (%d) must be greater than 100000 when more then one geo_location is used", *msp)
		}
		if mis := consistencyPolicy.MaxIntervalInSeconds; mis != nil && pointer.From(mis) < 300 {
			return fmt.Errorf("max_interval_in_seconds (%d) must be greater than 300 (5 minutes) when more then one geo_location is used", *mis)
		}
	}

	if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, id, account); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// NOTE: this is to work around the issue here: https://github.com/Azure/azure-rest-api-specs/issues/27596
	// Once the above issue is resolved we shouldn't need this check and update anymore
	if d.Get("create_mode").(string) == string(cosmosdb.CreateModeRestore) {
		if err = resourceCosmosDbAccountApiCreateOrUpdate(client, ctx, id, account); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceCosmosDbAccountRead(d, meta)
}
