package iotoperations

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iotoperations/armiotoperations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BrokerListenerResource struct{}

var _ sdk.ResourceWithUpdate = BrokerListenerResource{}

func (r BrokerListenerResource) ModelObject() interface{} {
    return &BrokerListenerModel{}
}

func (r BrokerListenerResource) ResourceType() string {
    return "azurerm_iotoperations_broker_listener"
}

func (r BrokerListenerResource) Arguments() map[string]*schema.Schema {
    return map[string]*schema.Schema{
        "name": {
            Type:     pluginsdk.TypeString,
            Required: true,
            ForceNew: true,
        },
        "resourceGroupName": {
            Type:     pluginsdk.TypeString,
            Required: true,
            ForceNew: true,
        },
        "instanceName": {
            Type:     pluginsdk.TypeString,
            Required: true,
            ForceNew: true,
        },
        "brokerName": {
            Type:     pluginsdk.TypeString,
            Required: true,
            ForceNew: true,
        },
        "location": {
            Type:     pluginsdk.TypeString,
            Optional: true,
        },
        "tags": {
            Type:     pluginsdk.TypeMap,
            Optional: true,
            Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
        },
        "extended_location": {
            Type:     pluginsdk.TypeList,
            Optional: true,
            MaxItems: 1,
            Elem: &pluginsdk.Resource{
                Schema: map[string]*pluginsdk.Schema{
                    "name": {Type: pluginsdk.TypeString, Required: true},
                    "type": {Type: pluginsdk.TypeString, Required: true},
                },
            },
        },
        "properties": {
            Type:     pluginsdk.TypeList,
            Required: true,
            MaxItems: 1,
            Elem: &pluginsdk.Resource{
                Schema: map[string]*pluginsdk.Schema{
                    "service_name": {
                        Type:     pluginsdk.TypeString,
                        Optional: true,
                    },
                    "service_type": {
                        Type:     pluginsdk.TypeString,
                        Optional: true,
                    },
                    "ports": {
                        Type:     pluginsdk.TypeList,
                        Required: true,
                        Elem: &pluginsdk.Resource{
                            Schema: map[string]*pluginsdk.Schema{
                                "port": {
                                    Type:     pluginsdk.TypeInt,
                                    Required: true,
                                },
                                "nodePort": {
                                    Type:     pluginsdk.TypeInt,
                                    Optional: true,
                                },
                                "protocol": {
                                    Type:     pluginsdk.TypeString,
                                    Optional: true,
                                },
                                "authenticationRef": {
                                    Type:     pluginsdk.TypeString,
                                    Optional: true,
								},
                                "authorizationRef": {
                                    Type:     pluginsdk.TypeString,
                                    Optional: true,
                                },
                                "tls": {
                                    Type:     pluginsdk.TypeList,
                                    Optional: true,
                                    MaxItems: 1,
                                    Elem: &pluginsdk.Resource{
                                        Schema: map[string]*pluginsdk.Schema{
                                            "mode": {
                                                Type:     pluginsdk.TypeString,
                                                Required: true,
                                            },
                                            "cert_manager_certificate_spec": {
                                                Type:     pluginsdk.TypeList,
                                                Optional: true,
                                                MaxItems: 1,
                                                Elem: &pluginsdk.Resource{
                                                    Schema: map[string]*pluginsdk.Schema{
                                                        "duration": {
                                                            Type:     pluginsdk.TypeString,
                                                            Optional: true,
                                                        },
                                                        "secret_name": {
                                                            Type:     pluginsdk.TypeString,
                                                            Optional: true,
                                                        },
                                                        "renew_before": {
                                                            Type:     pluginsdk.TypeString,
                                                            Optional: true,
                                                        },
                                                        "issuer_ref": {
                                                            Type:     pluginsdk.TypeList,
                                                            Optional: true,
                                                            MaxItems: 1,
                                                            Elem: &pluginsdk.Resource{
                                                                Schema: map[string]*pluginsdk.Schema{
                                                                    "group": {Type: pluginsdk.TypeString, Optional: true},
                                                                    "kind":  {Type: pluginsdk.TypeString, Optional: true},
                                                                    "name":  {Type: pluginsdk.TypeString, Optional: true},
                                                                },
                                                            },
                                                        },
                                                        "private_key": {
                                                            Type:     pluginsdk.TypeList,
                                                            Optional: true,
                                                            MaxItems: 1,
                                                            Elem: &pluginsdk.Resource{
                                                                Schema: map[string]*pluginsdk.Schema{
                                                                    "algorithm":       {Type: pluginsdk.TypeString, Optional: true},
                                                                    "rotation_policy": {Type: pluginsdk.TypeString, Optional: true},
                                                                },
                                                            },
                                                        },
                                                        "san": {
                                                            Type:     pluginsdk.TypeList,
                                                            Optional: true,
                                                            MaxItems: 1,
                                                            Elem: &pluginsdk.Resource{
                                                                Schema: map[string]*pluginsdk.Schema{
                                                                    "dns": {Type: pluginsdk.TypeList, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString}, Optional: true},
                                                                    "ip":  {Type: pluginsdk.TypeList, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString}, Optional: true},
                                                                },
                                                            },
                                                        },
                                                    },
                                                },
                                            },
                                            "manual": {
                                                Type:     pluginsdk.TypeList,
                                                Optional: true,
                                                MaxItems: 1,
                                                Elem: &pluginsdk.Resource{
                                                    Schema: map[string]*pluginsdk.Schema{
                                                        "secret_ref": {Type: pluginsdk.TypeString, Optional: true},
                                                    },
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    "provisioning_state": {
                        Type:     pluginsdk.TypeString,
                        Computed: true,
                    },
                },
            },
        },
        "subscription_id": {
            Type:     pluginsdk.TypeString,
            Required: true,
            ForceNew: true,
            ValidateFunc: pluginsdk.ValidateFunc(func(val interface{}, key string) (warns []string, errs []error) {
                v := val.(string)
                // Simple UUID format check
                if len(v) != 36 {
                    errs = append(errs, fmt.Errorf("%q must be a valid UUID", key))
                }
                return
            }),
        },
        "api_version": {
            Type:     pluginsdk.TypeString,
            Required: true,
            ForceNew: true,
            ValidateFunc: pluginsdk.ValidateFunc(func(val interface{}, key string) (warns []string, errs []error) {
                v := val.(string)
                if len(v) < 1 {
                    errs = append(errs, fmt.Errorf("%q must not be empty", key))
                }
                return
            }),
        },
    }
}

func (r BrokerListenerResource) Attributes() map[string]*schema.Schema {
    return map[string]*schema.Schema{
        "id": {
            Type:     pluginsdk.TypeString,
            Computed: true,
        },
    }
}

func (r BrokerListenerResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
            client := meta.Client.IotOperations.BrokerListenerClient

            var model map[string]interface{}
            if err := meta.Decode(&model); err != nil {
                return fmt.Errorf("decoding: %+v", err)
            }

            params := expandBrokerListenerResource(model)
            poller, err := client.BeginCreateOrUpdate(
                ctx,
                model["resource_group_name"].(string),
                model["instance_name"].(string),
                model["broker_name"].(string),
                model["name"].(string),
                *params,
                nil,
            )
            if err != nil {
                return fmt.Errorf("creating broker listener: %+v", err)
            }
            resp, err := poller.PollUntilDone(ctx, nil)
            if err != nil {
                return fmt.Errorf("waiting for broker listener create: %+v", err)
            }

            meta.SetID(utils.StringValue(resp.BrokerListenerResource.ID))
            return nil
        },
    }
}

func (r BrokerListenerResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
            client := meta.Client.IotOperations.BrokerListenerClient

            id := meta.ResourceData.Id()
            rg, instance, broker, listener, err := parseBrokerListenerID(id)
            if err != nil {
                return err
            }

            resp, err := client.Get(ctx, rg, instance, broker, listener, nil)
            if err != nil {
                return fmt.Errorf("retrieving broker listener: %+v", err)
            }

            state := flattenBrokerListenerResource(resp.BrokerListenerResource)
            meta.SetID(utils.StringValue(resp.BrokerListenerResource.ID))
            return meta.Encode(state)
        },
    }
}

func (r BrokerListenerResource) Update() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
            return r.Create().Func(ctx, meta)
        },
    }
}

