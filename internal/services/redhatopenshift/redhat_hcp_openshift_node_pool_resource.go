// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "managed_cluster_name:hcp_openshift_cluster_id,resource_group_name:hcp_openshift_cluster_id,subscription_id:hcp_openshift_cluster_id

package redhatopenshift

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/hcpopenshiftclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/nodepools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithUpdate        = RedHatHcpOpenShiftNodePool{}
	_ sdk.ResourceWithIdentity      = RedHatHcpOpenShiftNodePool{}
	_ sdk.ResourceWithCustomizeDiff = RedHatHcpOpenShiftNodePool{}
)

type RedHatHcpOpenShiftNodePool struct{}

type RedHatHcpOpenShiftNodePoolModel struct {
	Name                    string                            `tfschema:"name"`
	ClusterId               string                            `tfschema:"hcp_openshift_cluster_id"`
	Platform                []HcpOpenShiftNodePoolPlatform    `tfschema:"platform"`
	AutoRepair              bool                              `tfschema:"auto_repair"`
	Replicas                int64                             `tfschema:"replicas"`
	NodeDrainTimeoutMinutes int64                             `tfschema:"node_drain_timeout_minutes"`
	AutoScaling             []HcpOpenShiftNodePoolAutoscaling `tfschema:"auto_scaling"`
	Version                 []HcpOpenShiftNodePoolVersion     `tfschema:"version"`
	Labels                  map[string]string                 `tfschema:"labels"`
	Taints                  []HcpOpenShiftNodePoolTaint       `tfschema:"taints"`
	Tags                    map[string]string                 `tfschema:"tags"`
}

type HcpOpenShiftNodePoolPlatform struct {
	VMSize                 string                       `tfschema:"vm_size"`
	EnableEncryptionAtHost bool                         `tfschema:"enable_encryption_at_host"`
	OsDisk                 []HcpOpenShiftNodePoolOsDisk `tfschema:"os_disk"`
	SubnetId               string                       `tfschema:"subnet_id"`
	AvailabilityZone       string                       `tfschema:"availability_zone"`
}

type HcpOpenShiftNodePoolAutoscaling struct {
	Min int64 `tfschema:"min"`
	Max int64 `tfschema:"max"`
}

type HcpOpenShiftNodePoolVersion struct {
	Id           string `tfschema:"id"`
	ChannelGroup string `tfschema:"channel_group"`
}

type HcpOpenShiftNodePoolOsDisk struct {
	SizeGb                 int64  `tfschema:"size_gb"`
	DiskStorageAccountType string `tfschema:"disk_storage_account_type"`
	DiskType               string `tfschema:"disk_type"`
	EncryptionSetId        string `tfschema:"encryption_set_id"`
}

type HcpOpenShiftNodePoolTaint struct {
	Key    string `tfschema:"key"`
	Value  string `tfschema:"value"`
	Effect string `tfschema:"effect"`
}

func (r RedHatHcpOpenShiftNodePool) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{Timeout: 5 * time.Minute, Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		diff := metadata.ResourceDiff
		minPath, maxPath := "auto_scaling.0.min", "auto_scaling.0.max"
		if diff.NewValueKnown(minPath) && diff.NewValueKnown(maxPath) && diff.Get(minPath).(int) > diff.Get(maxPath).(int) {
			return fmt.Errorf("`auto_scaling.min` must not exceed `auto_scaling.max`")
		}
		zonePath := "platform.0.availability_zone"
		if diff.NewValueKnown(zonePath) && diff.Get(zonePath).(string) == "" {
			for _, path := range []string{minPath, maxPath} {
				if diff.NewValueKnown(path) && diff.Get(path).(int) > 200 {
					return fmt.Errorf("`%s` must not exceed 200 when `platform.availability_zone` is not specified", path)
				}
			}
		}
		return nil
	}}
}

