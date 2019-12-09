package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2019-05-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	azappplatform "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSpringCloudConfigServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSpringCloudConfigServerCreate,
		Read:   resourceArmSpringCloudConfigServerRead,
		Update: resourceArmSpringCloudConfigServerUpdate,
		Delete: resourceArmSpringCloudConfigServerDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"spring_cloud_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"uri": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azappplatform.ValidateConfigServerURI,
			},
			"host_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_key_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repositories": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"uri": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azappplatform.ValidateConfigServerURI,
						},
						"host_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"host_key_algorithm": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"label": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"pattern": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"private_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"search_paths": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"strict_host_key_checking": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"search_paths": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"strict_host_key_checking": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceArmSpringCloudConfigServerCreate(d *schema.ResourceData, meta interface{}) error {
	springCloudId := d.Get("spring_cloud_id").(string)
	gitProperty := expandArmSpringCloudConfigServerGitProperty(d)

	if err := modifySpringCloudConfigServer(d, meta, springCloudId, gitProperty); err != nil {
		return err
	}

	d.SetId(springCloudId)

	return resourceArmSpringCloudConfigServerRead(d, meta)
}

func resourceArmSpringCloudConfigServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).AppPlatform.ServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["Spring"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Spring Cloud %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("spring_cloud_id", resp.ID)
	if resp.Properties != nil && resp.Properties.ConfigServerProperties != nil && resp.Properties.ConfigServerProperties.ConfigServer != nil {
		if props := resp.Properties.ConfigServerProperties.ConfigServer.GitProperty; props != nil {
			d.Set("host_key", props.HostKey)
			d.Set("host_key_algorithm", props.HostKeyAlgorithm)
			d.Set("label", props.Label)
			d.Set("password", props.Password)
			d.Set("private_key", props.PrivateKey)
			d.Set("strict_host_key_checking", props.StrictHostKeyChecking)
			d.Set("uri", props.URI)
			d.Set("username", props.Username)
			d.Set("search_paths", props.SearchPaths)
			if err := d.Set("repositories", flattenArmSpringCloudGitPatternRepository(props.Repositories)); err != nil {
				return fmt.Errorf("Error setting `repositories`: %+v", err)
			}
		}
	}

	return nil
}

func resourceArmSpringCloudConfigServerUpdate(d *schema.ResourceData, meta interface{}) error {
	springCloudId := d.Get("spring_cloud_id").(string)
	gitProperty := expandArmSpringCloudConfigServerGitProperty(d)

	if err := modifySpringCloudConfigServer(d, meta, springCloudId, gitProperty); err != nil {
		return err
	}

	return resourceArmSpringCloudConfigServerRead(d, meta)
}

func resourceArmSpringCloudConfigServerDelete(d *schema.ResourceData, meta interface{}) error {
	return modifySpringCloudConfigServer(d, meta, d.Id(), nil)
}

