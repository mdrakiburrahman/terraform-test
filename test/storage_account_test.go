package test

import (
	// Native
	"context"
	"fmt"
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

func TestStorageAccount(t *testing.T) {

	t.Parallel()
	ctx := context.Background()

	// Grabs Terraform authentication environment variables, errors if any are missing
	EnvVars, err := setTerraformVariables()

	if err != nil {
		t.Fatal(err)
	}

	// Authenticate to Azure and initiate context
	cred, err := getAzureCred()

	if err != nil {
		t.Fatal(err)
	}

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
		EnvVars: EnvVars,

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
	TF_OutputLocation := terraform.Output(t, terraformOptions, "primary_location")
	TF_OutputStorageAccountName := terraform.Output(t, terraformOptions, "storage_account_name")

	// = = = = = = = = = =
	// Run tests
	// = = = = = = = = = =
	t.Run("storage_account_location_tf_input_matched_tf_output", func(t *testing.T) {
		assert.Equal(t, strings.ToLower(inputLocation), strings.ToLower(TF_OutputLocation), "Storage Account Location TF Input = TF Output")
	})

	// Call Azure SDK for ARM Storage to get back the Storage Account Location
	// https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/resourcemanager/storage/armstorage
	properties, err := storageAccountProperties(ctx, cred, inputResourceGroupName, TF_OutputStorageAccountName)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("storage_account_location_tf_input_matched_arm_output", func(t *testing.T) {
		assert.Equal(t, strings.ToLower(inputLocation), strings.ToLower(*properties.Location), "Storage Account Location TF Input = ARM Output")
	})
}
