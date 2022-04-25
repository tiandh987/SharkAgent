package store

// Factory defines the iam platform storage interface.
type Factory interface {
	Users() UserStore
	Close() error
}