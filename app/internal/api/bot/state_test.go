package bot

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserStates(t *testing.T) {
	states := NewUserStates()

	assert.NotNil(t, states)
	assert.Empty(t, states.data)
	assert.IsType(t, sync.RWMutex{}, states.mu)
}

func TestUserStates_Get(t *testing.T) {
	states := UserStates{
		mu: sync.RWMutex{},
		data: []*UserState{
			{Id: 1, State: ChatDoing},
			{Id: 2, State: ImageDoing},
		},
	}

	tests := []struct {
		name       string
		id         int64
		wantState  *UserState
		wantExists bool
	}{
		{
			name:       "existing chat state",
			id:         1,
			wantState:  &UserState{Id: 1, State: ChatDoing},
			wantExists: true,
		},
		{
			name:       "existing image state",
			id:         2,
			wantState:  &UserState{Id: 2, State: ImageDoing},
			wantExists: true,
		},
		{
			name:       "non-existing state",
			id:         3,
			wantState:  nil,
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, exists := states.Get(tt.id)

			assert.Equal(t, tt.wantExists, exists)
			if tt.wantState != nil {
				assert.Equal(t, tt.wantState.Id, state.Id)
				assert.Equal(t, tt.wantState.State, state.State)
			} else {
				assert.Nil(t, state)
			}
		})
	}
}

func TestUserStates_Set(t *testing.T) {
	t.Run("add new state", func(t *testing.T) {
		states := NewUserStates()

		states.Set(1, ChatDoing)

		state, exists := states.Get(1)
		assert.True(t, exists)
		assert.Equal(t, int64(1), state.Id)
		assert.Equal(t, ChatDoing, state.State)
	})

	t.Run("update existing state", func(t *testing.T) {
		states := UserStates{
			mu: sync.RWMutex{},
			data: []*UserState{
				{Id: 1, State: ChatDoing},
			},
		}

		states.Set(1, ImageDoing)

		state, exists := states.Get(1)
		assert.True(t, exists)
		assert.Equal(t, int64(1), state.Id)
		assert.Equal(t, ImageDoing, state.State)
	})
}

func TestUserStates_Concurrency(t *testing.T) {
	states := NewUserStates()
	var wg sync.WaitGroup

	// Test concurrent reads and writes
	for i := int64(1); i <= 100; i++ {
		wg.Add(2)

		id := i
		go func() {
			defer wg.Done()
			states.Set(id, ChatDoing)
		}()

		go func() {
			defer wg.Done()
			states.Get(id)
		}()
	}

	wg.Wait()

	// Verify all states were set correctly
	for i := int64(1); i <= 100; i++ {
		state, exists := states.Get(i)
		assert.True(t, exists)
		assert.Equal(t, i, state.Id)
		assert.Equal(t, ChatDoing, state.State)
	}
}
