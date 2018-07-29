package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2018-02-01/servicefabric"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceFabricCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceFabricClusterCreate,
		Read:   resourceArmServiceFabricClusterRead,
		Update: resourceArmServiceFabricClusterUpdate,
		Delete: resourceArmServiceFabricClusterDelete,

		Schema: map[string]*schema.Schema{
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reliability_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"upgrade_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"management_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"client_endpoint_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"http_endpoint_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protected_account_key_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"blob_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"queue_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"table_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vm_image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_thumbprint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_store_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_protection_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_certificate_thumbprint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": locationSchema(),
			"tags":     tagsSchema(),
			"cluster_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmServiceFabricClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceFabricClustersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Service Fabric creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("cluster_name").(string)
	reliabilityLevel := d.Get("reliability_level").(string)
	upgradeMode := d.Get("upgrade_mode").(string)
	managementEndpoint := d.Get("management_endpoint").(string)
	nodeName := d.Get("node_name").(string)
	instanceCount := utils.Int32(int32(d.Get("instance_count").(int)))
	isPrimary := true
	clientEndpointPort := utils.Int32(int32(d.Get("client_endpoint_port").(int)))
	httpEndpointPort := utils.Int32(int32(d.Get("http_endpoint_port").(int)))
	vmImage := d.Get("vm_image").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	storageAccountName := d.Get("storage_account_name").(string)
	protectedAccountKeyName := d.Get("protected_account_key_name").(string)
	blobEndpoint := d.Get("blob_endpoint").(string)
	queueEndpoint := d.Get("queue_endpoint").(string)
	tableEndpoint := d.Get("table_endpoint").(string)

	nodeTypeDescription := servicefabric.NodeTypeDescription{
		Name:                         &nodeName,
		VMInstanceCount:              instanceCount,
		IsPrimary:                    &isPrimary,
		ClientConnectionEndpointPort: clientEndpointPort,
		HTTPGatewayEndpointPort:      httpEndpointPort,
	}

	nodeTypes := make([]servicefabric.NodeTypeDescription, 0)
	nodeTypes = append(nodeTypes, nodeTypeDescription)

	diagnosticsStorageAccountConfig := servicefabric.DiagnosticsStorageAccountConfig{
		StorageAccountName:      &storageAccountName,
		ProtectedAccountKeyName: &protectedAccountKeyName,
		BlobEndpoint:            &blobEndpoint,
		QueueEndpoint:           &queueEndpoint,
		TableEndpoint:           &tableEndpoint,
	}

	clusterProperties := servicefabric.ClusterProperties{
		ReliabilityLevel:   servicefabric.ReliabilityLevel(reliabilityLevel),
		UpgradeMode:        servicefabric.UpgradeMode(upgradeMode),
		ManagementEndpoint: &managementEndpoint,
		NodeTypes:          &nodeTypes,
		VMImage:            &vmImage,
		DiagnosticsStorageAccountConfig: &diagnosticsStorageAccountConfig,
	}

	if v, ok := d.GetOk("certificate_thumbprint"); ok {
		certificate := servicefabric.CertificateDescription{}
		certificateThumbprint := v.(string)
		certificate.Thumbprint = &certificateThumbprint

		if v, ok := d.GetOk("certificate_store_value"); ok {
			certificateStoreValue := v.(string)
			certificate.X509StoreName = servicefabric.X509StoreName(certificateStoreValue)
		}

		clusterProperties.Certificate = &certificate
	}

	if v, ok := d.GetOk("admin_certificate_thumbprint"); ok {
		clientCertificate := servicefabric.ClientCertificateThumbprint{}
		adminCertificateThumbprint := v.(string)
		isAdmin := true
		clientCertificate.CertificateThumbprint = &adminCertificateThumbprint
		clientCertificate.IsAdmin = &isAdmin

		clientCertificateArray := make([]servicefabric.ClientCertificateThumbprint, 0)
		clientCertificateArray = append(clientCertificateArray, clientCertificate)

		clusterProperties.ClientCertificateThumbprints = &clientCertificateArray
	}

	if v, ok := d.GetOk("cluster_protection_level"); ok {
		fabricSettingsParameters := servicefabric.SettingsParameterDescription{}

		clusterProtectionLevelName := "ClusterProtectionLevel"

		clusterProtectionLevel := v.(string)
		fabricSettingsParameters.Name = &clusterProtectionLevelName
		fabricSettingsParameters.Value = &clusterProtectionLevel

		fabricSettingsParametersArray := make([]servicefabric.SettingsParameterDescription, 0)
		fabricSettingsParametersArray = append(fabricSettingsParametersArray, fabricSettingsParameters)

		fabricSettingsName := "Security"

		fabricSettings := servicefabric.SettingsSectionDescription{
			Name:       &fabricSettingsName,
			Parameters: &fabricSettingsParametersArray,
		}

		fabricSettingsArray := make([]servicefabric.SettingsSectionDescription, 0)
		fabricSettingsArray = append(fabricSettingsArray, fabricSettings)
		clusterProperties.FabricSettings = &fabricSettingsArray
	}

	cluster := servicefabric.Cluster{
		Location:          utils.String(location),
		ClusterProperties: &clusterProperties,
		Tags:              expandTags(tags),
	}

	future, err := client.Create(ctx, resourceGroup, name, cluster)
	if err != nil {
		return fmt.Errorf("Error creating Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Service Fabric %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceFabricClusterRead(d, meta)
}

func resourceArmServiceFabricClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceFabricClustersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Service Fabric Cluster update.")

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("cluster_name").(string)
	reliabilityLevel := d.Get("reliability_level").(string)
	upgradeMode := d.Get("upgrade_mode").(string)
	nodeName := d.Get("node_name").(string)
	instanceCount := utils.Int32(int32(d.Get("instance_count").(int)))
	clientEndpointPort := utils.Int32(int32(d.Get("client_endpoint_port").(int)))
	httpEndpointPort := utils.Int32(int32(d.Get("http_endpoint_port").(int)))
	tags := d.Get("tags").(map[string]interface{})
	isPrimary := true

	nodeTypeDescription := servicefabric.NodeTypeDescription{
		Name:                         &nodeName,
		VMInstanceCount:              instanceCount,
		IsPrimary:                    &isPrimary,
		ClientConnectionEndpointPort: clientEndpointPort,
		HTTPGatewayEndpointPort:      httpEndpointPort,
	}

	nodeTypes := make([]servicefabric.NodeTypeDescription, 0)

	nodeTypes = append(nodeTypes, nodeTypeDescription)

	clusterProperties := servicefabric.ClusterPropertiesUpdateParameters{
		ReliabilityLevel: servicefabric.ReliabilityLevel1(reliabilityLevel),
		UpgradeMode:      servicefabric.UpgradeMode1(upgradeMode),
		NodeTypes:        &nodeTypes,
	}

	clusterUpdateParameters := servicefabric.ClusterUpdateParameters{
		ClusterPropertiesUpdateParameters: &clusterProperties,
		Tags: expandTags(tags),
	}

	future, err := client.Update(ctx, resourceGroup, name, clusterUpdateParameters)
	if err != nil {
		return fmt.Errorf("Error updating Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for update of Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot ID of Service Fabric Cluster %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmServiceFabricClusterRead(d, meta)
}

func resourceArmServiceFabricClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceFabricClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading service fabric cluster %s", id)

	resourceGroup := id.ResourceGroup
	name := id.Path["clusters"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Service Fabric Cluster %q (Resource Group %q) was not found - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// TODO: set all the other properties here
	clusterEndpoint := resp.ClusterProperties.ClusterEndpoint

	d.Set("cluster_name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("cluster_endpoint", clusterEndpoint)

	return nil
}

func resourceArmServiceFabricClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).serviceFabricClustersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["clusters"]

	log.Printf("[DEBUG] Deleting Service Fabric Cluster %q (Resource Group %q)", name, resourceGroup)

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Service Fabric Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
