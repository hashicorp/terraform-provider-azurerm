package iotoperations

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthentication"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBrokerAuthentication() *schema.Resource {
	return &schema.Resource{
		Create: resourceBrokerAuthenticationCreate,
		Read:   resourceBrokerAuthenticationRead,
		Delete: resourceBrokerAuthenticationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"broker_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authentication_methods": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": {
							Type:     schema.TypeString,
							Required: true,
						},
						"custom_settings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auth": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"x509": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"secret_ref": {Type: schema.TypeString, Optional: true},
														},
													},
												},
											},
										},
									},
									"ca_cert_config_map": {Type: schema.TypeString, Optional: true},
									"endpoint":           {Type: schema.TypeString, Optional: true},
									"headers": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"service_account_token_settings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"audiences": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"x509_settings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authorization_attributes": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"attributes": {
													Type:     schema.TypeMap,
													Optional: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"subject": {Type: schema.TypeString, Optional: true},
											},
										},
									},
									"trusted_client_ca_cert": {Type: schema.TypeString, Optional: true},
								},
							},
						},
					},
				},
			},
			"extended_location": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {Type: schema.TypeString, Required: true},
						"type": {Type: schema.TypeString, Required: true},
					},
				},
			},
		},
	}
}

func resourceBrokerAuthenticationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTOperations.BrokerAuthenticationClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resourceGroupName := d.Get("resource_group_name").(string)
	instanceName := d.Get("instance_name").(string)
	brokerName := d.Get("broker_name").(string)
	authenticationName := d.Get("name").(string)

	id := brokerauthentication.NewAuthenticationID(subscriptionId, resourceGroupName, instanceName, brokerName, authenticationName)

	// Build the broker authentication resource
	payload := brokerauthentication.BrokerAuthenticationResource{
		Properties: &brokerauthentication.BrokerAuthenticatorProperties{},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceBrokerAuthenticationRead(d, meta)
}

func resourceBrokerAuthenticationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTOperations.BrokerAuthenticationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := brokerauthentication.ParseAuthenticationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.AuthenticationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("instance_name", id.InstanceName)
	d.Set("broker_name", id.BrokerName)

	if model := resp.Model; model != nil {
	}

	return nil
}

func resourceBrokerAuthenticationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTOperations.BrokerAuthenticationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := brokerauthentication.ParseAuthenticationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.AuthenticationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("instance_name", id.InstanceName)
	d.Set("broker_name", id.BrokerName)

	if model := resp.Model; model != nil {
		if model.Properties != nil {
			// Set authentication methods if available
			if model.Properties.AuthenticationMethods != nil {
				authMethods := make([]interface{}, 0)
				for _, method := range *model.Properties.AuthenticationMethods {
					authMethod := make(map[string]interface{})
					if method.Method != nil {
						authMethod["method"] = string(*method.Method)
					}
					// Add other method properties as needed based on the SDK structure
					authMethods = append(authMethods, authMethod)
				}
				d.Set("authentication_methods", authMethods)
			}
		}

		// Set extended location if available
		if model.ExtendedLocation != nil {
			extendedLocation := make([]interface{}, 0)
			extLocation := make(map[string]interface{})
			if model.ExtendedLocation.Name != nil {
				extLocation["name"] = *model.ExtendedLocation.Name
			}
			if model.ExtendedLocation.Type != nil {
				extLocation["type"] = string(*model.ExtendedLocation.Type)
			}
			extendedLocation = append(extendedLocation, extLocation)
			d.Set("extended_location", extendedLocation)
		}
	}

	return nil
}

func resourceBrokerAuthenticationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTOperations.BrokerAuthenticationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := brokerauthentication.ParseAuthenticationID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
