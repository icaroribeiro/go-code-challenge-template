package dbtrx

import (
	"context"
	"log"
	"net/http"

	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	"gorm.io/gorm"
)

var dbTrxCtxKey = &contextKey{"db_trx"}

type contextKey struct {
	name string
}

var (
	statusCodesList = []int{http.StatusOK}
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured.
type responseWriter struct {
	http.ResponseWriter
	wroteHeader bool
	statusCode  int
}

// NewContext is the function that returns a new Context that carries db_trx value.
func NewContext(ctx context.Context, dbTrx *gorm.DB) context.Context {
	return context.WithValue(ctx, dbTrxCtxKey, dbTrx)
}

// FromContext is the function that returns the db_trx value stored in context, if any.
func FromContext(ctx context.Context) (*gorm.DB, bool) {
	raw, ok := ctx.Value(dbTrxCtxKey).(*gorm.DB)
	return raw, ok
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

// DBTrx is the function that wraps a http.Handler to enable using a database transaction during an API incoming request.
func DBTrx(db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if db == nil {
				next.ServeHTTP(w, r)
				return
			}

			dbTrx := db.Begin()
			defer func() {
				if r := recover(); r != nil {
					var err error
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = customerror.Newf("%v", r)
					}
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("Transaction is being rolled back: %s \n", err.Error())
					dbTrx.Rollback()
					return
				}
			}()

			// It is necessary to set database transaction that can be used for performing operations with transaction.
			ctx := NewContext(r.Context(), dbTrx)
			r = r.WithContext(ctx)

			wrapped := wrapResponseWriter(w)

			next.ServeHTTP(wrapped, r)

			if isStatusCodeInList(wrapped.Status(), statusCodesList) {
				if err := dbTrx.Commit().Error; err != nil {
					log.Printf("failed to commit database transaction: %s", err.Error())
				}
			} else {
				log.Printf("database transaction is being rolled back due to status code: %d", wrapped.statusCode)
				if err := dbTrx.Rollback().Error; err != nil {
					log.Printf("failed to rollback database transaction: %s", err.Error())
				}
			}
		}
	}
}
