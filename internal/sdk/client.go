package improvmx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// NewClient constructs a Checly API client.
func NewClient(
	//checkly API's base url
	baseURL,
	//checkly's api key
	apiKey string,
	//optional, defaults to http.DefaultClient
	httpClient *http.Client,
) Client {
	c := &client{
		apiKey:     apiKey,
		url:        baseURL,
		httpClient: httpClient,
	}
	if httpClient != nil {
		c.httpClient = httpClient
	} else {
		c.httpClient = http.DefaultClient
	}
	return c
}

func (c *client) AddDomain(ctx context.Context, domain Domain) (*Domain, error) {
	var result DomainResponse
	if err := c.apiCall(
		ctx,
		http.MethodPost,
		"/domains/",
		domain,
		&result,
	); err != nil {
		return nil, err
	}
	return &result.Domain, nil
}

func (c *client) GetDomain(ctx context.Context, domain string) (*Domain, error) {
	var result DomainResponse
	if err := c.apiCall(
		ctx,
		http.MethodGet,
		fmt.Sprintf("/domains/%s", domain),
		nil,
		&result,
	); err != nil {
		return nil, err
	}
	return &result.Domain, nil
}

func (c *client) ListDomains(ctx context.Context) (*[]Domain, error) {
	var result DomainListResponse
	if err := c.apiCall(
		ctx,
		http.MethodGet,
		"/domains/",
		nil,
		&result,
	); err != nil {
		return nil, err
	}
	return &result.Domains, nil
}

func (c *client) handleResponseError(resp *http.Response) error {
	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("error decoding response error: %s", err)
	}
	if res.Errors != nil {
		for k, v := range res.Errors {
			fmt.Fprintf(c.debug, "error (%s): %s", k, strings.Join(v, "; "))
			fmt.Fprintln(c.debug)
		}
	}
	return fmt.Errorf("response error: %s", http.StatusText(resp.StatusCode))

}

func (c *client) apiCall(
	ctx context.Context,
	method string,
	URL string,
	body interface{},
	result interface{},
) error {
	requestURL := c.url + URL

	data, err := json.Marshal(&body)
	if err != nil {
		return fmt.Errorf("error generating request payload: %v", err)
	}

	req, err := http.NewRequest(method, requestURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Authorization", "Basic  api:"+c.apiKey)
	req.Header.Add("content-type", "application/json")

	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed with: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleResponseError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding response data: %s", err)
	}

	return nil
}
