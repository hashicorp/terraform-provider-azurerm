package accounts

type CreateOrUpdateComputePolicyProperties struct {
	MaxDegreeOfParallelismPerJob *int64        `json:"maxDegreeOfParallelismPerJob,omitempty"`
	MinPriorityPerJob            *int64        `json:"minPriorityPerJob,omitempty"`
	ObjectId                     string        `json:"objectId"`
	ObjectType                   AADObjectType `json:"objectType"`
}
