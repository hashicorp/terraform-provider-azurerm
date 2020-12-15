package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2020-05-01/media"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMediaContentKeyPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMediaContentKeyPolicyCreate,
		Read:   resourceMediaContentKeyPolicyRead,
		Update: resourceMediaContentKeyPolicyUpdate,
		Delete: resourceMediaContentKeyPolicyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ContentKeyPolicyID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9(_)]{1,128}$"),
					"Content Key Policy name must be 1 - 128 characters long, can contain letters, numbers, underscores, and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"media_services_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateMediaServicesAccountName,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"policy_option": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"clear_key_configuration_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"widevine_configuration_template": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"playready_configuration_license": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_test_devices": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"begin_date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.ValidateRFC3339TimeString,
									},
									"content_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.ContentKeyPolicyPlayReadyContentTypeUltraVioletDownload),
											string(media.ContentKeyPolicyPlayReadyContentTypeUltraVioletStreaming),
											string(media.ContentKeyPolicyPlayReadyContentTypeUnspecified),
											string(media.ContentKeyPolicyPlayReadyContentTypeUnknown),
										}, false),
									},
									"expiration_date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.ValidateRFC3339TimeString,
									},
									"grace_period": {
										Type:         schema.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"license_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.ContentKeyPolicyPlayReadyLicenseTypeNonPersistent),
											string(media.ContentKeyPolicyPlayReadyLicenseTypePersistent),
											string(media.ContentKeyPolicyPlayReadyLicenseTypeUnknown),
										}, false),
									},
									"play_right": { //TODO:Complete definition
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"playback_duration_seconds": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(1),
												},
												"storage_duration_seconds": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(1),
												},
											},
										},
									},
									"relative_begin_date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.ValidateRFC3339TimeString,
									},
									"relative_expiration_date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.ValidateRFC3339TimeString,
									},
								},
							},
						},
						"fairplay_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ask": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"pfx": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"pfx_password": {
										Type:         schema.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"offline_rental_configuration": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.ContentKeyPolicyRestrictionTokenTypeJwt),
											string(media.ContentKeyPolicyRestrictionTokenTypeSwt),
											string(media.ContentKeyPolicyRestrictionTokenTypeUnknown),
										}, false),
									},
									"rental_and_lease_key_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.DualExpiry),
											string(media.PersistentLimited),
											string(media.PersistentUnlimited),
											string(media.Undefined),
											string(media.Unknown),
										}, false),
									},
									"rental_duration": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"token_restriction": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"audience": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"issuer": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"token_type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(media.ContentKeyPolicyRestrictionTokenTypeJwt),
											string(media.ContentKeyPolicyRestrictionTokenTypeSwt),
											string(media.ContentKeyPolicyRestrictionTokenTypeUnknown),
										}, false),
									},
									"primary_symmetric_token_key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"alternative_symmetric_token_key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"primary_rsa_token_key_exponent": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"primary_rsa_token_key_modulus": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"alternate_rsa_token_key_exponent": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"alternate_rsa_token_key_modulus": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"primary_x509_token_key_raw": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"alternate_x509_token_key_raw": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"open_id_connect_discovery_document": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"required_claim": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"value": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
								},
							},
						},
						"open_restriction_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMediaContentKeyPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	streamingEndpointName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("media_services_account_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	scaleUnits := d.Get("scale_units").(int)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resourceId := parse.NewStreamingEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.MediaserviceName, resourceId.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of %s: %+v", resourceId, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_media_streaming_endpoint", resourceId.ID())
	}

	parameters := media.StreamingEndpoint{
		StreamingEndpointProperties: &media.StreamingEndpointProperties{
			ScaleUnits: utils.Int32(int32(scaleUnits)),
		},
		Location: utils.String(location),
	}

	autoStart := utils.Bool(false)
	if _, ok := d.GetOk("auto_start_enabled"); ok {
		autoStart = utils.Bool(d.Get("auto_start_enabled").(bool))
	}
	if _, ok := d.GetOk("access_control"); ok {
		accessControl, err := expandAccessControl(d)
		if err != nil {
			return err
		}
		parameters.StreamingEndpointProperties.AccessControl = accessControl
	}
	if cdnEnabled, ok := d.GetOk("cdn_enabled"); ok {
		parameters.StreamingEndpointProperties.CdnEnabled = utils.Bool(cdnEnabled.(bool))
	}

	if cdnProfile, ok := d.GetOk("cdn_profile"); ok {
		parameters.StreamingEndpointProperties.CdnProfile = utils.String(cdnProfile.(string))
	}

	if cdnProvider, ok := d.GetOk("cdn_provider"); ok {
		parameters.StreamingEndpointProperties.CdnProvider = utils.String(cdnProvider.(string))
	}

	if crossSite, ok := d.GetOk("cross_site_access_policy"); ok {
		parameters.StreamingEndpointProperties.CrossSiteAccessPolicies = expandCrossSiteAccessPolicies(crossSite.([]interface{}))
	}

	if _, ok := d.GetOk("custom_host_names"); ok {
		customHostNames := d.Get("custom_host_names").([]interface{})
		parameters.StreamingEndpointProperties.CustomHostNames = utils.ExpandStringSlice(customHostNames)
	}

	if description, ok := d.GetOk("description"); ok {
		parameters.StreamingEndpointProperties.Description = utils.String(description.(string))
	}

	if maxCacheAge, ok := d.GetOk("max_cache_age_seconds"); ok {
		parameters.StreamingEndpointProperties.MaxCacheAge = utils.Int64(int64(maxCacheAge.(int)))
	}

	future, err := client.Create(ctx, resourceGroup, accountName, streamingEndpointName, parameters, autoStart)
	if err != nil {
		return fmt.Errorf("Error creating Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
	}

	d.SetId(resourceId.ID())

	return resourceMediaContentKeyPolicyRead(d, meta)
}

func resourceMediaContentKeyPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingEndpointID(d.Id())
	if err != nil {
		return err
	}
	location := azure.NormalizeLocation(d.Get("location").(string))
	scaleUnits := d.Get("scale_units").(int)

	parameters := media.StreamingEndpoint{
		StreamingEndpointProperties: &media.StreamingEndpointProperties{
			ScaleUnits: utils.Int32(int32(scaleUnits)),
		},
		Location: utils.String(location),
	}

	if d.HasChange("scale_units") {
		scaleParamaters := media.StreamingEntityScaleUnit{
			ScaleUnit: utils.Int32(int32(scaleUnits)),
		}

		future, err := client.Scale(ctx, id.ResourceGroup, id.MediaserviceName, id.Name, scaleParamaters)
		if err != nil {
			return fmt.Errorf("Error scaling units in Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for scaling of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
		}
	}

	if _, ok := d.GetOk("access_control"); ok {
		accessControl, err := expandAccessControl(d)
		if err != nil {
			return err
		}
		parameters.StreamingEndpointProperties.AccessControl = accessControl
	}

	if cdnEnabled, ok := d.GetOk("cdn_enabled"); ok {
		parameters.StreamingEndpointProperties.CdnEnabled = utils.Bool(cdnEnabled.(bool))
	}

	if cdnProfile, ok := d.GetOk("cdn_profile"); ok {
		parameters.StreamingEndpointProperties.CdnProfile = utils.String(cdnProfile.(string))
	}

	if cdnProvider, ok := d.GetOk("cdn_provider"); ok {
		parameters.StreamingEndpointProperties.CdnProvider = utils.String(cdnProvider.(string))
	}

	if crossSitePolicies, ok := d.GetOk("cross_site_access_policy"); ok {
		parameters.StreamingEndpointProperties.CrossSiteAccessPolicies = expandCrossSiteAccessPolicies(crossSitePolicies.([]interface{}))
	}

	if _, ok := d.GetOk("custom_host_names"); ok {
		customHostNames := d.Get("custom_host_names").([]interface{})
		parameters.StreamingEndpointProperties.CustomHostNames = utils.ExpandStringSlice(customHostNames)
	}

	if description, ok := d.GetOk("description"); ok {
		parameters.StreamingEndpointProperties.Description = utils.String(description.(string))
	}

	if maxCacheAge, ok := d.GetOk("max_cache_age_seconds"); ok {
		parameters.StreamingEndpointProperties.MaxCacheAge = utils.Int64(int64(maxCacheAge.(int)))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.MediaserviceName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return resourceMediaContentKeyPolicyRead(d, meta)
}

func resourceMediaContentKeyPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Streaming Endpoint %q was not found in Media Services Account %q and Resource Group %q - removing from state", id.Name, id.MediaserviceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.MediaserviceName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.StreamingEndpointProperties; props != nil {
		if scaleUnits := props.ScaleUnits; scaleUnits != nil {
			d.Set("scale_units", scaleUnits)
		}

		accessControl := flattenAccessControl(props.AccessControl)
		if err := d.Set("access_control", accessControl); err != nil {
			return fmt.Errorf("Error flattening `access_control`: %s", err)
		}

		d.Set("cdn_enabled", props.CdnEnabled)
		d.Set("cdn_profile", props.CdnProfile)
		d.Set("cdn_provider", props.CdnProvider)

		crossSiteAccessPolicies := flattenCrossSiteAccessPolicies(resp.CrossSiteAccessPolicies)
		if err := d.Set("cross_site_access_policy", crossSiteAccessPolicies); err != nil {
			return fmt.Errorf("Error flattening `cross_site_access_policy`: %s", err)
		}

		d.Set("custom_host_names", props.CustomHostNames)
		d.Set("description", props.Description)

		maxCacheAge := 0
		if props.MaxCacheAge != nil {
			maxCacheAge = int(*props.MaxCacheAge)
		}
		d.Set("max_cache_age_seconds", maxCacheAge)
	}

	return nil
}

func resourceMediaContentKeyPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingEndpointID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.MediaserviceName, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.MediaserviceName, id.ResourceGroup, err)
	}

	return nil
}
