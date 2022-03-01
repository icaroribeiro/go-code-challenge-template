package datastore

// IDatastore interface is the datastore's contract.
type IDatastore interface {
	SetDB(interface{})
	GetDB() interface{}
	CheckStatus() error
	Close() error
}
