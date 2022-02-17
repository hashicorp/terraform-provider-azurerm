package signalr

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/webpubsub/mgmt/2021-10-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var defaultRequestTypes = []webpubsub.RequestType{
	webpubsub.RequestTypeRESTAPI,
	webpubsub.RequestTypeClientConnection,
	webpubsub.RequestTypeServerConnection,
	webpubsub.RequestTypeTrace,
}

func resourceWebpubsubNetworkACL() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceWebPubsubNetworkACLCreateUpdate,
		Read:   resourceWebPubsubNetworkACLRead,
		Update: resourceWebPubsubNetworkACLCreateUpdate,
		Delete: resourceWebpubsubNetworkACLDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.DefaultImporter(),

		Schema: map[string]*pluginsdk.Schema{
			"web_pubsub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WebPubsubID,
			},

			"default_action": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  webpubsub.ACLActionDeny,
				ValidateFunc: validation.StringInSlice([]string{
					string(webpubsub.ACLActionAllow),
					string(webpubsub.ACLActionDeny),
				}, false),
			},

			"public_network": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allowed_request_types": {
							Type:          pluginsdk.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"public_network.0.denied_request_types"},
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(webpubsub.RequestTypeClientConnection),
									string(webpubsub.RequestTypeRESTAPI),
									string(webpubsub.RequestTypeServerConnection),
									string(webpubsub.RequestTypeTrace),
								}, false),
							},
						},

						"denied_request_types": {
							Type:          pluginsdk.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"public_network.0.allowed_request_types"},
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(webpubsub.RequestTypeClientConnection),
									string(webpubsub.RequestTypeRESTAPI),
									string(webpubsub.RequestTypeServerConnection),
									string(webpubsub.RequestTypeTrace),
								}, false),
							},
						},
					},
				},
			},

			"private_endpoint": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: networkValidate.PrivateEndpointID,
						},

						"allowed_request_types": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(webpubsub.RequestTypeClientConnection),
									string(webpubsub.RequestTypeRESTAPI),
									string(webpubsub.RequestTypeServerConnection),
									string(webpubsub.RequestTypeTrace),
								}, false),
							},
						},

						"denied_request_types": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(webpubsub.RequestTypeClientConnection),
									string(webpubsub.RequestTypeRESTAPI),
									string(webpubsub.RequestTypeServerConnection),
									string(webpubsub.RequestTypeTrace),
								}, false),
							},
						},
					},
				},
			},
		},
	}
}

func resourceWebPubsubNetworkACLCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubsubClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebPubsubID(d.Get("web_pubsub_id").(string))
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		return fmt.Errorf("checking for present of existing %q: %+v", id, err)
	}

	locks.ByName(id.WebPubSubName, "azurerm_web_pubsub")
	defer locks.UnlockByName(id.WebPubSubName, "azurerm_web_pubsub")

	if d.IsNewResource() {
		if !isNewNetworkACL(existing) {
			return tf.ImportAsExistsError("azurerm_web_pubsub_network_acl", id.ID())
		}
	}

	if props := existing.Properties; props != nil {
		defaultAction := webpubsub.ACLAction(d.Get("default_action").(string))
		networkACL := webpubsub.NetworkACLs{
			DefaultAction:    defaultAction,
			PublicNetwork:    expandWebpubsubPublicNetwork(d.Get("public_network").([]interface{})),
			PrivateEndpoints: expandWebpubsubPrivateEndpoint(d.Get("private_endpoint").(*pluginsdk.Set).List(), props.PrivateEndpointConnections),
		}

		if defaultAction == webpubsub.ACLActionAllow && networkACL.PublicNetwork.Allow != nil && len(*networkACL.PublicNetwork.Allow) != 0 {
			return fmt.Errorf("when `default_action` is `Allow` for `public_network`, `allowed_request_types` cannot be specified")
		} else if defaultAction == webpubsub.ACLActionDeny && len(*networkACL.PublicNetwork.Deny) != 0 {
			return fmt.Errorf("when `default_action` is `Deny` for `public_network`, `denied_request_types` cannot be specified")
		}

		if networkACL.PrivateEndpoints != nil {
			for _, privateEndpoint := range *networkACL.PrivateEndpoints {
				if privateEndpoint.Allow != nil && privateEndpoint.Deny != nil && len(*privateEndpoint.Allow) != 0 && len(*privateEndpoint.Deny) != 0 {
					return fmt.Errorf("`allowed_request_types` and `denied_request_types` cannot be set together for `private_endpoint`")
				}

				if defaultAction == webpubsub.ACLActionAllow && privateEndpoint.Allow != nil && len(*privateEndpoint.Allow) != 0 {
					return fmt.Errorf("when `default_action` is `Allow` for `private_endpoint`, `allowed_request_types` cannot be specified")
				} else if defaultAction == webpubsub.ACLActionDeny && privateEndpoint.Deny != nil && len(*privateEndpoint.Deny) != 0 {
					return fmt.Errorf("when `default_action` is `Deny` for `private_endpoint`, `denied_request_types` cannot be specified")
				}
			}
		}
		existing.Properties.NetworkACLs = &networkACL
	}

	future, err := client.Update(ctx, existing, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		return fmt.Errorf("updating Network ACL configuration for %q: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of network ACL of %q：%+v", id, err)
	}

	d.SetId(id.ID())
	return resourceWebPubsubNetworkACLRead(d, meta)
}

func resourceWebPubsubNetworkACLRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubsubClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebPubsubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %q does not exists - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	d.Set("web_pubsub_id", id.ID())

	props := resp.Properties
	if props != nil && props.NetworkACLs != nil {
		defaultAction := ""
		if props.NetworkACLs.DefaultAction != "" {
			defaultAction = string(props.NetworkACLs.DefaultAction)
		}
		d.Set("default_action", defaultAction)

		if err := d.Set("public_network", flattenWebpubsubPublicNetwork(props.NetworkACLs.PublicNetwork)); err != nil {
			return fmt.Errorf("setting `public_network`: %+v", err)
		}

		if err := d.Set("private_endpoint", flattenWebpubsubPrivatEndpoint(props.NetworkACLs.PrivateEndpoints, props.PrivateEndpointConnections)); err != nil {
			return fmt.Errorf("setting `private_endpoint`: %+v", err)
		}

	}
	return nil
}

func resourceWebpubsubNetworkACLDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.WebPubsubClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebPubsubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	defaultAction := webpubsub.ACLActionDeny
	var denyRequestTypes []webpubsub.RequestType
	networkACL := &webpubsub.NetworkACLs{
		DefaultAction: defaultAction,
		PublicNetwork: &webpubsub.NetworkACL{
			Allow: &defaultRequestTypes,
			Deny:  &denyRequestTypes,
		},
	}

	if resp.Properties != nil && resp.Properties.NetworkACLs != nil && resp.Properties.NetworkACLs.PrivateEndpoints != nil {
		privateEndpoints := make([]webpubsub.PrivateEndpointACL, 0)
		for _, item := range *resp.Properties.NetworkACLs.PrivateEndpoints {
			privateEndpoints = append(privateEndpoints, webpubsub.PrivateEndpointACL{
				Allow: &defaultRequestTypes,
				Name:  item.Name,
			})
		}
		networkACL.PrivateEndpoints = &privateEndpoints
	}

	if resp.Properties != nil {
		resp.Properties.NetworkACLs = networkACL
	}

	future, err := client.Update(ctx, resp, id.ResourceGroup, id.WebPubSubName)
	if err != nil {
		return fmt.Errorf("resetting the default Network ACL configuration for %q: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of network ACL of %q：%+v", id, err)
	}

	return nil
}

func expandWebpubsubPublicNetwork(input []interface{}) *webpubsub.NetworkACL {
	allowRTs := make([]webpubsub.RequestType, 0)
	deniedRTs := make([]webpubsub.RequestType, 0)

	if len(input) == 0 || input[0] == nil {
		return &webpubsub.NetworkACL{}
	}

	v := input[0].(map[string]interface{})

	for _, item := range *(utils.ExpandStringSlice(v["allowed_request_types"].(*pluginsdk.Set).List())) {
		allowRTs = append(allowRTs, webpubsub.RequestType(item))
	}

	for _, item := range *(utils.ExpandStringSlice(v["denied_request_types"].(*pluginsdk.Set).List())) {
		deniedRTs = append(deniedRTs, webpubsub.RequestType(item))
	}

	return &webpubsub.NetworkACL{
		Allow: &allowRTs,
		Deny:  &deniedRTs,
	}
}

func flattenWebpubsubPublicNetwork(input *webpubsub.NetworkACL) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	allowRequestTypes := make([]string, 0)
	if input.Allow != nil {
		for _, item := range *input.Allow {
			allowRequestTypes = append(allowRequestTypes, (string)(item))
		}
	}
	allow := utils.FlattenStringSlice(&allowRequestTypes)

	deniedRequestTypes := make([]string, 0)
	if input.Deny != nil {
		for _, item := range *input.Deny {
			deniedRequestTypes = append(deniedRequestTypes, (string)(item))
		}
	}
	deny := utils.FlattenStringSlice(&deniedRequestTypes)

	return []interface{}{
		map[string]interface{}{
			"allowed_request_types": allow,
			"denied_request_types":  deny,
		},
	}
}

