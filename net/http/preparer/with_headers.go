//revive:disable-next-line:package-directory-mismatch
package httppreparer

import (
	"net/http"

	"github.com/Azure/go-autorest/autorest"
)

// WithHeaders returns a PrepareDecorator that sets the given headers on each request.
// The headers map is snapshotted at construction time.
func WithHeaders(headers map[string]string) autorest.PrepareDecorator {
	snapshot := make(map[string]string, len(headers))
	for k, v := range headers {
		snapshot[k] = v
	}
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				return r, err
			}
			for k, v := range snapshot {
				r.Header.Set(k, v)
			}
			return r, nil
		})
	}
}
