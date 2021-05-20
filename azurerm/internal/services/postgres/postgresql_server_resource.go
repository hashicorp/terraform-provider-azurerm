package postgres

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2020-01-01/postgresql"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

func resourcePostgreSQLServer() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostgreSQLServerCreate,
		Read:   resourcePostgreSQLServerRead,
		Update: resourcePostgreSQLServerUpdate,
		Delete: resourcePostgreSQLServerDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				client := meta.(*clients.Client).Postgres.ServersClient
				ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
				defer cancel()

				id, err := parse.ServerID(d.Id())
				if err != nil {
					return []*schema.ResourceData{d}, err
				}

				resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
				if err != nil {
					return []*schema.ResourceData{d}, fmt.Errorf("reading PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
				}

				d.Set("create_mode", "Default")
				if resp.ReplicationRole != nil && *resp.ReplicationRole != "Master" && *resp.ReplicationRole != "None" {
					d.Set("create_mode", resp.ReplicationRole)

					masterServerId, err := parse.ServerID(*resp.MasterServerID)
					if err != nil {
						return []*schema.ResourceData{d}, fmt.Errorf("parsing Postgres Master Server ID : %v", err)
					}
					d.Set("creation_source_server_id", masterServerId.ID())
				}

				return []*schema.ResourceData{d}, nil
			},
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
				ValidateFunc: validate.ServerName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(skuList, false),
			},

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(postgresql.NineFullStopFive),
					string(postgresql.NineFullStopSix),
					string(postgresql.OneOne),
					string(postgresql.OneZero),
					string(postgresql.OneZeroFullStopZero),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference, // TODO: make case sensitive in 3.0
			},

			"storage_profile": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				Deprecated: "all storage_profile properties have been move to the top level. This block will be removed in version 3.0 of the provider.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_mb": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"storage_mb"},
							Deprecated:    "this has been moved to the top level and will be removed in version 3.0 of the provider.",
							ValidateFunc: validation.All(
								validation.IntBetween(5120, 4194304),
								validation.IntDivisibleBy(1024),
							),
						},

						"backup_retention_days": {
							Type:          schema.TypeInt,
							Optional:      true,
							Default:       7,
							ConflictsWith: []string{"backup_retention_days"},
							Deprecated:    "this has been moved to the top level and will be removed in version 3.0 of the provider.",
							ValidateFunc:  validation.IntBetween(7, 35),
						},

						"auto_grow": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"auto_grow_enabled"},
							Deprecated:    "this has been moved to the top level and will be removed in version 3.0 of the provider.",
							ValidateFunc: validation.StringInSlice([]string{
								string(postgresql.StorageAutogrowEnabled),
								string(postgresql.StorageAutogrowDisabled),
							}, false),
						},

						"geo_redundant_backup": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ForceNew:      true,
							ConflictsWith: []string{"geo_redundant_backup_enabled"},
							Deprecated:    "this has been moved to the top level and will be removed in version 3.0 of the provider.",
							ValidateFunc: validation.StringInSlice([]string{
								"Enabled",
								"Disabled",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"administrator_login": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace, validate.AdminUsernames),
			},

			"administrator_login_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"auto_grow_enabled": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true, // TODO: remove in 3.0 and default to true
				ConflictsWith: []string{"storage_profile", "storage_profile.0.auto_grow"},
			},

			"backup_retention_days": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"storage_profile", "storage_profile.0.backup_retention_days"},
				ValidateFunc:  validation.IntBetween(7, 35),
			},

			"geo_redundant_backup_enabled": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				Computed:      true, // TODO: remove in 2.0 and default to false
				ConflictsWith: []string{"storage_profile", "storage_profile.0.geo_redundant_backup"},
			},

			"create_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(postgresql.CreateModeDefault),
				ValidateFunc: validation.StringInSlice([]string{
					string(postgresql.CreateModeDefault),
					string(postgresql.CreateModeGeoRestore),
					string(postgresql.CreateModePointInTimeRestore),
					string(postgresql.CreateModeReplica),
				}, false),
			},

			"creation_source_server_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.ServerID,
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
								string(postgresql.SystemAssigned),
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

			"infrastructure_encryption_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"public_network_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"restore_point_in_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"storage_mb": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"storage_profile", "storage_profile.0.storage_mb"},
				ValidateFunc: validation.All(
					validation.IntBetween(5120, 16777216),
					validation.IntDivisibleBy(1024),
				),
			},

			"ssl_minimal_tls_version_enforced": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(postgresql.TLSEnforcementDisabled),
				ValidateFunc: validation.StringInSlice([]string{
					string(postgresql.TLSEnforcementDisabled),
					string(postgresql.TLS10),
					string(postgresql.TLS11),
					string(postgresql.TLS12),
				}, false),
			},

			"ssl_enforcement_enabled": {
				Type:         schema.TypeBool,
				Optional:     true, // required in 3.0
				ExactlyOneOf: []string{"ssl_enforcement", "ssl_enforcement_enabled"},
			},

			"ssl_enforcement": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "this has been renamed to the boolean `ssl_enforcement_enabled` and will be removed in version 3.0 of the provider.",
				ExactlyOneOf: []string{"ssl_enforcement", "ssl_enforcement_enabled"},
				ValidateFunc: validation.StringInSlice([]string{
					string(postgresql.SslEnforcementEnumDisabled),
					string(postgresql.SslEnforcementEnumEnabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"threat_detection_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							AtLeastOneOf: []string{"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"disabled_alerts": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"Sql_Injection",
									"Sql_Injection_Vulnerability",
									"Access_Anomaly",
									"Data_Exfiltration",
									"Unsafe_Action",
								}, false),
							},
							AtLeastOneOf: []string{"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"email_account_admins": {
							Type:     schema.TypeBool,
							Optional: true,
							AtLeastOneOf: []string{"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"email_addresses": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								// todo email validation in code
							},
							Set: schema.HashString,
							AtLeastOneOf: []string{"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"retention_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							AtLeastOneOf: []string{"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"storage_account_access_key": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"storage_endpoint": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},
					},
				},
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
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
		),
	}
}

func resourcePostgreSQLServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	securityClient := meta.(*clients.Client).Postgres.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_postgresql_server", *existing.ID)
	}

	mode := postgresql.CreateMode(d.Get("create_mode").(string))
	tlsMin := postgresql.MinimalTLSVersionEnum(d.Get("ssl_minimal_tls_version_enforced").(string))
	source := d.Get("creation_source_server_id").(string)
	version := postgresql.ServerVersion(d.Get("version").(string))

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name` for PostgreSQL Server %s (Resource Group %q): %v", name, resourceGroup, err)
	}

	infraEncrypt := postgresql.InfrastructureEncryptionEnabled
	if v := d.Get("infrastructure_encryption_enabled"); !v.(bool) {
		infraEncrypt = postgresql.InfrastructureEncryptionDisabled
	}

	publicAccess := postgresql.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		publicAccess = postgresql.PublicNetworkAccessEnumDisabled
	}

	ssl := postgresql.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled"); !v.(bool) {
		ssl = postgresql.SslEnforcementEnumDisabled
	}

	storage := expandPostgreSQLStorageProfile(d)

	var props postgresql.BasicServerPropertiesForCreate
	switch mode {
	case postgresql.CreateModeDefault:
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
		props = &postgresql.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         &admin,
			AdministratorLoginPassword: &pass,
			CreateMode:                 mode,
			InfrastructureEncryption:   infraEncrypt,
			PublicNetworkAccess:        publicAccess,
			MinimalTLSVersion:          tlsMin,
			SslEnforcement:             ssl,
			StorageProfile:             storage,
			Version:                    version,
		}
	case postgresql.CreateModePointInTimeRestore:
		v, ok := d.GetOk("restore_point_in_time")
		if !ok || v.(string) == "" {
			return fmt.Errorf("restore_point_in_time must be set when create_mode is PointInTimeRestore")
		}
		time, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema

		props = &postgresql.ServerPropertiesForRestore{
			CreateMode:     mode,
			SourceServerID: &source,
			RestorePointInTime: &date.Time{
				Time: time,
			},
			InfrastructureEncryption: infraEncrypt,
			PublicNetworkAccess:      publicAccess,
			MinimalTLSVersion:        tlsMin,
			SslEnforcement:           ssl,
			StorageProfile:           storage,
			Version:                  version,
		}
	case postgresql.CreateModeGeoRestore:
		props = &postgresql.ServerPropertiesForGeoRestore{
			CreateMode:               mode,
			SourceServerID:           &source,
			InfrastructureEncryption: infraEncrypt,
			PublicNetworkAccess:      publicAccess,
			MinimalTLSVersion:        tlsMin,
			SslEnforcement:           ssl,
			StorageProfile:           storage,
			Version:                  version,
		}
	case postgresql.CreateModeReplica:
		props = &postgresql.ServerPropertiesForReplica{
			CreateMode:               mode,
			SourceServerID:           &source,
			InfrastructureEncryption: infraEncrypt,
			PublicNetworkAccess:      publicAccess,
			MinimalTLSVersion:        tlsMin,
			SslEnforcement:           ssl,
			Version:                  version,
		}
	}

	server := postgresql.ServerForCreate{
		Identity:   expandServerIdentity(d.Get("identity").([]interface{})),
		Location:   &location,
		Properties: props,
		Sku:        sku,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Create(ctx, resourceGroup, name, server)
	if err != nil {
		return fmt.Errorf("creating PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for PostgreSQL Server %q (Resource Group %q) to become available", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(postgresql.ServerStateInaccessible)},
		Target:     []string{string(postgresql.ServerStateReady)},
		Refresh:    postgreSqlStateRefreshFunc(ctx, client, resourceGroup, name),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for PostgreSQL Server %q (Resource Group %q)to become available: %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL Server %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	if v, ok := d.GetOk("threat_detection_policy"); ok {
		alert := expandSecurityAlertPolicy(v)
		if alert != nil {
			future, err := securityClient.CreateOrUpdate(ctx, resourceGroup, name, *alert)
			if err != nil {
				return fmt.Errorf("error updataing postgres server security alert policy: %v", err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("error waiting for creation/update of postgrest server security alert policy (server %q, resource group %q): %+v", name, resourceGroup, err)
			}
		}
	}

	// Issue tracking the REST API update failure: https://github.com/Azure/azure-rest-api-specs/issues/14117
	if mode == postgresql.CreateModeReplica {
		log.Printf("[INFO] changing `public_network_access_enabled` for AzureRM PostgreSQL Server %q (Resource Group %q)", name, resourceGroup)
		properties := postgresql.ServerUpdateParameters{
			ServerUpdateParametersProperties: &postgresql.ServerUpdateParametersProperties{
				PublicNetworkAccess: publicAccess,
			},
		}

		future, err := client.Update(ctx, resourceGroup, name, properties)
		if err != nil {
			return fmt.Errorf("updating PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for update of PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return resourcePostgreSQLServerRead(d, meta)
}

func resourcePostgreSQLServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	securityClient := meta.(*clients.Client).Postgres.ServerSecurityAlertPoliciesClient
	replicasClient := meta.(*clients.Client).Postgres.ReplicasClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// TODO: support for Delta updates

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server update.")

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Postgres Server ID : %v", err)
	}

	// Locks for upscaling of replicas
	mode := postgresql.CreateMode(d.Get("create_mode").(string))
	primaryID := id.String()
	if mode == postgresql.CreateModeReplica {
		primaryID = d.Get("creation_source_server_id").(string)

		// Wait for possible restarts triggered by scaling primary (and its replicas)
		log.Printf("[DEBUG] Waiting for PostgreSQL Server %q (Resource Group %q) to become available", id.Name, id.ResourceGroup)
		stateConf := &resource.StateChangeConf{
			Pending:    []string{string(postgresql.ServerStateInaccessible), "Restarting"},
			Target:     []string{string(postgresql.ServerStateReady)},
			Refresh:    postgreSqlStateRefreshFunc(ctx, client, id.ResourceGroup, id.Name),
			MinTimeout: 15 * time.Second,
			Timeout:    d.Timeout(schema.TimeoutCreate),
		}

		if _, err = stateConf.WaitForState(); err != nil {
			return fmt.Errorf("waiting for PostgreSQL Server %q (Resource Group %q)to become available: %+v", id.Name, id.ResourceGroup, err)
		}
	}
	locks.ByID(primaryID)
	defer locks.UnlockByID(primaryID)

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name` for PostgreSQL Server %s (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	if d.HasChange("sku_name") && mode != postgresql.CreateModeReplica {
		oldRaw, newRaw := d.GetChange("sku_name")
		old := oldRaw.(string)
		new := newRaw.(string)

		if indexOfSku(old) < indexOfSku(new) {
			listReplicas, err := replicasClient.ListByServer(ctx, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("listing replicas for PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}

			propertiesReplica := postgresql.ServerUpdateParameters{
				Sku: sku,
			}
			for _, replica := range *listReplicas.Value {
				replicaId, err := parse.ServerID(*replica.ID)
				if err != nil {
					return fmt.Errorf("parsing Postgres Server Replica ID : %v", err)
				}
				future, err := client.Update(ctx, replicaId.ResourceGroup, replicaId.Name, propertiesReplica)
				if err != nil {
					return fmt.Errorf("upscaling PostgreSQL Server Replica %q (Resource Group %q): %+v", replicaId.Name, replicaId.ResourceGroup, err)
				}

				if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
					return fmt.Errorf("waiting for update of PostgreSQL Server Replica %q (Resource Group %q): %+v", replicaId.Name, replicaId.ResourceGroup, err)
				}
			}
		}
	}

	publicAccess := postgresql.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		publicAccess = postgresql.PublicNetworkAccessEnumDisabled
	}

	ssl := postgresql.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled"); !v.(bool) {
		ssl = postgresql.SslEnforcementEnumDisabled
	}

	tlsMin := postgresql.MinimalTLSVersionEnum(d.Get("ssl_minimal_tls_version_enforced").(string))

	properties := postgresql.ServerUpdateParameters{
		Identity: expandServerIdentity(d.Get("identity").([]interface{})),
		ServerUpdateParametersProperties: &postgresql.ServerUpdateParametersProperties{
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
			PublicNetworkAccess:        publicAccess,
			SslEnforcement:             ssl,
			MinimalTLSVersion:          tlsMin,
			StorageProfile:             expandPostgreSQLStorageProfile(d),
			Version:                    postgresql.ServerVersion(d.Get("version").(string)),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if old, new := d.GetChange("create_mode"); postgresql.CreateMode(old.(string)) == postgresql.CreateModeReplica && postgresql.CreateMode(new.(string)) == postgresql.CreateModeDefault {
		properties.ServerUpdateParametersProperties.ReplicationRole = utils.String("None")
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, properties)
	if err != nil {
		return fmt.Errorf("updating PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if v, ok := d.GetOk("threat_detection_policy"); ok {
		alert := expandSecurityAlertPolicy(v)
		if alert != nil {
			future, err := securityClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, *alert)
			if err != nil {
				return fmt.Errorf("error updataing mssql server security alert policy: %v", err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("error waiting for creation/update of postgrest server security alert policy (server %q, resource group %q): %+v", id.Name, id.ResourceGroup, err)
			}
		}
	}

	return resourcePostgreSQLServerRead(d, meta)
}

func resourcePostgreSQLServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	securityClient := meta.(*clients.Client).Postgres.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Postgres Server ID : %v", err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] PostgreSQL Server %q was not found (resource group %q)", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Azure PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	tier := postgresql.Basic
	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
		tier = sku.Tier
	}

	if err := d.Set("identity", flattenServerIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.ServerProperties; props != nil {
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("ssl_enforcement", string(props.SslEnforcement))
		d.Set("ssl_minimal_tls_version_enforced", props.MinimalTLSVersion)
		d.Set("version", string(props.Version))

		d.Set("infrastructure_encryption_enabled", props.InfrastructureEncryption == postgresql.InfrastructureEncryptionEnabled)
		d.Set("public_network_access_enabled", props.PublicNetworkAccess == postgresql.PublicNetworkAccessEnumEnabled)
		d.Set("ssl_enforcement_enabled", props.SslEnforcement == postgresql.SslEnforcementEnumEnabled)

		if err := d.Set("storage_profile", flattenPostgreSQLStorageProfile(props.StorageProfile)); err != nil {
			return fmt.Errorf("setting `storage_profile`: %+v", err)
		}

		if storage := props.StorageProfile; storage != nil {
			d.Set("storage_mb", storage.StorageMB)
			d.Set("backup_retention_days", storage.BackupRetentionDays)
			d.Set("auto_grow_enabled", storage.StorageAutogrow == postgresql.StorageAutogrowEnabled)
			d.Set("geo_redundant_backup_enabled", storage.GeoRedundantBackup == postgresql.Enabled)
		}

		// Computed
		d.Set("fqdn", props.FullyQualifiedDomainName)
	}

	// the basic does not support threat detection policies
	if tier == postgresql.GeneralPurpose || tier == postgresql.MemoryOptimized {
		secResp, err := securityClient.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil && !utils.ResponseWasNotFound(secResp.Response) {
			return fmt.Errorf("error making read request to postgres server security alert policy: %+v", err)
		}

		if !utils.ResponseWasNotFound(secResp.Response) {
			block := flattenSecurityAlertPolicy(secResp.SecurityAlertPolicyProperties, d.Get("threat_detection_policy.0.storage_account_access_key").(string))
			if err := d.Set("threat_detection_policy", block); err != nil {
				return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourcePostgreSQLServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Postgres Server ID : %v", err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("deleting PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("waiting for deletion of PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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

func expandServerSkuName(skuName string) (*postgresql.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 3 {
		return nil, fmt.Errorf("sku_name (%s) has the wrong number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var tier postgresql.SkuTier
	switch parts[0] {
	case "B":
		tier = postgresql.Basic
	case "GP":
		tier = postgresql.GeneralPurpose
	case "MO":
		tier = postgresql.MemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	capacity, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("cannot convert skuname %s capcity %s to int", skuName, parts[2])
	}

	return &postgresql.Sku{
		Name:     utils.String(skuName),
		Tier:     tier,
		Capacity: utils.Int32(int32(capacity)),
		Family:   utils.String(parts[1]),
	}, nil
}

func expandPostgreSQLStorageProfile(d *schema.ResourceData) *postgresql.StorageProfile {
	storage := postgresql.StorageProfile{}
	if v, ok := d.GetOk("storage_profile"); ok {
		storageprofile := v.([]interface{})[0].(map[string]interface{})

		storage.BackupRetentionDays = utils.Int32(int32(storageprofile["backup_retention_days"].(int)))
		storage.StorageMB = utils.Int32(int32(storageprofile["storage_mb"].(int)))
		storage.StorageAutogrow = postgresql.StorageAutogrow(storageprofile["auto_grow"].(string))
		storage.GeoRedundantBackup = postgresql.GeoRedundantBackup(storageprofile["geo_redundant_backup"].(string))
	}

	// now override whatever we may have from the block with the top level properties
	if v, ok := d.GetOk("auto_grow_enabled"); ok {
		storage.StorageAutogrow = postgresql.StorageAutogrowDisabled
		if v.(bool) {
			storage.StorageAutogrow = postgresql.StorageAutogrowEnabled
		}
	}

	if v, ok := d.GetOk("backup_retention_days"); ok {
		storage.BackupRetentionDays = utils.Int32(int32(v.(int)))
	}

	if v, ok := d.GetOk("geo_redundant_backup_enabled"); ok {
		storage.GeoRedundantBackup = postgresql.Disabled
		if v.(bool) {
			storage.GeoRedundantBackup = postgresql.Enabled
		}
	}

	if v, ok := d.GetOk("storage_mb"); ok {
		storage.StorageMB = utils.Int32(int32(v.(int)))
	}

	return &storage
}

func flattenPostgreSQLStorageProfile(resp *postgresql.StorageProfile) []interface{} {
	values := map[string]interface{}{}

	values["storage_mb"] = nil
	if storageMB := resp.StorageMB; storageMB != nil {
		values["storage_mb"] = *storageMB
	}

	values["backup_retention_days"] = nil
	if backupRetentionDays := resp.BackupRetentionDays; backupRetentionDays != nil {
		values["backup_retention_days"] = *backupRetentionDays
	}

	values["auto_grow"] = string(resp.StorageAutogrow)
	values["geo_redundant_backup"] = string(resp.GeoRedundantBackup)

	return []interface{}{values}
}

func expandSecurityAlertPolicy(i interface{}) *postgresql.ServerSecurityAlertPolicy {
	slice := i.([]interface{})
	if len(slice) == 0 {
		return nil
	}

	block := slice[0].(map[string]interface{})

	state := postgresql.ServerSecurityAlertPolicyStateEnabled
	if !block["enabled"].(bool) {
		state = postgresql.ServerSecurityAlertPolicyStateDisabled
	}

	props := &postgresql.SecurityAlertPolicyProperties{
		State: state,
	}

	if v, ok := block["disabled_alerts"]; ok {
		props.DisabledAlerts = utils.ExpandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := block["email_addresses"]; ok {
		props.EmailAddresses = utils.ExpandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := block["email_account_admins"]; ok {
		props.EmailAccountAdmins = utils.Bool(v.(bool))
	}

	if v, ok := block["retention_days"]; ok {
		props.RetentionDays = utils.Int32(int32(v.(int)))
	}

	if v, ok := block["storage_account_access_key"]; ok && v.(string) != "" {
		props.StorageAccountAccessKey = utils.String(v.(string))
	}

	if v, ok := block["storage_endpoint"]; ok && v.(string) != "" {
		props.StorageEndpoint = utils.String(v.(string))
	}

	return &postgresql.ServerSecurityAlertPolicy{
		SecurityAlertPolicyProperties: props,
	}
}

func flattenSecurityAlertPolicy(props *postgresql.SecurityAlertPolicyProperties, accessKey string) interface{} {
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
		props.State == postgresql.ServerSecurityAlertPolicyStateDisabled {
		return nil
	}

	block := map[string]interface{}{}

	block["enabled"] = props.State == postgresql.ServerSecurityAlertPolicyStateEnabled

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

func expandServerIdentity(input []interface{}) *postgresql.ResourceIdentity {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &postgresql.ResourceIdentity{
		Type: postgresql.IdentityType(v["type"].(string)),
	}
}

func flattenServerIdentity(input *postgresql.ResourceIdentity) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	principalID := ""
	if input.PrincipalID != nil {
		principalID = input.PrincipalID.String()
	}

	tenantID := ""
	if input.TenantID != nil {
		tenantID = input.TenantID.String()
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"principal_id": principalID,
			"tenant_id":    tenantID,
		},
	}
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

func postgreSqlStateRefreshFunc(ctx context.Context, client *postgresql.ServersClient, resourceGroup string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)
		if !utils.ResponseWasNotFound(res.Response) && err != nil {
			return nil, "", fmt.Errorf("retrieving status of PostgreSQL Server %s (Resource Group %q): %+v", name, resourceGroup, err)
		}

		// This is an issue with the RP, there is a 10 to 15 second lag before the
		// service will actually return the server
		if utils.ResponseWasNotFound(res.Response) {
			return res, string(postgresql.ServerStateInaccessible), nil
		}

		if res.ServerProperties != nil && res.ServerProperties.UserVisibleState != "" {
			return res, string(res.ServerProperties.UserVisibleState), nil
		}

		return res, string(postgresql.ServerStateInaccessible), nil
	}
}
