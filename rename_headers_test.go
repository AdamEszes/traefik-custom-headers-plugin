package traefik_custom_headers_plugin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		desc          string
		renames       []renameData
		reqHeader     http.Header
		expRespHeader http.Header
	}{
		{
			desc: "Should rename headers while keeping their values",
			renames: []renameData{
				{
					ExistingHeaderName: "Foo",
					NewHeaderName:      "bar",
				},
				{
					ExistingHeaderName: "Another-Foo",
					NewHeaderName:      "another-bar",
				},
			},
			reqHeader: map[string][]string{
				"Foo":         {"fooval", "fooval2"},
				"Another-Foo": {"another-fooval", "another-fooval2"},
			},
			expRespHeader: map[string][]string{
				"bar":         {"fooval", "fooval2"},
				"another-bar": {"another-fooval", "another-fooval2"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			config := &Config{
				RenameData: test.renames,
			}

			next := func(rw http.ResponseWriter, req *http.Request) {
				for k, v := range test.reqHeader {
					for _, h := range v {
						rw.Header().Add(k, h)
					}
				}
				rw.WriteHeader(http.StatusOK)
			}

			rewriteBody, err := New(context.Background(), http.HandlerFunc(next), config, "rewriteHeader")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			rewriteBody.ServeHTTP(recorder, req)
			for k, expected := range test.expRespHeader {
				values := recorder.Result().Header[k]

				if !testEq(values, expected) {
					t.Errorf("Slice arent equals: expect: %+v, result: %+v", expected, values)
				}
			}
		})
	}
}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
