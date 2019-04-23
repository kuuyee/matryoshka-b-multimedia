package multimedia

import (
	"fmt"
	"net/http"
)

// HTTPCodeError is returned if remote returned an unexpected status code
type HTTPCodeError struct {
	Resp *http.Response
}

func (e HTTPCodeError) Error() string {
	return fmt.Sprintf("remote returned %d", e.Resp.StatusCode)
}
