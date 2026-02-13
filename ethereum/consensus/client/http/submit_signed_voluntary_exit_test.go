//nolint:revive // package name intentionally reflects domain, not directory name
package eth2http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubmitSignedVoluntaryExit(t *testing.T) {
	u, err := url.Parse("http://localhost")
	require.NoError(t, err)
	var (
		errMsg        = "Invalid voluntary exit, it will never pass validation so it's rejected"
		errStatusCode = http.StatusBadRequest
		errBodyMsg    = fmt.Sprintf(`{"code": %d,"message": %q}`, errStatusCode, errMsg)
		respError     = &http.Response{StatusCode: errStatusCode, Body: io.NopCloser(bytes.NewReader([]byte(errBodyMsg))), Request: &http.Request{Method: http.MethodPost, URL: u}}
	)

	respMsg, err := inspectSubmitSignedVoluntaryExitResponse(respError)
	require.Error(t, err)
	require.Contains(t, err.Error(), errMsg)
	require.Nil(t, respMsg)

	var (
		okMsg        = "all right"
		okStatusCode = http.StatusOK
		okBodyMsg    = fmt.Sprintf(`{"code": %d,"message": %q}`, okStatusCode, okMsg)
	)

	respOK := &http.Response{StatusCode: okStatusCode, Body: io.NopCloser(bytes.NewReader([]byte(okBodyMsg))), Request: &http.Request{Method: http.MethodPost, URL: u}}
	respMsg, err = inspectSubmitSignedVoluntaryExitResponse(respOK)
	require.NoError(t, err)
	require.NotNil(t, respMsg)
	assert.Equal(t, okStatusCode, respMsg.Code)
	assert.Equal(t, okMsg, respMsg.Message)
}
