package domainservices

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/domainservices/mgmt/2020-01-01/aad"
	"github.com/hashicorp/go-azure-helpers/lang/response"

	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const DomainServiceResourceName = "azurerm_active_directory_domain_service"

func resourceActiveDirectoryDomainService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceActiveDirectoryDomainServiceCreateUpdate,
		Read:   resourceActiveDirectoryDomainServiceRead,
		Update: resourceActiveDirectoryDomainServiceCreateUpdate,
		Delete: resourceActiveDirectoryDomainServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(3 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(2 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(1 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DomainServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty, // TODO: proper validation
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"domain_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty, // TODO: proper validation, first prefix must be 15 chars or less
			},

			"initial_replica_set": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"domain_controller_ip_addresses": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"external_access_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						// location is computed here
						"location": azure.SchemaLocationForDataSource(),

						"service_status": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: networkValidate.SubnetID,
						},
					},
				},
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"Enterprise",
					"Premium",
				}, false),
			},

			"filtered_sync_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"notifications": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"additional_recipients": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotWhiteSpace,
							},
						},

						"notify_dc_admins": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"notify_global_admins": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"secure_ldap": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},

						"external_access_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"pfx_certificate": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: azValidate.Base64EncodedString,
						},

						"pfx_certificate_password": {
							Type:      pluginsdk.TypeString,
							Required:  true,
							Sensitive: true,
						},

						"certificate_expiry": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"certificate_thumbprint": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"public_certificate": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"security": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ntlm_v1_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"sync_kerberos_passwords": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"sync_ntlm_passwords": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"sync_on_prem_passwords": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"tls_v1_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"tags": tags.Schema(),

			"deployment_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sync_owner": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceActiveDirectoryDomainServiceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	resourceErrorName := fmt.Sprintf("Domain Service (Name: %q, Resource Group: %q)", name, resourceGroup)

	locks.ByName(name, DomainServiceResourceName)
	defer locks.UnlockByName(name, DomainServiceResourceName)

	// If this is a new resource, we cannot determine the resource ID until after it has been created since we need to
	// know the ID of the first replica set.
	var id *parse.DomainServiceId

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", resourceErrorName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			// Parse the replica sets and assume the first one returned to be the initial replica set
			// This is a best effort and the user can choose any replica set if they structure their config accordingly
			props := existing.DomainServiceProperties
			if props == nil {
				return fmt.Errorf("checking for presence of existing %s: API response contained nil or missing properties", resourceErrorName)
			}
			replicaSets := flattenDomainServiceReplicaSets(props.ReplicaSets)
			if len(replicaSets) == 0 {
				return fmt.Errorf("checking for presence of existing %s: API response contained nil or missing replica set details", resourceErrorName)
			}
			initialReplicaSetId := replicaSets[0].(map[string]interface{})["id"].(string)
			id := parse.NewDomainServiceID(client.SubscriptionID, resourceGroup, name, initialReplicaSetId)

			return tf.ImportAsExistsError(DomainServiceResourceName, id.ID())
		}
	} else {
		var err error
		id, err = parse.DomainServiceID(d.Id())
		if err != nil {
			return fmt.Errorf("preparing update for %s: %+v", resourceErrorName, err)
		}
		if id == nil {
			return fmt.Errorf("preparing update for %s: resource ID could not be parsed", resourceErrorName)
		}
	}

	loc := location.Normalize(d.Get("location").(string))
	filteredSync := aad.FilteredSyncDisabled
	if d.Get("filtered_sync_enabled").(bool) {
		filteredSync = aad.FilteredSyncDisabled
	}

	domainService := aad.DomainService{
		DomainServiceProperties: &aad.DomainServiceProperties{
			DomainName:             utils.String(d.Get("domain_name").(string)),
			DomainSecuritySettings: expandDomainServiceSecurity(d.Get("security").([]interface{})),
			FilteredSync:           filteredSync,
			LdapsSettings:          expandDomainServiceLdaps(d.Get("secure_ldap").([]interface{})),
			NotificationSettings:   expandDomainServiceNotifications(d.Get("notifications").([]interface{})),
			Sku:                    utils.String(d.Get("sku").(string)),
		},
		Location: utils.String(loc),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if d.IsNewResource() {
		// On resource creation, specify the initial replica set.
		// No provision is made for changing the initial replica set, it should remain intact for the resource to function properly
		replicaSets := []aad.ReplicaSet{
			{
				Location: utils.String(loc),
				SubnetID: utils.String(d.Get("initial_replica_set.0.subnet_id").(string)),
			},
		}
		domainService.DomainServiceProperties.ReplicaSets = &replicaSets
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, domainService)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", resourceErrorName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s: %+v", resourceErrorName, err)
	}

	// Retrieve the domain service to discover the unique ID for the initial replica set, which should not subsequently change
	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("retrieving %s after creating: %+v", resourceErrorName, err)
		}
		props := resp.DomainServiceProperties
		if props == nil {
			return fmt.Errorf("%s returned with no properties", resourceErrorName)
		}
		if props.ReplicaSets == nil {
			return fmt.Errorf("%s returned with no replica set details", resourceErrorName)
		}

		replicaSets := flattenDomainServiceReplicaSets(props.ReplicaSets)
		if replicaSetCount := len(replicaSets); replicaSetCount != 1 {
			return fmt.Errorf("unexpected number of replica sets for %s: expected 1, saw %d", resourceErrorName, replicaSetCount)
		}

		// Once we know the initial replica set ID, we can build a resource ID
		initialReplicaSetId := replicaSets[0].(map[string]interface{})["id"].(string)
		newId := parse.NewDomainServiceID(client.SubscriptionID, resourceGroup, name, initialReplicaSetId)
		id = &newId
		d.SetId(id.ID())

		if err := d.Set("initial_replica_set", []interface{}{replicaSets[0]}); err != nil {
			return fmt.Errorf("setting `initial_replica_set` after creating resource: %+v", err)
		}
	}

	if id == nil {
		return fmt.Errorf("after creating/updating %s: id was unexpectedly nil", resourceErrorName)
	}

	// A fully deployed domain service has 2 domain controllers per replica set, but the create operation completes early before the DCs are online.
	// The domain service is still provisioning and further operations are blocked until both DCs are up and ready.
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      domainServiceControllerRefreshFunc(ctx, client, *id, false),
		Delay:        1 * time.Minute,
		PollInterval: 1 * time.Minute,
		Timeout:      time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for both domain controllers to become available in initial replica set for %s: %+v", id, err)
	}

	return resourceActiveDirectoryDomainServiceRead(d, meta)
}

func resourceActiveDirectoryDomainServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("resource_id", resp.ID)

	loc := location.NormalizeNilable(resp.Location)
	d.Set("location", loc)

	if props := resp.DomainServiceProperties; props != nil {
		d.Set("deployment_id", props.DeploymentID)
		d.Set("domain_name", props.DomainName)
		d.Set("sync_owner", props.SyncOwner)
		d.Set("tenant_id", props.TenantID)
		d.Set("version", props.Version)

		d.Set("filtered_sync_enabled", false)
		if props.FilteredSync == aad.FilteredSyncEnabled {
			d.Set("filtered_sync_enabled", true)
		}

		d.Set("sku", props.Sku)

		if err := d.Set("notifications", flattenDomainServiceNotifications(props.NotificationSettings)); err != nil {
			return fmt.Errorf("setting `notifications`: %+v", err)
		}

		var initialReplicaSet interface{}
		replicaSets := flattenDomainServiceReplicaSets(props.ReplicaSets)

		// Determine the initial replica set. This is why we need to include InitialReplicaSetId in the resource ID,
		// without it we would not be able to reliably support importing.
		for _, replicaSetRaw := range replicaSets {
			replicaSet := replicaSetRaw.(map[string]interface{})
			if replicaSet["id"].(string) == id.InitialReplicaSetIdName {
				initialReplicaSet = replicaSetRaw
				break
			}
		}
		if initialReplicaSet == nil {
			// It's safest to error out here, since we don't want to wipe the initial replica set from state if it was deleted manually
			return fmt.Errorf("reading %s: could not determine initial replica set from API response", id)
		}
		if err := d.Set("initial_replica_set", []interface{}{initialReplicaSet}); err != nil {
			return fmt.Errorf("setting `initial_replica_set`: %+v", err)
		}

		if err := d.Set("secure_ldap", flattenDomainServiceLdaps(d, props.LdapsSettings, false)); err != nil {
			return fmt.Errorf("setting `secure_ldap`: %+v", err)
		}

		if err := d.Set("security", flattenDomainServiceSecurity(props.DomainSecuritySettings)); err != nil {
			return fmt.Errorf("setting `security`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceActiveDirectoryDomainServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainServiceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
		}
	}

	return nil
}

