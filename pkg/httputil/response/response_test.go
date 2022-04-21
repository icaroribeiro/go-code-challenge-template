package response_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	errorhttputil "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/error"
	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestResponseUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestRespondWithJson() {
	res := &httptest.ResponseRecorder{}
	statusCode := 0
	var payload interface{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRespondingWithOKAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusOK
				text := "everything is up and running"
				payload = messagehttputilpkg.Message{Text: text}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailInRespondingAndJsonIfItIsNotPossibleToGetJsonEncodingOfPayload",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusInternalServerError
				payload = func() {
					customerror.New("failed")
				}
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			responsehttputilpkg.RespondWithJson(res, statusCode, payload)

			if !tc.WantError {
				assert.Equal(t, res.Result().Header.Get("Content-Type"), "application/json")
				assert.Equal(t, statusCode, res.Result().StatusCode)
				message := messagehttputilpkg.Message{}
				err := json.NewDecoder(res.Body).Decode(&message)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.Equal(t, payload, message)
			} else {
				assert.Equal(t, statusCode, res.Result().StatusCode)
			}
		})
	}
}

func (ts *TestSuite) TestRespondErrorAndJson() {
	res := &httptest.ResponseRecorder{}
	statusCode := 0
	err := customerror.New("failed")
	payload := errorhttputil.Error{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRespondingWithInternalServerErrorAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusInternalServerError
				text := "failed"
				err = customerror.New(text)
				payload = errorhttputil.Error{Text: text}
			},
		},
		{
			Context: "ItShouldSucceedInRespondingWithBadRequestAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusBadRequest
				text := "failed"
				err = customerror.BadRequest.New(text)
				payload = errorhttputil.Error{Text: text}
			},
		},
		{
			Context: "ItShouldSucceedInRespondingWithUnauthorizedAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusUnauthorized
				text := "failed"
				err = customerror.Unauthorized.New(text)
				payload = errorhttputil.Error{Text: text}
			},
		},
		{
			Context: "ItShouldSucceedInRespondingWithNotFoundAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusNotFound
				text := "failed"
				err = customerror.NotFound.New(text)
				payload = errorhttputil.Error{Text: text}
			},
		},
		{
			Context: "ItShouldSucceedInRespondingWithConflictAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusConflict
				text := "failed"
				err = customerror.Conflict.New(text)
				payload = errorhttputil.Error{Text: text}
			},
		},
		{
			Context: "ItShouldSucceedInRespondingWithUnprocessableEntityAndJsonBody",
			SetUp: func(t *testing.T) {
				res = httptest.NewRecorder()
				statusCode = http.StatusUnprocessableEntity
				text := "failed"
				err = customerror.UnprocessableEntity.New(text)
				payload = errorhttputil.Error{Text: text}
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			responsehttputilpkg.RespondErrorWithJson(res, err)

			assert.Equal(t, res.Result().Header.Get("Content-Type"), "application/json")
			assert.Equal(t, statusCode, res.Result().StatusCode)
			errMessage := errorhttputil.Error{}
			err := json.NewDecoder(res.Body).Decode(&errMessage)
			assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			assert.Equal(t, payload, errMessage)
		})
	}
}
