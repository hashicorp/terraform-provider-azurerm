package redhatopenshift

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/redhatopenshift/mgmt/2022-04-01/redhatopenshift"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2022-09-04/openshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	openShiftValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/redhatopenshift/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var randomDomainName = GenerateRandomDomainName()

func resourceOpenShiftCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceOpenShiftClusterCreate,
		Read:   resourceOpenShiftClusterRead,
		Update: resourceOpenShiftClusterUpdate,
		Delete: resourceOpenShiftClusterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := openshiftclusters.ParseProviderOpenShiftClusterID(id)
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

			// todo swap this
			"location": azure.SchemaLocation(),

			// todo swap this
			"resource_group_name": azure.SchemaResourceGroupName(),

			// todo check that this needs to computed
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
						// todo check that this is forcenew
						"domain": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						// todo default really false?
						"fips_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
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
						// todo see if this is returned
						"client_secret": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			// todo is this computed? maybe we can just pull each attribute out
			"network_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// todo check default
						"pod_cidr": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "10.128.0.0/14",
							ValidateFunc: validate.CIDR,
						},
						// todo check default
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

			// todo check required
			"main_profile": {
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
						// todo  validate list of available sizes?
						"vm_size": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.StringIsNotEmpty,
						},
						// todo check default
						"encryption_at_host_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"disk_encryption_set_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			// todo check required
			"worker_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// todo validate sizes?
						"vm_size": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.StringIsNotEmpty,
						},
						"disk_size_gb": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: openShiftValidate.DiskSizeGB,
						},
						"node_count": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						// todo check default
						"encryption_at_host_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"disk_encryption_set_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			// todo pull out attribute
			"api_server_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"visibility": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(redhatopenshift.VisibilityPublic),
						},
					},
				},
			},

			// todo pull out attribute
			"ingress_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"visibility": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(redhatopenshift.VisibilityPublic),
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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceOpenShiftClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedHatOpenshift.OpenShiftClustersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Red Hat Openshift Cluster create.")

	id := openshiftclusters.NewProviderOpenShiftClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id.ID(), err)
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_redhat_openshift_cluster", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	// clusterProfile := expandOpenshiftClusterProfile(d.Get("cluster_profile").([]interface{}), resourceParse.NewResourceGroupID(id.SubscriptionId, id.ResourceGroupName).ID())
	servicePrincipalProfile := expandOpenshiftServicePrincipalProfile(d.Get("service_principal").([]interface{}))
	networkProfile := expandOpenshiftNetworkProfile(d.Get("network_profile").([]interface{}))
	mainProfile := expandOpenshiftMasterProfile(d.Get("main_profile").([]interface{}))
	workerProfiles := expandOpenshiftWorkerProfiles(d.Get("worker_profile").([]interface{}))
	// apiServerProfile := expandOpenshiftApiServerProfile(d.Get("api_server_profile").([]interface{}))

	// ingressProfilesRaw := d.Get("ingress_profile").([]interface{})
	// ingressProfiles := expandOpenshiftIngressProfiles(ingressProfilesRaw)

	t := d.Get("tags").(map[string]interface{})

	parameters := openshiftclusters.OpenShiftCluster{
		Name:     &id.OpenShiftClusterName,
		Location: location,
		Properties: &openshiftclusters.OpenShiftClusterProperties{
			// ClusterProfile:          clusterProfile,
			ServicePrincipalProfile: servicePrincipalProfile,
			NetworkProfile:          networkProfile,
			MasterProfile:           mainProfile,
			WorkerProfiles:          workerProfiles,
			// ApiserverProfile:        apiServerProfile,
			// IngressProfiles: ingressProfiles,
		},
		Tags: tags.Expand(t),
	}

	if err = client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id.ID(), err)
	}

	d.SetId(id.ID())

	return resourceOpenShiftClusterRead(d, meta)
}

func resourceOpenShiftClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedHatOpenshift.OpenShiftClustersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Red Hat OpenShift Cluster update.")

	id, err := openshiftclusters.ParseProviderOpenShiftClusterID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", id.ID(), err)
	}

	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `model` was nil", id.ID())
	}

	/*
		if d.HasChange("cluster_profile") {
			clusterProfileRaw := d.Get("cluster_profile").([]interface{})
			clusterProfile := expandOpenshiftClusterProfile(clusterProfileRaw, resourceGroupId)
			existing.Model.Properties.ClusterProfile = clusterProfile
		}*/

	if d.HasChange("main_profile") {
		mainProfileRaw := d.Get("main_profile").([]interface{})
		mainProfile := expandOpenshiftMasterProfile(mainProfileRaw)
		existing.Model.Properties.MasterProfile = mainProfile
	}

	if d.HasChange("worker_profile") {
		workerProfilesRaw := d.Get("worker_profile").([]interface{})
		workerProfiles := expandOpenshiftWorkerProfiles(workerProfilesRaw)
		existing.Model.Properties.WorkerProfiles = workerProfiles
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id.ID(), err)
	}

	return resourceOpenShiftClusterRead(d, meta)
}

func resourceOpenShiftClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedHatOpenshift.OpenShiftClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := openshiftclusters.ParseProviderOpenShiftClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id.ID())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}

	d.Set("name", id.OpenShiftClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
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

			mainProfile := flattenOpenShiftMasterProfile(props.MasterProfile)
			if err := d.Set("main_profile", mainProfile); err != nil {
				return fmt.Errorf("setting `main_profile`: %+v", err)
			}

			workerProfiles := flattenOpenShiftWorkerProfiles(props.WorkerProfiles)
			if err := d.Set("worker_profile", workerProfiles); err != nil {
				return fmt.Errorf("setting `worker_profile`: %+v", err)
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
			d.Set("console_url", props.ConsoleProfile.Url)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceOpenShiftClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedHatOpenshift.OpenShiftClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := openshiftclusters.ParseProviderOpenShiftClusterID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id.ID(), err)
	}

	return nil
}

func flattenOpenShiftClusterProfile(profile *openshiftclusters.ClusterProfile) []interface{} {
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

	fipsEnabled := false
	if profile.FipsValidatedModules != nil {
		fipsEnabled = *profile.FipsValidatedModules == openshiftclusters.FipsValidatedModulesEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"pull_secret":  pullSecret,
			"domain":       clusterDomain,
			"fips_enabled": fipsEnabled,
		},
	}
}

func flattenOpenShiftServicePrincipalProfile(profile *openshiftclusters.ServicePrincipalProfile, d *pluginsdk.ResourceData) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	clientID := ""
	if profile.ClientId != nil {
		clientID = *profile.ClientId
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

func flattenOpenShiftNetworkProfile(profile *openshiftclusters.NetworkProfile) []interface{} {
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

func flattenOpenShiftMasterProfile(profile *openshiftclusters.MasterProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	subnetId := ""
	if profile.SubnetId != nil {
		subnetId = *profile.SubnetId
	}

	encryptionAtHostEnabled := false
	if profile.EncryptionAtHost != nil {
		encryptionAtHostEnabled = *profile.EncryptionAtHost == openshiftclusters.EncryptionAtHostEnabled
	}

	diskEncryptionSetId := ""

	if profile.DiskEncryptionSetId != nil {
		diskEncryptionSetId = *profile.DiskEncryptionSetId
	}

	return []interface{}{
		map[string]interface{}{
			"vm_size":                    profile.VMSize,
			"subnet_id":                  subnetId,
			"encryption_at_host_enabled": encryptionAtHostEnabled,
			"disk_encryption_set_id":     diskEncryptionSetId,
		},
	}
}

func flattenOpenShiftWorkerProfiles(profiles *[]openshiftclusters.WorkerProfile) []interface{} {
	if profiles == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	result := make(map[string]interface{})
	result["node_count"] = int32(len(*profiles))

	// todo what is going on here? Shouldn't we have only one worker?
	for _, profile := range *profiles {
		if result["disk_size_gb"] == nil && profile.DiskSizeGB != nil {
			result["disk_size_gb"] = profile.DiskSizeGB
		}

		if result["vm_size"] == nil && profile.VMSize != nil {
			result["vm_size"] = profile.VMSize
		}

		if result["subnet_id"] == nil && profile.SubnetId != nil {
			result["subnet_id"] = profile.SubnetId
		}

		if result["encryption_at_host_enabled"] == nil {
			encryptionAtHostEnabled := false
			if profile.EncryptionAtHost != nil {
				encryptionAtHostEnabled = *profile.EncryptionAtHost == openshiftclusters.EncryptionAtHostEnabled
			}
			result["encryption_at_host_enabled"] = encryptionAtHostEnabled
		}

		if result["disk_encryption_set_id"] == nil && profile.DiskEncryptionSetId != nil {
			result["disk_encryption_set_id"] = profile.DiskEncryptionSetId
		}
	}

	results = append(results, result)

	return results
}

func flattenOpenShiftAPIServerProfile(profile *openshiftclusters.APIServerProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	visibility := ""
	if profile.Visibility != nil {
		visibility = string(*profile.Visibility)
	}

	return []interface{}{
		map[string]interface{}{
			"visibility": visibility,
		},
	}
}

func flattenOpenShiftIngressProfiles(profiles *[]openshiftclusters.IngressProfile) []interface{} {
	if profiles == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	for _, profile := range *profiles {
		visibility := ""
		if profile.Visibility != nil {
			visibility = string(*profile.Visibility)
		}

		result := make(map[string]interface{})
		result["visibility"] = visibility

		results = append(results, result)
	}

	return results
}

/*
func expandOpenshiftClusterProfile(input []interface{}, resourceGroupId string) *openshiftclusters.ClusterProfile {
	if len(input) == 0 {
		return &openshiftclusters.ClusterProfile{
			ResourceGroupId:      utils.String(resourceGroupId),
			Domain:               utils.String(randomDomainName),
			FipsValidatedModules: openshiftclusters.FipsValidatedModulesDisabled,
		}
	}

	config := input[0].(map[string]interface{})

	pullSecret := config["pull_secret"].(string)

	domain := config["domain"].(string)
	if domain == "" {
		domain = randomDomainName
	}

	fipsValidatedModules := openshiftclusters.FipsValidatedModulesDisabled
	fipsEnabled := config["fips_enabled"].(bool)
	if fipsEnabled {
		fipsValidatedModules = openshiftclusters.FipsValidatedModulesEnabled
	}

	return &openshiftclusters.ClusterProfile{
		ResourceGroupId:      utils.String(resourceGroupId),
		Domain:               utils.String(domain),
		PullSecret:           utils.String(pullSecret),
		FipsValidatedModules: fipsValidatedModules,
	}
}*/

func expandOpenshiftServicePrincipalProfile(input []interface{}) *openshiftclusters.ServicePrincipalProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})
	return &openshiftclusters.ServicePrincipalProfile{
		ClientId:     utils.String(config["client_id"].(string)),
		ClientSecret: utils.String(config["client_secret"].(string)),
	}
}

