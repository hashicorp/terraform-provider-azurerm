// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/capacitypools"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetAppPool() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceNetAppPoolCreate,
		Read:   resourceNetAppPoolRead,
		Update: resourceNetAppPoolUpdate,
		Delete: resourceNetAppPoolDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := capacitypools.ParseCapacityPoolID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PoolName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			"service_level": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(capacitypools.ServiceLevelPremium),
					string(capacitypools.ServiceLevelStandard),
					string(capacitypools.ServiceLevelUltra),
				}, false),
			},

			"size_in_tb": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 2048),
			},

			"qos_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(capacitypools.QosTypeAuto),
				ValidateFunc: validation.StringInSlice([]string{
					string(capacitypools.QosTypeAuto),
					string(capacitypools.QosTypeManual),
				}, false),
			},

			"encryption_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  capacitypools.EncryptionTypeSingle,
				ValidateFunc: validation.StringInSlice([]string{
					string(capacitypools.EncryptionTypeSingle),
					string(capacitypools.EncryptionTypeDouble),
				}, false),
			},

			"cool_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			// `cool_access_enabled` cannot be disabled
			pluginsdk.ForceNewIfChange("cool_access_enabled", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(bool) && !new.(bool)
			}),
		),
	}

	return resource
}

func resourceNetAppPoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := capacitypools.NewCapacityPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.PoolsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_netapp_pool", id.ID())
		}
	}

	sizeInTB := int64(d.Get("size_in_tb").(int))
	sizeInMB := sizeInTB * 1024 * 1024
	sizeInBytes := sizeInMB * 1024 * 1024

	encryptionType := capacitypools.EncryptionType(d.Get("encryption_type").(string))

	capacityPoolParameters := capacitypools.CapacityPool{
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Properties: capacitypools.PoolProperties{
			ServiceLevel:   capacitypools.ServiceLevel(d.Get("service_level").(string)),
			Size:           sizeInBytes,
			EncryptionType: &encryptionType,
			CoolAccess:     pointer.To(d.Get("cool_access_enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if qosType, ok := d.GetOk("qos_type"); ok {
		qos := capacitypools.QosType(qosType.(string))
		capacityPoolParameters.Properties.QosType = &qos
	}

	if err := client.PoolsCreateOrUpdateThenPoll(ctx, id, capacityPoolParameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Wait for pool to complete create
	if err := waitForPoolCreateOrUpdate(ctx, client, id); err != nil {
		return err
	}

	d.SetId(id.ID())
	return resourceNetAppPoolRead(d, meta)
}

func resourceNetAppPoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacitypools.ParseCapacityPoolID(d.Id())
	if err != nil {
		return err
	}

	update := capacitypools.CapacityPoolPatch{
		Properties: &capacitypools.PoolPatchProperties{},
	}

	if d.HasChange("size_in_tb") {
		sizeInTB := int64(d.Get("size_in_tb").(int))
		sizeInMB := sizeInTB * 1024 * 1024
		sizeInBytes := sizeInMB * 1024 * 1024

		update.Properties.Size = utils.Int64(sizeInBytes)
	}

	if d.HasChange("qos_type") {
		qosType := capacitypools.QosType(d.Get("qos_type").(string))
		update.Properties.QosType = &qosType
	}

	if d.HasChange("cool_access_enabled") {
		update.Properties.CoolAccess = pointer.To(d.Get("cool_access_enabled").(bool))
	}

	if d.HasChange("tags") {
		tagsRaw := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(tagsRaw)
	}

	if err = client.PoolsUpdateThenPoll(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", id.ID(), err)
	}

	// Wait for pool to complete update
	if err = waitForPoolCreateOrUpdate(ctx, client, *id); err != nil {
		return err
	}

	return resourceNetAppPoolRead(d, meta)
}

func resourceNetAppPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacitypools.ParseCapacityPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.PoolsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CapacityPoolName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.NetAppAccountName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		poolProperties := model.Properties
		d.Set("service_level", poolProperties.ServiceLevel)

		sizeInBytes := poolProperties.Size
		sizeInMB := sizeInBytes / 1024 / 1024
		sizeInTB := sizeInMB / 1024 / 1024
		d.Set("size_in_tb", int(sizeInTB))
		qosType := ""
		if poolProperties.QosType != nil {
			qosType = string(*poolProperties.QosType)
		}
		d.Set("qos_type", qosType)
		d.Set("encryption_type", string(pointer.From(poolProperties.EncryptionType)))
		d.Set("cool_access_enabled", pointer.From(poolProperties.CoolAccess))

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceNetAppPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacitypools.ParseCapacityPoolID(d.Id())
	if err != nil {
		return err
	}

	if err := client.PoolsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// The resource NetApp Pool depends on the resource NetApp Account.
	// Although the delete API returns 404 which means the NetApp Pool resource has been deleted.
	// Then it tries to immediately delete NetApp Account but it still throws error `Can not delete resource before nested resources are deleted.`
	// In this case we're going to re-check status code again.
	// For more details, see related Bug: https://github.com/Azure/azure-sdk-for-go/issues/6374
	log.Printf("[DEBUG] Waiting for %s to be deleted", *id)
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context has no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"204", "404"},
		Refresh:                   netappPoolDeleteStateRefreshFunc(ctx, client, *id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
	}

	return nil
}

func netappPoolDeleteStateRefreshFunc(ctx context.Context, client *capacitypools.CapacityPoolsClient, id capacitypools.CapacityPoolId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.PoolsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
			}
		}

		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = strconv.Itoa(res.HttpResponse.StatusCode)
		}
		return res, statusCode, nil
	}
}

func waitForPoolCreateOrUpdate(ctx context.Context, client *capacitypools.CapacityPoolsClient, id capacitypools.CapacityPoolId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404"},
		Target:                    []string{"200", "202"},
		Refresh:                   netappPoolStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
	}

	return nil
}

func netappPoolStateRefreshFunc(ctx context.Context, client *capacitypools.CapacityPoolsClient, id capacitypools.CapacityPoolId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.PoolsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving %s: %s", id.ID(), err)
			}
		}

		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = strconv.Itoa(res.HttpResponse.StatusCode)
		}
		return res, statusCode, nil
	}
}