func expandWebpubsubPrivateEndpoint(input []interface{}, privateEndpointConnections *[]webpubsub.PrivateEndpointConnection) *[]webpubsub.PrivateEndpointACL {
	results := make([]webpubsub.PrivateEndpointACL, 0)
	if privateEndpointConnections == nil {
		return &results
	}

	for _, privateEndpointConnection := range *privateEndpointConnections {
		result := webpubsub.PrivateEndpointACL{
			Allow: &[]webpubsub.RequestType{},
			Deny:  &[]webpubsub.RequestType{},
		}

		if privateEndpointConnection.Name != nil {
			result.Name = privateEndpointConnection.Name
		}

		for _, item := range input {
			v := item.(map[string]interface{})
			privateEndpointId := v["id"].(string)

			if props := privateEndpointConnection.PrivateEndpointConnectionProperties; props != nil {
				if props.PrivateEndpoint == nil || props.PrivateEndpoint.ID == nil || privateEndpointId != *props.PrivateEndpoint.ID {
					continue
				}

				allowedRTs := make([]webpubsub.RequestType, 0)
				for _, item := range *(utils.ExpandStringSlice(v["allowed_request_types"].(*pluginsdk.Set).List())) {
					allowedRTs = append(allowedRTs, webpubsub.RequestType(item))
				}
				result.Allow = &allowedRTs

				deniedRTs := make([]webpubsub.RequestType, 0)
				for _, item := range *(utils.ExpandStringSlice(v["denied_request_types"].(*pluginsdk.Set).List())) {
					deniedRTs = append(deniedRTs, webpubsub.RequestType(item))
				}
				result.Deny = &deniedRTs

				break
			}
		}
		results = append(results, result)
	}
	return &results
}

func flattenWebpubsubPrivatEndpoint(input *[]webpubsub.PrivateEndpointACL, privateEndpointConnections *[]webpubsub.PrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if privateEndpointConnections != nil {
			for _, privateEndpointConnection := range *privateEndpointConnections {
				if privateEndpointConnection.Name == nil || privateEndpointConnection.PrivateEndpointConnectionProperties == nil {
					continue
				}
				if !strings.EqualFold(*privateEndpointConnection.Name, *item.Name) {
					continue
				}
				props := privateEndpointConnection.PrivateEndpointConnectionProperties
				if props.PrivateEndpoint == nil || props.PrivateEndpoint.ID == nil {
					continue
				}

				allowedRequestTypes := make([]string, 0)
				if item.Allow != nil {
					for _, item := range *item.Allow {
						allowedRequestTypes = append(allowedRequestTypes, string(item))
					}
				}
				allow := utils.FlattenStringSlice(&allowedRequestTypes)

				deniedRequestTypes := make([]string, 0)
				if item.Deny != nil {
					for _, item := range *item.Deny {
						deniedRequestTypes = append(deniedRequestTypes, (string)(item))
					}
				}
				deny := utils.FlattenStringSlice(&deniedRequestTypes)

				results = append(results, map[string]interface{}{
					"id":                    *props.PrivateEndpoint.ID,
					"allowed_request_types": allow,
					"denied_request_types":  deny,
				})

				break
			}
		}
	}

	return results
}

func isNewNetworkACL(existing webpubsub.ResourceType) bool {
	if existing.Properties == nil || existing.Properties.NetworkACLs == nil {
		return true
	}

	if existing.Properties.NetworkACLs.DefaultAction != webpubsub.ACLActionDeny {
		return false
	}

	if existing.Properties.NetworkACLs.PublicNetwork == nil ||
		existing.Properties.NetworkACLs.PublicNetwork.Deny != nil ||
		(existing.Properties.NetworkACLs.PublicNetwork.Allow == nil || len(*existing.Properties.NetworkACLs.PublicNetwork.Allow) != len(defaultRequestTypes)) {
		return false
	}

	defaultRequestTypeMap := map[webpubsub.RequestType]bool{
		webpubsub.RequestTypeClientConnection: true,
		webpubsub.RequestTypeRESTAPI:          true,
		webpubsub.RequestTypeServerConnection: true,
		webpubsub.RequestTypeTrace:            true,
	}

	for _, allowType := range *existing.Properties.NetworkACLs.PublicNetwork.Allow {
		if !defaultRequestTypeMap[allowType] {
			return false
		}
	}

	if existing.Properties.NetworkACLs.PrivateEndpoints != nil {
		for _, peItem := range *existing.Properties.NetworkACLs.PrivateEndpoints {
			if peItem.Allow == nil || len(*peItem.Allow) != len(defaultRequestTypes) || peItem.Deny != nil {
				return false
			}
			for _, allowType := range *peItem.Allow {
				if !defaultRequestTypeMap[allowType] {
					return false
				}
			}
		}
	}
	return true
}
