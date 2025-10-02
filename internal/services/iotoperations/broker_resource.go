package iotoperations

import (
    "context"
    "fmt"
    "regexp"
    "time"

    "github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/broker"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
    "github.com/hashicorp/terraform-provider-azurerm/internal/clients"
    "github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBroker() *schema.Resource {
    return &schema.Resource{
        Create: resourceBrokerCreate,
        Read:   resourceBrokerRead,
        Update: resourceBrokerUpdate,
        Delete: resourceBrokerDelete,

        Schema: map[string]*schema.Schema{
            "broker_name": {
                Type:         schema.TypeString,
                Required:     true,
                ValidateFunc: validation.All(
                    validation.StringLenBetween(3, 63),
                    validation.StringMatch(
                        regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
                        "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$",
                    ),
                ),
                Description: "Name of broker.",
            },
            "instance_name": {
                Type:         schema.TypeString,
                Required:     true,
                ValidateFunc: validation.All(
                    validation.StringLenBetween(3, 63),
                    validation.StringMatch(
                        regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
                        "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$",
                    ),
                ),
                Description: "Name of instance.",
            },
            "resource_group_name": {
                Type:         schema.TypeString,
                Required:     true,
                ValidateFunc: validation.StringLenBetween(1, 90),
                Description:  "The name of the resource group. The name is case insensitive.",
            },
            "subscription_id": {
                Type:         schema.TypeString,
                Required:     true,
                ValidateFunc: validation.IsUUID,
                Description:  "The ID of the target subscription. The value must be an UUID.",
            },
            "api_version": {
                Type:         schema.TypeString,
                Required:     true,
                ValidateFunc: validation.StringLenAtLeast(1),
                Description:  "The API version to use for this operation.",
            },
            "properties": {
                Type:     schema.TypeList,
                Required: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "advanced": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "clients": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "max_session_expiry_seconds": {Type: schema.TypeInt, Optional: true},
                                                "max_message_expiry_seconds": {Type: schema.TypeInt, Optional: true},
                                                "max_packet_size_bytes":      {Type: schema.TypeInt, Optional: true},
                                                "subscriber_queue_limit": {
                                                    Type:     schema.TypeList,
                                                    Optional: true,
                                                    Elem: &schema.Resource{
                                                        Schema: map[string]*schema.Schema{
                                                            "length":   {Type: schema.TypeInt, Optional: true},
                                                            "strategy": {Type: schema.TypeString, Optional: true},
                                                        },
                                                    },
                                                },
                                                "max_receive_maximum":   {Type: schema.TypeInt, Optional: true},
                                                "max_keep_alive_seconds": {Type: schema.TypeInt, Optional: true},
                                            },
                                        },
                                    },
                                    "encrypt_internal_traffic": {Type: schema.TypeString, Optional: true},
                                    "internal_certs": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "duration":     {Type: schema.TypeString, Optional: true},
                                                "renew_before": {Type: schema.TypeString, Optional: true},
                                                "private_key": {
                                                    Type:     schema.TypeList,
                                                    Optional: true,
                                                    Elem: &schema.Resource{
                                                        Schema: map[string]*schema.Schema{
                                                            "algorithm":       {Type: schema.TypeString, Optional: true},
                                                            "rotation_policy": {Type: schema.TypeString, Optional: true},
                                                        },
                                                    },
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                        "cardinality": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "backend_chain": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "partitions": {
                                                    Type:         schema.TypeInt,
                                                    Optional:     true,
                                                    ValidateFunc: validation.IntBetween(1, 16),
                                                    Description:  "The desired number of physical backend partitions.",
                                                },
                                                "redundancy_factor": {
                                                    Type:         schema.TypeInt,
                                                    Optional:     true,
                                                    ValidateFunc: validation.IntBetween(1, 5),
                                                    Description:  "The desired numbers of backend replicas (pods) in a physical partition.",
                                                },
                                                "workers": {
                                                    Type:         schema.TypeInt,
                                                    Optional:     true,
                                                    ValidateFunc: validation.IntBetween(1, 16),
                                                    Description:  "Number of logical backend workers per replica (pod).",
                                                },
                                            },
                                        },
                                    },
                                    "frontend": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "replicas": {
                                                    Type:         schema.TypeInt,
                                                    Optional:     true,
                                                    ValidateFunc: validation.IntBetween(1, 16),
                                                    Description:  "The desired number of frontend instances (pods).",
                                                },
                                                "workers": {
                                                    Type:         schema.TypeInt,
                                                    Optional:     true,
                                                    ValidateFunc: validation.IntBetween(1, 16),
                                                    Description:  "Number of logical frontend workers per instance (pod).",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                        "diagnostics": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "logs": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "level": {Type: schema.TypeString, Optional: true},
                                            },
                                        },
                                    },
                                    "metrics": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "prometheus_port": {
                                                    Type:         schema.TypeInt,
                                                    Optional:     true,
                                                    ValidateFunc: validation.IntBetween(0, 65535),
                                                    Description:  "The prometheus port to expose the metrics.",
                                                },
                                            },
                                        },
                                    },
                                    "self_check": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "mode":           {Type: schema.TypeString, Optional: true},
                                                "interval_seconds": {Type: schema.TypeInt, Optional: true},
                                                "timeout_seconds":  {Type: schema.TypeInt, Optional: true},
                                            },
                                        },
                                    },
                                    "traces": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "mode":                {Type: schema.TypeString, Optional: true},
                                                "cache_size_megabytes": {Type: schema.TypeInt, Optional: true},
                                                "self_tracing": {
                                                    Type:     schema.TypeList,
                                                    Optional: true,
                                                    Elem: &schema.Resource{
                                                        Schema: map[string]*schema.Schema{
                                                            "mode":            {Type: schema.TypeString, Optional: true},
                                                            "interval_seconds": {Type: schema.TypeInt, Optional: true},
                                                        },
                                                    },
                                                },
                                                "span_channel_capacity": {Type: schema.TypeInt, Optional: true},
                                            },
                                        },
                                    },
                                },
                            },
                        },
                        "disk_backed_message_buffer": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "max_size": {Type: schema.TypeString, Optional: true},
                                    "ephemeral_volume_claim_spec": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem:     volumeClaimSpecSchema(),
                                    },
                                    "persistent_volume_claim_spec": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem:     volumeClaimSpecSchema(),
                                    },
                                },
                            },
                        },
                        "generate_resource_limits": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "cpu": {Type: schema.TypeString, Optional: true},
                                },
                            },
                        },
                        "memory_profile": {Type: schema.TypeString, Optional: true},
                        "provisioning_state": {
                            Type:        schema.TypeString,
                            Computed:    true,
                            Description: "The status of the last operation.",
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

func volumeClaimSpecSchema() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "volume_name":        {Type: schema.TypeString, Optional: true},
            "volume_mode":        {Type: schema.TypeString, Optional: true},
            "storage_class_name": {Type: schema.TypeString, Optional: true},
            "access_modes": {
                Type:     schema.TypeList,
                Optional: true,
                Elem:     &schema.Schema{Type: schema.TypeString},
            },
            "data_source": {
                Type:     schema.TypeList,
                Optional: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "api_group": {Type: schema.TypeString, Optional: true},
                        "kind":      {Type: schema.TypeString, Optional: true},
                        "name":      {Type: schema.TypeString, Optional: true},
                    },
                },
            },
            "data_source_ref": {
                Type:     schema.TypeList,
                Optional: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "api_group": {Type: schema.TypeString, Optional: true},
                        "kind":      {Type: schema.TypeString, Optional: true},
                        "name":      {Type: schema.TypeString, Optional: true},
                        "namespace": {Type: schema.TypeString, Optional: true},
                    },
                },
            },
            "resources": {
                Type:     schema.TypeList,
                Optional: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "limits": {
                            Type:     schema.TypeMap,
                            Optional: true,
                            Elem:     &schema.Schema{Type: schema.TypeString},
                        },
                        "requests": {
                            Type:     schema.TypeMap,
                            Optional: true,
                            Elem:     &schema.Schema{Type: schema.TypeString},
                        },
                    },
                },
            },
            "selector": {
                Type:     schema.TypeList,
                Optional: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "match_expressions": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "key":      {Type: schema.TypeString, Optional: true},
                                    "operator": {Type: schema.TypeString, Optional: true},
                                    "values": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        Elem:     &schema.Schema{Type: schema.TypeString},
                                    },
                                },
                            },
                        },
                        "match_labels": {
                            Type:     schema.TypeMap,
                            Optional: true,
                            Elem:     &schema.Schema{Type: schema.TypeString},
                        },
                    },
                },
            },
        },
    }
}

