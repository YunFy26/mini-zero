package errorx

import (
	"errors"
	"sync"
)

type BatchError struct {
	errs []error
	lock sync.RWMutex
}

func (be *BatchError) Add(errs ...error) {
	be.lock.Lock()
	defer be.lock.Unlock()
	for _, err := range errs {
		if err != nil {
			be.errs = append(be.errs, err)
		}
	}
}

func (be *BatchError) Err() error {
	be.lock.RLock()
	defer be.lock.RUnlock()
	return errors.Join(be.errs...)
}

func (be *BatchError) NotNil() bool {
	be.lock.RLock()
	defer be.lock.RUnlock()
	return len(be.errs) > 0
}
