package media

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2020-05-01/media"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMediaStreamingEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMediaStreamingEndpointCreate,
		Read:   resourceArmMediaStreamingEndpointRead,
		Update: resourceArmMediaStreamingEndpointUpdate,
		Delete: resourceArmMediaStreamingEndpointDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.MediaStreamingEndpointID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateStreamingEnpointName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"media_services_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateMediaServicesAccountName,
			},

			"auto_start": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"location": azure.SchemaLocation(),

			"scale_units": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},

			"access_control": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"akamai": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"akamai_signature_header_authentication_key": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"base64_key": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"expiration": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.IsRFC3339Time,
												},
												"identifier": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},

						"ip": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"subnet_prefix_length": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"cdn_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"cdn_profile": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cdn_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cross_site_access_policies": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_access_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"cross_domain_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"custom_host_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"max_cache_age": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMediaStreamingEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	streamingEndpointName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("media_services_account_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	scaleUnits := d.Get("scale_units").(int)
	autoStart := utils.Bool(false)

	parameters := media.StreamingEndpoint{
		StreamingEndpointProperties: &media.StreamingEndpointProperties{
			ScaleUnits: utils.Int32(int32(scaleUnits)),
		},
		Location: utils.String(location),
	}

	if _, ok := d.GetOk("auto_start"); ok {
		autoStart = utils.Bool(d.Get("auto_start").(bool))
	}
	if _, ok := d.GetOk("access_control"); ok {
		parameters.StreamingEndpointProperties.AccessControl = expandAzureRmAccessControl(d)
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

	if _, ok := d.GetOk("cross_site_access_policies"); ok {
		parameters.StreamingEndpointProperties.CrossSiteAccessPolicies = expandAzureRmCrossSiteAccessPolicies(d)
	}

	if _, ok := d.GetOk("custom_host_names"); ok {
		customHostNames := d.Get("custom_host_names").([]interface{})
		parameters.StreamingEndpointProperties.CustomHostNames = utils.ExpandStringSlice(customHostNames)
	}

	if description, ok := d.GetOk("description"); ok {
		parameters.StreamingEndpointProperties.Description = utils.String(description.(string))
	}

	if maxCacheAge, ok := d.GetOk("max_cache_age"); ok {
		parameters.StreamingEndpointProperties.MaxCacheAge = utils.Int64(int64(maxCacheAge.(int)))
	}

	future, err := client.Create(ctx, resourceGroup, accountName, streamingEndpointName, parameters, autoStart)
	if err != nil {
		return fmt.Errorf("Error creating Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
	}

	endpoint, err := client.Get(ctx, resourceGroup, accountName, streamingEndpointName)
	if err != nil {
		return fmt.Errorf("Error retrieving Streaming Endpoint %q from Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
	}

	d.SetId(*endpoint.ID)

	return resourceArmMediaStreamingEndpointRead(d, meta)
}

func resourceArmMediaStreamingEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	streamingEndpointName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("media_services_account_name").(string)
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

		future, err := client.Scale(ctx, resourceGroup, accountName, streamingEndpointName, scaleParamaters)
		if err != nil {
			return fmt.Errorf("Error scaling units in Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for scaling of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
		}

	}

	if _, ok := d.GetOk("access_control"); ok {
		parameters.StreamingEndpointProperties.AccessControl = expandAzureRmAccessControl(d)
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

	if _, ok := d.GetOk("cross_site_access_policies"); ok {
		parameters.StreamingEndpointProperties.CrossSiteAccessPolicies = expandAzureRmCrossSiteAccessPolicies(d)
	}

	if _, ok := d.GetOk("custom_host_names"); ok {
		customHostNames := d.Get("custom_host_names").([]interface{})
		parameters.StreamingEndpointProperties.CustomHostNames = utils.ExpandStringSlice(customHostNames)
	}

	if description, ok := d.GetOk("description"); ok {
		parameters.StreamingEndpointProperties.Description = utils.String(description.(string))
	}

	if maxCacheAge, ok := d.GetOk("max_cache_age"); ok {
		parameters.StreamingEndpointProperties.MaxCacheAge = utils.Int64(int64(maxCacheAge.(int)))
	}

	future, err := client.Update(ctx, resourceGroup, accountName, streamingEndpointName, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
	}

	endpoint, err := client.Get(ctx, resourceGroup, accountName, streamingEndpointName)
	if err != nil {
		return fmt.Errorf("Error retrieving Streaming Endpoint %q from Media Services Account %q (Resource Group %q): %+v", streamingEndpointName, accountName, resourceGroup, err)
	}

	d.SetId(*endpoint.ID)

	return resourceArmMediaStreamingEndpointRead(d, meta)
}

func resourceArmMediaStreamingEndpointRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MediaStreamingEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Streaming Endpoint %q was not found in Media Services Account %q and Resource Group %q - removing from state", id.Name, id.AccountName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.AccountName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.AccountName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.StreamingEndpointProperties; props != nil {

		if scaleUnits := props.ScaleUnits; scaleUnits != nil {
			d.Set("scale_units", scaleUnits)
		}

		accessControl := flattenAzureRmAccessControl(props.AccessControl)
		if err := d.Set("access_control", accessControl); err != nil {
			return fmt.Errorf("Error flattening `access_control`: %s", err)
		}

		if cdnEnabled := props.CdnEnabled; cdnEnabled != nil {
			d.Set("cdn_enabled", cdnEnabled)
		}

		if cdnProfile := props.CdnProfile; cdnProfile != nil {
			d.Set("cdn_profile", cdnProfile)
		}

		if cdnProvider := props.CdnProvider; cdnProvider != nil {
			d.Set("cdn_provider", cdnProvider)
		}

		crossSiteAccessPolicies := flattenAzureRmCrossSiteAccessPolicies(resp.CrossSiteAccessPolicies)
		if err := d.Set("cross_site_access_policies", crossSiteAccessPolicies); err != nil {
			return fmt.Errorf("Error flattening `cross_site_access_policies`: %s", err)
		}

		if customHostNames := props.CustomHostNames; customHostNames != nil {
			d.Set("custom_host_names", customHostNames)
		}

		if description := props.Description; description != nil {
			d.Set("description", description)
		}

		if maxCacheAge := props.MaxCacheAge; maxCacheAge != nil {
			d.Set("max_cache_age", maxCacheAge)
		}
	}

	return nil
}

func resourceArmMediaStreamingEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.StreamingEndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MediaStreamingEndpointID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.AccountName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Streaming Endpoint %q in Media Services Account %q (Resource Group %q): %+v", id.Name, id.AccountName, id.ResourceGroup, err)
	}

	return nil
}

