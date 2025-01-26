package state

import (
	"sync"
	"todoist-tg/internal/abstractions"
)

type UserState struct {
	activeHandlers   map[int64]abstractions.CommandHandler
	activeHandlersMx sync.RWMutex
	state            map[int64]string
	stateMx          sync.RWMutex
}

func NewUserState() *UserState {
	return &UserState{
		activeHandlers: make(map[int64]abstractions.CommandHandler),
		state:          make(map[int64]string),
	}
}

func (u *UserState) SetActiveHandler(chatId int64, handler abstractions.CommandHandler) {
	u.activeHandlersMx.Lock()
	defer u.activeHandlersMx.Unlock()

	u.activeHandlers[chatId] = handler
}

func (u *UserState) GetActiveHandler(chatId int64) (abstractions.CommandHandler, bool) {
	u.activeHandlersMx.RLock()
	defer u.activeHandlersMx.RUnlock()

	handler, exists := u.activeHandlers[chatId]
	return handler, exists
}

func (u *UserState) DeleteActiveHandler(chatId int64) {
	u.activeHandlersMx.Lock()
	defer u.activeHandlersMx.Unlock()

	delete(u.activeHandlers, chatId)
}

func (u *UserState) SetUserState(chatId int64, state string) {
	u.stateMx.Lock()
	defer u.stateMx.Unlock()

	u.state[chatId] = state
}

func (u *UserState) GetUserState(chatId int64) (string, bool) {
	u.stateMx.RLock()
	defer u.stateMx.RUnlock()

	state, exists := u.state[chatId]
	return state, exists
}

func (u *UserState) DeleteUserState(chatId int64) {
	u.stateMx.Lock()
	defer u.stateMx.Unlock()

	delete(u.state, chatId)
}
