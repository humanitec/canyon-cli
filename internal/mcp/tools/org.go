package tools

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/humanitec/humanitec-go-autogen/client"

	"github.com/humanitec/canyon-cli/internal"
	"github.com/humanitec/canyon-cli/internal/clients/humanitec"
	"github.com/humanitec/canyon-cli/internal/mcp"
)

func NewListHumanitecOrgsAndSession() mcp.Tool {
	return mcp.Tool{
		Name: "list_humanitec_orgs_and_session",
		Description: `This tool checks whether the local humctl (Humanitec CLI) tool has a valid and non-expired session.
This tool should be used if you don't know whether the user has a valid session or if other related tool commands return errors indicating the user is not authenticated.
This tool also returns the list of Organizations that the user has access to including their role in the Organization.
`,
		InputSchema: map[string]interface{}{"type": "object", "additionalProperties": false},
		Callable: func(ctx context.Context, m map[string]interface{}) ([]mcp.CallToolResponseContent, error) {
			hc, err := humanitec.NewHumanitecClientWithCurrentToken(ctx)
			if err != nil {
				return nil, err
			}
			if r, err := humanitec.CheckResponse(func() (*client.GetCurrentUserResponse, error) {
				return hc.GetCurrentUserWithResponse(ctx)
			}).AndStatusCodeEq(http.StatusOK).RespAndError(); err != nil {
				return nil, err
			} else {
				out := make(map[string]string)
				seenOrgs := make(map[string]bool)
				for obj, role := range r.JSON200.Roles {
					parts := strings.Split(obj, "/")
					if len(parts) == 3 {
						org := parts[2]
						if _, ok := seenOrgs[org]; !ok {
							out[org] = role
							seenOrgs[org] = true
						}
					}
				}
				rawOrgs := internal.PrettyJson(out)
				return []mcp.CallToolResponseContent{mcp.NewTextToolResponseContent(`The user is currently logged in. The following JSON is map from Humanitec Organization to Role:
%s
'administrators' can take all actions in the Organization, 'managers' may create applications and manage users, 'members' only have access to an application level, 'org_viewers' have read access to the whole Organization.`,
					string(rawOrgs),
				)}, nil
			}
		},
	}
}

