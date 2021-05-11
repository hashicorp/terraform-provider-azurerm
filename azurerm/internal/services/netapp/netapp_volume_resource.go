package netapp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	netAppValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2020-09-01/netapp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceNetAppVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetAppVolumeCreateUpdate,
		Read:   resourceNetAppVolumeRead,
		Update: resourceNetAppVolumeCreateUpdate,
		Delete: resourceNetAppVolumeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VolumeID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.VolumeName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.AccountName,
			},

			"pool_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.PoolName,
			},

			"volume_path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.VolumePath,
			},

			"service_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(netapp.Premium),
					string(netapp.Standard),
					string(netapp.Ultra),
				}, false),
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"create_from_snapshot_resource_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.SnapshotID,
			},

			"protocols": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Computed: true,
				MaxItems: 2,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"NFSv3",
						"NFSv4.1",
						"CIFS",
					}, false),
				},
			},

			"storage_quota_in_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(100, 102400),
			},

			"export_policy_rule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_index": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 5),
						},

						"allowed_clients": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.CIDR,
							},
						},

						"protocols_enabled": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"NFSv3",
									"NFSv4.1",
									"CIFS",
								}, false),
							},
						},

						"cifs_enabled": {
							Type:       schema.TypeBool,
							Optional:   true,
							Computed:   true,
							Deprecated: "Deprecated in favour of `protocols_enabled`",
						},

						"nfsv3_enabled": {
							Type:       schema.TypeBool,
							Optional:   true,
							Computed:   true,
							Deprecated: "Deprecated in favour of `protocols_enabled`",
						},

						"nfsv4_enabled": {
							Type:       schema.TypeBool,
							Optional:   true,
							Computed:   true,
							Deprecated: "Deprecated in favour of `protocols_enabled`",
						},

						"unix_read_only": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"unix_read_write": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"root_access_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"tags": tags.Schema(),

			"mount_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"data_protection_replication": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "dst",
							ValidateFunc: validation.StringInSlice([]string{
								"dst",
							}, false),
						},

						"remote_volume_location": azure.SchemaLocation(),

						"remote_volume_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"replication_frequency": {
							Type:     schema.TypeString,
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
		},
	}
}

func resourceNetAppVolumeCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	poolName := d.Get("pool_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, accountName, poolName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_netapp_volume", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	volumePath := d.Get("volume_path").(string)
	serviceLevel := d.Get("service_level").(string)
	subnetID := d.Get("subnet_id").(string)
	protocols := d.Get("protocols").(*schema.Set).List()
	if len(protocols) == 0 {
		protocols = append(protocols, "NFSv3")
	}

	storageQuotaInGB := int64(d.Get("storage_quota_in_gb").(int) * 1073741824)

	exportPolicyRuleRaw := d.Get("export_policy_rule").([]interface{})
	exportPolicyRule := expandNetAppVolumeExportPolicyRule(exportPolicyRuleRaw)

	dataProtectionReplicationRaw := d.Get("data_protection_replication").([]interface{})
	dataProtectionReplication := expandNetAppVolumeDataProtectionReplication(dataProtectionReplicationRaw)

	authorizeReplication := false
	volumeType := ""
	if dataProtectionReplication != nil && dataProtectionReplication.Replication != nil && strings.ToLower(string(dataProtectionReplication.Replication.EndpointType)) == "dst" {
		authorizeReplication = true
		volumeType = "DataProtection"
	}

	// Handling volume creation from snapshot case
	snapshotResourceID := d.Get("create_from_snapshot_resource_id").(string)
	snapshotID := ""
	if snapshotResourceID != "" {
		// Get snapshot ID GUID value
		parsedSnapshotResourceID, err := parse.SnapshotID(snapshotResourceID)
		if err != nil {
			return fmt.Errorf("Error parsing snapshotResourceID %q: %+v", snapshotResourceID, err)
		}

		snapshotClient := meta.(*clients.Client).NetApp.SnapshotClient
		snapshotResponse, err := snapshotClient.Get(
			ctx,
			parsedSnapshotResourceID.ResourceGroup,
			parsedSnapshotResourceID.NetAppAccountName,
			parsedSnapshotResourceID.CapacityPoolName,
			parsedSnapshotResourceID.VolumeName,
			parsedSnapshotResourceID.Name,
		)
		if err != nil {
			return fmt.Errorf("Error getting snapshot from NetApp Volume %q (Resource Group %q): %+v", parsedSnapshotResourceID.VolumeName, parsedSnapshotResourceID.ResourceGroup, err)
		}
		snapshotID = *snapshotResponse.SnapshotID

		// Validate if properties that cannot be changed matches (protocols, subnet_id, location, resource group, account_name, pool_name, service_level)
		sourceVolume, err := client.Get(
			ctx,
			parsedSnapshotResourceID.ResourceGroup,
			parsedSnapshotResourceID.NetAppAccountName,
			parsedSnapshotResourceID.CapacityPoolName,
			parsedSnapshotResourceID.VolumeName,
		)
		if err != nil {
			return fmt.Errorf("Error getting source NetApp Volume (snapshot's parent resource) %q (Resource Group %q): %+v", parsedSnapshotResourceID.VolumeName, parsedSnapshotResourceID.ResourceGroup, err)
		}

		parsedVolumeID, err := parse.VolumeID(*sourceVolume.ID)
		if err != nil {
			return fmt.Errorf("parsing Source Volume ID: %s", err)
		}
		propertyMismatch := []string{}
		if !ValidateSlicesEquality(*sourceVolume.ProtocolTypes, *utils.ExpandStringSlice(protocols), false) {
			propertyMismatch = append(propertyMismatch, "protocols")
		}
		if !strings.EqualFold(*sourceVolume.SubnetID, subnetID) {
			propertyMismatch = append(propertyMismatch, "subnet_id")
		}
		if !strings.EqualFold(*sourceVolume.Location, location) {
			propertyMismatch = append(propertyMismatch, "location")
		}
		if !strings.EqualFold(string(sourceVolume.ServiceLevel), serviceLevel) {
			propertyMismatch = append(propertyMismatch, "service_level")
		}
		if !strings.EqualFold(parsedVolumeID.ResourceGroup, resourceGroup) {
			propertyMismatch = append(propertyMismatch, "resource_group_name")
		}
		if !strings.EqualFold(parsedVolumeID.NetAppAccountName, accountName) {
			propertyMismatch = append(propertyMismatch, "account_name")
		}
		if !strings.EqualFold(parsedVolumeID.CapacityPoolName, poolName) {
			propertyMismatch = append(propertyMismatch, "pool_name")
		}
		if len(propertyMismatch) > 0 {
			return fmt.Errorf("Following NetApp Volume properties on new Volume from Snapshot does not match Snapshot's source Volume %q (Resource Group %q): %+v", name, resourceGroup, propertyMismatch)
		}
	}

	parameters := netapp.Volume{
		Location: utils.String(location),
		VolumeProperties: &netapp.VolumeProperties{
			CreationToken:  utils.String(volumePath),
			ServiceLevel:   netapp.ServiceLevel(serviceLevel),
			SubnetID:       utils.String(subnetID),
			ProtocolTypes:  utils.ExpandStringSlice(protocols),
			UsageThreshold: utils.Int64(storageQuotaInGB),
			ExportPolicy:   exportPolicyRule,
			VolumeType:     utils.String(volumeType),
			SnapshotID:     utils.String(snapshotID),
			DataProtection: dataProtectionReplication,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, parameters, resourceGroup, accountName, poolName, name)
	if err != nil {
		return fmt.Errorf("Error creating NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// Waiting for volume be completely provisioned
	id := parse.NewVolumeID(client.SubscriptionID, resourceGroup, accountName, poolName, name)
	log.Printf("[DEBUG] Waiting for NetApp Volume Provisioning Service %q (Resource Group %q) to complete", id.Name, id.ResourceGroup)
	if err := waitForVolumeCreation(ctx, client, id, d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	// If this is a data replication secondary volume, authorize replication on primary volume
	if authorizeReplication {
		replVolID, err := parse.VolumeID(*dataProtectionReplication.Replication.RemoteVolumeResourceID)
		if err != nil {
			return err
		}

		future, err := client.AuthorizeReplication(
			ctx,
			replVolID.ResourceGroup,
			replVolID.NetAppAccountName,
			replVolID.CapacityPoolName,
			replVolID.Name,
			netapp.AuthorizeRequest{
				RemoteVolumeResourceID: utils.String(id.ID()),
			},
		)

		if err != nil {
			return fmt.Errorf("Cannot authorize volume replication: %v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Cannot get authorize volume replication future response: %v", err)
		}

		// Wait for volume replication authorization to complete
		log.Printf("[DEBUG] Waiting for replication authorization on NetApp Volume Provisioning Service %q (Resource Group %q) to complete", id.Name, id.ResourceGroup)
		if err := waitForReplAuthorization(ctx, client, id, d.Timeout(schema.TimeoutDelete)); err != nil {
			return err
		}
	}

	d.SetId(id.ID())

	return resourceNetAppVolumeRead(d, meta)
}

func resourceNetAppVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VolumeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] NetApp Volumes %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading NetApp Volumes %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.NetAppAccountName)
	d.Set("pool_name", id.CapacityPoolName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.VolumeProperties; props != nil {
		d.Set("volume_path", props.CreationToken)
		d.Set("service_level", props.ServiceLevel)
		d.Set("subnet_id", props.SubnetID)
		d.Set("protocols", props.ProtocolTypes)
		if props.UsageThreshold != nil {
			d.Set("storage_quota_in_gb", *props.UsageThreshold/1073741824)
		}
		if err := d.Set("export_policy_rule", flattenNetAppVolumeExportPolicyRule(props.ExportPolicy)); err != nil {
			return fmt.Errorf("Error setting `export_policy_rule`: %+v", err)
		}
		if err := d.Set("mount_ip_addresses", flattenNetAppVolumeMountIPAddresses(props.MountTargets)); err != nil {
			return fmt.Errorf("setting `mount_ip_addresses`: %+v", err)
		}
		if err := d.Set("data_protection_replication", flattenNetAppVolumeDataProtectionReplication(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_replication`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceNetAppVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VolumeID(d.Id())
	if err != nil {
		return err
	}

	// Removing replication if present
	dataProtectionReplicationRaw := d.Get("data_protection_replication").([]interface{})
	dataProtectionReplication := expandNetAppVolumeDataProtectionReplication(dataProtectionReplicationRaw)

	if replVolumeID := id; dataProtectionReplication != nil && dataProtectionReplication.Replication != nil {
		if strings.ToLower(string(dataProtectionReplication.Replication.EndpointType)) != "dst" {
			// This is the case where primary volume started the deletion, in this case, to be consistent we will remove replication from secondary
			replVolumeID, err = parse.VolumeID(*dataProtectionReplication.Replication.RemoteVolumeResourceID)
			if err != nil {
				return err
			}
		}

		// Checking replication status before deletion, it need to be broken before proceeding with deletion
		if res, err := client.ReplicationStatusMethod(ctx, replVolumeID.ResourceGroup, replVolumeID.NetAppAccountName, replVolumeID.CapacityPoolName, replVolumeID.Name); err == nil {
			// Wait for replication state = "mirrored"
			if strings.ToLower(string(res.MirrorState)) == "uninitialized" {
				if err := waitForReplMirrorState(ctx, client, *replVolumeID, d.Timeout(schema.TimeoutDelete), "mirrored"); err != nil {
					return err
				}
			}

			// Breaking replication
			_, err = client.BreakReplication(ctx,
				replVolumeID.ResourceGroup,
				replVolumeID.NetAppAccountName,
				replVolumeID.CapacityPoolName,
				replVolumeID.Name,
				&netapp.BreakReplicationRequest{
					ForceBreakReplication: utils.Bool(true),
				})

			if err != nil {
				return fmt.Errorf("Error deleting replication from NetApp Volume %q (Resource Group %q): %+v", replVolumeID.Name, replVolumeID.ResourceGroup, err)
			}

			// Waiting for replication be in broken state
			log.Printf("[DEBUG] Waiting for replication on NetApp Volume Provisioning Service %q (Resource Group %q) to be in broken state", replVolumeID.Name, replVolumeID.ResourceGroup)
			if err := waitForReplMirrorState(ctx, client, *replVolumeID, d.Timeout(schema.TimeoutDelete), "broken"); err != nil {
				return err
			}
		}

		// Deleting replication and waiting for it to fully complete the operation
		if _, err = client.DeleteReplication(ctx, replVolumeID.ResourceGroup, replVolumeID.NetAppAccountName, replVolumeID.CapacityPoolName, replVolumeID.Name); err != nil {
			return fmt.Errorf("Error deleting replication from NetApp Volume %q (Resource Group %q): %+v", replVolumeID.Name, replVolumeID.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Waiting for replication on NetApp Volume Provisioning Service %q (Resource Group %q) to be deleted", replVolumeID.Name, replVolumeID.ResourceGroup)
		if err := waitForReplicationDeletion(ctx, client, *replVolumeID, d.Timeout(schema.TimeoutDelete)); err != nil {
			return err
		}
	}

	// Deleting volume and waiting for it fo fully complete the operation
	if _, err = client.Delete(ctx, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.Name); err != nil {
		return fmt.Errorf("Error deleting NetApp Volume %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for NetApp Volume Provisioning Service %q (Resource Group %q) to be deleted", id.Name, id.ResourceGroup)
	if err := waitForVolumeDeletion(ctx, client, *id, d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	return nil
}

func waitForVolumeCreation(ctx context.Context, client *netapp.VolumesClient, id parse.VolumeId, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404"},
		Target:                    []string{"200", "202"},
		Refresh:                   netappVolumeStateRefreshFunc(ctx, client, id),
		Timeout:                   timeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting NetApp Volume Provisioning Service %q (Resource Group %q) to complete: %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func waitForReplAuthorization(ctx context.Context, client *netapp.VolumesClient, id parse.VolumeId, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404", "400"}, // TODO: Remove 400 when bug is fixed on RP side, where replicationStatus returns 400 at some point during authorization process
		Target:                    []string{"200", "202"},
		Refresh:                   netappVolumeReplicationStateRefreshFunc(ctx, client, id),
		Timeout:                   timeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for replication authorization NetApp Volume Provisioning Service %q (Resource Group %q) to complete: %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func waitForReplMirrorState(ctx context.Context, client *netapp.VolumesClient, id parse.VolumeId, timeout time.Duration, desiredState string) error {
	stateConf := &resource.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200"}, // 200 means mirror state is still Mirrored
		Target:                    []string{"204"}, // 204 means mirror state is <> than Mirrored
		Refresh:                   netappVolumeReplicationMirrorStateRefreshFunc(ctx, client, id, desiredState),
		Timeout:                   timeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NetApp Volume %q (Resource Group %q) to be in %s mirroring state: %+v", id.Name, id.ResourceGroup, desiredState, err)
	}

	return nil
}

func waitForReplicationDeletion(ctx context.Context, client *netapp.VolumesClient, id parse.VolumeId, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202", "400"}, // TODO: Remove 400 when bug is fixed on RP side, where replicationStatus returns 400 while it is in "Deleting" state
		Target:                    []string{"404"},
		Refresh:                   netappVolumeReplicationStateRefreshFunc(ctx, client, id),
		Timeout:                   timeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NetApp Volume replication %q (Resource Group %q) to be deleted: %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func waitForVolumeDeletion(ctx context.Context, client *netapp.VolumesClient, id parse.VolumeId, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"204", "404"},
		Refresh:                   netappVolumeStateRefreshFunc(ctx, client, id),
		Timeout:                   timeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NetApp Volume Provisioning Service %q (Resource Group %q) to be deleted: %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func netappVolumeStateRefreshFunc(ctx context.Context, client *netapp.VolumesClient, id parse.VolumeId) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("Error retrieving NetApp Volume %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
			}
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func netappVolumeReplicationMirrorStateRefreshFunc(ctx context.Context, client *netapp.VolumesClient, id parse.VolumeId, desiredState string) resource.StateRefreshFunc {
	validStates := []string{"mirrored", "broken", "uninitialized"}

	return func() (interface{}, string, error) {
		// Possible Mirror States to be used as desiredStates:
		// mirrored, broken or uninitialized
		if !utils.SliceContainsValue(validStates, strings.ToLower(desiredState)) {
			return nil, "", fmt.Errorf("Invalid desired mirror state was passed to check mirror replication state (%s), possible values: (%+v)", desiredState, netapp.PossibleMirrorStateValues())
		}

		// Setting 200 as default response
		response := 200

		res, err := client.ReplicationStatusMethod(ctx, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("Error retrieving replication status information from NetApp Volume %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
			}
		}

		if strings.EqualFold(string(res.MirrorState), desiredState) {
			// return 204 if state matches desired state
			response = 204
		}

		return res, strconv.Itoa(response), nil
	}
}

func netappVolumeReplicationStateRefreshFunc(ctx context.Context, client *netapp.VolumesClient, id parse.VolumeId) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.ReplicationStatusMethod(ctx, id.ResourceGroup, id.NetAppAccountName, id.CapacityPoolName, id.Name)
		if err != nil {
			if res.StatusCode == 400 && (strings.Contains(strings.ToLower(err.Error()), "deleting") || strings.Contains(strings.ToLower(err.Error()), "volume replication missing or deleted")) {
				// This error can be ignored until a bug is fixed on RP side that it is returning 400 while the replication is in "Deleting" process
				// TODO: remove this workaround when above bug is fixed
			} else if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("Error retrieving replication status from NetApp Volume %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
			}
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func expandNetAppVolumeExportPolicyRule(input []interface{}) *netapp.VolumePropertiesExportPolicy {
	results := make([]netapp.ExportPolicyRule, 0)
	for _, item := range input {
		if item != nil {
			v := item.(map[string]interface{})
			ruleIndex := int32(v["rule_index"].(int))
			allowedClients := strings.Join(*utils.ExpandStringSlice(v["allowed_clients"].(*schema.Set).List()), ",")

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
				} else {
					// TODO: Remove in next major version
					cifsEnabled = v["cifs_enabled"].(bool)
					nfsv3Enabled = v["nfsv3_enabled"].(bool)
					nfsv41Enabled = v["nfsv4_enabled"].(bool)
				}
			}

			unixReadOnly := v["unix_read_only"].(bool)
			unixReadWrite := v["unix_read_write"].(bool)
			rootAccessEnabled := v["root_access_enabled"].(bool)

			result := netapp.ExportPolicyRule{
				AllowedClients: utils.String(allowedClients),
				Cifs:           utils.Bool(cifsEnabled),
				Nfsv3:          utils.Bool(nfsv3Enabled),
				Nfsv41:         utils.Bool(nfsv41Enabled),
				RuleIndex:      utils.Int32(ruleIndex),
				UnixReadOnly:   utils.Bool(unixReadOnly),
				UnixReadWrite:  utils.Bool(unixReadWrite),
				HasRootAccess:  utils.Bool(rootAccessEnabled),
			}

			results = append(results, result)
		}
	}

	return &netapp.VolumePropertiesExportPolicy{
		Rules: &results,
	}
}

func expandNetAppVolumeDataProtectionReplication(input []interface{}) *netapp.VolumePropertiesDataProtection {
	if len(input) == 0 || input[0] == nil {
		return &netapp.VolumePropertiesDataProtection{}
	}

	replicationObject := netapp.ReplicationObject{}

	replicationRaw := input[0].(map[string]interface{})

	if v, ok := replicationRaw["endpoint_type"]; ok {
		replicationObject.EndpointType = netapp.EndpointType(v.(string))
	}
	if v, ok := replicationRaw["remote_volume_location"]; ok {
		replicationObject.RemoteVolumeRegion = utils.String(v.(string))
	}
	if v, ok := replicationRaw["remote_volume_resource_id"]; ok {
		replicationObject.RemoteVolumeResourceID = utils.String(v.(string))
	}
	if v, ok := replicationRaw["replication_frequency"]; ok {
		replicationObject.ReplicationSchedule = netapp.ReplicationSchedule(translateTFSchedule(v.(string)))
	}

	return &netapp.VolumePropertiesDataProtection{
		Replication: &replicationObject,
	}
}

func flattenNetAppVolumeExportPolicyRule(input *netapp.VolumePropertiesExportPolicy) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.Rules == nil {
		return results
	}

	for _, item := range *input.Rules {
		ruleIndex := int32(0)
		if v := item.RuleIndex; v != nil {
			ruleIndex = *v
		}
		allowedClients := []string{}
		if v := item.AllowedClients; v != nil {
			allowedClients = strings.Split(*v, ",")
		}
		// TODO: Start - Remove in next major version
		cifsEnabled := false
		nfsv3Enabled := false
		nfsv4Enabled := false
		// End - Remove in next major version
		protocolsEnabled := []string{}
		if v := item.Cifs; v != nil {
			if *v {
				protocolsEnabled = append(protocolsEnabled, "CIFS")
			}
			cifsEnabled = *v // TODO: Remove in next major version
		}
		if v := item.Nfsv3; v != nil {
			if *v {
				protocolsEnabled = append(protocolsEnabled, "NFSv3")
			}
			nfsv3Enabled = *v // TODO: Remove in next major version
		}
		if v := item.Nfsv41; v != nil {
			if *v {
				protocolsEnabled = append(protocolsEnabled, "NFSv4.1")
			}
			nfsv4Enabled = *v // TODO: Remove in next major version
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

		results = append(results, map[string]interface{}{
			"rule_index":          ruleIndex,
			"allowed_clients":     utils.FlattenStringSlice(&allowedClients),
			"unix_read_only":      unixReadOnly,
			"unix_read_write":     unixReadWrite,
			"root_access_enabled": rootAccessEnabled,
			"protocols_enabled":   utils.FlattenStringSlice(&protocolsEnabled),
			// TODO: Remove in next major version
			"cifs_enabled":  cifsEnabled,
			"nfsv3_enabled": nfsv3Enabled,
			"nfsv4_enabled": nfsv4Enabled,
		})
	}

	return results
}

func flattenNetAppVolumeMountIPAddresses(input *[]netapp.MountTargetProperties) []interface{} {
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

func flattenNetAppVolumeDataProtectionReplication(input *netapp.VolumePropertiesDataProtection) []interface{} {
	if input == nil || input.Replication == nil {
		return []interface{}{}
	}

	if strings.ToLower(string(input.Replication.EndpointType)) == "" || strings.ToLower(string(input.Replication.EndpointType)) != "dst" {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"endpoint_type":             strings.ToLower(string(input.Replication.EndpointType)),
			"remote_volume_location":    input.Replication.RemoteVolumeRegion,
			"remote_volume_resource_id": input.Replication.RemoteVolumeResourceID,
			"replication_frequency":     translateSDKSchedule(strings.ToLower(string(input.Replication.ReplicationSchedule))),
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
