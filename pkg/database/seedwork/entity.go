package seedwork

// Entity .
type Entity struct {
	_new             bool
	afterCommitHooks []func()
}

// MarkNew .
func (e *Entity) MarkNew() {
	e._new = true
}

// IsNew .
func (e *Entity) IsNew() bool {
	return e._new
}

// AfterCommit .
func (e *Entity) AfterCommit(hook func()) {
	e.afterCommitHooks = append(e.afterCommitHooks, hook)
}

// PopAfterCommitHooks .
func (e *Entity) PopAfterCommitHooks() []func() {
	hooks := e.afterCommitHooks
	e.afterCommitHooks = []func(){}
	return hooks
}
