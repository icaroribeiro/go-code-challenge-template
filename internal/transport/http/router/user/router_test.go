package user_test

// import (
// 	"net/http"
// 	"reflect"
// 	"runtime"
// 	"testing"

// 	fake "github.com/brianvoe/gofakeit/v5"
// 	"github.com/gorilla/mux"
// 	"github.com/icaroribeiro/go-code-challenge-template/internal/interfaces/httputil"
// 	routerpkg "github.com/icaroribeiro/go-code-challenge-template/internal/interfaces/router"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// func TestRouter(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestGetInstance() {
// 	r := &mux.Router{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInGettingTheRouterInstance",
// 			SetUp: func(t *testing.T) {
// 				r = mux.NewRouter()
// 			},
// 			WantError: false,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			router := routerpkg.New()

// 			returnedRouter := router.GetInstance()

// 			if !tc.WantError {
// 				assert.Equal(t, r, returnedRouter)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestSetRoutes() {
// 	route := httputil.Route{}

// 	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 	})

// 	routes := make(httputil.Routes, 0)

// 	r := &mux.Router{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSettingTheSwaggerRoute",
// 			SetUp: func(t *testing.T) {
// 				route = httputil.Route{
// 					Name:        "Swagger",
// 					Method:      fake.RandomString([]string{"GET", "POST", "PUT", "DELETE"}),
// 					PathPrefix:  fake.Word(),
// 					HandlerFunc: handlerFunc,
// 				}
// 				r = mux.NewRouter()
// 				r.Name(route.Name).
// 					Methods(route.Method).
// 					PathPrefix(route.PathPrefix).
// 					HandlerFunc(route.HandlerFunc)
// 				routes = append(routes, route)
// 			},
// 			WantError: false,
// 			TearDown: func(t *testing.T) {
// 				routes = make(httputil.Routes, 0)
// 			},
// 		},
// 		{
// 			Context: "ItShouldSucceedInSettingAnyOtherRoute",
// 			SetUp: func(t *testing.T) {
// 				route = httputil.Route{
// 					Name:        "AnyOtherRoute",
// 					Method:      fake.RandomString([]string{"GET", "POST", "PUT", "DELETE"}),
// 					Path:        fake.Word(),
// 					HandlerFunc: handlerFunc,
// 				}
// 				r = mux.NewRouter()
// 				r.Name(route.Name).
// 					Methods(route.Method).
// 					Path(route.Path).
// 					HandlerFunc(route.HandlerFunc)
// 				routes = append(routes, route)
// 			},
// 			WantError: false,
// 			TearDown: func(t *testing.T) {
// 				routes = make(httputil.Routes, 0)
// 			},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			router := routerpkg.New()

// 			router.SetRoutes(routes)

// 			returnedRouter := router.GetInstance()

// 			if !tc.WantError {
// 				assert.Equal(t, r.GetRoute(route.Name), returnedRouter.GetRoute(route.Name))
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestSetNotFoundHandler() {
// 	r := &mux.Router{}

// 	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusNotFound)
// 	})

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSettingTheNotFoundHandler",
// 			SetUp: func(t *testing.T) {
// 				r = mux.NewRouter()
// 				r.NotFoundHandler = handlerFunc
// 			},
// 			WantError: false,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			router := routerpkg.New()

// 			router.SetNotFoundHandler(handlerFunc)

// 			returnedRouter := router.GetInstance()

// 			// It is necessary to compare function "equality".
// 			handler1 := runtime.FuncForPC(reflect.ValueOf(r.NotFoundHandler).Pointer()).Name()
// 			handler2 := runtime.FuncForPC(reflect.ValueOf(returnedRouter.NotFoundHandler).Pointer()).Name()

// 			if !tc.WantError {
// 				assert.Equal(t, handler1, handler2)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestSetMethodNotAllowedHandler() {
// 	r := &mux.Router{}

// 	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	})

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSettingTheMethodNotAllowedHandler",
// 			SetUp: func(t *testing.T) {
// 				r = mux.NewRouter()
// 				r.MethodNotAllowedHandler = handlerFunc
// 			},
// 			WantError: false,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			router := routerpkg.New()

// 			router.SetMethodNotAllowedHandler(handlerFunc)

// 			returnedRouter := router.GetInstance()

// 			// It is necessary to compare function "equality".
// 			handler1 := runtime.FuncForPC(reflect.ValueOf(r.MethodNotAllowedHandler).Pointer()).Name()
// 			handler2 := runtime.FuncForPC(reflect.ValueOf(returnedRouter.MethodNotAllowedHandler).Pointer()).Name()

// 			if !tc.WantError {
// 				assert.Equal(t, handler1, handler2)
// 			}
// 		})
// 	}
// }
