package model

import (
	"github.com/bluele/factory-go/factory"
	fake "github.com/brianvoe/gofakeit/v5"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
)

// NewAuth is the function that returns an instance of the auth's domain model for performing tests.
func NewAuth(args map[string]interface{}) domainmodel.Chunk {

	chunkFactory := factory.NewFactory(
		domainmodel.Chunk{},
	).Attr("Offset", func(fArgs factory.Args) (interface{}, error) {
		offset := fake.Int64()
		if val, ok := args["offset"]; ok {
			offset = val.(int64)
		}

		return offset, nil
	}).Attr("Limit", func(fArgs factory.Args) (interface{}, error) {
		limit := fake.Int64()
		if val, ok := args["limit"]; ok {
			limit = val.(int64)
		}

		return limit, nil
	}).Attr("BytesRead", func(fArgs factory.Args) (interface{}, error) {
		bytesRead := fake.Int64()
		if val, ok := args["bytesRead"]; ok {
			bytesRead = val.(int64)
		}

		return bytesRead, nil
	}).Attr("BytestreamToString", func(fArgs factory.Args) (interface{}, error) {
		bytestreamToString := fake.Name()
		if val, ok := args["byteStreamToString"]; ok {
			bytestreamToString = val.(string)
		}

		return bytestreamToString, nil
	})

	return chunkFactory.MustCreate().(domainmodel.Auth)
}
