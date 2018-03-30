package matcher

import (
	"path/filepath"
	"testing"
)

type matcherTestScenario struct {
	expectedRouteID string
	expectNoMatch   bool
	reqAttributes   []*RequestAttributes
}

func TestRoutes(t *testing.T) {

	routesFile, err := filepath.Abs("./matcher_test.eskip")
	if err != nil {
		t.Error(err)
		return
	}

	tester, err := New(&Options{
		RoutesFile: routesFile,
		Verbose:    true,
	})

	if err != nil {
		t.Error(err)
		return
	}

	scenarios := []matcherTestScenario{
		{
			expectedRouteID: "foo",
			reqAttributes: []*RequestAttributes{
				{
					Method: "POST",
					Path:   "/foo",
				},
				{
					Method: "POST",
					Path:   "/foo/1",
				},
			},
		},
		{
			expectedRouteID: "foo_get",
			reqAttributes: []*RequestAttributes{
				{
					Method: "GET",
					Path:   "/foo",
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
			expectNoMatch: true,
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

			}
		})
	}
}
