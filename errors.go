package dvb

import "fmt"

type apiError struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}
