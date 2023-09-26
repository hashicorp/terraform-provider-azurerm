// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/replicas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/serversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	postgreSQLServerResourceName = "azurerm_postgresql_server"
)

var skuList = []string{
	"B_Gen4_1",
	"B_Gen4_2",
	"B_Gen5_1",
	"B_Gen5_2",
	"GP_Gen4_2",
	"GP_Gen4_4",
	"GP_Gen4_8",
	"GP_Gen4_16",
	"GP_Gen4_32",
	"GP_Gen5_2",
	"GP_Gen5_4",
	"GP_Gen5_8",
	"GP_Gen5_16",
	"GP_Gen5_32",
	"GP_Gen5_64",
	"MO_Gen5_2",
	"MO_Gen5_4",
	"MO_Gen5_8",
	"MO_Gen5_16",
	"MO_Gen5_32",
}

func resourcePostgreSQLServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgreSQLServerCreate,
		Read:   resourcePostgreSQLServerRead,
		Update: resourcePostgreSQLServerUpdate,
		Delete: resourcePostgreSQLServerDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := servers.ParseServerID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).Postgres.ServersClient

			id, err := servers.ParseServerID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, err
			}

			timeout, cancel := context.WithTimeout(ctx, d.Timeout(pluginsdk.TimeoutRead))
			defer cancel()

			resp, err := client.Get(timeout, *id)
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("reading %s: %+v", id, err)
			}

			d.Set("create_mode", "Default")
			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.ReplicationRole != nil && *props.ReplicationRole != "Master" && *props.ReplicationRole != "None" {
						d.Set("create_mode", props.ReplicationRole)

						sourceServerId, err := servers.ParseServerID(*props.MasterServerId)
						if err != nil {
							return []*pluginsdk.ResourceData{d}, fmt.Errorf("parsing Postgres Main Server ID : %v", err)
						}
						d.Set("creation_source_server_id", sourceServerId.ID())
					}
				}
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.PostgresqlServerV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(skuList, false),
			},

			"version": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(servers.PossibleValuesForServerVersion(), false),
			},

			"administrator_login": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace, validate.AdminUsernames),
			},

			"administrator_login_password": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"auto_grow_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"backup_retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(7, 35),
			},

			"geo_redundant_backup_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"create_mode": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(servers.CreateModeDefault),
				ValidateFunc: validation.StringInSlice(servers.PossibleValuesForCreateMode(), false),
			},

			"creation_source_server_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: servers.ValidateServerID,
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"infrastructure_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"restore_point_in_time": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"storage_mb": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.IntBetween(5120, 16777216),
					validation.IntDivisibleBy(1024),
				),
			},

			"ssl_minimal_tls_version_enforced": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(servers.MinimalTlsVersionEnumTLSOneTwo),
				ValidateFunc: validation.StringInSlice(servers.PossibleValuesForMinimalTlsVersionEnum(), false),
			},

			"ssl_enforcement_enabled": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"threat_detection_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"disabled_alerts": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Set:      pluginsdk.HashString,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"Sql_Injection",
									"Sql_Injection_Vulnerability",
									"Access_Anomaly",
									"Data_Exfiltration",
									"Unsafe_Action",
								}, false),
							},
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"email_account_admins": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"email_addresses": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								// todo email validation in code
							},
							Set: pluginsdk.HashString,
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"retention_days": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"storage_account_access_key": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"storage_endpoint": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},
					},
				},
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("sku_name", func(ctx context.Context, old, new, meta interface{}) bool {
				oldTier := strings.Split(old.(string), "_")
				newTier := strings.Split(new.(string), "_")
				// If the sku tier was not changed, we don't need ForceNew
				if oldTier[0] == newTier[0] {
					return false
				}
				// Basic tier could not be changed to other tiers
				if oldTier[0] == "B" || newTier[0] == "B" {
					return true
				}
				return false
			}),
			pluginsdk.ForceNewIfChange("create_mode", func(ctx context.Context, old, new, meta interface{}) bool {
				oldMode := servers.CreateMode(old.(string))
				newMode := servers.CreateMode(new.(string))
				// Instance could not be changed from Default to Replica
				if oldMode == servers.CreateModeDefault && newMode == servers.CreateModeReplica {
					return true
				}
				return false
			}),
		),
	}
}

func resourcePostgreSQLServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	securityClient := meta.(*clients.Client).Postgres.ServerSecurityAlertPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server creation.")

	id := servers.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_postgresql_server", id.ID())
	}

	mode := servers.CreateMode(d.Get("create_mode").(string))
	source := d.Get("creation_source_server_id").(string)
	version := servers.ServerVersion(d.Get("version").(string))

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name`: %+v", err)
	}

	infraEncrypt := servers.InfrastructureEncryptionEnabled
	if v := d.Get("infrastructure_encryption_enabled"); !v.(bool) {
		infraEncrypt = servers.InfrastructureEncryptionDisabled
	}

	publicAccess := servers.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		publicAccess = servers.PublicNetworkAccessEnumDisabled
	}

	ssl := servers.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled"); !v.(bool) {
		ssl = servers.SslEnforcementEnumDisabled
	}

	tlsMin := servers.MinimalTlsVersionEnum(d.Get("ssl_minimal_tls_version_enforced").(string))
	if ssl == servers.SslEnforcementEnumDisabled && tlsMin != servers.MinimalTlsVersionEnumTLSEnforcementDisabled {
		return fmt.Errorf("`ssl_minimal_tls_version_enforced` must be set to `TLSEnforcementDisabled` if `ssl_enforcement_enabled` is set to `false`")
	}

	storage := expandPostgreSQLStorageProfile(d)
	var props servers.ServerPropertiesForCreate
	switch mode {
	case servers.CreateModeDefault:
		admin := d.Get("administrator_login").(string)
		pass := d.Get("administrator_login_password").(string)
		if admin == "" {
			return fmt.Errorf("`administrator_login` must not be empty when `create_mode` is `default`")
		}
		if pass == "" {
			return fmt.Errorf("`administrator_login_password` must not be empty when `create_mode` is `default`")
		}

		if _, ok := d.GetOk("restore_point_in_time"); ok {
			return fmt.Errorf("`restore_point_in_time` cannot be set when `create_mode` is `default`")
		}

		// check admin
		props = servers.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         admin,
			AdministratorLoginPassword: pass,
			InfrastructureEncryption:   &infraEncrypt,
			PublicNetworkAccess:        &publicAccess,
			MinimalTlsVersion:          &tlsMin,
			SslEnforcement:             &ssl,
			StorageProfile:             storage,
			Version:                    &version,
		}
	case servers.CreateModePointInTimeRestore:
		v, ok := d.GetOk("restore_point_in_time")
		if !ok || v.(string) == "" {
			return fmt.Errorf("restore_point_in_time must be set when create_mode is PointInTimeRestore")
		}

		props = &servers.ServerPropertiesForRestore{
			SourceServerId:           source,
			RestorePointInTime:       v.(string),
			InfrastructureEncryption: &infraEncrypt,
			MinimalTlsVersion:        &tlsMin,
			SslEnforcement:           &ssl,
			StorageProfile:           storage,
			Version:                  &version,
		}
	case servers.CreateModeGeoRestore:
		props = &servers.ServerPropertiesForGeoRestore{
			SourceServerId:           source,
			InfrastructureEncryption: &infraEncrypt,
			PublicNetworkAccess:      &publicAccess,
			MinimalTlsVersion:        &tlsMin,
			SslEnforcement:           &ssl,
			StorageProfile:           storage,
			Version:                  &version,
		}
	case servers.CreateModeReplica:
		props = &servers.ServerPropertiesForReplica{
			SourceServerId:           source,
			InfrastructureEncryption: &infraEncrypt,
			PublicNetworkAccess:      &publicAccess,
			MinimalTlsVersion:        &tlsMin,
			SslEnforcement:           &ssl,
			Version:                  &version,
		}
	}

	expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	server := servers.ServerForCreate{
		Identity:   expandedIdentity,
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: props,
		Sku:        sku,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err = client.CreateThenPoll(ctx, id, server); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to become available", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(servers.ServerStateInaccessible)},
		Target:     []string{string(servers.ServerStateReady)},
		Refresh:    postgreSqlStateRefreshFunc(ctx, client, id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}

	d.SetId(id.ID())

	if v, ok := d.GetOk("threat_detection_policy"); ok {
		securityAlertId := serversecurityalertpolicies.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
		alert := expandSecurityAlertPolicy(v)
		if alert != nil {
			if err = securityClient.CreateOrUpdateThenPoll(ctx, securityAlertId, *alert); err != nil {
				return fmt.Errorf("updataing security alert policy for %s: %v", id, err)
			}
		}
	}

	// Issue tracking the REST API update failure: https://github.com/Azure/azure-rest-api-specs/issues/14117
	if mode == servers.CreateModeReplica {
		log.Printf("[INFO] updating `public_network_access_enabled` and `identity` for %s", id)
		properties := servers.ServerUpdateParameters{
			Identity: expandedIdentity,
			Properties: &servers.ServerUpdateParametersProperties{
				PublicNetworkAccess: &publicAccess,
			},
		}

		if err = client.UpdateThenPoll(ctx, id, properties); err != nil {
			return fmt.Errorf("updating Public Network Access for Replica %q: %+v", id, err)
		}
	}

	if mode == servers.CreateModePointInTimeRestore {
		log.Printf("[INFO] updating `public_network_access_enabled` for %s", id)
		properties := servers.ServerUpdateParameters{
			Properties: &servers.ServerUpdateParametersProperties{
				PublicNetworkAccess: &publicAccess,
			},
		}

		if err = client.UpdateThenPoll(ctx, id, properties); err != nil {
			return fmt.Errorf("updating Public Network Access for PointInTimeRestore %q: %+v", id, err)
		}
	}

	return resourcePostgreSQLServerRead(d, meta)
}

func resourcePostgreSQLServerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	securityClient := meta.(*clients.Client).Postgres.ServerSecurityAlertPoliciesClient
	replicasClient := meta.(*clients.Client).Postgres.ReplicasClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// TODO: support for Delta updates

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Postgres Server ID : %v", err)
	}

	// Locks for upscaling of replicas
	mode := servers.CreateMode(d.Get("create_mode").(string))
	primaryID := id.String()
	if mode == servers.CreateModeReplica {
		primaryID = d.Get("creation_source_server_id").(string)

		// Wait for possible restarts triggered by scaling primary (and its replicas)
		log.Printf("[DEBUG] Waiting for %s to become available", *id)
		stateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{string(servers.ServerStateInaccessible), "Restarting"},
			Target:     []string{string(servers.ServerStateReady)},
			Refresh:    postgreSqlStateRefreshFunc(ctx, client, *id),
			MinTimeout: 15 * time.Second,
			Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
		}

		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for %s to become available: %+v", *id, err)
		}
	}
	locks.ByID(primaryID)
	defer locks.UnlockByID(primaryID)

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name`: %v", err)
	}

	if d.HasChange("sku_name") && mode != servers.CreateModeReplica {
		oldRaw, newRaw := d.GetChange("sku_name")
		old := oldRaw.(string)
		new := newRaw.(string)

		if indexOfSku(old) < indexOfSku(new) {
			replicasId := replicas.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
			listReplicas, err := replicasClient.ListByServer(ctx, replicasId)
			if err != nil {
				return fmt.Errorf("listing replicas for %s: %+v", *id, err)
			}

			propertiesReplica := servers.ServerUpdateParameters{
				Sku: sku,
			}

			if listReplicas.Model != nil && listReplicas.Model.Value != nil {
				replicaList := *listReplicas.Model.Value
				for _, replica := range replicaList {
					replicaId, err := servers.ParseServerID(*replica.Id)
					if err != nil {
						return fmt.Errorf("parsing Postgres Server Replica ID : %v", err)
					}
					if err = client.UpdateThenPoll(ctx, *replicaId, propertiesReplica); err != nil {
						return fmt.Errorf("updating SKU for Replica %s: %+v", *replicaId, err)
					}
				}
			}
		}
	}

	ssl := servers.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled"); !v.(bool) {
		ssl = servers.SslEnforcementEnumDisabled
	}

	tlsMin := servers.MinimalTlsVersionEnum(d.Get("ssl_minimal_tls_version_enforced").(string))

	if ssl == servers.SslEnforcementEnumDisabled && tlsMin != servers.MinimalTlsVersionEnumTLSEnforcementDisabled {
		return fmt.Errorf("`ssl_minimal_tls_version_enforced` must be set to `TLSEnforcementDisabled` if `ssl_enforcement_enabled` is set to `false`")
	}

	expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	serverVersion := servers.ServerVersion(d.Get("version").(string))
	properties := servers.ServerUpdateParameters{
		Identity: expandedIdentity,
		Properties: &servers.ServerUpdateParametersProperties{
			SslEnforcement:    &ssl,
			MinimalTlsVersion: &tlsMin,
			StorageProfile:    expandPostgreSQLStorageProfile(d),
			Version:           &serverVersion,
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	publicAccess := servers.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		publicAccess = servers.PublicNetworkAccessEnumDisabled
	}
	properties.Properties.PublicNetworkAccess = &publicAccess

	oldCreateMode, newCreateMode := d.GetChange("create_mode")
	replicaUpdatedToDefault := servers.CreateMode(oldCreateMode.(string)) == servers.CreateModeReplica && servers.CreateMode(newCreateMode.(string)) == servers.CreateModeDefault
	if replicaUpdatedToDefault {
		properties.Properties.ReplicationRole = utils.String("None")
	}

	// Update Admin Password in the separate call when Replication is stopped: https://github.com/Azure/azure-rest-api-specs/issues/16898
	if d.HasChange("administrator_login_password") && !replicaUpdatedToDefault {
		properties.Properties.AdministratorLoginPassword = utils.String(d.Get("administrator_login_password").(string))
	}

	if err = client.UpdateThenPoll(ctx, *id, properties); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	// Update Admin Password in a separate call when Replication is stopped: https://github.com/Azure/azure-rest-api-specs/issues/16898
	if d.HasChange("administrator_login_password") && replicaUpdatedToDefault {
		properties.Properties.AdministratorLoginPassword = utils.String(d.Get("administrator_login_password").(string))

		if err = client.UpdateThenPoll(ctx, *id, properties); err != nil {
			return fmt.Errorf("updating Admin Password of %q: %+v", id, err)
		}
	}

	if v, ok := d.GetOk("threat_detection_policy"); ok {
		alert := expandSecurityAlertPolicy(v)
		securityId := serversecurityalertpolicies.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
		if alert != nil {
			if err = securityClient.CreateOrUpdateThenPoll(ctx, securityId, *alert); err != nil {
				return fmt.Errorf("updating security alert policy for %s: %+v", *id, err)
			}
		}
	}

	return resourcePostgreSQLServerRead(d, meta)
}

func resourcePostgreSQLServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	securityClient := meta.(*clients.Client).Postgres.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Postgres Server ID : %v", err)
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))

		tier := servers.SkuTierBasic
		if sku := model.Sku; sku != nil {
			d.Set("sku_name", sku.Name)
			if sku.Tier != nil {
				tier = *sku.Tier
			}
		}

		if err := d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("administrator_login", props.AdministratorLogin)
			d.Set("ssl_minimal_tls_version_enforced", string(pointer.From(props.MinimalTlsVersion)))

			version := ""
			if props.Version != nil {
				version = string(*props.Version)
			}
			d.Set("version", version)

			infrastructureEncryption := false
			if props.InfrastructureEncryption != nil {
				infrastructureEncryption = *props.InfrastructureEncryption == servers.InfrastructureEncryptionEnabled
			}
			d.Set("infrastructure_encryption_enabled", infrastructureEncryption)

			publicNetworkAccess := false
			if props.PublicNetworkAccess != nil {
				publicNetworkAccess = *props.PublicNetworkAccess == servers.PublicNetworkAccessEnumEnabled
			}
			d.Set("public_network_access_enabled", publicNetworkAccess)

			sslEnforcement := false
			if props.SslEnforcement != nil {
				sslEnforcement = *props.SslEnforcement == servers.SslEnforcementEnumEnabled
			}
			d.Set("ssl_enforcement_enabled", sslEnforcement)

			if storage := props.StorageProfile; storage != nil {
				d.Set("storage_mb", storage.StorageMB)
				d.Set("backup_retention_days", storage.BackupRetentionDays)

				autoGrow := false
				if storage.StorageAutogrow != nil {
					autoGrow = *storage.StorageAutogrow == servers.StorageAutogrowEnabled
				}
				d.Set("auto_grow_enabled", autoGrow)

				geoRedundant := false
				if storage.GeoRedundantBackup != nil {
					geoRedundant = *storage.GeoRedundantBackup == servers.GeoRedundantBackupEnabled
				}
				d.Set("geo_redundant_backup_enabled", geoRedundant)
			}

			// Computed
			d.Set("fqdn", props.FullyQualifiedDomainName)
		}

		// the basic does not support threat detection policies
		if tier == servers.SkuTierGeneralPurpose || tier == servers.SkuTierMemoryOptimized {
			securityId := serversecurityalertpolicies.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
			secResp, err := securityClient.Get(ctx, securityId)
			if err != nil && !response.WasNotFound(secResp.HttpResponse) {
				return fmt.Errorf("making read request to postgres server security alert policy: %+v", err)
			}

			if !response.WasNotFound(secResp.HttpResponse) {
				if secResp.Model != nil {
					block := flattenSecurityAlertPolicy(secResp.Model.Properties, d.Get("threat_detection_policy.0.storage_account_access_key").(string))
					if err := d.Set("threat_detection_policy", block); err != nil {
						return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
					}
				}
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourcePostgreSQLServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func indexOfSku(skuName string) int {
	for k, v := range skuList {
		if skuName == v {
			return k
		}
	}
	return -1 // not found.
}

func expandServerSkuName(skuName string) (*servers.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 3 {
		return nil, fmt.Errorf("sku_name (%s) has the wrong number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var tier servers.SkuTier
	switch parts[0] {
	case "B":
		tier = servers.SkuTierBasic
	case "GP":
		tier = servers.SkuTierGeneralPurpose
	case "MO":
		tier = servers.SkuTierMemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	capacity, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("cannot convert skuname %s capcity %s to int", skuName, parts[2])
	}

	return &servers.Sku{
		Name:     skuName,
		Tier:     &tier,
		Capacity: utils.Int64(int64(capacity)),
		Family:   utils.String(parts[1]),
	}, nil
}

func expandPostgreSQLStorageProfile(d *pluginsdk.ResourceData) *servers.StorageProfile {
	storage := servers.StorageProfile{}

	// now override whatever we may have from the block with the top level properties
	if v, ok := d.GetOk("auto_grow_enabled"); ok {
		autogrowEnabled := servers.StorageAutogrowDisabled
		if v.(bool) {
			autogrowEnabled = servers.StorageAutogrowEnabled
		}
		storage.StorageAutogrow = &autogrowEnabled
	}

	if v, ok := d.GetOk("backup_retention_days"); ok {
		storage.BackupRetentionDays = utils.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("geo_redundant_backup_enabled"); ok {
		geoRedundantBackup := servers.GeoRedundantBackupDisabled

		if v.(bool) {
			geoRedundantBackup = servers.GeoRedundantBackupEnabled
		}
		storage.GeoRedundantBackup = &geoRedundantBackup
	}

	if v, ok := d.GetOk("storage_mb"); ok {
		storage.StorageMB = utils.Int64(int64(v.(int)))
	}

	return &storage
}

func expandSecurityAlertPolicy(i interface{}) *serversecurityalertpolicies.ServerSecurityAlertPolicy {
	slice := i.([]interface{})
	if len(slice) == 0 {
		return nil
	}

	block := slice[0].(map[string]interface{})

	state := serversecurityalertpolicies.ServerSecurityAlertPolicyStateEnabled
	if !block["enabled"].(bool) {
		state = serversecurityalertpolicies.ServerSecurityAlertPolicyStateDisabled
	}

	props := &serversecurityalertpolicies.SecurityAlertPolicyProperties{
		State: state,
	}

	if v, ok := block["disabled_alerts"]; ok {
		props.DisabledAlerts = utils.ExpandStringSlice(v.(*pluginsdk.Set).List())
	}

	if v, ok := block["email_addresses"]; ok {
		props.EmailAddresses = utils.ExpandStringSlice(v.(*pluginsdk.Set).List())
	}

	if v, ok := block["email_account_admins"]; ok {
		props.EmailAccountAdmins = utils.Bool(v.(bool))
	}

	if v, ok := block["retention_days"]; ok {
		props.RetentionDays = utils.Int64(int64(v.(int)))
	}

	if v, ok := block["storage_account_access_key"]; ok && v.(string) != "" {
		props.StorageAccountAccessKey = utils.String(v.(string))
	}

	if v, ok := block["storage_endpoint"]; ok && v.(string) != "" {
		props.StorageEndpoint = utils.String(v.(string))
	}

	return &serversecurityalertpolicies.ServerSecurityAlertPolicy{
		Properties: props,
	}
}

func flattenSecurityAlertPolicy(props *serversecurityalertpolicies.SecurityAlertPolicyProperties, accessKey string) interface{} {
	if props == nil {
		return nil
	}

	// check if its an empty block as in its never been set before
	if props.DisabledAlerts != nil && len(*props.DisabledAlerts) == 1 && (*props.DisabledAlerts)[0] == "" &&
		props.EmailAddresses != nil && len(*props.EmailAddresses) == 1 && (*props.EmailAddresses)[0] == "" &&
		props.StorageAccountAccessKey != nil && *props.StorageAccountAccessKey == "" &&
		props.StorageEndpoint != nil && *props.StorageEndpoint == "" &&
		props.RetentionDays != nil && *props.RetentionDays == 0 &&
		props.EmailAccountAdmins != nil && !*props.EmailAccountAdmins &&
		props.State == serversecurityalertpolicies.ServerSecurityAlertPolicyStateDisabled {
		return nil
	}

	block := map[string]interface{}{}

	block["enabled"] = props.State == serversecurityalertpolicies.ServerSecurityAlertPolicyStateEnabled

	block["disabled_alerts"] = flattenSecurityAlertPolicySet(props.DisabledAlerts)
	block["email_addresses"] = flattenSecurityAlertPolicySet(props.EmailAddresses)

	if v := props.EmailAccountAdmins; v != nil {
		block["email_account_admins"] = *v
	}
	if v := props.RetentionDays; v != nil {
		block["retention_days"] = *v
	}
	if v := props.StorageEndpoint; v != nil {
		block["storage_endpoint"] = *v
	}

	block["storage_account_access_key"] = accessKey

	return []interface{}{block}
}

func flattenSecurityAlertPolicySet(input *[]string) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	// When empty, `disabledAlerts` and `emailAddresses` are returned as `[""]` by the api. We'll catch that here and return
	// an empty interface to set.
	attr := *input
	if len(attr) == 1 && attr[0] == "" {
		return make([]interface{}, 0)
	}

	return utils.FlattenStringSlice(input)
}

func postgreSqlStateRefreshFunc(ctx context.Context, client *servers.ServersClient, id servers.ServerId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if !response.WasNotFound(res.HttpResponse) && err != nil {
			return nil, "", fmt.Errorf("retrieving status of %s: %+v", id, err)
		}

		// For Replica servers, with enabled BYOK, state would be reported as 'Inaccessible', even when deployment was in 'Succeeded' state.
		// It is caused by a need to revalidate the key.
		if res.Model != nil && res.Model.Properties != nil &&
			res.Model.Properties.ReplicationRole != nil && *res.Model.Properties.ReplicationRole == "Replica" &&
			res.Model.Properties.ByokEnforcement != nil && *res.Model.Properties.ByokEnforcement == "Enabled" &&
			res.Model.Properties.UserVisibleState != nil && *res.Model.Properties.UserVisibleState == servers.ServerStateInaccessible {
			return res, string(servers.ServerStateReady), nil
		}

		// This is an issue with the RP, there is a 10 to 15 second lag before the
		// service will actually return the server
		if response.WasNotFound(res.HttpResponse) {
			return res, string(servers.ServerStateInaccessible), nil
		}

		if res.Model != nil && res.Model.Properties != nil && res.Model.Properties.UserVisibleState != nil && *res.Model.Properties.UserVisibleState != "" {
			return res, string(*res.Model.Properties.UserVisibleState), nil
		}

		return res, string(servers.ServerStateInaccessible), nil
	}
}
