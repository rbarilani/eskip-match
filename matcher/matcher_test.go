package matcher

import (
	"path/filepath"
	"testing"
)

type matcherTestScenario struct {
	expectedRouteID string
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
					Path: "/foo",
				},
				{
					Path: "/foo/1",
				},
			},
		},
		{
			expectedRouteID: "get_foo",
			reqAttributes: []*RequestAttributes{
				{
					Method: "GET",
					Path:   "/foo",
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
	}

	for _, s := range scenarios {
		t.Run(s.expectedRouteID, func(t *testing.T) {
			for _, a := range s.reqAttributes {
				route := tester.Test(a)
				if route == nil {
					t.Errorf("expected route id to be '%s' but no match\n request: %s", s.expectedRouteID, a.Path)
				} else if route.Id != s.expectedRouteID {
					t.Errorf("expected route id to be '%s' but got '%s'\n request: %s", s.expectedRouteID, route.Id, a.Path)
				}
			}
		})
	}
}
