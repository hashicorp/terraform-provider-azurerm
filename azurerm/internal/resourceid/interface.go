package resourceid

type Formatter interface {
	ID(subscriptionId string) string
}