func expandOpenshiftNetworkProfile(input []interface{}) *openshiftclusters.NetworkProfile {
	if len(input) == 0 {
		return &openshiftclusters.NetworkProfile{
			// todo see if we need this
			PodCidr:     utils.String("10.128.0.0/14"),
			ServiceCidr: utils.String("172.30.0.0/16"),
		}
	}

	config := input[0].(map[string]interface{})

	podCidr := config["pod_cidr"].(string)
	serviceCidr := config["service_cidr"].(string)

	return &openshiftclusters.NetworkProfile{
		PodCidr:     utils.String(podCidr),
		ServiceCidr: utils.String(serviceCidr),
	}
}

func expandOpenshiftMasterProfile(input []interface{}) *openshiftclusters.MasterProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})
	encryptionAtHost := openshiftclusters.EncryptionAtHostDisabled
	enableEncryptionAtHost := config["encryption_at_host_enabled"].(bool)
	if enableEncryptionAtHost {
		encryptionAtHost = openshiftclusters.EncryptionAtHostEnabled
	}

	return &openshiftclusters.MasterProfile{
		VMSize:              utils.String(config["vm_size"].(string)),
		SubnetId:            utils.String(config["subnet_id"].(string)),
		EncryptionAtHost:    &encryptionAtHost,
		DiskEncryptionSetId: utils.String(config["disk_encryption_set_id"].(string)),
	}
}

func expandOpenshiftWorkerProfiles(inputs []interface{}) *[]openshiftclusters.WorkerProfile {
	if len(inputs) == 0 {
		return nil
	}

	profiles := make([]openshiftclusters.WorkerProfile, 0)
	config := inputs[0].(map[string]interface{})

	encryptionAtHost := openshiftclusters.EncryptionAtHostDisabled
	enableEncryptionAtHost := config["encryption_at_host_enabled"].(bool)
	if enableEncryptionAtHost {
		encryptionAtHost = openshiftclusters.EncryptionAtHostEnabled
	}

	profile := openshiftclusters.WorkerProfile{
		Name:                utils.String("worker"),
		VMSize:              utils.String(config["vm_size"].(string)),
		DiskSizeGB:          utils.Int64(int64(config["disk_size_gb"].(int))),
		SubnetId:            utils.String(config["subnet_id"].(string)),
		Count:               utils.Int64(int64(config["node_count"].(int))),
		EncryptionAtHost:    &encryptionAtHost,
		DiskEncryptionSetId: utils.String(config["disk_encryption_set_id"].(string)),
	}

	profiles = append(profiles, profile)

	return &profiles
}

/* todo figure this out and add more attributes
func expandOpenshiftApiServerProfile(input []interface{}) *openshiftclusters.APIServerProfile {
	if len(input) == 0 {
		return &openshiftclusters.APIServerProfile{
			Visibility: redhatopenshift.VisibilityPublic,
		}
	}

	config := input[0].(map[string]interface{})

	visibility := config["visibility"].(string)

	return &openshiftclusters.APIServerProfile{
		Visibility: redhatopenshift.Visibility(visibility),
	}
}*/

func expandOpenshiftIngressProfiles(inputs []interface{}) *[]redhatopenshift.IngressProfile {
	profiles := make([]redhatopenshift.IngressProfile, 0)

	name := "default"
	visibility := string(redhatopenshift.VisibilityPublic)

	if len(inputs) > 0 {
		input := inputs[0].(map[string]interface{})
		visibility = input["visibility"].(string)
	}

	profile := redhatopenshift.IngressProfile{
		Name:       utils.String(name),
		Visibility: redhatopenshift.Visibility(visibility),
	}

	profiles = append(profiles, profile)

	return &profiles
}
