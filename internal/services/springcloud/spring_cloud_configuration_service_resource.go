// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudConfigurationService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSpringCloudConfigurationServiceCreateUpdate,
		Read:   resourceSpringCloudConfigurationServiceRead,
		Update: resourceSpringCloudConfigurationServiceCreateUpdate,
		Delete: resourceSpringCloudConfigurationServiceDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SpringCloudConfigurationServiceV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudConfigurationServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"default",
				}, false),
			},

			"spring_cloud_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudServiceID,
			},

			"generation": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(appplatform.ConfigurationServiceGenerationGen1),
					string(appplatform.ConfigurationServiceGenerationGen2),
				}, false),
			},

			"repository": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"label": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"patterns": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"uri": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"ca_certificate_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.SpringCloudCertificateID,
						},

						"host_key": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"host_key_algorithm": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"password": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"private_key": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"search_paths": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"strict_host_key_checking": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"username": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}
func resourceSpringCloudConfigurationServiceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).AppPlatform.ConfigurationServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	springId, err := parse.SpringCloudServiceID(d.Get("spring_cloud_service_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewSpringCloudConfigurationServiceID(subscriptionId, springId.ResourceGroup, springId.SpringName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ConfigurationServiceName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_configuration_service", id.ID())
		}
	}

	configurationServiceResource := appplatform.ConfigurationServiceResource{
		Properties: &appplatform.ConfigurationServiceProperties{
			Generation: appplatform.ConfigurationServiceGeneration(d.Get("generation").(string)),
			Settings: &appplatform.ConfigurationServiceSettings{
				GitProperty: &appplatform.ConfigurationServiceGitProperty{
					Repositories: expandConfigurationServiceConfigurationServiceGitRepositoryArray(d.Get("repository").([]interface{})),
				},
			},
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.ConfigurationServiceName, configurationServiceResource)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpringCloudConfigurationServiceRead(d, meta)
}

func resourceSpringCloudConfigurationServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ConfigurationServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudConfigurationServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ConfigurationServiceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] appplatform %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.ConfigurationServiceName)
	d.Set("spring_cloud_service_id", parse.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroup, id.SpringName).ID())
	if props := resp.Properties; props != nil {
		d.Set("generation", props.Generation)
		if props.Settings != nil && props.Settings.GitProperty != nil {
			d.Set("repository", flattenConfigurationServiceConfigurationServiceGitRepositoryArray(props.Settings.GitProperty.Repositories, d.Get("repository").([]interface{})))
		}
	}
	return nil
}

func resourceSpringCloudConfigurationServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.ConfigurationServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudConfigurationServiceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.ConfigurationServiceName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}
	return nil
}

func expandConfigurationServiceConfigurationServiceGitRepositoryArray(input []interface{}) *[]appplatform.ConfigurationServiceGitRepository {
	if len(input) == 0 {
		return nil
	}
	results := make([]appplatform.ConfigurationServiceGitRepository, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		repo := appplatform.ConfigurationServiceGitRepository{
			Name:                  utils.String(v["name"].(string)),
			Patterns:              utils.ExpandStringSlice(v["patterns"].(*pluginsdk.Set).List()),
			URI:                   utils.String(v["uri"].(string)),
			Label:                 utils.String(v["label"].(string)),
			SearchPaths:           utils.ExpandStringSlice(v["search_paths"].(*pluginsdk.Set).List()),
			Username:              utils.String(v["username"].(string)),
			Password:              utils.String(v["password"].(string)),
			HostKey:               utils.String(v["host_key"].(string)),
			HostKeyAlgorithm:      utils.String(v["host_key_algorithm"].(string)),
			PrivateKey:            utils.String(v["private_key"].(string)),
			StrictHostKeyChecking: utils.Bool(v["strict_host_key_checking"].(bool)),
		}
		if caCertificatedId := v["ca_certificate_id"].(string); caCertificatedId != "" {
			repo.CaCertResourceID = utils.String(caCertificatedId)
		}
		results = append(results, repo)
	}
	return &results
}

func flattenConfigurationServiceConfigurationServiceGitRepositoryArray(input *[]appplatform.ConfigurationServiceGitRepository, old []interface{}) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	oldItems := make(map[string]map[string]interface{})
	for _, item := range old {
		v := item.(map[string]interface{})
		if name, ok := v["name"]; ok {
			oldItems[name.(string)] = v
		}
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}
		var label string
		if item.Label != nil {
			label = *item.Label
		}
		var uri string
		if item.URI != nil {
			uri = *item.URI
		}

		var strictHostKeyChecking bool
		if item.StrictHostKeyChecking != nil {
			strictHostKeyChecking = *item.StrictHostKeyChecking
		}

		var hostKey string
		var hostKeyAlgorithm string
		var privateKey string
		var username string
		var password string
		if oldItem, ok := oldItems[name]; ok {
			if value, ok := oldItem["host_key"]; ok {
				hostKey = value.(string)
			}
			if value, ok := oldItem["host_key_algorithm"]; ok {
				hostKeyAlgorithm = value.(string)
			}
			if value, ok := oldItem["password"]; ok {
				password = value.(string)
			}
			if value, ok := oldItem["private_key"]; ok {
				privateKey = value.(string)
			}
			if value, ok := oldItem["username"]; ok {
				username = value.(string)
			}
		}

		var caCertificateId string
		if item.CaCertResourceID != nil {
			certificatedId, err := parse.SpringCloudCertificateIDInsensitively(*item.CaCertResourceID)
			if err == nil {
				caCertificateId = certificatedId.ID()
			}
		}
		results = append(results, map[string]interface{}{
			"ca_certificate_id":        caCertificateId,
			"name":                     name,
			"label":                    label,
			"patterns":                 utils.FlattenStringSlice(item.Patterns),
			"uri":                      uri,
			"host_key":                 hostKey,
			"host_key_algorithm":       hostKeyAlgorithm,
			"password":                 password,
			"private_key":              privateKey,
			"search_paths":             utils.FlattenStringSlice(item.SearchPaths),
			"strict_host_key_checking": strictHostKeyChecking,
			"username":                 username,
		})
	}
	return results
}
