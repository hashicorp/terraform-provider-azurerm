package redhatopenshift

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2023-09-04/openshiftclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	openShiftValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/redhatopenshift/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

var randomDomainName = GenerateRandomDomainName()

var _ sdk.ResourceWithUpdate = RedHatOpenShiftCluster{}

type RedHatOpenShiftCluster struct{}

type RedHatOpenShiftClusterModel struct {
	Name             string             `tfschema:"name"`
	Location         string             `tfschema:"location"`
	ResourceGroup    string             `tfschema:"resource_group_name"`
	Version          string             `tfschema:"version"`
	ConsoleUrl       string             `tfschema:"console_url"`
	ServicePrincipal []ServicePrincipal `tfschema:"service_principal"`
	ClusterProfile   []ClusterProfile   `tfschema:"cluster_profile"`
	NetworkProfile   []NetworkProfile   `tfschema:"network_profile"`
	MainProfile      []MainProfile      `tfschema:"main_profile"`
	WorkerProfile    []WorkerProfile    `tfschema:"worker_profile"`
	ApiServerProfile []ApiServerProfile `tfschema:"api_server_profile"`
	IngressProfile   []IngressProfile   `tfschema:"ingress_profile"`
	Tags             map[string]string  `tfschema:"tags"`
}

type ServicePrincipal struct {
	ClientId     string `tfschema:"client_id"`
	ClientSecret string `tfschema:"client_secret"`
}

type ClusterProfile struct {
	PullSecret  string `tfschema:"pull_secret"`
	Domain      string `tfschema:"domain"`
	FipsEnabled bool   `tfschema:"fips_enabled"`
}

type NetworkProfile struct {
	PodCidr     string `tfschema:"pod_cidr"`
	ServiceCidr string `tfschema:"service_cidr"`
}

type MainProfile struct {
	SubnetId                string `tfschema:"subnet_id"`
	VmSize                  string `tfschema:"vm_size"`
	EncryptionAtHostEnabled bool   `tfschema:"encryption_at_host_enabled"`
	DiskEncryptionSetId     string `tfschema:"disk_encryption_set_id"`
}

type WorkerProfile struct {
	VmSize                  string `tfschema:"vm_size"`
	DiskSizeGb              int64  `tfschema:"disk_size_gb"`
	NodeCount               int64  `tfschema:"node_count"`
	SubnetId                string `tfschema:"subnet_id"`
	EncryptionAtHostEnabled bool   `tfschema:"encryption_at_host_enabled"`
	DiskEncryptionSetId     string `tfschema:"disk_encryption_set_id"`
}

type IngressProfile struct {
	Visibility string `tfschema:"visibility"`
	Ip         string `tfschema:"ip"`
}

type ApiServerProfile struct {
	Visibility string `tfschema:"visibility"`
	Url        string `tfschema:"url"`
}

func (r RedHatOpenShiftCluster) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.IntBetween(3, 60),
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

		"tags": commonschema.Tags(),
	}
}

func (r RedHatOpenShiftCluster) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"console_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r RedHatOpenShiftCluster) ModelObject() interface{} {
	return &RedHatOpenShiftClusterModel{}
}

func (r RedHatOpenShiftCluster) ResourceType() string {
	return "azurerm_redhat_openshift_cluster"
}

func (r RedHatOpenShiftCluster) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return openshiftclusters.ValidateProviderOpenShiftClusterID
}

func (r RedHatOpenShiftCluster) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedHatOpenshift.OpenShiftClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config RedHatOpenShiftClusterModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := openshiftclusters.NewProviderOpenShiftClusterID(subscriptionId, config.ResourceGroup, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id.ID(), err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			parameters := openshiftclusters.OpenShiftCluster{
				Name:     &id.OpenShiftClusterName,
				Location: azure.NormalizeLocation(config.Location),
				Properties: &openshiftclusters.OpenShiftClusterProperties{
					ClusterProfile:          expandOpenshiftClusterProfile(config.ClusterProfile, id.SubscriptionId),
					ServicePrincipalProfile: expandOpenshiftServicePrincipalProfile(config.ServicePrincipal),
					NetworkProfile:          expandOpenshiftNetworkProfile(config.NetworkProfile),
					MasterProfile:           expandOpenshiftMasterProfile(config.MainProfile),
					WorkerProfiles:          expandOpenshiftWorkerProfiles(config.WorkerProfile),
					ApiserverProfile:        expandOpenshiftApiServerProfile(config.ApiServerProfile),
					IngressProfiles:         expandOpenshiftIngressProfiles(config.IngressProfile),
				},
				Tags: pointer.To(config.Tags),
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id.ID(), err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r RedHatOpenShiftCluster) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedHatOpenshift.OpenShiftClustersClient

			id, err := openshiftclusters.ParseProviderOpenShiftClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] %s was not found - removing from state!", id.ID())
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id.ID(), err)
			}

			model := RedHatOpenShiftClusterModel{
				Name:          id.OpenShiftClusterName,
				ResourceGroup: id.ResourceGroup,
			}

			if model := resp.Model; model != nil {
				model.Location = location.Normalize(model.Location)
				model.Tags = pointer.From(model.Tags)

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
		},
	}
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

