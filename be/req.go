package be

import (
	"net/http"
)

type Req struct {
	*http.Request
}
