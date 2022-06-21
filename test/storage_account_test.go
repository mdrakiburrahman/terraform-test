package test

import (
	// Native

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
)

func TestStorageAccountExample(t *testing.T) {

	t.Parallel()

	// Grabs Service Principal environment variables, errors if any are missing
	EnvVars, err := setTerraformVariables()

	if err != nil {
		t.Fatal(err)
	}

	// Set input values for the test
	inputUniquePostfix := strings.ToLower(random.UniqueId())
	// inputResourceGroupName := fmt.Sprintf("%s-%s", "terratest-storage-account", inputUniquePostfix) // <---- Local switch
	inputResourceGroupName := "terratest-storage-account-debug-rg" // <---- Local switch
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
		EnvVars: EnvVars,

		// State backend - we start empty
		BackendConfig: globalBackendConf,

		// Colors in Terraform commands
		NoColor: false,
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	// defer terraform.Destroy(t, terraformOptions) // <---- Local switch

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	outputLocation := terraform.Output(t, terraformOptions, "primary_location")

	// = = = = = = = = = =
	// Run tests
	// = = = = = = = = = =
	t.Run("storage_account_location_tf_input_matched_output", func(t *testing.T) {
		assert.Equal(t, strings.ToLower(inputLocation), strings.ToLower(outputLocation), "Storage Account Location TF Input = TF Output")
	})

	// Call Azure SDK for ARM Storage to get back the Storage Account Location
	// https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/resourcemanager/storage/armstorage

	t.Run("test_1", func(t *testing.T) {
		assert.Equal(t, "one", "one", "Should always pass")
	})

	t.Run("test_2", func(t *testing.T) {
		assert.Equal(t, "one", "two", "Should always fail")
	})
}
