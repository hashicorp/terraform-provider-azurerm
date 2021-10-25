package redhatopenshift

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/redhatopenshift/mgmt/2020-04-30/redhatopenshift"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redhatopenshift/parse"
	openShiftValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/redhatopenshift/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var (
	randomDomainName = GenerateRandomDomainName()
)

func resourceOpenShiftCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceOpenShiftClusterCreate,
		Read:   resourceOpenShiftClusterRead,
		Update: resourceOpenShiftClusterUpdate,
		Delete: resourceOpenShiftClusterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ClusterID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"cluster_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"pull_secret": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"domain": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"resource_group_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"service_principal": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"client_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: openShiftValidate.ClientID,
						},
						"client_secret": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"network_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"pod_cidr": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "10.128.0.0/14",
							ValidateFunc: validate.CIDR,
						},
						"service_cidr": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "172.30.0.0/16",
							ValidateFunc: validate.CIDR,
						},
					},
				},
			},

			"master_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"vm_size": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.StringIsNotEmpty,
						},
					},
				},
			},

			"worker_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"vm_size": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.StringIsNotEmpty,
						},
						"disk_size_gb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: openShiftValidate.DiskSizeGB,
						},
						"node_count": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(3, 20),
						},
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"api_server_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"visibility": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      redhatopenshift.Public,
							ValidateFunc: validate.CIDR,
						},
					},
				},
			},

			"ingress_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "default",
						},
						"visibility": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      redhatopenshift.Visibility1Public,
							ValidateFunc: validate.CIDR,
						},
					},
				},
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"console_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceOpenShiftClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedHatOpenshift.OpenShiftClustersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Red Hat Openshift Cluster create.")

	resourceGroupName := d.Get("resource_group_name").(string)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	name := d.Get("name").(string)

	existing, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Red Hat Openshift Cluster %q (Resource Group %q): %s", name, resourceGroupName, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_redhatopenshift_cluster", *existing.ID)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	clusterProfileRaw := d.Get("cluster_profile").([]interface{})
	clusterProfile := expandOpenshiftClusterProfile(clusterProfileRaw, subscriptionId)

	consoleProfile := &redhatopenshift.ConsoleProfile{}

	servicePrincipalProfileRaw := d.Get("service_principal").([]interface{})
	servicePrincipalProfile := expandOpenshiftServicePrincipalProfile(servicePrincipalProfileRaw)

	networkProfileRaw := d.Get("network_profile").([]interface{})
	networkProfile := expandOpenshiftNetworkProfile(networkProfileRaw)

	masterProfileRaw := d.Get("master_profile").([]interface{})
	masterProfile := expandOpenshiftMasterProfile(masterProfileRaw)

	workerProfilesRaw := d.Get("worker_profile").([]interface{})
	workerProfiles := expandOpenshiftWorkerProfiles(workerProfilesRaw)

	apiServerProfileRaw := d.Get("api_server_profile").([]interface{})
	apiServerProfile := expandOpenshiftApiServerProfile(apiServerProfileRaw)

	ingressProfilesRaw := d.Get("ingress_profile").([]interface{})
	ingressProfiles := expandOpenshiftIngressProfiles(ingressProfilesRaw)

	t := d.Get("tags").(map[string]interface{})

	parameters := redhatopenshift.OpenShiftCluster{
		Name:     &name,
		Location: &location,
		OpenShiftClusterProperties: &redhatopenshift.OpenShiftClusterProperties{
			ClusterProfile:          clusterProfile,
			ConsoleProfile:          consoleProfile,
			ServicePrincipalProfile: servicePrincipalProfile,
			NetworkProfile:          networkProfile,
			MasterProfile:           masterProfile,
			WorkerProfiles:          workerProfiles,
			ApiserverProfile:        apiServerProfile,
			IngressProfiles:         ingressProfiles,
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Red Hat OpenShift Cluster %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Red Hat OpenShift Cluster %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("retrieving Red Hat OpenShift Cluster %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read ID for Red Hat OpenShift Cluster %q (Resource Group %q)", name, resourceGroupName)
	}

	d.SetId(*read.ID)

	return resourceOpenShiftClusterRead(d, meta)
}

func resourceOpenShiftClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedHatOpenshift.OpenShiftClustersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Red Hat OpenShift Cluster update.")

	resourceGroupName := d.Get("resource_group_name").(string)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resourceGroupId := ResourceGroupID(subscriptionId, resourceGroupName)

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)

	existing, err := client.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		return fmt.Errorf("retrieving existing Red Hat OpenShift Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
	}
	if existing.OpenShiftClusterProperties == nil {
		return fmt.Errorf("retrieving existing Red Hat OpenShift Cluster %q (Resource Group %q): `properties` was nil", id.ManagedClusterName, id.ResourceGroup)
	}

	if d.HasChange("cluster_profile") {
		clusterProfileRaw := d.Get("cluster_profile").([]interface{})
		clusterProfile := expandOpenshiftClusterProfile(clusterProfileRaw, resourceGroupId)
		existing.OpenShiftClusterProperties.ClusterProfile = clusterProfile
	}

	if d.HasChange("master_profile") {
		masterProfileRaw := d.Get("master_profile").([]interface{})
		masterProfile := expandOpenshiftMasterProfile(masterProfileRaw)
		existing.OpenShiftClusterProperties.MasterProfile = masterProfile
	}

	if d.HasChange("worker_profile") {
		workerProfilesRaw := d.Get("worker_profile").([]interface{})
		workerProfiles := expandOpenshiftWorkerProfiles(workerProfilesRaw)
		existing.OpenShiftClusterProperties.WorkerProfiles = workerProfiles
	}

	d.Partial(false)

	return resourceOpenShiftClusterRead(d, meta)
}

func resourceOpenShiftClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedHatOpenshift.OpenShiftClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Red Hat OpenShift Cluster %q was not found in Resource Group %q - removing from state!", id.ManagedClusterName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Red Hat OpenShift Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.OpenShiftClusterProperties; props != nil {
		clusterProfile := flattenOpenShiftClusterProfile(props.ClusterProfile)
		if err := d.Set("cluster_profile", clusterProfile); err != nil {
			return fmt.Errorf("setting `cluster_profile`: %+v", err)
		}

		servicePrincipalProfile := flattenOpenShiftServicePrincipalProfile(props.ServicePrincipalProfile, d)
		if err := d.Set("service_principal", servicePrincipalProfile); err != nil {
			return fmt.Errorf("setting `service_principal`: %+v", err)
		}

		networkProfile := flattenOpenShiftNetworkProfile(props.NetworkProfile)
		if err := d.Set("network_profile", networkProfile); err != nil {
			return fmt.Errorf("setting `network_profile`: %+v", err)
		}

		masterProfile := flattenOpenShiftMasterProfile(props.MasterProfile)
		if err := d.Set("master_profile", masterProfile); err != nil {
			return fmt.Errorf("setting `master_profile`: %+v", err)
		}

		apiServerProfile := flattenOpenShiftAPIServerProfile(props.ApiserverProfile)
		if err := d.Set("api_server_profile", apiServerProfile); err != nil {
			return fmt.Errorf("setting `api_server_profile`: %+v", err)
		}

		ingressProfiles := flattenOpenShiftIngressProfiles(props.IngressProfiles)
		if err := d.Set("ingress_profile", ingressProfiles); err != nil {
			return fmt.Errorf("setting `ingress_profile`: %+v", err)
		}

		d.Set("version", props.ClusterProfile.Version)
		d.Set("console_url", props.ConsoleProfile.URL)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceOpenShiftClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedHatOpenshift.OpenShiftClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		return fmt.Errorf("deleting Red Hat Openshift Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of Red Hat Openshift Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
	}

	return nil
}

