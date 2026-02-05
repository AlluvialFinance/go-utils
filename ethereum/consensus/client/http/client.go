package eth2http

import (
	"context"
	"io"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/sirupsen/logrus"

	kilnhttp "github.com/kilnfi/go-utils/net/http"
	httppreparer "github.com/kilnfi/go-utils/net/http/preparer"
	"github.com/kilnfi/go-utils/tracing"
)

var silentLog = &logrus.Logger{
	Out:       io.Discard,
	Formatter: &logrus.TextFormatter{DisableTimestamp: true},
	Level:     logrus.PanicLevel,
}

// Client provides methods to connect to an Ethereum 2.0 Beacon chain node
type Client struct {
	client autorest.Sender

	logger logrus.FieldLogger
}

func NewClientFromClient(s autorest.Sender) *Client {
	return &Client{
		client: s,
		logger: silentLog, // Disabled (silent) logger by default
	}
}

// NewClient creates a client connecting to an Ethereum 2.0 Beacon chain node at given addr
func NewClient(cfg *Config) (*Client, error) {
	httpc, err := kilnhttp.NewClient(cfg.HTTP)
	if err != nil {
		return nil, err
	}

	c := NewClientFromClient(
		autorest.Client{
			Sender:           httpc,
			RequestInspector: httppreparer.WithBaseURL(cfg.Address),
		},
	)

	if cfg.DisableLog {
		return c, nil
	}

	c.SetLogger(logrus.StandardLogger())
	return c, nil
}

func (c *Client) Logger() logrus.FieldLogger {
	return c.logger
}

func (c *Client) SetLogger(logger logrus.FieldLogger) {
	c.logger = logger.WithField("component", "eth.consensus.client")
}

// do performs an HTTP request and logs the trace ID from the response if present.
func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if resp != nil {
		c.logResponseTraceID(resp)
	}
	return resp, err
}

// logResponseTraceID logs the trace ID from the response header if present.
func (c *Client) logResponseTraceID(resp *http.Response) {
	if traceID := resp.Header.Get(tracing.HeaderTraceID); traceID != "" {
		c.logger.WithField(tracing.FieldTraceID, traceID).Debug("received response from CL node")
	}
}

func newRequest(ctx context.Context) *http.Request {
	req, _ := http.NewRequestWithContext(ctx, "", "", http.NoBody)
	return req
}

func inspectResponse(resp *http.Response, msg interface{}) error {
	return autorest.Respond(
		resp,
		WithBeaconErrorUnlessOK(),
		autorest.ByUnmarshallingJSON(msg),
		autorest.ByClosing(),
	)
}