func expandOpenshiftClusterProfile(input []ClusterProfile, subscriptionId string) *openshiftclusters.ClusterProfile {
	if len(input) == 0 {
		return nil
	}

	fipsValidatedModules := openshiftclusters.FipsValidatedModulesDisabled
	if input[0].FipsEnabled {
		fipsValidatedModules = openshiftclusters.FipsValidatedModulesEnabled
	}

	return &openshiftclusters.ClusterProfile{
		// the api needs a ResourceGroupId value and the portal doesn't allow you to set it but the portal returns the
		// resource id being `aro-{domain}` so we'll follow that here.
		ResourceGroupId:      pointer.To(commonids.NewResourceGroupID(subscriptionId, fmt.Sprintf("aro-%s", input[0].Domain)).ID()),
		Domain:               pointer.To(input[0].Domain),
		PullSecret:           pointer.To(input[0].PullSecret),
		FipsValidatedModules: pointer.To(fipsValidatedModules),
	}
}

func expandOpenshiftServicePrincipalProfile(input []ServicePrincipal) *openshiftclusters.ServicePrincipalProfile {
	if len(input) == 0 {
		return nil
	}

	return &openshiftclusters.ServicePrincipalProfile{
		ClientId:     pointer.To(input[0].ClientId),
		ClientSecret: pointer.To(input[0].ClientSecret),
	}
}

func expandOpenshiftNetworkProfile(input []NetworkProfile) *openshiftclusters.NetworkProfile {
	if len(input) == 0 {
		return &openshiftclusters.NetworkProfile{
			PodCidr:     pointer.To("10.128.0.0/14"),
			ServiceCidr: pointer.To("172.30.0.0/16"),
		}
	}

	return &openshiftclusters.NetworkProfile{
		PodCidr:     pointer.To(input[0].PodCidr),
		ServiceCidr: pointer.To(input[0].ServiceCidr),
	}
}

func expandOpenshiftMasterProfile(input []MainProfile) *openshiftclusters.MasterProfile {
	if len(input) == 0 {
		return nil
	}

	encryptionAtHost := openshiftclusters.EncryptionAtHostDisabled
	if input[0].EncryptionAtHostEnabled {
		encryptionAtHost = openshiftclusters.EncryptionAtHostEnabled
	}

	return &openshiftclusters.MasterProfile{
		VMSize:              pointer.To(input[0].VmSize),
		SubnetId:            pointer.To(input[0].SubnetId),
		EncryptionAtHost:    pointer.To(encryptionAtHost),
		DiskEncryptionSetId: pointer.To(input[0].DiskEncryptionSetId),
	}
}

func expandOpenshiftWorkerProfiles(input []WorkerProfile) *[]openshiftclusters.WorkerProfile {
	if len(input) == 0 {
		return nil
	}

	profiles := make([]openshiftclusters.WorkerProfile, 0)

	encryptionAtHost := openshiftclusters.EncryptionAtHostDisabled
	if input[0].EncryptionAtHostEnabled {
		encryptionAtHost = openshiftclusters.EncryptionAtHostEnabled
	}

	profile := openshiftclusters.WorkerProfile{
		Name:                pointer.To("worker"),
		VMSize:              pointer.To(input[0].VmSize),
		DiskSizeGB:          pointer.To(input[0].DiskSizeGb),
		SubnetId:            pointer.To(input[0].SubnetId),
		Count:               pointer.To(input[0].NodeCount),
		EncryptionAtHost:    pointer.To(encryptionAtHost),
		DiskEncryptionSetId: pointer.To(input[0].DiskEncryptionSetId),
	}

	profiles = append(profiles, profile)

	return &profiles
}

func expandOpenshiftApiServerProfile(input []ApiServerProfile) *openshiftclusters.APIServerProfile {
	if len(input) == 0 {
		return nil
	}

	visibility := openshiftclusters.Visibility(input[0].Visibility)

	return &openshiftclusters.APIServerProfile{
		Visibility: &visibility,
	}
}

func expandOpenshiftIngressProfiles(input []IngressProfile) *[]openshiftclusters.IngressProfile {
	if len(input) == 0 {
		return nil
	}

	profiles := make([]openshiftclusters.IngressProfile, 0)

	profile := openshiftclusters.IngressProfile{
		Name:       pointer.To("default"),
		Visibility: pointer.To(openshiftclusters.Visibility(input[0].Visibility)),
	}

	profiles = append(profiles, profile)

	return &profiles
}
