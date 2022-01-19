package parse

type FlowLogIdShim struct {
	subscriptionId     string
	resourceGroup      string
	networkWatcherName string
	name               string
	nsgId              NetworkSecurityGroupId
}

func NewFlowLogIDShim(subscriptionId, resourceGroup, networkWatcherName, name string, nsgId NetworkSecurityGroupId) FlowLogIdShim {
	return FlowLogIdShim{
		subscriptionId:     subscriptionId,
		resourceGroup:      resourceGroup,
		networkWatcherName: networkWatcherName,
		name:               name,
		nsgId:              nsgId,
	}
}

func (id FlowLogIdShim) Id() interface{} {
	if id.name == "" {
		return NewFlowLogIDLegacy(id.subscriptionId, id.resourceGroup, id.networkWatcherName, id.nsgId)
	}
	return NewFlowLogID(id.subscriptionId, id.resourceGroup, id.networkWatcherName, id.name)
}

func (id FlowLogIdShim) String() string {
	switch id := id.Id().(type) {
	case FlowLogId:
		return id.String()
	case FlowLogIdLegacy:
		return id.String()
	default:
		return ""
	}
}

func (id FlowLogIdShim) ID() string {
	switch id := id.Id().(type) {
	case FlowLogId:
		return id.ID()
	case FlowLogIdLegacy:
		return id.ID()
	default:
		return ""
	}
}

func (id FlowLogIdShim) SubscriptionId() string {
	switch id := id.Id().(type) {
	case FlowLogId:
		return id.SubscriptionId
	case FlowLogIdLegacy:
		return id.SubscriptionId
	default:
		return ""
	}
}

func (id FlowLogIdShim) ResourceGroup() string {
	switch id := id.Id().(type) {
	case FlowLogId:
		return id.ResourceGroup
	case FlowLogIdLegacy:
		return id.ResourceGroupName
	default:
		return ""
	}
}

func (id FlowLogIdShim) NetworkWatcherName() string {
	switch id := id.Id().(type) {
	case FlowLogId:
		return id.NetworkWatcherName
	case FlowLogIdLegacy:
		return id.NetworkWatcherName
	default:
		return ""
	}
}

func (id FlowLogIdShim) Name() string {
	switch id := id.Id().(type) {
	case FlowLogId:
		return id.Name
	case FlowLogIdLegacy:
		return id.Name()
	default:
		return ""
	}
}

func FlowLogIDShim(id string) (*FlowLogIdShim, error) {
	legacyId, err := FlowLogIDLegacy(id)
	if err == nil {
		return &FlowLogIdShim{
			subscriptionId:     legacyId.SubscriptionId,
			resourceGroup:      legacyId.ResourceGroupName,
			networkWatcherName: legacyId.NetworkWatcherName,
			name:               "",
			nsgId:              legacyId.nsgId,
		}, nil
	}
	newId, err := FlowLogID(id)
	if err != nil {
		return nil, err
	}
	return &FlowLogIdShim{
		subscriptionId:     newId.SubscriptionId,
		resourceGroup:      newId.ResourceGroup,
		networkWatcherName: newId.NetworkWatcherName,
		name:               newId.Name,
		nsgId:              NetworkSecurityGroupId{},
	}, nil
}
