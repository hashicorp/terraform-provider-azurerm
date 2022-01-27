package origingroups

type ResponseBasedOriginErrorDetectionParameters struct {
	HttpErrorRanges                          *[]HttpErrorRangeParameters      `json:"httpErrorRanges,omitempty"`
	ResponseBasedDetectedErrorTypes          *ResponseBasedDetectedErrorTypes `json:"responseBasedDetectedErrorTypes,omitempty"`
	ResponseBasedFailoverThresholdPercentage *int64                           `json:"responseBasedFailoverThresholdPercentage,omitempty"`
}
