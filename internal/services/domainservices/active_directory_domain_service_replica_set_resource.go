package domainservices

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/domainservices/mgmt/2020-01-01/aad"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceActiveDirectoryDomainServiceReplicaSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceActiveDirectoryDomainServiceReplicaSetCreate,
		Read:   resourceActiveDirectoryDomainServiceReplicaSetRead,
		Delete: resourceActiveDirectoryDomainServiceReplicaSetDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(3 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(2 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(1 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DomainServiceReplicaSetID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"domain_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DomainServiceID,
			},

			"location": azure.SchemaLocation(),

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.SubnetID,
			},

			"domain_controller_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"external_access_ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceActiveDirectoryDomainServiceReplicaSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	domainServiceId, err := parse.DomainServiceID(d.Get("domain_service_id").(string))
	if err != nil {
		return err
	}
	if domainServiceId == nil {
		return fmt.Errorf("parsing ID for Domain Service Replica Set")
	}

	locks.ByName(domainServiceId.Name, DomainServiceResourceName)
	defer locks.UnlockByName(domainServiceId.Name, DomainServiceResourceName)

	domainService, err := client.Get(ctx, domainServiceId.ResourceGroup, domainServiceId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(domainService.Response) {
			return fmt.Errorf("could not find %s: %s", domainServiceId, err)
		}
		return fmt.Errorf("reading %s: %s", domainServiceId, err)
	}

	if domainService.DomainServiceProperties.ReplicaSets == nil || len(*domainService.DomainServiceProperties.ReplicaSets) == 0 {
		return fmt.Errorf("reading %s: returned with missing replica set information, expected at least 1 replica set: %s", domainServiceId, err)
	}

	subnetId := d.Get("subnet_id").(string)
	replicaSets := *domainService.DomainServiceProperties.ReplicaSets

	for _, r := range replicaSets {
		if r.ReplicaSetID == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing ReplicaSetID", domainServiceId)
		}
		if r.SubnetID == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing SubnetID", domainServiceId)
		}

		// We assume that two replica sets cannot coexist in the same subnet
		if strings.EqualFold(subnetId, *r.SubnetID) {
			// Generate an ID here since we only know it once we know the ReplicaSetID
			id := parse.NewDomainServiceReplicaSetID(domainServiceId.SubscriptionId, domainServiceId.ResourceGroup, domainServiceId.Name, *r.ReplicaSetID)
			return tf.ImportAsExistsError("azurerm_active_directory_domain_service_replica_set", id.ID())
		}
	}

	loc := location.Normalize(d.Get("location").(string))
	replicaSets = append(replicaSets, aad.ReplicaSet{
		Location: utils.String(loc),
		SubnetID: utils.String(subnetId),
	})

	domainService.DomainServiceProperties.ReplicaSets = &replicaSets

	future, err := client.CreateOrUpdate(ctx, domainServiceId.ResourceGroup, domainServiceId.Name, domainService)
	if err != nil {
		return fmt.Errorf("creating/updating Replica Sets for %s: %+v", domainServiceId, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for Replica Sets for %s: %+v", domainServiceId, err)
	}

	// We need to retrieve the domain service again to find out the new replica set ID
	domainService, err = client.Get(ctx, domainServiceId.ResourceGroup, domainServiceId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(domainService.Response) {
			return fmt.Errorf("could not find %s: %s", domainServiceId, err)
		}
		return fmt.Errorf("reading %s: %s", domainServiceId, err)
	}

	if domainService.DomainServiceProperties.ReplicaSets == nil || len(*domainService.DomainServiceProperties.ReplicaSets) == 0 {
		return fmt.Errorf("reading %s: returned with missing replica set information, expected at least 1 replica set: %s", domainServiceId, err)
	}

	var id parse.DomainServiceReplicaSetId
	// Assuming that two replica sets cannot coexist in the same subnet, we identify our new replica set by its SubnetID
	for _, r := range *domainService.DomainServiceProperties.ReplicaSets {
		if r.ReplicaSetID == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing ReplicaSetID", domainServiceId)
		}
		if r.SubnetID == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing SubnetID", domainServiceId)
		}

		if strings.EqualFold(subnetId, *r.SubnetID) {
			// We found it!
			id = parse.NewDomainServiceReplicaSetID(domainServiceId.SubscriptionId, domainServiceId.ResourceGroup, domainServiceId.Name, *r.ReplicaSetID)
		}
	}

	if id.ReplicaSetName == "" {
		return fmt.Errorf("reading %s: the new replica set was not returned", domainServiceId)
	}

	// Wait for all replica sets to become available with two domain controllers each before proceeding
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      domainServiceControllerRefreshFunc(ctx, client, *domainServiceId, false),
		Delay:        1 * time.Minute,
		PollInterval: 1 * time.Minute,
		Timeout:      time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for both domain controllers to become available in all replica sets for %s: %+v", domainServiceId, err)
	}

	d.SetId(id.ID())

	return resourceActiveDirectoryDomainServiceReplicaSetRead(d, meta)
}

