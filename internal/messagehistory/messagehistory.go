package messagehistory

import (
	"fmt"
	"sync"
)

type Service interface {
	Add(message string)
	Clear()
	GetCurrent() string
	NavigateDown() (string, bool)
	NavigateUp(currentMessage string) (string, bool)
	Reset()
	SetCurrent(message string)
	Size() int
}

type messageHistoryService struct {
	currentMessage string
	history        []string
	historyIndex   int
	mu             sync.RWMutex
}

var globalMessageHistoryService *messageHistoryService

func InitService() error {
	if globalMessageHistoryService != nil {
		return fmt.Errorf("message history service already initialized")
	}

	globalMessageHistoryService = &messageHistoryService{
		currentMessage: "",
		history:        make([]string, 0),
		historyIndex:   0,
	}

	return nil
}

func GetService() *messageHistoryService {
	if globalMessageHistoryService == nil {
		panic("message history service not initialized. Call messagehistory.InitService() first.")
	}
	return globalMessageHistoryService
}

func (mh *messageHistoryService) Add(message string) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	if message == "" {
		return
	}

	// Don't add if it's a duplicate of the last entry
	if len(mh.history) > 0 && mh.history[len(mh.history)-1] == message {
		return
	}

	mh.history = append(mh.history, message)
	mh.historyIndex = len(mh.history)
	mh.currentMessage = ""
}

func (mh *messageHistoryService) Clear() {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	mh.history = make([]string, 0)
	mh.historyIndex = 0
	mh.currentMessage = ""
}

func (mh *messageHistoryService) GetCurrent() string {
	mh.mu.RLock()
	defer mh.mu.RUnlock()

	return mh.currentMessage
}

func (mh *messageHistoryService) NavigateUp(currentMessage string) (string, bool) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	if len(mh.history) == 0 {
		return "", false
	}

	if mh.historyIndex == len(mh.history) {
		mh.currentMessage = currentMessage
	}

	if mh.historyIndex > 0 {
		mh.historyIndex--
		return mh.history[mh.historyIndex], true
	}

	return "", false
}

func (mh *messageHistoryService) NavigateDown() (string, bool) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	if mh.historyIndex < len(mh.history)-1 {
		// Go to next message in history
		mh.historyIndex++
		return mh.history[mh.historyIndex], true
	} else if mh.historyIndex == len(mh.history)-1 {
		// Return to the current message being composed
		mh.historyIndex = len(mh.history)
		return mh.currentMessage, true
	}

	return "", false
}

func (mh *messageHistoryService) Reset() {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	mh.historyIndex = len(mh.history)
	mh.currentMessage = ""
}

func (mh *messageHistoryService) SetCurrent(message string) {
	mh.mu.Lock()
	defer mh.mu.Unlock()

	mh.currentMessage = message
}

func (mh *messageHistoryService) Size() int {
	mh.mu.RLock()
	defer mh.mu.RUnlock()

	return len(mh.history)
}
