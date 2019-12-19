package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-10-01/containerservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKubernetesClusterNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKubernetesClusterNodePoolCreate,
		Read:   resourceArmKubernetesClusterNodePoolRead,
		Update: resourceArmKubernetesClusterNodePoolUpdate,
		Delete: resourceArmKubernetesClusterNodePoolDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := containers.ParseKubernetesNodePoolID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KubernetesAgentPoolName,
			},

			"kubernetes_cluster_id": containers.KubernetesClusterIDSchema(),

			"node_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"vm_size": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			// Optional
			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"enable_auto_scaling": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enable_node_public_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"max_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"max_pods": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"min_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"node_taints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"os_disk_size_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(containerservice.Linux),
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.Linux),
					string(containerservice.Windows),
				}, false),
			},

			"vnet_subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmKubernetesClusterNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).Containers.KubernetesClustersClient
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	kubernetesClusterId, err := containers.ParseKubernetesClusterID(d.Get("kubernetes_cluster_id").(string))
	if err != nil {
		return err
	}

	resourceGroup := kubernetesClusterId.ResourceGroup
	clusterName := kubernetesClusterId.Name
	name := d.Get("name").(string)

	log.Printf("[DEBUG] Retrieving Kubernetes Cluster %q (Resource Group %q)..", clusterName, resourceGroup)
	cluster, err := clustersClient.Get(ctx, resourceGroup, clusterName)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			return fmt.Errorf("Kubernetes Cluster %q was not found in Resource Group %q!", clusterName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving existing Kubernetes Cluster %q (Resource Group %q): %+v", clusterName, resourceGroup, err)
	}

	// try to provide a more helpful error here
	defaultPoolIsVMSS := false
	if props := cluster.ManagedClusterProperties; props != nil {
		if pools := props.AgentPoolProfiles; pools != nil {
			for _, p := range *pools {
				if p.Type == containerservice.VirtualMachineScaleSets {
					defaultPoolIsVMSS = true
					break
				}
			}
		}
	}
	if !defaultPoolIsVMSS {
		return fmt.Errorf("The Default Node Pool for Kubernetes Cluster %q (Resource Group %q) must be a VirtualMachineScaleSet to attach multiple node pools!", clusterName, resourceGroup)
	}

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := poolsClient.Get(ctx, resourceGroup, clusterName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Agent Pool %q (Kubernetes Cluster %q / Resource Group %q): %s", name, clusterName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_kubernetes_cluster_node_pool", *existing.ID)
		}
	}

	count := d.Get("node_count").(int)
	enableAutoScaling := d.Get("enable_auto_scaling").(bool)
	osType := d.Get("os_type").(string)
	vmSize := d.Get("vm_size").(string)

	profile := containerservice.ManagedClusterAgentPoolProfileProperties{
		OsType:             containerservice.OSType(osType),
		EnableAutoScaling:  utils.Bool(enableAutoScaling),
		EnableNodePublicIP: utils.Bool(d.Get("enable_node_public_ip").(bool)),
		Type:               containerservice.VirtualMachineScaleSets,
		VMSize:             containerservice.VMSizeTypes(vmSize),

		// this must always be sent during creation, but is optional for auto-scaled clusters during update
		Count: utils.Int32(int32(count)),
	}

	availabilityZonesRaw := d.Get("availability_zones").([]interface{})
	if availabilityZones := utils.ExpandStringSlice(availabilityZonesRaw); len(*availabilityZones) > 0 {
		profile.AvailabilityZones = availabilityZones
	}

	if maxPods := int32(d.Get("max_pods").(int)); maxPods > 0 {
		profile.MaxPods = utils.Int32(maxPods)
	}

	nodeTaintsRaw := d.Get("node_taints").([]interface{})
	if nodeTaints := utils.ExpandStringSlice(nodeTaintsRaw); len(*nodeTaints) > 0 {
		profile.NodeTaints = nodeTaints
	}

	if osDiskSizeGB := d.Get("os_disk_size_gb").(int); osDiskSizeGB > 0 {
		profile.OsDiskSizeGB = utils.Int32(int32(osDiskSizeGB))
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

		if maxCount > 0 {
			profile.MaxCount = utils.Int32(int32(maxCount))
		} else {
			return fmt.Errorf("`max_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > 0 {
			profile.MinCount = utils.Int32(int32(minCount))
		} else {
			return fmt.Errorf("`min_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > maxCount {
			return fmt.Errorf("`max_count` must be >= `min_count`")
		}
	} else if minCount > 0 || maxCount > 0 {
		return fmt.Errorf("`max_count` and `min_count` must be set to `0` when enable_auto_scaling is set to `false`")
	}

	parameters := containerservice.AgentPool{
		Name:                                     &name,
		ManagedClusterAgentPoolProfileProperties: &profile,
	}

	future, err := poolsClient.CreateOrUpdate(ctx, resourceGroup, clusterName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Managed Kubernetes Cluster Node Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, poolsClient.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Managed Kubernetes Cluster Node Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := poolsClient.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Managed Kubernetes Cluster Node Pool %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Managed Kubernetes Cluster Node Pool %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmKubernetesClusterNodePoolRead(d, meta)
}

func resourceArmKubernetesClusterNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseKubernetesNodePoolID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	clusterName := id.ClusterName
	name := id.Name

	d.Partial(true)

	log.Printf("[DEBUG] Retrieving existing Node Pool %q (Kubernetes Cluster %q / Resource Group %q)..", name, clusterName, resourceGroup)
	existing, err := client.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("[DEBUG] Node Pool %q was not found in Managed Kubernetes Cluster %q / Resource Group %q!", name, clusterName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
	}
	if existing.ManagedClusterAgentPoolProfileProperties == nil {
		return fmt.Errorf("Error retrieving Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): `properties` was nil", name, clusterName, resourceGroup)
	}

	props := existing.ManagedClusterAgentPoolProfileProperties

	// store the existing value should the user have opted to ignore it
	enableAutoScaling := false
	if props.EnableAutoScaling != nil {
		enableAutoScaling = *props.EnableAutoScaling
	}

	log.Printf("[DEBUG] Determining delta for existing Node Pool %q (Kubernetes Cluster %q / Resource Group %q)..", name, clusterName, resourceGroup)

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

	if d.HasChange("enable_node_public_ip") {
		props.EnableNodePublicIP = utils.Bool(d.Get("enable_node_public_ip").(bool))
	}

	if d.HasChange("max_count") {
		props.MaxCount = utils.Int32(int32(d.Get("max_count").(int)))
	}

	if d.HasChange("min_count") {
		props.MinCount = utils.Int32(int32(d.Get("min_count").(int)))
	}

	if d.HasChange("node_count") {
		props.Count = utils.Int32(int32(d.Get("node_count").(int)))
	}

	if d.HasChange("node_taints") {
		nodeTaintsRaw := d.Get("node_taints").([]interface{})
		nodeTaints := utils.ExpandStringSlice(nodeTaintsRaw)
		props.NodeTaints = nodeTaints
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
		if minCount == 0 {
			return fmt.Errorf("`min_count` must be configured when `enable_auto_scaling` is set to `true`")
		}

		if minCount > maxCount {
			return fmt.Errorf("`max_count` must be >= `min_count`")
		}
	} else if minCount > 0 || maxCount > 0 {
		return fmt.Errorf("`max_count` and `min_count` must be set to `0` when enable_auto_scaling is set to `false`")
	}

	log.Printf("[DEBUG] Updating existing Node Pool %q (Kubernetes Cluster %q / Resource Group %q)..", name, clusterName, resourceGroup)
	existing.ManagedClusterAgentPoolProfileProperties = props
	future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, name, existing)
	if err != nil {
		return fmt.Errorf("Error updating Node Pool %q (Kubernetes Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Node Pool %q (Kubernetes Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
	}

	d.Partial(false)

	return resourceArmKubernetesClusterNodePoolRead(d, meta)
}

func resourceArmKubernetesClusterNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).Containers.KubernetesClustersClient
	poolsClient := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseKubernetesNodePoolID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	clusterName := id.ClusterName
	name := id.Name

	// if the parent cluster doesn't exist then the node pool won't
	cluster, err := clustersClient.Get(ctx, resourceGroup, clusterName)
	if err != nil {
		if utils.ResponseWasNotFound(cluster.Response) {
			log.Printf("[DEBUG] Managed Kubernetes Cluster %q was not found in Resource Group %q - removing from state!", clusterName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Managed Kubernetes Cluster %q (Resource Group %q): %+v", clusterName, resourceGroup, err)
	}

	resp, err := poolsClient.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Node Pool %q was not found in Managed Kubernetes Cluster %q / Resource Group %q - removing from state!", name, clusterName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): %+v", name, clusterName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("kubernetes_cluster_id", cluster.ID)

	if props := resp.ManagedClusterAgentPoolProfileProperties; props != nil {
		if err := d.Set("availability_zones", utils.FlattenStringSlice(props.AvailabilityZones)); err != nil {
			return fmt.Errorf("Error setting `availability_zones`: %+v", err)
		}

		d.Set("enable_auto_scaling", props.EnableAutoScaling)
		d.Set("enable_node_public_ip", props.EnableNodePublicIP)

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

		count := 0
		if props.Count != nil {
			count = int(*props.Count)
		}
		d.Set("node_count", count)

		if err := d.Set("node_taints", utils.FlattenStringSlice(props.NodeTaints)); err != nil {
			return fmt.Errorf("Error setting `node_taints`: %+v", err)
		}

		osDiskSizeGB := 0
		if props.OsDiskSizeGB != nil {
			osDiskSizeGB = int(*props.OsDiskSizeGB)
		}
		d.Set("os_disk_size_gb", osDiskSizeGB)
		d.Set("os_type", string(props.OsType))
		d.Set("vnet_subnet_id", props.VnetSubnetID)
		d.Set("vm_size", string(props.VMSize))
	}

	return nil
}

func resourceArmKubernetesClusterNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.AgentPoolsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseKubernetesNodePoolID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ClusterName, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): %+v", id.Name, id.ClusterName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Node Pool %q (Managed Kubernetes Cluster %q / Resource Group %q): %+v", id.Name, id.ClusterName, id.ResourceGroup, err)
	}

	return nil
}
