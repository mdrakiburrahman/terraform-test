package test

import (
	// Native
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	// Terragrunt
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"

	// Testing
	"github.com/stretchr/testify/assert"
)

// Global variable declaration block
var (
	globalBackendConf = make(map[string]interface{})
	globalEnvVars     = make(map[string]string)
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

func TestStorageAccountExample(t *testing.T) {

	t.Parallel()

	// Grabs Service Principal environment variables
	setTerraformVariables()

	// Set input values for the test
	inputUniquePostfix := strings.ToLower(random.UniqueId())
	inputResourceGroupName := fmt.Sprintf("%s-%s", "terratest-storage-account", inputUniquePostfix)
	inputLocation := "canadacentral"
	inputTags := map[string]string{
		"Source":  "terratest",
		"Owner":   "Raki Rahman",
		"Project": "Terraform CI testing",
	}
	inputAccountTier := "Standard"
	inputAccountKind := "StorageV2"

	// Use Terratest to call Terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{

		// Set the path to the Terraform code that will be tested.
		TerraformDir: "../examples/storage-account",

		// Variables to pass to our Terraform code using -var options.
		Vars: map[string]interface{}{
			"resource_group_name":    inputResourceGroupName,
			"location":               inputLocation,
			"tags":                   inputTags,
			"account_tier":           inputAccountTier,
			"account_kind":           inputAccountKind,
			"account_unique_postfix": inputUniquePostfix,
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
	assert.Equal(t, strings.ToLower(inputLocation), strings.ToLower(outputLocation))

	// Run tests
	t.Run("storage_account_location_matched", func(t *testing.T) {
		assert.Equal(t, strings.ToLower(inputLocation), strings.ToLower(outputLocation), "Storage Account Location matched input")
	})

	t.Run("test_1", func(t *testing.T) {
		assert.Equal(t, "one", "one", "Should always pass")
	})

	t.Run("test_2", func(t *testing.T) {
		assert.Equal(t, "one", "two", "Should always fail")
	})
}
