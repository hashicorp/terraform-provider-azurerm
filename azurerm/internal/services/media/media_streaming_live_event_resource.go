package media

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMediaLiveEvent() *schema.Resource {
	return &schema.Resource{
		Create: resourceMediaLiveEventCreateUpdate,
		Read:   resourceMediaLiveEventRead,
		Update: resourceMediaLiveEventCreateUpdate,
		Delete: resourceMediaLiveEventDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.LiveEventID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateLiveEventName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"media_services_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateMediaServicesAccountName,
			},

			"auto_start_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"location": azure.SchemaLocation(),

			"input": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_access_control_allow": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"subnet_prefix_length": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},

						"access_token": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"endpoint": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"uri": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"key_frame_interval": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"streaming_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"cross_site_access_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_access_policy": {
							Type:         schema.TypeString,
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"cross_domain_policy": {
							Type:         schema.TypeString,
							Computed:     true,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"encoding": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"key_frame_interval": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"preset_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"stretch_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"hostname_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"preview": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_access_control_allow": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"subnet_prefix_length": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},

						"alternative_media_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"endpoint": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"uri": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"preview_locator": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"streaming_policy_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"transcription_language": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"use_static_hostname": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMediaLiveEventCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.LiveEventsClient
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceID := parse.NewLiveEventID(subscriptionID, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", resourceID, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_media_live_event", resourceID.ID())
		}
	}
	/*
		location := azure.NormalizeLocation(d.Get("location").(string))
		t := d.Get("tags").(map[string]interface{})

		parameters := media.LiveEvent{
			LiveEventProperties: &media.LiveEventProperties{},
			Location:            utils.String(location),
			Tags:                tags.Expand(t),
		}

		autoStart := utils.Bool(false)
		if _, ok := d.GetOk("auto_start_enabled"); ok {
			autoStart = utils.Bool(d.Get("auto_start_enabled").(bool))
		}

		if input, ok := d.GetOk("input"); ok {
			parameters.LiveEventProperties.Input = expandLiveEventInput(input.([]interface{}))
		}

		if crossSitePolicies, ok := d.GetOk("cross_site_access_policy"); ok {
			parameters.LiveEventProperties.CrossSiteAccessPolicies = expandCrossSiteAccessPolicies(crossSitePolicies.([]interface{}))
		}

		if description, ok := d.GetOk("description"); ok {
			parameters.LiveEventProperties.Description = utils.String(description.(string))
		}

		if encoding, ok := d.GetOk("encoding"); ok {
			parameters.LiveEventProperties.Encoding = expandEncoding(encoding.([]interface{}))
		}

		if hostNamePrefix, ok := d.GetOk("hostname_prefix"); ok {
			parameters.LiveEventProperties.HostnamePrefix = utils.String(hostNamePrefix.(string))
		}

		if preview, ok := d.GetOk("preview"); ok {
			parameters.LiveEventProperties.Preview = expandPreview(preview.([]interface{}))
		}

		if transcriptionLanguage, ok := d.GetOk("transcription_language"); ok {
			parameters.LiveEventProperties.Transcriptions = utils.String(transcriptionLanguage.(string))
		}

		if useStaticHostName, ok := d.GetOk("use_static_hostname"); ok {
			parameters.LiveEventProperties.UseStaticHostname = utils.String(useStaticHostName.(string))
		}

		if d.IsNewResource() {
			future, err := client.Create(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name, parameters, autoStart)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", resourceID, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation %s: %+v", resourceID, err)
			}
		} else {
			future, err := client.Update(ctx, resourceID.ResourceGroup, resourceID.MediaserviceName, resourceID.Name, parameters)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", resourceID, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for %s to update: %+v", resourceID, err)
			}
		}

		d.SetId(resourceID.ID())*/

	return resourceMediaLiveEventRead(d, meta)
}

func resourceMediaLiveEventRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.LiveEventsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LiveEventID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.MediaserviceName)
	/*
		if location := resp.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if props := resp.LiveEventProperties; props != nil {
			input := flattenLiveEventInput(props.Input)
			if err := d.Set("input", input); err != nil {
				return fmt.Errorf("Error flattening `input`: %s", err)
			}

			crossSiteAccessPolicies := flattenCrossSiteAccessPolicies(resp.CrossSiteAccessPolicies)
			if err := d.Set("cross_site_access_policy", crossSiteAccessPolicies); err != nil {
				return fmt.Errorf("Error flattening `cross_site_access_policy`: %s", err)
			}

			encoding := flattenEncoding(resp.Encoding)
			if err := d.Set("encoding", encoding); err != nil {
				return fmt.Errorf("Error flattening `cross_site_access_policy`: %s", err)
			}

			d.Set("description", props.Description)
			d.Set("hostname_prefix", props.HostnamePrefix)

			preview := flattenEncoding(resp.Preview)
			if err := d.Set("preview", preview); err != nil {
				return fmt.Errorf("Error flattening `preview`: %s", err)
			}

			d.Set("transcription_language", props.Transcriptions)

			useStaticHostName := false
			if props.UseStaticHostName != nil {
				useStaticHostName = bool(*props.UseStaticHostname)
			}
			d.Set("use_static_hostname", useStaticHostName)
		}*/

	return nil
}

func resourceMediaLiveEventDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.LiveEventsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LiveEventID(d.Id())
	if err != nil {
		return err
	}

	// Stop Live Event before we attempt to delete it.
	/*resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return fmt.Errorf("reading %s: %+v", id, err)
	}
	if props := resp.LiveEventProperties; props != nil {
		if props.ResourceState == media.LiveEventResourceStateRunning {
			stopFuture, err := client.Stop(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
			if err != nil {
				return fmt.Errorf("stopping %s: %+v", id, err)
			}

			if err = stopFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for %s to stop: %+v", id, err)
			}
		}
	}*/

	future, err := client.Delete(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to delete: %+v", id, err)
	}

	return nil
}
