package processor

import (
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 60 * time.Second,
}
