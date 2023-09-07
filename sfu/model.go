package sfu

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

// check verifies whether all provided parameters are valid. It will
// return a slice with all the errors found or nil if no errors exist.
func (p CreateSessionParams) check() []string {
	var errors []string
	if err := isNotBlank("name", p.Name); err != nil {
		errors = append(errors, err.Error())
	}
	return errors
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

// check verifies whether all provided parameters are valid. It will
// return a slice with all the errors found or nil if no errors exist.
func (p GetSessionParams) check() []string {
	var errors []string
	if err := isId("id", p.Id); err != nil {
		errors = append(errors, err.Error())
	}
	return errors
}

// GetSessionResult holds the result of GetSession
// operations.
type GetSessionResult struct {
	Session *Session `json:"session,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// DeleteSessionParams holds the parameters required to
// delete an existing live view session.
type DeleteSessionParams struct {
	Id string `json:"id"`
}

// check verifies whether all provided parameters are valid. It will
// return a slice with all the errors found or nil if no errors exist.
func (p DeleteSessionParams) check() []string {
	var errors []string
	if err := isId("id", p.Id); err != nil {
		errors = append(errors, err.Error())
	}
	return errors
}

// DeleteSessionResult holds the result of DeleteSession
// operations.
type DeleteSessionResult struct {
	Session *Session `json:"session,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// AddParticipantParams encapsulates the parameters
// used to add a new participant to an existing live
// view session.
type AddParticipantParams struct {
	SessionId string `json:"sessionId"`
	Name      string `json:"name"`
}

// check verifies whether all provided parameters are valid. It will
// return a slice with all the errors found or nil if no errors exist.
func (p AddParticipantParams) check() []string {
	var errors []string
	if err := isId("sessionId", p.SessionId); err != nil {
		errors = append(errors, err.Error())
	}
	if err := isNotBlank("name", p.Name); err != nil {
		errors = append(errors, err.Error())
	}
	return errors
}

// AddParticipantResult holds the result of AddParticipant API
// calls.
type AddParticipantResult struct {
	Participant *Participant `json:"participant,omitempty"`
	Errors      []string     `json:"errors,omitempty"`
}

// GetParticipantParams hold the required parameters to
// locate and retrieve a participant of a live view
// session.
type GetParticipantParams struct {
	SessionId     string `json:"sessionId"`
	ParticipantId string `json:"participantId"`
}

// check verifies whether all provided parameters are valid. It will
// return a slice with all the errors found or nil if no errors exist.
func (p GetParticipantParams) check() []string {
	var errors []string
	if err := isId("sessionId", p.SessionId); err != nil {
		errors = append(errors, err.Error())
	}
	if err := isId("participantId", p.ParticipantId); err != nil {
		errors = append(errors, err.Error())
	}
	return errors
}

// GetParticipantResult holds the result of GetParticipant
// API calls.
type GetParticipantResult struct {
	Participant *Participant `json:"participant,omitempty"`
	Errors      []string     `json:"errors,omitempty"`
}

// UpdateParticipantParams holds the parameters to
// locate and update an existing live view session
// participant.
type UpdateParticipantParams struct {
	SessionId     string `json:"sessionId"`
	ParticipantId string `json:"participantId"`
	Name          string `json:"name"`
}

// check verifies whether all provided parameters are valid. It will
// return a slice with all the errors found or nil if no errors exist.
func (p UpdateParticipantParams) check() []string {
	var errors []string
	if err := isId("sessionId", p.SessionId); err != nil {
		errors = append(errors, err.Error())
	}
	if err := isId("participantId", p.ParticipantId); err != nil {
		errors = append(errors, err.Error())
	}
	if err := isNotBlank("name", p.Name); err != nil {
		errors = append(errors, err.Error())
	}
	return errors
}

// UpdateParticipantResult returns the result of UpdateParticipant
// API calls.
type UpdateParticipantResult struct {
	Participant *Participant `json:"participant,omitempty"`
	Errors      []string     `json:"errors,omitempty"`
}

// DeleteParticipantParams holds all required parameters
// to locate and remove an existing participant from a
// live view session.
type DeleteParticipantParams struct {
	SessionId     string `json:"sessionId"`
	ParticipantId string `json:"participantId"`
}

// check verifies whether all provided parameters are valid. It will
// return a slice with all the errors found or nil if no errors exist.
func (p DeleteParticipantParams) check() []string {
	var errors []string
	if err := isId("sessionId", p.SessionId); err != nil {
		errors = append(errors, err.Error())
	}
	if err := isId("participantId", p.ParticipantId); err != nil {
		errors = append(errors, err.Error())
	}
	return errors
}

// DeleteParticipantResult returns the result of DeleteParticipant
// API calls.
type DeleteParticipantResult struct {
	Participant *Participant `json:"participant,omitempty"`
	Errors      []string     `json:"errors,omitempty"`
}

// GetParticipantsParams holds the required parameters to
// retrieve all participants of an existing live view session.
type GetParticipantsParams struct {
	SessionId string `json:"sessionId"`
}

func (p GetParticipantsParams) check() []string {
	var errors []string
	if err := isId("sessionId", p.SessionId); err != nil {
		errors = append(errors, err.Error())
	}
	return errors
}

// GetParticipantsResult returns the result of GetParticipants
// API calls
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
