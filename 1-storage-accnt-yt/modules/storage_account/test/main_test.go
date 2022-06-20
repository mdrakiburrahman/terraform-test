package main

import (
	// Native
	"os"
	"strings"
	"testing"
	"errors"

	// Terragrunt
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/random"

	// Testing
	"github.com/stretchr/testify/assert"
)

// Global variable declaration block
var (

	resourceGroup 	  = "terratest-hackathon-1-youtube-rg"
	globalBackendConf = make(map[string]interface{})
	globalEnvVars     = make(map[string]string)
	uniquePostfix     = strings.ToLower(random.UniqueId())

)

const (

	apiVersion              = "2019-06-01"
	resourceProvisionStatus = "Succeeded"

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
		return nil, errors.New("Missing one or more of the following environment variables: spnClientId, spnClientSecret, spnTenantId, subscriptionId")
	}

	// Creating globalEnVars for terraform call through Terratest
	if ARM_CLIENT_ID != "" {
		globalEnvVars["ARM_CLIENT_ID"] = ARM_CLIENT_ID
		globalEnvVars["ARM_CLIENT_SECRET"] = ARM_CLIENT_SECRET
		globalEnvVars["ARM_TENANT_ID"] = ARM_TENANT_ID
		globalEnvVars["ARM_SUBSCRIPTION_ID"] = ARM_SUBSCRIPTION_ID
	}

	return globalEnvVars, nil
}

func TestTerraform_storage_account(t *testing.T) {

	t.Parallel()

	// Grabs Service Principal environment variables
	setTerraformVariables()

	// Sets expected values for the test
	expectedLocation := "eastus"
	expectedStorageAccountTier := "Standard"
	expectedStorageAccountKind := "StorageV2"

	// Use Terratest to call Terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		
		// Set the path to the Terraform code that will be tested.
		TerraformDir: "../provision",

		// Variables to pass to our Terraform code using -var options.
		Vars: map[string]interface{}{
			"resource_group_name": 	resourceGroup,
			"location":         	expectedLocation,   
			"account_tier":         expectedStorageAccountTier,
			"unique_postfix":       uniquePostfix,
			"account_kind":         expectedStorageAccountKind,
		},

		// Service Principal creds from Environment Variables
		EnvVars: globalEnvVars,

		// State backend - we start empty
		BackendConfig: globalBackendConf,

		// Colors in Terraform commands
		NoColor: false,

	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	outputLocation := terraform.Output(t, terraformOptions, "primary_location")

	// Convert outputLocation to lower and compare with expectedLocation
	assert.Equal(t, strings.ToLower(expectedLocation), strings.ToLower(outputLocation))
}