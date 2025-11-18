package upload

import (
	"picstagsbot/internal/service"
	"picstagsbot/pkg/constants"
	"sync"
	"time"
)

type SessionState string

const (
	StateAwaitingPhoto       SessionState = "awaiting_photo"
	StateAwaitingDescription SessionState = "awaiting_description"
)

type UploadedPhoto struct {
	FileID   string
	FileSize int64
	Width    int
	Height   int
}

type UploadSession struct {
	State           SessionState
	Photos          []UploadedPhoto
	LastMediaGroup  string
	PendingResponse bool
	LastActivity    time.Time
}

type UploadHandler struct {
	uploadService *service.UploadService
	sessions      map[int64]*UploadSession
	mu            sync.RWMutex
	stopCleanup   chan struct{}
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	uh := &UploadHandler{
		uploadService: uploadService,
		sessions:      make(map[int64]*UploadSession),
		stopCleanup:   make(chan struct{}),
	}

	go uh.cleanupSessions()

	return uh
}

func (h *UploadHandler) clearSession(userID int64) {
	h.mu.Lock()
	delete(h.sessions, userID)
	h.mu.Unlock()
}

func (h *UploadHandler) getSession(userID int64) *UploadSession {
	h.mu.RLock()
	defer h.mu.RUnlock()

	session, ok := h.sessions[userID]
	if !ok {
		return nil
	}

	sessionCopy := &UploadSession{
		State:           session.State,
		LastMediaGroup:  session.LastMediaGroup,
		PendingResponse: session.PendingResponse,
		LastActivity:    session.LastActivity,
		Photos:          make([]UploadedPhoto, len(session.Photos)),
	}
	copy(sessionCopy.Photos, session.Photos)

	return sessionCopy
}

func (h *UploadHandler) setSession(userID int64, session *UploadSession) {
	h.mu.Lock()
	session.LastActivity = time.Now()
	h.sessions[userID] = session
	h.mu.Unlock()
}

func (h *UploadHandler) addPhotoToSession(userID int64, photo UploadedPhoto) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	session, ok := h.sessions[userID]
	if !ok || session.State != StateAwaitingPhoto {
		return false
	}

	if len(session.Photos) >= constants.MaxPhotosPerSession {
		return false
	}

	session.Photos = append(session.Photos, photo)
	session.LastActivity = time.Now()
	return true
}

func (h *UploadHandler) updateSessionState(userID int64, state SessionState) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if session, ok := h.sessions[userID]; ok {
		session.State = state
		session.LastActivity = time.Now()
	}
}

func (h *UploadHandler) cleanupSessions() {
	ticker := time.NewTicker(constants.SessionCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.mu.Lock()
			now := time.Now()
			for userID, session := range h.sessions {
				if now.Sub(session.LastActivity) > constants.SessionTimeout {
					delete(h.sessions, userID)
				}
			}
			h.mu.Unlock()
		case <-h.stopCleanup:
			return
		}
	}
}

func (h *UploadHandler) Stop() {
	close(h.stopCleanup)
}
