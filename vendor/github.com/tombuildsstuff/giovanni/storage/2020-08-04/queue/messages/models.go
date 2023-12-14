package messages

import (
	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type QueueMessage struct {
	MessageText string `xml:"MessageText"`
}

type QueueMessagesListResponse struct {
	HttpResponse *client.Response

	QueueMessages *[]QueueMessageResponse `xml:"QueueMessage"`
}

type QueueMessageResponse struct {
	MessageId       string `xml:"MessageId"`
	InsertionTime   string `xml:"InsertionTime"`
	ExpirationTime  string `xml:"ExpirationTime"`
	PopReceipt      string `xml:"PopReceipt"`
	TimeNextVisible string `xml:"TimeNextVisible"`
}