func expandAzureRmAccessControl(d *schema.ResourceData) *media.StreamingEndpointAccessControl {
	accessControls := d.Get("access_control").(*schema.Set).List()
	if len(accessControls) == 0 {
		return nil
	}
	accessControlResult := new(media.StreamingEndpointAccessControl)
	accessControl := accessControls[0].(map[string]interface{})
	// Get IP information
	if raw, ok := accessControl["ip"]; ok {
		ipsList := raw.(*schema.Set).List()
		if len(ipsList) > 0 {
			ip := ipsList[0].(map[string]interface{})
			ipAllowsList := ip["allow"].([]interface{})
			ipRanges := make([]media.IPRange, len(ipAllowsList))
			for index, ipAllow := range ipAllowsList {
				allow := ipAllow.(map[string]interface{})
				address := allow["address"].(string)
				name := allow["name"].(string)

				ipRange := media.IPRange{
					Name:    utils.String(name),
					Address: utils.String(address),
				}
				subnetPrefixLengthRaw := allow["subnet_prefix_length"]
				if subnetPrefixLengthRaw != "" {
					ipRange.SubnetPrefixLength = utils.Int32(int32(subnetPrefixLengthRaw.(int)))
				}
				ipRanges[index] = ipRange
			}
			accessControlResult.IP = &media.IPAccessControl{
				Allow: &ipRanges,
			}
		}

	}
	// Get Akamai information
	if raw, ok := accessControl["akamai"]; ok {
		akamaiList := raw.(*schema.Set).List()
		if len(akamaiList) > 0 {
			akamai := akamaiList[0].(map[string]interface{})
			akamaiSignatureKeyList := akamai["akamai_signature_header_authentication_key"].([]interface{})
			akamaiSignatureHeaderAuthenticationKeyList := make([]media.AkamaiSignatureHeaderAuthenticationKey, len(akamaiSignatureKeyList))
			for index, akamaiSignatureKey := range akamaiSignatureKeyList {
				akamaiKey := akamaiSignatureKey.(map[string]interface{})
				base64Key := akamaiKey["base64_key"].(string)
				expirationRaw := akamaiKey["expiration"].(string)
				identifier := akamaiKey["identifier"].(string)

				akamaiSignatureHeaderAuthenticationKey := media.AkamaiSignatureHeaderAuthenticationKey{
					Base64Key:  utils.String(base64Key),
					Identifier: utils.String(identifier),
				}
				if expirationRaw != "" {
					expiration, _ := date.ParseTime(time.RFC3339, expirationRaw)
					akamaiSignatureHeaderAuthenticationKey.Expiration = &date.Time{
						Time: expiration,
					}
				}
				akamaiSignatureHeaderAuthenticationKeyList[index] = akamaiSignatureHeaderAuthenticationKey

			}
			accessControlResult.Akamai = &media.AkamaiAccessControl{
				AkamaiSignatureHeaderAuthenticationKeyList: &akamaiSignatureHeaderAuthenticationKeyList,
			}
		}

	}

	return accessControlResult
}