func (r RedHatHcpOpenShiftNodePool) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z]([-a-zA-Z0-9]{0,13}[a-zA-Z0-9])?$`),
				"must begin with a letter, end with a letter or number, contain only letters, numbers and hyphens, and be between 1 and 15 characters in length",
			),
		},

		"hcp_openshift_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: hcpopenshiftclusters.ValidateHcpOpenShiftClusterID,
		},

		"platform": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"vm_size": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"availability_zone": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"enable_encryption_at_host": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},
					"os_disk": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"disk_storage_account_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(nodepools.PossibleValuesForDiskStorageAccountType(), false),
									Default:      nodepools.DiskStorageAccountTypePremiumLRS,
								},
								"disk_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(nodepools.PossibleValuesForOsDiskType(), false),
									Default:      nodepools.OsDiskTypeManaged,
								},
								"encryption_set_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: commonids.ValidateDiskEncryptionSetID,
								},
								"size_gb": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(1),
									Default:      64,
								},
							},
						},
					},
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"version": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^\d+\.\d+\.\d+`),
							"must be specified as X.Y.Z with optional pre-release",
						),
					},
					"channel_group": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "stable",
						ValidateFunc: validation.StringInSlice([]string{
							"nightly",
							"candidate",
							"fast",
							"stable",
						}, false),
					},
				},
			},
		},

		"auto_repair": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  true,
		},

		"auto_scaling": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"replicas"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"max": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"min": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"labels": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"node_drain_timeout_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 10080),
			Default:      0,
		},

		"replicas": {
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"auto_scaling"},
			ValidateFunc:  validation.IntAtLeast(0),
		},

		"taints": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"effect": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(nodepools.PossibleValuesForEffect(), false),
					},
					"key": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"value": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r RedHatHcpOpenShiftNodePool) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r RedHatHcpOpenShiftNodePool) ModelObject() any {
	return &RedHatHcpOpenShiftNodePoolModel{}
}

func (r RedHatHcpOpenShiftNodePool) ResourceType() string {
	return "azurerm_redhat_hcp_openshift_node_pool"
}

func (r RedHatHcpOpenShiftNodePool) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nodepools.ValidateNodePoolID
}

func (r RedHatHcpOpenShiftNodePool) Identity() resourceids.ResourceId {
	return &nodepools.NodePoolId{}
}

func (r RedHatHcpOpenShiftNodePool) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{Timeout: 60 * time.Minute, Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		var config RedHatHcpOpenShiftNodePoolModel
		if err := metadata.Decode(&config); err != nil {
			return fmt.Errorf("decoding: %+v", err)
		}
		parent, err := hcpopenshiftclusters.ParseHcpOpenShiftClusterID(config.ClusterId)
		if err != nil {
			return err
		}
		id := nodepools.NewNodePoolID(parent.SubscriptionId, parent.ResourceGroupName, parent.HcpOpenShiftClusterName, config.Name)
		client := metadata.Client.RedHatOpenShift.NodePoolsClient
		if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}
		}
		cluster, err := metadata.Client.RedHatOpenShift.HcpOpenShiftClustersClient.Get(ctx, *parent)
		if err != nil {
			return fmt.Errorf("retrieving parent %s: %+v", parent, err)
		}
		if cluster.Model == nil {
			return fmt.Errorf("retrieving parent %s: `model` was nil", parent)
		}
		parameters, err := expandHcpOpenShiftNodePool(config, cluster.Model.Location)
		if err != nil {
			return err
		}
		if _, ok := metadata.ResourceData.GetOkExists("replicas"); !ok {
			parameters.Properties.Replicas = nil
		}
		if _, ok := metadata.ResourceData.GetOkExists("platform.0.enable_encryption_at_host"); !ok {
			parameters.Properties.Platform.EnableEncryptionAtHost = nil
		}
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *parameters, metadata.SetIDAndIdentityCallback(&id)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
		metadata.SetID(id)
		return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
	}}
}

