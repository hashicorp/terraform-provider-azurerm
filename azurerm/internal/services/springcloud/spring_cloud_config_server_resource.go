package springcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/appplatform/mgmt/2020-07-01/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSpringCloudConfigServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpringCloudConfigServerCreateUpdate,
		Read:   resourceSpringCloudConfigServerRead,
		Update: resourceSpringCloudConfigServerCreateUpdate,
		Delete: resourceSpringCloudConfigServerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SpringCloudServiceID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"spring_cloud_service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudServiceID,
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

			"search_paths": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"http_basic_auth": SchemaConfigServerHttpBasicAuth(false, "ssh_auth"),

			"ssh_auth": SchemaConfigServerSSHAuth(false, "http_basic_auth"),

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

						"http_basic_auth": SchemaConfigServerHttpBasicAuth(false),

						"ssh_auth": SchemaConfigServerSSHAuth(false),
					},
				},
			},
		},
	}
}

func resourceSpringCloudConfigServerCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ConfigServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serviceId, err := parse.SpringCloudServiceID(d.Get("spring_cloud_service_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewSpringCloudServiceID(subscriptionId, serviceId.ResourceGroup, serviceId.SpringName)
	existing, err := client.Get(ctx, serviceId.ResourceGroup, serviceId.SpringName)
	if err != nil {
		return fmt.Errorf("making Read request on %s: %+v", serviceId, err)
	}
	if d.IsNewResource() {
		if existing.Properties != nil && existing.Properties.ConfigServer != nil && existing.Properties.ConfigServer.GitProperty != nil && existing.Properties.ConfigServer.GitProperty.URI != nil {
			return tf.ImportAsExistsError("azurerm_spring_cloud_config_server", id.ID())
		}
	}

	gitProperty := appplatform.ConfigServerGitProperty{
		URI:         utils.String(d.Get("uri").(string)),
		Label:       utils.String(d.Get("label").(string)),
		SearchPaths: utils.ExpandStringSlice(d.Get("search_paths").([]interface{})),
	}

	httpBasicAuth := d.Get("http_basic_auth").([]interface{})
	sshAuth := d.Get("ssh_auth").([]interface{})
	if len(httpBasicAuth) > 0 && len(sshAuth) > 0 {
		return fmt.Errorf("can not set both `http_basic_auth` and `ssh_auth`")
	}
	if len(httpBasicAuth) > 0 {
		v := httpBasicAuth[0].(map[string]interface{})
		gitProperty.Username = utils.String(v["username"].(string))
		gitProperty.Password = utils.String(v["password"].(string))
	}
	if len(sshAuth) > 0 {
		v := sshAuth[0].(map[string]interface{})
		gitProperty.PrivateKey = utils.String(v["private_key"].(string))
		gitProperty.StrictHostKeyChecking = utils.Bool(v["strict_host_key_checking_enabled"].(bool))

		if hostKey := v["host_key"].(string); hostKey != "" {
			gitProperty.HostKey = utils.String(hostKey)
		}
		if hostKeyAlgorithm := v["host_key_algorithm"].(string); hostKeyAlgorithm != "" {
			gitProperty.HostKeyAlgorithm = utils.String(hostKeyAlgorithm)
		}
	}

	repositories, err := expandSpringCloudGitPatternRepository(d.Get("repository").([]interface{}))
	if err != nil {
		return err
	}
	gitProperty.Repositories = repositories

	existing.Properties = &appplatform.ConfigServerProperties{
		ConfigServer: &appplatform.ConfigServerSettings{
			GitProperty: &gitProperty,
		},
	}

	future, err := client.UpdatePut(ctx, id.ResourceGroup, id.SpringName, existing)
	if err != nil {
		return fmt.Errorf("updating config server for %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for updation of config server for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSpringCloudConfigServerRead(d, meta)
}

func resourceSpringCloudConfigServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ConfigServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud Config Server %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("unable to read Config Server for %s: %+v", id, err)
	}

	if resp.Properties == nil || resp.Properties.ConfigServer == nil || resp.Properties.ConfigServer.GitProperty == nil {
		log.Printf("[INFO] Spring Cloud Config Server %q does not exist - removing from state", d.Id())
		d.SetId("")
		return nil
	}

	prop := resp.Properties.ConfigServer.GitProperty

	d.Set("spring_cloud_service_id", id.ID())
	d.Set("uri", prop.URI)
	d.Set("label", prop.Label)
	d.Set("search_paths", utils.FlattenStringSlice(prop.SearchPaths))

	if err := d.Set("http_basic_auth", flattenSpringCloudConfigServerHttpBasicAuth(prop, d.Get("http_basic_auth").([]interface{}))); err != nil {
		return fmt.Errorf("setting `http_basic_auth`: %+v", err)
	}

	if err := d.Set("ssh_auth", flattenSpringCloudConfigServerSSHAuth(prop, d.Get("ssh_auth").([]interface{}))); err != nil {
		return fmt.Errorf("setting `ssh_auth`: %+v", err)
	}

	if err := d.Set("repository", flattenSpringCloudGitPatternRepository(prop.Repositories, d.Get("repository").([]interface{}))); err != nil {
		return fmt.Errorf("setting `config_server_git_setting`: %+v", err)
	}
	return nil
}

func resourceSpringCloudConfigServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ConfigServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudServiceID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		return fmt.Errorf("making Read request on Config Server for %s: %+v", id, err)
	}

	existing.Properties = nil

	future, err := client.UpdatePut(ctx, id.ResourceGroup, id.SpringName, existing)
	if err != nil {
		return fmt.Errorf("failure deleting Config Server for %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failure waiting for deleting Config Server for %s: %+v", id, err)
	}

	return nil
}

// username and password returned by API are *
// to avoid state diff, we get the props from old state
func flattenSpringCloudConfigServerHttpBasicAuth(input *appplatform.ConfigServerGitProperty, httpBasicAuthOldState []interface{}) []interface{} {
	if input == nil || input.Username == nil || input.Password == nil {
		return []interface{}{}
	}

	username := ""
	password := ""
	if len(httpBasicAuthOldState) > 0 {
		oldItem := httpBasicAuthOldState[0].(map[string]interface{})
		username = oldItem["username"].(string)
		password = oldItem["password"].(string)
	}
	return []interface{}{
		map[string]interface{}{
			"username": username,
			"password": password,
		},
	}
}

// private_key, host_key and host_key_algorithm returned by API are *
// to avoid state diff, we get the props from old state
func flattenSpringCloudConfigServerSSHAuth(input *appplatform.ConfigServerGitProperty, sshAuthOldState []interface{}) []interface{} {
	if input == nil || input.PrivateKey == nil {
		return []interface{}{}
	}

	privateKey := ""
	hostKey := ""
	hostKeyAlgorithm := ""
	if len(sshAuthOldState) > 0 {
		oldItem := sshAuthOldState[0].(map[string]interface{})
		privateKey = oldItem["private_key"].(string)
		hostKey = oldItem["host_key"].(string)
		hostKeyAlgorithm = oldItem["host_key_algorithm"].(string)
	}

	strictHostKeyChecking := false
	if input.StrictHostKeyChecking != nil {
		strictHostKeyChecking = *input.StrictHostKeyChecking
	}
	return []interface{}{
		map[string]interface{}{
			"private_key":                      privateKey,
			"host_key":                         hostKey,
			"host_key_algorithm":               hostKeyAlgorithm,
			"strict_host_key_checking_enabled": strictHostKeyChecking,
		},
	}
}
