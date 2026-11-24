// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replicated_vm

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryReplicatedItemRead(d *pluginsdk.ResourceData, meta any) error {
	id, err := replicationprotecteditems.ParseReplicationProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationProtectedItemsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery replicated vm %s: %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("making Read request on site recovery replicated vm %s: model is nil", id.String())
	}

	return flattenSiteRecoveryReplicatedItem(d, model)
}

func networkInterfaceResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"source_network_interface_id": {
				Type:         pluginsdk.TypeString,
				Computed:     true, // azignore:AZS007 - pre-existing violation
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"ip_configuration": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true, // azignore:AZS007 - pre-existing violation
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"failover_test_static_ip": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"target_static_ip": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"failover_test_subnet_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"target_subnet_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"failover_test_public_ip_address_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"recovery_load_balancer_backend_address_pool_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
							},
						},

						"recovery_public_ip_address_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"primary": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							// NOTE: O+C - Azure marks the sole IP configuration as primary when omitted.
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func diskEncryptionResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"disk_encryption_key": {
				Type:       pluginsdk.TypeList,
				Required:   true,
				MaxItems:   1,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"secret_url": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret),
						},

						"vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
					},
				},
			},

			"key_encryption_key": {
				Type:       pluginsdk.TypeList,
				Optional:   true,
				MaxItems:   1,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_url": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
						},

						"vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
					},
				},
			},
		},
	}
}

func findNicId(state *replicationprotecteditems.ReplicationProtectedItem, sourceNicId string) *string {
	if a2aDetails, isA2a := state.Properties.ProviderSpecificDetails.(replicationprotecteditems.A2AReplicationDetails); isA2a {
		if a2aDetails.VMNics != nil {
			for _, nic := range *a2aDetails.VMNics {
				if nic.SourceNicArmId != nil && *nic.SourceNicArmId == sourceNicId {
					return nic.NicId
				}
			}
		}
	}
	return nil
}

func resourceSiteRecoveryReplicatedVMDiskHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		if v, ok := m["disk_id"]; ok {
			buf.WriteString(strings.ToLower(v.(string)))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func waitForReplicationToBeHealthy(ctx context.Context, d *pluginsdk.ResourceData, meta any) (*replicationprotecteditems.ReplicationProtectedItem, error) {
	log.Printf("Waiting for Site Recover to replicate VM.")
	stateConf := &pluginsdk.StateChangeConf{
		Target:       []string{"Protected"},
		Refresh:      waitForReplicationToBeHealthyRefreshFunc(d, meta),
		PollInterval: time.Minute,
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("context had no deadline")
	}
	stateConf.Timeout = time.Until(deadline)

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("waiting for site recovery to replicate vm: %+v", err)
	}

	protectedItem, ok := result.(replicationprotecteditems.ReplicationProtectedItem)
	if ok {
		return &protectedItem, nil
	} else {
		return nil, fmt.Errorf("waiting for site recovery return incompatible type")
	}
}

func waitForReplicationToBeHealthyRefreshFunc(d *pluginsdk.ResourceData, meta any) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		id, err := replicationprotecteditems.ParseReplicationProtectedItemID(d.Id())
		if err != nil {
			return nil, "", err
		}

		client := meta.(*clients.Client).RecoveryServices.ReplicationProtectedItemsClient

		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := client.Get(ctx, *id)
		if err != nil {
			return nil, "", fmt.Errorf("making Read request on site recovery replicated vm %s : %+v", id.String(), err)
		}

		if resp.Model == nil {
			return nil, "", fmt.Errorf("missing Model in response when making Read request on site recovery replicated vm %s  %+v", id.String(), err)
		}

		if resp.Model.Properties == nil {
			return nil, "", fmt.Errorf("missing Properties in response when making Read request on site recovery replicated vm %s  %+v", id.String(), err)
		}

		if resp.Model.Properties.ProviderSpecificDetails == nil {
			return nil, "", fmt.Errorf("missing Properties.ProviderSpecificDetails in response when making Read request on site recovery replicated vm %s : %+v", id.String(), err)
		}

		// Find first disk that is not fully replicated yet
		if a2aDetails, isA2a := resp.Model.Properties.ProviderSpecificDetails.(replicationprotecteditems.A2AReplicationDetails); isA2a {
			if a2aDetails.MonitoringPercentageCompletion != nil {
				log.Printf("Waiting for Site Recover to replicate VM, %d%% complete.", *a2aDetails.MonitoringPercentageCompletion)
			}
			if a2aDetails.VMProtectionState != nil {
				return *resp.Model, *a2aDetails.VMProtectionState, nil
			}
		}

		if resp.Model.Properties.ReplicationHealth == nil {
			return nil, "", fmt.Errorf("missing ReplicationHealth in response when making Read request on site recovery replicated vm %s : %+v", id.String(), err)
		}
		return *resp.Model, *resp.Model.Properties.ReplicationHealth, nil
	}
}