func (r RedHatHcpOpenShiftNodePool) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{Timeout: 60 * time.Minute, Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := nodepools.ParseNodePoolID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		var config RedHatHcpOpenShiftNodePoolModel
		if err := metadata.Decode(&config); err != nil {
			return fmt.Errorf("decoding: %+v", err)
		}
		client := metadata.Client.RedHatOpenShift.NodePoolsClient
		existing, err := client.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if existing.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil", id)
		}
		parameters, err := expandHcpOpenShiftNodePool(config, existing.Model.Location)
		if err != nil {
			return err
		}
		if err := client.CreateOrUpdateThenPoll(ctx, *id, *parameters); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
		return nil
	}}
}

func (r RedHatHcpOpenShiftNodePool) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{Timeout: 5 * time.Minute, Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := nodepools.ParseNodePoolID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		resp, err := metadata.Client.RedHatOpenShift.NodePoolsClient.Get(ctx, *id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return metadata.MarkAsGone(id)
			}
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: `model` was nil", id)
		}
		state, err := flattenHcpOpenShiftNodePool(*id, *resp.Model)
		if err != nil {
			return err
		}
		if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
			return err
		}
		return metadata.Encode(state)
	}}
}

func (r RedHatHcpOpenShiftNodePool) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{Timeout: 60 * time.Minute, Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := nodepools.ParseNodePoolID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		if err := metadata.Client.RedHatOpenShift.NodePoolsClient.DeleteThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
		return nil
	}}
}

func expandHcpOpenShiftNodePool(input RedHatHcpOpenShiftNodePoolModel, parentLocation string) (*nodepools.NodePool, error) {
	if len(input.Platform) != 1 {
		return nil, fmt.Errorf("expanding `platform`: exactly one block is required")
	}
	platform := input.Platform[0]
	props := &nodepools.NodePoolProperties{
		AutoRepair:              pointer.To(input.AutoRepair),
		NodeDrainTimeoutMinutes: pointer.To(input.NodeDrainTimeoutMinutes),
		Platform:                nodepools.NodePoolPlatformProfile{VMSize: platform.VMSize, EnableEncryptionAtHost: pointer.To(platform.EnableEncryptionAtHost)},
	}
	if platform.SubnetId != "" {
		props.Platform.SubnetId = pointer.To(platform.SubnetId)
	}
	if platform.AvailabilityZone != "" {
		props.Platform.AvailabilityZone = pointer.To(platform.AvailabilityZone)
	}
	if len(input.AutoScaling) > 0 {
		props.AutoScaling = &nodepools.NodePoolAutoScaling{Min: pointer.To(input.AutoScaling[0].Min), Max: pointer.To(input.AutoScaling[0].Max)}
	} else {
		props.Replicas = pointer.To(input.Replicas)
	}
	if len(input.Version) > 0 {
		props.Version = &nodepools.NodePoolVersionProfile{Id: input.Version[0].Id, ChannelGroup: pointer.To(input.Version[0].ChannelGroup)}
	}
	if len(platform.OsDisk) > 0 {
		disk := platform.OsDisk[0]
		props.Platform.OsDisk = &nodepools.OsDiskProfile{}
		if disk.SizeGb > 0 {
			props.Platform.OsDisk.SizeGiB = pointer.To(disk.SizeGb)
		}
		if disk.DiskStorageAccountType != "" {
			props.Platform.OsDisk.DiskStorageAccountType = pointer.ToEnum[nodepools.DiskStorageAccountType](disk.DiskStorageAccountType)
		}
		if disk.DiskType != "" {
			props.Platform.OsDisk.DiskType = pointer.ToEnum[nodepools.OsDiskType](disk.DiskType)
		}
		if disk.EncryptionSetId != "" {
			props.Platform.OsDisk.EncryptionSetId = pointer.To(disk.EncryptionSetId)
		}
	}
	labels := make([]nodepools.Label, 0, len(input.Labels))
	keys := make([]string, 0, len(input.Labels))
	for key := range input.Labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		labels = append(labels, nodepools.Label{Key: key, Value: pointer.To(input.Labels[key])})
	}
	props.Labels = &labels
	taints := make([]nodepools.Taint, 0, len(input.Taints))
	for _, taint := range input.Taints {
		taints = append(taints, nodepools.Taint{Key: taint.Key, Value: pointer.To(taint.Value), Effect: nodepools.Effect(taint.Effect)})
	}
	props.Taints = &taints
	return &nodepools.NodePool{Location: location.Normalize(parentLocation), Properties: props, Tags: pointer.To(input.Tags)}, nil
}

