package netapp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/volumes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/volumesreplication"
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
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.VolumeName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

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
				ValidateFunc: azure.ValidateResourceID,
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
				ForceNew: true,
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
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Unix", // Using hardcoded values instead of SDK enum since no matter what case is passed,
					"Ntfs", // ANF changes casing to Pascal case in the backend. Please refer to https://github.com/Azure/azure-sdk-for-go/issues/14684
				}, false),
			},

			"storage_quota_in_gb": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 102400),
			},

			"throughput_in_mibps": {
				Type:     pluginsdk.TypeFloat,
				Optional: true,
				Computed: true,
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

						"remote_volume_location": azure.SchemaLocation(),

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
	volumePath := d.Get("volume_path").(string)
	serviceLevel := volumes.ServiceLevel(d.Get("service_level").(string))
	subnetID := d.Get("subnet_id").(string)

	var networkFeatures volumes.NetworkFeatures
	networkFeaturesString := d.Get("network_features").(string)
	if networkFeaturesString == "" {
		networkFeatures = volumes.NetworkFeaturesBasic
	}
	networkFeatures = volumes.NetworkFeatures(networkFeaturesString)

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
		_, err = snapshotClient.Get(ctx, *parsedSnapshotResourceID)
		if err != nil {
			return fmt.Errorf("getting snapshot from %s: %+v", id, err)
		}

		sourceVolumeId := volumes.NewVolumeID(parsedSnapshotResourceID.SubscriptionId, parsedSnapshotResourceID.ResourceGroupName, parsedSnapshotResourceID.AccountName, parsedSnapshotResourceID.PoolName, parsedSnapshotResourceID.VolumeName)
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
			if !strings.EqualFold(sourceVolumeId.AccountName, id.AccountName) {
				propertyMismatch = append(propertyMismatch, "account_name")
			}
			if !strings.EqualFold(sourceVolumeId.PoolName, id.PoolName) {
				propertyMismatch = append(propertyMismatch, "pool_name")
			}
			if len(propertyMismatch) > 0 {
				return fmt.Errorf("following NetApp Volume properties on new Volume from Snapshot does not match Snapshot's source %s: %s", id, strings.Join(propertyMismatch, ", "))
			}
		}
	}

	parameters := volumes.Volume{
		Location: location,
		Properties: volumes.VolumeProperties{
			CreationToken:   volumePath,
			ServiceLevel:    &serviceLevel,
			SubnetId:        subnetID,
			NetworkFeatures: &networkFeatures,
			ProtocolTypes:   utils.ExpandStringSlice(protocols),
			SecurityStyle:   &securityStyle,
			UsageThreshold:  storageQuotaInGB,
			ExportPolicy:    exportPolicyRule,
			VolumeType:      utils.String(volumeType),
			SnapshotId:      utils.String(snapshotID),
			DataProtection: &volumes.VolumePropertiesDataProtection{
				Replication: dataProtectionReplication.Replication,
				Snapshot:    dataProtectionSnapshotPolicy.Snapshot,
			},
			SnapshotDirectoryVisible: utils.Bool(snapshotDirectoryVisible),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if throughputMibps, ok := d.GetOk("throughput_in_mibps"); ok {
		parameters.Properties.ThroughputMibps = utils.Float(throughputMibps.(float64))
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
	d.Set("account_name", id.AccountName)
	d.Set("pool_name", id.PoolName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		props := model.Properties
		d.Set("volume_path", props.CreationToken)
		d.Set("service_level", props.ServiceLevel)
		d.Set("subnet_id", props.SubnetId)
		d.Set("network_features", props.NetworkFeatures)
		d.Set("protocols", props.ProtocolTypes)
		d.Set("security_style", props.SecurityStyle)
		d.Set("snapshot_directory_visible", props.SnapshotDirectoryVisible)
		d.Set("throughput_in_mibps", props.ThroughputMibps)
		d.Set("storage_quota_in_gb", props.UsageThreshold/1073741824)
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

	// Removing replication if present
	dataProtectionReplicationRaw := d.Get("data_protection_replication").([]interface{})
	dataProtectionReplication := expandNetAppVolumeDataProtectionReplication(dataProtectionReplicationRaw)

	if dataProtectionReplication != nil && dataProtectionReplication.Replication != nil {
		replicaVolumeId, err := volumesreplication.ParseVolumeID(id.ID())
		if err != nil {
			return err
		}
		if dataProtectionReplication.Replication.EndpointType != nil && strings.ToLower(string(*dataProtectionReplication.Replication.EndpointType)) != "dst" {
			// This is the case where primary volume started the deletion, in this case, to be consistent we will remove replication from secondary
			replicaVolumeId, err = volumesreplication.ParseVolumeID(dataProtectionReplication.Replication.RemoteVolumeResourceId)
			if err != nil {
				return err
			}
		}

		replicationClient := meta.(*clients.Client).NetApp.VolumeReplicationClient
		// Checking replication status before deletion, it need to be broken before proceeding with deletion
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
		if err = replicationClient.VolumesDeleteReplicationThenPoll(ctx, *replicaVolumeId); err != nil {
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

func waitForVolumeCreateOrUpdate(ctx context.Context, client *volumes.VolumesClient, id volumes.VolumeId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404"},
		Target:                    []string{"200", "202"},
		Refresh:                   netappVolumeStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish creating: %+v", id, err)
	}

	return nil
}

func waitForReplAuthorization(ctx context.Context, client *volumesreplication.VolumesReplicationClient, id volumesreplication.VolumeId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404", "400"}, // TODO: Remove 400 when bug is fixed on RP side, where replicationStatus returns 400 at some point during authorization process
		Target:                    []string{"200", "202"},
		Refresh:                   netappVolumeReplicationStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for replication authorization %s to complete: %+v", id, err)
	}

	return nil
}

func waitForReplMirrorState(ctx context.Context, client *volumesreplication.VolumesReplicationClient, id volumesreplication.VolumeId, desiredState string) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200"}, // 200 means mirror state is still Mirrored
		Target:                    []string{"204"}, // 204 means mirror state is <> than Mirrored
		Refresh:                   netappVolumeReplicationMirrorStateRefreshFunc(ctx, client, id, desiredState),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be in the state %q: %+v", id, desiredState, err)
	}

	return nil
}

func waitForReplicationDeletion(ctx context.Context, client *volumesreplication.VolumesReplicationClient, id volumesreplication.VolumeId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}

	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202", "400"}, // TODO: Remove 400 when bug is fixed on RP side, where replicationStatus returns 400 while it is in "Deleting" state
		Target:                    []string{"404"},
		Refresh:                   netappVolumeReplicationStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Replication of %s to be deleted: %+v", id, err)
	}

	return nil
}

