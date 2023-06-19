package DistributedLock

type Lock interface {
	Lock() error
	UnLock() error
}
