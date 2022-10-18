package containers

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-08-02-preview/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-08-02-preview/managedclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKubernetesClusterNodePool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKubernetesClusterNodePoolCreate,
		Read:   resourceKubernetesClusterNodePoolRead,
		Update: resourceKubernetesClusterNodePoolUpdate,
		Delete: resourceKubernetesClusterNodePoolDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NodePoolID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KubernetesClusterNodePoolV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: containerValidate.KubernetesAgentPoolName,
			},

			"kubernetes_cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: containerValidate.ClusterID,
			},

			"node_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},

			"tags": commonschema.Tags(),

			"vm_size": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"host_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.HostGroupID,
			},

			// Optional
			"capacity_reservation_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.CapacityReservationGroupID,
			},

			"enable_auto_scaling": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"enable_host_encryption": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"enable_node_public_ip": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"eviction_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(agentpools.ScaleSetEvictionPolicyDelete),
					string(agentpools.ScaleSetEvictionPolicyDeallocate),
				}, false),
			},

			"kubelet_config": schemaNodePoolKubeletConfig(),

			"linux_os_config": schemaNodePoolLinuxOSConfig(),

			"fips_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"kubelet_disk_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(agentpools.KubeletDiskTypeOS),
					string(agentpools.KubeletDiskTypeTemporary),
				}, false),
			},

			"max_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},

			"max_pods": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"message_of_the_day": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(agentpools.AgentPoolModeUser),
				ValidateFunc: validation.StringInSlice([]string{
					string(agentpools.AgentPoolModeSystem),
					string(agentpools.AgentPoolModeUser),
				}, false),
			},

			"min_count": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				// NOTE: rather than setting `0` users should instead pass `null` here
				ValidateFunc: validation.IntBetween(0, 1000),
			},

			"node_labels": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"node_public_ip_prefix_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"enable_node_public_ip"},
			},

			// Node Taints control the behaviour of the Node Pool, as such they should not be computed and
			// must be specified/reconciled as required
			"node_taints": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"orchestrator_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"os_disk_size_gb": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"os_disk_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  agentpools.OSDiskTypeManaged,
				ValidateFunc: validation.StringInSlice([]string{
					string(agentpools.OSDiskTypeEphemeral),
					string(agentpools.OSDiskTypeManaged),
				}, false),
			},

			"os_sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true, // defaults to Ubuntu if using Linux
				ValidateFunc: validation.StringInSlice([]string{
					string(agentpools.OSSKUUbuntu),
					string(agentpools.OSSKUCBLMariner),
				}, false),
			},

			"os_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(agentpools.OSTypeLinux),
				ValidateFunc: validation.StringInSlice([]string{
					string(agentpools.OSTypeLinux),
					string(agentpools.OSTypeWindows),
				}, false),
			},

			"pod_subnet_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.SubnetID,
			},

			"priority": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(agentpools.ScaleSetPriorityRegular),
				ValidateFunc: validation.StringInSlice([]string{
					string(agentpools.ScaleSetPriorityRegular),
					string(agentpools.ScaleSetPrioritySpot),
				}, false),
			},

			"proximity_placement_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
			},

			"spot_max_price": {
				Type:         pluginsdk.TypeFloat,
				Optional:     true,
				ForceNew:     true,
				Default:      -1.0,
				ValidateFunc: computeValidate.SpotMaxPrice,
			},

			"scale_down_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(agentpools.ScaleDownModeDelete),
				ValidateFunc: validation.StringInSlice([]string{
					string(agentpools.ScaleDownModeDeallocate),
					string(agentpools.ScaleDownModeDelete),
				}, false),
			},

			"ultra_ssd_enabled": {
				Type:     pluginsdk.TypeBool,
				ForceNew: true,
				Default:  false,
				Optional: true,
			},

			"vnet_subnet_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"upgrade_settings": upgradeSettingsSchema(),

			"workload_runtime": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(agentpools.WorkloadRuntimeOCIContainer),
					string(agentpools.WorkloadRuntimeWasmWasi),
				}, false),
			},
			"zones": commonschema.ZonesMultipleOptionalForceNew(),
		},
	}
}

func resourceKubernetesClusterNodePoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	containersClient := meta.(*clients.Client).Containers
	clustersClient := meta.(*clients.Client).Containers.ManagedClustersClient
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId, err := agentpools.ParseManagedClusterID(d.Get("kubernetes_cluster_id").(string))
	mcClusterId, err := managedclusters.ParseManagedClusterID(d.Get("kubernetes_cluster_id").(string))
	if err != nil {
		return err
	}

	id := agentpools.NewAgentPoolID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.ResourceName, d.Get("name").(string))

	log.Printf("[DEBUG] Retrieving %s...", *clusterId)
	resp, err := clustersClient.Get(ctx, *mcClusterId)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", *clusterId)
		}

		return fmt.Errorf("retrieving %s: %+v", *clusterId, err)
	}

	cluster := resp.Model

	// try to provide a more helpful error here
	defaultPoolIsVMSS := false
	if props := cluster.Properties; props != nil {
		if pools := props.AgentPoolProfiles; pools != nil {
			for _, p := range *pools {
				if *(p).Type == managedclusters.AgentPoolTypeVirtualMachineScaleSets {
					defaultPoolIsVMSS = true
					break
				}
			}
		}
	}
	if !defaultPoolIsVMSS {
		return fmt.Errorf("multiple node pools are only supported when the Default Node Pool uses a VMScaleSet (but %s doesn't)", *clusterId)
	}

	existing, err := poolsClient.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_kubernetes_cluster_node_pool", id.ID())
	}

	count := d.Get("node_count").(int)
	enableAutoScaling := d.Get("enable_auto_scaling").(bool)
	evictionPolicy := agentpools.ScaleSetEvictionPolicy(d.Get("eviction_policy").(string))
	mode := agentpools.AgentPoolMode(d.Get("mode").(string))
	osType := agentpools.OSType(d.Get("os_type").(string))
	priority := agentpools.ScaleSetPriority(d.Get("priority").(string))
	spotMaxPrice := d.Get("spot_max_price").(float64)
	kubeletDiskType := agentpools.KubeletDiskType(d.Get("kubelet_disk_type").(string))
	t := d.Get("tags").(map[string]interface{})
	poolType := agentpools.AgentPoolTypeVirtualMachineScaleSets
	profile := agentpools.ManagedClusterAgentPoolProfileProperties{
		OsType:                 &osType,
		EnableAutoScaling:      utils.Bool(enableAutoScaling),
		EnableFIPS:             utils.Bool(d.Get("fips_enabled").(bool)),
		EnableEncryptionAtHost: utils.Bool(d.Get("enable_host_encryption").(bool)),
		EnableUltraSSD:         utils.Bool(d.Get("ultra_ssd_enabled").(bool)),
		EnableNodePublicIP:     utils.Bool(d.Get("enable_node_public_ip").(bool)),
		KubeletDiskType:        &kubeletDiskType,
		Mode:                   &mode,
		ScaleSetPriority:       &priority,
		Tags:                   tags.Expand(t),
		Type:                   &poolType,
		VmSize:                 utils.String(d.Get("vm_size").(string)),
		UpgradeSettings:        expandUpgradeSettingsForNodepool(d.Get("upgrade_settings").([]interface{})),

		// this must always be sent during creation, but is optional for auto-scaled clusters during update
		Count: utils.Int64(int64(count)),
	}

	if osSku := agentpools.OSSKU(d.Get("os_sku").(string)); osSku != "" {
		profile.OsSKU = &osSku
	}

	if scaleDownMode := agentpools.ScaleDownMode(d.Get("scale_down_mode").(string)); scaleDownMode != "" {
		profile.ScaleDownMode = &scaleDownMode
	}
	if workloadRuntime := agentpools.WorkloadRuntime(d.Get("workload_runtime").(string)); workloadRuntime != "" {
		profile.WorkloadRuntime = &workloadRuntime
	}

	if priority == agentpools.ScaleSetPrioritySpot {
		profile.ScaleSetEvictionPolicy = &evictionPolicy
		profile.SpotMaxPrice = utils.Float(spotMaxPrice)
	} else {
		if evictionPolicy != "" {
			return fmt.Errorf("`eviction_policy` can only be set when `priority` is set to `Spot`")
		}

		if spotMaxPrice != -1.0 {
			return fmt.Errorf("`spot_max_price` can only be set when `priority` is set to `Spot`")
		}
	}

	orchestratorVersion := d.Get("orchestrator_version").(string)
	if orchestratorVersion != "" {
		if err := validateNodePoolSupportsVersion(ctx, containersClient, "", *clusterId, id.AgentPoolName, orchestratorVersion); err != nil {
			return err
		}

		profile.OrchestratorVersion = utils.String(orchestratorVersion)
	}

	zones := zones.Expand(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		profile.AvailabilityZones = &zones
	}

	if maxPods := int64(d.Get("max_pods").(int)); maxPods > 0 {
		profile.MaxPods = utils.Int64(maxPods)
	}

	nodeLabelsRaw := d.Get("node_labels").(map[string]interface{})
	nodeLabels := make(map[string]string)
	for k, v := range nodeLabelsRaw {
		nodeLabels[k] = v.(string)
	}
	if len(nodeLabels) > 0 {
		profile.NodeLabels = &nodeLabels
	}

	if nodePublicIPPrefixID := d.Get("node_public_ip_prefix_id").(string); nodePublicIPPrefixID != "" {
		profile.NodePublicIPPrefixID = utils.String(nodePublicIPPrefixID)
	}

	nodeTaintsRaw := d.Get("node_taints").([]interface{})
	if nodeTaints := utils.ExpandStringSlice(nodeTaintsRaw); len(*nodeTaints) > 0 {
		profile.NodeTaints = nodeTaints
	}

	if v := d.Get("message_of_the_day").(string); v != "" {
		if *profile.OsType == agentpools.OSTypeWindows {
			return fmt.Errorf("`message_of_the_day` cannot be specified for Windows nodes and must be a static string (i.e. will be printed raw and not executed as a script)")
		}
		messageOfTheDayEncoded := base64.StdEncoding.EncodeToString([]byte(v))
		profile.MessageOfTheDay = &messageOfTheDayEncoded
	}

	if osDiskSizeGB := d.Get("os_disk_size_gb").(int); osDiskSizeGB > 0 {
		profile.OsDiskSizeGB = utils.Int64(int64(osDiskSizeGB))
	}

	proximityPlacementGroupId := d.Get("proximity_placement_group_id").(string)
	if proximityPlacementGroupId != "" {
		profile.ProximityPlacementGroupID = &proximityPlacementGroupId
	}

	if osDiskType := agentpools.OSDiskType(d.Get("os_disk_type").(string)); osDiskType != "" {
		profile.OsDiskType = &osDiskType
	}

	if podSubnetID := d.Get("pod_subnet_id").(string); podSubnetID != "" {
		profile.PodSubnetID = utils.String(podSubnetID)
	}

	if vnetSubnetID := d.Get("vnet_subnet_id").(string); vnetSubnetID != "" {
		profile.VnetSubnetID = utils.String(vnetSubnetID)
	}

	if hostGroupID := d.Get("host_group_id").(string); hostGroupID != "" {
		profile.HostGroupID = utils.String(hostGroupID)
	}

	if capacityReservationGroupId := d.Get("capacity_reservation_group_id").(string); capacityReservationGroupId != "" {
		profile.CapacityReservationGroupID = utils.String(capacityReservationGroupId)
	}

	maxCount := d.Get("max_count").(int)
	minCount := d.Get("min_count").(int)

	if enableAutoScaling {
		// handle count being optional
		if count == 0 {
			profile.Count = utils.Int64(int64(minCount))
		}

		if maxCount >= 0 {
			profile.MaxCount = utils.Int64(int64(maxCount))
		} else {
			return fmt.Errorf("`max_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount >= 0 {
			profile.MinCount = utils.Int64(int64(minCount))
		} else {
			return fmt.Errorf("`min_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > maxCount {
			return fmt.Errorf("`max_count` must be >= `min_count`")
		}
	} else if minCount > 0 || maxCount > 0 {
		return fmt.Errorf("`max_count` and `min_count` must be set to `null` when enable_auto_scaling is set to `false`")
	}

	if kubeletConfig := d.Get("kubelet_config").([]interface{}); len(kubeletConfig) > 0 {
		profile.KubeletConfig = expandAgentPoolKubeletConfig(kubeletConfig)
	}

	if linuxOSConfig := d.Get("linux_os_config").([]interface{}); len(linuxOSConfig) > 0 {
		if osType != agentpools.OSTypeLinux {
			return fmt.Errorf("`linux_os_config` can only be configured when `os_type` is set to `linux`")
		}
		linuxOSConfig, err := expandAgentPoolLinuxOSConfig(linuxOSConfig)
		if err != nil {
			return err
		}
		profile.LinuxOSConfig = linuxOSConfig
	}

	parameters := agentpools.AgentPool{
		Name:       utils.String(id.AgentPoolName),
		Properties: &profile,
	}

	_, err = poolsClient.CreateOrUpdate(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKubernetesClusterNodePoolRead(d, meta)
}

func resourceKubernetesClusterNodePoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	containersClient := meta.(*clients.Client).Containers
	client := containersClient.AgentPoolsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agentpools.ParseAgentPoolID(d.Id())
	if err != nil {
		return err
	}

	clusterId, err := agentpools.ParseManagedClusterID(d.Get("kubernetes_cluster_id").(string))
	if err != nil {
		return err
	}

	d.Partial(true)

	log.Printf("[DEBUG] Retrieving existing %s..", *id)
	existing, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	props := existing.Model.Properties

	// store the existing value should the user have opted to ignore it
	enableAutoScaling := false
	if props.EnableAutoScaling != nil {
		enableAutoScaling = *props.EnableAutoScaling
	}

	log.Printf("[DEBUG] Determining delta for existing %s..", *id)

	// delta patching
	if d.HasChange("enable_auto_scaling") {
		enableAutoScaling = d.Get("enable_auto_scaling").(bool)
		props.EnableAutoScaling = utils.Bool(enableAutoScaling)
	}

	if d.HasChange("enable_host_encryption") {
		props.EnableEncryptionAtHost = utils.Bool(d.Get("enable_host_encryption").(bool))
	}

	if d.HasChange("enable_node_public_ip") {
		props.EnableNodePublicIP = utils.Bool(d.Get("enable_node_public_ip").(bool))
	}

	if d.HasChange("max_count") || d.Get("enable_auto_scaling").(bool) {
		props.MaxCount = utils.Int64(int64(d.Get("max_count").(int)))
	}

	if d.HasChange("mode") {
		*props.Mode = agentpools.AgentPoolMode(d.Get("mode").(string))
	}

	if d.HasChange("min_count") || d.Get("enable_auto_scaling").(bool) {
		props.MinCount = utils.Int64(int64(d.Get("min_count").(int)))
	}

	if d.HasChange("node_count") {
		props.Count = utils.Int64(int64(d.Get("node_count").(int)))
	}

	if d.HasChange("node_public_ip_prefix_id") {
		props.NodePublicIPPrefixID = utils.String(d.Get("node_public_ip_prefix_id").(string))
	}

	if d.HasChange("orchestrator_version") {
		existingNodePool, err := client.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving Node Pool %s: %+v", *id, err)
		}
		orchestratorVersion := d.Get("orchestrator_version").(string)
		currentOrchestratorVersion := ""
		if v := existingNodePool.Model.Properties.OrchestratorVersion; v != nil {
			currentOrchestratorVersion = *v
		}
		if err := validateNodePoolSupportsVersion(ctx, containersClient, currentOrchestratorVersion, *clusterId, id.AgentPoolName, orchestratorVersion); err != nil {
			return err
		}

		props.OrchestratorVersion = utils.String(orchestratorVersion)
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})
		props.Tags = tags.Expand(t)
	}

	if d.HasChange("upgrade_settings") {
		upgradeSettingsRaw := d.Get("upgrade_settings").([]interface{})
		props.UpgradeSettings = expandUpgradeSettingsForNodepool(upgradeSettingsRaw)
	}

	if d.HasChange("scale_down_mode") {
		*props.ScaleDownMode = agentpools.ScaleDownMode(d.Get("scale_down_mode").(string))
	}
	if d.HasChange("workload_runtime") {
		*props.WorkloadRuntime = agentpools.WorkloadRuntime(d.Get("workload_runtime").(string))
	}

	if d.HasChange("node_labels") {
		result := make(map[string]string)
		for k, v := range d.Get("node_labels").(map[string]interface{}) {
			result[k] = v.(string)
		}
		props.NodeLabels = &result
	}

	// validate the auto-scale fields are both set/unset to prevent a continual diff
	maxCount := 0
	if props.MaxCount != nil {
		maxCount = int(*props.MaxCount)
	}
	minCount := 0
	if props.MinCount != nil {
		minCount = int(*props.MinCount)
	}
	if enableAutoScaling {
		if maxCount == 0 {
			return fmt.Errorf("`max_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > maxCount {
			return fmt.Errorf("`max_count` must be >= `min_count`")
		}
	} else {
		if minCount > 0 || maxCount > 0 {
			return fmt.Errorf("`max_count` and `min_count` must be set to `nil` when enable_auto_scaling is set to `false`")
		}

		// @tombuildsstuff: as of API version 2019-11-01 we need to explicitly nil these out
		props.MaxCount = nil
		props.MinCount = nil
	}

	log.Printf("[DEBUG] Updating existing %s..", *id)
	existing.Model.Properties = props

	_, err = client.CreateOrUpdate(ctx, *id, *existing.Model)
	if err != nil {
		return fmt.Errorf("updating Node Pool %s: %+v", *id, err)
	}

	d.Partial(false)

	return resourceKubernetesClusterNodePoolRead(d, meta)
}

func resourceKubernetesClusterNodePoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).Containers.ManagedClustersClient
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agentpools.ParseAgentPoolID(d.Id())
	if err != nil {
		return err
	}

	// if the parent cluster doesn't exist then the node pool won't
	clusterId := managedclusters.NewManagedClusterID(id.SubscriptionId, id.ResourceGroupName, id.ResourceName)
	cluster, err := clustersClient.Get(ctx, clusterId)
	if err != nil {
		if response.WasNotFound(cluster.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", clusterId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", clusterId, err)
	}

	resp, err := poolsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AgentPoolName)
	d.Set("kubernetes_cluster_id", clusterId.ID())

	if props := resp.Model.Properties; props != nil {
		d.Set("zones", zones.Flatten(props.AvailabilityZones))
		d.Set("enable_auto_scaling", props.EnableAutoScaling)
		d.Set("enable_node_public_ip", props.EnableNodePublicIP)
		d.Set("enable_host_encryption", props.EnableEncryptionAtHost)
		d.Set("fips_enabled", props.EnableFIPS)
		d.Set("ultra_ssd_enabled", props.EnableUltraSSD)
		d.Set("kubelet_disk_type", string(*props.KubeletDiskType))
		scaleDownMode := string(agentpools.ScaleDownModeDelete)
		if v := *props.ScaleDownMode; v != "" {
			scaleDownMode = string(v)
		}
		d.Set("scale_down_mode", scaleDownMode)
		d.Set("workload_runtime", string(*props.WorkloadRuntime))

		evictionPolicy := ""
		if *props.ScaleSetEvictionPolicy != "" {
			evictionPolicy = string(*props.ScaleSetEvictionPolicy)
		}
		d.Set("eviction_policy", evictionPolicy)

		if err := d.Set("kubelet_config", flattenAgentPoolKubeletConfig(props.KubeletConfig)); err != nil {
			return fmt.Errorf("setting `kubelet_config`: %+v", err)
		}

		linuxOSConfig, err := flattenAgentPoolLinuxOSConfig(props.LinuxOSConfig)
		if err != nil {
			return err
		}
		if err := d.Set("linux_os_config", linuxOSConfig); err != nil {
			return fmt.Errorf("setting `linux_os_config`: %+v", err)
		}

		maxCount := 0
		if props.MaxCount != nil {
			maxCount = int(*props.MaxCount)
		}
		d.Set("max_count", maxCount)

		messageOfTheDay := ""
		if props.MessageOfTheDay != nil {
			messageOfTheDayDecoded, err := base64.StdEncoding.DecodeString(*props.MessageOfTheDay)
			if err != nil {
				return fmt.Errorf("setting `message_of_the_day`: %+v", err)
			}
			messageOfTheDay = string(messageOfTheDayDecoded)
		}
		d.Set("message_of_the_day", messageOfTheDay)

		maxPods := 0
		if props.MaxPods != nil {
			maxPods = int(*props.MaxPods)
		}
		d.Set("max_pods", maxPods)

		minCount := 0
		if props.MinCount != nil {
			minCount = int(*props.MinCount)
		}
		d.Set("min_count", minCount)

		mode := string(agentpools.AgentPoolModeUser)
		if *props.Mode != "" {
			mode = string(*props.Mode)
		}
		d.Set("mode", mode)

		count := 0
		if props.Count != nil {
			count = int(*props.Count)
		}
		d.Set("node_count", count)

		if err := d.Set("node_labels", props.NodeLabels); err != nil {
			return fmt.Errorf("setting `node_labels`: %+v", err)
		}

		d.Set("node_public_ip_prefix_id", props.NodePublicIPPrefixID)

		if err := d.Set("node_taints", utils.FlattenStringSlice(props.NodeTaints)); err != nil {
			return fmt.Errorf("setting `node_taints`: %+v", err)
		}

		// NOTE: workaround for migration from 2022-01-02-preview (<3.12.0) to 2022-03-02-preview (>=3.12.0). Before terraform apply is run against the new API, Azure will respond only with currentOrchestratorVersion, orchestratorVersion will be absent. More details: https://github.com/hashicorp/terraform-provider-azurerm/issues/17833#issuecomment-1227583353
		if props.OrchestratorVersion != nil {
			d.Set("orchestrator_version", props.OrchestratorVersion)
		} else {
			d.Set("orchestrator_version", props.CurrentOrchestratorVersion)
		}

		osDiskSizeGB := 0
		if props.OsDiskSizeGB != nil {
			osDiskSizeGB = int(*props.OsDiskSizeGB)
		}
		d.Set("os_disk_size_gb", osDiskSizeGB)

		osDiskType := agentpools.OSDiskTypeManaged
		if *props.OsDiskType != "" {
			osDiskType = *props.OsDiskType
		}
		d.Set("os_disk_type", osDiskType)
		d.Set("os_type", string(*props.OsType))
		d.Set("os_sku", string(*props.OsSKU))
		d.Set("pod_subnet_id", props.PodSubnetID)

		// not returned from the API if not Spot
		priority := string(agentpools.ScaleSetPriorityRegular)
		if *props.ScaleSetPriority != "" {
			priority = string(*props.ScaleSetPriority)
		}
		d.Set("priority", priority)

		d.Set("proximity_placement_group_id", props.ProximityPlacementGroupID)

		spotMaxPrice := -1.0
		if props.SpotMaxPrice != nil {
			spotMaxPrice = *props.SpotMaxPrice
		}
		d.Set("spot_max_price", spotMaxPrice)

		d.Set("vnet_subnet_id", props.VnetSubnetID)
		d.Set("vm_size", props.VmSize)
		d.Set("host_group_id", props.HostGroupID)
		d.Set("capacity_reservation_group_id", props.CapacityReservationGroupID)

		if err := d.Set("upgrade_settings", flattenUpgradeSettingsFromNodePool(props.UpgradeSettings)); err != nil {
			return fmt.Errorf("setting `upgrade_settings`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Model.Properties.Tags)
}

func resourceKubernetesClusterNodePoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agentpools.ParseAgentPoolID(d.Id())
	if err != nil {
		return err
	}

	ignorePodDisruptionBudget := agentpools.DeleteOperationOptions{
		IgnorePodDisruptionBudget: utils.Bool(true),
	}

	err = client.DeleteThenPoll(ctx, *id, ignorePodDisruptionBudget)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func upgradeSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_surge": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
			},
		},
	}
}

func upgradeSettingsForDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_surge": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func expandUpgradeSettingsForNodepool(input []interface{}) *agentpools.AgentPoolUpgradeSettings {
	setting := &agentpools.AgentPoolUpgradeSettings{}
	if len(input) == 0 || input[0] == nil {
		return setting
	}

	v := input[0].(map[string]interface{})
	if maxSurgeRaw := v["max_surge"].(string); maxSurgeRaw != "" {
		setting.MaxSurge = utils.String(maxSurgeRaw)
	}
	return setting
}

func flattenUpgradeSettingsFromNodePool(input *agentpools.AgentPoolUpgradeSettings) []interface{} {
	maxSurge := ""
	if input != nil && input.MaxSurge != nil {
		maxSurge = *input.MaxSurge
	}

	if maxSurge == "" {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"max_surge": maxSurge,
		},
	}
}
