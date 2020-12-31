package springcloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2019-05-01-preview/appplatform"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSpringCloudService() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpringCloudServiceCreate,
		Read:   resourceSpringCloudServiceRead,
		Update: resourceSpringCloudServiceUpdate,
		Delete: resourceSpringCloudServiceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SpringCloudServiceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudServiceName,
			},

			// Spring Cloud Service only supports following locations, we are still supporting more locations (Wednesday, November 20, 2019 4:20 PM):
			// `East US`, `Southeast Asia`, `West Europe`, `West US 2`
			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "S0",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"B0",
					"S0",
				}, false),
			},

			"network": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: networkValidate.SubnetID,
						},

						"service_runtime_subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: networkValidate.SubnetID,
						},

						"cidr_ranges": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MinItems: 3,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"app_network_resource_group": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"service_runtime_network_resource_group": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},

			"config_server_git_setting": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uri": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.ConfigServerURI,
						},

						"label": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"search_paths": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"http_basic_auth": SchemaConfigServerHttpBasicAuth("config_server_git_setting.0.ssh_auth"),

						"ssh_auth": SchemaConfigServerSSHAuth("config_server_git_setting.0.http_basic_auth"),

						"repository": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"uri": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.ConfigServerURI,
									},

									"label": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"pattern": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},

									"search_paths": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},

									"http_basic_auth": SchemaConfigServerHttpBasicAuth(),

									"ssh_auth": SchemaConfigServerSSHAuth(),
								},
							},
						},
					},
				},
			},

			"trace": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instrumentation_key": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"outbound_public_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSpringCloudServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resourceId := parse.NewSpringCloudServiceID(subscriptionId, resourceGroup, name).ID()
	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Spring Cloud %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_spring_cloud_service", resourceId)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	resource := appplatform.ServiceResource{
		Location: utils.String(location),
		Properties: &appplatform.ClusterResourceProperties{
			Trace:          expandSpringCloudTrace(d.Get("trace").([]interface{})),
			NetworkProfile: expandSpringCloudNetwork(d.Get("network").([]interface{})),
		},
		Sku: &appplatform.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	gitProperty, err := expandSpringCloudConfigServerGitProperty(d.Get("config_server_git_setting").([]interface{}))
	if err != nil {
		return err
	}

	// current create api doesn't take care parameters of config server.
	// so we need to invoke create api first and then update api
	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, resource)
	if err != nil {
		return fmt.Errorf("creating Spring Cloud %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Spring Cloud %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if gitProperty != nil {
		resource.Properties = &appplatform.ClusterResourceProperties{
			ConfigServerProperties: &appplatform.ConfigServerProperties{
				ConfigServer: &appplatform.ConfigServerSettings{
					GitProperty: gitProperty,
				},
			},
		}

		updateFuture, err := client.Update(ctx, resourceGroup, name, resource)
		if err != nil {
			return fmt.Errorf("failure updating config server of Spring Cloud Service %q  (Resource Group %q): %+v", name, resourceGroup, err)
		}
		if err = updateFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("failure waiting for setting config server of Spring Cloud Service %q config server (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("unable to retrieve Spring Cloud Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.Properties != nil && resp.Properties.ConfigServerProperties != nil {
		if err := resp.Properties.ConfigServerProperties.Error; err != nil {
			return fmt.Errorf("failure setting config server of Spring Cloud Service %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	d.SetId(resourceId)

	return resourceSpringCloudServiceRead(d, meta)
}

func resourceSpringCloudServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	gitPropertyRaw := d.Get("config_server_git_setting").([]interface{})
	gitProperty, err := expandSpringCloudConfigServerGitProperty(gitPropertyRaw)
	if err != nil {
		return err
	}

	springCloudService := appplatform.ServiceResource{
		Properties: &appplatform.ClusterResourceProperties{
			ConfigServerProperties: &appplatform.ConfigServerProperties{
				ConfigServer: &appplatform.ConfigServerSettings{
					GitProperty: gitProperty,
				},
			},
			Trace: expandSpringCloudTrace(d.Get("trace").([]interface{})),
		},
		Sku: &appplatform.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Update(ctx, resourceGroup, name, springCloudService)
	if err != nil {
		return fmt.Errorf("updating Spring Cloud Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Spring Cloud Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := future.Result(*client)
	if err != nil {
		return fmt.Errorf("failure getting result of Spring Cloud Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.Properties != nil && resp.Properties.ConfigServerProperties != nil {
		if err := resp.Properties.ConfigServerProperties.Error; err != nil {
			return fmt.Errorf("failure setting config server of Spring Cloud Service %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return resourceSpringCloudServiceRead(d, meta)
}

func resourceSpringCloudServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud Service %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to read Spring Cloud Service %q (Resource Group %q): %+v", id.SpringName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if resp.Sku != nil {
		d.Set("sku_name", resp.Sku.Name)
	}
	if props := resp.Properties; props != nil {
		if err := d.Set("trace", flattenSpringCloudTrace(props.Trace)); err != nil {
			return fmt.Errorf("failure setting `trace`: %+v", err)
		}
		if err := d.Set("network", flattenSpringCloudNetwork(props.NetworkProfile)); err != nil {
			return fmt.Errorf("setting `network`: %+v", err)
		}

		outboundPublicIPAddresses := flattenOutboundPublicIPAddresses(props.NetworkProfile)
		if err := d.Set("outbound_public_ip_addresses", outboundPublicIPAddresses); err != nil {
			return fmt.Errorf("setting `outbound_public_ip_addresses`: %+v", err)
		}
		if err := d.Set("config_server_git_setting", flattenSpringCloudConfigServerGitProperty(props.ConfigServerProperties, d)); err != nil {
			return fmt.Errorf("setting `config_server_git_setting`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSpringCloudServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudServiceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		return fmt.Errorf("failure deleting Spring Cloud Service %q (Resource Group %q): %+v", id.SpringName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("failure waiting for deleting Spring Cloud Service %q (Resource Group %q): %+v", id.SpringName, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandSpringCloudNetwork(input []interface{}) *appplatform.NetworkProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	cidrRanges := utils.ExpandStringSlice(v["cidr_ranges"].([]interface{}))
	network := &appplatform.NetworkProfile{
		ServiceRuntimeSubnetID: utils.String(v["service_runtime_subnet_id"].(string)),
		AppSubnetID:            utils.String(v["app_subnet_id"].(string)),
		ServiceCidr:            utils.String(strings.Join(*cidrRanges, ",")),
	}
	if serviceRuntimeNetworkResourceGroup := v["service_runtime_network_resource_group"].(string); serviceRuntimeNetworkResourceGroup != "" {
		network.ServiceRuntimeNetworkResourceGroup = utils.String(serviceRuntimeNetworkResourceGroup)
	}
	if appNetworkResourceGroup := v["app_network_resource_group"].(string); appNetworkResourceGroup != "" {
		network.AppNetworkResourceGroup = utils.String(appNetworkResourceGroup)
	}
	return network
}

func expandSpringCloudConfigServerGitProperty(input []interface{}) (*appplatform.ConfigServerGitProperty, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	v := input[0].(map[string]interface{})
	if v == nil {
		return nil, nil
	}

	result := appplatform.ConfigServerGitProperty{
		URI: utils.String(v["uri"].(string)),
	}

	if label := v["label"].(string); label != "" {
		result.Label = utils.String(label)
	}
	if searchPaths := v["search_paths"].([]interface{}); len(searchPaths) > 0 {
		result.SearchPaths = utils.ExpandStringSlice(searchPaths)
	}

	httpBasicAuth := v["http_basic_auth"].([]interface{})
	sshAuth := v["ssh_auth"].([]interface{})
	if len(httpBasicAuth) > 0 && len(sshAuth) > 0 {
		return nil, fmt.Errorf("can not set both `http_basic_auth` and `ssh_auth`")
	}
	if len(httpBasicAuth) > 0 {
		v := httpBasicAuth[0].(map[string]interface{})
		result.Username = utils.String(v["username"].(string))
		result.Password = utils.String(v["password"].(string))
	}
	if len(sshAuth) > 0 {
		v := sshAuth[0].(map[string]interface{})
		result.PrivateKey = utils.String(v["private_key"].(string))
		result.StrictHostKeyChecking = utils.Bool(v["strict_host_key_checking_enabled"].(bool))

		if hostKey := v["host_key"].(string); hostKey != "" {
			result.HostKey = utils.String(hostKey)
		}
		if hostKeyAlgorithm := v["host_key_algorithm"].(string); hostKeyAlgorithm != "" {
			result.HostKeyAlgorithm = utils.String(hostKeyAlgorithm)
		}
	}

	if v, ok := v["repository"]; ok {
		repositories, err := expandSpringCloudGitPatternRepository(v.([]interface{}))
		if err != nil {
			return nil, err
		}
		result.Repositories = repositories
	}

	return &result, nil
}

func expandSpringCloudGitPatternRepository(input []interface{}) (*[]appplatform.GitPatternRepository, error) {
	results := make([]appplatform.GitPatternRepository, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		result := appplatform.GitPatternRepository{
			Name: utils.String(v["name"].(string)),
			URI:  utils.String(v["uri"].(string)),
		}

		if label := v["label"].(string); len(label) > 0 {
			result.Label = utils.String(label)
		}
		if pattern := v["pattern"].([]interface{}); len(pattern) > 0 {
			result.Pattern = utils.ExpandStringSlice(pattern)
		}
		if searchPaths := v["search_paths"].([]interface{}); len(searchPaths) > 0 {
			result.SearchPaths = utils.ExpandStringSlice(searchPaths)
		}

		httpBasicAuth := v["http_basic_auth"].([]interface{})
		sshAuth := v["ssh_auth"].([]interface{})
		if len(httpBasicAuth) > 0 && len(sshAuth) > 0 {
			return nil, fmt.Errorf("can not set both `http_basic_auth` and `ssh_auth` for the same repository")
		}
		if len(httpBasicAuth) > 0 {
			v := httpBasicAuth[0].(map[string]interface{})
			result.Username = utils.String(v["username"].(string))
			result.Password = utils.String(v["password"].(string))
		}
		if len(sshAuth) > 0 {
			v := sshAuth[0].(map[string]interface{})
			result.PrivateKey = utils.String(v["private_key"].(string))
			result.StrictHostKeyChecking = utils.Bool(v["strict_host_key_checking_enabled"].(bool))

			if hostKey := v["host_key"].(string); hostKey != "" {
				result.HostKey = utils.String(hostKey)
			}
			if hostKeyAlgorithm := v["host_key_algorithm"].(string); hostKeyAlgorithm != "" {
				result.HostKeyAlgorithm = utils.String(hostKeyAlgorithm)
			}
		}

		results = append(results, result)
	}
	return &results, nil
}

func expandSpringCloudTrace(input []interface{}) *appplatform.TraceProperties {
	if len(input) == 0 || input[0] == nil {
		return &appplatform.TraceProperties{
			Enabled: utils.Bool(false),
		}
	}
	v := input[0].(map[string]interface{})
	return &appplatform.TraceProperties{
		Enabled:                      utils.Bool(true),
		AppInsightInstrumentationKey: utils.String(v["instrumentation_key"].(string)),
	}
}

func flattenSpringCloudConfigServerGitProperty(input *appplatform.ConfigServerProperties, d *schema.ResourceData) []interface{} {
	if input == nil || input.ConfigServer == nil || input.ConfigServer.GitProperty == nil {
		return []interface{}{}
	}

	gitProperty := input.ConfigServer.GitProperty

	// prepare old state to find sensitive props not returned by API.
	oldGitSetting := make(map[string]interface{})
	if oldGitSettings := d.Get("config_server_git_setting").([]interface{}); len(oldGitSettings) > 0 {
		oldGitSetting = oldGitSettings[0].(map[string]interface{})
	}

	uri := ""
	if gitProperty.URI != nil {
		uri = *gitProperty.URI
	}

	label := ""
	if gitProperty.Label != nil {
		label = *gitProperty.Label
	}

	searchPaths := utils.FlattenStringSlice(gitProperty.SearchPaths)

	httpBasicAuth := make([]interface{}, 0)
	if gitProperty.Username != nil && gitProperty.Password != nil {
		// username and password returned by API are *
		// to avoid state diff, we get the props from old state
		username := ""
		password := ""
		if v, ok := oldGitSetting["http_basic_auth"]; ok {
			oldHTTPBasicAuth := v.([]interface{})
			if len(oldHTTPBasicAuth) > 0 {
				oldItem := oldHTTPBasicAuth[0].(map[string]interface{})
				username = oldItem["username"].(string)
				password = oldItem["password"].(string)
			}
		}

		httpBasicAuth = []interface{}{
			map[string]interface{}{
				"username": username,
				"password": password,
			},
		}
	}

	sshAuth := []interface{}{}
	if gitProperty.PrivateKey != nil {
		// private_key, host_key and host_key_algorithm returned by API are *
		// to avoid state diff, we get the props from old state
		privateKey := ""
		hostKey := ""
		hostKeyAlgorithm := ""
		if v, ok := oldGitSetting["ssh_auth"]; ok {
			sshAuth := v.([]interface{})
			if len(sshAuth) > 0 {
				oldItem := sshAuth[0].(map[string]interface{})
				privateKey = oldItem["private_key"].(string)
				hostKey = oldItem["host_key"].(string)
				hostKeyAlgorithm = oldItem["host_key_algorithm"].(string)
			}
		}

		strictHostKeyChecking := false
		if gitProperty.StrictHostKeyChecking != nil {
			strictHostKeyChecking = *gitProperty.StrictHostKeyChecking
		}

		sshAuth = []interface{}{
			map[string]interface{}{
				"private_key":                      privateKey,
				"host_key":                         hostKey,
				"host_key_algorithm":               hostKeyAlgorithm,
				"strict_host_key_checking_enabled": strictHostKeyChecking,
			},
		}
	}

	return []interface{}{
		map[string]interface{}{
			"uri":             uri,
			"label":           label,
			"search_paths":    searchPaths,
			"http_basic_auth": httpBasicAuth,
			"ssh_auth":        sshAuth,
			"repository":      flattenSpringCloudGitPatternRepository(gitProperty.Repositories, d),
		},
	}
}

func flattenSpringCloudGitPatternRepository(input *[]appplatform.GitPatternRepository, d *schema.ResourceData) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	// prepare old state to find sensitive props not returned by API.
	oldGitPatternRepositories := []interface{}{}
	if oldGitSettings := d.Get("config_server_git_setting").([]interface{}); len(oldGitSettings) > 0 {
		oldGitSetting := oldGitSettings[0].(map[string]interface{})
		oldGitPatternRepositories = oldGitSetting["repository"].([]interface{})
	}

	for i, item := range *input {
		// prepare old state to find sensitive props not returned by API.
		oldGitPatternRepository := make(map[string]interface{})
		if len(oldGitPatternRepositories) > 0 {
			oldGitPatternRepository = oldGitPatternRepositories[i].(map[string]interface{})
		}

		name := ""
		if item.Name != nil {
			name = *item.Name
		}

		uri := ""
		if item.URI != nil {
			uri = *item.URI
		}

		label := ""
		if item.Label != nil {
			label = *item.Label
		}

		pattern := utils.FlattenStringSlice(item.Pattern)
		searchPaths := utils.FlattenStringSlice(item.SearchPaths)

		httpBasicAuth := []interface{}{}
		if item.Username != nil && item.Password != nil {
			// username and password returned by API are *
			// to avoid state diff, we get the props from old state
			username := ""
			password := ""
			if v, ok := oldGitPatternRepository["http_basic_auth"]; ok {
				oldHTTPBasicAuth := v.([]interface{})
				if len(oldHTTPBasicAuth) > 0 {
					oldItem := oldHTTPBasicAuth[0].(map[string]interface{})
					username = oldItem["username"].(string)
					password = oldItem["password"].(string)
				}
			}

			httpBasicAuth = []interface{}{
				map[string]interface{}{
					"username": username,
					"password": password,
				},
			}
		}

		sshAuth := []interface{}{}
		if item.PrivateKey != nil {
			// private_key, host_key and host_key_algorithm returned by API are *
			// to avoid state diff, we get the props from old state
			privateKey := ""
			hostKey := ""
			hostKeyAlgorithm := ""
			if v, ok := oldGitPatternRepository["ssh_auth"]; ok {
				sshAuth := v.([]interface{})
				if len(sshAuth) > 0 {
					oldItem := sshAuth[0].(map[string]interface{})
					privateKey = oldItem["private_key"].(string)
					hostKey = oldItem["host_key"].(string)
					hostKeyAlgorithm = oldItem["host_key_algorithm"].(string)
				}
			}

			strictHostKeyChecking := false
			if item.StrictHostKeyChecking != nil {
				strictHostKeyChecking = *item.StrictHostKeyChecking
			}

			sshAuth = []interface{}{
				map[string]interface{}{
					"private_key":                      privateKey,
					"host_key":                         hostKey,
					"host_key_algorithm":               hostKeyAlgorithm,
					"strict_host_key_checking_enabled": strictHostKeyChecking,
				},
			}
		}

		results = append(results, map[string]interface{}{
			"name":            name,
			"uri":             uri,
			"label":           label,
			"pattern":         pattern,
			"search_paths":    searchPaths,
			"http_basic_auth": httpBasicAuth,
			"ssh_auth":        sshAuth,
		})
	}

	return results
}

func flattenSpringCloudTrace(input *appplatform.TraceProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	enabled := false
	instrumentationKey := ""
	if input.Enabled != nil {
		enabled = *input.Enabled
	}
	if input.AppInsightInstrumentationKey != nil {
		instrumentationKey = *input.AppInsightInstrumentationKey
	}

	if !enabled {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"instrumentation_key": instrumentationKey,
		},
	}
}

func flattenSpringCloudNetwork(input *appplatform.NetworkProfile) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var serviceRuntimeSubnetID, appSubnetID, serviceRuntimeNetworkResourceGroup, appNetworkResourceGroup string
	var cidrRanges []interface{}
	if input.ServiceRuntimeSubnetID != nil {
		serviceRuntimeSubnetID = *input.ServiceRuntimeSubnetID
	}
	if input.AppSubnetID != nil {
		appSubnetID = *input.AppSubnetID
	}
	if input.ServiceCidr != nil {
		cidrs := strings.Split(*input.ServiceCidr, ",")
		cidrRanges = utils.FlattenStringSlice(&cidrs)
	}
	if input.ServiceRuntimeNetworkResourceGroup != nil {
		serviceRuntimeNetworkResourceGroup = *input.ServiceRuntimeNetworkResourceGroup
	}
	if input.AppNetworkResourceGroup != nil {
		appNetworkResourceGroup = *input.AppNetworkResourceGroup
	}

	if serviceRuntimeSubnetID == "" && appSubnetID == "" && serviceRuntimeNetworkResourceGroup == "" && appNetworkResourceGroup == "" && len(cidrRanges) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"app_subnet_id":                          appSubnetID,
			"service_runtime_subnet_id":              serviceRuntimeSubnetID,
			"cidr_ranges":                            cidrRanges,
			"app_network_resource_group":             appNetworkResourceGroup,
			"service_runtime_network_resource_group": serviceRuntimeNetworkResourceGroup,
		},
	}
}

func flattenOutboundPublicIPAddresses(input *appplatform.NetworkProfile) []interface{} {
	if input == nil || input.OutboundIPs == nil {
		return []interface{}{}
	}

	return utils.FlattenStringSlice(input.OutboundIPs.PublicIPs)
}
