package cache

const (
	StateNone = iota
	WaitingChannelName
	WaitingKeyWord
)

type UserCache struct {
	State uint32
	/* add mutex */
}

func CreateUser() UserCache {
	return UserCache{State: StateNone}
}
