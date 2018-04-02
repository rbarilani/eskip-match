package matcher

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatcherError(t *testing.T) {
	_, err := New(&Options{
		RoutesFile: "",
	})
	assert.Error(t, err)
}

func TestMacherTestSetHeaders(t *testing.T) {
	routesFile, err := filepath.Abs("./testdata/routes.eskip")
	if err != nil {
		t.Error(err)
		return
	}
	tester, err := New(&Options{
		RoutesFile:  routesFile,
		MockFilters: []string{"customfilter"},
	})

	assert.NoError(t, err)

	res := tester.Test(&RequestAttributes{
		Method: "POST",
		Path:   "/foo",
		Headers: map[string]string{
			"X-Bar": "bar",
		},
	})

	assert.NotNil(t, res.Request())
	assert.Equal(t, "bar", res.Request().Header.Get("X-Bar"))
}

func TestRoutes(t *testing.T) {

	routesFile, err := filepath.Abs("./testdata/routes.eskip")
	if err != nil {
		t.Error(err)
		return
	}

	tester, err := New(&Options{
		RoutesFile:          routesFile,
		MockFilters:         []string{"customfilter"},
		Verbose:             true,
		IgnoreTrailingSlash: true,
	})

	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		routeID string
		nomatch bool
		attrs   []*RequestAttributes
	}{
		{
			routeID: "foo",
			attrs: []*RequestAttributes{
				{
					Method: "POST",
					Path:   "/foo",
				},
				{
					Method: "post",
					Path:   "/foo/1",
				},
			},
		},
		{
			routeID: "foo_get",
			attrs: []*RequestAttributes{
				{
					Path: "/foo",
				},
				{
					Path: "foo",
				},
			},
		},
		{
			routeID: "query_param",
			attrs: []*RequestAttributes{
				{
					Path: "/abdc?q=bar",
				},
			},
		},
		{
			routeID: "bar",
			attrs: []*RequestAttributes{
				{
					Path: "/bar",
				},
			},
		},
		{
			routeID: "foo_header",
			attrs: []*RequestAttributes{
				{
					Path: "/foo",
					Headers: map[string]string{
						"Accept": "application/json",
					},
				},
			},
		},
		{
			routeID: "customfilter",
			attrs: []*RequestAttributes{
				{
					Path: "/customfilter",
				},
			},
		},
		{
			routeID: "no-match",
			nomatch: true,
			attrs: []*RequestAttributes{
				{
					Path: "/blobblob",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.routeID, func(t *testing.T) {
			for _, a := range tt.attrs {
				result := tester.Test(a)

				route := result.Route()
				req := result.Request()
				attrs := result.Attributes()

				assert.NotNil(t, req)
				assert.NotNil(t, attrs)

				if route != nil {
					assert.Contains(t, result.PrettyPrint(), "matching")
				} else {
					assert.NotContains(t, result.PrettyPrint(), "matching")
				}

				if tt.nomatch == true && route != nil {
					t.Errorf("request: %s %s shouldn't match but matches route id: %s", req.Method, a.Path, route.Id)
					return
				}

				if tt.nomatch == true && route == nil {
					return
				}

				if route == nil {
					t.Errorf("expected route id to be '%s' but no match\n request: %s %s", tt.routeID, req.Method, a.Path)
				} else if route.Id != tt.routeID {
					t.Errorf("expected route id to be '%s' but got '%s'\n request: %s %s", tt.routeID, route.Id, req.Method, a.Path)
				}

			}
		})
	}
}

func Example() {
	m, err := New(&Options{
		RoutesFile: "./testdata/routes.eskip",
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	result := m.Test(&RequestAttributes{
		Method: "GET",
		Path:   "/bar",
		Headers: map[string]string{
			"Accept": "application/json",
		},
	})

	route := result.Route()

	if route != nil {
		fmt.Println(route.Id)
		// Output: bar
	}
}
