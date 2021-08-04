package signalr

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr/sdk/signalr"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSignalRServiceNetworkACL() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSignalRServiceNetworkACLCreateUpdate,
		Read:   resourceSignalRServiceNetworkACLRead,
		Update: resourceSignalRServiceNetworkACLCreateUpdate,
		Delete: resourceSignalRServiceNetworkACLDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.DefaultImporter(),

		Schema: map[string]*pluginsdk.Schema{
			"signalr_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServiceID,
			},

			"default_action": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(signalr.ACLActionAllow),
					string(signalr.ACLActionDeny),
				}, false),
			},

			"public_network": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// API response includes the `Trace` type but it isn't in rest api client.
						// https://github.com/Azure/azure-rest-api-specs/issues/14923
						"allowed_request_types": {
							Type:          pluginsdk.TypeSet,
							Optional:      true,
							ConflictsWith: []string{"public_network.0.denied_request_types"},
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(signalr.SignalRRequestTypeClientConnection),
									string(signalr.SignalRRequestTypeRESTAPI),
									string(signalr.SignalRRequestTypeServerConnection),
									"Trace",
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
									string(signalr.SignalRRequestTypeClientConnection),
									string(signalr.SignalRRequestTypeRESTAPI),
									string(signalr.SignalRRequestTypeServerConnection),
									"Trace",
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

						// API response includes the `Trace` type but it isn't in rest api client.
						// https://github.com/Azure/azure-rest-api-specs/issues/14923
						"allowed_request_types": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(signalr.SignalRRequestTypeClientConnection),
									string(signalr.SignalRRequestTypeRESTAPI),
									string(signalr.SignalRRequestTypeServerConnection),
									"Trace",
								}, false),
							},
						},

						"denied_request_types": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(signalr.SignalRRequestTypeClientConnection),
									string(signalr.SignalRRequestTypeRESTAPI),
									string(signalr.SignalRRequestTypeServerConnection),
									"Trace",
								}, false),
							},
						},
					},
				},
			},
		},
	}
}

func resourceSignalRServiceNetworkACLCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.Client
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Get("signalr_service_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.SignalRName, "azurerm_signalr_service")
	defer locks.UnlockByName(id.SignalRName, "azurerm_signalr_service")

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", *id)
	}

	model := *resp.Model
	if props := model.Properties; props != nil {
		defaultAction := signalr.ACLAction(d.Get("default_action").(string))
		networkACL := signalr.SignalRNetworkACLs{
			DefaultAction: &defaultAction,
			PublicNetwork: expandSignalRServicePublicNetwork(d.Get("public_network").([]interface{})),
		}

		if v, ok := d.GetOk("private_endpoint"); ok {
			networkACL.PrivateEndpoints = expandSignalRServicePrivateEndpoint(v.(*pluginsdk.Set).List(), props.PrivateEndpointConnections)
		}

		if defaultAction == signalr.ACLActionAllow && len(*networkACL.PublicNetwork.Allow) != 0 {
			return fmt.Errorf("when `default_action` is `Allow` for `public_network`, `allowed_request_types` cannot be specified")
		} else if defaultAction == signalr.ACLActionDeny && len(*networkACL.PublicNetwork.Deny) != 0 {
			return fmt.Errorf("when `default_action` is `Deny` for `public_network`, `denied_request_types` cannot be specified")
		}

		if networkACL.PrivateEndpoints != nil {
			for _, privateEndpoint := range *networkACL.PrivateEndpoints {
				if len(*privateEndpoint.Allow) != 0 && len(*privateEndpoint.Deny) != 0 {
					return fmt.Errorf("`allowed_request_types` and `denied_request_types` cannot be set together for `private_endpoint`")
				}

				if defaultAction == signalr.ACLActionAllow && len(*privateEndpoint.Allow) != 0 {
					return fmt.Errorf("when `default_action` is `Allow` for `private_endpoint`, `allowed_request_types` cannot be specified")
				} else if defaultAction == signalr.ACLActionDeny && len(*privateEndpoint.Deny) != 0 {
					return fmt.Errorf("when `default_action` is `Deny` for `private_endpoint`, `denied_request_types` cannot be specified")
				}
			}
		}

		model.Properties.NetworkACLs = &networkACL
	}

	if err := client.UpdateThenPoll(ctx, *id, model); err != nil {
		return fmt.Errorf("creating/updating NetworkACL for %s: %v", id, err)
	}

	d.SetId(id.ID())
	return resourceSignalRServiceNetworkACLRead(d, meta)
}

func resourceSignalRServiceNetworkACLRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("signalr_service_id", id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.NetworkACLs != nil {
			defaultAction := ""
			if props.NetworkACLs.DefaultAction != nil {
				defaultAction = string(*props.NetworkACLs.DefaultAction)
			}
			d.Set("default_action", defaultAction)

			if err := d.Set("public_network", flattenSignalRServicePublicNetwork(props.NetworkACLs.PublicNetwork)); err != nil {
				return fmt.Errorf("setting `public_network`: %+v", err)
			}

			if err := d.Set("private_endpoint", flattenSignalRServicePrivateEndpoint(props.NetworkACLs.PrivateEndpoints, props.PrivateEndpointConnections)); err != nil {
				return fmt.Errorf("setting `private_endpoint`: %+v", err)
			}
		}
	}

	return nil
}

func resourceSignalRServiceNetworkACLDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.Client
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.SignalRName, "azurerm_signalr_service")
	defer locks.UnlockByName(id.SignalRName, "azurerm_signalr_service")

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", *id)
	}

	model := *resp.Model

	defaultAction := signalr.ACLActionDeny
	defaultRequestTypes := []signalr.SignalRRequestType{
		signalr.SignalRRequestTypeClientConnection,
		signalr.SignalRRequestTypeRESTAPI,
		signalr.SignalRRequestTypeServerConnection,
		"Trace",
	}
	networkACL := &signalr.SignalRNetworkACLs{
		DefaultAction: &defaultAction,
		PublicNetwork: &signalr.NetworkACL{
			Allow: &defaultRequestTypes,
		},
	}

	if model.Properties != nil && model.Properties.NetworkACLs != nil && model.Properties.NetworkACLs.PrivateEndpoints != nil {
		privateEndpoints := make([]signalr.PrivateEndpointACL, 0)
		for _, item := range *model.Properties.NetworkACLs.PrivateEndpoints {
			privateEndpoints = append(privateEndpoints, signalr.PrivateEndpointACL{
				Allow: &defaultRequestTypes,
				Name:  item.Name,
			})
		}
		networkACL.PrivateEndpoints = &privateEndpoints
	}

	if model.Properties != nil {
		model.Properties.NetworkACLs = networkACL
	}

	if err := client.UpdateThenPoll(ctx, *id, model); err != nil {
		return fmt.Errorf("resetting the default Network ACL configuration for %s: %+v", *id, err)
	}

	return nil
}

func expandSignalRServicePublicNetwork(input []interface{}) *signalr.NetworkACL {
	allowedRTs := make([]signalr.SignalRRequestType, 0)
	deniedRTs := make([]signalr.SignalRRequestType, 0)

	if len(input) != 0 && input[0] != nil {
		v := input[0].(map[string]interface{})

		for _, item := range *(utils.ExpandStringSlice(v["allowed_request_types"].(*pluginsdk.Set).List())) {
			allowedRTs = append(allowedRTs, signalr.SignalRRequestType(item))
		}

		for _, item := range *(utils.ExpandStringSlice(v["denied_request_types"].(*pluginsdk.Set).List())) {
			deniedRTs = append(deniedRTs, signalr.SignalRRequestType(item))
		}
	}

	return &signalr.NetworkACL{
		Allow: &allowedRTs,
		Deny:  &deniedRTs,
	}
}

func expandSignalRServicePrivateEndpoint(input []interface{}, privateEndpointConnections *[]signalr.PrivateEndpointConnection) *[]signalr.PrivateEndpointACL {
	results := make([]signalr.PrivateEndpointACL, 0)
	if privateEndpointConnections == nil {
		return &results
	}

	for _, privateEndpointConnection := range *privateEndpointConnections {
		result := signalr.PrivateEndpointACL{
			Allow: &[]signalr.SignalRRequestType{},
			Deny:  &[]signalr.SignalRRequestType{},
		}

		if privateEndpointConnection.Name != nil {
			result.Name = *privateEndpointConnection.Name
		}

		for _, item := range input {
			v := item.(map[string]interface{})
			privateEndpointId := v["id"].(string)

			if props := privateEndpointConnection.Properties; props != nil {
				if props.PrivateEndpoint == nil || props.PrivateEndpoint.Id == nil || privateEndpointId != *props.PrivateEndpoint.Id {
					continue
				}

				allowedRTs := make([]signalr.SignalRRequestType, 0)
				for _, item := range *(utils.ExpandStringSlice(v["allowed_request_types"].(*pluginsdk.Set).List())) {
					allowedRTs = append(allowedRTs, signalr.SignalRRequestType(item))
				}
				result.Allow = &allowedRTs

				deniedRTs := make([]signalr.SignalRRequestType, 0)
				for _, item := range *(utils.ExpandStringSlice(v["denied_request_types"].(*pluginsdk.Set).List())) {
					deniedRTs = append(deniedRTs, signalr.SignalRRequestType(item))
				}
				result.Deny = &deniedRTs

				break
			}
		}

		results = append(results, result)
	}

	return &results
}

func flattenSignalRServicePublicNetwork(input *signalr.NetworkACL) []interface{} {
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

func flattenSignalRServicePrivateEndpoint(input *[]signalr.PrivateEndpointACL, privateEndpointConnections *[]signalr.PrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if privateEndpointConnections != nil {
			for _, privateEndpointConnection := range *privateEndpointConnections {
				if privateEndpointConnection.Name == nil || privateEndpointConnection.Properties == nil {
					continue
				}
				if *privateEndpointConnection.Name != item.Name {
					continue
				}
				props := privateEndpointConnection.Properties
				if props.PrivateEndpoint == nil || props.PrivateEndpoint.Id == nil {
					continue
				}

				allowedRequestTypes := make([]string, 0)
				if item.Allow != nil {
					for _, item := range *item.Allow {
						allowedRequestTypes = append(allowedRequestTypes, (string)(item))
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
					"id":                    *props.PrivateEndpoint.Id,
					"allowed_request_types": allow,
					"denied_request_types":  deny,
				})

				break
			}
		}
	}

	return results
}
