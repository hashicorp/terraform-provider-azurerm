package netapp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2021-10-01/netapp"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetAppPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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
			_, err := parse.CapacityPoolID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PoolName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

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
					string(netapp.ServiceLevelPremium),
					string(netapp.ServiceLevelStandard),
					string(netapp.ServiceLevelUltra),
				}, false),
			},

			"size_in_tb": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(4, 500),
			},

			"qos_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(netapp.QosTypeAuto),
					string(netapp.QosTypeManual),
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceNetAppPoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewCapacityPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_netapp_pool", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	serviceLevel := d.Get("service_level").(string)
	sizeInTB := int64(d.Get("size_in_tb").(int))
	sizeInMB := sizeInTB * 1024 * 1024
	sizeInBytes := sizeInMB * 1024 * 1024

	capacityPoolParameters := netapp.CapacityPool{
		Location: utils.String(location),
		PoolProperties: &netapp.PoolProperties{
			ServiceLevel: netapp.ServiceLevel(serviceLevel),
			Size:         utils.Int64(sizeInBytes),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if qosType, ok := d.GetOk("qos_type"); ok {
		capacityPoolParameters.PoolProperties.QosType = netapp.QosType(qosType.(string))
	}

	future, err := client.CreateOrUpdate(ctx, capacityPoolParameters, id.ResourceGroup, id.NetAppAccountName, id.Name)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	// Wait for pool to complete update
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

	id, err := parse.CapacityPoolID(d.Id())
	if err != nil {
		return err
	}

	shouldUpdate := false
	update := netapp.CapacityPoolPatch{
		PoolPatchProperties: &netapp.PoolPatchProperties{},
	}

	if d.HasChange("size_in_tb") {
		shouldUpdate = true

		sizeInTB := int64(d.Get("size_in_tb").(int))
		sizeInMB := sizeInTB * 1024 * 1024
		sizeInBytes := sizeInMB * 1024 * 1024

		update.PoolPatchProperties.Size = utils.Int64(sizeInBytes)
	}

	if d.HasChange("qos_type") {
		shouldUpdate = true
		qosType := d.Get("qos_type")
		update.PoolPatchProperties.QosType = netapp.QosType(qosType.(string))
	}

	if d.HasChange("tags") {
		shouldUpdate = true
		tagsRaw := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(tagsRaw)
	}

	if shouldUpdate {
		future, err := client.Update(ctx, update, id.ResourceGroup, id.NetAppAccountName, id.Name)
		if err != nil {
			return fmt.Errorf("updating Capacity Pool %q: %+v", id.Name, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for the update of %s: %+v", id, err)
		}

		// Wait for pool to complete update
		if err := waitForPoolCreateOrUpdate(ctx, client, *id); err != nil {
			return err
		}
	}

	return resourceNetAppPoolRead(d, meta)
}

func resourceNetAppPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CapacityPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.NetAppAccountName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if poolProperties := resp.PoolProperties; poolProperties != nil {
		d.Set("service_level", poolProperties.ServiceLevel)

		sizeInTB := int64(0)
		if poolProperties.Size != nil {
			sizeInBytes := *poolProperties.Size
			sizeInMB := sizeInBytes / 1024 / 1024
			sizeInTB = sizeInMB / 1024 / 1024
		}
		d.Set("size_in_tb", int(sizeInTB))
		d.Set("qos_type", string(poolProperties.QosType))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceNetAppPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CapacityPoolID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.NetAppAccountName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
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

func netappPoolDeleteStateRefreshFunc(ctx context.Context, client *netapp.PoolsClient, id parse.CapacityPoolId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
			}
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func waitForPoolCreateOrUpdate(ctx context.Context, client *netapp.PoolsClient, id parse.CapacityPoolId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
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

func netappPoolStateRefreshFunc(ctx context.Context, client *netapp.PoolsClient, id parse.CapacityPoolId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("retrieving NetApp Capacity Pool %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
			}
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
