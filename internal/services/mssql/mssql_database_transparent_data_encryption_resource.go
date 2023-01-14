package mssql

import (
	"fmt"
	"time"

	// nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlDatabaseTransparentDataEncryption() *pluginsdk.Resource {
	return &pluginsdk.Resource{

		Create: resourceMsSqlDatabaseTransparentDataEncryptionCreateUpdate,
		Read:   resourceMsSqlDatabaseTransparentDataEncryptionRead,
		Update: resourceMsSqlDatabaseTransparentDataEncryptionCreateUpdate,
		Delete: resourceMsSqlDatabaseTransparentDataEncryptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.DatabaseID(id)
			return err
		}, resourceMsSqlDatabaseImporter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"server_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"state": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"database_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
		},
	}
}

func resourceMsSqlDatabaseTransparentDataEncryptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

func resourceMsSqlDatabaseTransparentDataEncryptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	transparentDataEncryptionsClient := meta.(*clients.Client).MSSQL.TransparentDataEncryptionsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseTransparentDataEncryptionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := transparentDataEncryptionsClient.Get(ctx, id.ResourceGroup, id.ServerName, id.DatabaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request for %s: %v", id, err)
	}

	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)
	d.Set("server_id", serverId.ID())

	databaseId := parse.NewDatabaseID(id.SubscriptionId, id.ResourceGroup, id.ServerName, id.DatabaseName)
	d.Set("database_id", databaseId.ID())

	state := ""
	if resp.TransparentDataEncryptionProperties != nil && resp.TransparentDataEncryptionProperties.Status != "" {
		state = string(resp.TransparentDataEncryptionProperties.Status)
	}
	if err := d.Set("state", state); err != nil {
		return fmt.Errorf("setting state`: %+v", err)
	}

	return nil
}

func resourceMsSqlDatabaseTransparentDataEncryptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	return nil
}

type TransparentDataEncryptionId struct {
	SubscriptionId string
	ResourceGroup  string
	ServerName     string
	DatabaseName   string
	Name           string
}

func NewTransparentDataEncryptiontID(subscriptionId, resourceGroup, serverName, database, name string) TransparentDataEncryptionId {
	return TransparentDataEncryptionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServerName:     serverName,
		DatabaseName:   database,
		Name:           name,
	}
}

func parseTransparentDataEncryptionID(input string) (*TransparentDataEncryptionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := TransparentDataEncryptionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if resourceId.DatabaseName, err = id.PopSegment("databases"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("transparentDataEncryption"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
