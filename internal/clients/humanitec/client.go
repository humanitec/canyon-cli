package humanitec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"slices"

	"github.com/humanitec/humanitec-go-autogen/client"
	"gopkg.in/yaml.v3"
)

func GetCurrentHumanitecToken() (string, error) {
	if v := os.Getenv("HUMANITEC_TOKEN"); v != "" {
		return v, nil
	}
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to identify the users home directory: %w", err)
	}
	if content, err := os.ReadFile(filepath.Join(hd, ".humctl")); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read the humctl file: %w", err)
	} else {
		s := struct {
			Token string `yaml:"token"`
		}{}
		if err := yaml.Unmarshal(content, &s); err != nil {
			return "", fmt.Errorf("failed to unmarshal the humctl file: %w", err)
		}
		return s.Token, nil
	}
}

type contextKey int

const (
	overrideHumanitecClientKey contextKey = iota
)

type WrappedHumanitecClient interface {
	client.ClientWithResponsesInterface
}

type WrappedHumanitecClientImpl struct {
	client.ClientWithResponsesInterface
	apiPrefix     string
	httpClient    client.HttpRequestDoer
	requestEditor client.RequestEditorFn
}

func NewHumanitecClientWithCurrentToken(ctx context.Context) (*WrappedHumanitecClientImpl, error) {
	token, err := GetCurrentHumanitecToken()
	if err != nil {
		return nil, err
	} else if token == "" {
		return nil, fmt.Errorf("The user is not currently logged in and should be prompted to run 'humctl login' to fix this.")
	}
	apiPrefix := "https://api.humanitec.io"
	if v := os.Getenv("HUMANITEC_API_PREFIX"); v != "" {
		apiPrefix = v
	}
	bi, _ := debug.ReadBuildInfo()
	wci := &WrappedHumanitecClientImpl{
		apiPrefix:  apiPrefix,
		httpClient: http.DefaultClient,
		requestEditor: func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Humanitec-User-Agent", fmt.Sprintf("app %s/%s; sdk humanitec-go-autogen/latest", filepath.Base(bi.Main.Path), bi.Main.Version))
			return nil
		},
	}
	if v, ok := ctx.Value(overrideHumanitecClientKey).(client.HttpRequestDoer); ok {
		wci.httpClient = v
	}
	wci.ClientWithResponsesInterface, err = client.NewClientWithResponses(apiPrefix, client.WithHTTPClient(wci.httpClient), client.WithRequestEditorFn(wci.requestEditor))
	return wci, err
}

type checkableResponse interface {
	StatusCode() int
}

type CheckedResponse[k checkableResponse] struct {
	Err      error
	Response k
}

func (ac *CheckedResponse[k]) AndStatusCodeEq(code int, codes ...int) *CheckedResponse[k] {
	var r checkableResponse = ac.Response
	if r != nil && code != r.StatusCode() && !slices.Contains(codes, r.StatusCode()) {
		if r.StatusCode() == http.StatusForbidden {
			ac.Err = errors.Join(ac.Err, fmt.Errorf("The user is not currently logged in and should be prompted to run 'humctl login' to fix this."))
		} else if r.StatusCode() == http.StatusNotFound {
			ac.Err = errors.Join(ac.Err, fmt.Errorf("The API request returned a 404 (Not Found) error which may indicate that the resource does not exist. The user may have misspelt something or the state may have changed."))
		} else {
			body := make([]byte, 0)
			v := reflect.ValueOf(ac.Response)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			f := v.FieldByName("Body")
			if f.Kind() == reflect.Slice {
				body = f.Interface().([]byte)
			}
			anon := struct {
				Message string `yaml:"message"`
			}{}
			_ = json.Unmarshal(body, &anon)
			bodyText := string(body)
			if anon.Message != "" {
				bodyText = anon.Message
			}
			ac.Err = errors.Join(ac.Err, fmt.Errorf(
				"The API request to Humanitec returned an unexpected status code %d (%s). The content of the error response is '%s' and may provide a hint as to what went wrong.", ac.Response.StatusCode(), http.StatusText(ac.Response.StatusCode()), bodyText))
		}
	}
	return ac
}

func (ac *CheckedResponse[k]) RespAndError() (k, error) {
	return ac.Response, ac.Err
}

func CheckResponse[k checkableResponse](requester func() (k, error)) *CheckedResponse[k] {
	resp, err := requester()
	if err != nil {
		if ne := (net.Error)(nil); errors.As(err, &ne) {
			err = fmt.Errorf("The API request to Humanitec hit a temporary network error '%s'. The request may work if the user requests it again.", ne.Error())
		} else {
			err = fmt.Errorf("The API request to Humanitec hit an unexpected error '%s'.", ne.Error())
		}
	}
	return &CheckedResponse[k]{Response: resp, Err: err}
}
