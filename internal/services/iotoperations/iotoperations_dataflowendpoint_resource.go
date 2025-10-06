package iotoperations

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

func resourceIotOperationsDataflowEndpoint() *schema.Resource {
    return &schema.Resource{
        Create: resourceIotOperationsDataflowEndpointCreate,
        Read:   resourceIotOperationsDataflowEndpointRead,
        Update: resourceIotOperationsDataflowEndpointUpdate,
        Delete: resourceIotOperationsDataflowEndpointDelete,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "resourceGroupName": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "instanceName": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "subscription_id": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
                ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
                    v := val.(string)
                    if len(v) != 36 {
                        errs = append(errs, fmt.Errorf("%q must be a valid UUID", key))
                    }
                    return
                },
            },
            "api_version": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
                ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
                    v := val.(string)
                    if len(v) < 1 {
                        errs = append(errs, fmt.Errorf("%q must not be empty", key))
                    }
                    return
                },
            },
            "properties": {
                Type:     schema.TypeList,
                Required: true,
                MaxItems: 1,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "endpointType": {
                            Type:     schema.TypeString,
                            Required: true,
                        },
                        "dataExplorerSettings": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "authentication": authSchema(),
                                    "database": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "host": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "batching": batchingSchema(),
                                },
                            },
                        },
                        "dataLakeStorageSettings": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "authentication": authSchemaWithAccessToken(),
                                    "host": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "batching": batchingSchema(),
                                },
                            },
                        },
                        "fabricOneLakeSettings": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "authentication": authSchema(),
                                    "names": {
                                        Type:     schema.TypeList,
                                        Required: true,
                                        MaxItems: 1,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "lakehouse_name": {
                                                    Type:     schema.TypeString,
                                                    Required: true,
                                                },
                                                "workspace_name": {
                                                    Type:     schema.TypeString,
                                                    Required: true,
                                                },
                                            },
                                        },
                                    },
                                    "one_lake_path_type": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "host": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "batching": batchingSchema(),
                                },
                            },
                        },
                        "kafkaSettings": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "authentication": kafkaAuthSchema(),
                                    "consumer_group_id": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "host": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "batching": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        MaxItems: 1,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "mode": {
                                                    Type:     schema.TypeString,
                                                    Required: true,
                                                },
                                                "latency_ms": {
                                                    Type:     schema.TypeInt,
                                                    Required: true,
                                                },
                                                "max_bytes": {
                                                    Type:     schema.TypeInt,
                                                    Required: true,
                                                },
                                                "max_messages": {
                                                    Type:     schema.TypeInt,
                                                    Required: true,
                                                },
                                            },
                                        },
                                    },
                                    "copy_mqtt_properties": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "compression": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "kafka_acks": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "partition_strategy": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "tls": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        MaxItems: 1,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "mode": {
                                                    Type:     schema.TypeString,
                                                    Required: true,
                                                },
                                                "trusted_ca_certificate_config_map_ref": {
                                                    Type:     schema.TypeString,
                                                    Required: true,
                                                },
                                            },
                                        },
                                    },
                                    "cloud_event_attributes": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                },
                            },
                        },
                        "localStorageSettings": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "persistent_volume_claim_ref": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                },
                            },
                        },
                        "mqttSettings": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "authentication": mqttAuthSchema(),
                                    "client_id_prefix": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "host": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "protocol": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "keep_alive_seconds": {
                                        Type:     schema.TypeInt,
                                        Required: true,
                                    },
                                    "retain": {
                                        Type:     schema.TypeString,
                                        Required: true,
                                    },
                                    "max_inflight_messages": {
                                        Type:     schema.TypeInt,
                                        Required: true,
                                    },
                                    "qos": {
                                        Type:     schema.TypeInt,
                                        Required: true,
                                    },
                                    "session_expiry_seconds": {
                                        Type:     schema.TypeInt,
                                        Required: true,
                                    },
                                    "tls": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        MaxItems: 1,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "mode": {
                                                    Type:     schema.TypeString,
                                                    Required: true,
                                                },
                                                "trusted_ca_certificate_config_map_ref": {
                                                    Type:     schema.TypeString,
                                                    Required: true,
                                                },
                                            },
                                        },
                                    },
                                    "cloud_event_attributes": {
                                        Type:     schema.TypeString,
                                        Optional: true,
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
                MaxItems: 1,
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
        },
    }
}

