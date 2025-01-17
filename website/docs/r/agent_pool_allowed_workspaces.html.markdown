---
layout: "tfe"
page_title: "Terraform Enterprise: tfe_agent_pool_allowed_workspaces"
description: |-
  Manages allowed workspaces on agent pools
---

# tfe_agent_pool_allowed_workspaces

Adds and removes allowed workspaces on an agent pool

~> **NOTE:** This resource requires using the provider with Terraform Cloud and a Terraform Cloud
for Business account.
[Learn more about Terraform Cloud pricing here](https://www.hashicorp.com/products/terraform/pricing).

## Example Usage

Basic usage:

```hcl
resource "tfe_organization" "test-organization" {
  name  = "my-org-name"
  email = "admin@company.com"
}

resource "tfe_workspace" "test-workspace" {
  name         = "my-workspace-name"
  organization = tfe_organization.test-organization.name
}

resource "tfe_agent_pool" "test-agent-pool" {
  name                = "my-agent-pool-name"
  organization        = tfe_organization.test-organization.name
  organization_scoped = false
}

resource "tfe_agent_pool_allowed_workspaces" "test-allowed-workspaces" {
  agent_pool_id         = tfe_agent_pool.test-agent-pool.id
  allowed_workspace_ids = [tfe_workspace.test-workspace.id]
}
```

## Argument Reference

The following arguments are supported:

* `agent_pool_id` - (Required) The ID of the agent pool.
* `allowed_workspace_ids` - (Required) IDs of workspaces to be added as allowed workspaces on the agent pool.


## Import

A resource can be imported; use `<AGENT POOL ID>` as the import ID. For example:

```shell
terraform import tfe_agent_pool_allowed_workspaces.foobar apool-rW0KoLSlnuNb5adB
```

