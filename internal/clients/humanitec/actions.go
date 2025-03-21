package humanitec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ActionPipelineSummary struct {
	OrgId       string `json:"org_id"`
	Id          string `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Type        string `json:"type"`
}

type ListActionPipelineSummariesResponse struct {
	HTTPResponse *http.Response
	Body         []byte
	JSON200      []ActionPipelineSummary
}

func (r ListActionPipelineSummariesResponse) StatusCode() int {
	return r.HTTPResponse.StatusCode
}

func (w *WrappedHumanitecClientImpl) ListActionPipelineSummaries(ctx context.Context, orgId string) (*ListActionPipelineSummariesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.apiPrefix+fmt.Sprintf("/orgs/%s/action-pipelines", orgId), nil)
	if err != nil {
		return &ListActionPipelineSummariesResponse{}, err
	}
	if err := w.requestEditor(ctx, req); err != nil {
		return &ListActionPipelineSummariesResponse{}, err
	}
	var out ListActionPipelineSummariesResponse
	out.HTTPResponse, err = w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if out.HTTPResponse.Body != nil {
		defer out.HTTPResponse.Body.Close()
		out.Body, err = io.ReadAll(out.HTTPResponse.Body)
		if err != nil {
			return &out, err
		}
	}
	if out.StatusCode() == 200 && out.Body != nil {
		var js200 []ActionPipelineSummary
		err = json.Unmarshal(out.Body, &js200)
		out.JSON200 = js200
	}
	return &out, err
}

type ActionPipeline struct {
	OrgId            string                 `json:"org_id"`
	Id               string                 `json:"id"`
	Description      string                 `json:"description"`
	CreatedAt        string                 `json:"created_at"`
	Type             string                 `json:"type"`
	PipelineId       string                 `json:"pipeline_id"`
	PipelineVersion  string                 `json:"pipeline_version"`
	Inputs           map[string]interface{} `json:"inputs"`
	InputsJsonSchema map[string]interface{} `json:"inputs_jsonschema"`
}

type GetActionPipelineResponse struct {
	HTTPResponse *http.Response
	Body         []byte
	JSON200      *ActionPipeline
}

func (r GetActionPipelineResponse) StatusCode() int {
	return r.HTTPResponse.StatusCode
}

func (w *WrappedHumanitecClientImpl) GetActionPipeline(ctx context.Context, orgId, id string) (*GetActionPipelineResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, w.apiPrefix+fmt.Sprintf("/orgs/%s/action-pipelines/%s", orgId, id), nil)
	if err != nil {
		return &GetActionPipelineResponse{}, err
	}
	if err := w.requestEditor(ctx, req); err != nil {
		return &GetActionPipelineResponse{}, err
	}
	var out GetActionPipelineResponse
	out.HTTPResponse, err = w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if out.HTTPResponse.Body != nil {
		defer out.HTTPResponse.Body.Close()
		out.Body, err = io.ReadAll(out.HTTPResponse.Body)
		if err != nil {
			return &out, err
		}
	}
	if out.StatusCode() == 200 && out.Body != nil {
		var js200 ActionPipeline
		err = json.Unmarshal(out.Body, &js200)
		out.JSON200 = &js200
	}
	return &out, err
}

type CallActionPipelineRequestBody struct {
	Inputs map[string]interface{} `json:"inputs"`
}

type CallActionPipelineParams struct {
	IdempotencyKey string
}

type CallActionPipelineResult struct {
	Outputs map[string]interface{} `json:"outputs"`
}

type CallActionPipelineResponse struct {
	HTTPResponse *http.Response
	Body         []byte
	JSON200      *CallActionPipelineResult
}

func (r CallActionPipelineResponse) StatusCode() int {
	return r.HTTPResponse.StatusCode
}

func (w *WrappedHumanitecClientImpl) CallActionPipeline(ctx context.Context, orgId, id string, params *CallActionPipelineParams, body CallActionPipelineRequestBody) (*CallActionPipelineResponse, error) {
	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.apiPrefix+fmt.Sprintf("/orgs/%s/action-pipelines/%s/calls", orgId, id), bytes.NewReader(raw))
	if err != nil {
		return &CallActionPipelineResponse{}, err
	}
	if params != nil {
		if params.IdempotencyKey != "" {
			req.Header.Add("Idempotency-Key", params.IdempotencyKey)
		}
	}
	if err := w.requestEditor(ctx, req); err != nil {
		return &CallActionPipelineResponse{}, err
	}
	var out CallActionPipelineResponse
	out.HTTPResponse, err = w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if out.HTTPResponse.Body != nil {
		defer out.HTTPResponse.Body.Close()
		out.Body, err = io.ReadAll(out.HTTPResponse.Body)
		if err != nil {
			return &out, err
		}
	}
	if out.StatusCode() == 200 && out.Body != nil {
		var js200 CallActionPipelineResult
		err = json.Unmarshal(out.Body, &js200)
		out.JSON200 = &js200
	}
	return &out, err
}
