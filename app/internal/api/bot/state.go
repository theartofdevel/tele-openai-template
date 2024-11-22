package bot

import (
	"sync"
)

type Doing byte

const (
	ChatDoing Doing = iota
	ImageDoing
)

type UserState struct {
	Id    int64
	State Doing
}

type UserStates struct {
	mu   sync.RWMutex
	data []*UserState
}

func NewUserStates() UserStates {
	return UserStates{
		mu:   sync.RWMutex{},
		data: []*UserState{},
	}
}

func (uss *UserStates) Get(id int64) (*UserState, bool) {
	uss.mu.RLock()
	defer uss.mu.RUnlock()

	for _, us := range uss.data {
		if us.Id == id {
			return us, true
		}
	}

	return nil, false
}

// Set add state or update existing.
func (uss *UserStates) Set(id int64, state Doing) {
	st, exist := uss.Get(id)
	if exist {
		uss.mu.Lock()
		st.State = state
		uss.mu.Unlock()

		return
	}

	uss.mu.Lock()
	uss.data = append(uss.data, &UserState{
		Id:    id,
		State: state,
	})
	uss.mu.Unlock()
}