// CRUD functions
func resourceBrokerCreate(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Iot.BrokerClient 
    ctx := meta.(*clients.Client).StopContext

    brokerName := d.Get("broker_name").(string)
    resourceGroupName := d.Get("resource_group_name").(string)
    instanceName := d.Get("instance_name").(string)
    subscriptionId := d.Get("subscription_id").(string)
    apiVersion := d.Get("api_version").(string)
    properties := d.Get("properties").([]interface{})

    // Build the broker ID using SDk
    id := broker.NewBrokerID(subscriptionId, resourceGroupName, instanceName, brokerName)

    // Build the request model using SDK
    req := broker.BrokerResource{
        Properties: expandBrokerProperties(properties),
    }

    // Call the Azure API
     resp, err := client.Create(ctx, id, req)
     if err != nil {
        return fmt.Errorf("error creating broker: %+v", err)
     }

    // Set resource ID and fields
    d.SetId(id.ID())
    d.Set("provisioning_state", resp.Properties.ProvisioningState)
     d.Set("properties", flattenBrokerProperties(resp.Properties))

    // For now, set a dummy ID
    d.SetId(brokerName)
    return resourceBrokerRead(d, meta)
}

func resourceBrokerRead(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Iot.BrokerClient // Adjust to your client
    ctx := meta.(*clients.Client).StopContext

    brokerName := d.Get("broker_name").(string)
    resourceGroupName := d.Get("resource_group_name").(string)
    instanceName := d.Get("instance_name").(string)
    subscriptionId := d.Get("subscription_id").(string)
    apiVersion := d.Get("api_version").(string)

    // Build the broker ID (adjust to your SDK)
    id := broker.NewBrokerID(subscriptionId, resourceGroupName, instanceName, brokerName)

    // Call the Azure API
    resp, err := client.Get(ctx, id)
    if err != nil {
        if response.WasNotFound(resp.HttpResponse) {
            d.SetId("")
            return nil
        }
        return fmt.Errorf("error reading broker: %+v", err)
    }

    // Set fields from response
    d.Set("provisioning_state", resp.Properties.ProvisioningState)
    d.Set("properties", flattenBrokerProperties(resp.Properties))

    return nil
}

