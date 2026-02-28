package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/broker"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthentication"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerauthorization"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/brokerlistener"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflow"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowendpoint"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowprofile"
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/instance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BrokerAuthenticationClient *brokerauthentication.BrokerAuthenticationClient
	BrokerAuthorizationClient  *brokerauthorization.BrokerAuthorizationClient
	BrokerClient               *broker.BrokerClient
	BrokerListenerClient       *brokerlistener.BrokerListenerClient
	DataflowClient             *dataflow.DataflowClient
	DataflowEndpointClient     *dataflowendpoint.DataflowEndpointClient
	DataflowProfileClient      *dataflowprofile.DataflowProfileClient
	InstanceClient             *instance.InstanceClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	brokerAuthenticationClient, err := brokerauthentication.NewBrokerAuthenticationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BrokerAuthentication client: %+v", err)
	}
	o.Configure(brokerAuthenticationClient.Client, o.Authorizers.ResourceManager)

	brokerAuthorizationClient, err := brokerauthorization.NewBrokerAuthorizationClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BrokerAuthorization client: %+v", err)
	}
	o.Configure(brokerAuthorizationClient.Client, o.Authorizers.ResourceManager)

	brokerClient, err := broker.NewBrokerClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Broker client: %+v", err)
	}
	o.Configure(brokerClient.Client, o.Authorizers.ResourceManager)

	brokerListenerClient, err := brokerlistener.NewBrokerListenerClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BrokerListener client: %+v", err)
	}
	o.Configure(brokerListenerClient.Client, o.Authorizers.ResourceManager)

	dataflowClient, err := dataflow.NewDataflowClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Dataflow client: %+v", err)
	}
	o.Configure(dataflowClient.Client, o.Authorizers.ResourceManager)

	dataflowEndpointClient, err := dataflowendpoint.NewDataflowEndpointClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DataflowEndpoint client: %+v", err)
	}
	o.Configure(dataflowEndpointClient.Client, o.Authorizers.ResourceManager)

	dataflowProfileClient, err := dataflowprofile.NewDataflowProfileClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DataflowProfile client: %+v", err)
	}
	o.Configure(dataflowProfileClient.Client, o.Authorizers.ResourceManager)

	instanceClient, err := instance.NewInstanceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Instance client: %+v", err)
	}
	o.Configure(instanceClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		BrokerAuthenticationClient: brokerAuthenticationClient,
		BrokerAuthorizationClient:  brokerAuthorizationClient,
		BrokerClient:               brokerClient,
		BrokerListenerClient:       brokerListenerClient,
		DataflowClient:             dataflowClient,
		DataflowEndpointClient:     dataflowEndpointClient,
		DataflowProfileClient:      dataflowProfileClient,
		InstanceClient:             instanceClient,
	}, nil
}
