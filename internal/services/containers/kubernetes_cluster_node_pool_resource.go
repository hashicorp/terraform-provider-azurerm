package containers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerservice/mgmt/2022-03-02-preview/containerservice"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
					string(containerservice.ScaleSetEvictionPolicyDelete),
					string(containerservice.ScaleSetEvictionPolicyDeallocate),
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
					string(containerservice.KubeletDiskTypeOS),
					string(containerservice.KubeletDiskTypeTemporary),
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

			"mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(containerservice.AgentPoolModeUser),
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.AgentPoolModeSystem),
					string(containerservice.AgentPoolModeUser),
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
				Default:  containerservice.OSDiskTypeManaged,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.OSDiskTypeEphemeral),
					string(containerservice.OSDiskTypeManaged),
				}, false),
			},

			"os_sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true, // defaults to Ubuntu if using Linux
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.OSSKUUbuntu),
					string(containerservice.OSSKUCBLMariner),
				}, false),
			},

			"os_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(containerservice.OSTypeLinux),
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.OSTypeLinux),
					string(containerservice.OSTypeWindows),
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
				Default:  string(containerservice.ScaleSetPriorityRegular),
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.ScaleSetPriorityRegular),
					string(containerservice.ScaleSetPrioritySpot),
				}, false),
			},

			"proximity_placement_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.ProximityPlacementGroupID,
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
				Default:  string(containerservice.ScaleDownModeDelete),
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.ScaleDownModeDeallocate),
					string(containerservice.ScaleDownModeDelete),
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
					string(containerservice.WorkloadRuntimeOCIContainer),
					string(containerservice.WorkloadRuntimeWasmWasi),
				}, false),
			},
			"zones": commonschema.ZonesMultipleOptionalForceNew(),
		},
	}
}

func resourceKubernetesClusterNodePoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	containersClient := meta.(*clients.Client).Containers
	clustersClient := containersClient.KubernetesClustersClient
	poolsClient := containersClient.AgentPoolsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId, err := parse.ClusterID(d.Get("kubernetes_cluster_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewNodePoolID(poolsClient.SubscriptionID, clusterId.ResourceGroup, clusterId.ManagedClusterName, d.Get("name").(string))

	log.Printf("[DEBUG] Retrieving %s...", *clusterId)
	cluster, err := clustersClient.Get(ctx, clusterId.ResourceGroup, clusterId.ManagedClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			return fmt.Errorf("%s was not found", *clusterId)
		}

		return fmt.Errorf("retrieving %s: %+v", *clusterId, err)
	}

	// try to provide a more helpful error here
	defaultPoolIsVMSS := false
	if props := cluster.ManagedClusterProperties; props != nil {
		if pools := props.AgentPoolProfiles; pools != nil {
			for _, p := range *pools {
				if p.Type == containerservice.AgentPoolTypeVirtualMachineScaleSets {
					defaultPoolIsVMSS = true
					break
				}
			}
		}
	}
	if !defaultPoolIsVMSS {
		return fmt.Errorf("multiple node pools are only supported when the Default Node Pool uses a VMScaleSet (but %s doesn't)", *clusterId)
	}

	existing, err := poolsClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_kubernetes_cluster_node_pool", id.ID())
	}

	count := d.Get("node_count").(int)
	enableAutoScaling := d.Get("enable_auto_scaling").(bool)
	evictionPolicy := d.Get("eviction_policy").(string)
	mode := containerservice.AgentPoolMode(d.Get("mode").(string))
	osType := d.Get("os_type").(string)
	priority := d.Get("priority").(string)
	spotMaxPrice := d.Get("spot_max_price").(float64)
	t := d.Get("tags").(map[string]interface{})

	profile := containerservice.ManagedClusterAgentPoolProfileProperties{
		OsType:                 containerservice.OSType(osType),
		EnableAutoScaling:      utils.Bool(enableAutoScaling),
		EnableFIPS:             utils.Bool(d.Get("fips_enabled").(bool)),
		EnableEncryptionAtHost: utils.Bool(d.Get("enable_host_encryption").(bool)),
		EnableUltraSSD:         utils.Bool(d.Get("ultra_ssd_enabled").(bool)),
		EnableNodePublicIP:     utils.Bool(d.Get("enable_node_public_ip").(bool)),
		KubeletDiskType:        containerservice.KubeletDiskType(d.Get("kubelet_disk_type").(string)),
		Mode:                   mode,
		ScaleSetPriority:       containerservice.ScaleSetPriority(priority),
		Tags:                   tags.Expand(t),
		Type:                   containerservice.AgentPoolTypeVirtualMachineScaleSets,
		VMSize:                 utils.String(d.Get("vm_size").(string)),
		UpgradeSettings:        expandUpgradeSettings(d.Get("upgrade_settings").([]interface{})),

		// this must always be sent during creation, but is optional for auto-scaled clusters during update
		Count: utils.Int32(int32(count)),
	}

	if osSku := d.Get("os_sku").(string); osSku != "" {
		profile.OsSKU = containerservice.OSSKU(osSku)
	}

	if scaleDownMode := d.Get("scale_down_mode").(string); scaleDownMode != "" {
		profile.ScaleDownMode = containerservice.ScaleDownMode(scaleDownMode)
	}
	if workloadRuntime := d.Get("workload_runtime").(string); workloadRuntime != "" {
		profile.WorkloadRuntime = containerservice.WorkloadRuntime(workloadRuntime)
	}

	if priority == string(containerservice.ScaleSetPrioritySpot) {
		profile.ScaleSetEvictionPolicy = containerservice.ScaleSetEvictionPolicy(evictionPolicy)
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
		if err := validateNodePoolSupportsVersion(ctx, containersClient, "", id, orchestratorVersion); err != nil {
			return err
		}

		profile.OrchestratorVersion = utils.String(orchestratorVersion)
	}

	zones := zones.Expand(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		profile.AvailabilityZones = &zones
	}

	if maxPods := int32(d.Get("max_pods").(int)); maxPods > 0 {
		profile.MaxPods = utils.Int32(maxPods)
	}

	nodeLabelsRaw := d.Get("node_labels").(map[string]interface{})
	if nodeLabels := utils.ExpandMapStringPtrString(nodeLabelsRaw); len(nodeLabels) > 0 {
		profile.NodeLabels = nodeLabels
	}

	if nodePublicIPPrefixID := d.Get("node_public_ip_prefix_id").(string); nodePublicIPPrefixID != "" {
		profile.NodePublicIPPrefixID = utils.String(nodePublicIPPrefixID)
	}

	nodeTaintsRaw := d.Get("node_taints").([]interface{})
	if nodeTaints := utils.ExpandStringSlice(nodeTaintsRaw); len(*nodeTaints) > 0 {
		profile.NodeTaints = nodeTaints
	}

	if osDiskSizeGB := d.Get("os_disk_size_gb").(int); osDiskSizeGB > 0 {
		profile.OsDiskSizeGB = utils.Int32(int32(osDiskSizeGB))
	}

	proximityPlacementGroupId := d.Get("proximity_placement_group_id").(string)
	if proximityPlacementGroupId != "" {
		profile.ProximityPlacementGroupID = &proximityPlacementGroupId
	}

	if osDiskType := d.Get("os_disk_type").(string); osDiskType != "" {
		profile.OsDiskType = containerservice.OSDiskType(osDiskType)
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
			profile.Count = utils.Int32(int32(minCount))
		}

		if maxCount >= 0 {
			profile.MaxCount = utils.Int32(int32(maxCount))
		} else {
			return fmt.Errorf("`max_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount >= 0 {
			profile.MinCount = utils.Int32(int32(minCount))
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
		if osType != string(containerservice.OSTypeLinux) {
			return fmt.Errorf("`linux_os_config` can only be configured when `os_type` is set to `linux`")
		}
		linuxOSConfig, err := expandAgentPoolLinuxOSConfig(linuxOSConfig)
		if err != nil {
			return err
		}
		profile.LinuxOSConfig = linuxOSConfig
	}

	parameters := containerservice.AgentPool{
		Name:                                     utils.String(id.AgentPoolName),
		ManagedClusterAgentPoolProfileProperties: &profile,
	}

	future, err := poolsClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, poolsClient.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKubernetesClusterNodePoolRead(d, meta)
}

func resourceKubernetesClusterNodePoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	containersClient := meta.(*clients.Client).Containers
	client := containersClient.AgentPoolsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NodePoolID(d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)

	log.Printf("[DEBUG] Retrieving existing %s..", *id)
	existing, err := client.Get(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.ManagedClusterAgentPoolProfileProperties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	props := existing.ManagedClusterAgentPoolProfileProperties

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
		props.MaxCount = utils.Int32(int32(d.Get("max_count").(int)))
	}

	if d.HasChange("mode") {
		props.Mode = containerservice.AgentPoolMode(d.Get("mode").(string))
	}

	if d.HasChange("min_count") || d.Get("enable_auto_scaling").(bool) {
		props.MinCount = utils.Int32(int32(d.Get("min_count").(int)))
	}

	if d.HasChange("node_count") {
		props.Count = utils.Int32(int32(d.Get("node_count").(int)))
	}

	if d.HasChange("node_public_ip_prefix_id") {
		props.NodePublicIPPrefixID = utils.String(d.Get("node_public_ip_prefix_id").(string))
	}

	if d.HasChange("orchestrator_version") {
		// Spot Node pool's can't be updated - Azure Docs: https://docs.microsoft.com/en-us/azure/aks/spot-node-pool
		//   > You can't upgrade a spot node pool since spot node pools can't guarantee cordon and drain.
		//   > You must replace your existing spot node pool with a new one to do operations such as upgrading
		//   > the Kubernetes version. To replace a spot node pool, create a new spot node pool with a different
		//   > version of Kubernetes, wait until its status is Ready, then remove the old node pool.
		if strings.EqualFold(string(props.ScaleSetPriority), string(containerservice.ScaleSetPrioritySpot)) {
			// ^ the Scale Set Priority isn't returned when Regular
			return fmt.Errorf("the Orchestrator Version cannot be updated when using a Spot Node Pool")
		}

		existingNodePool, err := client.Get(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName)
		if err != nil {
			return fmt.Errorf("retrieving Node Pool %s: %+v", *id, err)
		}
		orchestratorVersion := d.Get("orchestrator_version").(string)
		currentOrchestratorVersion := ""
		if v := existingNodePool.OrchestratorVersion; v != nil {
			currentOrchestratorVersion = *v
		}
		if err := validateNodePoolSupportsVersion(ctx, containersClient, currentOrchestratorVersion, *id, orchestratorVersion); err != nil {
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
		props.UpgradeSettings = expandUpgradeSettings(upgradeSettingsRaw)
	}

	if d.HasChange("scale_down_mode") {
		props.ScaleDownMode = containerservice.ScaleDownMode(d.Get("scale_down_mode").(string))
	}
	if d.HasChange("workload_runtime") {
		props.WorkloadRuntime = containerservice.WorkloadRuntime(d.Get("workload_runtime").(string))
	}

	if d.HasChange("node_labels") {
		props.NodeLabels = utils.ExpandMapStringPtrString(d.Get("node_labels").(map[string]interface{}))
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
	existing.ManagedClusterAgentPoolProfileProperties = props
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName, existing)
	if err != nil {
		return fmt.Errorf("updating Node Pool %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	d.Partial(false)

	return resourceKubernetesClusterNodePoolRead(d, meta)
}

func resourceKubernetesClusterNodePoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).Containers.KubernetesClustersClient
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NodePoolID(d.Id())
	if err != nil {
		return err
	}

	// if the parent cluster doesn't exist then the node pool won't
	clusterId := parse.NewClusterID(id.SubscriptionId, id.ResourceGroup, id.ManagedClusterName)
	cluster, err := clustersClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", clusterId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", clusterId, err)
	}

	resp, err := poolsClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %q was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AgentPoolName)
	d.Set("kubernetes_cluster_id", clusterId.ID())

	if props := resp.ManagedClusterAgentPoolProfileProperties; props != nil {
		d.Set("zones", zones.Flatten(props.AvailabilityZones))
		d.Set("enable_auto_scaling", props.EnableAutoScaling)
		d.Set("enable_node_public_ip", props.EnableNodePublicIP)
		d.Set("enable_host_encryption", props.EnableEncryptionAtHost)
		d.Set("fips_enabled", props.EnableFIPS)
		d.Set("ultra_ssd_enabled", props.EnableUltraSSD)
		d.Set("kubelet_disk_type", string(props.KubeletDiskType))
		scaleDownMode := string(containerservice.ScaleDownModeDelete)
		if v := props.ScaleDownMode; v != "" {
			scaleDownMode = string(v)
		}
		d.Set("scale_down_mode", scaleDownMode)
		d.Set("workload_runtime", string(props.WorkloadRuntime))

		evictionPolicy := ""
		if props.ScaleSetEvictionPolicy != "" {
			evictionPolicy = string(props.ScaleSetEvictionPolicy)
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

		mode := string(containerservice.AgentPoolModeUser)
		if props.Mode != "" {
			mode = string(props.Mode)
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

		d.Set("orchestrator_version", props.OrchestratorVersion)
		osDiskSizeGB := 0
		if props.OsDiskSizeGB != nil {
			osDiskSizeGB = int(*props.OsDiskSizeGB)
		}
		d.Set("os_disk_size_gb", osDiskSizeGB)

		osDiskType := containerservice.OSDiskTypeManaged
		if props.OsDiskType != "" {
			osDiskType = props.OsDiskType
		}
		d.Set("os_disk_type", osDiskType)
		d.Set("os_type", string(props.OsType))
		d.Set("os_sku", string(props.OsSKU))
		d.Set("pod_subnet_id", props.PodSubnetID)

		// not returned from the API if not Spot
		priority := string(containerservice.ScaleSetPriorityRegular)
		if props.ScaleSetPriority != "" {
			priority = string(props.ScaleSetPriority)
		}
		d.Set("priority", priority)

		d.Set("proximity_placement_group_id", props.ProximityPlacementGroupID)

		spotMaxPrice := -1.0
		if props.SpotMaxPrice != nil {
			spotMaxPrice = *props.SpotMaxPrice
		}
		d.Set("spot_max_price", spotMaxPrice)

		d.Set("vnet_subnet_id", props.VnetSubnetID)
		d.Set("vm_size", props.VMSize)
		d.Set("host_group_id", props.HostGroupID)
		d.Set("capacity_reservation_group_id", props.CapacityReservationGroupID)

		if err := d.Set("upgrade_settings", flattenUpgradeSettings(props.UpgradeSettings)); err != nil {
			return fmt.Errorf("setting `upgrade_settings`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceKubernetesClusterNodePoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NodePoolID(d.Id())
	if err != nil {
		return err
	}

	ignorePodDisruptionBudget := true

	future, err := client.Delete(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName, &ignorePodDisruptionBudget)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", id, err)
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

func expandUpgradeSettings(input []interface{}) *containerservice.AgentPoolUpgradeSettings {
	setting := &containerservice.AgentPoolUpgradeSettings{}
	if len(input) == 0 {
		return setting
	}

	v := input[0].(map[string]interface{})
	if maxSurgeRaw := v["max_surge"].(string); maxSurgeRaw != "" {
		setting.MaxSurge = utils.String(maxSurgeRaw)
	}
	return setting
}

func flattenUpgradeSettings(input *containerservice.AgentPoolUpgradeSettings) []interface{} {
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
