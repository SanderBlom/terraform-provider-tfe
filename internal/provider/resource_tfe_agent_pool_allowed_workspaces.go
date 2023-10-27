// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTFEAgentPoolAllowedWorkspaces() *schema.Resource {
	return &schema.Resource{
		Create: resourceTFEAgentPoolAllowedWorkspacesCreate,
		Read:   resourceTFEAgentPoolAllowedWorkspacesRead,
		Update: resourceTFEAgentPoolAllowedWorkspacesUpdate,
		Delete: resourceTFEAgentPoolAllowedWorkspacesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"agent_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"allowed_workspace_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceTFEAgentPoolAllowedWorkspacesCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(ConfiguredClient)

	apID := d.Get("agent_pool_id").(string)

	// Before executing, make an API call to check the organization_scoped bool
	// If organization_scoped is true, fail fast
	// The purpose of this resource is to allow workspaces access to the agent pool
	// When no workspaces have been given access
	agentPool, err := config.Client.AgentPools.Read(ctx, apID)
	if err != nil {
		if errors.Is(err, tfe.ErrResourceNotFound) {
			log.Printf("[DEBUG] agent pool %s no longer exists: ", apID)
			return nil
		}
		return fmt.Errorf("Error reading configuration of agent pool: %s %w", apID, err)
	}

	// Create a new options struct.
	options := tfe.AgentPoolAllowedWorkspacesUpdateOptions{}

	if agentPool.OrganizationScoped {
		return fmt.Errorf("error updating allowed workspaces on agent pool, workspaces already scoped for access to organization: %s", agentPool.Organization.Name)
	} else if !agentPool.OrganizationScoped {
		if allowedWorkspaceIDs, allowedWorkspaceSet := d.GetOk("allowed_workspace_ids"); allowedWorkspaceSet {
			options.AllowedWorkspaces = []*tfe.Workspace{}
			for _, workspaceID := range allowedWorkspaceIDs.(*schema.Set).List() {
				if val, ok := workspaceID.(string); ok {
					options.AllowedWorkspaces = append(options.AllowedWorkspaces, &tfe.Workspace{ID: val})
				}
			}
		}

		log.Printf("[DEBUG] Update agent pool: %s", apID)
		_, err := config.Client.AgentPools.UpdateAllowedWorkspaces(ctx, apID, options)
		if err != nil {
			return fmt.Errorf("error updating agent pool %s: %w", apID, err)
		}

		d.SetId(apID)
	}

	return resourceTFEAgentPoolAllowedWorkspacesRead(d, meta)
}

func resourceTFEAgentPoolAllowedWorkspacesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(ConfiguredClient)

	agentPool, err := config.Client.AgentPools.Read(ctx, d.Id())
	if err != nil {
		if errors.Is(err, tfe.ErrResourceNotFound) {
			log.Printf("[DEBUG] agent pool %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading configuration of agent pool %s: %w", d.Id(), err)
	}

	var allowedWorkspaceIDs []string
	for _, workspace := range agentPool.AllowedWorkspaces {
		allowedWorkspaceIDs = append(allowedWorkspaceIDs, workspace.ID)
	}
	d.Set("allowed_workspace_ids", allowedWorkspaceIDs)
	d.Set("agent_pool_id", agentPool.ID)

	return nil
}

func resourceTFEAgentPoolAllowedWorkspacesUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(ConfiguredClient)

	apID := d.Get("agent_pool_id").(string)

	// See comments at resourceTFEAgentPoolAllowedWorkspacesCreate()
	agentPool, err := config.Client.AgentPools.Read(ctx, apID)
	if err != nil {
		if errors.Is(err, tfe.ErrResourceNotFound) {
			log.Printf("[DEBUG] agent pool %s no longer exists: ", apID)
			return nil
		}
		return fmt.Errorf("Error reading configuration of agent pool: %s %w", apID, err)
	}

	// Create a new options struct.
	options := tfe.AgentPoolAllowedWorkspacesUpdateOptions{
		AllowedWorkspaces: []*tfe.Workspace{},
	}

	if agentPool.OrganizationScoped {
		return fmt.Errorf("error updating allowed workspaces on agent pool, workspaces already scoped for access to organization: %s", agentPool.Organization.Name)
	} else if !agentPool.OrganizationScoped {
		if allowedWorkspaceIDs, allowedWorkspaceSet := d.GetOk("allowed_workspace_ids"); allowedWorkspaceSet {
			options.AllowedWorkspaces = []*tfe.Workspace{}
			for _, workspaceID := range allowedWorkspaceIDs.(*schema.Set).List() {
				if val, ok := workspaceID.(string); ok {
					options.AllowedWorkspaces = append(options.AllowedWorkspaces, &tfe.Workspace{ID: val})
				}
			}
		}

		log.Printf("[DEBUG] Update agent pool: %s", apID)
		_, err := config.Client.AgentPools.UpdateAllowedWorkspaces(ctx, apID, options)
		if err != nil {
			return fmt.Errorf("error updating agent pool %s: %w", apID, err)
		}

		d.SetId(apID)
	}

	return resourceTFEAgentPoolAllowedWorkspacesRead(d, meta)
}

func resourceTFEAgentPoolAllowedWorkspacesDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(ConfiguredClient)

	apID := d.Get("agent_pool_id").(string)

	// Create a new options struct.
	options := tfe.AgentPoolAllowedWorkspacesUpdateOptions{
		AllowedWorkspaces: []*tfe.Workspace{},
	}

	log.Printf("[DEBUG] Update agent pool: %s", apID)
	_, err := config.Client.AgentPools.UpdateAllowedWorkspaces(ctx, apID, options)
	if err != nil {
		return fmt.Errorf("error updating agent pool %s: %w", apID, err)
	}

	return nil
}
