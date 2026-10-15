package gallery

import (
	"fmt"
)

type APIError struct{ Code int }

func (e *APIError) Error() string { return fmt.Sprintf("gallery API result %d", e.Code) }
