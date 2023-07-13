// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSqlDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSqlDatabaseCreateUpdate,
		Read:   resourceSqlDatabaseRead,
		Update: resourceSqlDatabaseCreateUpdate,
		Delete: resourceSqlDatabaseDelete,

		DeprecationMessage: "The `azurerm_sql_database` resource is deprecated and will be removed in version 4.0 of the AzureRM provider. Please use the `azurerm_mssql_database` resource instead.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DatabaseID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: resourceSqlDatabaseSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
			threatDetection, hasThreatDetection := diff.GetOk("threat_detection_policy")
			if hasThreatDetection {
				if tl := threatDetection.([]interface{}); len(tl) > 0 {
					t := tl[0].(map[string]interface{})

					state := strings.ToLower(t["state"].(string))
					_, hasStorageEndpoint := t["storage_endpoint"]
					_, hasStorageAccountAccessKey := t["storage_account_access_key"]
					if state == "enabled" && !hasStorageEndpoint && !hasStorageAccountAccessKey {
						return fmt.Errorf("`storage_endpoint` and `storage_account_access_key` are required when `state` is `Enabled`")
					}
				}
			}

			return nil
		}),
	}
}

func resourceSqlDatabaseCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.DatabasesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	threatClient := meta.(*clients.Client).Sql.DatabaseThreatDetectionPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDatabaseID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_sql_database", id.ID())
		}
	}

	createMode := sql.CreateMode(d.Get("create_mode").(string))

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	properties := sql.Database{
		Location: utils.String(location),
		DatabaseProperties: &sql.DatabaseProperties{
			CreateMode:    createMode,
			ZoneRedundant: utils.Bool(d.Get("zone_redundant").(bool)),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("source_database_id"); ok {
		sourceDatabaseID := v.(string)
		properties.DatabaseProperties.SourceDatabaseID = utils.String(sourceDatabaseID)
	}

	if v, ok := d.GetOk("edition"); ok {
		edition := v.(string)
		properties.DatabaseProperties.Edition = sql.DatabaseEdition(edition)
	}

	if v, ok := d.GetOk("collation"); ok {
		collation := v.(string)
		properties.DatabaseProperties.Collation = utils.String(collation)
	}

	if v, ok := d.GetOk("max_size_bytes"); ok {
		maxSizeBytes := v.(string)
		properties.DatabaseProperties.MaxSizeBytes = utils.String(maxSizeBytes)
	}

	if v, ok := d.GetOk("source_database_deletion_date"); ok {
		sourceDatabaseDeletionString := v.(string)
		sourceDatabaseDeletionDate, err2 := date.ParseTime(time.RFC3339, sourceDatabaseDeletionString)
		if err2 != nil {
			return fmt.Errorf("`source_database_deletion_date` wasn't a valid RFC3339 date %q: %+v", sourceDatabaseDeletionString, err2)
		}

		properties.DatabaseProperties.SourceDatabaseDeletionDate = &date.Time{
			Time: sourceDatabaseDeletionDate,
		}
	}

	if v, ok := d.GetOk("requested_service_objective_id"); ok {
		requestedServiceObjectiveID := v.(string)
		id, err2 := uuid.FromString(requestedServiceObjectiveID)
		if err2 != nil {
			return fmt.Errorf("`requested_service_objective_id` wasn't a valid UUID %q: %+v", requestedServiceObjectiveID, err2)
		}
		properties.DatabaseProperties.RequestedServiceObjectiveID = &id
	}

	if v, ok := d.GetOk("elastic_pool_name"); ok {
		elasticPoolName := v.(string)
		properties.DatabaseProperties.ElasticPoolName = utils.String(elasticPoolName)
	}

	if v, ok := d.GetOk("requested_service_objective_name"); ok {
		requestedServiceObjectiveName := v.(string)
		properties.DatabaseProperties.RequestedServiceObjectiveName = sql.ServiceObjectiveName(requestedServiceObjectiveName)
	}

	if v, ok := d.GetOk("restore_point_in_time"); ok {
		restorePointInTime := v.(string)
		restorePointInTimeDate, err2 := date.ParseTime(time.RFC3339, restorePointInTime)
		if err2 != nil {
			return fmt.Errorf("`restore_point_in_time` wasn't a valid RFC3339 date %q: %+v", restorePointInTime, err2)
		}

		properties.DatabaseProperties.RestorePointInTime = &date.Time{
			Time: restorePointInTimeDate,
		}
	}

	readScale := d.Get("read_scale").(bool)
	if readScale {
		properties.DatabaseProperties.ReadScale = sql.ReadScaleEnabled
	} else {
		properties.DatabaseProperties.ReadScale = sql.ReadScaleDisabled
	}

	// The requested Service Objective Name does not match the requested Service Objective Id.
	if d.HasChange("requested_service_objective_name") && !d.HasChange("requested_service_objective_id") {
		properties.DatabaseProperties.RequestedServiceObjectiveID = nil
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, properties)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %+v", id, err)
	}

	if _, ok := d.GetOk("import"); ok {
		if createMode != sql.CreateModeDefault {
			return fmt.Errorf("import can only be used when create_mode is Default")
		}
		importParameters := expandAzureRmSqlDatabaseImport(d)
		importFuture, err := client.CreateImportOperation(ctx, id.ResourceGroup, id.ServerName, id.Name, importParameters)
		if err != nil {
			return err
		}

		if err = importFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return err
		}
	}

	d.SetId(id.ID())
	if _, err = threatClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, *expandArmSqlServerThreatDetectionPolicy(d, location)); err != nil {
		return fmt.Errorf("setting database threat detection policy: %+v", err)
	}

	return resourceSqlDatabaseRead(d, meta)
}

func resourceSqlDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	threatClient := meta.(*clients.Client).Sql.DatabaseThreatDetectionPoliciesClient
	threat, err := threatClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err == nil {
		if err := d.Set("threat_detection_policy", flattenArmSqlServerThreatDetectionPolicy(d, threat)); err != nil {
			return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
		}
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("server_name", id.ServerName)

	if props := resp.DatabaseProperties; props != nil {
		// TODO: set `create_mode` & `source_database_id` once this issue is fixed:
		// https://github.com/Azure/azure-rest-api-specs/issues/1604

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
		d.Set("read_scale", props.ReadScale == sql.ReadScaleEnabled)
		d.Set("zone_redundant", props.ZoneRedundant)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSqlDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting SQL Database: %+v", err)
	}

	return nil
}

func flattenEncryptionStatus(encryption *[]sql.TransparentDataEncryption) string {
	if encryption != nil {
		encrypted := *encryption
		if len(encrypted) > 0 {
			if props := encrypted[0].TransparentDataEncryptionProperties; props != nil {
				return string(props.Status)
			}
		}
	}

	return ""
}

func flattenArmSqlServerThreatDetectionPolicy(d *pluginsdk.ResourceData, policy sql.DatabaseSecurityAlertPolicy) []interface{} {
	// The SQL database threat detection API always returns the default value even if never set.
	// If the values are on their default one, threat it as not set.
	properties := policy.DatabaseSecurityAlertPolicyProperties
	if properties == nil {
		return []interface{}{}
	}

	threatDetectionPolicy := make(map[string]interface{})

	threatDetectionPolicy["state"] = string(properties.State)
	threatDetectionPolicy["email_account_admins"] = string(properties.EmailAccountAdmins)

	if disabledAlerts := properties.DisabledAlerts; disabledAlerts != nil {
		flattenedAlerts := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
		if v := *disabledAlerts; v != "" {
			parsedAlerts := strings.Split(v, ";")
			for _, a := range parsedAlerts {
				flattenedAlerts.Add(a)
			}
		}
		threatDetectionPolicy["disabled_alerts"] = flattenedAlerts
	}
	if emailAddresses := properties.EmailAddresses; emailAddresses != nil {
		flattenedEmails := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
		if v := *emailAddresses; v != "" {
			parsedEmails := strings.Split(*emailAddresses, ";")
			for _, e := range parsedEmails {
				flattenedEmails.Add(e)
			}
		}
		threatDetectionPolicy["email_addresses"] = flattenedEmails
	}
	if properties.StorageEndpoint != nil {
		threatDetectionPolicy["storage_endpoint"] = *properties.StorageEndpoint
	}
	if properties.RetentionDays != nil {
		threatDetectionPolicy["retention_days"] = int(*properties.RetentionDays)
	}

	// If storage account access key is in state read it to the new state, as the API does not return it for security reasons
	if v, ok := d.GetOk("threat_detection_policy.0.storage_account_access_key"); ok {
		threatDetectionPolicy["storage_account_access_key"] = v.(string)
	}

	return []interface{}{threatDetectionPolicy}
}

func expandAzureRmSqlDatabaseImport(d *pluginsdk.ResourceData) sql.ImportExtensionRequest {
	v := d.Get("import")
	dbimportRefs := v.([]interface{})
	dbimportRef := dbimportRefs[0].(map[string]interface{})
	return sql.ImportExtensionRequest{
		Name: utils.String("terraform"),
		ImportExtensionProperties: &sql.ImportExtensionProperties{
			StorageKeyType:             sql.StorageKeyType(dbimportRef["storage_key_type"].(string)),
			StorageKey:                 utils.String(dbimportRef["storage_key"].(string)),
			StorageURI:                 utils.String(dbimportRef["storage_uri"].(string)),
			AdministratorLogin:         utils.String(dbimportRef["administrator_login"].(string)),
			AdministratorLoginPassword: utils.String(dbimportRef["administrator_login_password"].(string)),
			AuthenticationType:         sql.AuthenticationType(dbimportRef["authentication_type"].(string)),
			OperationMode:              utils.String(dbimportRef["operation_mode"].(string)),
		},
	}
}

