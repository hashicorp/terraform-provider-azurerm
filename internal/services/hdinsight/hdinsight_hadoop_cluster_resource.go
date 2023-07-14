// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2018-06-01/hdinsight" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// NOTE: this isn't a recommended way of building resources in Terraform
// this pattern is used to work around a generic but pedantic API endpoint
var hdInsightHadoopClusterHeadNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         2,
	MaxInstanceCount:         utils.Int(2),
	CanSpecifyDisks:          false,
	FixedMinInstanceCount:    utils.Int32(int32(1)),
	FixedTargetInstanceCount: utils.Int32(int32(2)),
}

var hdInsightHadoopClusterWorkerNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	CanSpecifyDisks:         false,
	CanAutoScaleByCapacity:  true,
	CanAutoScaleOnSchedule:  true,
}

var hdInsightHadoopClusterZookeeperNodeDefinition = HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         3,
	MaxInstanceCount:         utils.Int(3),
	CanSpecifyDisks:          false,
	FixedMinInstanceCount:    utils.Int32(int32(1)),
	FixedTargetInstanceCount: utils.Int32(int32(3)),
}

func resourceHDInsightHadoopCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHDInsightHadoopClusterCreate,
		Read:   resourceHDInsightHadoopClusterRead,
		Update: hdinsightClusterUpdate("Hadoop", resourceHDInsightHadoopClusterRead),
		Delete: hdinsightClusterDelete("Hadoop"),

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
			"name": SchemaHDInsightName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"cluster_version": SchemaHDInsightClusterVersion(),

			"tier": SchemaHDInsightTier(),

			"tls_min_version": SchemaHDInsightTls(),

			"component_version": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"hadoop": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"disk_encryption": SchemaHDInsightsDiskEncryptionProperties(),

			"compute_isolation": SchemaHDInsightsComputeIsolation(),

			"gateway": SchemaHDInsightsGateway(),

			"metastores": SchemaHDInsightsExternalMetastores(),

			"network": SchemaHDInsightsNetwork(),

			"security_profile": SchemaHDInsightsSecurityProfile(),

			"storage_account": SchemaHDInsightsStorageAccounts(),

			"storage_account_gen2": SchemaHDInsightsGen2StorageAccounts(),

			"roles": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"head_node": SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightHadoopClusterHeadNodeDefinition, true),

						"worker_node": SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightHadoopClusterWorkerNodeDefinition, true),

						"zookeeper_node": SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightHadoopClusterZookeeperNodeDefinition, true),

						"edge_node": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"target_instance_count": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 25),
									},

									"vm_size": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(validate.NodeDefinitionVMSize, false),
									},

									"install_script_action": {
										Type:     pluginsdk.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"name": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"uri": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"parameters": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},

									"https_endpoints": SchemaHDInsightsHttpsEndpoints(),

									"uninstall_script_actions": SchemaHDInsightsScriptActions(),
								},
							},
						},
					},
				},
			},

			"tags": tags.Schema(),

			"https_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ssh_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"monitor": SchemaHDInsightsMonitor(),

			"extension": SchemaHDInsightsExtension(),
		},
	}
}

func resourceHDInsightHadoopClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HDInsight.ClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	extensionsClient := meta.(*clients.Client).HDInsight.ExtensionsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := parse.NewClusterID(subscriptionId, resourceGroup, name)
	location := azure.NormalizeLocation(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	t := d.Get("tags").(map[string]interface{})
	tier := hdinsight.Tier(d.Get("tier").(string))
	tls := d.Get("tls_min_version").(string)

	componentVersionsRaw := d.Get("component_version").([]interface{})
	componentVersions := expandHDInsightHadoopComponentVersion(componentVersionsRaw)

	gatewayRaw := d.Get("gateway").([]interface{})
	configurations := ExpandHDInsightsConfigurations(gatewayRaw)

	metastoresRaw := d.Get("metastores").([]interface{})
	metastores := expandHDInsightsMetastore(metastoresRaw)
	for k, v := range metastores {
		configurations[k] = v
	}

	networkPropertiesRaw := d.Get("network").([]interface{})
	networkProperties := ExpandHDInsightsNetwork(networkPropertiesRaw)

	computeIsolationProperties := ExpandHDInsightComputeIsolationProperties(d.Get("compute_isolation").([]interface{}))

	storageAccountsRaw := d.Get("storage_account").([]interface{})
	storageAccountsGen2Raw := d.Get("storage_account_gen2").([]interface{})
	storageAccounts, identity, err := ExpandHDInsightsStorageAccounts(storageAccountsRaw, storageAccountsGen2Raw)
	if err != nil {
		return fmt.Errorf("expanding `storage_account`: %s", err)
	}

	rolesRaw := d.Get("roles").([]interface{})
	hadoopRoles := hdInsightRoleDefinition{
		HeadNodeDef:      hdInsightHadoopClusterHeadNodeDefinition,
		WorkerNodeDef:    hdInsightHadoopClusterWorkerNodeDefinition,
		ZookeeperNodeDef: hdInsightHadoopClusterZookeeperNodeDefinition,
	}
	roles, err := expandHDInsightRoles(rolesRaw, hadoopRoles)
	if err != nil {
		return fmt.Errorf("expanding `roles`: %+v", err)
	}

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_hdinsight_hadoop_cluster", id.ID())
	}

	params := hdinsight.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &hdinsight.ClusterCreateProperties{
			Tier:                   tier,
			OsType:                 hdinsight.OSTypeLinux,
			ClusterVersion:         utils.String(clusterVersion),
			MinSupportedTLSVersion: utils.String(tls),
			NetworkProperties:      networkProperties,
			ClusterDefinition: &hdinsight.ClusterDefinition{
				Kind:             utils.String("Hadoop"),
				ComponentVersion: componentVersions,
				Configurations:   configurations,
			},
			StorageProfile: &hdinsight.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &hdinsight.ComputeProfile{
				Roles: roles,
			},
			ComputeIsolationProperties: computeIsolationProperties,
		},
		Tags:     tags.Expand(t),
		Identity: identity,
	}

	if diskEncryptionPropertiesRaw, ok := d.GetOk("disk_encryption"); ok {
		diskEncryptionProperties, err := ExpandHDInsightsDiskEncryptionProperties(diskEncryptionPropertiesRaw.([]interface{}))
		if err != nil {
			return err
		}
		params.Properties.DiskEncryptionProperties = diskEncryptionProperties
	}

	if v, ok := d.GetOk("security_profile"); ok {
		params.Properties.SecurityProfile = ExpandHDInsightSecurityProfile(v.([]interface{}))

		params.Identity = &hdinsight.ClusterIdentity{
			Type:                   hdinsight.ResourceIdentityTypeUserAssigned,
			UserAssignedIdentities: make(map[string]*hdinsight.ClusterIdentityUserAssignedIdentitiesValue),
		}

		if params.Properties.SecurityProfile != nil && params.Properties.SecurityProfile.MsiResourceID != nil {
			params.Identity.UserAssignedIdentities[*params.Properties.SecurityProfile.MsiResourceID] = &hdinsight.ClusterIdentityUserAssignedIdentitiesValue{}
		}
	}

	future, err := client.Create(ctx, resourceGroup, name, params)
	if err != nil {
		return fmt.Errorf("creating HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("reading ID for HDInsight Hadoop Cluster %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(id.ID())

	// We can only add an edge node after creation
	if v, ok := d.GetOk("roles.0.edge_node"); ok {
		edgeNodeRaw := v.([]interface{})
		applicationsClient := meta.(*clients.Client).HDInsight.ApplicationsClient
		edgeNodeConfig := edgeNodeRaw[0].(map[string]interface{})

		err := createHDInsightEdgeNodes(ctx, applicationsClient, resourceGroup, name, edgeNodeConfig)
		if err != nil {
			return err
		}

		// we can't rely on the use of the Future here due to the node being successfully completed but now the cluster is applying those changes.
		log.Printf("[DEBUG] Waiting for Hadoop Cluster to %q (Resource Group %q) to finish applying edge node", name, resourceGroup)
		stateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{"AzureVMConfiguration", "Accepted", "HdInsightConfiguration"},
			Target:     []string{"Running"},
			Refresh:    hdInsightWaitForReadyRefreshFunc(ctx, client, resourceGroup, name),
			MinTimeout: 15 * time.Second,
			Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for HDInsight Cluster %q (Resource Group %q) to be running: %s", name, resourceGroup, err)
		}
	}

	// We can only enable monitoring after creation
	if v, ok := d.GetOk("monitor"); ok {
		monitorRaw := v.([]interface{})
		if err := enableHDInsightMonitoring(ctx, extensionsClient, resourceGroup, name, monitorRaw); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("extension"); ok {
		extensionRaw := v.([]interface{})
		if err := enableHDInsightAzureMonitor(ctx, extensionsClient, resourceGroup, name, extensionRaw); err != nil {
			return err
		}
	}

	return resourceHDInsightHadoopClusterRead(d, meta)
}

func resourceHDInsightHadoopClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).HDInsight.ClustersClient
	configurationsClient := meta.(*clients.Client).HDInsight.ConfigurationsClient
	extensionsClient := meta.(*clients.Client).HDInsight.ExtensionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Name

	resp, err := clustersClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] HDInsight Hadoop Cluster %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configurations, err := configurationsClient.List(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Configuration for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	gateway, exists := configurations.Configurations["gateway"]
	if !exists {
		return fmt.Errorf("retrieving gateway for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	// storage_account isn't returned so I guess we just leave it ¯\_(ツ)_/¯
	if props := resp.Properties; props != nil {
		tier := ""
		// the Azure API is inconsistent here, so rewrite this into the casing we expect
		for _, v := range hdinsight.PossibleTierValues() {
			if strings.EqualFold(string(v), string(props.Tier)) {
				tier = string(v)
			}
		}
		d.Set("tier", tier)
		d.Set("cluster_version", props.ClusterVersion)
		d.Set("tls_min_version", props.MinSupportedTLSVersion)

		if def := props.ClusterDefinition; def != nil {
			if err := d.Set("component_version", flattenHDInsightHadoopComponentVersion(def.ComponentVersion)); err != nil {
				return fmt.Errorf("flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", FlattenHDInsightsConfigurations(gateway, d)); err != nil {
				return fmt.Errorf("flattening `gateway`: %+v", err)
			}

			flattenHDInsightsMetastores(d, configurations.Configurations)

			if props.NetworkProperties != nil {
				if err := d.Set("network", FlattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
					return fmt.Errorf("flattening `network`: %+v", err)
				}
			}
		}

		hadoopRoles := hdInsightRoleDefinition{
			HeadNodeDef:      hdInsightHadoopClusterHeadNodeDefinition,
			WorkerNodeDef:    hdInsightHadoopClusterWorkerNodeDefinition,
			ZookeeperNodeDef: hdInsightHadoopClusterZookeeperNodeDefinition,
		}
		flattenedRoles := flattenHDInsightRoles(d, props.ComputeProfile, hadoopRoles)

		applicationsClient := meta.(*clients.Client).HDInsight.ApplicationsClient

		edgeNode, err := applicationsClient.Get(ctx, resourceGroup, name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(edgeNode.Response) {
				return fmt.Errorf("reading edge node for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if edgeNodeProps := edgeNode.Properties; edgeNodeProps != nil {
			flattenedRoles = flattenHDInsightEdgeNode(flattenedRoles, edgeNodeProps)
		}

		if props.DiskEncryptionProperties != nil {
			diskEncryptionProps, err := FlattenHDInsightsDiskEncryptionProperties(*props.DiskEncryptionProperties)
			if err != nil {
				return err
			}
			if err := d.Set("disk_encryption", diskEncryptionProps); err != nil {
				return fmt.Errorf("flattening `disk_encryption`: %+v", err)
			}
		}

		if props.ComputeIsolationProperties != nil {
			if err := d.Set("compute_isolation", FlattenHDInsightComputeIsolationProperties(*props.ComputeIsolationProperties)); err != nil {
				return fmt.Errorf("failed setting `compute_isolation`: %+v", err)
			}
		}

		if err := d.Set("roles", flattenedRoles); err != nil {
			return fmt.Errorf("flattening `roles`: %+v", err)
		}

		httpEndpoint := FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
		d.Set("https_endpoint", httpEndpoint)
		sshEndpoint := FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
		d.Set("ssh_endpoint", sshEndpoint)

		monitor, err := extensionsClient.GetMonitoringStatus(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("reading monitor configuration for HDInsight Hadoop Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		d.Set("monitor", flattenHDInsightMonitoring(monitor))

		extension, err := extensionsClient.GetAzureMonitorStatus(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("reading extension configuration for HDInsight Hadoop Cluster %q (Resource Group %q) %+v", name, resourceGroup, err)
		}

		d.Set("extension", flattenHDInsightAzureMonitor(extension))

		if err := d.Set("security_profile", flattenHDInsightSecurityProfile(props.SecurityProfile, d)); err != nil {
			return fmt.Errorf("setting `security_profile`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenHDInsightEdgeNode(roles []interface{}, props *hdinsight.ApplicationProperties) []interface{} {
	if len(roles) == 0 || props == nil {
		return roles
	}

	role := roles[0].(map[string]interface{})

	edgeNode := make(map[string]interface{})
	if computeProfile := props.ComputeProfile; computeProfile != nil {
		if roles := computeProfile.Roles; roles != nil {
			for _, role := range *roles {
				if targetInstanceCount := role.TargetInstanceCount; targetInstanceCount != nil {
					edgeNode["target_instance_count"] = targetInstanceCount
				}
				if hardwareProfile := role.HardwareProfile; hardwareProfile != nil && hardwareProfile.VMSize != nil {
					vmSize := ""
					// the Azure API is inconsistent here, so rewrite this into the casing we expect
					for _, v := range validate.NodeDefinitionVMSize {
						if strings.EqualFold(v, *hardwareProfile.VMSize) {
							vmSize = v
						}
					}
					edgeNode["vm_size"] = vmSize
				}
			}
		}
	}

	actions := make(map[string]interface{})
	if installScriptActions := props.InstallScriptActions; installScriptActions != nil {
		for _, action := range *installScriptActions {
			actions["name"] = action.Name
			actions["uri"] = action.URI
			actions["parameters"] = action.Parameters
		}
	}

	if uninstallScriptActions := props.UninstallScriptActions; uninstallScriptActions != nil && len(*uninstallScriptActions) != 0 {
		uninstallActions := make(map[string]interface{})
		for _, uninstallAction := range *uninstallScriptActions {
			actions["name"] = uninstallAction.Name
			actions["uri"] = uninstallAction.URI
			actions["parameters"] = uninstallAction.Parameters
		}
		edgeNode["uninstall_script_actions"] = []interface{}{uninstallActions}
	}

	if HTTPSEndpoints := props.HTTPSEndpoints; HTTPSEndpoints != nil && len(*HTTPSEndpoints) != 0 {
		httpsEndpoints := make(map[string]interface{})
		for _, HTTPSEndpoint := range *HTTPSEndpoints {
			httpsEndpoints["access_modes"] = HTTPSEndpoint.AccessModes
			httpsEndpoints["destination_port"] = HTTPSEndpoint.DestinationPort
			httpsEndpoints["disable_gateway_auth"] = HTTPSEndpoint.DisableGatewayAuth
			httpsEndpoints["private_ip_address"] = HTTPSEndpoint.PrivateIPAddress
			httpsEndpoints["sub_domain_suffix"] = HTTPSEndpoint.SubDomainSuffix
		}
		edgeNode["https_endpoints"] = []interface{}{httpsEndpoints}
	}

	edgeNode["install_script_action"] = []interface{}{actions}

	role["edge_node"] = []interface{}{edgeNode}

	return []interface{}{role}
}

func expandHDInsightHadoopComponentVersion(input []interface{}) map[string]*string {
	vs := input[0].(map[string]interface{})
	return map[string]*string{
		"Hadoop": utils.String(vs["hadoop"].(string)),
	}
}

func flattenHDInsightHadoopComponentVersion(input map[string]*string) []interface{} {
	hadoopVersion := ""
	if v, ok := input["Hadoop"]; ok {
		if v != nil {
			hadoopVersion = *v
		}
	}
	return []interface{}{
		map[string]interface{}{
			"hadoop": hadoopVersion,
		},
	}
}

func expandHDInsightApplicationEdgeNodeInstallScriptActions(input []interface{}) *[]hdinsight.RuntimeScriptAction {
	actions := make([]hdinsight.RuntimeScriptAction, 0)

	for _, v := range input {
		val := v.(map[string]interface{})

		name := val["name"].(string)
		uri := val["uri"].(string)
		parameters := val["parameters"].(string)

		action := hdinsight.RuntimeScriptAction{
			Name: utils.String(name),
			URI:  utils.String(uri),
			// The only role available for edge nodes is edgenode
			Parameters: utils.String(parameters),
			Roles:      &[]string{"edgenode"},
		}

		actions = append(actions, action)
	}

	return &actions
}

func expandHDInsightApplicationEdgeNodeHttpsEndpoints(input []interface{}) *[]hdinsight.ApplicationGetHTTPSEndpoint {
	endpoints := make([]hdinsight.ApplicationGetHTTPSEndpoint, 0)
	if len(input) == 0 || input[0] == nil {
		return &endpoints
	}

	for _, v := range input {
		val := v.(map[string]interface{})

		accessModes := val["access_modes"].([]string)
		destinationPort := val["destination_port"].(int32)
		disableGatewayAuth := val["disable_gateway_auth"].(bool)
		privateIpAddress := val["private_ip_address"].(string)
		subDomainSuffix := val["sub_domain_suffix"].(string)

		endPoint := hdinsight.ApplicationGetHTTPSEndpoint{
			AccessModes:        &accessModes,
			DestinationPort:    utils.Int32(destinationPort),
			PrivateIPAddress:   utils.String(privateIpAddress),
			SubDomainSuffix:    utils.String(subDomainSuffix),
			DisableGatewayAuth: utils.Bool(disableGatewayAuth),
		}

		endpoints = append(endpoints, endPoint)
	}

	return &endpoints
}

func expandHDInsightApplicationEdgeNodeUninstallScriptActions(input []interface{}) *[]hdinsight.RuntimeScriptAction {
	actions := make([]hdinsight.RuntimeScriptAction, 0)
	if len(input) == 0 || input[0] == nil {
		return &actions
	}

	for _, v := range input {
		val := v.(map[string]interface{})

		name := val["name"].(string)
		uri := val["uri"].(string)
		parameters := val["parameters"].(string)

		action := hdinsight.RuntimeScriptAction{
			Name:       utils.String(name),
			URI:        utils.String(uri),
			Parameters: utils.String(parameters),
			Roles:      &[]string{"edgenode"},
		}

		actions = append(actions, action)
	}

	return &actions
}

func hdInsightWaitForReadyRefreshFunc(ctx context.Context, client *hdinsight.ClustersClient, resourceGroupName string, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			return nil, "Error", fmt.Errorf("issuing read request in hdInsightWaitForReadyRefreshFunc to Hadoop Cluster %q (Resource Group %q): %s", name, resourceGroupName, err)
		}
		if props := res.Properties; props != nil {
			if state := props.ClusterState; state != nil {
				return res, *state, nil
			}
		}

		return res, "Pending", nil
	}
}
