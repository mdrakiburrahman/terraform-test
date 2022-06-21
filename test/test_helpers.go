package test

import (
	"errors"
	"os"
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
