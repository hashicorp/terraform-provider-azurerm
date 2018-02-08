package azurerm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/hdinsight/mgmt/2015-03-01-preview/hdinsight"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmHDInsight() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHDInsightCreate,
		Read:   resourceArmHDInsightRead,
		//	Update: resourceArmHDInsightUpdate,
		Delete: resourceArmHDInsightDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAppServiceName,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			// TODO: (tombuildsstuff) support Update once the API is fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/1697
			"tags": tagsSchema(),

			"cluster_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tier": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Standard",
			},
			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				//TODO:ValidateFunc: validateHDInsightKind,
			},
			"gateway": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
			},
			"cluster_identity": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
			},
			"core_site": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
			},
			"security": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"domain": {
							Type:     schema.TypeString,
							Required: true,
						},
						"organizational_unit_dn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ldaps_urls": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ldaps_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
							Set: resourceArmHDInsightldapsHash,
						},
						"domain_username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"domain_userpassword": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cluster_users_group_dns": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"users_group_dn": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
							Set: resourceArmHDInsightUserGroupDnHash,
						},
					},
				},
				Set: resourceArmHDInsightSecurityHash,
			},
			"node_role": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"count": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validateArmContainerServiceAgentPoolProfileCount,
						},
						"size": {
							Type:     schema.TypeString,
							Required: true,
						},
						"admin_username": {
							Type:     schema.TypeString,
							Required: true,
						},

						"admin_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"ssh_key": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_data": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"vnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subnet_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"numberofdisks": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"scripts": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"uri": {
										Type:     schema.TypeString,
										Required: true,
									},
									"params": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
							Set: resourceArmHDInsightScriptActionHash,
						},
					},
				},
				Set: resourceArmHDInsightnodeHash,
			},
			"storageaccount_profile": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"isdefault": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"container": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceArmHDInsightstorageaccountHash,
			},
		},
	}
}

func resourceArmHDInsightCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	hdiClusterClient := client.hdiClusterClient
	log.Printf("[INFO] preparing arguments for AzureRM HDI creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	clusterVersion := d.Get("cluster_version").(string)
	tags := d.Get("tags").(map[string]interface{})
	RolesProfile, err := expandAzureHDInsightRoleProfile(d)
	storageProfile := expandAzureHDInsightStorageProfile(d)

	tier := hdinsight.Premium
	switch d.Get("tier").(string) {
	case "Prenium":
		tier = hdinsight.Premium
		break
	case "Standard":
		tier = hdinsight.Standard
		break
	}

	clusterDefinition, err := expandAzureHDInsightClusterDefinition(d)

	clusterCreateProperties := hdinsight.ClusterCreateProperties{
		ClusterVersion:    &clusterVersion,
		OsType:            hdinsight.Linux,
		Tier:              tier,
		ClusterDefinition: clusterDefinition,
		ComputeProfile: &hdinsight.ComputeProfile{
			Roles: RolesProfile,
		},
	}

	securityProfile, err := expandAzureHDInsightSecurityProfile(d)
	if securityProfile != nil {
		print(securityProfile)
		clusterCreateProperties.SecurityProfile = securityProfile
	}

	if storageProfile != nil {
		clusterCreateProperties.StorageProfile = storageProfile
	}
	//	fmt.Printf("clusterCreateProperties :%v", clusterCreateProperties)
	parameters := hdinsight.ClusterCreateParametersExtended{
		Location:   &location,
		Tags:       expandTags(tags),
		Properties: &clusterCreateProperties,
	}

	fmt.Printf("Parameters: %+v\n", clusterCreateProperties)
	//ttlInSeconds := "60"
	ctx := meta.(*ArmClient).StopContext

	createFuture, err := hdiClusterClient.Create(ctx, resGroup, name, parameters)

	//fmt.Printf("createFuture: %+v\n", createFuture.Future.Response)
	if err != nil {
		return err
	}

	err = createFuture.WaitForCompletion(ctx, hdiClusterClient.Client)
	if err != nil {
		return err
	}

	read, err := hdiClusterClient.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read HDI %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmHDInsightRead(d, meta)
}

