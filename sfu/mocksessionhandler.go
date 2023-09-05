package sfu

import (
	"fmt"
	"sync"
	"time"
)

const timeFormat = "2006-01-02T15:04:05 -070000"

type RegisteredSession struct {
	Session
	participants map[string]*Participant
}

type MockSessionHandler struct {
	sessions map[string]*RegisteredSession
	mutex    sync.Mutex
}

func (h *MockSessionHandler) CreateSession(p CreateSessionParams) (CreateSessionResult, error) {
	if err := checkNotBlank(p.Name); err != nil {
		return CreateSessionResult{
			Errors: []string{err.Error()},
		}, nil
	}
	if h.sessions == nil {
		h.sessions = make(map[string]*RegisteredSession)
	}
	session := Session{
		Name:             p.Name,
		Id:               generateSessionId(),
		CreationDateTime: time.Now().Format(timeFormat),
	}
	registeredSession := RegisteredSession{
		Session: session,
	}
	h.mutex.Lock()
	h.sessions[session.Id] = &registeredSession
	h.mutex.Unlock()
	return CreateSessionResult{Session: &registeredSession.Session}, nil
}

func (h *MockSessionHandler) GetSession(p GetSessionParams) (GetSessionResult, error) {
	if err := checkSessionId(p.Id); err != nil {
		return GetSessionResult{Errors: []string{err.Error()}}, nil
	}
	h.mutex.Lock()
	registeredSession := h.sessions[p.Id]
	h.mutex.Unlock()
	return GetSessionResult{Session: &registeredSession.Session}, nil
}

func (h *MockSessionHandler) DeleteSession(p DeleteSessionParams) (DeleteSessionResult, error) {
	if err := checkSessionId(p.Id); err != nil {
		return DeleteSessionResult{Errors: []string{err.Error()}}, nil
	}
	h.mutex.Lock()
	registeredSession := h.sessions[p.Id]
	delete(h.sessions, p.Id)
	h.mutex.Unlock()
	return DeleteSessionResult{Session: &registeredSession.Session}, nil
}

func (h *MockSessionHandler) AddParticipant(p AddParticipantParams) (AddParticipantResult, error) {
	if err := checkSessionId(p.SessionId); err != nil {
		return AddParticipantResult{Errors: []string{err.Error()}}, nil
	}
	if err := checkNotBlank(p.Name); err != nil {
		return AddParticipantResult{Errors: []string{err.Error()}}, nil
	}
	h.mutex.Lock()
	defer h.mutex.Unlock()
	registeredSession := h.sessions[p.SessionId]
	if registeredSession == nil {
		e := fmt.Errorf("session %s does not exist", p.SessionId)
		return AddParticipantResult{Errors: []string{e.Error()}}, nil
	}
	if registeredSession.participants == nil {
		registeredSession.participants = make(map[string]*Participant)
	}
	participant := Participant{
		SessionId:        registeredSession.Id,
		Id:               generateSessionId(),
		Name:             p.Name,
		CreationDateTime: time.Now().Format(timeFormat),
	}
	registeredSession.participants[participant.Id] = &participant
	return AddParticipantResult{Participant: &participant}, nil
}

func (h *MockSessionHandler) GetParticipant(p GetParticipantParams) (GetParticipantResult, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	registeredSession := h.sessions[p.SessionId]
	if registeredSession == nil {
		e := fmt.Errorf("session %s does not exist", p.SessionId)
		return GetParticipantResult{Errors: []string{e.Error()}}, nil
	}
	participant := registeredSession.participants[p.ParticipantId]
	return GetParticipantResult{Participant: participant}, nil
}

func (h *MockSessionHandler) UpdateParticipant(p UpdateParticipantParams) (UpdateParticipantResult, error) {
	if err := checkNotBlank(p.Name); err != nil {
		return UpdateParticipantResult{Errors: []string{err.Error()}}, nil
	}
	h.mutex.Lock()
	defer h.mutex.Unlock()
	registeredSession := h.sessions[p.SessionId]
	if registeredSession == nil {
		e := fmt.Errorf("session %s does not exist", p.SessionId)
		return UpdateParticipantResult{Errors: []string{e.Error()}}, nil
	}
	participant := registeredSession.participants[p.ParticipantId]
	if participant == nil {
		return UpdateParticipantResult{}, nil
	}
	participant.Name = p.Name
	return UpdateParticipantResult{Participant: participant}, nil
}

func (h *MockSessionHandler) DeleteParticipant(p DeleteParticipantParams) (DeleteParticipantResult, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	registeredSession := h.sessions[p.SessionId]
	if registeredSession == nil {
		e := fmt.Errorf("session %s does not exist", p.SessionId)
		return DeleteParticipantResult{Errors: []string{e.Error()}}, nil
	}
	participant := registeredSession.participants[p.ParticipantId]
	if participant == nil {
		return DeleteParticipantResult{}, nil
	}
	delete(registeredSession.participants, participant.Id)
	return DeleteParticipantResult{Participant: participant}, nil
}

func (h *MockSessionHandler) GetParticipants(p GetParticipantsParams) (GetParticipantsResult, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	registeredSession := h.sessions[p.SessionId]
	if registeredSession == nil {
		e := fmt.Errorf("session %s does not exist", p.SessionId)
		return GetParticipantsResult{Errors: []string{e.Error()}}, nil
	}
	participants := make([]*Participant, len(registeredSession.participants))
	i := 0
	for _, v := range registeredSession.participants {
		participants[i] = v
		i++
	}
	return GetParticipantsResult{Participants: participants}, nil
}
