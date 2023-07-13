// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/disasterrecoveryconfigs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ServiceBusNamespaceDisasterRecoveryConfigResource struct{}

func resourceServiceBusNamespaceDisasterRecoveryConfig() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceServiceBusNamespaceDisasterRecoveryConfigCreate,
		Read:   resourceServiceBusNamespaceDisasterRecoveryConfigRead,
		Update: resourceServiceBusNamespaceDisasterRecoveryConfigUpdate,
		Delete: resourceServiceBusNamespaceDisasterRecoveryConfigDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := disasterrecoveryconfigs.ParseDisasterRecoveryConfigID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"primary_namespace_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"partner_namespace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty, // nolint: staticcheck
			},

			"alias_authorization_rule_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"primary_connection_string_alias": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string_alias": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceServiceBusNamespaceDisasterRecoveryConfigCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	namespaceId, err := disasterrecoveryconfigs.ParseNamespaceID(d.Get("primary_namespace_id").(string))
	if err != nil {
		return err
	}

	partnerNamespaceId := d.Get("partner_namespace_id").(string)

	id := disasterrecoveryconfigs.NewDisasterRecoveryConfigID(namespaceId.SubscriptionId, namespaceId.ResourceGroupName, namespaceId.NamespaceName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_servicebus_namespace_disaster_recovery_config", id.ID())
		}
	}

	parameters := disasterrecoveryconfigs.ArmDisasterRecovery{
		Properties: &disasterrecoveryconfigs.ArmDisasterRecoveryProperties{
			PartnerNamespace: utils.String(partnerNamespaceId),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, id); err != nil {
		return fmt.Errorf("waiting for replication to complete for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceServiceBusNamespaceDisasterRecoveryConfigRead(d, meta)
}

func resourceServiceBusNamespaceDisasterRecoveryConfigUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disasterrecoveryconfigs.ParseDisasterRecoveryConfigID(d.State().ID)
	if err != nil {
		return err
	}

	locks.ByName(id.NamespaceName, serviceBusNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, serviceBusNamespaceResourceName)

	if d.HasChange("partner_namespace_id") {
		if _, err := client.BreakPairing(ctx, *id); err != nil {
			return fmt.Errorf("breaking the pairing for %s: %+v", *id, err)
		}
		if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
			return fmt.Errorf("waiting for the pairing to break for %s: %+v", *id, err)
		}
	}

	parameters := disasterrecoveryconfigs.ArmDisasterRecovery{
		Properties: &disasterrecoveryconfigs.ArmDisasterRecoveryProperties{
			PartnerNamespace: utils.String(d.Get("partner_namespace_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", *id, err)
	}
	if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for %s to finish replicating: %+v", *id, err)
	}

	return resourceServiceBusNamespaceDisasterRecoveryConfigRead(d, meta)
}

func resourceServiceBusNamespaceDisasterRecoveryConfigRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disasterrecoveryconfigs.ParseDisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	primaryId := disasterrecoveryconfigs.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)

	d.Set("name", id.DisasterRecoveryConfigName)
	d.Set("primary_namespace_id", primaryId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("partner_namespace_id", props.PartnerNamespace)
		}
	}

	// the auth rule cannot be retrieved by dr config name, the shared access policy should either be specified by user or using the default one which is `RootManageSharedAccessKey`
	authRuleId := disasterrecoveryconfigs.NewDisasterRecoveryConfigAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.DisasterRecoveryConfigName, serviceBusNamespaceDefaultAuthorizationRule)
	if input := d.Get("alias_authorization_rule_id").(string); input != "" {
		ruleId, err := disasterrecoveryconfigs.ParseDisasterRecoveryConfigAuthorizationRuleID(input)
		if err != nil {
			return fmt.Errorf("parsing primary namespace auth rule id error: %+v", err)
		}
		authRuleId = *ruleId
	}

	keys, err := client.ListKeys(ctx, authRuleId)

	if err != nil {
		log.Printf("[WARN] listing default keys for %s: %+v", id, err)
	} else {
		if keysModel := keys.Model; keysModel != nil {
			d.Set("primary_connection_string_alias", keysModel.AliasPrimaryConnectionString)
			d.Set("secondary_connection_string_alias", keysModel.AliasSecondaryConnectionString)
			d.Set("default_primary_key", keysModel.PrimaryKey)
			d.Set("default_secondary_key", keysModel.SecondaryKey)
		}
	}

	return nil
}

func resourceServiceBusNamespaceDisasterRecoveryConfigDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disasterrecoveryconfigs.ParseDisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	// @tombuildsstuff: whilst we previously checked the 200 response, since that's the only valid status
	// code defined in the Swagger, anything else would raise an error thus the check is superfluous
	if _, err := client.BreakPairing(ctx, *id); err != nil {
		return fmt.Errorf("breaking pairing %s: %+v", id, err)
	}

	if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for the pairing to break for %s: %+v", *id, err)
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// no future for deletion so wait for it to vanish
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	deleteWait := &pluginsdk.StateChangeConf{
		Pending:    []string{"200"},
		Target:     []string{"404"},
		MinTimeout: 30 * time.Second,
		Timeout:    time.Until(deadline),
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, *id)
			statusCode := "dropped connection"
			if resp.HttpResponse != nil {
				statusCode = strconv.Itoa(resp.HttpResponse.StatusCode)
			}

			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return resp, statusCode, nil
				}
				return nil, "nil", fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return resp, statusCode, nil
		},
	}

	if _, err := deleteWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting the deletion of %s: %v", *id, err)
	}

	namespaceId := disasterrecoveryconfigs.NewNamespaceID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName)
	// it can take some time for the name to become available again
	// this is mainly here 	to enable updating the resource in place
	deadline, ok = ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	nameFreeWait := &pluginsdk.StateChangeConf{
		Pending:    []string{"NameInUse"},
		Target:     []string{"None"},
		MinTimeout: 30 * time.Second,
		Timeout:    time.Until(deadline),
		Refresh: func() (interface{}, string, error) {
			resp, err := client.CheckNameAvailability(ctx, namespaceId, disasterrecoveryconfigs.CheckNameAvailability{
				Name: id.DisasterRecoveryConfigName,
			})
			if err != nil {
				return resp, "Error", fmt.Errorf("checking for the status of %s: %+v", *id, err)
			}

			reason := ""
			if model := resp.Model; model != nil {
				if v := model.Reason; v != nil {
					reason = string(*v)
				}
			}
			return resp, reason, nil
		},
	}

	if _, err := nameFreeWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("checking if the name for %s has become free: %v", *id, err)
	}

	return nil
}

func resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx context.Context, client *disasterrecoveryconfigs.DisasterRecoveryConfigsClient, id disasterrecoveryconfigs.DisasterRecoveryConfigId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(disasterrecoveryconfigs.ProvisioningStateDRAccepted)},
		Target:     []string{string(disasterrecoveryconfigs.ProvisioningStateDRSucceeded)},
		MinTimeout: 30 * time.Second,
		Timeout:    time.Until(deadline),
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, id)
			if err != nil {
				return nil, "error", fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if *props.ProvisioningState == disasterrecoveryconfigs.ProvisioningStateDRFailed {
						return resp, "failed", fmt.Errorf("replication Failed for %s: %+v", id, err)
					}
					return resp, string(*props.ProvisioningState), nil
				}
			}

			return resp, "nil", fmt.Errorf("waiting on replication of %s: %+v", id, err)
		},
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
