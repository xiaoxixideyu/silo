package redlock

import (
	"context"
	"silo/pkg/errs"

	"github.com/go-co-op/gocron/v2"
)

//go:generate mockgen -source=gocron.go -destination=mocks/gocron.go

// Locker represents the required interface to lock jobs when running multiple schedulers.
// The lock is held for the duration of the job's run, and it is expected that the
// locker implementation handles time splay between schedulers.
// The lock key passed is the job's name - which, if not set, defaults to the
// go function's name, e.g. "pkg.myJob" for func myJob() {} in pkg
type Locker interface {
	// Lock if an error is returned by lock, the job will not be scheduled.
	Lock(ctx context.Context, key string) (gocron.Lock, error)
}

// Lock represents an obtained lock. The lock is released after the execution of the job
// by the scheduler.
type Lock interface {
	Unlock(ctx context.Context) error
}

// LockerMutex .
type LockerMutex struct {
	mutex Mutex
}

// Unlock .
func (m *LockerMutex) Unlock(ctx context.Context) error {
	ok, err := m.mutex.UnlockContext(ctx)
	if err != nil {
		return err
	}
	if !ok {
		return errs.NewInternalServerErrorf("unlock key not exist, key %s", m.mutex.Name())
	}

	return nil
}

// RedisLocker .
type RedisLocker struct {
	r   *RedisLock
	ops []*Options
}

// Lock .
func (l *RedisLocker) Lock(ctx context.Context, key string) (gocron.Lock, error) {
	mutex := l.r.NewMutex(key, l.ops...)
	lock := &LockerMutex{mutex: mutex}

	return lock, lock.mutex.TryLockContext(ctx)
}