func waitForVolumeDeletion(ctx context.Context, client *volumes.VolumesClient, id volumes.VolumeId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"204", "404"},
		Refresh:                   netappVolumeStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func netappVolumeStateRefreshFunc(ctx context.Context, client *volumes.VolumesClient, id volumes.VolumeId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving %s: %s", id, err)
			}
		}

		return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
	}
}

func netappVolumeReplicationMirrorStateRefreshFunc(ctx context.Context, client *volumesreplication.VolumesReplicationClient, id volumesreplication.VolumeId, desiredState string) pluginsdk.StateRefreshFunc {
	validStates := []string{"mirrored", "broken", "uninitialized"}

	return func() (interface{}, string, error) {
		// Possible Mirror States to be used as desiredStates:
		// mirrored, broken or uninitialized
		if !utils.SliceContainsValue(validStates, strings.ToLower(desiredState)) {
			return nil, "", fmt.Errorf("Invalid desired mirror state was passed to check mirror replication state (%s), possible values: (%+v)", desiredState, volumesreplication.PossibleValuesForMirrorState())
		}

		res, err := client.VolumesReplicationStatus(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving replication status information from %s: %s", id, err)
			}
		}

		// TODO: fix this refresh function to use strings instead of fake status codes
		// Setting 200 as default response
		response := 200
		if res.Model != nil && res.Model.MirrorState != nil && strings.EqualFold(string(*res.Model.MirrorState), desiredState) {
			// return 204 if state matches desired state
			response = 204
		}

		return res, strconv.Itoa(response), nil
	}
}

