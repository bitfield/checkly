package checkly

import (
	"io"
	"net/http"
	"time"
)

// Client represents a Checkly client. If the Debug field is set to an io.Writer
// (for example os.Stdout), then the client will dump API requests and responses
// to it.  To use a non-default HTTP client (for example, for testing, or to set
// a timeout), assign to the HTTPClient field. To set a non-default URL (for
// example, for testing), assign to the URL field.
type Client struct {
	apiKey     string
	URL        string
	HTTPClient *http.Client
	Debug      io.Writer
}

// Check type constants

// TypeBrowser is used to identify a browser check.
const TypeBrowser = "BROWSER"

// TypeAPI is used to identify an API check.
const TypeAPI = "API"

// Escalation type constants

// RunBased identifies a run-based escalation type, for use with an AlertSettings.
const RunBased = "RUN_BASED"

// TimeBased identifies a time-based escalation type, for use with an AlertSettings.
const TimeBased = "TIME_BASED"

// Assertion source constants

// StatusCode identifies the HTTP status code as an assertion source.
const StatusCode = "STATUS_CODE"

// JSONBody identifies the JSON body data as an assertion source.
const JSONBody = "JSON_BODY"

// TextBody identifies the response body text as an assertion source.
const TextBody = "TEXT_BODY"

// Headers identifies the HTTP headers as an assertion source.
const Headers = "HEADERS"

// ResponseTime identifies the response time as an assertion source.
const ResponseTime = "RESPONSE_TIME"

// Assertion comparison constants

// Equals asserts that the source and target are equal.
const Equals = "EQUALS"

// NotEquals asserts that the source and target are not equal.
const NotEquals = "NOT_EQUALS"

// IsEmpty asserts that the source is empty.
const IsEmpty = "IS_EMPTY"

// NotEmpty asserts that the source is not empty.
const NotEmpty = "NOT_EMPTY"

// GreaterThan asserts that the source is greater than the target.
const GreaterThan = "GREATER_THAN"

// LessThan asserts that the source is less than the target.
const LessThan = "LESS_THAN"

// Contains asserts that the source contains a specified value.
const Contains = "CONTAINS"

// NotContains asserts that the source does not contain a specified value.
const NotContains = "NOT_CONTAINS"

// Check represents the parameters for an existing check.
type Check struct {
	ID                        string                `json:"id"`
	Name                      string                `json:"name"`
	Type                      string                `json:"checkType"`
	Frequency                 int                   `json:"frequency"`
	Activated                 bool                  `json:"activated"`
	Muted                     bool                  `json:"muted"`
	ShouldFail                bool                  `json:"shouldFail"`
	Locations                 []string              `json:"locations"`
	DegradedResponseTime      int                   `json:"degradedResponseTime"`
	MaxResponseTime           int                   `json:"maxResponseTime"`
	Script                    string                `json:"script,omitempty"`
	CreatedAt                 time.Time             `json:"created_at,omitempty"`
	UpdatedAt                 time.Time             `json:"updated_at,omitempty"`
	EnvironmentVariables      []EnvironmentVariable `json:"environmentVariables"`
	DoubleCheck               bool                  `json:"doubleCheck"`
	Tags                      []string              `json:"tags,omitempty"`
	SSLCheck                  bool                  `json:"sslCheck"`
	SSLCheckDomain            string                `json:"sslCheckDomain"`
	SetupSnippetID            int64                 `json:"setupSnippetId,omitempty"`
	TearDownSnippetID         int64                 `json:"tearDownSnippetId,omitempty"`
	LocalSetupScript          string                `json:"localSetupScript,omitempty"`
	LocalTearDownScript       string                `json:"localTearDownScript,omitempty"`
	AlertSettings             AlertSettings         `json:"alertSettings,omitempty"`
	UseGlobalAlertSettings    bool                  `json:"useGlobalAlertSettings"`
	Request                   Request               `json:"request"`
	AlertChannelSubscriptions []Subscription        `json:"alertChannelSubscriptions"`
}

// Request represents the parameters for the request made by the check.
type Request struct {
	Method          string      `json:"method"`
	URL             string      `json:"url"`
	FollowRedirects bool        `json:"followRedirects"`
	Body            string      `json:"body"`
	BodyType        string      `json:"bodyType,omitempty"`
	Headers         []KeyValue  `json:"headers"`
	QueryParameters []KeyValue  `json:"queryParameters"`
	Assertions      []Assertion `json:"assertions"`
	BasicAuth       BasicAuth   `json:"basicAuth,omitempty"`
}

// Assertion represents an assertion about an API response, which will be
// verified as part of the check.
type Assertion struct {
	Edit          bool   `json:"edit"`
	Order         int    `json:"order"`
	ArrayIndex    int    `json:"arrayIndex"`
	ArraySelector int    `json:"arraySelector"`
	Source        string `json:"source"`
	Property      string `json:"property"`
	Comparison    string `json:"comparison"`
	Target        string `json:"target"`
}

// BasicAuth represents the HTTP basic authentication credentials for a request.
type BasicAuth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// KeyValue represents a key-value pair, for example a request header setting,
// or a query parameter.
type KeyValue struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Locked bool   `json:"locked"`
}

// EnvironmentVariable represents a key-value pair for setting environment
// values during check execution.
type EnvironmentVariable struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Locked bool   `json:"locked"`
}

// AlertSettings represents an alert configuration.
type AlertSettings struct {
	EscalationType      string              `json:"escalationType,omitempty"`
	RunBasedEscalation  RunBasedEscalation  `json:"runBasedEscalation,omitempty"`
	TimeBasedEscalation TimeBasedEscalation `json:"timeBasedEscalation,omitempty"`
	Reminders           Reminders           `json:"reminders,omitempty"`
	SSLCertificates     SSLCertificates     `json:"sslCertificates,omitempty"`
}

// RunBasedEscalation represents an alert escalation based on a number of failed
// check runs.
type RunBasedEscalation struct {
	FailedRunThreshold int `json:"failedRunThreshold,omitempty"`
}

// TimeBasedEscalation represents an alert escalation based on the number of
// minutes after a check first starts failing.
type TimeBasedEscalation struct {
	MinutesFailingThreshold int `json:"minutesFailingThreshold,omitempty"`
}

// Reminders represents the number of reminders to send after an alert
// notification, and the time interval between them.
type Reminders struct {
	Amount   int `json:"amount,omitempty"`
	Interval int `json:"interval,omitempty"`
}

// SSLCertificates represents alert settings for expiring SSL certificates.
type SSLCertificates struct {
	Enabled        bool `json:"enabled"`
	AlertThreshold int  `json:"alertThreshold"`
}

// AlertChannel represents an alert channel and its subscribed checks. The API
// defines this data as read-only.
type AlertChannel struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty"`
	CreatedAt time.Time              `json:"created_at,omitempty"`
	UpdatedAt time.Time              `json:"updated_at,omitempty"`
}

// Subscription represents a subscription to an alert channel. The API defines
// this data as read-only.
type Subscription struct {
	ID             string `json:"id,omitempty"`
	CheckID        string `json:"checkId,omitempty"`
	AlertChannelID int64  `json:"alertChannelId,omitempty"`
	Activated      bool   `json:"activated"`
}