func flattenOpenShiftClusterProfile(profile *redhatopenshift.ClusterProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	pullSecret := ""
	if profile.PullSecret != nil {
		pullSecret = *profile.PullSecret
	}

	clusterDomain := ""
	if profile.Domain != nil {
		clusterDomain = *profile.Domain
	}

	clusterResourceGroupId := ""
	if profile.ResourceGroupID != nil {
		clusterResourceGroupId = *profile.ResourceGroupID
	}

	return []interface{}{
		map[string]interface{}{
			"pull_secret":       pullSecret,
			"domain":            clusterDomain,
			"resource_group_id": clusterResourceGroupId,
		},
	}
}

func flattenOpenShiftServicePrincipalProfile(profile *redhatopenshift.ServicePrincipalProfile, d *pluginsdk.ResourceData) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	clientID := ""
	if profile.ClientID != nil {
		clientID = *profile.ClientID
	}

	// client secret isn't returned by the API so pass the existing value along
	clientSecret := ""
	if sp, ok := d.GetOk("service_principal"); ok {
		var val []interface{}

		// prior to 1.34 this was a *pluginsdk.Set, now it's a List - try both
		if v, ok := sp.([]interface{}); ok {
			val = v
		} else if v, ok := sp.(*pluginsdk.Set); ok {
			val = v.List()
		}

		if len(val) > 0 && val[0] != nil {
			raw := val[0].(map[string]interface{})
			clientSecret = raw["client_secret"].(string)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"client_id":     clientID,
			"client_secret": clientSecret,
		},
	}
}

func flattenOpenShiftNetworkProfile(profile *redhatopenshift.NetworkProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	podCidr := ""
	if profile.PodCidr != nil {
		podCidr = *profile.PodCidr
	}

	serviceCidr := ""
	if profile.ServiceCidr != nil {
		serviceCidr = *profile.ServiceCidr
	}

	return []interface{}{
		map[string]interface{}{
			"pod_cidr":     podCidr,
			"service_cidr": serviceCidr,
		},
	}
}

func flattenOpenShiftMasterProfile(profile *redhatopenshift.MasterProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	subnetId := ""
	if profile.SubnetID != nil {
		subnetId = *profile.SubnetID
	}

	return []interface{}{
		map[string]interface{}{
			"vm_size":   string(profile.VMSize),
			"subnet_id": subnetId,
		},
	}
}

func flattenOpenShiftAPIServerProfile(profile *redhatopenshift.APIServerProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"visibility": string(profile.Visibility),
		},
	}
}

func flattenOpenShiftIngressProfiles(profiles *[]redhatopenshift.IngressProfile) []interface{} {
	if profiles == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, profile := range *profiles {
		result := make(map[string]interface{})
		result["visibility"] = string(profile.Visibility)

		results = append(results, result)
	}

	return results
}

func expandOpenshiftClusterProfile(input []interface{}, subscriptionId string) *redhatopenshift.ClusterProfile {
	if len(input) == 0 {
		randomClusterResourceGroupName := fmt.Sprintf("aro-%s", randomDomainName)
		randomClusterResourceGroupId := ResourceGroupID(subscriptionId, randomClusterResourceGroupName)

		return &redhatopenshift.ClusterProfile{
			ResourceGroupID: utils.String(randomClusterResourceGroupId),
			Domain:          utils.String(randomDomainName),
		}
	}

	config := input[0].(map[string]interface{})

	pullSecret := config["pull_secret"].(string)

	domain := config["domain"].(string)
	if domain == "" {
		domain = randomDomainName
	}

	resourceGroupId := config["resource_group_id"].(string)
	if resourceGroupId == "" {
		randomClusterResourceGroupName := fmt.Sprintf("aro-%s", randomDomainName)
		resourceGroupId = ResourceGroupID(subscriptionId, randomClusterResourceGroupName)
	}

	return &redhatopenshift.ClusterProfile{
		ResourceGroupID: utils.String(resourceGroupId),
		Domain:          utils.String(domain),
		PullSecret:      utils.String(pullSecret),
	}
}

