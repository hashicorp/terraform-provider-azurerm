package kusto

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2022-02-01/kusto"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoCluster() *pluginsdk.Resource {
	s := &pluginsdk.Resource{
		Create: resourceKustoClusterCreateUpdate,
		Read:   resourceKustoClusterRead,
		Update: resourceKustoClusterCreateUpdate,
		Delete: resourceKustoClusterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ClusterID(id)
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
				ValidateFunc: validate.ClusterName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(kusto.AzureSkuNameDevNoSLAStandardD11V2),
								string(kusto.AzureSkuNameDevNoSLAStandardE2aV4),
								string(kusto.AzureSkuNameStandardD11V2),
								string(kusto.AzureSkuNameStandardD12V2),
								string(kusto.AzureSkuNameStandardD13V2),
								string(kusto.AzureSkuNameStandardD14V2),
								string(kusto.AzureSkuNameStandardDS13V21TBPS),
								string(kusto.AzureSkuNameStandardDS13V22TBPS),
								string(kusto.AzureSkuNameStandardDS14V23TBPS),
								string(kusto.AzureSkuNameStandardDS14V24TBPS),
								string(kusto.AzureSkuNameStandardE16asV43TBPS),
								string(kusto.AzureSkuNameStandardE16asV44TBPS),
								string(kusto.AzureSkuNameStandardE16aV4),
								string(kusto.AzureSkuNameStandardE2aV4),
								string(kusto.AzureSkuNameStandardE4aV4),
								string(kusto.AzureSkuNameStandardE64iV3),
								string(kusto.AzureSkuNameStandardE8asV41TBPS),
								string(kusto.AzureSkuNameStandardE8asV42TBPS),
								string(kusto.AzureSkuNameStandardE8aV4),
								string(kusto.AzureSkuNameStandardL16s),
								string(kusto.AzureSkuNameStandardL4s),
								string(kusto.AzureSkuNameStandardL8s),
								string(kusto.AzureSkuNameStandardL16sV2),
								string(kusto.AzureSkuNameStandardL8sV2),
							}, false),
						},

						"capacity": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 1000),
						},
					},
				},
			},

			"trusted_external_tenants": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty, validation.StringInSlice([]string{"*"}, false)),
				},
			},

			"optimized_auto_scale": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"minimum_instances": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"maximum_instances": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
					},
				},
			},

			"virtual_network_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"engine_public_ip_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"data_management_public_ip_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"language_extensions": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(kusto.LanguageExtensionNamePYTHON),
						string(kusto.LanguageExtensionNameR),
					}, false),
				},
			},

			"engine": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.EngineTypeV2),
					string(kusto.EngineTypeV3),
				}, false),
			},

			"uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"data_ingestion_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_ip_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(kusto.PublicIPTypeIPv4),
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.PublicIPTypeIPv4),
					string(kusto.PublicIPTypeDualStack),
				}, false),
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"double_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"auto_stop_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"disk_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"streaming_ingestion_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"purge_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"tags": tags.Schema(),
		},
	}

	if features.FourPointOhBeta() {
		s.Schema["engine"].Default = string(kusto.EngineTypeV3)
	} else {
		s.Schema["engine"].Default = string(kusto.EngineTypeV2)
	}

	return s
}

func resourceKustoClusterCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Cluster creation.")

	id := parse.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		server, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(server.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(server.Response) {
			return tf.ImportAsExistsError("azurerm_kusto_cluster", id.ID())
		}
	}

	locks.ByID(id.Name)
	defer locks.UnlockByID(id.Name)

	sku, err := expandKustoClusterSku(d.Get("sku").([]interface{}))
	if err != nil {
		return err
	}

	optimizedAutoScale := expandOptimizedAutoScale(d.Get("optimized_auto_scale").([]interface{}))

	if optimizedAutoScale != nil && *optimizedAutoScale.IsEnabled {
		// Ensure that requested Capcity is always between min and max to support updating to not overlapping autoscale ranges
		if *sku.Capacity < *optimizedAutoScale.Minimum {
			sku.Capacity = utils.Int32(*optimizedAutoScale.Minimum)
		}
		if *sku.Capacity > *optimizedAutoScale.Maximum {
			sku.Capacity = utils.Int32(*optimizedAutoScale.Maximum)
		}

		// Capacity must be set for the initial creation when using OptimizedAutoScaling but cannot be updated
		if d.HasChange("sku.0.capacity") && !d.IsNewResource() {
			return fmt.Errorf("cannot change `sku.capacity` when `optimized_auto_scaling.enabled` is set to `true`")
		}

		if *optimizedAutoScale.Minimum > *optimizedAutoScale.Maximum {
			return fmt.Errorf("`optimized_auto_scaling.maximum_instances` must be >= `optimized_auto_scaling.minimum_instances`")
		}
	}

	engine := kusto.EngineType(d.Get("engine").(string))

	publicNetworkAccess := kusto.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = kusto.PublicNetworkAccessDisabled
	}

	clusterProperties := kusto.ClusterProperties{
		OptimizedAutoscale:     optimizedAutoScale,
		EnableAutoStop:         utils.Bool(d.Get("auto_stop_enabled").(bool)),
		EnableDiskEncryption:   utils.Bool(d.Get("disk_encryption_enabled").(bool)),
		EnableDoubleEncryption: utils.Bool(d.Get("double_encryption_enabled").(bool)),
		EnableStreamingIngest:  utils.Bool(d.Get("streaming_ingestion_enabled").(bool)),
		EnablePurge:            utils.Bool(d.Get("purge_enabled").(bool)),
		EngineType:             engine,
		PublicNetworkAccess:    publicNetworkAccess,
		PublicIPType:           kusto.PublicIPType(d.Get("public_ip_type").(string)),
		TrustedExternalTenants: expandTrustedExternalTenants(d.Get("trusted_external_tenants").([]interface{})),
	}

	if v, ok := d.GetOk("virtual_network_configuration"); ok {
		vnet := expandKustoClusterVNET(v.([]interface{}))
		clusterProperties.VirtualNetworkConfiguration = vnet
	}

	expandedIdentity, err := expandClusterIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	kustoCluster := kusto.Cluster{
		Name:              utils.String(id.Name),
		Location:          utils.String(location.Normalize(d.Get("location").(string))),
		Identity:          expandedIdentity,
		Sku:               sku,
		ClusterProperties: &clusterProperties,
		Tags:              tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	zones := zones.Expand(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		kustoCluster.Zones = &zones
	}

	ifMatch := ""
	ifNoneMatch := ""
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, kustoCluster, ifMatch, ifNoneMatch)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if v, ok := d.GetOk("language_extensions"); ok {
		languageExtensions := expandKustoClusterLanguageExtensions(v.([]interface{}))

		currentLanguageExtensions, err := client.ListLanguageExtensions(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving the language extensions on %s: %+v", id, err)
		}

		languageExtensionsToAdd := diffLanguageExtensions(*languageExtensions.Value, *currentLanguageExtensions.Value)
		if len(languageExtensionsToAdd) > 0 {
			languageExtensionsListToAdd := kusto.LanguageExtensionsList{
				Value: &languageExtensionsToAdd,
			}

			future, err := client.AddLanguageExtensions(ctx, id.ResourceGroup, id.Name, languageExtensionsListToAdd)
			if err != nil {
				return fmt.Errorf("adding language extensions to %s: %+v", id, err)
			}
			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the addition of language extensions on %s: %+v", id, err)
			}
		}

		languageExtensionsToRemove := diffLanguageExtensions(*currentLanguageExtensions.Value, *languageExtensions.Value)
		if len(languageExtensionsToRemove) > 0 {
			languageExtensionsListToRemove := kusto.LanguageExtensionsList{
				Value: &languageExtensionsToRemove,
			}

			removeLanguageExtensionsFuture, err := client.RemoveLanguageExtensions(ctx, id.ResourceGroup, id.Name, languageExtensionsListToRemove)
			if err != nil {
				return fmt.Errorf("removing language extensions from %s: %+v", id, err)
			}
			if err = removeLanguageExtensionsFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the removal of language extensions from %s: %+v", id, err)
			}
		}
	}

	return resourceKustoClusterRead(d, meta)
}

func resourceKustoClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("zones", zones.Flatten(resp.Zones))

	d.Set("public_network_access_enabled", resp.PublicNetworkAccess == kusto.PublicNetworkAccessEnabled)

	identity, err := flattenClusterIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	if err := d.Set("sku", flattenKustoClusterSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `sku`: %+v", err)
	}

	if err := d.Set("optimized_auto_scale", flattenOptimizedAutoScale(resp.OptimizedAutoscale)); err != nil {
		return fmt.Errorf("setting `optimized_auto_scale`: %+v", err)
	}

	if props := resp.ClusterProperties; props != nil {
		d.Set("double_encryption_enabled", props.EnableDoubleEncryption)
		d.Set("trusted_external_tenants", flattenTrustedExternalTenants(props.TrustedExternalTenants))
		d.Set("auto_stop_enabled", props.EnableAutoStop)
		d.Set("disk_encryption_enabled", props.EnableDiskEncryption)
		d.Set("streaming_ingestion_enabled", props.EnableStreamingIngest)
		d.Set("purge_enabled", props.EnablePurge)
		d.Set("virtual_network_configuration", flattenKustoClusterVNET(props.VirtualNetworkConfiguration))
		d.Set("language_extensions", flattenKustoClusterLanguageExtensions(props.LanguageExtensions))
		d.Set("uri", props.URI)
		d.Set("data_ingestion_uri", props.DataIngestionURI)
		d.Set("engine", props.EngineType)
		d.Set("public_ip_type", props.PublicIPType)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceKustoClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandOptimizedAutoScale(input []interface{}) *kusto.OptimizedAutoscale {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]interface{})
	optimizedAutoScale := &kusto.OptimizedAutoscale{
		Version:   utils.Int32(1),
		IsEnabled: utils.Bool(true),
		Minimum:   utils.Int32(int32(config["minimum_instances"].(int))),
		Maximum:   utils.Int32(int32(config["maximum_instances"].(int))),
	}

	return optimizedAutoScale
}