func NewListAppsAndEnvsForOrganization() mcp.Tool {
	return mcp.Tool{
		Name: "list_apps_and_envs_for_humanitec_organization",
		Description: `This tool returns the Applications within the specified Humanitec Organization. It also includes the Environments within each Application including the latest deployment state and status.
An optional app_id regex argument can filter Application Ids, while the env_type argument can filter by Environment Type (eg: development, staging, production).
`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"org_id":   map[string]interface{}{"type": "string", "description": "The Humanitec Organization (org) ID to work with."},
				"app_id":   map[string]interface{}{"type": "string", "description": "Optional regex pattern to filter for app id"},
				"env_type": map[string]interface{}{"type": "string", "description": "Optional filter for a specific environment type"},
			},
			"required":             []string{"org_id"},
			"additionalProperties": false,
		},
		Callable: func(ctx context.Context, m map[string]interface{}) ([]mcp.CallToolResponseContent, error) {
			hc, err := humanitec.NewHumanitecClientWithCurrentToken(ctx)
			if err != nil {
				return nil, fmt.Errorf("unable to create Humanitec client: %w", err)
			}
			orgId, _ := m["org_id"].(string)

			var appIdPattern *regexp.Regexp
			if v, ok := m["app_id"].(string); ok {
				appIdPattern, err = regexp.CompilePOSIX(v)
				if err != nil {
					return nil, fmt.Errorf("invalid app_id  regex: %w", err)
				}
			}
			envTypeFilter, _ := m["env_type"].(string)

			if r, err := humanitec.CheckResponse(func() (*client.ListApplicationsResponse, error) {
				return hc.ListApplicationsWithResponse(ctx, orgId)
			}).AndStatusCodeEq(http.StatusOK).RespAndError(); err != nil {
				return nil, err
			} else {

				type envstate struct {
					Name               string    `json:"name"`
					Type               string    `json:"type"`
					CreatedTime        time.Time `json:"createdTime"`
					LastDeploymentId   string    `json:"lastDeploymentId,omitempty"`
					LastDeploymentSet  string    `json:"lastDeploymentSetId,omitempty"`
					LastDeploymentTime time.Time `json:"lastDeploymentTime,omitempty"`
				}

				type appstate struct {
					Name         string              `json:"name"`
					Environments map[string]envstate `json:"environments"`
					CreatedTime  string              `json:"createdTime"`
				}

				apps := new(sync.Map)
				{
					wg := new(sync.WaitGroup)
					sem := make(chan struct{}, 10)
					for _, app := range *r.JSON200 {
						if appIdPattern != nil && !appIdPattern.MatchString(app.Id) {
							continue
						}
						wg.Add(1)
						sem <- struct{}{}
						go func() {
							defer wg.Done()
							defer func() { <-sem }()
							if r, err := humanitec.CheckResponse(func() (*client.ListEnvironmentsResponse, error) {
								return hc.ListEnvironmentsWithResponse(ctx, orgId, app.Id)
							}).AndStatusCodeEq(http.StatusOK).RespAndError(); err != nil {
								apps.Store(app.Id, err.Error())
							} else {
								envs := make(map[string]envstate)
								for _, e := range *r.JSON200 {
									if envTypeFilter != "" && e.Type != envTypeFilter {
										continue
									}
									es := envstate{
										Name:        e.Name,
										Type:        e.Type,
										CreatedTime: e.CreatedAt,
									}
									if e.LastDeploy != nil {
										es.LastDeploymentId = e.LastDeploy.Id
										es.LastDeploymentSet = e.LastDeploy.SetId
										es.LastDeploymentTime = e.LastDeploy.CreatedAt
									}
									envs[e.Id] = es
								}
								apps.Store(app.Id, appstate{
									Name:         app.Name,
									CreatedTime:  app.CreatedAt,
									Environments: envs,
								})
							}
						}()
					}
					wg.Wait()
				}

				out := make(map[string]appstate)
				apps.Range(func(key, value any) bool {
					if e, ok := value.(error); ok {
						err = errors.Join(err, fmt.Errorf("failed to fetch app '%s': %w", key, e))
					} else if a, ok := value.(appstate); ok {
						out[key.(string)] = a
					}
					return true
				})

				if err != nil {
					return nil, err
				}

				rawApps := internal.PrettyJson(out)
				return []mcp.CallToolResponseContent{mcp.NewTextToolResponseContent("The user is has access to the following Humanitec Applications with Organization '%s' in JSON format: %s", orgId, string(rawApps))}, nil
			}
		},
	}
}

func NewGetWorkloadProfileSchema() mcp.Tool {
	return mcp.Tool{
		Name: "get_humanitec_workload_profile_schema",
		Description: `This tool returns information including the JSON schema used to define the workload profile with the specific id.
Multiple workload profiles exist.
The humanitec/ prefix is part of the workload profile id.
The profile schema includes the set of properties supported in Workloads specs that use this profile.`,
		InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{
			"org_id":              map[string]interface{}{"type": "string", "description": "The Humanitec Organization (org) ID to work with."},
			"workload_profile_id": map[string]interface{}{"type": "string", "description": "The Humanitec Workload Profile (profile) ID to work with."},
		}, "required": []string{"org_id", "workload_profile_id"}},
		Callable: func(ctx context.Context, arguments map[string]interface{}) ([]mcp.CallToolResponseContent, error) {
			orgId := arguments["org_id"].(string)
			workloadProfileId := arguments["workload_profile_id"].(string)
			hc, err := humanitec.NewHumanitecClientWithCurrentToken(ctx)
			if err != nil {
				return nil, err
			}
			if r, err := humanitec.CheckResponse(func() (*client.GetWorkloadProfileResponse, error) {
				return hc.GetWorkloadProfileWithResponse(ctx, orgId, workloadProfileId)
			}).AndStatusCodeEq(http.StatusOK).RespAndError(); err != nil {
				return nil, err
			} else {
				profileSchema := internal.PrettyJson(r.JSON200.SpecSchema)
				return []mcp.CallToolResponseContent{mcp.NewTextToolResponseContent(`The humanitec workload profile has the following JSON schema for the spec of a deployment set module: %s`, string(profileSchema))}, nil
			}
		},
	}
}