func expandOpenshiftServicePrincipalProfile(input []interface{}) *redhatopenshift.ServicePrincipalProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

	clientId := config["client_id"].(string)
	clientSecret := config["client_secret"].(string)

	return &redhatopenshift.ServicePrincipalProfile{
		ClientID:     utils.String(clientId),
		ClientSecret: utils.String(clientSecret),
	}
}

func expandOpenshiftNetworkProfile(input []interface{}) *redhatopenshift.NetworkProfile {
	if len(input) == 0 {
		return &redhatopenshift.NetworkProfile{
			PodCidr:     utils.String("10.128.0.0/14"),
			ServiceCidr: utils.String("172.30.0.0/16"),
		}
	}

	config := input[0].(map[string]interface{})

	podCidr := config["pod_cidr"].(string)
	serviceCidr := config["service_cidr"].(string)

	return &redhatopenshift.NetworkProfile{
		PodCidr:     utils.String(podCidr),
		ServiceCidr: utils.String(serviceCidr),
	}
}

func expandOpenshiftMasterProfile(input []interface{}) *redhatopenshift.MasterProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

	vmSize := config["vm_size"].(string)
	subnetId := config["subnet_id"].(string)

	return &redhatopenshift.MasterProfile{
		VMSize:   redhatopenshift.VMSize(vmSize),
		SubnetID: utils.String(subnetId),
	}
}

func expandOpenshiftWorkerProfiles(inputs []interface{}) *[]redhatopenshift.WorkerProfile {
	if len(inputs) == 0 {
		return nil
	}

	profiles := make([]redhatopenshift.WorkerProfile, 0)
	config := inputs[0].(map[string]interface{})

	// Hardcoded name required by ARO interface
	workerName := "worker"

	vmSize := config["vm_size"].(string)
	if vmSize == "" {
		vmSize = "Standard_D4s_v3"
	}

	diskSizeGb := int32(config["disk_size_gb"].(int))
	if diskSizeGb == 0 {
		diskSizeGb = 128
	}

	nodeCount := int32(config["node_count"].(int))
	if nodeCount == 0 {
		nodeCount = 3
	}

	subnetId := config["subnet_id"].(string)

	profile := redhatopenshift.WorkerProfile{
		Name:       utils.String(workerName),
		VMSize:     redhatopenshift.VMSize1(vmSize),
		DiskSizeGB: utils.Int32(diskSizeGb),
		SubnetID:   utils.String(subnetId),
		Count:      utils.Int32(nodeCount),
	}

	profiles = append(profiles, profile)

	return &profiles
}

func expandOpenshiftApiServerProfile(input []interface{}) *redhatopenshift.APIServerProfile {
	if len(input) == 0 {
		return &redhatopenshift.APIServerProfile{
			Visibility: redhatopenshift.Public,
		}
	}
	config := input[0].(map[string]interface{})

	visibility := config["visibility"].(string)

	return &redhatopenshift.APIServerProfile{
		Visibility: redhatopenshift.Visibility(visibility),
	}
}

func expandOpenshiftIngressProfiles(inputs []interface{}) *[]redhatopenshift.IngressProfile {
	profiles := make([]redhatopenshift.IngressProfile, 0)

	if len(inputs) == 0 {
		profile := redhatopenshift.IngressProfile{
			Name:       utils.String("default"),
			Visibility: redhatopenshift.Visibility1Public,
		}

		profiles = append(profiles, profile)

		return &profiles
	}

	for index := range inputs {
		config := inputs[index].(map[string]interface{})

		name := config["name"].(string)
		visibility := config["visibility"].(string)

		profile := redhatopenshift.IngressProfile{
			Name:       utils.String(name),
			Visibility: redhatopenshift.Visibility1(visibility),
		}

		profiles = append(profiles, profile)
	}

	return &profiles
}
