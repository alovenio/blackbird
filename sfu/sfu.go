package sfu

import (
	"alovenio.com/blackbird/logger"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

// Session holds all information related to a single
// live view session.
type Session struct {
	Name             string `json:"name"`
	Id               string `json:"id"`
	CreationDateTime string `json:"creationDateTime"`
}

// Participant holds all information related to a single
// participant of a live view session.
type Participant struct {
	Name             string `json:"name"`
	Id               string `json:"id"`
	SessionId        string `json:"sessionId"`
	CreationDateTime string `json:"creationDateTime"`
}

// CreateSessionParams holds all parameters required
// to create a new live view session.
type CreateSessionParams struct {
	Name string `json:"name"`
}

// CreateSessionResult holds the result of CreateSession
// operations.
type CreateSessionResult struct {
	// Pointer to the created session. A nil value means no session was created.
	Session *Session `json:"session,omitempty"`
	// Slices with all errors that prevented a session to be created. Can be nil.
	Errors []string `json:"errors,omitempty"`
}

// GetSessionParams holds all parameters required to
// locate and retrieve an existing live view session.
type GetSessionParams struct {
	Id string `json:"id"`
}

// GetSessionResult holds the result of GetSession
// operations.
type GetSessionResult struct {
	Session *Session `json:"session,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

type DeleteSessionParams struct {
	Id string `json:"id"`
}

type DeleteSessionResult struct {
	Session *Session `json:"session,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

type AddParticipantParams struct {
	SessionId string `json:"sessionId"`
	Name      string `json:"name"`
}

type AddParticipantResult struct {
	Participant *Participant `json:"participant,omitempty"`
	Errors      []string     `json:"errors,omitempty"`
}

type GetParticipantParams struct {
	SessionId     string `json:"sessionId"`
	ParticipantId string `json:"participantId"`
}

type GetParticipantResult struct {
	Participant *Participant `json:"participant,omitempty"`
	Errors      []string     `json:"errors,omitempty"`
}

type UpdateParticipantParams struct {
	SessionId     string `json:"sessionId"`
	ParticipantId string `json:"participantId"`
	Name          string `json:"name"`
}

type UpdateParticipantResult struct {
	Participant *Participant `json:"participant,omitempty"`
	Errors      []string     `json:"errors,omitempty"`
}

type DeleteParticipantParams struct {
	SessionId     string `json:"sessionId"`
	ParticipantId string `json:"participantId"`
}

type DeleteParticipantResult struct {
	Participant *Participant `json:"participant,omitempty"`
	Errors      []string     `json:"errors,omitempty"`
}

type GetParticipantsParams struct {
	SessionId string `json:"sessionId"`
}

type GetParticipantsResult struct {
	Participants []*Participant `json:"participants,omitempty"`
	Errors       []string       `json:"errors,omitempty"`
}

// SessionHandler defines the interface for implementors
// of live view sessions.
type SessionHandler interface {
	// CreateSession creates a new live view session. On success,
	// a pointer to the newly created session will be available
	// inside results object. If session creation fails due to an
	// expected error, the results object will have its Errors
	// property populated. Returning an error outside the results
	// object will be the case in the presence of unexpected
	// conditions, and should be interpreted as an internal server
	// error.
	CreateSession(p CreateSessionParams) (CreateSessionResult, error)
	// GetSession locates and retrieves an existing live view session.
	// The located session pointer will be available inside results object.
	// If no such session exists, the pointer will be nil. If session
	// retrieval fails due to an expected error, the results object will
	// have its Errors property populated. Returning an error outside the
	// results object will be the case when unexpected conditions are detected,
	// and should be interpreted as an internal server error.
	GetSession(p GetSessionParams) (GetSessionResult, error)
	// DeleteSession locates and deletes an existing live view session.
	// The located session will be added to the operation's result, unless
	// the session does not exist. In the presence of any expected error, the
	// results object will have its Errors property populated. Returning an
	// error outside the results object will be the case when unexpected conditions
	// are detected, and should be interpreted as an internal server error.
	DeleteSession(p DeleteSessionParams) (DeleteSessionResult, error)
	// AddParticipant adds a new participant to an existing live view session.
	// On success, a pointer to the newly added participant will be available
	// inside results object. If participant addition fails due to an expected
	// error, the results object will have its Errors property populated. Returning
	// an error outside the results object will be the case when unexpected conditions
	// are detected and should be interpreted as an internal server error.
	AddParticipant(p AddParticipantParams) (AddParticipantResult, error)
	// GetParticipant locates and retrieves an existing participant of a live view
	// session. The located participant pointer will be available inside the results
	// object. If no such participant exists, the pointer will be nil. If participant
	// retrieval fails due to an unexpected error, the results object will have its
	// Errors property populated. Returning an error outside the results object will
	// be the case when unexpected conditions are detected, and should be interpreted
	// as an internal server error.
	GetParticipant(p GetParticipantParams) (GetParticipantResult, error)
	// UpdateParticipant locates and updates an existing participant of a live view
	// session. On success, a pointer to the updated participant will be present
	// in the results object. If update fails due to expected conditions, the results
	// object will have its errors slice populated. If an unexpected error is encountered,
	// this call will return an error which should be interpreted as an internal
	// server error.
	UpdateParticipant(p UpdateParticipantParams) (UpdateParticipantResult, error)
	// DeleteParticipant deletes an existing participant from an existing live view
	// session. On success, a pointer to the deleted participant will be present in
	// the results object. If deletion fails due to expected conditions, the results
	// object will have its errors slice populated. If an unexpected error is encountered,
	// this call will return an error which should be interpreted as an internal
	// server error.
	DeleteParticipant(p DeleteParticipantParams) (DeleteParticipantResult, error)
	// GetParticipants retrieves all participants of an existing live view session.
	// On success, the participants slice in the results object will be populated. If
	// expected errors are detected, the Errors property of the results object will be
	// populated. If an unexpected error is encountered,
	// this call will return an error which should be interpreted as an internal
	// server error.
	GetParticipants(p GetParticipantsParams) (GetParticipantsResult, error)
}

// Server objects represent instances of Blackbird's SFU
// server.
type Server struct {
	handler       *SessionHandler
	startDateTime time.Time
	address       string
}

func (s *Server) Start(addr string, handler SessionHandler) error {
	if err := checkAddr(addr); err != nil {
		return err
	}
	s.startDateTime = time.Now()
	s.address = addr
	s.handler = &handler
	router := mux.NewRouter()
	router.HandleFunc("/{version}/sessions", s.onSessionsRequest)
	router.HandleFunc("/{version}/sessions/{sessionId}", s.onSessionRequest)
	router.HandleFunc("/{version}/sessions/{sessionId}/participants", s.onSessionParticipantsRequest)
	router.HandleFunc("/{version}/sessions/{sessionId}/participants/{participantId}", s.onSessionParticipantRequest)
	router.Use(contentTypeMiddleware)
	logger.LogInfoF("Starting Blackbird SFU server on %s...", addr)
	logger.LogFatalF(http.ListenAndServe(addr, router))
	return nil
}

// checkAddr checks whether the given addr parameter is a valid server
// address. In the case the given address is found invalid, an error
// will be returned.
func checkAddr(addr string) error {
	if len(strings.TrimSpace(addr)) == 0 {
		return fmt.Errorf("server address must not be null or blank")
	}
	return nil
}

// contentTypeMiddleware is called before handling any http request.
// It sets the response's Content-Type header to application/json.
func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// onSessionsRequest is called for every request to /{version}/sessions API
func (s *Server) onSessionsRequest(w http.ResponseWriter, r *http.Request) {
	if isPutOrPost(r) == false {
		logger.LogWarnF(requestAwareMsg(r, "operation not supported: %s", r.Method))
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	params := CreateSessionParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.LogWarnF(requestAwareMsg(r, "decoding error: %s", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := (*s.handler).CreateSession(params)
	if err != nil {
		logger.LogErrorF(requestAwareMsg(r, "handling error: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.Errors != nil {
		logger.LogWarnF(requestAwareMsg(r, "bad request: %s", result.Errors))
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogWarnF(requestAwareMsg(r, "failed to encode result: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// onSessionRequest is called for every request to /{version}/sessions/{sessionId}
func (s *Server) onSessionRequest(w http.ResponseWriter, r *http.Request) {
	if isGet(r) {
		s.onGetSessionRequest(w, r)
	} else if isDelete(r) {
		s.onDeleteSessionRequest(w, r)
	} else {
		logger.LogWarnF(requestAwareMsg(r, "operation not supported"))
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// onGetSessionRequest is called for every GET request to /{version}/sessions/{sessionId}
func (s *Server) onGetSessionRequest(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	sessionId := vars["sessionId"]
	params := GetSessionParams{Id: sessionId}
	result, err := (*s.handler).GetSession(params)
	if err != nil {
		logger.LogErrorF(requestAwareMsg(r, "handling error: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.Errors != nil {
		logger.LogWarnF(requestAwareMsg(r, "bad request: %s", result.Errors))
		w.WriteHeader(http.StatusBadRequest)
	} else if result.Session == nil {
		logger.LogDebugF(requestAwareMsg(r, "no such session: %s", sessionId))
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogWarnF(requestAwareMsg(r, "failed to encode result: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// onDeleteSessionRequest is called for every DELETE request to /{version}/sessions/{sessionId}
func (s *Server) onDeleteSessionRequest(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	sessionId := vars["sessionId"]
	params := DeleteSessionParams{Id: sessionId}
	result, err := (*s.handler).DeleteSession(params)
	if err != nil {
		logger.LogErrorF(requestAwareMsg(r, "handling error: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.Errors != nil {
		logger.LogWarnF(requestAwareMsg(r, "bad request: %s", result.Errors))
		w.WriteHeader(http.StatusBadRequest)
	} else if result.Session == nil {
		logger.LogDebugF(requestAwareMsg(r, "no such session: %s", sessionId))
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogWarnF(requestAwareMsg(r, "failed to encode result: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// onSessionParticipantsRequest is called for every request to /{version}/sessions/{sessionId}/participants
func (s *Server) onSessionParticipantsRequest(w http.ResponseWriter, r *http.Request) {
	if isGet(r) == true {
		s.onGetSessionParticipantsRequest(w, r)
	} else if isPutOrPost(r) == true {
		s.onPostSessionParticipantsRequest(w, r)
	} else {
		logger.LogWarnF(requestAwareMsg(r, "operation not supported"))
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// onGetSessionParticipantsRequest is called for every GET request to /{version}/sessions/{sessionId}/participants
func (s *Server) onGetSessionParticipantsRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionId := vars["sessionId"]
	result, err := (*s.handler).GetParticipants(GetParticipantsParams{SessionId: sessionId})
	if err != nil {
		logger.LogErrorF(requestAwareMsg(r, "handling error: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.Errors != nil {
		logger.LogWarnF(requestAwareMsg(r, "bad request: %s", result.Errors))
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogWarnF(requestAwareMsg(r, "failed to encode result: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// onGetSessionParticipantsRequest is called for every POST request to /{version}/sessions/{sessionId}/participants
func (s *Server) onPostSessionParticipantsRequest(w http.ResponseWriter, r *http.Request) {
	params := AddParticipantParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.LogWarnF(requestAwareMsg(r, "decoding error: %s", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	params.SessionId = mux.Vars(r)["sessionId"]
	result, err := (*s.handler).AddParticipant(params)
	if err != nil {
		logger.LogErrorF(requestAwareMsg(r, "handling error: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.Errors != nil {
		logger.LogWarnF(requestAwareMsg(r, "bad request: %s", result.Errors))
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogWarnF(requestAwareMsg(r, "failed to encode result: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// onSessionParticipantRequest is called for every request to
// /{version}/sessions/{sessionId}/participants/{participantId}
func (s *Server) onSessionParticipantRequest(w http.ResponseWriter, r *http.Request) {
	if isGet(r) == true {
		s.onGetSessionParticipantRequest(w, r)
	} else if isPutOrPost(r) == true {
		s.onUpdateSessionParticipantRequest(w, r)
	} else if isDelete(r) == true {
		s.onDeleteSessionParticipantRequest(w, r)
	} else {
		logger.LogWarnF(requestAwareMsg(r, "operation not supported"))
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// onGetSessionParticipantRequest is called for every GET request to
// /{version}/sessions/{sessionId}/participants/{participantId}
func (s *Server) onGetSessionParticipantRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionId := vars["sessionId"]
	participantId := vars["participantId"]
	result, err := (*s.handler).GetParticipant(GetParticipantParams{
		SessionId:     sessionId,
		ParticipantId: participantId,
	})
	if err != nil {
		logger.LogErrorF(requestAwareMsg(r, "handling error: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.Errors != nil {
		logger.LogWarnF(requestAwareMsg(r, "bad request: %s", result.Errors))
		w.WriteHeader(http.StatusBadRequest)
	} else if result.Participant == nil {
		logger.LogDebugF(requestAwareMsg(r, "no such participant %q in session %q", participantId, sessionId))
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogWarnF(requestAwareMsg(r, "failed to encode result: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// onUpdateSessionParticipantRequest is called for every POST/PUT request to
// /{version}/sessions/{sessionId}/participants/{participantId}
func (s *Server) onUpdateSessionParticipantRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionId := vars["sessionId"]
	participantId := vars["participantId"]
	params := UpdateParticipantParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.LogWarnF(requestAwareMsg(r, "decoding error: %s", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	params.SessionId = sessionId
	params.ParticipantId = participantId
	result, err := (*s.handler).UpdateParticipant(params)
	if err != nil {
		logger.LogErrorF(requestAwareMsg(r, "handling error: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.Errors != nil {
		logger.LogWarnF(requestAwareMsg(r, "bad request: %s", result.Errors))
		w.WriteHeader(http.StatusBadRequest)
	} else if result.Participant == nil {
		logger.LogDebugF(requestAwareMsg(r, "no such participant %q in session %q", participantId, sessionId))
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogWarnF(requestAwareMsg(r, "failed to encode result: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// onDeleteSessionParticipantRequest is called for every DELETE request to
// /{version}/sessions/{sessionId}/participants/{participantId}
func (s *Server) onDeleteSessionParticipantRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionId := vars["sessionId"]
	participantId := vars["participantId"]
	result, err := (*s.handler).DeleteParticipant(DeleteParticipantParams{
		SessionId:     sessionId,
		ParticipantId: participantId,
	})
	if err != nil {
		logger.LogErrorF(requestAwareMsg(r, "handling error: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.Errors != nil {
		logger.LogWarnF(requestAwareMsg(r, "bad request: %s", result.Errors))
		w.WriteHeader(http.StatusBadRequest)
	} else if result.Participant == nil {
		logger.LogDebugF(requestAwareMsg(r, "no such participant %q in session %q", participantId, sessionId))
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogWarnF(requestAwareMsg(r, "failed to encode result: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// isPutOrPost returns whether a given request object refers to a PUT or POST http method.
func isPutOrPost(r *http.Request) bool {
	return r.Method == "PUT" || r.Method == "POST"
}

// isGet returns whether a given request object refers to a GET http method.
func isGet(r *http.Request) bool {
	return r.Method == "GET"
}

// isDelete returns whether a given request object refers to a DELETE http method.
func isDelete(r *http.Request) bool {
	return r.Method == "DELETE"
}

// requestAwareMsg creates a message in the context of a given request
func requestAwareMsg(r *http.Request, format string, args ...any) string {
	prefix := fmt.Sprintf("%s %s:", r.RequestURI, r.Method) + format
	return fmt.Sprintf(prefix, args...)
}
