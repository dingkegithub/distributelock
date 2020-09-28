package lock

//
//
type Locker interface {
	Lock(...Option) (bool, error)

	UnLock() (bool, error)
}
