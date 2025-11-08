package interfaces

//go:generate mockgen -source=entity.go -destination=mocks/entity.go

// Entity .
type Entity interface {
	IsNew() bool
	PopAfterCommitHooks() []func()
}
