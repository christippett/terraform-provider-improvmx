package improvmx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
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
	out io.Writer,
) Client {
	c := &client{
		apiKey:     apiKey,
		url:        baseURL,
		httpClient: httpClient,
		out:        out,
	}
	if httpClient != nil {
		c.httpClient = httpClient
	} else {
		c.httpClient = http.DefaultClient
	}
	return c
}

/* ACCOUNT ⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁ */

func (c *client) GetAccount(ctx context.Context) (*Account, error) {
	var result struct {
		Account *Account `json:"account,omitempty"`
		Response
	}
	if err := c.apiCall(ctx, http.MethodGet, "/account/", nil, &result); err != nil {
		return nil, err
	}
	return result.Account, nil
}

func (c *client) GetWhitelabels(ctx context.Context) (*[]Whitelabel, error) {
	var result struct {
		Whitelabels *[]Whitelabel `json:"whitelabels,omitempty"`
		Response
	}

	url := "/account/whitelabels"
	if err := c.apiCall(ctx, http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return result.Whitelabels, nil
}

/* DOMAIN ⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁ */

func (c *client) ListDomains(ctx context.Context, query *QueryDomain) (*[]Domain, error) {
	var result struct {
		Domains *[]Domain `json:"domains,omitempty"`
		Response
	}
	if err := c.apiCall(ctx, http.MethodGet, "/domains/", query, &result); err != nil {
		return nil, err
	}
	return result.Domains, nil
}

func (c *client) AddDomain(ctx context.Context, domain *Domain) (*Domain, error) {
	var result struct {
		Domain *Domain `json:"domain,omitempty"`
		Response
	}

	url := "/domains/"
	if err := c.apiCall(ctx, http.MethodPost, url, domain, &result); err != nil {
		return nil, err
	}
	return result.Domain, nil
}

func (c *client) GetDomain(ctx context.Context, domain string) (*Domain, error) {
	var result struct {
		Domain *Domain `json:"domain,omitempty"`
		Response
	}

	url := fmt.Sprintf("/domains/%s", domain)
	if err := c.apiCall(ctx, http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return result.Domain, nil
}

func (c *client) UpdateDomain(ctx context.Context, domain *Domain) (*Domain, error) {
	var result struct {
		Domain *Domain `json:"domain,omitempty"`
		Response
	}

	url := fmt.Sprintf("/domains/%s", domain.Domain)
	if err := c.apiCall(ctx, http.MethodPut, url, domain, &result); err != nil {
		return nil, err
	}
	return result.Domain, nil
}

func (c *client) DeleteDomain(ctx context.Context, domain *Domain) error {
	var result Response
	url := fmt.Sprintf("/domains/%s", domain.Domain)
	if err := c.apiCall(ctx, http.MethodDelete, url, nil, &result); err != nil {
		return err
	}
	return nil
}

func (c *client) CheckDomain(ctx context.Context, domain *Domain) (*Check, error) {
	var result struct {
		Records *Check `json:"records,omitempty"`
		Response
	}

	url := fmt.Sprintf("/domains/%s/check", domain.Domain)
	if err := c.apiCall(ctx, http.MethodGet, url, domain, &result); err != nil {
		return nil, err
	}
	return result.Records, nil
}

/* ALIAS ⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁ */

func (c *client) ListAliases(ctx context.Context, domain string) (*[]Alias, error) {
	var result struct {
		Aliases *[]Alias `json:"aliases,omitempty"`
		Response
	}

	url := fmt.Sprintf("/domains/%s/aliases/", domain)
	if err := c.apiCall(ctx, http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return result.Aliases, nil
}

func (c *client) CreateAlias(ctx context.Context, domain string, alias *Alias) (*Alias, error) {
	var result struct {
		Alias *Alias `json:"alias,omitempty"`
		Response
	}

	url := fmt.Sprintf("/domains/%s/aliases/", domain)
	if err := c.apiCall(ctx, http.MethodPost, url, alias, &result); err != nil {
		return nil, err
	}
	return result.Alias, nil
}

func (c *client) UpdateAlias(ctx context.Context, domain string, alias *Alias) (*Alias, error) {
	var result struct {
		Alias *Alias `json:"alias,omitempty"`
		Response
	}

	url := fmt.Sprintf("/domains/%s/aliases/%s", domain, alias.Alias)
	if err := c.apiCall(ctx, http.MethodPut, url, alias, &result); err != nil {
		return nil, err
	}
	return result.Alias, nil
}

func (c *client) DeleteAlias(ctx context.Context, domain string, alias *Alias) error {
	var result Response

	url := fmt.Sprintf("/domains/%s/aliases/%s", domain, alias.Alias)
	if err := c.apiCall(ctx, http.MethodDelete, url, nil, &result); err != nil {
		return err
	}
	return nil
}

/* SMTP CREDENTIAL ⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁ */

func (c *client) ListSMTPCredentials(ctx context.Context, domain string) (*[]SMTPCredential, error) {
	var result struct {
		Credentials *[]SMTPCredential `json:"credentials,omitempty"`
		Response
	}

	url := fmt.Sprintf("/domains/%s/credentials/", domain)
	if err := c.apiCall(ctx, http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return result.Credentials, nil
}

func (c *client) CreateSMTPCredential(ctx context.Context, domain string, credential *WriteSMTPCredential) (*SMTPCredential, error) {
	var result struct {
		Credential *SMTPCredential `json:"credential,omitempty"`
		Response
	}

	url := fmt.Sprintf("/domains/%s/credentials/", domain)
	if err := c.apiCall(ctx, http.MethodPost, url, credential, &result); err != nil {
		return nil, err
	}
	return result.Credential, nil
}

func (c *client) UpdateSMTPCredential(ctx context.Context, domain string, credential *WriteSMTPCredential) (*SMTPCredential, error) {
	var result struct {
		Credential *SMTPCredential `json:"credential,omitempty"`
		Response
	}

	url := fmt.Sprintf("/domains/%s/credentials/%s", domain, credential.Username)
	if err := c.apiCall(ctx, http.MethodPut, url, credential, &result); err != nil {
		return nil, err
	}
	return result.Credential, nil
}

func (c *client) DeleteSMTPCredential(ctx context.Context, domain string, credential *SMTPCredential) error {
	var result Response

	url := fmt.Sprintf("/domains/%s/credentials/%s", domain, credential.Username)
	if err := c.apiCall(ctx, http.MethodDelete, url, nil, &result); err != nil {
		return err
	}
	return nil
}

/* DOMAIN / ALIAS LOG ⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁ */

func (c *client) GetLogs(ctx context.Context, query *QueryLog) (*[]Log, error) {
	var result struct {
		Logs *[]Log `json:"logs,omitempty"`
		Response
	}

	var url string
	if query.Alias != nil {
		url = fmt.Sprintf("/domains/%s/logs/%s", *query.Domain, *query.Alias)
	} else {
		url = fmt.Sprintf("/domains/%s/logs", *query.Domain)
	}

	if err := c.apiCall(ctx, http.MethodGet, url, nil, &result); err != nil {
		return nil, err
	}
	return result.Logs, nil
}

/* Misc ⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁⌁ */

func (c *client) handleResponseError(resp *http.Response) error {
	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("error decoding response error: %s", err)
	}
	if res.Errors != nil {
		for k, v := range res.Errors {
			fmt.Fprintf(c.out, "error (%s): %s", k, strings.Join(v, "; "))
			fmt.Fprintln(c.out)
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

	// Log request output to stdout
	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return fmt.Errorf("error dumping HTTP request: %v", err)
	}
	fmt.Fprintln(c.out, string(requestDump))
	fmt.Fprintln(c.out)

	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed with: %v", err)
	}

	defer resp.Body.Close()

	// Log response output to stdout
	responseDump, _ := httputil.DumpResponse(resp, true)
	fmt.Fprintln(c.out, string(responseDump))
	fmt.Fprintln(c.out)

	if resp.StatusCode != http.StatusOK {
		return c.handleResponseError(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding response data: %s", err)
	}

	return nil
}
