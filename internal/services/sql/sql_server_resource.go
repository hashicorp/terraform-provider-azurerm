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
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSqlServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSqlServerCreateUpdate,
		Read:   resourceSqlServerRead,
		Update: resourceSqlServerCreateUpdate,
		Delete: resourceSqlServerDelete,

		DeprecationMessage: "The `azurerm_sql_server` resource is deprecated and will be removed in version 4.0 of the AzureRM provider. Please use the `azurerm_mssql_server` resource instead.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ServerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlServerName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"version": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"2.0",
					"12.0",
				}, false),
			},

			"administrator_login": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"administrator_login_password": {
				Type:      pluginsdk.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"connection_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(sql.ServerConnectionTypeDefault),
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.ServerConnectionTypeDefault),
					string(sql.ServerConnectionTypeProxy),
					string(sql.ServerConnectionTypeRedirect),
				}, false),
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"fully_qualified_domain_name": {
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
									"Data_Exfiltration",
									"Unsafe_Action",
								}, false),
							},
						},

						"email_account_admins": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Computed: true,
						},

						"email_addresses": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Computed: true,
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
								string(sql.SecurityAlertPolicyStateNew), // Only kept for backward compatibility - TODO investigate if we can remove this in 4.0
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

			"tags": tags.Schema(),
		},

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

func resourceSqlServerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	connectionClient := meta.(*clients.Client).Sql.ServerConnectionPoliciesClient
	secPolicyClient := meta.(*clients.Client).Sql.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_sql_server", id.ID())
		}
	}

	adminUsername := d.Get("administrator_login").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	tags := tags.Expand(d.Get("tags").(map[string]interface{}))
	version := d.Get("version").(string)

	parameters := sql.Server{
		Location: utils.String(location),
		Tags:     tags,
		ServerProperties: &sql.ServerProperties{
			Version:            utils.String(version),
			AdministratorLogin: utils.String(adminUsername),
		},
	}

	if _, ok := d.GetOk("identity"); ok {
		sqlServerIdentity, err := expandAzureRmSqlServerIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		parameters.Identity = sqlServerIdentity
	}

	if d.HasChange("administrator_login_password") {
		adminPassword := d.Get("administrator_login_password").(string)
		parameters.ServerProperties.AdministratorLoginPassword = utils.String(adminPassword)
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasConflict(future.Response()) {
			return fmt.Errorf("SQL Server names need to be globally unique and %q is already in use.", id.Name)
		}

		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	connection := sql.ServerConnectionPolicy{
		ServerConnectionPolicyProperties: &sql.ServerConnectionPolicyProperties{
			ConnectionType: sql.ServerConnectionType(d.Get("connection_policy").(string)),
		},
	}
	if _, err = connectionClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, connection); err != nil {
		return fmt.Errorf("creating/updating Connection Policy for %s: %+v", id, err)
	}

	policyInput := expandSqlServerThreatDetectionPolicy(d)
	policyFuture, err := secPolicyClient.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, policyInput)
	if err != nil {
		return fmt.Errorf("updating database threat detection policy for %s: %+v", id, err)
	}
	if err := policyFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of database threat detection policy for %s: %+v", id, err)
	}

	return resourceSqlServerRead(d, meta)
}

func resourceSqlServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ServersClient
	connectionClient := meta.(*clients.Client).Sql.ServerConnectionPoliciesClient
	secPolicyClient := meta.(*clients.Client).Sql.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	connection, err := connectionClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Blob Connection Policy for %s: %+v", *id, err)
	}

	secPolicy, err := secPolicyClient.Get(ctx, id.ResourceGroup, id.Name)
	if err == nil {
		if err := d.Set("threat_detection_policy", flattenSqlServerThreatDetectionPolicy(d, secPolicy)); err != nil {
			return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
		}
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("identity", flattenAzureRmSqlServerIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.ServerProperties; props != nil {
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("fully_qualified_domain_name", props.FullyQualifiedDomainName)
		d.Set("version", props.Version)
	}

	if props := connection.ServerConnectionPolicyProperties; props != nil {
		d.Set("connection_policy", string(props.ConnectionType))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSqlServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandAzureRmSqlServerIdentity(input []interface{}) (*sql.ResourceIdentity, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	if expanded.Type == identity.TypeNone {
		return &sql.ResourceIdentity{}, nil
	}

	return &sql.ResourceIdentity{
		Type: sql.IdentityType(string(expanded.Type)),
	}, nil
}

func flattenAzureRmSqlServerIdentity(input *sql.ResourceIdentity) []interface{} {
	var transform *identity.SystemAssigned

	if input != nil {
		transform = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}

		if input.PrincipalID != nil {
			transform.PrincipalId = input.PrincipalID.String()
		}
		if input.TenantID != nil {
			transform.TenantId = input.TenantID.String()
		}
	}

	return identity.FlattenSystemAssigned(transform)
}

func flattenSqlServerThreatDetectionPolicy(d *pluginsdk.ResourceData, policy sql.ServerSecurityAlertPolicy) []interface{} {
	properties := policy.SecurityAlertPolicyProperties

	securityAlertPolicy := make(map[string]interface{})

	securityAlertPolicy["state"] = string(properties.State)
	emailAccountAdmins := false
	if properties.EmailAccountAdmins != nil {
		emailAccountAdmins = *properties.EmailAccountAdmins
	}
	securityAlertPolicy["email_account_admins"] = emailAccountAdmins

	if disabledAlerts := properties.DisabledAlerts; disabledAlerts != nil {
		flattenedAlerts := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
		if v := *disabledAlerts; v != nil {
			for _, a := range v {
				flattenedAlerts.Add(a)
			}
		}
		securityAlertPolicy["disabled_alerts"] = flattenedAlerts
	}
	if emailAddresses := properties.EmailAddresses; emailAddresses != nil {
		flattenedEmails := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
		if v := *emailAddresses; v != nil {
			for _, a := range v {
				flattenedEmails.Add(a)
			}
		}
		securityAlertPolicy["email_addresses"] = flattenedEmails
	}
	if properties.StorageEndpoint != nil {
		securityAlertPolicy["storage_endpoint"] = *properties.StorageEndpoint
	}
	if properties.RetentionDays != nil {
		securityAlertPolicy["retention_days"] = int(*properties.RetentionDays)
	}

	// If storage account access key is in state read it to the new state, as the API does not return it for security reasons
	if v, ok := d.GetOk("threat_detection_policy.0.storage_account_access_key"); ok {
		securityAlertPolicy["storage_account_access_key"] = v.(string)
	}

	return []interface{}{securityAlertPolicy}
}

func expandSqlServerThreatDetectionPolicy(d *pluginsdk.ResourceData) sql.ServerSecurityAlertPolicy {
	policy := sql.ServerSecurityAlertPolicy{
		SecurityAlertPolicyProperties: &sql.SecurityAlertPolicyProperties{
			State: sql.SecurityAlertPolicyStateDisabled,
		},
	}
	properties := policy.SecurityAlertPolicyProperties

	td, ok := d.GetOk("threat_detection_policy")
	if !ok {
		return policy
	}

	if tdl := td.([]interface{}); len(tdl) > 0 {
		securityAlert := tdl[0].(map[string]interface{})

		properties.State = sql.SecurityAlertPolicyState(securityAlert["state"].(string))
		properties.EmailAccountAdmins = utils.Bool(securityAlert["email_account_admins"].(bool))

		if v, ok := securityAlert["disabled_alerts"]; ok {
			alerts := v.(*pluginsdk.Set).List()
			expandedAlerts := make([]string, len(alerts))
			for i, a := range alerts {
				expandedAlerts[i] = a.(string)
			}
			properties.DisabledAlerts = &expandedAlerts
		}
		if v, ok := securityAlert["email_addresses"]; ok {
			emails := v.(*pluginsdk.Set).List()
			expandedEmails := make([]string, len(emails))
			for i, e := range emails {
				expandedEmails[i] = e.(string)
			}
			properties.EmailAddresses = &expandedEmails
		}
		if v, ok := securityAlert["retention_days"]; ok {
			properties.RetentionDays = utils.Int32(int32(v.(int)))
		}
		if v, ok := securityAlert["storage_account_access_key"]; ok {
			properties.StorageAccountAccessKey = utils.String(v.(string))
		}
		if v, ok := securityAlert["storage_endpoint"]; ok {
			properties.StorageEndpoint = utils.String(v.(string))
		}
	}

	return policy
}
