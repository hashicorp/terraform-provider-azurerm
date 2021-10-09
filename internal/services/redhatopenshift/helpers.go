package redhatopenshift

import (
	"fmt"
	"math/rand"
	"time"
)

func randomString(acceptedChars string, size int) string {
	charSet := []rune(acceptedChars)
	randomChars := make([]rune, size)

	rand.Seed(time.Now().UnixNano())

	for i := range randomChars {
		randomChars[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(randomChars)
}

func GenerateRandomDomainName() string {
	randomPrefix := randomString("abcdefghijklmnopqrstuvwxyz", 1)
	randomName := randomString("abcdefghijklmnopqrstuvwxyz1234567890", 7)

	return fmt.Sprintf("%s%s", randomPrefix, randomName)
}

func ResourceGroupID(subscriptionId string, resourceGroupName string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s"
	return fmt.Sprintf(fmtString, subscriptionId, resourceGroupName)
}
