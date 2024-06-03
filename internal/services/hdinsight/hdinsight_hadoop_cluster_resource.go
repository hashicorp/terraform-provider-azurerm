// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/applications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01/clusters"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/custompollers"
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
	MaxInstanceCount:         pointer.To(2),
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
	MaxInstanceCount:         pointer.To(3),
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
			_, err := commonids.ParseHDInsightClusterID(id)
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

			"private_link_configuration": SchemaHDInsightPrivateLinkConfigurations(),

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
	client := meta.(*clients.Client).HDInsight.Clusters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	extensionsClient := meta.(*clients.Client).HDInsight.Extensions
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewHDInsightClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := location.Normalize(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	t := d.Get("tags").(map[string]interface{})
	tier := clusters.Tier(d.Get("tier").(string))
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

	privateLinkConfigurationsRaw := d.Get("private_link_configuration").([]interface{})
	privateLinkConfigurations := ExpandHDInsightPrivateLinkConfigurations(privateLinkConfigurationsRaw)

	computeIsolationProperties := ExpandHDInsightComputeIsolationProperties(d.Get("compute_isolation").([]interface{}))

	storageAccountsRaw := d.Get("storage_account").([]interface{})
	storageAccountsGen2Raw := d.Get("storage_account_gen2").([]interface{})
	storageAccounts, expandedIdentity, err := ExpandHDInsightsStorageAccounts(storageAccountsRaw, storageAccountsGen2Raw)
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

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing Hadoop %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_hdinsight_hadoop_cluster", id.ID())
	}

	var configurationsRaw interface{} = configurations
	payload := clusters.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &clusters.ClusterCreateProperties{
			Tier:                      pointer.To(tier),
			OsType:                    pointer.To(clusters.OSTypeLinux),
			ClusterVersion:            utils.String(clusterVersion),
			MinSupportedTlsVersion:    utils.String(tls),
			NetworkProperties:         networkProperties,
			PrivateLinkConfigurations: privateLinkConfigurations,
			ClusterDefinition: &clusters.ClusterDefinition{
				Kind:             pointer.To(clusters.ClusterKindHadoop),
				ComponentVersion: pointer.To(componentVersions),
				Configurations:   pointer.To(configurationsRaw),
			},
			StorageProfile: &clusters.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &clusters.ComputeProfile{
				Roles: roles,
			},
			ComputeIsolationProperties: computeIsolationProperties,
		},
		Tags:     tags.Expand(t),
		Identity: expandedIdentity,
	}

	if diskEncryptionPropertiesRaw, ok := d.GetOk("disk_encryption"); ok {
		diskEncryptionProperties, err := ExpandHDInsightsDiskEncryptionProperties(diskEncryptionPropertiesRaw.([]interface{}))
		if err != nil {
			return err
		}
		payload.Properties.DiskEncryptionProperties = diskEncryptionProperties
	}

	if v, ok := d.GetOk("security_profile"); ok {
		payload.Properties.SecurityProfile = ExpandHDInsightSecurityProfile(v.([]interface{}))

		// @tombuildsstuff: this behaviour is likely wrong and wants reevaluating - users should need to explicitly define this in the config?
		payload.Identity = &identity.SystemAndUserAssignedMap{
			Type:        identity.TypeUserAssigned,
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if payload.Properties.SecurityProfile != nil && payload.Properties.SecurityProfile.MsiResourceId != nil {
			payload.Identity.IdentityIds[*payload.Properties.SecurityProfile.MsiResourceId] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}
	}

	if err := client.CreateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating Hadoop %s: %+v", id, err)
	}
	d.SetId(id.ID())

	// We can only add an edge node after creation
	if v, ok := d.GetOk("roles.0.edge_node"); ok {
		edgeNodeRaw := v.([]interface{})
		applicationsClient := meta.(*clients.Client).HDInsight.Applications
		edgeNodeConfig := edgeNodeRaw[0].(map[string]interface{})
		applicationId := applications.NewApplicationID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ClusterName) // 2 id.ClusterName's are intentional

		err := createHDInsightEdgeNodes(ctx, applicationsClient, applicationId, edgeNodeConfig)
		if err != nil {
			return err
		}

		// we can't rely on the use of the Future here due to the node being successfully completed but now the cluster is applying those changes.
		log.Printf("[DEBUG] Waiting for the Hadoop %s to finish applying edge nodes..", id)
		pollerType := custompollers.NewEdgeNodePoller(client, id)
		poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the Hadoop %s to finish applying the edge nodes: %+v", id, err)
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
	clustersClient := meta.(*clients.Client).HDInsight.Clusters
	configurationsClient := meta.(*clients.Client).HDInsight.Configurations
	extensionsClient := meta.(*clients.Client).HDInsight.Extensions
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseHDInsightClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := clustersClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Hadoop %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Hadoop %s: %+v", id, err)
	}

	// Each call to configurationsClient methods is HTTP request. Getting all settings in one operation
	configurationsResp, err := configurationsClient.List(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Configurations for HBase %s: %+v", id, err)
	}

	configurations := make(map[string]map[string]string)
	if model := configurationsResp.Model; model != nil && model.Configurations != nil {
		configurations = *model.Configurations
	}

	gateway, exists := configurations["gateway"]
	if !exists {
		return fmt.Errorf("retrieving Gateway Configuration for Hadoop %s: %+v", id, err)
	}

	applicationsClient := meta.(*clients.Client).HDInsight.Applications
	applicationId := applications.NewApplicationID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ClusterName) // 2 id.ClusterName's are intentional
	edgeNode, err := applicationsClient.Get(ctx, applicationId)
	if err != nil {
		if !response.WasNotFound(edgeNode.HttpResponse) {
			return fmt.Errorf("reading Edge Node for Hadoop %s: %+v", id, err)
		}
	}

	extension, err := extensionsClient.GetAzureMonitorStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Azure Monitor Status for Hadoop %s: %+v", id, err)
	}
	monitor, err := extensionsClient.GetMonitoringStatus(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Monitoring Status for Hadoop %s: %+v", id, err)
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		// storage_account isn't returned so I guess we just leave it ¯\_(ツ)_/¯
		if props := model.Properties; props != nil {
			d.Set("tier", string(pointer.From(props.Tier)))
			d.Set("cluster_version", props.ClusterVersion)
			d.Set("tls_min_version", props.MinSupportedTlsVersion)

			if err := d.Set("component_version", flattenHDInsightHadoopComponentVersion(props.ClusterDefinition.ComponentVersion)); err != nil {
				return fmt.Errorf("flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", FlattenHDInsightsConfigurations(gateway, d)); err != nil {
				return fmt.Errorf("flattening `gateway`: %+v", err)
			}

			flattenHDInsightsMetastores(d, configurations)

			if err := d.Set("network", flattenHDInsightsNetwork(props.NetworkProperties)); err != nil {
				return fmt.Errorf("flattening `network`: %+v", err)
			}

			if err := d.Set("private_link_configuration", flattenHDInsightPrivateLinkConfigurations(props.PrivateLinkConfigurations)); err != nil {
				return fmt.Errorf("flattening `private_link_configuration`: %+v", err)
			}

			hadoopRoles := hdInsightRoleDefinition{
				HeadNodeDef:      hdInsightHadoopClusterHeadNodeDefinition,
				WorkerNodeDef:    hdInsightHadoopClusterWorkerNodeDefinition,
				ZookeeperNodeDef: hdInsightHadoopClusterZookeeperNodeDefinition,
			}
			flattenedRoles := flattenHDInsightRoles(d, props.ComputeProfile, hadoopRoles)
			if edgeNode.Model != nil && edgeNode.Model.Properties != nil {
				flattenedRoles = flattenHDInsightEdgeNode(flattenedRoles, edgeNode.Model.Properties, d)
			}

			diskEncryptionProps, err := flattenHDInsightsDiskEncryptionProperties(props.DiskEncryptionProperties)
			if err != nil {
				return err
			}
			if err := d.Set("disk_encryption", diskEncryptionProps); err != nil {
				return fmt.Errorf("flattening `disk_encryption`: %+v", err)
			}

			if err := d.Set("compute_isolation", flattenHDInsightComputeIsolationProperties(props.ComputeIsolationProperties)); err != nil {
				return fmt.Errorf("failed setting `compute_isolation`: %+v", err)
			}

			if err := d.Set("roles", flattenedRoles); err != nil {
				return fmt.Errorf("flattening `roles`: %+v", err)
			}

			httpEndpoint := findHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
			d.Set("https_endpoint", httpEndpoint)
			sshEndpoint := findHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
			d.Set("ssh_endpoint", sshEndpoint)

			if err := d.Set("security_profile", flattenHDInsightSecurityProfile(props.SecurityProfile, d)); err != nil {
				return fmt.Errorf("setting `security_profile`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	d.Set("extension", flattenHDInsightAzureMonitor(extension.Model))

	d.Set("monitor", flattenHDInsightMonitoring(monitor.Model))

	return nil
}

func flattenHDInsightEdgeNode(roles []interface{}, props *applications.ApplicationProperties, d *pluginsdk.ResourceData) []interface{} {
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

	installScriptActions := make([]interface{}, 0)
	if props.InstallScriptActions != nil {
		for _, action := range *props.InstallScriptActions {
			installScriptActions = append(installScriptActions, map[string]interface{}{
				"name":       action.Name,
				"uri":        action.Uri,
				"parameters": d.Get("roles.0.edge_node.0.install_script_action.0.parameters").(string),
			})
		}
	}

	uninstallScriptActions := make([]interface{}, 0)
	if props.UninstallScriptActions != nil {
		for _, uninstallAction := range *props.UninstallScriptActions {
			uninstallScriptActions = append(uninstallScriptActions, map[string]interface{}{
				"name":       uninstallAction.Name,
				"uri":        uninstallAction.Uri,
				"parameters": d.Get("roles.0.edge_node.0.uninstall_script_actions.0.parameters").(string),
			})
		}
	}

	httpsEndpoints := make([]interface{}, 0)
	if HTTPSEndpoints := props.HTTPSEndpoints; HTTPSEndpoints != nil && len(*HTTPSEndpoints) != 0 {
		for _, HTTPSEndpoint := range *HTTPSEndpoints {
			httpsEndpoints = append(httpsEndpoints, map[string]interface{}{
				"access_modes":         pointer.From(HTTPSEndpoint.AccessModes),
				"destination_port":     pointer.From(HTTPSEndpoint.DestinationPort),
				"disable_gateway_auth": pointer.From(HTTPSEndpoint.DisableGatewayAuth),
				"private_ip_address":   pointer.From(HTTPSEndpoint.PrivateIPAddress),
				"sub_domain_suffix":    pointer.From(HTTPSEndpoint.SubDomainSuffix),
			})
		}
	}

	role["edge_node"] = []interface{}{
		map[string]interface{}{
			"install_script_action":    installScriptActions,
			"https_endpoints":          httpsEndpoints,
			"uninstall_script_actions": uninstallScriptActions,
		},
	}

	return []interface{}{role}
}

func expandHDInsightHadoopComponentVersion(input []interface{}) map[string]string {
	vs := input[0].(map[string]interface{})
	return map[string]string{
		"Hadoop": vs["hadoop"].(string),
	}
}

func flattenHDInsightHadoopComponentVersion(input *map[string]string) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		hadoopVersion := ""
		if v, ok := (*input)["Hadoop"]; ok {
			hadoopVersion = v
		}
		output = append(output, map[string]interface{}{
			"hadoop": hadoopVersion,
		})
	}

	return output
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
			Parameters: utils.String(parameters),
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

		accessModes := val["access_modes"].([]string)
		destinationPort := val["destination_port"].(int64)
		disableGatewayAuth := val["disable_gateway_auth"].(bool)
		privateIpAddress := val["private_ip_address"].(string)
		subDomainSuffix := val["sub_domain_suffix"].(string)

		endPoint := applications.ApplicationGetHTTPSEndpoint{
			AccessModes:        &accessModes,
			DestinationPort:    pointer.To(destinationPort),
			PrivateIPAddress:   utils.String(privateIpAddress),
			SubDomainSuffix:    utils.String(subDomainSuffix),
			DisableGatewayAuth: utils.Bool(disableGatewayAuth),
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

		name := val["name"].(string)
		uri := val["uri"].(string)
		parameters := val["parameters"].(string)

		action := applications.RuntimeScriptAction{
			Name:       name,
			Uri:        uri,
			Parameters: utils.String(parameters),
			Roles:      []string{"edgenode"},
		}

		actions = append(actions, action)
	}

	return &actions
}