func domainServiceControllerRefreshFunc(ctx context.Context, client *aad.DomainServicesClient, id parse.DomainServiceId, deleting bool) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Waiting for domain controllers to deploy...")
		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil, "error", err
		}
		if resp.DomainServiceProperties == nil || resp.DomainServiceProperties.ReplicaSets == nil || len(*resp.DomainServiceProperties.ReplicaSets) == 0 {
			return nil, "error", fmt.Errorf("API error: `replicaSets` was not returned")
		}
		// Loop through all replica sets and ensure they are running and each have two available domain controllers
		for _, repl := range *resp.DomainServiceProperties.ReplicaSets {
			if repl.ServiceStatus == nil {
				return resp, "pending", nil
			}
			switch {
			case !deleting && strings.EqualFold(*repl.ServiceStatus, "TearingDown"):
				// Sometimes a service error will cause the replica set, or resource, to self destruct
				return resp, "error", fmt.Errorf("service error: a replica set is unexpectedly tearing down")
			case strings.EqualFold(*repl.ServiceStatus, "Failed"):
				// If a replica set enters a failed state, it needs manual intervention
				return resp, "error", fmt.Errorf("service error: a replica set has entered a Failed state and must be recovered or deleted manually")
			case !strings.EqualFold(*repl.ServiceStatus, "Running"):
				// If it's not yet running, it isn't ready
				return resp, "pending", nil
			case repl.DomainControllerIPAddress == nil || len(*repl.DomainControllerIPAddress) < 2:
				// When a domain controller is online, its IP address will be returned. We're looking for 2 active domain controllers.
				return resp, "pending", nil
			}
		}
		return resp, "available", nil
	}
}

func expandDomainServiceLdaps(input []interface{}) (ldaps *aad.LdapsSettings) {
	ldaps = &aad.LdapsSettings{
		Ldaps: aad.LdapsDisabled,
	}

	if len(input) > 0 {
		v := input[0].(map[string]interface{})
		if v["enabled"].(bool) {
			ldaps.Ldaps = aad.LdapsEnabled
		}
		ldaps.PfxCertificate = utils.String(v["pfx_certificate"].(string))
		ldaps.PfxCertificatePassword = utils.String(v["pfx_certificate_password"].(string))
		if v["external_access_enabled"].(bool) {
			ldaps.ExternalAccess = aad.Enabled
		} else {
			ldaps.ExternalAccess = aad.Disabled
		}
	}

	return
}

func expandDomainServiceNotifications(input []interface{}) *aad.NotificationSettings {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	additionalRecipients := make([]string, 0)
	if ar, ok := v["additional_recipients"]; ok {
		for _, r := range ar.(*pluginsdk.Set).List() {
			additionalRecipients = append(additionalRecipients, r.(string))
		}
	}

	notifyDcAdmins := aad.NotifyDcAdminsDisabled
	if n, ok := v["notify_dc_admins"]; ok && n.(bool) {
		notifyDcAdmins = aad.NotifyDcAdminsEnabled
	}

	notifyGlobalAdmins := aad.NotifyGlobalAdminsDisabled
	if n, ok := v["notify_global_admins"]; ok && n.(bool) {
		notifyGlobalAdmins = aad.NotifyGlobalAdminsEnabled
	}

	return &aad.NotificationSettings{
		AdditionalRecipients: &additionalRecipients,
		NotifyDcAdmins:       notifyDcAdmins,
		NotifyGlobalAdmins:   notifyGlobalAdmins,
	}
}

func expandDomainServiceSecurity(input []interface{}) *aad.DomainSecuritySettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	ntlmV1 := aad.NtlmV1Disabled
	syncKerberosPasswords := aad.SyncKerberosPasswordsDisabled
	syncNtlmPasswords := aad.SyncNtlmPasswordsDisabled
	syncOnPremPasswords := aad.SyncOnPremPasswordsDisabled
	tlsV1 := aad.TLSV1Disabled

	if v["ntlm_v1_enabled"].(bool) {
		ntlmV1 = aad.NtlmV1Enabled
	}
	if v["sync_kerberos_passwords"].(bool) {
		syncKerberosPasswords = aad.SyncKerberosPasswordsEnabled
	}
	if v["sync_ntlm_passwords"].(bool) {
		syncNtlmPasswords = aad.SyncNtlmPasswordsEnabled
	}
	if v["sync_on_prem_passwords"].(bool) {
		syncOnPremPasswords = aad.SyncOnPremPasswordsEnabled
	}
	if v["tls_v1_enabled"].(bool) {
		tlsV1 = aad.TLSV1Enabled
	}

	return &aad.DomainSecuritySettings{
		NtlmV1:                ntlmV1,
		SyncKerberosPasswords: syncKerberosPasswords,
		SyncNtlmPasswords:     syncNtlmPasswords,
		SyncOnPremPasswords:   syncOnPremPasswords,
		TLSV1:                 tlsV1,
	}
}

