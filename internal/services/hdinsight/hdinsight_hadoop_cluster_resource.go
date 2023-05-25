package hdinsight

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/applications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/validate"
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
	FixedMinInstanceCount:    pointer.To(int64(1)),
	FixedTargetInstanceCount: pointer.To(int64(2)),
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
	FixedMinInstanceCount:    pointer.To(int64(1)),
	FixedTargetInstanceCount: pointer.To(int64(3)),
}

func resourceHDInsightHadoopCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHDInsightHadoopClusterCreate,
		Read:   resourceHDInsightHadoopClusterRead,
		Update: hdinsightClusterUpdate("Hadoop", resourceHDInsightHadoopClusterRead),
		Delete: hdinsightClusterDelete("Hadoop"),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
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

			"tags": commonschema.Tags(),

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

	id := clusters.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_hdinsight_hadoop_cluster", id.ID())
	}

	gatewayRaw := d.Get("gateway").([]interface{})
	configurations := expandHDInsightsConfigurations(gatewayRaw)

	metastoresRaw := d.Get("metastores").([]interface{})
	metastores := expandHDInsightsMetastore(metastoresRaw)
	for k, v := range metastores {
		configurations[k] = v
	}

	storageAccountsRaw := d.Get("storage_account").([]interface{})
	storageAccountsGen2Raw := d.Get("storage_account_gen2").([]interface{})
	storageAccounts, i, err := expandHDInsightsStorageAccounts(storageAccountsRaw, storageAccountsGen2Raw)
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

	params := clusters.ClusterCreateParametersExtended{
		Location: pointer.To(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &clusters.ClusterCreateProperties{
			Tier:                   pointer.To(clusters.Tier(d.Get("tier").(string))),
			OsType:                 pointer.To(clusters.OSTypeLinux),
			ClusterVersion:         pointer.To(d.Get("cluster_version").(string)),
			MinSupportedTlsVersion: pointer.To(d.Get("tls_min_version").(string)),
			NetworkProperties:      expandHDInsightsNetwork(d.Get("network").([]interface{})),
			ClusterDefinition: &clusters.ClusterDefinition{
				Kind:             pointer.To("Hadoop"),
				ComponentVersion: expandHDInsightHadoopComponentVersion(d.Get("component_version").([]interface{})),
				Configurations:   pointer.To(interface{}(configurations)),
			},
			StorageProfile: &clusters.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &clusters.ComputeProfile{
				Roles: roles,
			},
			ComputeIsolationProperties: expandHDInsightComputeIsolationProperties(d.Get("compute_isolation").([]interface{})),
		},
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Identity: i,
	}

	if diskEncryptionPropertiesRaw, ok := d.GetOk("disk_encryption"); ok {
		diskEncryptionProperties, err := expandHDInsightsDiskEncryptionProperties(diskEncryptionPropertiesRaw.([]interface{}))
		if err != nil {
			return err
		}
		params.Properties.DiskEncryptionProperties = diskEncryptionProperties
	}

	if v, ok := d.GetOk("security_profile"); ok {
		params.Properties.SecurityProfile = expandHDInsightSecurityProfile(v.([]interface{}))

		params.Identity = &identity.SystemAndUserAssignedMap{
			Type:        identity.TypeUserAssigned,
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}

		if params.Properties.SecurityProfile != nil && params.Properties.SecurityProfile.MsiResourceId != nil {
			params.Identity.IdentityIds[*params.Properties.SecurityProfile.MsiResourceId] = identity.UserAssignedIdentityDetails{}
		}
	}

	if err := client.CreateThenPoll(ctx, id, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// We can only add an edge node after creation
	if v, ok := d.GetOk("roles.0.edge_node"); ok {
		edgeNodeRaw := v.([]interface{})
		applicationsClient := meta.(*clients.Client).HDInsight.ApplicationsClient
		edgeNodeConfig := edgeNodeRaw[0].(map[string]interface{})

		err := createHDInsightEdgeNodes(ctx, applicationsClient, id, edgeNodeConfig)
		if err != nil {
			return err
		}

		// we can't rely on the use of the Future here due to the node being successfully completed but now the cluster is applying those changes.
		log.Printf("[DEBUG] Waiting for %s to finish applying edge node", id)
		stateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{"AzureVMConfiguration", "Accepted", "HdInsightConfiguration"},
			Target:     []string{"Running"},
			Refresh:    hdInsightWaitForReadyRefreshFunc(ctx, client, id),
			MinTimeout: 15 * time.Second,
			Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for %s to be running: %s", id, err)
		}
	}

	// We can only enable monitoring after creation
	if v, ok := d.GetOk("monitor"); ok {
		monitorRaw := v.([]interface{})
		if err := enableHDInsightMonitoring(ctx, extensionsClient, id, monitorRaw); err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("extension"); ok {
		extensionRaw := v.([]interface{})
		if err := enableHDInsightAzureMonitor(ctx, extensionsClient, id, extensionRaw); err != nil {
			return err
		}
	}

	return resourceHDInsightHadoopClusterRead(d, meta)
}

func resourceHDInsightHadoopClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	clustersClient := meta.(*clients.Client).HDInsight.ClustersClient
	configurationsClient := meta.(*clients.Client).HDInsight.ConfigurationsClient
	extensionsClient := meta.(*clients.Client).HDInsight.ExtensionsClient
	applicationsClient := meta.(*clients.Client).HDInsight.ApplicationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := clustersClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving HDInsight Hadoop Cluster %q (Resource Group %q): %+v", *id, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configId := configurations.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)
	configs, err := configurationsClient.List(ctx, configId)
	if err != nil {
		return fmt.Errorf("retrieving Configuration for %s: %+v", *id, err)
	}
	if model := configs.Model; model != nil {
		if config := configs.Model.Configurations; config != nil {
			flattenAndSetHDInsightsMetastores(d, *config)

			gateway, exists := (*config)["gateway"]
			if !exists {
				return fmt.Errorf("retrieving gateway for %s: %+v", id, err)
			}

			if err := d.Set("gateway", flattenHDInsightsConfigurations(gateway, d)); err != nil {
				return fmt.Errorf("flattening `gateway`: %+v", err)
			}

		}
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		// storage_account isn't returned so I guess we just leave it ¯\_(ツ)_/¯
		if props := model.Properties; props != nil {
			tier := ""
			// the Azure API is inconsistent here, so rewrite this into the casing we expect
			for _, v := range clusters.PossibleValuesForTier() {
				if strings.EqualFold(v, string(pointer.From(props.Tier))) {
					tier = v
				}
			}
			d.Set("tier", tier)
			d.Set("cluster_version", props.ClusterVersion)
			d.Set("tls_min_version", props.MinSupportedTlsVersion)

			def := props.ClusterDefinition
			if err := d.Set("component_version", flattenHDInsightHadoopComponentVersion(def.ComponentVersion)); err != nil {
				return fmt.Errorf("flattening `component_version`: %+v", err)
			}

			if props.NetworkProperties != nil {
				if err := d.Set("network", FlattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
					return fmt.Errorf("flattening `network`: %+v", err)
				}
			}

			hadoopRoles := hdInsightRoleDefinition{
				HeadNodeDef:      hdInsightHadoopClusterHeadNodeDefinition,
				WorkerNodeDef:    hdInsightHadoopClusterWorkerNodeDefinition,
				ZookeeperNodeDef: hdInsightHadoopClusterZookeeperNodeDefinition,
			}
			flattenedRoles := flattenHDInsightRoles(d, props.ComputeProfile, hadoopRoles)

			appId := applications.NewApplicationID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ClusterName)
			edgeNode, err := applicationsClient.Get(ctx, appId)
			if err != nil {
				if !response.WasNotFound(edgeNode.HttpResponse) {
					return fmt.Errorf("reading edge node for %s: %+v", *id, err)
				}
			}

			if model := edgeNode.Model; model != nil {
				if edgeNodeProps := model.Properties; edgeNodeProps != nil {
					flattenedRoles = flattenHDInsightEdgeNode(flattenedRoles, edgeNodeProps)
				}
			}

			if props.DiskEncryptionProperties != nil {
				diskEncryptionProps, err := flattenHDInsightsDiskEncryptionProperties(*props.DiskEncryptionProperties)
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

			monId := extensions.NewClusterID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)
			monitor, err := extensionsClient.GetMonitoringStatus(ctx, monId)
			if err != nil {
				return fmt.Errorf("reading monitor configuration for %s: %+v", *id, err)
			}

			d.Set("monitor", flattenHDInsightMonitoring(monitor.Model))

			extension, err := extensionsClient.GetAzureMonitorStatus(ctx, monId)
			if err != nil {
				return fmt.Errorf("reading extension configuration for %s: %+v", *id, err)
			}

			d.Set("extension", flattenHDInsightAzureMonitor(extension.Model))

			if err := d.Set("security_profile", flattenHDInsightSecurityProfile(props.SecurityProfile, d)); err != nil {
				return fmt.Errorf("setting `security_profile`: %+v", err)
			}
		}

		if err = tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting tags: %+v", err)
		}
	}

	return nil
}

