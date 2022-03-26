package dbtrx

import (
	"context"
	"log"
	"net/http"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
	"gorm.io/gorm"
)

var (
	statusCodesList = []int{http.StatusOK}
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured.
type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.statusCode
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if rw.wroteHeader {
		return
	}

	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.wroteHeader = true
}

func isStatusCodeInList(statusCode int, statusCodeList []int) bool {
	for _, value := range statusCodeList {
		if statusCode == value {
			return true
		}
	}

	return false
}

// DBTrx is the function that  wraps a http.Handler to enable using a database transaction during an API incoming request.
func DBTrx(db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if db == nil {
				next.ServeHTTP(w, r)
			}

			dbTrx := db.Begin()
			defer func() {
				if r := recover(); r != nil {
					dbTrx.Rollback()
				}
			}()

			// It is necessary to set database transaction that can be used for performing operations with transaction.
			ctx := r.Context()
			var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
			ctx = context.WithValue(ctx, dbTrxKey, dbTrx)
			r = r.WithContext(ctx)

			wrapped := wrapResponseWriter(w)

			next.ServeHTTP(wrapped, r)

			if isStatusCodeInList(wrapped.statusCode, statusCodesList) {
				if err := dbTrx.Commit().Error; err != nil {
					errorMessage := "database transaction commit failed: "
					responsehttputilpkg.RespondErrorWithJson(w, customerror.Newf(errorMessage, err.Error()))
					return
				}
			} else {
				log.Printf("rolling back database transaction due to status code: %d", wrapped.statusCode)
				dbTrx.Rollback()
			}
		}
	}
}
