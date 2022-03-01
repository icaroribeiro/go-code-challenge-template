package datastore

// IDatastore interface is the datastore's contract.
type IDatastore interface {
	Close() error
}
