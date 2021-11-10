package accounts

type UpdateComputePolicyProperties struct {
	MaxDegreeOfParallelismPerJob *int64         `json:"maxDegreeOfParallelismPerJob,omitempty"`
	MinPriorityPerJob            *int64         `json:"minPriorityPerJob,omitempty"`
	ObjectId                     *string        `json:"objectId,omitempty"`
	ObjectType                   *AADObjectType `json:"objectType,omitempty"`
}