func netappVolumeReplicationStateRefreshFunc(ctx context.Context, client *volumesreplication.VolumesReplicationClient, id volumesreplication.VolumeId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.VolumesReplicationStatus(ctx, id)
		if err != nil {
			if httpResponse := res.HttpResponse; httpResponse != nil {
				if httpResponse.StatusCode == 400 && (strings.Contains(strings.ToLower(err.Error()), "deleting") || strings.Contains(strings.ToLower(err.Error()), "volume replication missing or deleted")) {
					// This error can be ignored until a bug is fixed on RP side that it is returning 400 while the replication is in "Deleting" process
					// TODO: remove this workaround when above bug is fixed
				} else if !response.WasNotFound(httpResponse) {
					return nil, "", fmt.Errorf("retrieving replication status from %s: %s", id, err)
				}
			}
		}

		return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
	}
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

	return &volumes.VolumePatchPropertiesExportPolicy{
		Rules: &results,
	}
}

func expandNetAppVolumeDataProtectionReplication(input []interface{}) *volumes.VolumePropertiesDataProtection {
	if len(input) == 0 || input[0] == nil {
		return &volumes.VolumePropertiesDataProtection{}
	}

	replicationObject := volumes.ReplicationObject{}

	replicationRaw := input[0].(map[string]interface{})

	if v, ok := replicationRaw["endpoint_type"]; ok {
		endpointType := volumes.EndpointType(v.(string))
		replicationObject.EndpointType = &endpointType
	}
	if v, ok := replicationRaw["remote_volume_location"]; ok {
		replicationObject.RemoteVolumeRegion = utils.String(v.(string))
	}
	if v, ok := replicationRaw["remote_volume_resource_id"]; ok {
		replicationObject.RemoteVolumeResourceId = v.(string)
	}
	if v, ok := replicationRaw["replication_frequency"]; ok {
		replicationSchedule := volumes.ReplicationSchedule(translateTFSchedule(v.(string)))
		replicationObject.ReplicationSchedule = &replicationSchedule
	}

	return &volumes.VolumePropertiesDataProtection{
		Replication: &replicationObject,
	}
}

func expandNetAppVolumeDataProtectionSnapshotPolicy(input []interface{}) *volumes.VolumePropertiesDataProtection {
	if len(input) == 0 || input[0] == nil {
		return &volumes.VolumePropertiesDataProtection{}
	}

	snapshotObject := volumes.VolumeSnapshotProperties{}

	snapshotRaw := input[0].(map[string]interface{})

	if v, ok := snapshotRaw["snapshot_policy_id"]; ok {
		snapshotObject.SnapshotPolicyId = utils.String(v.(string))
	}

	return &volumes.VolumePropertiesDataProtection{
		Snapshot: &snapshotObject,
	}
}

func expandNetAppVolumeDataProtectionSnapshotPolicyPatch(input []interface{}) *volumes.VolumePatchPropertiesDataProtection {
	if len(input) == 0 || input[0] == nil {
		return &volumes.VolumePatchPropertiesDataProtection{}
	}

	snapshotObject := volumes.VolumeSnapshotProperties{}

	snapshotRaw := input[0].(map[string]interface{})

	if v, ok := snapshotRaw["snapshot_policy_id"]; ok {
		snapshotObject.SnapshotPolicyId = utils.String(v.(string))
	}

	return &volumes.VolumePatchPropertiesDataProtection{
		Snapshot: &snapshotObject,
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
		if item.IpAddress != nil {
			results = append(results, item.IpAddress)
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
	if input == nil || input.Snapshot == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"snapshot_policy_id": input.Snapshot.SnapshotPolicyId,
		},
	}
}

func translateTFSchedule(scheduleName string) string {
	if strings.EqualFold(scheduleName, "10minutes") {
		return "_10minutely"
	}

	return scheduleName
}

func translateSDKSchedule(scheduleName string) string {
	if strings.EqualFold(scheduleName, "_10minutely") {
		return "10minutes"
	}

	return scheduleName
}
