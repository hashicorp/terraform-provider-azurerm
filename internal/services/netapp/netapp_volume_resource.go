// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/volumes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/volumesreplication"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetAppVolume() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetAppVolumeCreate,
		Read:   resourceNetAppVolumeRead,
		Update: resourceNetAppVolumeUpdate,
		Delete: resourceNetAppVolumeDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := volumes.ParseVolumeID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.VolumeName,
			},

			"location": commonschema.Location(),

			"zone": commonschema.ZoneSingleOptionalForceNew(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.AccountName,
			},

			"pool_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.PoolName,
			},

			"volume_path": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.VolumePath,
			},

			"service_level": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(volumes.ServiceLevelPremium),
					string(volumes.ServiceLevelStandard),
					string(volumes.ServiceLevelUltra),
				}, false),
			},

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"create_from_snapshot_resource_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: snapshots.ValidateSnapshotID,
			},

			"network_features": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(volumes.NetworkFeaturesBasic),
					string(volumes.NetworkFeaturesStandard),
				}, false),
			},

			"protocols": {
				Type:     pluginsdk.TypeSet,
				ForceNew: true,
				Optional: true,
				Computed: true,
				MaxItems: 2,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"NFSv3",
						"NFSv4.1",
						"CIFS",
					}, false),
				},
			},

			"security_style": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForSecurityStyle(), false),
			},

			"storage_quota_in_gb": {
				Type:     pluginsdk.TypeInt,
				Required: true,
			},

			"throughput_in_mibps": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.FloatAtLeast(1.0),
			},

			"export_policy_rule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 5,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"rule_index": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 5),
						},

						"allowed_clients": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.CIDR,
							},
						},

						"protocols_enabled": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"NFSv3",
									"NFSv4.1",
									"CIFS",
								}, false),
							},
						},

						"unix_read_only": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"unix_read_write": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"root_access_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"tags": commonschema.Tags(),

			"mount_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"snapshot_directory_visible": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"data_protection_replication": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"endpoint_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "dst",
							ValidateFunc: validation.StringInSlice([]string{
								"dst",
							}, false),
						},

						"remote_volume_location": commonschema.Location(),

						"remote_volume_resource_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"replication_frequency": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"10minutes",
								"daily",
								"hourly",
							}, false),
						},
					},
				},
			},

			"data_protection_snapshot_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"snapshot_policy_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"azure_vmware_data_store_enabled": {
				Type:     pluginsdk.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},

			"encryption_key_source": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(volumes.PossibleValuesForEncryptionKeySource(), false),
			},

			"key_vault_private_endpoint_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: azure.ValidateResourceID,
				RequiredWith: []string{"encryption_key_source"},
			},

			"smb_non_browsable_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"smb_access_based_enumeration_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"is_large_volume": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNetAppVolumeCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := volumes.NewVolumeID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("pool_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_netapp_volume", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	zones := &[]string{}
	if v, ok := d.GetOk("zone"); ok {
		zones = &[]string{
			v.(string),
		}
	}

	volumePath := d.Get("volume_path").(string)
	serviceLevel := volumes.ServiceLevel(d.Get("service_level").(string))
	subnetID := d.Get("subnet_id").(string)

	var networkFeatures volumes.NetworkFeatures
	networkFeaturesString := d.Get("network_features").(string)
	if networkFeaturesString == "" {
		networkFeatures = volumes.NetworkFeaturesBasic
	}
	networkFeatures = volumes.NetworkFeatures(networkFeaturesString)

	smbNonBrowsable := volumes.SmbNonBrowsableDisabled
	if d.Get("smb_non_browsable_enabled").(bool) {
		smbNonBrowsable = volumes.SmbNonBrowsableEnabled
	}

	smbAccessBasedEnumeration := volumes.SmbAccessBasedEnumerationDisabled
	if d.Get("smb_access_based_enumeration_enabled").(bool) {
		smbAccessBasedEnumeration = volumes.SmbAccessBasedEnumerationEnabled
	}

	protocols := d.Get("protocols").(*pluginsdk.Set).List()
	if len(protocols) == 0 {
		protocols = append(protocols, "NFSv3")
	}

	// Handling security style property
	securityStyle := volumes.SecurityStyle(d.Get("security_style").(string))
	if strings.EqualFold(string(securityStyle), "unix") && len(protocols) == 1 && strings.EqualFold(protocols[0].(string), "cifs") {
		return fmt.Errorf("unix security style cannot be used in a CIFS enabled volume for %s", id)
	}
	if strings.EqualFold(string(securityStyle), "ntfs") && len(protocols) == 1 && (strings.EqualFold(protocols[0].(string), "nfsv3") || strings.EqualFold(protocols[0].(string), "nfsv4.1")) {
		return fmt.Errorf("ntfs security style cannot be used in a NFSv3/NFSv4.1 enabled volume for %s", id)
	}

	storageQuotaInGB := int64(d.Get("storage_quota_in_gb").(int) * 1073741824)

	exportPolicyRuleRaw := d.Get("export_policy_rule").([]interface{})
	exportPolicyRule := expandNetAppVolumeExportPolicyRule(exportPolicyRuleRaw)

	dataProtectionReplicationRaw := d.Get("data_protection_replication").([]interface{})
	dataProtectionSnapshotPolicyRaw := d.Get("data_protection_snapshot_policy").([]interface{})

	dataProtectionReplication := expandNetAppVolumeDataProtectionReplication(dataProtectionReplicationRaw)
	dataProtectionSnapshotPolicy := expandNetAppVolumeDataProtectionSnapshotPolicy(dataProtectionSnapshotPolicyRaw)

	authorizeReplication := false
	volumeType := ""
	if dataProtectionReplication != nil && dataProtectionReplication.Replication != nil {
		endpointType := ""
		if dataProtectionReplication.Replication.EndpointType != nil {
			endpointType = string(*dataProtectionReplication.Replication.EndpointType)
		}
		if strings.ToLower(endpointType) == "dst" {
			authorizeReplication = true
			volumeType = "DataProtection"
		}
	}

	// Validating that snapshot policies are not being created in a data protection volume
	if dataProtectionSnapshotPolicy.Snapshot != nil && volumeType != "" {
		return fmt.Errorf("snapshot policy cannot be enabled on a data protection volume, NetApp Volume %q (Resource Group %q)", id.VolumeName, id.ResourceGroupName)
	}

	snapshotDirectoryVisible := d.Get("snapshot_directory_visible").(bool)

	// Handling volume creation from snapshot case
	snapshotResourceID := d.Get("create_from_snapshot_resource_id").(string)
	snapshotID := ""
	if snapshotResourceID != "" {
		// Get snapshot ID GUID value
		parsedSnapshotResourceID, err := snapshots.ParseSnapshotID(snapshotResourceID)
		if err != nil {
			return fmt.Errorf("parsing snapshotResourceID %q: %+v", snapshotResourceID, err)
		}

		snapshotClient := meta.(*clients.Client).NetApp.SnapshotClient
		snapshotResponse, err := snapshotClient.Get(ctx, *parsedSnapshotResourceID)
		if err != nil {
			return fmt.Errorf("getting snapshot from %s: %+v", id, err)
		}
		if model := snapshotResponse.Model; model != nil && model.Id != nil {
			snapshotID = *model.Id
		}

		sourceVolumeId := volumes.NewVolumeID(parsedSnapshotResourceID.SubscriptionId, parsedSnapshotResourceID.ResourceGroupName, parsedSnapshotResourceID.NetAppAccountName, parsedSnapshotResourceID.CapacityPoolName, parsedSnapshotResourceID.VolumeName)
		// Validate if properties that cannot be changed matches (protocols, subnet_id, location, resource group, account_name, pool_name, service_level)
		sourceVolume, err := client.Get(ctx, sourceVolumeId)
		if err != nil {
			return fmt.Errorf("getting source NetApp Volume (snapshot's parent resource) %q (Resource Group %q): %+v", parsedSnapshotResourceID.VolumeName, parsedSnapshotResourceID.ResourceGroupName, err)
		}

		propertyMismatch := []string{}
		if model := sourceVolume.Model; model != nil {
			props := model.Properties
			if !ValidateSlicesEquality(*props.ProtocolTypes, *utils.ExpandStringSlice(protocols), false) {
				propertyMismatch = append(propertyMismatch, "protocols")
			}
			if !strings.EqualFold(props.SubnetId, subnetID) {
				propertyMismatch = append(propertyMismatch, "subnet_id")
			}
			if !strings.EqualFold(model.Location, location) {
				propertyMismatch = append(propertyMismatch, "location")
			}
			if volumeServiceLevel := props.ServiceLevel; volumeServiceLevel != nil {
				if !strings.EqualFold(string(*props.ServiceLevel), string(serviceLevel)) {
					propertyMismatch = append(propertyMismatch, "service_level")
				}
			}
			if !strings.EqualFold(sourceVolumeId.ResourceGroupName, id.ResourceGroupName) {
				propertyMismatch = append(propertyMismatch, "resource_group_name")
			}
			if !strings.EqualFold(sourceVolumeId.NetAppAccountName, id.NetAppAccountName) {
				propertyMismatch = append(propertyMismatch, "account_name")
			}
			if !strings.EqualFold(sourceVolumeId.CapacityPoolName, id.CapacityPoolName) {
				propertyMismatch = append(propertyMismatch, "pool_name")
			}
			if len(propertyMismatch) > 0 {
				return fmt.Errorf("following NetApp Volume properties on new Volume from Snapshot does not match Snapshot's source %s: %s", id, strings.Join(propertyMismatch, ", "))
			}
		}
	}

	avsDataStoreEnabled := volumes.AvsDataStoreDisabled
	if d.Get("azure_vmware_data_store_enabled").(bool) {
		avsDataStoreEnabled = volumes.AvsDataStoreEnabled
	}

	isLargeVolume := d.Get("is_large_volume").(bool)

	parameters := volumes.Volume{
		Location: location,
		Properties: volumes.VolumeProperties{
			CreationToken:             volumePath,
			ServiceLevel:              &serviceLevel,
			SubnetId:                  subnetID,
			NetworkFeatures:           &networkFeatures,
			SmbNonBrowsable:           &smbNonBrowsable,
			SmbAccessBasedEnumeration: &smbAccessBasedEnumeration,
			ProtocolTypes:             utils.ExpandStringSlice(protocols),
			SecurityStyle:             &securityStyle,
			UsageThreshold:            storageQuotaInGB,
			ExportPolicy:              exportPolicyRule,
			VolumeType:                utils.String(volumeType),
			SnapshotId:                utils.String(snapshotID),
			DataProtection: &volumes.VolumePropertiesDataProtection{
				Replication: dataProtectionReplication.Replication,
				Snapshot:    dataProtectionSnapshotPolicy.Snapshot,
			},
			AvsDataStore:             &avsDataStoreEnabled,
			SnapshotDirectoryVisible: utils.Bool(snapshotDirectoryVisible),
			IsLargeVolume:            &isLargeVolume,
		},
		Tags:  tags.Expand(d.Get("tags").(map[string]interface{})),
		Zones: zones,
	}

	if throughputMibps, ok := d.GetOk("throughput_in_mibps"); ok {
		parameters.Properties.ThroughputMibps = utils.Float(throughputMibps.(float64))
	}

	if encryptionKeySource, ok := d.GetOk("encryption_key_source"); ok {
		// Validating Microsoft.KeyVault encryption key provider is enabled only on Standard network features
		if volumes.EncryptionKeySource(encryptionKeySource.(string)) == volumes.EncryptionKeySourceMicrosoftPointKeyVault && networkFeatures == volumes.NetworkFeaturesBasic {
			return fmt.Errorf("volume encryption cannot be enabled when network features is set to basic: %s", id.ID())
		}

		parameters.Properties.EncryptionKeySource = pointer.To(volumes.EncryptionKeySource(encryptionKeySource.(string)))
	}

	if keyVaultPrivateEndpointID, ok := d.GetOk("key_vault_private_endpoint_id"); ok {
		parameters.Properties.KeyVaultPrivateEndpointResourceId = pointer.To(keyVaultPrivateEndpointID.(string))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Waiting for volume be completely provisioned
	if err := waitForVolumeCreateOrUpdate(ctx, client, id); err != nil {
		return err
	}

	// If this is a data replication secondary volume, authorize replication on primary volume
	if authorizeReplication {
		replicationClient := meta.(*clients.Client).NetApp.VolumeReplicationClient
		replVolID, err := volumesreplication.ParseVolumeID(dataProtectionReplication.Replication.RemoteVolumeResourceId)
		if err != nil {
			return err
		}

		if err = replicationClient.VolumesAuthorizeReplicationThenPoll(ctx, *replVolID, volumesreplication.AuthorizeRequest{
			RemoteVolumeResourceId: utils.String(id.ID()),
		},
		); err != nil {
			return fmt.Errorf("cannot authorize volume replication: %v", err)
		}

		// Wait for volume replication authorization to complete
		log.Printf("[DEBUG] Waiting for replication authorization on %s to complete", id)
		if err := waitForReplAuthorization(ctx, replicationClient, *replVolID); err != nil {
			return err
		}
	}

	d.SetId(id.ID())

	return resourceNetAppVolumeRead(d, meta)
}

func resourceNetAppVolumeUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := volumes.ParseVolumeID(d.Id())
	if err != nil {
		return err
	}

	shouldUpdate := false
	update := volumes.VolumePatch{
		Properties: &volumes.VolumePatchProperties{},
	}

	if d.HasChange("zones") {
		return fmt.Errorf("zone changes are not supported after volume is already created, %s", id)
	}

	if d.HasChange("storage_quota_in_gb") {
		shouldUpdate = true
		storageQuotaInBytes := int64(d.Get("storage_quota_in_gb").(int) * 1073741824)
		update.Properties.UsageThreshold = utils.Int64(storageQuotaInBytes)
	}

	if d.HasChange("export_policy_rule") {
		shouldUpdate = true
		exportPolicyRuleRaw := d.Get("export_policy_rule").([]interface{})
		exportPolicyRule := expandNetAppVolumeExportPolicyRulePatch(exportPolicyRuleRaw)
		update.Properties.ExportPolicy = exportPolicyRule
	}

	if d.HasChange("data_protection_snapshot_policy") {
		// Validating that snapshot policies are not being created in a data protection volume
		dataProtectionReplicationRaw := d.Get("data_protection_replication").([]interface{})
		dataProtectionReplication := expandNetAppVolumeDataProtectionReplication(dataProtectionReplicationRaw)

		if dataProtectionReplication != nil && dataProtectionReplication.Replication != nil && dataProtectionReplication.Replication.EndpointType != nil && strings.ToLower(string(*dataProtectionReplication.Replication.EndpointType)) == "dst" {
			return fmt.Errorf("snapshot policy cannot be enabled on a data protection volume, %s", id)
		}

		shouldUpdate = true
		dataProtectionSnapshotPolicyRaw := d.Get("data_protection_snapshot_policy").([]interface{})
		dataProtectionSnapshotPolicy := expandNetAppVolumeDataProtectionSnapshotPolicyPatch(dataProtectionSnapshotPolicyRaw)
		update.Properties.DataProtection = dataProtectionSnapshotPolicy
	}

	if d.HasChange("throughput_in_mibps") {
		shouldUpdate = true
		throughputMibps := d.Get("throughput_in_mibps")
		update.Properties.ThroughputMibps = utils.Float(throughputMibps.(float64))
	}

	if d.HasChange("smb_non_browsable_enabled") {
		shouldUpdate = true
		smbNonBrowsable := volumes.SmbNonBrowsableDisabled
		update.Properties.SmbNonBrowsable = &smbNonBrowsable
		if d.Get("smb_non_browsable_enabled").(bool) {
			smbNonBrowsable := volumes.SmbNonBrowsableEnabled
			update.Properties.SmbNonBrowsable = &smbNonBrowsable
		}
	}

	if d.HasChange("smb_access_based_enumeration_enabled") {
		shouldUpdate = true
		smbAccessBasedEnumeration := volumes.SmbAccessBasedEnumerationDisabled
		update.Properties.SmbAccessBasedEnumeration = &smbAccessBasedEnumeration
		if d.Get("smb_access_based_enumeration_enabled").(bool) {
			smbAccessBasedEnumeration := volumes.SmbAccessBasedEnumerationEnabled
			update.Properties.SmbAccessBasedEnumeration = &smbAccessBasedEnumeration
		}
	}

	if d.HasChange("tags") {
		shouldUpdate = true
		tagsRaw := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(tagsRaw)
	}

	if shouldUpdate {
		if err = client.UpdateThenPoll(ctx, *id, update); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		// Wait for volume to complete update
		if err := waitForVolumeCreateOrUpdate(ctx, client, *id); err != nil {
			return err
		}
	}

	return resourceNetAppVolumeRead(d, meta)
}

func resourceNetAppVolumeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := volumes.ParseVolumeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.VolumeName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.NetAppAccountName)
	d.Set("pool_name", id.CapacityPoolName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		zone := ""
		if model.Zones != nil {
			if zones := *model.Zones; len(zones) > 0 {
				zone = zones[0]
			}
		}
		d.Set("zone", zone)

		props := model.Properties
		d.Set("volume_path", props.CreationToken)
		d.Set("service_level", string(pointer.From(props.ServiceLevel)))
		d.Set("subnet_id", props.SubnetId)
		d.Set("network_features", string(pointer.From(props.NetworkFeatures)))
		d.Set("protocols", props.ProtocolTypes)
		d.Set("security_style", string(pointer.From(props.SecurityStyle)))
		d.Set("snapshot_directory_visible", props.SnapshotDirectoryVisible)
		d.Set("throughput_in_mibps", props.ThroughputMibps)
		d.Set("storage_quota_in_gb", props.UsageThreshold/1073741824)
		d.Set("encryption_key_source", string(pointer.From(props.EncryptionKeySource)))
		d.Set("key_vault_private_endpoint_id", props.KeyVaultPrivateEndpointResourceId)

		smbNonBrowsable := false
		if props.SmbNonBrowsable != nil {
			smbNonBrowsable = strings.EqualFold(string(*props.SmbNonBrowsable), string(volumes.SmbNonBrowsableEnabled))
		}
		d.Set("smb_non_browsable_enabled", smbNonBrowsable)

		smbAccessBasedEnumeration := false
		if props.SmbAccessBasedEnumeration != nil {
			smbAccessBasedEnumeration = strings.EqualFold(string(*props.SmbAccessBasedEnumeration), string(volumes.SmbAccessBasedEnumerationEnabled))
		}
		d.Set("smb_access_based_enumeration_enabled", smbAccessBasedEnumeration)

		d.Set("is_large_volume", props.IsLargeVolume)

		avsDataStore := false
		if props.AvsDataStore != nil {
			avsDataStore = strings.EqualFold(string(*props.AvsDataStore), string(volumes.AvsDataStoreEnabled))
		}
		d.Set("azure_vmware_data_store_enabled", avsDataStore)

		if err := d.Set("export_policy_rule", flattenNetAppVolumeExportPolicyRule(props.ExportPolicy)); err != nil {
			return fmt.Errorf("setting `export_policy_rule`: %+v", err)
		}
		if err := d.Set("mount_ip_addresses", flattenNetAppVolumeMountIPAddresses(props.MountTargets)); err != nil {
			return fmt.Errorf("setting `mount_ip_addresses`: %+v", err)
		}
		if err := d.Set("data_protection_replication", flattenNetAppVolumeDataProtectionReplication(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_replication`: %+v", err)
		}
		if err := d.Set("data_protection_snapshot_policy", flattenNetAppVolumeDataProtectionSnapshotPolicy(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_snapshot_policy`: %+v", err)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceNetAppVolumeDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := volumes.ParseVolumeID(d.Id())
	if err != nil {
		return err
	}

	netApp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("fetching netapp error: %+v", err)
	}

	if netApp.Model != nil && netApp.Model.Properties.DataProtection != nil {
		dataProtectionReplication := netApp.Model.Properties.DataProtection
		replicaVolumeId, err := volumesreplication.ParseVolumeID(id.ID())
		if err != nil {
			return err
		}
		if dataProtectionReplication.Replication != nil && dataProtectionReplication.Replication.EndpointType != nil && strings.ToLower(string(*dataProtectionReplication.Replication.EndpointType)) != "dst" {
			// This is the case where primary volume started the deletion, in this case, to be consistent we will remove replication from secondary
			replicaVolumeId, err = volumesreplication.ParseVolumeID(dataProtectionReplication.Replication.RemoteVolumeResourceId)
			if err != nil {
				return err
			}
		}

		replicationClient := meta.(*clients.Client).NetApp.VolumeReplicationClient
		// Checking replication status before deletion, it needs to be broken before proceeding with deletion
		if res, err := replicationClient.VolumesReplicationStatus(ctx, *replicaVolumeId); err == nil {
			// Wait for replication state = "mirrored"
			if model := res.Model; model != nil {
				if model.MirrorState != nil && strings.ToLower(string(*model.MirrorState)) == "uninitialized" {
					if err := waitForReplMirrorState(ctx, replicationClient, *replicaVolumeId, "mirrored"); err != nil {
						return fmt.Errorf("waiting for replica %s to become 'mirrored': %+v", *replicaVolumeId, err)
					}
				}
			}

			// Breaking replication
			if err = replicationClient.VolumesBreakReplicationThenPoll(ctx, *replicaVolumeId, volumesreplication.BreakReplicationRequest{
				ForceBreakReplication: utils.Bool(true),
			}); err != nil {
				return fmt.Errorf("breaking replication for %s: %+v", *replicaVolumeId, err)
			}

			// Waiting for replication be in broken state
			log.Printf("[DEBUG] Waiting for the replication of %s to be in broken state", *replicaVolumeId)
			if err := waitForReplMirrorState(ctx, replicationClient, *replicaVolumeId, "broken"); err != nil {
				return fmt.Errorf("waiting for the breaking of replication for %s: %+v", *replicaVolumeId, err)
			}
		}

		// Deleting replication and waiting for it to fully complete the operation
		if _, err = replicationClient.VolumesDeleteReplication(ctx, *replicaVolumeId); err != nil {
			return fmt.Errorf("deleting replicate %s: %+v", *replicaVolumeId, err)
		}

		if err := waitForReplicationDeletion(ctx, replicationClient, *replicaVolumeId); err != nil {
			return fmt.Errorf("waiting for the replica %s to be deleted: %+v", *replicaVolumeId, err)
		}
	}

	// Deleting volume and waiting for it fo fully complete the operation
	if err = client.DeleteThenPoll(ctx, *id, volumes.DeleteOperationOptions{
		ForceDelete: utils.Bool(true),
	}); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = waitForVolumeDeletion(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandNetAppVolumeExportPolicyRule(input []interface{}) *volumes.VolumePropertiesExportPolicy {
	results := make([]volumes.ExportPolicyRule, 0)
	for _, item := range input {
		if item != nil {
			v := item.(map[string]interface{})
			ruleIndex := int64(v["rule_index"].(int))
			allowedClients := strings.Join(*utils.ExpandStringSlice(v["allowed_clients"].(*pluginsdk.Set).List()), ",")

			cifsEnabled := false
			nfsv3Enabled := false
			nfsv41Enabled := false

			if vpe := v["protocols_enabled"]; vpe != nil {
				protocolsEnabled := vpe.([]interface{})
				if len(protocolsEnabled) != 0 {
					for _, protocol := range protocolsEnabled {
						if protocol != nil {
							switch strings.ToLower(protocol.(string)) {
							case "cifs":
								cifsEnabled = true
							case "nfsv3":
								nfsv3Enabled = true
							case "nfsv4.1":
								nfsv41Enabled = true
							}
						}
					}
				}
			}

			unixReadOnly := v["unix_read_only"].(bool)
			unixReadWrite := v["unix_read_write"].(bool)
			rootAccessEnabled := v["root_access_enabled"].(bool)

			result := volumes.ExportPolicyRule{
				AllowedClients: utils.String(allowedClients),
				Cifs:           utils.Bool(cifsEnabled),
				Nfsv3:          utils.Bool(nfsv3Enabled),
				Nfsv41:         utils.Bool(nfsv41Enabled),
				RuleIndex:      utils.Int64(ruleIndex),
				UnixReadOnly:   utils.Bool(unixReadOnly),
				UnixReadWrite:  utils.Bool(unixReadWrite),
				HasRootAccess:  utils.Bool(rootAccessEnabled),
			}

			results = append(results, result)
		}
	}

	return &volumes.VolumePropertiesExportPolicy{
		Rules: &results,
	}
}

func expandNetAppVolumeExportPolicyRulePatch(input []interface{}) *volumes.VolumePatchPropertiesExportPolicy {
	results := make([]volumes.ExportPolicyRule, 0)
	for _, item := range input {
		if item != nil {
			v := item.(map[string]interface{})
			ruleIndex := int64(v["rule_index"].(int))
			allowedClients := strings.Join(*utils.ExpandStringSlice(v["allowed_clients"].(*pluginsdk.Set).List()), ",")

			nfsv3Enabled := false
			nfsv41Enabled := false
			cifsEnabled := false

			if vpe := v["protocols_enabled"]; vpe != nil {
				protocolsEnabled := vpe.([]interface{})
				if len(protocolsEnabled) != 0 {
					for _, protocol := range protocolsEnabled {
						if protocol != nil {
							switch strings.ToLower(protocol.(string)) {
							case "cifs":
								cifsEnabled = true
							case "nfsv3":
								nfsv3Enabled = true
							case "nfsv4.1":
								nfsv41Enabled = true
							}
						}
					}
				}
			}

			unixReadOnly := v["unix_read_only"].(bool)
			unixReadWrite := v["unix_read_write"].(bool)
			rootAccessEnabled := v["root_access_enabled"].(bool)

			result := volumes.ExportPolicyRule{
				AllowedClients: utils.String(allowedClients),
				Cifs:           utils.Bool(cifsEnabled),
				Nfsv3:          utils.Bool(nfsv3Enabled),
				Nfsv41:         utils.Bool(nfsv41Enabled),
				RuleIndex:      utils.Int64(ruleIndex),
				UnixReadOnly:   utils.Bool(unixReadOnly),
				UnixReadWrite:  utils.Bool(unixReadWrite),
				HasRootAccess:  utils.Bool(rootAccessEnabled),
			}

			results = append(results, result)
		}
	}

	return &volumes.VolumePatchPropertiesExportPolicy{
		Rules: &results,
	}
}

func flattenNetAppVolumeExportPolicyRule(input *volumes.VolumePropertiesExportPolicy) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.Rules == nil {
		return results
	}

	for _, item := range *input.Rules {
		ruleIndex := int64(0)
		if v := item.RuleIndex; v != nil {
			ruleIndex = *v
		}
		allowedClients := []string{}
		if v := item.AllowedClients; v != nil {
			allowedClients = strings.Split(*v, ",")
		}

		protocolsEnabled := []string{}
		if v := item.Cifs; v != nil {
			if *v {
				protocolsEnabled = append(protocolsEnabled, "CIFS")
			}
		}
		if v := item.Nfsv3; v != nil {
			if *v {
				protocolsEnabled = append(protocolsEnabled, "NFSv3")
			}
		}
		if v := item.Nfsv41; v != nil {
			if *v {
				protocolsEnabled = append(protocolsEnabled, "NFSv4.1")
			}
		}
		unixReadOnly := false
		if v := item.UnixReadOnly; v != nil {
			unixReadOnly = *v
		}
		unixReadWrite := false
		if v := item.UnixReadWrite; v != nil {
			unixReadWrite = *v
		}
		rootAccessEnabled := false
		if v := item.HasRootAccess; v != nil {
			rootAccessEnabled = *v
		}

		result := map[string]interface{}{
			"rule_index":          ruleIndex,
			"allowed_clients":     utils.FlattenStringSlice(&allowedClients),
			"unix_read_only":      unixReadOnly,
			"unix_read_write":     unixReadWrite,
			"root_access_enabled": rootAccessEnabled,
			"protocols_enabled":   utils.FlattenStringSlice(&protocolsEnabled),
		}
		results = append(results, result)
	}

	return results
}

func flattenNetAppVolumeMountIPAddresses(input *[]volumes.MountTargetProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.IPAddress != nil {
			results = append(results, item.IPAddress)
		}
	}

	return results
}

func flattenNetAppVolumeDataProtectionReplication(input *volumes.VolumePropertiesDataProtection) []interface{} {
	if input == nil || input.Replication == nil || input.Replication.EndpointType == nil {
		return []interface{}{}
	}

	if strings.ToLower(string(*input.Replication.EndpointType)) == "" || strings.ToLower(string(*input.Replication.EndpointType)) != "dst" {
		return []interface{}{}
	}

	replicationFrequency := ""
	if input.Replication.ReplicationSchedule != nil {
		replicationFrequency = translateSDKSchedule(strings.ToLower(string(*input.Replication.ReplicationSchedule)))
	}

	return []interface{}{
		map[string]interface{}{
			"endpoint_type":             strings.ToLower(string(*input.Replication.EndpointType)),
			"remote_volume_location":    location.NormalizeNilable(input.Replication.RemoteVolumeRegion),
			"remote_volume_resource_id": input.Replication.RemoteVolumeResourceId,
			"replication_frequency":     replicationFrequency,
		},
	}
}

func flattenNetAppVolumeDataProtectionSnapshotPolicy(input *volumes.VolumePropertiesDataProtection) []interface{} {
	if input == nil || input.Snapshot == nil || input.Snapshot.SnapshotPolicyId == nil || *input.Snapshot.SnapshotPolicyId == "" {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"snapshot_policy_id": input.Snapshot.SnapshotPolicyId,
		},
	}
}
