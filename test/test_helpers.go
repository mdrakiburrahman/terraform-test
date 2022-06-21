package test

import (
	"errors"
	"log"
	"os"

	// Azure
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// Injects environment variables into map for Terraform authentication with Azure
func setTerraformVariables() (map[string]string, error) {

	// Grab from devcontainer environment variables
	ARM_CLIENT_ID := os.Getenv("spnClientId")
	ARM_CLIENT_SECRET := os.Getenv("spnClientSecret")
	ARM_TENANT_ID := os.Getenv("spnTenantId")
	ARM_SUBSCRIPTION_ID := os.Getenv("subscriptionId")

	// If any of the above variables are empty, return an error
	if ARM_CLIENT_ID == "" || ARM_CLIENT_SECRET == "" || ARM_TENANT_ID == "" || ARM_SUBSCRIPTION_ID == "" {
		return nil, errors.New("missing one or more of the following environment variables: spnClientId, spnClientSecret, spnTenantId, subscriptionId")
	}

	// Creating globalEnVars for terraform call through Terratest
	EnvVars := make(map[string]string)

	if ARM_CLIENT_ID != "" {
		EnvVars["ARM_CLIENT_ID"] = ARM_CLIENT_ID
		EnvVars["ARM_CLIENT_SECRET"] = ARM_CLIENT_SECRET
		EnvVars["ARM_TENANT_ID"] = ARM_TENANT_ID
		EnvVars["ARM_SUBSCRIPTION_ID"] = ARM_SUBSCRIPTION_ID
	}

	return EnvVars, nil
}

// Injects environment variables into map for Azure SDK authentication with Azure
// https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication?tabs=bash
func setARMVariables() error {

	// Set environment variables for Azure SDK authentication
	os.Setenv("AZURE_CLIENT_ID", os.Getenv("spnClientId"))
	os.Setenv("AZURE_CLIENT_SECRET", os.Getenv("spnClientSecret"))
	os.Setenv("AZURE_TENANT_ID", os.Getenv("spnTenantId"))
	os.Setenv("AZURE_SUBSCRIPTION_ID", os.Getenv("subscriptionId"))

	// If any of the above variables are empty, return an error
	if os.Getenv("AZURE_CLIENT_ID") == "" || os.Getenv("AZURE_CLIENT_SECRET") == "" || os.Getenv("AZURE_TENANT_ID") == "" || os.Getenv("AZURE_SUBSCRIPTION_ID") == "" {
		return errors.New("missing one or more of the following environment variables: spnClientId, spnClientSecret, spnTenantId, subscriptionId")
	}

	return nil
}

// Authenticates to Azure and initiates context
func getAzureCred() (azcore.TokenCredential, error) {

	// Grabs Azure SDK authentication environment variables, errors if any are missing
	err := setARMVariables()

	if err != nil {
		return nil, err
	}

	// Authenticates using Environment variables
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		log.Fatal(err)
	}

	return cred, nil
}
