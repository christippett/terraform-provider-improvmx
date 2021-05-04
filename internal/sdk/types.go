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
	SetUserAgent(agent string) error
	SetHTTPClient(client *http.Client)

	GetAccount(ctx context.Context) (*Account, error)
	GetWhitelabels(ctx context.Context) (*[]Whitelabel, error)

	ListDomains(ctx context.Context, query *QueryDomain) (*[]Domain, error)
	AddDomain(ctx context.Context, domain *Domain) (*Domain, error)
	GetDomain(ctx context.Context, domain string) (*Domain, error)
	UpdateDomain(ctx context.Context, domain *Domain) (*Domain, error)
	DeleteDomain(ctx context.Context, domain *Domain) error
	CheckDomain(ctx context.Context, domain string) (*Check, error)

	ListAliases(ctx context.Context, domain string) (*[]Alias, error)
	CreateAlias(ctx context.Context, domain string, alias *Alias) (*Alias, error)
	UpdateAlias(ctx context.Context, domain string, alias *Alias) (*Alias, error)
	DeleteAlias(ctx context.Context, domain string, alias *Alias) error

	ListSMTPCredentials(ctx context.Context, domain string) (*[]SMTPCredential, error)
	CreateSMTPCredential(ctx context.Context, domain string, credential *WriteSMTPCredential) (*SMTPCredential, error)
	UpdateSMTPCredential(ctx context.Context, domain string, credential *WriteSMTPCredential) (*SMTPCredential, error)
	DeleteSMTPCredential(ctx context.Context, domain string, credential *SMTPCredential) error

	GetLogs(ctx context.Context, query *QueryLog) (*[]Log, error)
}

type client struct {
	apiKey     string
	url        string
	userAgent  *string
	httpClient *http.Client
	out        io.Writer
}

type PaginationOptions struct {
	Limit int `json:"limit,omitempty"`
	Page  int `json:"page,omitempty"`
}

type QueryDomain struct {
	Query    string `json:"q,omitempty"`
	IsActive bool   `json:"is_active,omitempty"`
	PaginationOptions
}

type Response struct {
	Success bool                `json:"success"`
	Errors  map[string][]string `json:"errors,omitempty"`
	Total   int                 `json:"total,omitempty"`
	PaginationOptions
}

type Account struct {
	BillingEmail   string        `json:"billing_email"`
	CancelsOn      interface{}   `json:"cancels_on"`
	CardBrand      string        `json:"card_brand"`
	CompanyDetails string        `json:"company_details"`
	CompanyName    string        `json:"company_name"`
	CompanyVat     interface{}   `json:"company_vat"`
	Country        string        `json:"country"`
	Created        int64         `json:"created"`
	Email          string        `json:"email"`
	EmailHash      string        `json:"email_hash"`
	IsOtpEnabled   bool          `json:"is_otp_enabled"`
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
	API          int `json:"api"`
	Credentials  int `json:"credentials"`
	DailyQuota   int `json:"daily_quota"`
	DailySend    int `json:"daily_send"`
	Destinations int `json:"destinations"`
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

type Whitelabel struct {
	Name string `json:"name"`
}

type Domain struct {
	Domain            string   `json:"domain"`
	Active            bool     `json:"active,omitempty"`
	Display           string   `json:"display,omitempty"`
	DkimSelector      string   `json:"dkim_selector,omitempty"`
	NotificationEmail string   `json:"notification_email,omitempty"`
	Webhook           string   `json:"webhook,omitempty"`
	Whitelabel        string   `json:"whitelabel,omitempty"`
	Added             int64    `json:"added,omitempty"`
	Aliases           *[]Alias `json:"aliases,omitempty"`
}

type Alias struct {
	Alias   string `json:"alias"`
	Forward string `json:"forward,omitempty"`
	ID      int    `json:"id,omitempty"`
}

type SMTPCredential struct {
	Username string `json:"username"`
	Usage    int    `json:"usage,omitempty"`
	Created  int64  `json:"created,omitempty"`
}

type WriteSMTPCredential struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type Check struct {
	Provider string      `json:"provider"`
	Advanced bool        `json:"advanced"`
	Dkim1    *Record     `json:"dkim1"`
	Dkim2    *Record     `json:"dkim2"`
	Dmarc    *Record     `json:"dmarc"`
	Mx       *Record     `json:"mx"`
	Spf      *Record     `json:"spf"`
	Valid    bool        `json:"valid"`
	Error    interface{} `json:"error,omitempty"`
}

type Record struct {
	Expected *RecordValues `json:"expected"`
	Valid    bool          `json:"valid"`
	Values   *RecordValues `json:"values"`
}

type RecordValues []string

func (values RecordValues) Interface() []interface{} {
	s := make([]interface{}, len(values))
	for i, v := range values {
		s[i] = v
	}
	return s
}

func (values RecordValues) MarshalJSON() ([]byte, error) {
	if len(values) == 1 {
		return []byte(fmt.Sprintf("%v", values[0])), nil
	}
	return []byte(fmt.Sprintf("[%v]", strings.Join(values, ","))), nil
}

type Log struct {
	Created    string `json:"created"`
	CreatedRaw string `json:"created_raw"`
	Events     []struct {
		Code    int    `json:"code"`
		Created string `json:"created"`
		ID      string `json:"id"`
		Local   string `json:"local"`
		Message string `json:"message"`
		Server  string `json:"server"`
		Status  string `json:"status"`
	} `json:"events,omitempty"`
	Forward struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"forward,omitempty"`
	Hostname  string `json:"hostname"`
	ID        string `json:"id"`
	MessageID string `json:"messageId"`
	Recipient struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"recipient,omitempty"`
	Sender struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"sender,omitempty"`
	Subject   string `json:"subject"`
	Transport string `json:"transport"`
}

type QueryLog struct {
	Domain *string `json:"domain"`
	Alias  *string `json:"alias"`
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
