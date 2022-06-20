# Terratest Hackathon

## Example 1: Simple Storage Account
> https://www.youtube.com/watch?v=UKitApmIHFM&ab_channel=Bee-a-Learner
> https://rakesh-suryawanshi.medium.com/test-azure-terraform-code-with-terratest-6c1b1249aea2
> https://github.com/bee-a-learner/terraform-test/blob/main/virtual_network/test/main_test.go

Once the code is writtern, perform the following to run the test:
```bash
cd /workspaces/terraform-test/1-storage-accnt-yt/modules/storage_account/test

# Can call it whatever we want - like terratestmodule
go mod init github.com/mdrakiburrahman/terraform-test/1-storage-accnt-yt

# This creates a go.sum file with all our dependencies linked to git commits, and cleans up ones not required
go mod tidy

# Runs all tests
go test -v -timeout 60m
```