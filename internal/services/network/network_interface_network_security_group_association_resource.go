package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetworkInterfaceSecurityGroupAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceSecurityGroupAssociationCreate,
		Read:   resourceNetworkInterfaceSecurityGroupAssociationRead,
		Delete: resourceNetworkInterfaceSecurityGroupAssociationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			splitId := strings.Split(id, "|")
			if _, err := parse.NetworkInterfaceID(splitId[0]); err != nil {
				return err
			}
			if _, err := parse.NetworkSecurityGroupID(splitId[1]); err != nil {
				return err
			}
			return nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"network_interface_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"network_security_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceNetworkInterfaceSecurityGroupAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Network Interface <-> Network Security Group Association creation.")

	networkInterfaceId := d.Get("network_interface_id").(string)
	networkSecurityGroupId := d.Get("network_security_group_id").(string)

	nicId, err := parse.NetworkInterfaceID(networkInterfaceId)
	if err != nil {
		return err
	}

	locks.ByName(nicId.Name, networkInterfaceResourceName)
	defer locks.UnlockByName(nicId.Name, networkInterfaceResourceName)

	nsgId, err := parse.NetworkSecurityGroupID(networkSecurityGroupId)
	if err != nil {
		return err
	}

	locks.ByName(nsgId.Name, networkSecurityGroupResourceName)
	defer locks.UnlockByName(nsgId.Name, networkSecurityGroupResourceName)

	read, err := client.Get(ctx, nicId.ResourceGroup, nicId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("%s was not found!", *nicId)
		}

		return fmt.Errorf("retrieving %s: %+v", *nicId, err)
	}

	props := read.InterfacePropertiesFormat
	if props == nil {
		return fmt.Errorf("Error: `properties` was nil for %s", *nicId)
	}

	// first double-check it doesn't exist
	resourceId := fmt.Sprintf("%s|%s", networkInterfaceId, networkSecurityGroupId)
	if props.NetworkSecurityGroup != nil {
		return tf.ImportAsExistsError("azurerm_network_interface_security_group_association", resourceId)
	}

	props.NetworkSecurityGroup = &network.SecurityGroup{
		ID: utils.String(networkSecurityGroupId),
	}

	future, err := client.CreateOrUpdate(ctx, nicId.ResourceGroup, nicId.Name, read)
	if err != nil {
		return fmt.Errorf("updating Security Group Association for %s: %+v", *nicId, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of Security Group Association for %s: %+v", *nicId, err)
	}

	d.SetId(resourceId)

	return resourceNetworkInterfaceSecurityGroupAssociationRead(d, meta)
}

func resourceNetworkInterfaceSecurityGroupAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	splitId := strings.Split(d.Id(), "|")
	if len(splitId) != 2 {
		return fmt.Errorf("Expected ID to be in the format {networkInterfaceId}|{networkSecurityGroupId} but got %q", d.Id())
	}

	nicID, err := parse.NetworkInterfaceID(splitId[0])
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, nicID.ResourceGroup, nicID.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("%s was not found - removing from state!", *nicID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *nicID, err)
	}

	props := read.InterfacePropertiesFormat
	if props == nil {
		return fmt.Errorf("Error: `properties` was nil for %s", *nicID)
	}

	if props.NetworkSecurityGroup == nil || props.NetworkSecurityGroup.ID == nil {
		log.Printf("%s doesn't have a Security Group attached - removing from state!", *nicID)
		d.SetId("")
		return nil
	}

	d.Set("network_interface_id", read.ID)

	// nil-checked above
	d.Set("network_security_group_id", props.NetworkSecurityGroup.ID)

	return nil
}

func resourceNetworkInterfaceSecurityGroupAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	splitId := strings.Split(d.Id(), "|")
	if len(splitId) != 2 {
		return fmt.Errorf("Expected ID to be in the format {networkInterfaceId}/{networkSecurityGroup} but got %q", d.Id())
	}

	nicID, err := parse.NetworkInterfaceID(splitId[0])
	if err != nil {
		return err
	}

	locks.ByName(nicID.Name, networkInterfaceResourceName)
	defer locks.UnlockByName(nicID.Name, networkInterfaceResourceName)

	read, err := client.Get(ctx, nicID.ResourceGroup, nicID.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf(" %s was not found!", *nicID)
		}

		return fmt.Errorf("retrieving %s: %+v", *nicID, err)
	}

	props := read.InterfacePropertiesFormat
	if props == nil {
		return fmt.Errorf("Error: `properties` was nil for %s", *nicID)
	}

	props.NetworkSecurityGroup = nil
	read.InterfacePropertiesFormat = props

	future, err := azuresdkhacks.UpdateNetworkInterfaceAllowingRemovalOfNSG(ctx, client, nicID.ResourceGroup, nicID.Name, read)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *nicID, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *nicID, err)
	}

	return nil
}
