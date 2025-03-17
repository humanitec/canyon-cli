package humanitec

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type QueryAiDocsResponseJSON200 struct {
	Answer      string `json:"answer"`
	IsUncertain bool   `json:"is_uncertain"`
}

type QueryAiDocsResponse struct {
	HTTPResponse *http.Response
	Body         []byte
	JSON200      *QueryAiDocsResponseJSON200
}

func (q QueryAiDocsResponse) StatusCode() int {
	return q.HTTPResponse.StatusCode
}

func (w *WrappedHumanitecClientImpl) QueryAiDocs(ctx context.Context, query string) (*QueryAiDocsResponse, error) {
	type requestBody struct {
		Query string `json:"query"`
	}
	body, _ := json.Marshal(&requestBody{Query: query})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.apiPrefix+"/experimental/query-ai-documentation", bytes.NewReader(body))
	if err != nil {
		return &QueryAiDocsResponse{}, err
	}
	if err := w.requestEditor(ctx, req); err != nil {
		return &QueryAiDocsResponse{}, err
	}
	var out QueryAiDocsResponse
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
		var js200 QueryAiDocsResponseJSON200
		err = json.Unmarshal(out.Body, &js200)
		out.JSON200 = &js200
	}
	return &out, err
}
