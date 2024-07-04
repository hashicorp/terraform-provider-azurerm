package clients

// NOTE: this file is generated - manual changes will be overwritten.

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	chaosstudio "github.com/hashicorp/terraform-provider-azurerm/internal/services/chaosstudio/client"
	containers "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
	devcenter "github.com/hashicorp/terraform-provider-azurerm/internal/services/devcenter/client"
	managedidentity "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedidentity/client"
)

type autoClient struct {
	ChaosStudio      *chaosstudio.AutoClient
	ContainerService *containers.AutoClient
	DevCenter        *devcenter.AutoClient
	ManagedIdentity  *managedidentity.AutoClient
}

func buildAutoClients(client *autoClient, o *common.ClientOptions) (err error) {

	if client.ChaosStudio, err = chaosstudio.NewClient(o); err != nil {
		return fmt.Errorf("building client for ChaosStudio: %+v", err)
	}

	if client.ContainerService, err = containers.NewClient(o); err != nil {
		return fmt.Errorf("building client for ContainerService: %+v", err)
	}

	if client.DevCenter, err = devcenter.NewClient(o); err != nil {
		return fmt.Errorf("building client for DevCenter: %+v", err)
	}

	if client.ManagedIdentity, err = managedidentity.NewClient(o); err != nil {
		return fmt.Errorf("building client for ManagedIdentity: %+v", err)
	}

	return nil
}
