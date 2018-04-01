package matcher

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zalando/skipper/eskip"
	"github.com/zalando/skipper/eskipfile"
	"github.com/zalando/skipper/filters"
	"github.com/zalando/skipper/filters/builtin"
	"github.com/zalando/skipper/filters/filtertest"
	"github.com/zalando/skipper/loadbalancer"
	"github.com/zalando/skipper/logging/loggingtest"
	"github.com/zalando/skipper/predicates/cookie"
	"github.com/zalando/skipper/predicates/interval"
	"github.com/zalando/skipper/predicates/query"
	"github.com/zalando/skipper/predicates/source"
	"github.com/zalando/skipper/predicates/traffic"
	"github.com/zalando/skipper/routing"
)

// Matcher helps testing eskip routing logic
type Matcher interface {
	// Given request attributes test if a route matches
	Test(attributes *RequestAttributes) TestResult
}

// TestResult result of a Matcher.Test operation
type TestResult interface {
	// Matching route if there was match nil if no match
	Route() *eskip.Route
	// The http request that was used to perform the test
	Request() *http.Request
	// Normalized request attributes after test
	Attributes() *RequestAttributes
	// Nice string representation
	PrettyPrint() string
	// Nice string representation line by line
	PrettyPrintLines() []string
}

// RequestAttributes represents the http request attributes to test
type RequestAttributes struct {
	Method  string
	Path    string
	Headers map[string]string
}

type matcher struct {
	routing *routing.Routing
}

type testResult struct {
	route      *eskip.Route
	req        *http.Request
	attributes *RequestAttributes
}

func (t *testResult) Route() *eskip.Route {
	return t.route
}

func (t *testResult) Request() *http.Request {
	return t.req
}

func (t *testResult) Attributes() *RequestAttributes {
	return t.attributes
}

// PrettyPrint return a nice string output representing the result
func (t *testResult) PrettyPrint() string {
	out := t.PrettyPrintLines()
	return strings.Join(out, "\n")
}

// PrettyPrintLines like PrettyPrint but return a list of strings
// representing each line that forms the final output
func (t *testResult) PrettyPrintLines() []string {
	attrs := t.Attributes()
	out := []string{}
	out = append(out, fmt.Sprintf("request: %s %s", attrs.Method, attrs.Path))
	if len(attrs.Headers) > 0 {
		pairs := make([]string, 0, len(attrs.Headers))
		for key, value := range attrs.Headers {
			pairs = append(pairs, fmt.Sprintf(`"%s"="%s"`, key, value))
		}
		out = append(out, fmt.Sprintf("request headers: %s", strings.Join(pairs, ", ")))
	}

	route := t.Route()
	if route != nil {
		out = append(out, fmt.Sprintf("matching route id: %s", route.Id))
		out = append(out, fmt.Sprintf("matching route:\n```%s```", t.prettyPrintRoute()))
	}
	return out
}

// prettyPrintRoute return a nice string representation of the resulting route if any
func (t *testResult) prettyPrintRoute() string {
	def := t.route.Print(eskip.PrettyPrintInfo{
		Pretty:    true,
		IndentStr: "  ",
	})
	return fmt.Sprintf("%s: %s\n", t.route.Id, def)
}

// Options when creating a NewMatcher
type Options struct {
	// Path to a .eskip file defining routes
	RoutesFile string

	// CustomPredicates list of of custom Skipper predicate specs
	CustomPredicates []routing.PredicateSpec

	// CustomFilters lister of custom Skipper filter specs
	CustomFilters []filters.Spec

	// MockFilters list of custom Skipper filters to mock by name
	MockFilters []string

	// IgnoreTrailingSlash Skipper option
	IgnoreTrailingSlash bool

	// Verbose verbose debug output
	Verbose bool
}

// New create a new Matcher
func New(o *Options) (Matcher, error) {
	// creates data clients
	dataClients, err := createDataClients(o.RoutesFile)

	if err != nil {
		return nil, err
	}

	routing := createRouting(dataClients, o)

	return &matcher{
		routing,
	}, nil
}

// Test check if incoming request attributes are matching any eskip route
// Return is nil if there isn't a match
func (f *matcher) Test(attributes *RequestAttributes) TestResult {
	req, _ := createHTTPRequest(attributes)

	// find a match
	route, _ := f.routing.Route(req)
	var eroute eskip.Route

	if route != nil {
		eroute = route.Route
	}

	if eroute.Id == "" {
		return &testResult{
			nil,
			req,
			attributes,
		}
	}

	result := &testResult{
		&eroute,
		req,
		attributes,
	}

	// transform literal to pointer to use eskip.Route methods
	return result
}

func createHTTPRequest(attributes *RequestAttributes) (*http.Request, error) {
	if strings.HasPrefix(attributes.Path, "/") == false {
		attributes.Path = "/" + attributes.Path
	}

	u, _ := url.Parse("http://localhost" + attributes.Path)
	if attributes.Method == "" {
		attributes.Method = "GET"
	}

	httpReq := &http.Request{
		Method: strings.ToUpper(attributes.Method),
		URL:    u,
		Header: make(http.Header),
	}
	for key, value := range attributes.Headers {
		httpReq.Header.Set(key, value)
	}

	return httpReq, nil
}

func createRouting(dataClients []routing.DataClient, o *Options) *routing.Routing {
	l := loggingtest.New()

	if o.Verbose == true {
		l.Unmute() // unmute skipper logging
	}

	// create a filter registry with the available filter specs registered,
	// and register the mock and custom filters
	registry := builtin.MakeRegistry()
	customFilters := append(mockFilters(o.MockFilters), o.CustomFilters...)
	for _, f := range customFilters {
		registry.Register(f)
	}

	// create routing
	// create the proxy instance
	var mo routing.MatchingOptions
	if o.IgnoreTrailingSlash {
		mo = routing.IgnoreTrailingSlash
	}

	// include bundled custom predicates
	o.CustomPredicates = append(o.CustomPredicates,
		source.New(),
		source.NewFromLast(),
		interval.NewBetween(),
		interval.NewBefore(),
		interval.NewAfter(),
		cookie.New(),
		query.New(),
		traffic.New(),
		loadbalancer.NewGroup(),
		loadbalancer.NewMember(),
	)

	routingOptions := routing.Options{
		DataClients:     dataClients,
		Log:             l,
		FilterRegistry:  registry,
		MatchingOptions: mo,
		Predicates:      o.CustomPredicates,
	}

	router := routing.New(routingOptions)
	defer router.Close()

	// wait for "route settings applied"
	time.Sleep(120 * time.Millisecond)

	return router
}

func createDataClients(path string) ([]routing.DataClient, error) {
	client, err := eskipfile.Open(path)
	if err != nil {
		return nil, err
	}
	DataClients := []routing.DataClient{
		client,
	}
	return DataClients, nil
}

// mockFilters creates a list of mocked filters givane a list of filterNames
func mockFilters(filterNames []string) []filters.Spec {
	fs := make([]filters.Spec, len(filterNames))
	for i, filterName := range filterNames {
		fs[i] = &filtertest.Filter{
			FilterName: filterName,
		}
	}
	return fs
}
