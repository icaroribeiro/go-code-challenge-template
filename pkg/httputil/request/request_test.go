package request_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/mux"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestRequestUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestPrepareRequestBody() {
	var inputBody interface{}
	var reqBody io.Reader

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInPreparingRequestBodyIfInputBodyIsAnEmptyString",
			SetUp: func(t *testing.T) {
				inputBody = ""
				reqBody = nil
			},
		},
		{
			Context: "ItShouldSucceedInPreparingRequestBodyIfInputBodyIsAJsonStringWithEscapeSequencies",
			SetUp: func(t *testing.T) {
				inputBody = `
				{
					"testing":	"testing"
				}`
				reqBody = strings.NewReader(`{"testing":"testing"}`)
			},
		},
		{
			Context: "ItShouldSucceedInPreparingRequestBodyIfInputBodyIsAVariableSizedBufferOfBytes",
			SetUp: func(t *testing.T) {
				inputBody = new(bytes.Buffer)
				reqBody = new(bytes.Buffer)
			},
		},
		{
			Context: "ItShouldSucceedInPreparingRequestBodyIfInputBodyIsNeitherAStringNorAVariableSizedBufferOfBytes",
			SetUp: func(t *testing.T) {
				inputBody = nil
				reqBody = nil
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedBody := requesthttputilpkg.PrepareRequestBody(inputBody)

			assert.Equal(t, reqBody, returnedBody)
		})
	}
}

func (ts *TestSuite) TestSetRequestHeaders() {
	key := ""
	value := ""
	headers := make(map[string][]string)

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInNotSettingRequestHeadersIfHeadersMapIsEmpty",
			SetUp: func(t *testing.T) {
				key = "Content-Type"
				headers = map[string][]string{}
			},
		},
		{
			Context: "ItShouldSucceedInSettingRequestHeaders",
			SetUp: func(t *testing.T) {
				key = "Content-Type"
				value = "application/json"
				headers = map[string][]string{
					key: {value},
				}
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			route := routehttputilpkg.Route{
				Name:   "Testing",
				Method: http.MethodGet,
				Path:   "/testing",
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
			}

			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

			requesthttputilpkg.SetRequestHeaders(req, headers)

			sort.Strings(headers[key])
			sort.Strings(req.Header.Values(key))
			assert.Equal(t, headers[key], req.Header.Values(key))
		})
	}
}

const (
	_count requesthttputilpkg.ContextKeyType = "count"
)

func (ts *TestSuite) TestSetRequestContext() {
	var expectedCount int64
	statusCode := 0
	contextMap := make(map[interface{}]interface{})

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInNotSettingRequestContextIfContextMapIsEmpty",
			SetUp: func(t *testing.T) {
				statusCode = http.StatusInternalServerError
				contextMap = make(map[interface{}]interface{})
			},
		},
		{
			Context: "ItShouldSucceedInNotSettingRequestContextIfContextMapIsEmpty",
			SetUp: func(t *testing.T) {
				expectedCount = fake.Int64()
				statusCode = http.StatusInternalServerError
				contextMap[_count] = expectedCount
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			countHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
				i := r.Context().Value(_count)

				count, ok := i.(int)
				if !ok {
					responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
					return
				}

				assert.Equal(t, expectedCount, count)
			}

			route := routehttputilpkg.Route{
				Name:        "Testing",
				Method:      http.MethodGet,
				Path:        "/testing",
				HandlerFunc: countHandlerFunc,
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
			}

			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

			requesthttputilpkg.SetRequestContext(req, contextMap)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			assert.Equal(t, statusCode, resprec.Result().StatusCode)
		})
	}
}