func expandArmSqlServerThreatDetectionPolicy(d *pluginsdk.ResourceData, location string) *sql.DatabaseSecurityAlertPolicy {
	policy := sql.DatabaseSecurityAlertPolicy{
		Location: utils.String(location),
		DatabaseSecurityAlertPolicyProperties: &sql.DatabaseSecurityAlertPolicyProperties{
			State: sql.SecurityAlertPolicyStateDisabled,
		},
	}
	properties := policy.DatabaseSecurityAlertPolicyProperties

	td, ok := d.GetOk("threat_detection_policy")
	if !ok {
		return &policy
	}

	if tdl := td.([]interface{}); len(tdl) > 0 {
		threatDetection := tdl[0].(map[string]interface{})

		properties.State = sql.SecurityAlertPolicyState(threatDetection["state"].(string))
		properties.EmailAccountAdmins = sql.SecurityAlertPolicyEmailAccountAdmins(threatDetection["email_account_admins"].(string))

		if v, ok := threatDetection["disabled_alerts"]; ok {
			alerts := v.(*pluginsdk.Set).List()
			expandedAlerts := make([]string, len(alerts))
			for i, a := range alerts {
				expandedAlerts[i] = a.(string)
			}
			properties.DisabledAlerts = utils.String(strings.Join(expandedAlerts, ";"))
		}
		if v, ok := threatDetection["email_addresses"]; ok {
			emails := v.(*pluginsdk.Set).List()
			expandedEmails := make([]string, len(emails))
			for i, e := range emails {
				expandedEmails[i] = e.(string)
			}
			properties.EmailAddresses = utils.String(strings.Join(expandedEmails, ";"))
		}
		if v, ok := threatDetection["retention_days"]; ok {
			properties.RetentionDays = utils.Int32(int32(v.(int)))
		}
		if v, ok := threatDetection["storage_account_access_key"]; ok {
			properties.StorageAccountAccessKey = utils.String(v.(string))
		}
		if v, ok := threatDetection["storage_endpoint"]; ok {
			properties.StorageEndpoint = utils.String(v.(string))
		}

		return &policy
	}

	return &policy
}

func resourceSqlDatabaseSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlDatabaseName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"server_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlServerName,
		},

		"create_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(sql.Default),
			ValidateFunc: validation.StringInSlice([]string{
				string(sql.CreateModeCopy),
				string(sql.CreateModeDefault),
				string(sql.CreateModeNonReadableSecondary),
				string(sql.CreateModeOnlineSecondary),
				string(sql.CreateModePointInTimeRestore),
				string(sql.CreateModeRecovery),
				string(sql.CreateModeRestore),
				string(sql.CreateModeRestoreLongTermRetentionBackup),
			}, false),
		},

		"import": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"storage_uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"storage_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"storage_key_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"StorageAccessKey",
							"SharedAccessKey",
						}, false),
					},
					"administrator_login": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"administrator_login_password": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"authentication_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"ADPassword",
							"SQL",
						}, false),
					},
					"operation_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "Import",
						ValidateFunc: validation.StringInSlice([]string{
							"Import",
						}, false),
					},
				},
			},
		},

		"source_database_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"restore_point_in_time": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"edition": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sql.Basic),
				string(sql.Business),
				string(sql.BusinessCritical),
				string(sql.DataWarehouse),
				string(sql.Free),
				string(sql.GeneralPurpose),
				string(sql.Hyperscale),
				string(sql.Premium),
				string(sql.PremiumRS),
				string(sql.Standard),
				string(sql.Stretch),
				string(sql.System),
				string(sql.System2),
				string(sql.Web),
			}, false),
		},

		"collation": {
			Type:             pluginsdk.TypeString,
			DiffSuppressFunc: suppress.CaseDifference,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
		},

		"max_size_bytes": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"max_size_gb": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"requested_service_objective_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IsUUID,
		},

		"requested_service_objective_name": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validation.StringIsNotEmpty,
			// TODO: add validation once the Enum's complete
			// https://github.com/Azure/azure-rest-api-specs/issues/1609
		},

		"source_database_deletion_date": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"elastic_pool_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"encryption": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"creation_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"default_secondary_location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"threat_detection_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
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
							}, false),
						},
					},

					"email_account_admins": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(sql.SecurityAlertPolicyEmailAccountAdminsDisabled),
						ValidateFunc: validation.StringInSlice([]string{
							string(sql.SecurityAlertPolicyEmailAccountAdminsDisabled),
							string(sql.SecurityAlertPolicyEmailAccountAdminsEnabled),
						}, false),
					},

					"email_addresses": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},

					"retention_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"state": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(sql.SecurityAlertPolicyStateDisabled),
						ValidateFunc: validation.StringInSlice([]string{
							string(sql.SecurityAlertPolicyStateDisabled),
							string(sql.SecurityAlertPolicyStateEnabled),
							string(sql.SecurityAlertPolicyStateNew),
						}, false),
					},

					"storage_account_access_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"storage_endpoint": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"read_scale": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"zone_redundant": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tags": tags.Schema(),
	}
}