func flattenDomainServiceLdaps(d *pluginsdk.ResourceData, input *aad.LdapsSettings, dataSource bool) []interface{} {
	result := map[string]interface{}{
		"enabled":                 false,
		"external_access_enabled": false,
		"certificate_expiry":      "",
		"certificate_thumbprint":  "",
		"public_certificate":      "",
	}

	if !dataSource {
		// Read pfx_certificate and pfx_certificate_password from existing state since it's not returned
		result["pfx_certificate"] = ""
		if v, ok := d.GetOk("secure_ldap.0.pfx_certificate"); ok {
			result["pfx_certificate"] = v.(string)
		}
		result["pfx_certificate_password"] = ""
		if v, ok := d.GetOk("secure_ldap.0.pfx_certificate_password"); ok {
			result["pfx_certificate_password"] = v.(string)
		}
	}

	if input != nil {
		if input.ExternalAccess == aad.Enabled {
			result["external_access_enabled"] = true
		}
		if input.Ldaps == aad.LdapsEnabled {
			result["enabled"] = true
		}
		if v := input.CertificateNotAfter; v != nil {
			result["certificate_expiry"] = v.Format(time.RFC3339)
		}
		if v := input.CertificateThumbprint; v != nil {
			result["certificate_thumbprint"] = *v
		}
		if v := input.PublicCertificate; v != nil {
			result["public_certificate"] = *v
		}
	}

	return []interface{}{result}
}

func flattenDomainServiceNotifications(input *aad.NotificationSettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := map[string]interface{}{
		"additional_recipients": make([]string, 0),
		"notify_dc_admins":      false,
		"notify_global_admins":  false,
	}
	if input.AdditionalRecipients != nil {
		result["additional_recipients"] = *input.AdditionalRecipients
	}
	if input.NotifyDcAdmins == aad.NotifyDcAdminsEnabled {
		result["notify_dc_admins"] = true
	}
	if input.NotifyGlobalAdmins == aad.NotifyGlobalAdminsEnabled {
		result["notify_global_admins"] = true
	}

	return []interface{}{result}
}

func flattenDomainServiceReplicaSets(input *[]aad.ReplicaSet) (ret []interface{}) {
	if input == nil {
		return
	}

	for _, in := range *input {
		repl := map[string]interface{}{
			"domain_controller_ip_addresses": make([]string, 0),
			"external_access_ip_address":     "",
			"location":                       location.NormalizeNilable(in.Location),
			"id":                             "",
			"service_status":                 "",
			"subnet_id":                      "",
		}
		if in.DomainControllerIPAddress != nil {
			repl["domain_controller_ip_addresses"] = *in.DomainControllerIPAddress
		}
		if in.ExternalAccessIPAddress != nil {
			repl["external_access_ip_address"] = *in.ExternalAccessIPAddress
		}
		if in.ReplicaSetID != nil {
			repl["id"] = *in.ReplicaSetID
		}
		if in.ServiceStatus != nil {
			repl["service_status"] = *in.ServiceStatus
		}
		if in.SubnetID != nil {
			repl["subnet_id"] = *in.SubnetID
		}
		ret = append(ret, repl)
	}

	return
}

func flattenDomainServiceSecurity(input *aad.DomainSecuritySettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := map[string]bool{
		"ntlm_v1_enabled":         false,
		"sync_kerberos_passwords": false,
		"sync_ntlm_passwords":     false,
		"sync_on_prem_passwords":  false,
		"tls_v1_enabled":          false,
	}
	if input.NtlmV1 == aad.NtlmV1Enabled {
		result["ntlm_v1_enabled"] = true
	}
	if input.SyncKerberosPasswords == aad.SyncKerberosPasswordsEnabled {
		result["sync_kerberos_passwords"] = true
	}
	if input.SyncNtlmPasswords == aad.SyncNtlmPasswordsEnabled {
		result["sync_ntlm_passwords"] = true
	}
	if input.SyncOnPremPasswords == aad.SyncOnPremPasswordsEnabled {
		result["sync_on_prem_passwords"] = true
	}
	if input.TLSV1 == aad.TLSV1Enabled {
		result["tls_v1_enabled"] = true
	}

	return []interface{}{result}
}
