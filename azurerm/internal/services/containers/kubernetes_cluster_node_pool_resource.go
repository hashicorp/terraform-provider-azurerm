package containers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-03-01/containerservice"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	computeValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/parse"
	containerValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

			"tags": tags.Schema(),

			"vm_size": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Optional
			"availability_zones": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
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
				ForceNew: true,
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

			"vnet_subnet_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"upgrade_settings": upgradeSettingsSchema(),
		},
	}
}

func resourceKubernetesClusterNodePoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	containersClient := meta.(*clients.Client).Containers
	clustersClient := containersClient.KubernetesClustersClient
	poolsClient := containersClient.AgentPoolsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	kubernetesClusterId, err := parse.ClusterID(d.Get("kubernetes_cluster_id").(string))
	if err != nil {
		return err
	}

	resourceGroup := kubernetesClusterId.ResourceGroup
	clusterName := kubernetesClusterId.ManagedClusterName
	name := d.Get("name").(string)

	log.Printf("[DEBUG] Retrieving Kubernetes Cluster %q (Resource Group %q)..", clusterName, resourceGroup)
	cluster, err := clustersClient.Get(ctx, resourceGroup, clusterName)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			return fmt.Errorf("Kubernetes Cluster %q was not found in Resource Group %q!", clusterName, resourceGroup)
		}

		return fmt.Errorf("retrieving existing Kubernetes Cluster %q (Resource Group %q): %+v", clusterName, resourceGroup, err)
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
		return fmt.Errorf("The Default Node Pool for Kubernetes Cluster %q (Resource Group %q) must be a VirtualMachineScaleSet to attach multiple node pools!", clusterName, resourceGroup)
	}

	existing, err := poolsClient.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Agent Pool %q (Kubernetes Cluster %q / Resource Group %q): %s", name, clusterName, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_kubernetes_cluster_node_pool", *existing.ID)
	}

	count := d.Get("node_count").(int)
	enableAutoScaling := d.Get("enable_auto_scaling").(bool)
	evictionPolicy := d.Get("eviction_policy").(string)
	mode := containerservice.AgentPoolMode(d.Get("mode").(string))
	osType := d.Get("os_type").(string)
	priority := d.Get("priority").(string)
	spotMaxPrice := d.Get("spot_max_price").(float64)
	t := d.Get("tags").(map[string]interface{})
	vmSize := d.Get("vm_size").(string)
	enableHostEncryption := d.Get("enable_host_encryption").(bool)

	profile := containerservice.ManagedClusterAgentPoolProfileProperties{
		OsType:                 containerservice.OSType(osType),
		EnableAutoScaling:      utils.Bool(enableAutoScaling),
		EnableFIPS:             utils.Bool(d.Get("fips_enabled").(bool)),
		EnableNodePublicIP:     utils.Bool(d.Get("enable_node_public_ip").(bool)),
		KubeletDiskType:        containerservice.KubeletDiskType(d.Get("kubelet_disk_type").(string)),
		Mode:                   mode,
		ScaleSetPriority:       containerservice.ScaleSetPriority(priority),
		Tags:                   tags.Expand(t),
		Type:                   containerservice.AgentPoolTypeVirtualMachineScaleSets,
		VMSize:                 utils.String(vmSize),
		EnableEncryptionAtHost: utils.Bool(enableHostEncryption),
		UpgradeSettings:        expandUpgradeSettings(d.Get("upgrade_settings").([]interface{})),

		// this must always be sent during creation, but is optional for auto-scaled clusters during update
		Count: utils.Int32(int32(count)),
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
		if err := validateNodePoolSupportsVersion(ctx, containersClient, resourceGroup, clusterName, name, orchestratorVersion); err != nil {
			return err
		}

		profile.OrchestratorVersion = utils.String(orchestratorVersion)
	}

	availabilityZonesRaw := d.Get("availability_zones").([]interface{})
	if availabilityZones := utils.ExpandStringSlice(availabilityZonesRaw); len(*availabilityZones) > 0 {
		profile.AvailabilityZones = availabilityZones
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

	if vnetSubnetID := d.Get("vnet_subnet_id").(string); vnetSubnetID != "" {
		profile.VnetSubnetID = utils.String(vnetSubnetID)
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
		Name:                                     &name,
		ManagedClusterAgentPoolProfileProperties: &profile,
	}

	future, err := poolsClient.CreateOrUpdate(ctx, resourceGroup, clusterName, name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating Managed Kubernetes Cluster Node Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, poolsClient.Client); err != nil {
		return fmt.Errorf("waiting for completion of Managed Kubernetes Cluster Node Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := poolsClient.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		return fmt.Errorf("retrieving Managed Kubernetes Cluster Node Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Managed Kubernetes Cluster Node Pool %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*read.ID)

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

	log.Printf("[DEBUG] Retrieving existing Node Pool %q (Kubernetes Cluster %q / Resource Group %q)..", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup)
	existing, err := client.Get(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Node Pool %q was not found in Managed Kubernetes Cluster %q / Resource Group %q!", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup)
		}

		return fmt.Errorf("retrieving Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): %+v", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup, err)
	}
	if existing.ManagedClusterAgentPoolProfileProperties == nil {
		return fmt.Errorf("retrieving Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): `properties` was nil", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup)
	}

	props := existing.ManagedClusterAgentPoolProfileProperties

	// store the existing value should the user have opted to ignore it
	enableAutoScaling := false
	if props.EnableAutoScaling != nil {
		enableAutoScaling = *props.EnableAutoScaling
	}

	log.Printf("[DEBUG] Determining delta for existing Node Pool %q (Kubernetes Cluster %q / Resource Group %q)..", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup)

	// delta patching
	if d.HasChange("availability_zones") {
		availabilityZonesRaw := d.Get("availability_zones").([]interface{})
		availabilityZones := utils.ExpandStringSlice(availabilityZonesRaw)
		props.AvailabilityZones = availabilityZones
	}

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

	if d.HasChange("max_count") {
		props.MaxCount = utils.Int32(int32(d.Get("max_count").(int)))
	}

	if d.HasChange("mode") {
		props.Mode = containerservice.AgentPoolMode(d.Get("mode").(string))
	}

	if d.HasChange("min_count") {
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

		orchestratorVersion := d.Get("orchestrator_version").(string)
		if err := validateNodePoolSupportsVersion(ctx, containersClient, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName, orchestratorVersion); err != nil {
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

	log.Printf("[DEBUG] Updating existing Node Pool %q (Kubernetes Cluster %q / Resource Group %q)..", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup)
	existing.ManagedClusterAgentPoolProfileProperties = props
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName, existing)
	if err != nil {
		return fmt.Errorf("updating Node Pool %q (Kubernetes Cluster %q / Resource Group %q): %+v", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Node Pool %q (Kubernetes Cluster %q / Resource Group %q): %+v", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup, err)
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
	cluster, err := clustersClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			log.Printf("[DEBUG] Managed Kubernetes Cluster %q was not found in Resource Group %q - removing from state!", id.ManagedClusterName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
	}

	resp, err := poolsClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Node Pool %q was not found in Managed Kubernetes Cluster %q / Resource Group %q - removing from state!", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): %+v", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup, err)
	}

	d.Set("name", id.AgentPoolName)
	d.Set("kubernetes_cluster_id", cluster.ID)

	if props := resp.ManagedClusterAgentPoolProfileProperties; props != nil {
		if err := d.Set("availability_zones", utils.FlattenStringSlice(props.AvailabilityZones)); err != nil {
			return fmt.Errorf("setting `availability_zones`: %+v", err)
		}

		d.Set("enable_auto_scaling", props.EnableAutoScaling)
		d.Set("enable_node_public_ip", props.EnableNodePublicIP)
		d.Set("enable_host_encryption", props.EnableEncryptionAtHost)
		d.Set("fips_enabled", props.EnableFIPS)
		d.Set("kubelet_disk_type", string(props.KubeletDiskType))

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

	future, err := client.Delete(ctx, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName)
	if err != nil {
		return fmt.Errorf("deleting Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): %+v", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): %+v", id.AgentPoolName, id.ManagedClusterName, id.ResourceGroup, err)
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
