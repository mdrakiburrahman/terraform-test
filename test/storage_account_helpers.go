package test

import (
	"context"
	"os"

	// Azure
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

func storageAccountProperties(ctx context.Context, cred azcore.TokenCredential, resourceGroupName, storageAccountName string) (*armstorage.Account, error) {

	storageAccountClient, err := armstorage.NewAccountsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), cred, nil)
	if err != nil {
		return nil, err
	}

	storageAccountResponse, err := storageAccountClient.GetProperties(
		ctx,
		resourceGroupName,
		storageAccountName,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &storageAccountResponse.Account, nil
}