func flattenHDInsightEdgeNode(roles []interface{}, props *applications.ApplicationProperties) []interface{} {
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
			actions["uri"] = action.Uri
			actions["parameters"] = action.Parameters
		}
	}

	if uninstallScriptActions := props.UninstallScriptActions; uninstallScriptActions != nil && len(*uninstallScriptActions) != 0 {
		uninstallActions := make(map[string]interface{})
		for _, uninstallAction := range *uninstallScriptActions {
			actions["name"] = uninstallAction.Name
			actions["uri"] = uninstallAction.Uri
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

func expandHDInsightHadoopComponentVersion(input []interface{}) *map[string]string {
	vs := input[0].(map[string]interface{})
	return pointer.To(map[string]string{
		"Hadoop": vs["hadoop"].(string),
	})
}

func flattenHDInsightHadoopComponentVersion(input *map[string]string) []interface{} {
	hadoopVersion := ""
	if input != nil {
		if v, ok := (*input)["Hadoop"]; ok {
			hadoopVersion = v
		}
	}

	return []interface{}{
		map[string]interface{}{
			"hadoop": hadoopVersion,
		},
	}
}

func expandHDInsightApplicationEdgeNodeInstallScriptActions(input []interface{}) *[]applications.RuntimeScriptAction {
	actions := make([]applications.RuntimeScriptAction, 0)

	for _, v := range input {
		val := v.(map[string]interface{})

		name := val["name"].(string)
		uri := val["uri"].(string)
		parameters := val["parameters"].(string)

		action := applications.RuntimeScriptAction{
			Name: name,
			Uri:  uri,
			// The only role available for edge nodes is edgenode
			Parameters: pointer.To(parameters),
			Roles:      []string{"edgenode"},
		}

		actions = append(actions, action)
	}

	return &actions
}

func expandHDInsightApplicationEdgeNodeHttpsEndpoints(input []interface{}) *[]applications.ApplicationGetHTTPSEndpoint {
	endpoints := make([]applications.ApplicationGetHTTPSEndpoint, 0)
	if len(input) == 0 || input[0] == nil {
		return &endpoints
	}

	for _, v := range input {
		val := v.(map[string]interface{})

		endPoint := applications.ApplicationGetHTTPSEndpoint{
			AccessModes:        pointer.To(val["access_modes"].([]string)),
			DestinationPort:    pointer.To(val["destination_port"].(int64)),
			PrivateIPAddress:   pointer.To(val["private_ip_address"].(string)),
			SubDomainSuffix:    pointer.To(val["sub_domain_suffix"].(string)),
			DisableGatewayAuth: pointer.To(val["disable_gateway_auth"].(bool)),
		}

		endpoints = append(endpoints, endPoint)
	}

	return &endpoints
}

func expandHDInsightApplicationEdgeNodeUninstallScriptActions(input []interface{}) *[]applications.RuntimeScriptAction {
	actions := make([]applications.RuntimeScriptAction, 0)
	if len(input) == 0 || input[0] == nil {
		return &actions
	}

	for _, v := range input {
		val := v.(map[string]interface{})

		action := applications.RuntimeScriptAction{
			Name:       val["name"].(string),
			Uri:        val["uri"].(string),
			Parameters: pointer.To(val["parameters"].(string)),
			Roles:      []string{"edgenode"},
		}

		actions = append(actions, action)
	}

	return &actions
}

func hdInsightWaitForReadyRefreshFunc(ctx context.Context, client *clusters.ClustersClient, id clusters.ClusterId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "Error", fmt.Errorf("issuing read request in hdInsightWaitForReadyRefreshFunc to %s: %s", id, err)
		}

		if model := res.Model; model != nil {
			if props := model.Properties; props != nil {
				if state := props.ClusterState; state != nil {
					return res, *state, nil
				}
			}
		}

		return res, "Pending", nil
	}
}