func resourceBrokerUpdate(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Iot.BrokerClient // Adjust to your client
    ctx := meta.(*clients.Client).StopContext

    brokerName := d.Get("broker_name").(string)
    resourceGroupName := d.Get("resource_group_name").(string)
    instanceName := d.Get("instance_name").(string)
    subscriptionId := d.Get("subscription_id").(string)
    apiVersion := d.Get("api_version").(string)
    properties := d.Get("properties").([]interface{})

    // Build the broker ID and request model (adjust to your SDK)
    id := broker.NewBrokerID(subscriptionId, resourceGroupName, instanceName, brokerName)
    req := broker.BrokerModel{
        Name:       brokerName,
        Properties: flattenBrokerProperties(properties),
    }

    // Call the Azure API
    resp, err := client.Update(ctx, id, req)
    if err != nil {
        return fmt.Errorf("error updating broker: %+v", err)
    }

    d.Set("provisioning_state", resp.Properties.ProvisioningState)
    d.Set("properties", flattenBrokerProperties(resp.Properties))

    return resourceBrokerRead(d, meta)
}

func resourceBrokerDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Iot.BrokerClient
    ctx := meta.(*clients.Client).StopContext

    brokerName := d.Get("broker_name").(string)
    resourceGroupName := d.Get("resource_group_name").(string)
    instanceName := d.Get("instance_name").(string)
    subscriptionId := d.Get("subscription_id").(string)
    apiVersion := d.Get("api_version").(string)

    // Build the broker ID (adjust to your SDK)
    id := broker.NewBrokerID(subscriptionId, resourceGroupName, instanceName, brokerName)

    // Call the Azure API
    err := client.Delete(ctx, id)
    if err != nil {
        return fmt.Errorf("error deleting broker: %+v", err)
    }

    d.SetId("")
    return nil
}

func expandBrokerProperties(input []interface{}) *broker.BrokerProperties {
    if len(input) == 0 || input[0] == nil {
        return nil
    }

    v := input[0].(map[string]interface{})
    
    properties := &broker.BrokerProperties{}
    
    if memoryProfile, ok := v["memory_profile"].(string); ok {
        properties.MemoryProfile = &memoryProfile
    }
    
    return properties
}