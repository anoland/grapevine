# grapevine
Repo for the grapevine project

## End-to-End (E2E) Tests

This project includes end-to-end tests that provision a real Kubernetes cluster on Rackspace Spot and verify its functionality.

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18 or newer)
- [Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)
- A Rackspace API token

### Running the Tests

1. **Set your Rackspace API token as an environment variable:**
   ```sh
   export RACKSPACE_TOKEN="your-token-here"
   ```

2. **Navigate to the test directory and run the tests:**
   ```sh
   cd tests/e2e
   go test -v
   ```

   The tests will:
   - Initialize and apply the Terraform configuration to create a Kubernetes cluster.
   - Run Ginkgo tests to verify node status and deploy a sample application.
   - Automatically destroy the Terraform resources upon completion.