func (r BrokerListenerResource) Delete() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
            client := meta.Client.IotOperations.BrokerListenerClient

            id := meta.ResourceData.Id()
            rg, instance, broker, listener, err := parseBrokerListenerID(id)
            if err != nil {
                return err
            }

            poller, err := client.BeginDelete(ctx, rg, instance, broker, listener, nil)
            if err != nil {
                return fmt.Errorf("deleting broker listener: %+v", err)
            }
            _, err = poller.PollUntilDone(ctx, nil)
            return err
        },
    }
}

// --- Expand/Flatten helpers ---

func expandBrokerListenerResource(d map[string]interface{}) *armiotoperations.BrokerListenerResource {
    // TODO: Implement full expand logic for all nested properties
    return &armiotoperations.BrokerListenerResource{}
}

func flattenBrokerListenerResource(resource *armiotoperations.BrokerListenerResource) map[string]interface{} {
    // TODO: Implement full flatten logic for all nested properties
    return map[string]interface{}{}
}

func parseBrokerListenerID(id string) (rg, instance, broker, listener string, err error) {
    parts := utils.SplitResourceID(id)
    for i := 0; i < len(parts)-1; i++ {
        switch parts[i] {
        case "resourceGroups":
            rg = parts[i+1]
        case "instances":
            instance = parts[i+1]
        case "brokers":
            broker = parts[i+1]
        case "listeners":
            listener = parts[i+1]
        }
    }
    if rg == "" || instance == "" || broker == "" || listener == "" {
        return "", "", "", "", fmt.Errorf("failed to parse broker listener id: %s", id)
    }
    return rg, instance, broker, listener, nil
}