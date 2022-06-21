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
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"

	// Testing
	"github.com/stretchr/testify/assert"
)

const (
	storageAccountExampleGitDir = "../examples/storage-account"
)

// Test run that has skippable stages built in
func TestStorageAccountExampleWithStages(t *testing.T) {
	t.Parallel()

	// Deploy the Storage Account
	defer test_structure.RunTestStage(t, "teardown_storageAccount", func() { teardownStorageAccount(t, storageAccountExampleGitDir) })
	test_structure.RunTestStage(t, "deploy_storageAccount", func() { deployStorageAccount(t, storageAccountExampleGitDir) })

	// Validate the Storage Account is deployed correctly
	test_structure.RunTestStage(t, "validate_storageAccount", func() { validateStorageAccount(t, storageAccountExampleGitDir) })
}

// Creates Terraform Options for Storage Account module with remote state backend
func createStorageAccountOpts(t *testing.T, terraformDir string) *terraform.Options {
	uniqueId := strings.ToLower(random.UniqueId())

	// State backend environment variables
	stateBlobAccountNameForTesting := GetRequiredEnvVar(t, TerraformStateBlobStoreNameForTestEnvVarName)
	stateBlobAccountContainerForTesting := GetRequiredEnvVar(t, TerraformStateBlobStoreContainerForTestEnvVarName)
	stateBlobAccountKeyForTesting := GetRequiredEnvVar(t, TerraformStateBlobStoreKeyForTestEnvVarName)

	storageAccountStateKey := fmt.Sprintf("%s/%s/terraform.tfstate", t.Name(), uniqueId)

	return &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDir,

		// Variables to pass to our Terraform code using -var options.
		Vars: map[string]interface{}{
			"resource_group_name":    fmt.Sprintf("%s-%s", "terratest-storage-account", uniqueId),
			"location":               "canadacentral",
			"account_tier":           "Standard",
			"account_kind":           "StorageV2",
			"replication_type":       "GRS",
			"account_unique_postfix": uniqueId,
			"tags": map[string]string{
				"Source":  "terratest",
				"Owner":   "Raki Rahman",
				"Project": "Terraform CI testing",
			},
		},

		BackendConfig: map[string]interface{}{
			"storage_account_name": stateBlobAccountNameForTesting,
			"container_name":       stateBlobAccountContainerForTesting,
			"access_key":           stateBlobAccountKeyForTesting,
			"key":                  storageAccountStateKey,
		},

		// Service Principal creds from Environment Variables
		EnvVars: setTerraformVariables(t),

		// Colors in Terraform commands - we like colors
		NoColor: false,
	}
}

// Tears down the Storage Account
func teardownStorageAccount(t *testing.T, storageAccountDir string) {
	storageAccountOpts := test_structure.LoadTerraformOptions(t, storageAccountDir)
	defer terraform.Destroy(t, storageAccountOpts)
}

// Deploys storage account and stores Terraform Options locally for skipping stages
func deployStorageAccount(t *testing.T, storageAccountDir string) {
	storageAccountOpts := createStorageAccountOpts(t, storageAccountDir)

	// Save data to disk so that other test stages executed at a later time can read the data back in
	test_structure.SaveTerraformOptions(t, storageAccountDir, storageAccountOpts)

	terraform.InitAndApply(t, storageAccountOpts)
}

// Validates the Storage Account is deployed correctly
func validateStorageAccount(t *testing.T, storageAccountDir string) {
	storageAccountOpts := test_structure.LoadTerraformOptions(t, storageAccountDir)

	// 1. Validate Terraform location input variable and location output match
	TF_InputLocation := storageAccountOpts.Vars["location"].(string)
	TF_OutputLocation := terraform.Output(t, storageAccountOpts, "primary_location")

	t.Run("storage_account_location_tf_input_matched_tf_output", func(t *testing.T) {
		assert.Equal(t, strings.ToLower(TF_InputLocation), strings.ToLower(TF_OutputLocation), "Storage Account Location TF Input = TF Output")
	})

	// 2. Another such validation test is querying ARM - and we can have several of these - like file push/pull etc.
	validateStorageAccountWithARM(t, storageAccountOpts)
}

// Function calls ARM to validate the Storage Account is deployed correctly
func validateStorageAccountWithARM(t *testing.T, storageAccountOpts *terraform.Options) {
	// Authenticate to Azure and initiate context
	cred := getAzureCred(t)
	ctx := context.Background()

	TF_InputLocation := storageAccountOpts.Vars["location"].(string)
	TF_InputResourceGroupName := storageAccountOpts.Vars["resource_group_name"].(string)
	TF_OutputStorageAccountName := terraform.Output(t, storageAccountOpts, "storage_account_name")

	// Call Azure SDK for ARM Storage to get back the Storage Account Location
	// https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/resourcemanager/storage/armstorage
	properties, err := storageAccountProperties(ctx, cred, TF_InputResourceGroupName, TF_OutputStorageAccountName)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("storage_account_location_tf_input_matched_arm_output", func(t *testing.T) {
		assert.Equal(t, strings.ToLower(TF_InputLocation), strings.ToLower(*properties.Location), "Storage Account Location TF Input = ARM Output")
	})
}
