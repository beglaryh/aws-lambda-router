package context

type Context struct {
	Body            string
	PathParameters  map[string]string
	QueryParameters map[string]string
	Headers         map[string]string
}
