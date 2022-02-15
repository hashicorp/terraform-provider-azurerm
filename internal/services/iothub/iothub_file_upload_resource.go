package iothub

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2021-07-02/devices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	msivalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIotHubFileUpload() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubFileUploadCreateUpdate,
		Read:   resourceIotHubFileUploadRead,
		Update: resourceIotHubFileUploadCreateUpdate,
		Delete: resourceIotHubFileUploadDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FileUploadID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"iothub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IotHubID,
			},

			"connection_string": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: fileUploadConnectionStringDiffSuppress,
				Sensitive:        true,
			},

			"container_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"authentication_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(devices.AuthenticationTypeKeyBased),
				ValidateFunc: validation.StringInSlice([]string{
					string(devices.AuthenticationTypeKeyBased),
					string(devices.AuthenticationTypeIdentityBased),
				}, false),
			},

			"identity_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: msivalidate.UserAssignedIdentityID,
			},

			"notifications": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"max_delivery_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"sas_ttl": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT1H",
				ValidateFunc: validate.ISO8601DurationBetween("PT1M", "P1D"),
			},

			"default_ttl": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT1H",
				ValidateFunc: validate.ISO8601DurationBetween("PT1M", "P2D"),
			},

			"lock_duration": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "PT1M",
				ValidateFunc: validate.ISO8601DurationBetween("PT5S", "PT5M"),
			},
		},
	}
}

func resourceIotHubFileUploadCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	subscriptionId := meta.(*clients.Client).IoTHub.ResourceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	iothubId, err := parse.IotHubID(d.Get("iothub_id").(string))
	if err != nil {
		return err
	}
	iotHubName := iothubId.Name
	iotHubRG := iothubId.ResourceGroup

	id := parse.NewFileUploadID(subscriptionId, iotHubRG, iotHubName, "default")

	locks.ByName(iotHubName, IothubResourceName)
	defer locks.UnlockByName(iotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, iotHubRG, iotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iotHubName, iotHubRG)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", iotHubName, iotHubRG, err)
	}

	authenticationType := devices.AuthenticationType(d.Get("authentication_type").(string))
	identityId := d.Get("identity_id").(string)

	storageEndpoints := iothub.Properties.StorageEndpoints
	if storageEndpoints == nil {
		storageEndpoints = make(map[string]*devices.StorageEndpointProperties)
	}

	storageEndpoints["$default"] = &devices.StorageEndpointProperties{
		SasTTLAsIso8601:    utils.String(d.Get("sas_ttl").(string)),
		AuthenticationType: authenticationType,
		ConnectionString:   utils.String(d.Get("connection_string").(string)),
		ContainerName:      utils.String(d.Get("container_name").(string)),
	}

	if identityId != "" {
		if authenticationType != devices.AuthenticationTypeIdentityBased {
			return fmt.Errorf("`identity_id` can only be specified when `authentication_type` is `identityBased`")
		}
		storageEndpoints["$default"].Identity = &devices.ManagedIdentity{
			UserAssignedIdentity: &identityId,
		}
	}

	messagingEndpoints := iothub.Properties.MessagingEndpoints
	if messagingEndpoints == nil {
		messagingEndpoints = make(map[string]*devices.MessagingEndpointProperties)
	}

	messagingEndpoints["fileNotifications"] = &devices.MessagingEndpointProperties{
		LockDurationAsIso8601: utils.String(d.Get("lock_duration").(string)),
		TTLAsIso8601:          utils.String(d.Get("default_ttl").(string)),
		MaxDeliveryCount:      utils.Int32(int32(d.Get("max_delivery_count").(int))),
	}

	iothub.Properties.EnableFileUploadNotifications = utils.Bool(d.Get("notifications").(bool))

	future, err := client.CreateOrUpdate(ctx, iotHubRG, iotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubFileUploadRead(d, meta)
}

func resourceIotHubFileUploadRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FileUploadID(d.Id())
	if err != nil {
		return err
	}

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	iotHubId := parse.NewIotHubID(id.SubscriptionId, id.ResourceGroup, id.IotHubName)
	d.Set("iothub_id", iotHubId.ID())

	if props := iothub.Properties; props != nil {
		connectionString := ""
		containerName := ""
		sasTTL := ""
		authenticationType := string(devices.AuthenticationTypeKeyBased)
		identityId := ""
		if storageEndpointProperties, ok := props.StorageEndpoints["$default"]; ok {
			if v := storageEndpointProperties.ConnectionString; v != nil {
				connectionString = *v
			}

			if v := storageEndpointProperties.ContainerName; v != nil {
				containerName = *v
			}

			if v := storageEndpointProperties.SasTTLAsIso8601; v != nil {
				sasTTL = *v
			}

			if v := string(storageEndpointProperties.AuthenticationType); v != "" {
				authenticationType = v
			}

			if storageEndpointProperties.Identity != nil && storageEndpointProperties.Identity.UserAssignedIdentity != nil {
				identityId = *storageEndpointProperties.Identity.UserAssignedIdentity
			}
		}
		d.Set("connection_string", connectionString)
		d.Set("container_name", containerName)
		d.Set("sas_ttl", sasTTL)
		d.Set("authentication_type", authenticationType)
		d.Set("identity_id", identityId)

		lockDuration := ""
		defaultTTL := ""
		var maxDeliveryCount int32 = 10
		if messagingEndpointProperties, ok := props.MessagingEndpoints["fileNotifications"]; ok {
			if v := messagingEndpointProperties.LockDurationAsIso8601; v != nil {
				lockDuration = *v
			}

			if v := messagingEndpointProperties.TTLAsIso8601; v != nil {
				defaultTTL = *v
			}

			if v := messagingEndpointProperties.MaxDeliveryCount; v != nil {
				maxDeliveryCount = *v
			}
		}
		d.Set("lock_duration", lockDuration)
		d.Set("default_ttl", defaultTTL)
		d.Set("max_delivery_count", maxDeliveryCount)

		enableFileUploadNotifications := false
		if v := props.EnableFileUploadNotifications; v != nil {
			enableFileUploadNotifications = *v
		}
		d.Set("notifications", enableFileUploadNotifications)
	}

	return nil
}

func resourceIotHubFileUploadDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FileUploadID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", id.IotHubName, id.ResourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	if iothub.Properties == nil {
		return nil
	}

	shouldUpdate := false
	if iothub.Properties.StorageEndpoints != nil {
		shouldUpdate = true
		iothub.Properties.StorageEndpoints = make(map[string]*devices.StorageEndpointProperties)
	}

	if iothub.Properties.MessagingEndpoints != nil {
		shouldUpdate = true
		iothub.Properties.MessagingEndpoints["fileNotifications"] = &devices.MessagingEndpointProperties{}
	}

	if iothub.Properties.EnableFileUploadNotifications != nil {
		shouldUpdate = true
		iothub.Properties.EnableFileUploadNotifications = nil
	}

	if !shouldUpdate {
		return nil
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
	}

	return nil
}