func flattenHcpOpenShiftNodePool(id nodepools.NodePoolId, input nodepools.NodePool) (*RedHatHcpOpenShiftNodePoolModel, error) {
	parent := hcpopenshiftclusters.NewHcpOpenShiftClusterID(id.SubscriptionId, id.ResourceGroupName, id.HcpOpenShiftClusterName)
	state := &RedHatHcpOpenShiftNodePoolModel{Name: id.NodePoolName, ClusterId: parent.ID(), Tags: pointer.From(input.Tags), Labels: make(map[string]string)}
	if props := input.Properties; props != nil {
		platform := HcpOpenShiftNodePoolPlatform{
			VMSize:                 props.Platform.VMSize,
			AvailabilityZone:       pointer.From(props.Platform.AvailabilityZone),
			EnableEncryptionAtHost: pointer.From(props.Platform.EnableEncryptionAtHost),
		}
		state.AutoRepair = pointer.From(props.AutoRepair)
		state.Replicas = pointer.From(props.Replicas)
		state.NodeDrainTimeoutMinutes = pointer.From(props.NodeDrainTimeoutMinutes)
		if subnet := pointer.From(props.Platform.SubnetId); subnet != "" {
			parsed, err := commonids.ParseSubnetIDInsensitively(subnet)
			if err != nil {
				return nil, fmt.Errorf("parsing `platform.subnet_id`: %+v", err)
			}
			platform.SubnetId = parsed.ID()
		}
		if props.AutoScaling != nil {
			state.AutoScaling = []HcpOpenShiftNodePoolAutoscaling{{Min: pointer.From(props.AutoScaling.Min), Max: pointer.From(props.AutoScaling.Max)}}
		}
		if props.Version != nil {
			state.Version = []HcpOpenShiftNodePoolVersion{{Id: props.Version.Id, ChannelGroup: pointer.From(props.Version.ChannelGroup)}}
		}
		if disk := props.Platform.OsDisk; disk != nil {
			encryptionSetId := pointer.From(disk.EncryptionSetId)
			if encryptionSetId != "" {
				parsed, err := commonids.ParseDiskEncryptionSetIDInsensitively(encryptionSetId)
				if err != nil {
					return nil, fmt.Errorf("parsing `platform.os_disk.encryption_set_id`: %+v", err)
				}
				encryptionSetId = parsed.ID()
			}
			platform.OsDisk = []HcpOpenShiftNodePoolOsDisk{{SizeGb: pointer.From(disk.SizeGiB), DiskStorageAccountType: pointer.FromEnum(disk.DiskStorageAccountType), DiskType: pointer.FromEnum(disk.DiskType), EncryptionSetId: encryptionSetId}}
		}
		state.Platform = []HcpOpenShiftNodePoolPlatform{platform}
		for _, label := range pointer.From(props.Labels) {
			state.Labels[label.Key] = pointer.From(label.Value)
		}
		for _, taint := range pointer.From(props.Taints) {
			state.Taints = append(state.Taints, HcpOpenShiftNodePoolTaint{Key: taint.Key, Value: pointer.From(taint.Value), Effect: string(taint.Effect)})
		}
	}
	return state, nil
}
