package redlock

import (
	"context"
	"errors"
	"log/slog"
	"silo/pkg/errs"
	"silo/pkg/utils"
	"time"

	"github.com/go-redsync/redsync/v4"
)

func newRedMutex(mutex *redsync.Mutex, expiry time.Duration) *RedMutex {
	return &RedMutex{
		mutex:  mutex,
		expiry: expiry,
	}
}

type RedMutex struct {
	mutex  *redsync.Mutex
	expiry time.Duration
}

func (m *RedMutex) Name() string {
	return m.mutex.Name()
}

// // Value returns the current random value. The value will be empty until a lock is acquired (or WithValue option is used).
// Value() string
// // Until returns the time of validity of acquired lock. The value will be zero value until a lock is acquired.
// Until() time.Time
// TryLock only attempts to lock m once and returns immediately regardless of success or failure without retrying.
func (m *RedMutex) TryLock() error {
	return errorHandler(m.mutex.TryLock())
}

// TryLockContext only attempts to lock m once and returns immediately regardless of success or failure without retrying.
func (m *RedMutex) TryLockContext(ctx context.Context) error {
	return errorHandler(m.mutex.TryLockContext(ctx))
}

// Lock locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *RedMutex) Lock() error {
	return errorHandler(m.mutex.Lock())
}

// LockContext locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *RedMutex) LockContext(ctx context.Context) error {
	return errorHandler(m.mutex.LockContext(ctx))
}

// Unlock unlocks m and returns the status of unlock.
func (m *RedMutex) Unlock() (bool, error) {
	return m.mutex.Unlock()
}

// UnlockContext unlocks m and returns the status of unlock.
func (m *RedMutex) UnlockContext(ctx context.Context) (bool, error) {
	return m.mutex.UnlockContext(ctx)
}

// Extend resets the mutex's expiry and returns the status of expiry extension.
func (m *RedMutex) Extend() (bool, error) {
	return m.mutex.Extend()
}

// ExtendContext resets the mutex's expiry and returns the status of expiry extension.
func (m *RedMutex) ExtendContext(ctx context.Context) (bool, error) {
	return m.mutex.ExtendContext(ctx)
}

// LockFunc locks the mutex and runs the function.
func (m *RedMutex) LockFunc(ctx context.Context, f func(ctx context.Context) error) (err error) {
	if err := m.LockContext(ctx); err != nil {
		return err
	}

	defer m.Unlock()

	// Create a context with cancel to properly cleanup goroutine
	extendCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	ticker := time.NewTicker(m.expiry / 4)
	defer ticker.Stop()

	done := false

	// Goroutine to extend lock periodically
	go func() {
		for {
			select {
			case <-ticker.C:
				if done {
					return
				}

				if _, err := m.mutex.ExtendContext(extendCtx); err != nil {
					if done {
						return
					}
					if errors.Is(err, context.Canceled) {
						return
					}
					slog.Error("[redlock] Failed to extend lock", "name", m.Name(), "err", err)
					return
				}
			case <-extendCtx.Done():
				return
			}
		}
	}()

	// Run the job with error handling
	err = utils.PContextCall(extendCtx, f)
	if err != nil {
		return err
	}

	done = true

	return nil
}

func errorHandler(err error) error {
	if err == nil {
		return nil
	}

	switch err.(type) {
	case *redsync.ErrTaken, redsync.ErrTaken:
		return errs.NewCustomError(errs.StatusTryLockError, err.Error())
	case *redsync.ErrNodeTaken, redsync.ErrNodeTaken:
		return errs.NewCustomError(errs.StatusTryLockError, err.Error())
	}

	if errors.Is(err, redsync.ErrFailed) {
		return errs.NewCustomError(errs.StatusTryLockError, err.Error())
	}

	return errs.NewInternalServerError(err.Error())
}
