package iotoperations

import (
    "context"
    "fmt"
    "time"

    "github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthorization"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-provider-azurerm/internal/clients"
    "github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBrokerAuthorization() *schema.Resource {
    return &schema.Resource{
        Create: resourceBrokerAuthorizationCreate,
        Read:   resourceBrokerAuthorizationRead,
        Delete: resourceBrokerAuthorizationDelete,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "resource_group_name": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "instance_name": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "broker_name": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "authorization_policies": {
                Type:     schema.TypeList,
                Required: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "cache": {
                            Type:     schema.TypeString,
                            Optional: true,
                        },
                        "rules": {
                            Type:     schema.TypeList,
                            Required: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "broker_resources": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "method": {
                                                    Type:     schema.TypeString,
                                                    Required: true,
                                                },
                                                "client_ids": {
                                                    Type:     schema.TypeList,
                                                    Elem:     &schema.Schema{Type: schema.TypeString},
                                                    Optional: true,
                                                },
                                                "topics": {
                                                    Type:     schema.TypeList,
                                                    Elem:     &schema.Schema{Type: schema.TypeString},
                                                    Optional: true,
                                                },
                                            },
                                        },
                                    },
                                    "principals": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "attributes": {
                                                    Type:     schema.TypeList,
                                                    Elem:     &schema.Schema{Type: schema.TypeMap},
                                                    Optional: true,
                                                },
                                                "client_ids": {
                                                    Type:     schema.TypeList,
                                                    Elem:     &schema.Schema{Type: schema.TypeString},
                                                    Optional: true,
                                                },
                                                "usernames": {
                                                    Type:     schema.TypeList,
                                                    Elem:     &schema.Schema{Type: schema.TypeString},
                                                    Optional: true,
                                                },
                                            },
                                        },
                                    },
                                    "state_store_resources": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "key_type": {
                                                    Type:     schema.TypeString,
                                                    Required: true,
                                                },
                                                "keys": {
                                                    Type:     schema.TypeList,
                                                    Elem:     &schema.Schema{Type: schema.TypeString},
                                                    Required: true,
                                                },
                                                "method": {
                                                    Type:     schema.TypeString,
                                                    Required: true,
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
            "extended_location": {
                Type:     schema.TypeList,
                Optional: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "name": {
                            Type:     schema.TypeString,
                            Required: true,
                        },
                        "type": {
                            Type:     schema.TypeString,
                            Required: true,
                        },
                    },
                },
            },
            "provisioning_state": {
                Type:     schema.TypeString,
                Computed: true,
            },
        },
    }
}

func resourceBrokerAuthorizationCreate(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).IoTOperations.BrokerAuthorizationClient
    ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
    defer cancel()

    subscriptionId := meta.(*clients.Client).Account.SubscriptionId
    resourceGroupName := d.Get("resource_group_name").(string)
    instanceName := d.Get("instance_name").(string)
    brokerName := d.Get("broker_name").(string)
    authorizationName := d.Get("name").(string)

    id := brokerauthorization.NewAuthorizationID(subscriptionId, resourceGroupName, instanceName, brokerName, authorizationName)

    // Build the broker authorization resource
    payload := brokerauthorization.BrokerAuthorizationResource{
        Properties: &brokerauthorization.BrokerAuthorizationProperties{},
    }

    if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
        return fmt.Errorf("creating %s: %+v", id, err)
    }

    d.SetId(id.ID())
    return resourceBrokerAuthorizationRead(d, meta)
}

func resourceBrokerAuthorizationRead(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).IoTOperations.BrokerAuthorizationClient
    ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
    defer cancel()

    id, err := brokerauthorization.ParseAuthorizationID(d.Id())
    if err != nil {
        return err
    }

    resp, err := client.Get(ctx, *id)
    if err != nil {
        return fmt.Errorf("reading %s: %+v", *id, err)
    }

    d.Set("name", id.AuthorizationName)
    d.Set("resource_group_name", id.ResourceGroupName)
    d.Set("instance_name", id.InstanceName)
    d.Set("broker_name", id.BrokerName)

    if model := resp.Model; model != nil {
        // Set properties when needed
    }

    return nil
}

func resourceBrokerAuthorizationDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).IoTOperations.BrokerAuthorizationClient
    ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
    defer cancel()

    id, err := brokerauthorization.ParseAuthorizationID(d.Id())
    if err != nil {
        return err
    }

    if err := client.DeleteThenPoll(ctx, *id); err != nil {
        return fmt.Errorf("deleting %s: %+v", *id, err)
    }

    return nil
}

    d.SetId("")
    return nil
}