func resourceIotOperationsDataflowEndpointCreate(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client)
    ctx := context.Background()

    // TODO: Build the request payload from d and call the Azure SDK client
    // Example:
    // payload := expandDataflowEndpoint(d)
    // _, err := client.DataflowEndpointClient.CreateOrUpdate(ctx, ...params..., payload)
    // if err != nil {
    //     return fmt.Errorf("error creating DataflowEndpoint: %+v", err)
    // }

    // Set the resource ID
    d.SetId(fmt.Sprintf("%s/%s/%s", d.Get("resource_group_name").(string), d.Get("instance_name").(string), d.Get("name").(string)))
    return resourceIotOperationsDataflowEndpointRead(d, meta)
}

// TODO: Implement expandDataflowEndpoint and other CRUD functions

// Helper schemas for nested authentication and batching
func authSchema() *schema.Schema {
    return &schema.Schema{
        Type:     schema.TypeList,
        Required: true,
        MaxItems: 1,
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "method": {
                    Type:     schema.TypeString,
                    Required: true,
                },
                "system_assigned_managed_identity_settings": {
                    Type:     schema.TypeList,
                    Optional: true,
                    MaxItems: 1,
                    Elem: &schema.Resource{
                        Schema: map[string]*schema.Schema{
                            "audience": {
                                Type:     schema.TypeString,
                                Required: true,
                            },
                        },
                    },
                },
                "user_assigned_managed_identity_settings": {
                    Type:     schema.TypeList,
                    Optional: true,
                    MaxItems: 1,
                    Elem: &schema.Resource{
                        Schema: map[string]*schema.Schema{
                            "client_id": {
                                Type:     schema.TypeString,
                                Required: true,
                            },
                            "scope": {
                                Type:     schema.TypeString,
                                Required: true,
                            },
                            "tenant_id": {
                                Type:     schema.TypeString,
                                Required: true,
                            },
                        },
                    },
                },
            },
        },
    }
}

func authSchemaWithAccessToken() *schema.Schema {
    s := authSchema()
    s.Elem.(*schema.Resource).Schema["access_token_settings"] = &schema.Schema{
        Type:     schema.TypeList,
        Optional: true,
        MaxItems: 1,
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "secret_ref": {
                    Type:     schema.TypeString,
                    Required: true,
                },
            },
        },
    }
    return s
}

func batchingSchema() *schema.Schema {
    return &schema.Schema{
        Type:     schema.TypeList,
        Optional: true,
        MaxItems: 1,
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "latency_seconds": {
                    Type:     schema.TypeInt,
                    Required: true,
                },
                "max_messages": {
                    Type:     schema.TypeInt,
                    Required: true,
                },
            },
        },
    }
}

func kafkaAuthSchema() *schema.Schema {
    s := authSchema()
    s.Elem.(*schema.Resource).Schema["sasl_settings"] = &schema.Schema{
        Type:     schema.TypeList,
        Optional: true,
        MaxItems: 1,
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "sasl_type": {
                    Type:     schema.TypeString,
                    Required: true,
                },
                "secret_ref": {
                    Type:     schema.TypeString,
                    Required: true,
                },
            },
        },
    }
    s.Elem.(*schema.Resource).Schema["x509_certificate_settings"] = &schema.Schema{
        Type:     schema.TypeList,
        Optional: true,
        MaxItems: 1,
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "secret_ref": {
                    Type:     schema.TypeString,
                    Required: true,
                },
            },
        },
    }
    return s
}

func mqttAuthSchema() *schema.Schema {
    s := authSchema()
    s.Elem.(*schema.Resource).Schema["service_account_token_settings"] = &schema.Schema{
        Type:     schema.TypeList,
        Optional: true,
        MaxItems: 1,
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "audience": {
                    Type:     schema.TypeString,
                    Required: true,
                },
            },
        },
    }
    s.Elem.(*schema.Resource).Schema["x509_certificate_settings"] = &schema.Schema{
        Type:     schema.TypeList,
        Optional: true,
        MaxItems: 1,
        Elem: &schema.Resource{
            Schema: map[string]*schema.Schema{
                "secret_ref": {
                    Type:     schema.TypeString,
                    Required: true,
                },
            },
        },
    }
    return s
}
