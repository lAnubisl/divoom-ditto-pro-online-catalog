package gallery

type Options struct {
	Email, Password string
	// Optional endpoint overrides support isolated HTTP tests.
	APIURL, FileURL string
}