/*
func resourceArmHDInsightUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdiClusterClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["cluster"]

	// if d.HasChange("app_settings") || d.HasChange("version") {
	// 	appSettings := expandHDInsightAppSettings(d)
	// 	settings := web.StringDictionary{
	// 		Properties: appSettings,
	// 	}
	//client.
	// 	_, err := client.UpdateApplicationSettings(ctx, resGroup, name, settings)
	// 	if err != nil {
	// 		return fmt.Errorf("Error updating Application Settings for Function App %q: %+v", name, err)
	// 	}
	// }

	return resourceArmHDInsightRead(d, meta)

}
*/
func resourceArmHDInsightRead(d *schema.ResourceData, meta interface{}) error {
	hdinsightClusterClient := meta.(*ArmClient).hdiClusterClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["clusters"]

	resp, err := hdinsightClusterClient.Get(ctx, resGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure HDInsight %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if resp.Properties != nil {
		if resp.Properties.ClusterVersion != nil {
			d.Set("cluster_version", resp.Properties.ClusterVersion)
		}
		d.Set("tier", resp.Properties.Tier)

		if resp.Properties.ClusterDefinition != nil {
			d.Set("kind", resp.Properties.ClusterDefinition.Kind)
		}

		if resp.Properties.SecurityProfile != nil {
			if err := d.Set("security", schema.NewSet(resourceArmHDInsightSecurityHash, flattenHDInsightSecurityProfile(resp.Properties.SecurityProfile))); err != nil {
				//	return fmt.Errorf()
			}
		}

		if resp.Properties.ComputeProfile != nil && len(*resp.Properties.ComputeProfile.Roles) > 0 {
			if err := d.Set("role_node", schema.NewSet(resourceArmHDInsightnodeHash, flattenHDInsightRoleNode(resp.Properties.ComputeProfile.Roles))); err != nil {
				//	return fmt.Errorf()
			}
		}
		// not implemented in azure-sdk
		// if resp.Properties.StorageProfile != nil {
		// 	if err := d.Set("storageaccount_profile", schema.NewSet(resourceArmHDInsightstorageaccountHash, flattenHDInsightStorageProfile(resp.Properties.StorageProfile))); err != nil {
		// 		//	return fmt.Errorf()
		// 	}
		// }
	} else {
		//return fmt.Errorf()
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmHDInsightDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdiClusterClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["cluster"]

	log.Printf("[DEBUG] Deleting HDi Cluster %q (resource group %q)", name, resGroup)

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		print(resp.Future.Response)
		// if !utils.ResponseWasNotFound(resp.Future.Response) {
		// 	return err
		// }
	}

	return nil
}

func validateHDInsightKind() {

}

func flattenHDInsightStorageProfile(storageProfile *hdinsight.StorageProfile) []interface{} {
	result := make([]interface{}, 0, 1)
	sp := make(map[string]interface{})
	if *storageProfile.Storageaccounts != nil {
		storageAccounts := make([]interface{}, 0, len(*storageProfile.Storageaccounts))
		for _, sto := range *storageProfile.Storageaccounts {
			storage := make(map[string]interface{})
			storage["name"] = sto.Name
			storage["isdefault"] = sto.Name
			storage["container"] = sto.Container
			storage["key"] = sto.Key
			storageAccounts = append(storageAccounts, storage)
		}
		result = append(result, sp)
	}
	return result
}

func flattenHDInsightSecurityProfile(securityProfile *hdinsight.SecurityProfile) []interface{} {
	result := make([]interface{}, 0, 1)

	sp := make(map[string]interface{})

	sp["domain"] = *securityProfile.Domain
	sp["organizational_unit_dn"] = *securityProfile.OrganizationalUnitDN
	sp["domain_username"] = *securityProfile.DomainUsername
	sp["domain_userpassword"] = *securityProfile.DomainUserPassword

	if *securityProfile.LdapsUrls != nil && len(*securityProfile.LdapsUrls) > 0 {
		ldapsUrls := make([]interface{}, 0, len(*securityProfile.LdapsUrls))
		for _, i := range *securityProfile.LdapsUrls {
			url := make(map[string]interface{})
			url["ldaps_url"] = i
			ldapsUrls = append(ldapsUrls, url)
		}
		sp["ldaps_urls"] = ldapsUrls
	}
	if *securityProfile.ClusterUsersGroupDNS != nil && len(*securityProfile.ClusterUsersGroupDNS) > 0 {

		usergroupDNs := make([]interface{}, 0, len(*securityProfile.ClusterUsersGroupDNS))
		for _, i := range *securityProfile.ClusterUsersGroupDNS {
			dn := make(map[string]interface{})
			dn["users_group_dn"] = i
			usergroupDNs = append(usergroupDNs, dn)
		}
		sp["cluster_users_group_dns"] = usergroupDNs
	}
	result = append(result, sp)
	return result

}

func flattenHDInsightRoleNode(roles *[]hdinsight.Role) []interface{} {
	result := make([]interface{}, len(*roles))
	for i, role := range *roles {
		n := make(map[string]interface{})
		n["name"] = *role.Name
		n["count"] = *role.TargetInstanceCount
		n["size"] = *role.HardwareProfile.VMSize
		if role.OsProfile != nil && role.OsProfile.LinuxOperatingSystemProfile != nil {
			n["admin_username"] = *role.OsProfile.LinuxOperatingSystemProfile.Username
			n["admin_password"] = *role.OsProfile.LinuxOperatingSystemProfile.Password

			if role.OsProfile.LinuxOperatingSystemProfile.SSHProfile != nil && len(*role.OsProfile.LinuxOperatingSystemProfile.SSHProfile.PublicKeys) > 0 {

				ssh_keys := make([]map[string]interface{}, 0, len(*role.OsProfile.LinuxOperatingSystemProfile.SSHProfile.PublicKeys))
				for _, u := range *role.OsProfile.LinuxOperatingSystemProfile.SSHProfile.PublicKeys {
					key := make(map[string]interface{})
					key["key_data"] = u
					ssh_keys = append(ssh_keys, key)
				}

				n["ssh_key"] = ssh_keys
			}
		}
		if role.VirtualNetworkProfile != nil {
			n["vnet_id"] = *role.VirtualNetworkProfile.ID
			n["subnet_name"] = *role.VirtualNetworkProfile.Subnet
		}
		if role.DataDisksGroups != nil && len((*role.DataDisksGroups)) > 0 {
			n["numberofdisks"] = (*role.DataDisksGroups)[0].DisksPerNode
		}
		result[i] = n
	}

	return result
}

func expandAzureHDInsightSecurityProfile(d *schema.ResourceData) (*hdinsight.SecurityProfile, error) {
	configs := d.Get("security").(*schema.Set).List()
	if len(configs) > 0 {
		data := configs[0].(map[string]interface{})
		domain := data["domain"].(string)
		oudn := data["organizational_unit_dn"].(string)
		username := data["domain_username"].(string)
		password := data["domain_password"].(string)
		profile := hdinsight.SecurityProfile{
			Domain:               &domain,
			DirectoryType:        hdinsight.ActiveDirectory,
			OrganizationalUnitDN: &oudn,
			DomainUsername:       &username,
			DomainUserPassword:   &password,
		}

		if _, ok := d.GetOk("ldaps_urls"); ok {
			ldapsurl, err := expandAzureRmHDInsightldaps(d)
			if err != nil {
				return nil, err
			}
			if ldapsurl != nil {
				profile.LdapsUrls = ldapsurl
			}
		}

		if _, ok := d.GetOk("cluster_users_group_dns"); ok {
			clusterUsersGroupDNs, err := expandAzureRmHDInsightUserGroupDNs(d)
			if err != nil {
				return nil, err
			}
			if clusterUsersGroupDNs != nil {
				profile.ClusterUsersGroupDNS = clusterUsersGroupDNs
			}
		}
		return &profile, nil
	}
	return nil, nil

}

func expandAzureHDInsightClusterDefinition(d *schema.ResourceData) (*hdinsight.ClusterDefinition, error) {
	kind := d.Get("kind").(string)
	configuration, err := expandAzureHDInsightConfiguration(d)
	if err != nil {
		return nil, err
	}
	clusterDefinition := hdinsight.ClusterDefinition{
		Kind:           &kind,
		Configurations: configuration,
	}
	return &clusterDefinition, nil
}

func ValueToString(v interface{}) (string, error) {
	switch value := v.(type) {
	case string:
		return value, nil
	case int:
		return fmt.Sprintf("%d", value), nil
	default:
		return "", fmt.Errorf("unknown type %T in tag value", value)
	}
}

func expandAzureHDInsightConfiguration(d *schema.ResourceData) (*map[string]interface{}, error) {
	result := make(map[string]interface{})

	gateway := d.Get("gateway").(map[string]interface{})
	gatewayArray := make(map[string]*string, len(gateway))

	for k, v := range gateway {
		value, err := ValueToString(v)
		if err != nil {
			return nil, err
		}

		gatewayArray[k] = &value
	}
	gatewayByte, _ := json.Marshal(gatewayArray)
	result["gateways"] = string(gatewayByte)

	clusterIdentity := d.Get("cluster_identity").(map[string]interface{})
	if clusterIdentity != nil {
		clusterIdentityArray := make(map[string]*string, len(clusterIdentity))
		for k, v := range clusterIdentity {
			value, _ := ValueToString(v)
			clusterIdentityArray[k] = &value
		}
		clusterIdentityByte, _ := json.Marshal(clusterIdentityArray)
		result["clusterIdentity"] = string(clusterIdentityByte)
	}

	coreSite := d.Get("core_site").(map[string]interface{})
	if coreSite != nil {
		coreSiteArray := make(map[string]*string, len(coreSite))
		for k, v := range coreSite {
			value, _ := ValueToString(v)
			coreSiteArray[k] = &value
		}
		coreSiteByte, err := json.Marshal(coreSiteArray)
		if err != nil {
			return nil, err
		}
		result["core-site"] = string(coreSiteByte)
	}
	return &result, nil
}

func expandAzureHDInsightRoleProfile(d *schema.ResourceData) (*[]hdinsight.Role, error) {
	configs := d.Get("node_role").(*schema.Set).List()
	roleprofiles := make([]hdinsight.Role, 0, len(configs))

	for _, roleconf := range configs {
		config := roleconf.(map[string]interface{})

		name := config["name"].(string)

		linuxKeys := config["ssh_key"].(*schema.Set).List()
		sshPublicKeys := []hdinsight.SSHPublicKey{}
		if len(linuxKeys) > 0 {
			key := linuxKeys[0].(map[string]interface{})
			keyData := key["key_data"].(string)

			sshPublicKey := hdinsight.SSHPublicKey{
				CertificateData: &keyData,
			}

			sshPublicKeys = append(sshPublicKeys, sshPublicKey)
		}

		numberofdisks := int32(config["numberofdisks"].(int))
		dataDiskGroups := []hdinsight.DataDisksGroups{}
		if numberofdisks > 0 {
			dataDiskGroup := hdinsight.DataDisksGroups{
				DisksPerNode: &numberofdisks,
			}
			dataDiskGroups = append(dataDiskGroups, dataDiskGroup)
		}
		minInstance := int32(1)
		targetInstance := int32(config["count"].(int))

		vnetid := config["vnet_id"].(string)
		fmt.Printf("----vnetid : %s", vnetid)
		subnet := config["subnet_name"].(string)
		fmt.Printf("----subnet : %s", subnet)
		vmsize := config["size"].(string)
		adminUsername := config["admin_username"].(string)
		adminPassword := config["admin_password"].(string)
		profile := hdinsight.Role{
			Name:                &name,
			MinInstanceCount:    &minInstance,
			TargetInstanceCount: &targetInstance,
			HardwareProfile: &hdinsight.HardwareProfile{
				VMSize: &vmsize,
			},
			OsProfile: &hdinsight.OsProfile{
				LinuxOperatingSystemProfile: &hdinsight.LinuxOperatingSystemProfile{
					Username: &adminUsername,
					Password: &adminPassword,
				},
			},
			VirtualNetworkProfile: &hdinsight.VirtualNetworkProfile{
				ID:     &vnetid,
				Subnet: &subnet,
			},
		}

		if len(dataDiskGroups) > 0 {
			profile.DataDisksGroups = &dataDiskGroups
		}
		if len(sshPublicKeys) > 0 {
			profile.OsProfile.LinuxOperatingSystemProfile.SSHProfile = &hdinsight.SSHProfile{
				PublicKeys: &sshPublicKeys,
			}
		}

		if _, ok := d.GetOk("scripts"); ok {
			scriptactions, err := expandAzureRmHDInsightScriptAction(d)
			if err != nil {
				return nil, err
			}
			profile.ScriptActions = &scriptactions
		}
		roleprofiles = append(roleprofiles, profile)
	}

	return &roleprofiles, nil
}

func expandAzureRmHDInsightUserGroupDNs(d *schema.ResourceData) (*[]string, error) {
	userGroupDNs := d.Get("cluster_users_group_dns").([]interface{})
	clusterUserGroupDNs := make([]string, 0, len(userGroupDNs))
	for _, userGroupDNsConf := range userGroupDNs {
		config := userGroupDNsConf.(map[string]interface{})
		clusterUserGroupDNs = append(clusterUserGroupDNs, config["users_group_dn"].(string))
	}
	return &clusterUserGroupDNs, nil
}

func expandAzureRmHDInsightldaps(d *schema.ResourceData) (*[]string, error) {
	ldaps := d.Get("ldaps_urls").([]interface{})
	ldapsUrls := make([]string, 0, len(ldaps))

	for _, ldapconf := range ldaps {
		config := ldapconf.(map[string]interface{})
		ldapsUrls = append(ldapsUrls, config["ldaps_url"].(string))
	}
	return &ldapsUrls, nil
}

func expandAzureRmHDInsightScriptAction(d *schema.ResourceData) ([]hdinsight.ScriptAction, error) {
	scripts := d.Get("scripts").([]interface{})
	scriptActions := make([]hdinsight.ScriptAction, 0, len(scripts))
	for _, scriptconfig := range scripts {
		config := scriptconfig.(map[string]interface{})

		name := config["name"].(string)
		uri := config["uri"].(string)
		params := config["params"].(string)

		scriptAction := hdinsight.ScriptAction{
			Name:       &name,
			URI:        &uri,
			Parameters: &params,
		}
		scriptActions = append(scriptActions, scriptAction)
	}
	return scriptActions, nil
}

func expandAzureHDInsightStorageProfile(d *schema.ResourceData) *hdinsight.StorageProfile {
	configs := d.Get("storageaccount_profile").(*schema.Set).List()
	storageAccounts := make([]hdinsight.StorageAccount, 0, len(configs))
	for _, storageconfig := range configs {
		config := storageconfig.(map[string]interface{})
		name := config["name"].(string)
		container := config["container"].(string)
		key := config["key"].(string)
		isdefault := bool(config["isdefault"].(bool))
		storageAccountprofile := hdinsight.StorageAccount{
			Container: &container,
			Key:       &key,
			Name:      &name,
			IsDefault: &isdefault,
		}
		storageAccounts = append(storageAccounts, storageAccountprofile)
	}

	if len(storageAccounts) > 0 {
		storageProfile := hdinsight.StorageProfile{
			Storageaccounts: &storageAccounts,
		}
		return &storageProfile
	}
	return nil
}

/* func getBasicHDInsightAppSettings(d *schema.ResourceData) []web.NameValuePair {
	dashboardPropName := "AzureWebJobsDashboard"
	storagePropName := "AzureWebJobsStorage"
	functionVersionPropName := "FUNCTIONS_EXTENSION_VERSION"
	contentSharePropName := "WEBSITE_CONTENTSHARE"
	contentFileConnStringPropName := "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"

	storageConnection := d.Get("storage_connection_string").(string)
	functionVersion := d.Get("version").(string)
	contentShare := d.Get("name").(string) + "-content"

	return []web.NameValuePair{
		{Name: &dashboardPropName, Value: &storageConnection},
		{Name: &storagePropName, Value: &storageConnection},
		{Name: &functionVersionPropName, Value: &functionVersion},
		{Name: &contentSharePropName, Value: &contentShare},
		{Name: &contentFileConnStringPropName, Value: &storageConnection},
	}
}*/

func resourceArmHDInsightvirtualnetworkprofileHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	vnetid := m["vnet_id"].(string)
	subnetname := m["subnet_name"].(string)
	buf.WriteString(fmt.Sprintf("%s-", vnetid))
	buf.WriteString(fmt.Sprintf("%s-", subnetname))
	return hashcode.String(buf.String())
}

func resourceArmHDInsightUserGroupDnHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	usersGroupDN := m["users_group_dn"].(string)
	buf.WriteString(fmt.Sprintf("%s-", usersGroupDN))
	return hashcode.String(buf.String())
}

func resourceArmHDInsightldapsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	url := m["ldaps_url"].(string)
	buf.WriteString(fmt.Sprintf("%s-", url))
	return hashcode.String(buf.String())
}

func resourceArmHDInsightSecurityHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	domain := m["domain"].(string)
	ou := m["organizational_unit_dn"].(string)
	buf.WriteString(fmt.Sprintf("%s-", domain))
	buf.WriteString(fmt.Sprintf("%s-", ou))
	return hashcode.String(buf.String())
}

func resourceArmHDInsightOSProfileHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	adminUsername := m["admin_username"].(string)

	buf.WriteString(fmt.Sprintf("%s-", adminUsername))

	return hashcode.String(buf.String())
}

func resourceArmHDInsightnodeHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	count := m["count"].(int)
	name := m["name"].(string)

	buf.WriteString(fmt.Sprintf("%d-", count))
	buf.WriteString(fmt.Sprintf("%s-", name))

	return hashcode.String(buf.String())
}

func resourceArmHDInsightScriptActionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	count := m["count"].(int)
	name := m["name"].(string)

	buf.WriteString(fmt.Sprintf("%d-", count))
	buf.WriteString(fmt.Sprintf("%s-", name))

	return hashcode.String(buf.String())
}

func resourceArmHDInsightstorageaccountHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	count := m["count"].(int)
	name := m["name"].(string)
	container := m["container"].(string)

	buf.WriteString(fmt.Sprintf("%d-", count))
	buf.WriteString(fmt.Sprintf("%s-", name))
	buf.WriteString(fmt.Sprintf("%s-", container))

	return hashcode.String(buf.String())
}
