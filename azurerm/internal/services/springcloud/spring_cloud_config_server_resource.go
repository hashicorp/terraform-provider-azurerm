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

			"http_basic_auth": SchemaConfigServerHttpBasicAuth("ssh_auth"),

			"ssh_auth": SchemaConfigServerSSHAuth("http_basic_auth"),

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
	resourceId := parse.NewSpringCloudServiceID(subscriptionId, serviceId.ResourceGroup, serviceId.SpringName).ID()

	existing, err := client.Get(ctx, serviceId.ResourceGroup, serviceId.SpringName)
	if err != nil {
		return fmt.Errorf("making Read request on Spring Cloud Service %q (Resource Group %q): %+v", serviceId.SpringName, serviceId.ResourceGroup, err)
	}
	if d.IsNewResource() {
		if existing.Properties != nil && existing.Properties.ConfigServer != nil && existing.Properties.ConfigServer.GitProperty != nil && existing.Properties.ConfigServer.GitProperty.URI != nil {
			return tf.ImportAsExistsError("azurerm_spring_cloud_config_server", resourceId)
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

	future, err := client.UpdatePut(ctx, serviceId.ResourceGroup, serviceId.SpringName, existing)
	if err != nil {
		return fmt.Errorf("failure updating config server of Spring Cloud Service %q  (Resource Group %q): %+v", serviceId.SpringName, serviceId.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failure waiting for setting config server of Spring Cloud Service %q config server (Resource Group %q): %+v", serviceId.SpringName, serviceId.ResourceGroup, err)
	}

	d.SetId(resourceId)

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
		return fmt.Errorf("unable to read Spring Cloud Config Server %q (Resource Group %q): %+v", id.SpringName, id.ResourceGroup, err)
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

	// username and password returned by API are *
	// to avoid state diff, we get the props from old state
	httpBasicAuth := make([]interface{}, 0)
	if prop.Username != nil && prop.Password != nil {
		username := ""
		password := ""
		if oldHTTPBasicAuth := d.Get("http_basic_auth").([]interface{}); len(oldHTTPBasicAuth) > 0 {
			oldItem := oldHTTPBasicAuth[0].(map[string]interface{})
			username = oldItem["username"].(string)
			password = oldItem["password"].(string)
		}
		httpBasicAuth = []interface{}{
			map[string]interface{}{
				"username": username,
				"password": password,
			},
		}
	}
	d.Set("http_basic_auth", httpBasicAuth)

	sshAuth := []interface{}{}
	if prop.PrivateKey != nil {
		privateKey := ""
		hostKey := ""
		hostKeyAlgorithm := ""
		if sshAuth := d.Get("ssh_auth").([]interface{}); len(sshAuth) > 0 {
			oldItem := sshAuth[0].(map[string]interface{})
			privateKey = oldItem["private_key"].(string)
			hostKey = oldItem["host_key"].(string)
			hostKeyAlgorithm = oldItem["host_key_algorithm"].(string)
		}

		strictHostKeyChecking := false
		if prop.StrictHostKeyChecking != nil {
			strictHostKeyChecking = *prop.StrictHostKeyChecking
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
	d.Set("ssh_auth", sshAuth)

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
		return fmt.Errorf("making Read request on Spring Cloud Config Server %q (Resource Group %q): %+v", id.SpringName, id.ResourceGroup, err)
	}

	existing.Properties = nil

	future, err := client.UpdatePut(ctx, id.ResourceGroup, id.SpringName, existing)
	if err != nil {
		return fmt.Errorf("failure deleting Spring Cloud Config Server %q (Resource Group %q): %+v", id.SpringName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failure waiting for deleting Spring Cloud Config Server %q (Resource Group %q): %+v", id.SpringName, id.ResourceGroup, err)
	}

	return nil
}