func flattenAzureRmAccessControl(accessControl *media.StreamingEndpointAccessControl) []interface{} {
	if accessControl == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if accessControl.IP != nil {
		allows := make([]interface{}, len(*accessControl.IP.Allow))
		for i, ipAllow := range *accessControl.IP.Allow {
			allow := make(map[string]interface{})
			if ipAllow.Name != nil {
				allow["name"] = ipAllow.Name
			}
			if ipAllow.Address != nil {
				allow["address"] = ipAllow.Address
			}
			if ipAllow.SubnetPrefixLength != nil {
				allow["subnet_prefix_length"] = ipAllow.SubnetPrefixLength
			}
			if allow != nil {
				allows[i] = allow
			}
		}

		result["ip"] = []interface{}{map[string]interface{}{
			"allow": allows,
		}}
	}

	if accessControl.Akamai != nil {
		akamaiSignatureKeyList := make([]interface{}, len(*accessControl.Akamai.AkamaiSignatureHeaderAuthenticationKeyList))
		for i, key := range *accessControl.Akamai.AkamaiSignatureHeaderAuthenticationKeyList {
			akamaiSignatureHeaderKey := make(map[string]interface{})
			if key.Base64Key != nil {
				akamaiSignatureHeaderKey["base64_key"] = key.Base64Key
			}
			if key.Expiration != nil {
				akamaiSignatureHeaderKey["expiration"] = key.Expiration.Format(time.RFC3339)
			}
			if key.Identifier != nil {
				akamaiSignatureHeaderKey["identifier"] = key.Identifier
			}
			if akamaiSignatureKeyList != nil {
				akamaiSignatureKeyList[i] = akamaiSignatureHeaderKey
			}
		}

		result["akamai"] = []interface{}{map[string]interface{}{
			"akamai_signature_header_authentication_key": akamaiSignatureKeyList,
		}}

	}

	return []interface{}{result}
}

func expandAzureRmCrossSiteAccessPolicies(d *schema.ResourceData) *media.CrossSiteAccessPolicies {
	crossSiteAccessPolicies := d.Get("cross_site_access_policies").([]interface{})
	crossSiteAccessPolicy := crossSiteAccessPolicies[0].(map[string]interface{})
	clientAccessPolicy := crossSiteAccessPolicy["client_access_policy"].(string)
	crossDomainPolicy := crossSiteAccessPolicy["cross_domain_policy"].(string)
	return &media.CrossSiteAccessPolicies{
		ClientAccessPolicy: &clientAccessPolicy,
		CrossDomainPolicy:  &crossDomainPolicy,
	}
}

func flattenAzureRmCrossSiteAccessPolicies(crossSiteAccessPolicies *media.CrossSiteAccessPolicies) []interface{} {
	if crossSiteAccessPolicies == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if crossSiteAccessPolicies.ClientAccessPolicy != nil {
		result["client_access_policy"] = *crossSiteAccessPolicies.ClientAccessPolicy
	}
	if crossSiteAccessPolicies.CrossDomainPolicy != nil {
		result["cross_domain_policy"] = *crossSiteAccessPolicies.CrossDomainPolicy
	}

	return []interface{}{result}
}
