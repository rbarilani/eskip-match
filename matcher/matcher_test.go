package matcher

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

type testRoutesScenario struct {
	expectedRouteID string
	expectNoMatch   bool
	reqAttributes   []*RequestAttributes
}

func TestMatcherError(t *testing.T) {
	_, err := New(&Options{
		RoutesFile: "",
	})
	assert.Error(t, err)
}

func TestMacherTestSetHeaders(t *testing.T) {
	routesFile, err := filepath.Abs("./matcher_test.eskip")
	if err != nil {
		t.Error(err)
		return
	}
	tester, err := New(&Options{
		RoutesFile:    routesFile,
		CustomFilters: MockFilters([]string{"customfilter"}),
	})

	assert.NoError(t, err)

	res, err := tester.Test(&RequestAttributes{
		Method: "POST",
		Path:   "/foo",
		Headers: map[string]string{
			"X-Bar": "bar",
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, res.Request())

	assert.Equal(t, "bar", res.Request().Header.Get("X-Bar"))
}

func TestRoutes(t *testing.T) {

	routesFile, err := filepath.Abs("./matcher_test.eskip")
	if err != nil {
		t.Error(err)
		return
	}

	tester, err := New(&Options{
		RoutesFile:    routesFile,
		CustomFilters: MockFilters([]string{"customfilter"}),
		Verbose:       true,
	})

	if err != nil {
		t.Error(err)
		return
	}

	scenarios := []testRoutesScenario{
		{
			expectedRouteID: "foo",
			reqAttributes: []*RequestAttributes{
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
			expectedRouteID: "foo_get",
			reqAttributes: []*RequestAttributes{
				{
					Path: "/foo",
				},
				{
					Path: "foo",
				},
			},
		},
		{
			expectedRouteID: "query_param",
			reqAttributes: []*RequestAttributes{
				{
					Path: "/abdc?q=bar",
				},
			},
		},
		{
			expectedRouteID: "bar",
			reqAttributes: []*RequestAttributes{
				{
					Path: "/bar",
				},
			},
		},
		{
			expectedRouteID: "customfilter",
			reqAttributes: []*RequestAttributes{
				{
					Path: "/customfilter",
				},
			},
		},
		{
			expectedRouteID: "no-match",
			expectNoMatch:   true,
			reqAttributes: []*RequestAttributes{
				{
					Path: "/blobblob",
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.expectedRouteID, func(t *testing.T) {
			for _, a := range s.reqAttributes {
				result, err := tester.Test(a)
				if err != nil {
					t.Error(err)
					return
				}
				route := result.Route()
				req := result.Request()
				attrs := result.Attributes()

				assert.NotNil(t, req)
				assert.NotNil(t, attrs)

				if s.expectNoMatch == true && route != nil {
					t.Errorf("request: %s %s shouldn't match but matches route id: %s", req.Method, a.Path, route.Id)
					return
				}

				if s.expectNoMatch == true && route == nil {
					return
				}

				if route == nil {
					t.Errorf("expected route id to be '%s' but no match\n request: %s %s", s.expectedRouteID, req.Method, a.Path)
				} else if route.Id != s.expectedRouteID {
					t.Errorf("expected route id to be '%s' but got '%s'\n request: %s %s", s.expectedRouteID, route.Id, req.Method, a.Path)
				}

				if route != nil {
					assert.NotEmpty(t, result.PrettyPrintRoute())
				} else {
					assert.Empty(t, result.PrettyPrintRoute())
				}
			}
		})
	}
}