func modifySpringCloudConfigServer(d *schema.ResourceData, meta interface{}, springCloudId string, gitProperty *appplatform.ConfigServerGitProperty) error {
	client := meta.(*ArmClient).AppPlatform.ServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(springCloudId)
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	springCloudName := id.Path["Spring"]

	resource := appplatform.ServiceResource{
		Properties: &appplatform.ClusterResourceProperties{
			ConfigServerProperties: &appplatform.ConfigServerProperties{
				ConfigServer: &appplatform.ConfigServerSettings{
					GitProperty: gitProperty,
				},
			},
		},
	}

	if resp, err := client.Get(ctx, resourceGroup, springCloudName); err != nil {
		return fmt.Errorf("Error reading Spring Cloud %q (Resource Group %q): %+v", springCloudName, resourceGroup, err)
	} else {
		resource.Tags = resp.Tags
	}

	future, err := client.Update(ctx, resourceGroup, springCloudName, &resource)
	if err != nil {
		return fmt.Errorf("Error setting Spring Cloud %q config server (Resource Group %q): %+v", springCloudName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for setting of Spring Cloud %q config server (Resource Group %q): %+v", springCloudName, resourceGroup, err)
	}
	springCloudService, err := future.Result(*client)
	if err != nil {
		return fmt.Errorf("Error getting result of Spring Cloud %q config server (Resource Group %q): %+v", springCloudName, resourceGroup, err)
	}
	if springCloudService.Properties != nil && springCloudService.Properties.ConfigServerProperties != nil {
		if err := springCloudService.Properties.ConfigServerProperties.Error; err != nil {
			return fmt.Errorf("Error setting of Spring Cloud %q config server (Resource Group %q): %+v", springCloudName, resourceGroup, err)
		}
	}

	return nil
}

func expandArmSpringCloudConfigServerGitProperty(d *schema.ResourceData) *appplatform.ConfigServerGitProperty {
	result := appplatform.ConfigServerGitProperty{}

	if repositories, ok := d.GetOk("repositories"); ok {
		result.Repositories = expandArmSpringCloudGitPatternRepository(repositories.([]interface{}))
	}
	if uri, ok := d.GetOk("uri"); ok {
		result.URI = utils.String(uri.(string))
	}
	if label, ok := d.GetOk("label"); ok {
		result.Label = utils.String(label.(string))
	}
	if searchPaths, ok := d.GetOk("search_paths"); ok {
		result.SearchPaths = utils.ExpandStringSlice(searchPaths.([]interface{}))
	}
	if username, ok := d.GetOk("username"); ok {
		result.Username = utils.String(username.(string))
	}
	if password, ok := d.GetOk("password"); ok {
		result.Password = utils.String(password.(string))
	}
	if hostKey, ok := d.GetOk("host_key"); ok {
		result.HostKey = utils.String(hostKey.(string))
	}
	if hostKeyAlgorithm, ok := d.GetOk("host_key_algorithm"); ok {
		result.HostKeyAlgorithm = utils.String(hostKeyAlgorithm.(string))
	}
	if privateKey, ok := d.GetOk("private_key"); ok {
		result.PrivateKey = utils.String(privateKey.(string))
	}
	if strictHostKeyChecking, ok := d.GetOk("strict_host_key_checking"); ok {
		result.StrictHostKeyChecking = utils.Bool(strictHostKeyChecking.(bool))
	}

	return &result
}

func expandArmSpringCloudGitPatternRepository(input []interface{}) *[]appplatform.GitPatternRepository {
	results := make([]appplatform.GitPatternRepository, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		name := v["name"].(string)
		pattern := v["pattern"].([]interface{})
		uri := v["uri"].(string)
		label := v["label"].(string)
		searchPaths := v["search_paths"].([]interface{})

		result := appplatform.GitPatternRepository{
			Label:       utils.String(label),
			Name:        utils.String(name),
			Pattern:     utils.ExpandStringSlice(pattern),
			SearchPaths: utils.ExpandStringSlice(searchPaths),
			URI:         utils.String(uri),
		}

		if username := v["username"].(string); len(username) > 0 {
			result.Username = utils.String(username)
		}
		if password := v["password"].(string); len(password) > 0 {
			result.Password = utils.String(password)
		}
		if hostKey := v["host_key"].(string); len(hostKey) > 0 {
			result.HostKey = utils.String(hostKey)
		}
		if hostKeyAlgorithm := v["host_key_algorithm"].(string); len(hostKeyAlgorithm) > 0 {
			result.HostKeyAlgorithm = utils.String(hostKeyAlgorithm)
		}
		if privateKey := v["private_key"].(string); len(privateKey) > 0 {
			result.PrivateKey = utils.String(privateKey)
		}
		if strictHostKeyChecking := v["strict_host_key_checking"].(bool); strictHostKeyChecking {
			result.StrictHostKeyChecking = utils.Bool(strictHostKeyChecking)
		}

		results = append(results, result)
	}
	return &results
}

func flattenArmSpringCloudGitPatternRepository(input *[]appplatform.GitPatternRepository) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if name := item.Name; name != nil {
			v["name"] = *name
		}
		if hostKey := item.HostKey; hostKey != nil {
			v["host_key"] = *hostKey
		}
		if hostKeyAlgorithm := item.HostKeyAlgorithm; hostKeyAlgorithm != nil {
			v["host_key_algorithm"] = *hostKeyAlgorithm
		}
		if label := item.Label; label != nil {
			v["label"] = *label
		}
		if password := item.Password; password != nil {
			v["password"] = *password
		}
		v["pattern"] = utils.FlattenStringSlice(item.Pattern)
		if privateKey := item.PrivateKey; privateKey != nil {
			v["private_key"] = *privateKey
		}
		v["search_paths"] = utils.FlattenStringSlice(item.SearchPaths)
		if strictHostKeyChecking := item.StrictHostKeyChecking; strictHostKeyChecking != nil {
			v["strict_host_key_checking"] = *strictHostKeyChecking
		}
		if uri := item.URI; uri != nil {
			v["uri"] = *uri
		}
		if username := item.Username; username != nil {
			v["username"] = *username
		}

		results = append(results, v)
	}

	return results
}
