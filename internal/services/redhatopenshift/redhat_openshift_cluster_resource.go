package redhatopenshift

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2022-09-04/openshiftclusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
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
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"cluster_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"domain": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"fips_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"pull_secret": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"resource_group_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"service_principal": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"client_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: openShiftValidate.ClientID,
						},
						"client_secret": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ForceNew:     true,
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

			"main_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"vm_size": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.StringIsNotEmpty,
						},
						"encryption_at_host_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"disk_encryption_set_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"worker_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"vm_size": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.StringIsNotEmpty,
						},
						"disk_size_gb": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: openShiftValidate.DiskSizeGB,
						},
						"node_count": {
							Type:     pluginsdk.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"encryption_at_host_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"disk_encryption_set_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"api_server_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"visibility": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(openshiftclusters.PossibleValuesForVisibility(), false),
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ingress_profile": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"visibility": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(openshiftclusters.PossibleValuesForVisibility(), false),
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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

			"tags": commonschema.TagsForceNew(),
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
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_redhat_openshift_cluster", id.ID())
	}

	parameters := openshiftclusters.OpenShiftCluster{
		Name:     &id.OpenShiftClusterName,
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Properties: &openshiftclusters.OpenShiftClusterProperties{
			ClusterProfile:          expandOpenshiftClusterProfile(d.Get("cluster_profile").([]interface{}), id.SubscriptionId),
			ServicePrincipalProfile: expandOpenshiftServicePrincipalProfile(d.Get("service_principal").([]interface{})),
			NetworkProfile:          expandOpenshiftNetworkProfile(d.Get("network_profile").([]interface{})),
			MasterProfile:           expandOpenshiftMasterProfile(d.Get("main_profile").([]interface{})),
			WorkerProfiles:          expandOpenshiftWorkerProfiles(d.Get("worker_profile").([]interface{})),
			ApiserverProfile:        expandOpenshiftApiServerProfile(d.Get("api_server_profile").([]interface{})),
			IngressProfiles:         expandOpenshiftIngressProfiles(d.Get("ingress_profile").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err = client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id.ID(), err)
	}

	d.SetId(id.ID())

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

	resourceGroupId := ""
	if profile.ResourceGroupId != nil {
		resourceGroupId = *profile.ResourceGroupId
	}

	clusterDomain := ""
	if profile.Domain != nil {
		clusterDomain = *profile.Domain
	}

	pullSecret := ""
	if profile.PullSecret != nil {
		pullSecret = *profile.PullSecret
	}

	fipsEnabled := false
	if profile.FipsValidatedModules != nil {
		fipsEnabled = *profile.FipsValidatedModules == openshiftclusters.FipsValidatedModulesEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"pull_secret":       pullSecret,
			"domain":            clusterDomain,
			"fips_enabled":      fipsEnabled,
			"resource_group_id": resourceGroupId,
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
	if profiles == nil || len(*profiles) == 0 {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	result := make(map[string]interface{})
	result["node_count"] = int32(len(*profiles))

	rawProfiles := *profiles
	profile := rawProfiles[0]

	nodeCount := int64(0)
	if profile.Count != nil {
		nodeCount = *profile.Count
	}

	diskSizeGB := int64(0)
	if profile.DiskSizeGB != nil {
		diskSizeGB = *profile.DiskSizeGB
	}

	vmSize := ""
	if profile.VMSize != nil {
		vmSize = *profile.VMSize
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

	results = append(results, result)

	return []interface{}{
		map[string]interface{}{
			"node_count":                 nodeCount,
			"vm_size":                    vmSize,
			"disk_size_gb":               diskSizeGB,
			"subnet_id":                  subnetId,
			"encryption_at_host_enabled": encryptionAtHostEnabled,
			"disk_encryption_set_id":     diskEncryptionSetId,
		},
	}
}

func flattenOpenShiftAPIServerProfile(profile *openshiftclusters.APIServerProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	visibility := ""
	if profile.Visibility != nil {
		visibility = string(*profile.Visibility)
	}

	url := ""
	if profile.Url != nil {
		url = *profile.Url
	}

	ipAddress := ""
	if profile.IP != nil {
		ipAddress = *profile.IP
	}

	return []interface{}{
		map[string]interface{}{
			"visibility": visibility,
			"url":        url,
			"ip_address": ipAddress,
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

		name := ""
		if profile.Name != nil {
			name = *profile.Name
		}

		ipAddress := ""
		if profile.IP != nil {
			ipAddress = *profile.IP
		}

		result := make(map[string]interface{})
		result["visibility"] = visibility
		result["name"] = name
		result["ip_address"] = ipAddress

		results = append(results, result)
	}

	return results
}

func expandOpenshiftClusterProfile(input []interface{}, subscriptionId string) *openshiftclusters.ClusterProfile {
	config := input[0].(map[string]interface{})

	fipsValidatedModules := openshiftclusters.FipsValidatedModulesDisabled
	if config["fips_enabled"].(bool) {
		fipsValidatedModules = openshiftclusters.FipsValidatedModulesEnabled
	}

	return &openshiftclusters.ClusterProfile{
		// the api needs a ResourceGroupId value and the portal doesn't allow you to set it but the portal returns the
		// resource id being `aro-{domain}` so we'll follow that here.
		ResourceGroupId:      utils.String(commonids.NewResourceGroupID(subscriptionId, fmt.Sprintf("aro-%s", config["domain"].(string))).ID()),
		Domain:               utils.String(config["domain"].(string)),
		PullSecret:           utils.String(config["pull_secret"].(string)),
		FipsValidatedModules: &fipsValidatedModules,
	}
}

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

func expandOpenshiftApiServerProfile(input []interface{}) *openshiftclusters.APIServerProfile {
	config := input[0].(map[string]interface{})

	visibility := openshiftclusters.Visibility(config["visibility"].(string))

	return &openshiftclusters.APIServerProfile{
		Visibility: &visibility,
	}
}

func expandOpenshiftIngressProfiles(input []interface{}) *[]openshiftclusters.IngressProfile {
	config := input[0].(map[string]interface{})

	profiles := make([]openshiftclusters.IngressProfile, 0)

	visibility := openshiftclusters.Visibility(config["visibility"].(string))

	profile := openshiftclusters.IngressProfile{
		Name:       utils.String("default"),
		Visibility: &visibility,
	}

	profiles = append(profiles, profile)

	return &profiles
}
