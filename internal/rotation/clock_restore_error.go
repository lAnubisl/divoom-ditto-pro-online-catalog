package rotation

type clockRestoreError struct{ err error }

func (e *clockRestoreError) Error() string { return "clock restoration failed: " + e.err.Error() }

func (e *clockRestoreError) Unwrap() error { return e.err }
