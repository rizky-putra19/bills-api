package dynamox

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// Config defines the gosqs configuration
// private key to access dynamox
// secret to access dynamox
// region for dynamox and used for determining the topic ARN
// provided automatically by dynamox, but must be set for emulators or local testing
// account ID of the dynamox account, used for determining the topic ARN
// environment name, used for determinig the topic ARN
type Config struct {
	Key            string
	Secret         string
	Region         string
	Hostname       string
	AWSAccountID   string
	Env            string
	RetryCount     int
	OrderTableName string
	VATableName    string
	TableSuffix    string
}

type dataType string

func (dt dataType) String() string {
	return string(dt)
}

// DataTypeNumber represents the Number datatype, use it when creating custom attributes
const DataTypeNumber = dataType("Number")

// DataTypeString represents the String datatype, use it when creating custom attributes
const DataTypeString = dataType("String")

type retryer struct {
	client.DefaultRetryer
	retryCount int
}

func NewDynamo(c Config) (*dynamo.DB, error) {
	sess, err := newSession(c)
	if err != nil {
		return nil, err
	}

	return dynamo.New(sess), nil
}

// MaxRetries sets the total exponential back off attempts to 10 retries
func (r retryer) MaxRetries() int {
	if r.retryCount > 0 {
		return r.retryCount
	}

	return 10
}

// newSession creates a new dynamox session
func newSession(c Config) (*session.Session, error) {
	r := &retryer{retryCount: c.RetryCount}

	cfg := request.WithRetryer(aws.NewConfig().WithRegion(c.Region), r)

	// This will set the default AWS URL to a hostname of your choice. Perfect for testing, or mocking functionality
	if c.Env == "local" {
		cfg.Endpoint = &c.Hostname
	}

	return session.NewSession(cfg)
}
