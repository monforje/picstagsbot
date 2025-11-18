package search

import (
	"picstagsbot/internal/service"
	"picstagsbot/pkg/constants"
	"sync"
	"time"
)

const (
	albumSize = 10
)

type SearchSession struct {
	LastActivity time.Time
}

type SearchHandler struct {
	searchService *service.SearchService
	activeSearch  map[int64]*SearchSession
	mu            sync.RWMutex
	stopCleanup   chan struct{}
}

func NewSearchHandler(searchService *service.SearchService) *SearchHandler {
	sh := &SearchHandler{
		searchService: searchService,
		activeSearch:  make(map[int64]*SearchSession),
		stopCleanup:   make(chan struct{}),
	}

	go sh.cleanupSessions()

	return sh
}

func (h *SearchHandler) clearSession(userID int64) {
	h.mu.Lock()
	delete(h.activeSearch, userID)
	h.mu.Unlock()
}

func (h *SearchHandler) getSession(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.activeSearch[userID]
	return ok
}

func (h *SearchHandler) setSession(userID int64) {
	h.mu.Lock()
	h.activeSearch[userID] = &SearchSession{
		LastActivity: time.Now(),
	}
	h.mu.Unlock()
}

func (h *SearchHandler) cleanupSessions() {
	ticker := time.NewTicker(constants.SessionCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.mu.Lock()
			now := time.Now()
			for userID, session := range h.activeSearch {
				if now.Sub(session.LastActivity) > constants.SessionTimeout {
					delete(h.activeSearch, userID)
				}
			}
			h.mu.Unlock()
		case <-h.stopCleanup:
			return
		}
	}
}

func (h *SearchHandler) Stop() {
	close(h.stopCleanup)
}
