package messages

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/queue/queues"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

var _ StorageQueueMessage = Client{}

func TestLifeCycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	queueName := fmt.Sprintf("queue-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storageaccounts.KindStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	domainSuffix, ok := client.Environment.Storage.DomainSuffix()
	if !ok {
		t.Fatalf("storage didn't return a domain suffix for this environment")
	}
	queuesClient, err := queues.NewWithBaseUri(fmt.Sprintf("https://%s.%s.%s", accountName, "queue", *domainSuffix))
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(queuesClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	messagesClient, err := NewWithBaseUri(fmt.Sprintf("https://%s.%s.%s", accountName, "queue", *domainSuffix))
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(messagesClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	_, err = queuesClient.Create(ctx, queueName, queues.CreateInput{MetaData: map[string]string{}})
	if err != nil {
		t.Fatalf("Error creating queue: %s", err)
	}
	defer queuesClient.Delete(ctx, queueName)

	input := PutInput{
		Message: "ohhai",
	}
	putResp, err := messagesClient.Put(ctx, queueName, input)
	if err != nil {
		t.Fatalf("Error putting message in queue: %s", err)
	}

	messageId := (*putResp.QueueMessages)[0].MessageId
	popReceipt := (*putResp.QueueMessages)[0].PopReceipt

	_, err = messagesClient.Update(ctx, queueName, messageId, UpdateInput{
		PopReceipt:        popReceipt,
		Message:           "Updated message",
		VisibilityTimeout: 65,
	})
	if err != nil {
		t.Fatalf("Error updating: %s", err)
	}

	for i := 0; i < 5; i++ {
		input := PutInput{
			Message: fmt.Sprintf("Message %d", i),
		}
		_, err := messagesClient.Put(ctx, queueName, input)
		if err != nil {
			t.Fatalf("Error putting message %d in queue: %s", i, err)
		}
	}

	peakedMessages, err := messagesClient.Peek(ctx, queueName, PeekInput{NumberOfMessages: 3})
	if err != nil {
		t.Fatalf("Error peaking messages: %s", err)
	}

	for _, v := range *peakedMessages.QueueMessages {
		t.Logf("Message: %q", v.MessageId)
	}

	retrievedMessages, err := messagesClient.Get(ctx, queueName, GetInput{NumberOfMessages: 6})
	if err != nil {
		t.Fatalf("Error retrieving messages: %s", err)
	}

	for _, v := range *retrievedMessages.QueueMessages {
		t.Logf("Message: %q", v.MessageId)

		_, err = messagesClient.Delete(ctx, queueName, v.MessageId, DeleteInput{PopReceipt: v.PopReceipt})
		if err != nil {
			t.Fatalf("Error deleting message from queue: %s", err)
		}
	}
}