func flattenOptimizedAutoScale(optimizedAutoScale *kusto.OptimizedAutoscale) []interface{} {
	if optimizedAutoScale == nil {
		return []interface{}{}
	}

	maxInstances := 0
	if optimizedAutoScale.Maximum != nil {
		maxInstances = int(*optimizedAutoScale.Maximum)
	}

	minInstances := 0
	if optimizedAutoScale.Minimum != nil {
		minInstances = int(*optimizedAutoScale.Minimum)
	}

	return []interface{}{
		map[string]interface{}{
			"maximum_instances": maxInstances,
			"minimum_instances": minInstances,
		},
	}
}

func expandKustoClusterSku(input []interface{}) (*kusto.AzureSku, error) {
	sku := input[0].(map[string]interface{})
	name := sku["name"].(string)

	skuNamePrefixToTier := map[string]string{
		"Dev(No SLA)": "Basic",
		"Standard":    "Standard",
	}

	skuNamePrefix := strings.Split(sku["name"].(string), "_")[0]
	tier, ok := skuNamePrefixToTier[skuNamePrefix]
	if !ok {
		return nil, fmt.Errorf("sku name begins with invalid tier, possible are Dev(No SLA) and Standard but is: %q", skuNamePrefix)
	}
	capacity := sku["capacity"].(int)

	azureSku := &kusto.AzureSku{
		Name:     kusto.AzureSkuName(name),
		Tier:     kusto.AzureSkuTier(tier),
		Capacity: utils.Int32(int32(capacity)),
	}

	return azureSku, nil
}

func expandKustoClusterVNET(input []interface{}) *kusto.VirtualNetworkConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	vnet := input[0].(map[string]interface{})
	subnetID := vnet["subnet_id"].(string)
	enginePublicIPID := vnet["engine_public_ip_id"].(string)
	dataManagementPublicIPID := vnet["data_management_public_ip_id"].(string)

	return &kusto.VirtualNetworkConfiguration{
		SubnetID:                 &subnetID,
		EnginePublicIPID:         &enginePublicIPID,
		DataManagementPublicIPID: &dataManagementPublicIPID,
	}
}

func expandKustoClusterLanguageExtensions(input []interface{}) *kusto.LanguageExtensionsList {
	if len(input) == 0 {
		return nil
	}

	extensions := make([]kusto.LanguageExtension, 0)
	for _, language := range input {
		v := kusto.LanguageExtension{
			LanguageExtensionName: kusto.LanguageExtensionName(language.(string)),
		}
		extensions = append(extensions, v)
	}

	return &kusto.LanguageExtensionsList{
		Value: &extensions,
	}
}

func expandClusterIdentity(input []interface{}) (*kusto.Identity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := kusto.Identity{
		Type: kusto.IdentityType(string(expanded.Type)),
	}

	if expanded.Type == identity.TypeUserAssigned || expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.UserAssignedIdentities = make(map[string]*kusto.IdentityUserAssignedIdentitiesValue)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &kusto.IdentityUserAssignedIdentitiesValue{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func flattenClusterIdentity(input *kusto.Identity) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
		if input.UserAssignedIdentities != nil {
			for k, v := range input.UserAssignedIdentities {
				transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
					ClientId:    v.ClientID,
					PrincipalId: v.PrincipalID,
				}
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func flattenKustoClusterSku(sku *kusto.AzureSku) []interface{} {
	if sku == nil {
		return []interface{}{}
	}

	s := map[string]interface{}{
		"name": string(sku.Name),
	}

	if sku.Capacity != nil {
		s["capacity"] = int(*sku.Capacity)
	}

	return []interface{}{s}
}

func flattenKustoClusterVNET(vnet *kusto.VirtualNetworkConfiguration) []interface{} {
	if vnet == nil {
		return []interface{}{}
	}

	subnetID := ""
	if vnet.SubnetID != nil {
		subnetID = *vnet.SubnetID
	}

	enginePublicIPID := ""
	if vnet.EnginePublicIPID != nil {
		enginePublicIPID = *vnet.EnginePublicIPID
	}

	dataManagementPublicIPID := ""
	if vnet.DataManagementPublicIPID != nil {
		dataManagementPublicIPID = *vnet.DataManagementPublicIPID
	}

	output := map[string]interface{}{
		"subnet_id":                    subnetID,
		"engine_public_ip_id":          enginePublicIPID,
		"data_management_public_ip_id": dataManagementPublicIPID,
	}

	return []interface{}{output}
}

func flattenKustoClusterLanguageExtensions(extensions *kusto.LanguageExtensionsList) []interface{} {
	if extensions == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, v := range *extensions.Value {
		output = append(output, v.LanguageExtensionName)
	}

	return output
}

func diffLanguageExtensions(a, b []kusto.LanguageExtension) []kusto.LanguageExtension {
	target := make(map[string]bool)
	for _, x := range b {
		target[string(x.LanguageExtensionName)] = true
	}

	diff := make([]kusto.LanguageExtension, 0)
	for _, x := range a {
		if _, ok := target[string(x.LanguageExtensionName)]; !ok {
			diff = append(diff, x)
		}
	}

	return diff
}
