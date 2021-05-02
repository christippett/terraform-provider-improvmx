package improvmx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client interface {
	ListDomains(ctx context.Context) (*[]Domain, error)
	AddDomain(ctx context.Context, domain *Domain) (*Domain, error)
	GetDomain(ctx context.Context, domain *string) (*Domain, error)
	UpdateDomain(ctx context.Context, domain *Domain) (*Domain, error)
	DeleteDomain(ctx context.Context, domain *Domain) error
	CheckDomain(ctx context.Context, domain *Domain) (*Check, error)
}

type client struct {
	apiKey     string
	url        string
	httpClient *http.Client
	out        io.Writer
}

type DomainQuery struct {
	Query    string `json:"q,omitempty"`
	IsActive bool   `json:"is_active,omitempty"`
}

type PaginationOptions struct {
	Limit int `json:"limit,omitempty"`
	Page  int `json:"page,omitempty"`
}

type Response struct {
	Errors  map[string][]string `json:"errors,omitempty"`
	Total   int                 `json:"total,omitempty"`
	Limit   int                 `json:"limit,omitempty"`
	Page    int                 `json:"page,omitempty"`
	Success bool                `json:"success"`
}

type AccountResponse struct {
	account Account `json:"account,omitempty"`
	Response
}

type Account struct {
	BillingEmail   string        `json:"billing_email"`
	CancelsOn      int64         `json:"cancels_on"`
	CardBrand      string        `json:"card_brand"`
	CompanyDetails string        `json:"company_details"`
	CompanyName    string        `json:"company_name"`
	CompanyVat     string        `json:"company_vat"`
	Country        string        `json:"country"`
	Created        int64         `json:"created"`
	Email          string        `json:"email"`
	Last4          string        `json:"last4"`
	Limits         *AccountLimit `json:"limits"`
	LockReason     string        `json:"lock_reason"`
	Locked         bool          `json:"locked"`
	Password       bool          `json:"password"`
	Plan           *AccountPlan  `json:"plan"`
	Premium        bool          `json:"premium"`
	PrivacyLevel   int           `json:"privacy_level"`
	RenewDate      int64         `json:"renew_date"`
}

type AccountLimit struct {
	Aliases      int `json:"aliases"`
	DailyQuota   int `json:"daily_quota"`
	Domains      int `json:"domains"`
	Ratelimit    int `json:"ratelimit"`
	Redirections int `json:"redirections"`
	Subdomains   int `json:"subdomains"`
}

type AccountPlan struct {
	AliasesLimit int    `json:"aliases_limit"`
	DailyQuota   int    `json:"daily_quota"`
	Display      string `json:"display"`
	DomainsLimit int    `json:"domains_limit"`
	Kind         string `json:"kind"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Yearly       bool   `json:"yearly"`
}

type WhitelabelResponse struct {
	Whitelabels []struct {
		Name string `json:"name"`
	} `json:"whitelabels,omitempty"`
	Response
}

type Domain struct {
	Active            bool     `json:"active,omitempty"`
	Domain            string   `json:"domain"`
	Display           string   `json:"display,omitempty"`
	DkimSelector      string   `json:"dkim_selector,omitempty"`
	NotificationEmail string   `json:"notification_email,omitempty"`
	Webhook           string   `json:"webhook,omitempty"`
	Whitelabel        string   `json:"whitelabel,omitempty"`
	Added             int64    `json:"added,omitempty"`
	Aliases           []*Alias `json:"aliases,omitempty"`
}

type Alias struct {
	Forward string `json:"forward"`
	Alias   string `json:"alias"`
	ID      int    `json:"id"`
}

type SMTPCredential struct {
	Created  int64  `json:"created"`
	Usage    int    `json:"usage"`
	Username string `json:"username"`
}

type CheckResponse struct {
	Records []*Check `json:"records"`
	Response
}

type Check struct {
	Provider string      `json:"provider"`
	Advanced bool        `json:"advanced"`
	Dkim1    Record      `json:"dkim1"`
	Dkim2    Record      `json:"dkim2"`
	Dmarc    Record      `json:"dmarc"`
	Mx       Record      `json:"mx"`
	Spf      Record      `json:"spf"`
	Valid    bool        `json:"valid"`
	Error    interface{} `json:"error"`
}

type Record struct {
	Expected string       `json:"expected"`
	Valid    bool         `json:"valid"`
	Values   RecordValues `json:"values"`
}

type RecordValues []string

func (values RecordValues) MarshalJSON() ([]byte, error) {
	if len(values) == 1 {
		return []byte(fmt.Sprintf("%v", values[0])), nil
	}
	return []byte(fmt.Sprintf("[%v]", strings.Join(values, ","))), nil
}

func (values *RecordValues) UnmarshalJSON(b []byte) error {
	// Try array of strings first.
	var valueArr []string
	err := json.Unmarshal(b, &valueArr)
	if err != nil {
		// Convert single value to slice
		var value string
		if err := json.Unmarshal(b, &value); err != nil {
			return err
		}
		valueArr = append(valueArr, value)
	}
	*values = valueArr
	return nil
}
