package sfu

import (
	"fmt"
	"sync"
)

type webRtcSession struct {
	Session
	participants map[string]*webRtcParticipant
}

type webRtcParticipant struct {
	Participant
}

// WebRtcSessionHandler handles live view streaming
// sessions between multiple live view session
// participants.
type WebRtcSessionHandler struct {
	sessions map[string]*webRtcSession
	locker   sync.Mutex
}

// NewWebRtcSessionHandler creates and returns a properly
// initialized WebRtcSessionHandler instance.
func NewWebRtcSessionHandler() *WebRtcSessionHandler {
	h := &WebRtcSessionHandler{
		sessions: make(map[string]*webRtcSession),
	}
	return h
}

/**
========================================
     SessionHandler interface
========================================
*/

func (h *WebRtcSessionHandler) CreateSession(params CreateSessionParams) (CreateSessionResult, error) {
	if errors := params.check(); errors != nil {
		return CreateSessionResult{Errors: errors}, nil
	}
	s := newSession(params)
	h.locker.Lock()
	h.sessions[s.Id] = s
	h.locker.Unlock()
	return CreateSessionResult{Session: &s.Session}, nil
}

func newSession(params CreateSessionParams) *webRtcSession {
	return &webRtcSession{
		Session: Session{
			Id:               generateSessionId(),
			Name:             params.Name,
			CreationDateTime: generateCreationDateTime(),
		},
		participants: make(map[string]*webRtcParticipant),
	}
}

func (h *WebRtcSessionHandler) GetSession(params GetSessionParams) (GetSessionResult, error) {
	if errors := params.check(); errors != nil {
		return GetSessionResult{Errors: errors}, nil
	}
	var session *Session
	h.doActionOnSession(params.Id, func(s *webRtcSession) {
		session = &s.Session
	})
	return GetSessionResult{Session: session}, nil
}

func (h *WebRtcSessionHandler) DeleteSession(params DeleteSessionParams) (DeleteSessionResult, error) {
	if errors := params.check(); errors != nil {
		return DeleteSessionResult{Errors: errors}, nil
	}
	var session *Session
	h.doActionOnSession(params.Id, func(s *webRtcSession) {
		session = &s.Session
		delete(h.sessions, params.Id)
	})
	return DeleteSessionResult{Session: session}, nil
}

// doActionOnSession locates and executes a given action safely. It returns true
// if the action was executed, false if no such session exists.
func (h *WebRtcSessionHandler) doActionOnSession(sessionId string, action func(s *webRtcSession)) bool {
	h.locker.Lock()
	defer h.locker.Unlock()
	s := h.sessions[sessionId]
	if s == nil {
		return false
	}
	action(s)
	return true
}

func (h *WebRtcSessionHandler) AddParticipant(params AddParticipantParams) (AddParticipantResult, error) {
	if errors := params.check(); errors != nil {
		return AddParticipantResult{Errors: errors}, nil
	}
	participant := newParticipant(params)
	action := func(s *webRtcSession) {
		s.participants[participant.Id] = participant
	}
	if ok := h.doActionOnSession(params.SessionId, action); !ok {
		errorMsg := fmt.Sprintf("session %s does not exist", params.SessionId)
		return AddParticipantResult{Errors: []string{errorMsg}}, nil
	}
	return AddParticipantResult{Participant: &participant.Participant}, nil
}

func newParticipant(p AddParticipantParams) *webRtcParticipant {
	return &webRtcParticipant{
		Participant: Participant{
			Id:               generateParticipantId(),
			SessionId:        p.SessionId,
			CreationDateTime: generateCreationDateTime(),
			Name:             p.Name,
		},
	}
}

func (h *WebRtcSessionHandler) GetParticipant(params GetParticipantParams) (GetParticipantResult, error) {
	if errors := params.check(); errors != nil {
		return GetParticipantResult{Errors: errors}, nil
	}
	var participant *Participant
	action := func(s *webRtcSession) {
		p := s.participants[params.ParticipantId]
		if p != nil {
			participant = &p.Participant
		}
	}
	if ok := h.doActionOnSession(params.SessionId, action); !ok {
		errorMsg := fmt.Sprintf("session %s does not exist", params.SessionId)
		return GetParticipantResult{Errors: []string{errorMsg}}, nil
	}
	return GetParticipantResult{Participant: participant}, nil
}

func (h *WebRtcSessionHandler) UpdateParticipant(params UpdateParticipantParams) (UpdateParticipantResult, error) {
	if errors := params.check(); errors != nil {
		return UpdateParticipantResult{Errors: errors}, nil
	}
	var participant *Participant
	action := func(s *webRtcSession) {
		p := s.participants[params.ParticipantId]
		if p == nil {
			return
		}
		p.Name = params.Name
		participant = &p.Participant
	}
	if ok := h.doActionOnSession(params.SessionId, action); !ok {
		errorMsg := fmt.Sprintf("session %s does not exist", params.SessionId)
		return UpdateParticipantResult{Errors: []string{errorMsg}}, nil
	}
	return UpdateParticipantResult{Participant: participant}, nil
}

func (h *WebRtcSessionHandler) DeleteParticipant(params DeleteParticipantParams) (DeleteParticipantResult, error) {
	if errors := params.check(); errors != nil {
		return DeleteParticipantResult{Errors: errors}, nil
	}
	var participant *Participant
	action := func(s *webRtcSession) {
		p := s.participants[params.ParticipantId]
		if p == nil {
			return
		}
		delete(s.participants, params.ParticipantId)
		participant = &p.Participant
	}
	if ok := h.doActionOnSession(params.SessionId, action); !ok {
		errorMsg := fmt.Sprintf("session %s does not exist", params.SessionId)
		return DeleteParticipantResult{Errors: []string{errorMsg}}, nil
	}
	return DeleteParticipantResult{Participant: participant}, nil
}

func (h *WebRtcSessionHandler) GetParticipants(params GetParticipantsParams) (GetParticipantsResult, error) {
	if errors := params.check(); errors != nil {
		return GetParticipantsResult{Errors: errors}, nil
	}
	var participants []*Participant
	action := func(s *webRtcSession) {
		participants = make([]*Participant, len(s.participants))
		i := 0
		for _, v := range s.participants {
			participants[i] = &v.Participant
			i++
		}
	}
	if ok := h.doActionOnSession(params.SessionId, action); !ok {
		errorMsg := fmt.Sprintf("session %s does not exist", params.SessionId)
		return GetParticipantsResult{Errors: []string{errorMsg}}, nil
	}
	return GetParticipantsResult{Participants: participants}, nil
}
