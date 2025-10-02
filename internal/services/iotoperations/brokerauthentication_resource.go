package iotoperations

import (
    "context"
    "fmt"
    "time"

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

func resourceBrokerAuthenticationDelete(d *schema.ResourceData, meta interface{}) error {
                satMap := sat[0].(map[string]interface{})
                satSettings := &armiotoperations.BrokerAuthenticatorMethodSat{}
                if v, ok := satMap["audiences"].([]interface{}); ok && len(v) > 0 {
                    var audiences []*string
                    for _, a := range v {
                        audiences = append(audiences, to.Ptr(a.(string)))
                    }
                    satSettings.Audiences = audiences
                }
                method.ServiceAccountTokenSettings = satSettings
            }

            // X509Settings
            if x509, ok := mMap["x509_settings"].([]interface{}); ok && len(x509) > 0 {
                x509Map := x509[0].(map[string]interface{})
                x509Settings := &armiotoperations.BrokerAuthenticatorMethodX509{}
                if attrs, ok := x509Map["authorization_attributes"].(map[string]interface{}); ok && len(attrs) > 0 {
                    authAttrs := make(map[string]*armiotoperations.BrokerAuthenticatorMethodX509Attributes)
                    for k, v := range attrs {
                        attrMap := v.(map[string]interface{})
                        x509Attr := &armiotoperations.BrokerAuthenticatorMethodX509Attributes{}
                        if a, ok := attrMap["attributes"].(map[string]interface{}); ok && len(a) > 0 {
                            attributes := make(map[string]*string)
                            for ak, av := range a {
                                attributes[ak] = to.Ptr(av.(string))
                            }
                            x509Attr.Attributes = attributes
                        }
                        if s, ok := attrMap["subject"].(string); ok && s != "" {
                            x509Attr.Subject = to.Ptr(s)
                        }
                        authAttrs[k] = x509Attr
                    }
                    x509Settings.AuthorizationAttributes = authAttrs
                }
                if v, ok := x509Map["trusted_client_ca_cert"].(string); ok && v != "" {
                    x509Settings.TrustedClientCaCert = to.Ptr(v)
                }
                method.X509Settings = x509Settings
            }

            methods = append(methods, method)
        }
    }

    // Build ExtendedLocation if present
    var extLoc *armiotoperations.ExtendedLocation
    if v, ok := d.GetOk("extended_location"); ok && len(v.([]interface{})) > 0 {
        ext := v.([]interface{})[0].(map[string]interface{})
        extLoc = &armiotoperations.ExtendedLocation{
            Name: to.Ptr(ext["name"].(string)),
            Type: to.Ptr(armiotoperations.ExtendedLocationType(ext["type"].(string))),
        }
    }

    resource := armiotoperations.BrokerAuthenticationResource{
        Properties: &armiotoperations.BrokerAuthenticationProperties{
            AuthenticationMethods: methods,
        },
        ExtendedLocation: extLoc,
    }

    poller, err := c.BrokerAuthenticationClient.BeginCreateOrUpdate(
        ctx, resourceGroup, instanceName, brokerName, name, resource, nil)
    if err != nil {
        return fmt.Errorf("creating BrokerAuthenticationResource: %+v", err)
    }

    _, err = poller.PollUntilDone(ctx, nil)
    if err != nil {
        return fmt.Errorf("waiting for BrokerAuthenticationResource creation: %+v", err)
    }

    d.SetId(name)
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
        // Set properties when needed
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