func resourceActiveDirectoryDomainServiceReplicaSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainServiceReplicaSetID(d.Id())
	if err != nil {
		return err
	}

	domainService, err := client.Get(ctx, id.ResourceGroup, id.DomainServiceName)
	if err != nil {
		if utils.ResponseWasNotFound(domainService.Response) {
			d.SetId("")
			return nil
		}
		return err
	}

	if domainService.DomainServiceProperties.ReplicaSets == nil || len(*domainService.DomainServiceProperties.ReplicaSets) == 0 {
		return fmt.Errorf("reading %s: domain service returned with missing replica set information, expected at least 1 replica set: %s", id, err)
	}

	var (
		domainControllerIpAddresses []string
		externalAccessIpAddress     string
		loc                         string
		serviceStatus               string
		subnetId                    string
	)

	replicaSets := *domainService.DomainServiceProperties.ReplicaSets

	for _, r := range replicaSets {
		if r.ReplicaSetID == nil {
			return fmt.Errorf("reading %s: a replica set was returned with a missing ReplicaSetID", id)
		}

		// ReplicaSetName in the ID struct is really the replica set ID
		if *r.ReplicaSetID == id.ReplicaSetName {
			if r.DomainControllerIPAddress != nil {
				domainControllerIpAddresses = *r.DomainControllerIPAddress
			}
			if r.ExternalAccessIPAddress != nil {
				externalAccessIpAddress = *r.ExternalAccessIPAddress
			}
			if r.Location != nil {
				loc = location.NormalizeNilable(r.Location)
			}
			if r.ServiceStatus != nil {
				serviceStatus = *r.ServiceStatus
			}
			if r.SubnetID != nil {
				subnetId = *r.SubnetID
			}
		}
	}

	d.Set("domain_controller_ip_addresses", domainControllerIpAddresses)
	d.Set("external_access_ip_address", externalAccessIpAddress)
	d.Set("location", loc)
	d.Set("service_status", serviceStatus)
	d.Set("subnet_id", subnetId)

	return nil
}

func resourceActiveDirectoryDomainServiceReplicaSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainServiceReplicaSetID(d.Id())
	if err != nil {
		return err
	}

	domainService, err := client.Get(ctx, id.ResourceGroup, id.DomainServiceName)
	if err != nil {
		if utils.ResponseWasNotFound(domainService.Response) {
			return fmt.Errorf("deleting %s: domain service was not found: %s", id, err)
		}
		return err
	}

	if domainService.DomainServiceProperties.ReplicaSets == nil || len(*domainService.DomainServiceProperties.ReplicaSets) == 0 {
		return fmt.Errorf("deleting %s: domain service returned with missing replica set information, expected at least 1 replica set: %s", id, err)
	}

	replicaSets := *domainService.DomainServiceProperties.ReplicaSets

	newReplicaSets := make([]aad.ReplicaSet, 0)
	for _, r := range replicaSets {
		if r.ReplicaSetID == nil {
			return fmt.Errorf("deleting %s: a replica set was returned with a missing ReplicaSetID", id)
		}

		if *r.ReplicaSetID == id.ReplicaSetName {
			continue
		}

		newReplicaSets = append(newReplicaSets, r)
	}

	if len(replicaSets) == len(newReplicaSets) {
		return fmt.Errorf("deleting %s: could not determine which replica set to remove", id)
	}

	domainService.DomainServiceProperties.ReplicaSets = &newReplicaSets

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.DomainServiceName, domainService)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	// Wait for all replica sets to become available with two domain controllers each before proceeding
	// Generate a partial DomainServiceId since we don't need to know the initial replica set ID here
	domainServiceId := parse.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName, "")
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      domainServiceControllerRefreshFunc(ctx, client, domainServiceId, true),
		Delay:        1 * time.Minute,
		PollInterval: 1 * time.Minute,
		Timeout:      time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for replica sets to finish updating for %s: %+v", domainServiceId, err)
	}

	return nil
}
