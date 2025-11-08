package redlock

import (
	"context"
	"silo/pkg/database/interfaces"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

//go:generate mockgen -source=redlock.go -destination=mocks/redlock.go

// Redlock lock creator
type Redlock interface {
	NewMutex(name string, options ...*Options) Mutex
	AsGoCronLocker(options ...*Options) Locker
	// NewReentrantMutex(name string, options ...*Options) ReentrantMutex
}

// A Mutex is a distributed mutual exclusion lock.
type Mutex interface {
	// Name returns mutex name (i.e. the Redis key).
	Name() string
	// // Value returns the current random value. The value will be empty until a lock is acquired (or WithValue option is used).
	// Value() string
	// // Until returns the time of validity of acquired lock. The value will be zero value until a lock is acquired.
	// Until() time.Time
	// TryLock only attempts to lock m once and returns immediately regardless of success or failure without retrying.
	TryLock() error
	// TryLockContext only attempts to lock m once and returns immediately regardless of success or failure without retrying.
	TryLockContext(ctx context.Context) error
	// Lock locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
	Lock() error
	// LockContext locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
	LockContext(ctx context.Context) error
	// Unlock unlocks m and returns the status of unlock.
	Unlock() (bool, error)
	// UnlockContext unlocks m and returns the status of unlock.
	UnlockContext(ctx context.Context) (bool, error)
	// Extend resets the mutex's expiry and returns the status of expiry extension.
	Extend() (bool, error)
	// ExtendContext resets the mutex's expiry and returns the status of expiry extension.
	ExtendContext(ctx context.Context) (bool, error)
	LockFunc(ctx context.Context, f func(ctx context.Context) error) error
}

// Options .
type Options struct {
	Expiry time.Duration // The default is 8s.
	Rries  int           // The default value is 32.
	// WithFailFast can be used to quickly acquire and release the lock.
	// When some Redis servers are blocking, we do not need to wait for responses from all the Redis servers response.
	// As long as the quorum is met, we can assume the lock is acquired. The effect of this parameter is to achieve low
	// latency, avoid Redis blocking causing Lock/Unlock to not return for a long time.
	FailFast bool
	// WithSetNXOnExtend improves extending logic to extend the key if exist
	// and if not, tries to set a new key in redis
	// Useful if your redises restart often and you want to reduce the chances of losing the lock
	SetNXOnExtend bool
}

func WithExpiry(expiry time.Duration) *Options {
	return &Options{
		Expiry: expiry,
	}
}

func WithFailFast(failFast bool) *Options {
	return &Options{
		FailFast: failFast,
	}
}

func WithSetNXOnExtend(setNXOnExtend bool) *Options {
	return &Options{
		SetNXOnExtend: setNXOnExtend,
	}
}

var redisLock *RedisLock
var redisLockOnce sync.Once

// NewRedlock create a redlock instance
func NewRedlock(cacheClient interfaces.CacheClient) Redlock {
	redisLockOnce.Do(func() {
		pool := goredis.NewPool(cacheClient)
		rs := redsync.New(pool)

		redisLock = &RedisLock{
			rs:          rs,
			cacheClient: cacheClient,
			lockCount:   &sync.Map{},
		}
	})

	return redisLock
}

// RedisLock .
type RedisLock struct {
	rs          *redsync.Redsync
	cacheClient interfaces.CacheClient
	lockCount   *sync.Map
	lockMutex   *sync.RWMutex
}

// NewMutex create a mutex for lock
func (r *RedisLock) NewMutex(name string, ops ...*Options) Mutex {
	expiry := 8 * time.Second
	options := []redsync.Option{}
	for _, o := range ops {
		if o.Expiry > 0 {
			expiry = o.Expiry
			options = append(options, redsync.WithExpiry(o.Expiry))
		}
		if o.Rries > 0 {
			options = append(options, redsync.WithTries(o.Rries))
		}
		if o.FailFast {
			options = append(options, redsync.WithFailFast(o.FailFast))
		}
		if o.SetNXOnExtend {
			options = append(options, redsync.WithSetNXOnExtend())
		}
	}

	// apply key prefix
	mutex := r.rs.NewMutex(name, options...)

	return newRedMutex(mutex, expiry)
}

// AsGoCronLocker locker is for gocron
func (r *RedisLock) AsGoCronLocker(options ...*Options) Locker {
	return &RedisLocker{
		r:   r,
		ops: options,
	}
}

// type ReentrantMutex interface {
// 	ReentrantLock() error
// 	ReentrantLockContext(ctx context.Context) error
// 	UnReentrantLock() (bool, error)
// 	UnReentrantLockContext(ctx context.Context) (bool, error)
// }

// func (r *RedisLock) NewReentrantMutex(name string, ops ...*Options) ReentrantMutex {
// 	options := []redsync.Option{}
// 	for _, o := range ops {
// 		if o.Expiry > 0 {
// 			options = append(options, redsync.WithExpiry(o.Expiry))
// 		}
// 		if o.Rries > 0 {
// 			options = append(options, redsync.WithTries(o.Rries))
// 		}
// 		if o.FailFast {
// 			options = append(options, redsync.WithFailFast(o.FailFast))
// 		}
// 		if o.SetNXOnExtend {
// 			options = append(options, redsync.WithSetNXOnExtend())
// 		}
// 	}

// 	mutex := r.rs.NewMutex(name, options...)

// 	return &ReentrantMutexImpl{
// 		Mutex:     mutex,
// 		lockCount: r.lockCount,
// 		lockMutex: r.lockMutex,
// 	}
// }

// type ReentrantMutexImpl struct {
// 	Mutex
// 	lockCount *sync.Map
// 	lockMutex *sync.RWMutex
// }

// func (rm *ReentrantMutexImpl) ReentrantLock() error {
// 	rm.lockMutex.Lock()
// 	defer rm.lockMutex.Unlock()

// 	count, _ := rm.lockCount.LoadOrStore(rm.Name(), int32(0))
// 	if atomic.LoadInt32(count.(*int32)) > 0 {
// 		atomic.AddInt32(count.(*int32), 1)
// 		return nil
// 	}

// 	if err := rm.Lock(); err != nil {
// 		return err
// 	}

// 	atomic.StoreInt32(count.(*int32), 1)
// 	return nil
// }

// func (rm *ReentrantMutexImpl) ReentrantLockContext(ctx context.Context) error {
// 	rm.lockMutex.Lock()
// 	defer rm.lockMutex.Unlock()

// 	count, _ := rm.lockCount.LoadOrStore(rm.Name(), int32(0))
// 	if atomic.LoadInt32(count.(*int32)) > 0 {
// 		atomic.AddInt32(count.(*int32), 1)
// 		return nil
// 	}

// 	if err := rm.TryLockContext(ctx); err != nil {
// 		return err
// 	}

// 	atomic.StoreInt32(count.(*int32), 1)
// 	return nil
// }

// func (rm *ReentrantMutexImpl) UnReentrantLock() (bool, error) {
// 	rm.lockMutex.Lock()
// 	defer rm.lockMutex.Unlock()

// 	count, _ := rm.lockCount.Load(rm.Name())
// 	if atomic.LoadInt32(count.(*int32)) > 1 {
// 		atomic.AddInt32(count.(*int32), -1)
// 		return true, nil
// 	}

// 	success, err := rm.Mutex.Unlock()
// 	if err == nil {
// 		rm.lockCount.Delete(rm.Name())
// 	}
// 	return success, err
// }

// func (rm *ReentrantMutexImpl) UnReentrantLockContext(ctx context.Context) (bool, error) {
// 	rm.lockMutex.Lock()
// 	defer rm.lockMutex.Unlock()

// 	count, _ := rm.lockCount.Load(rm.Name())
// 	if atomic.LoadInt32(count.(*int32)) > 1 {
// 		atomic.AddInt32(count.(*int32), -1)
// 		return true, nil
// 	}

// 	success, err := rm.Mutex.UnlockContext(ctx)
// 	if err == nil {
// 		rm.lockCount.Delete(rm.Name())
// 	}
// 	return success, err
// }
