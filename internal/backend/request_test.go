package backend

import (
	"errors"
	"testing"

	"github.com/snippetaccumulator/snac/internal/backend/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetTeamByID(teamID string) (model.Team, error) {
	args := m.Called(teamID)
	return args.Get(0).(model.Team), args.Error(1)
}

func (m *MockDatabase) GetByID(id model.ID) (model.Snippet, error) {
	args := m.Called(id)
	return args.Get(0).(model.Snippet), args.Error(1)
}

func (m *MockDatabase) GetByTeamID(teamID string) ([]model.PartialSnippet, error) {
	args := m.Called(teamID)
	return args.Get(0).([]model.PartialSnippet), args.Error(1)
}

func (m *MockDatabase) InsertSnippet(snippet model.Snippet) (model.Snippet, error) {
	args := m.Called(snippet)
	return args.Get(0).(model.Snippet), args.Error(1)
}

func (m *MockDatabase) UpdateSnippet(snippet model.Snippet) error {
	args := m.Called(snippet)
	return args.Error(0)
}

func (m *MockDatabase) DeleteSnippet(id model.ID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDatabase) InsertTeam(name, displayName, passwordHash, adminHash string) error {
	args := m.Called(name, displayName, passwordHash, adminHash)
	return args.Error(0)
}

func (m *MockDatabase) UpdateTeam(team model.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

func (m *MockDatabase) DeleteTeam(teamID string) error {
	args := m.Called(teamID)
	return args.Error(0)
}

func (m *MockDatabase) Close() {
	m.Called()
}

func TestRequestExecute_Get(t *testing.T) {
	db := new(MockDatabase)
	snippet := model.Snippet{ID: "1", Content: "Sample"}
	snippetID := model.ID("1")

	// Setup mock for successful snippet retrieval
	db.On("GetByID", snippetID).Return(snippet, nil)

	// Test case for successful Get
	req := NewRequestBuilder().Get(snippetID).Build()
	result, retType, err := req.Execute(db)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, snippet, result)
	assert.Equal(t, ReturnSingleSnippet, retType)

	db.AssertExpectations(t)
}

func TestRequestExecute_GetAllPartials(t *testing.T) {
	db := new(MockDatabase)
	teamID := "team1"
	partials := []model.PartialSnippet{{ID: "1", Title: "Sample Partial"}}

	// Setup the mock response
	db.On("GetByTeamID", teamID).Return(partials, nil) // Successful case

	// Create the request
	req := NewRequestBuilder().ForTeamByID(teamID, "", false).GetAllPartials().Build()

	// Execute the request
	result, retType, err := req.Execute(db)
	assert.Nil(t, err)
	assert.Equal(t, ReturnPartials, retType)
	assert.Equal(t, partials, result)

	db.AssertExpectations(t)
}

func TestRequestExecute_Insert(t *testing.T) {
	db := new(MockDatabase)
	snippet := model.Snippet{ID: "1", Content: "Sample"}
	db.On("InsertSnippet", snippet).Return(snippet, nil) // Mock successful insert

	req := NewRequestBuilder().Insert(snippet).Build()

	// Test successful Insert
	result, retType, err := req.Execute(db)
	assert.Nil(t, err)
	assert.Equal(t, snippet, result)
	assert.Equal(t, ReturnSingleSnippet, retType)

	db.AssertExpectations(t)
}

func TestRequestExecute_Update(t *testing.T) {
	db := new(MockDatabase)
	snippet := model.Snippet{ID: "1", Content: "Updated Sample"}
	db.On("UpdateSnippet", snippet).Return(nil) // Mock successful update

	req := NewRequestBuilder().Update(snippet).Build()

	// Test successful Update
	_, retType, err := req.Execute(db)
	assert.Nil(t, err)
	assert.Equal(t, ReturnBoolean, retType)

	db.AssertExpectations(t)
}

func TestRequestExecute_Delete(t *testing.T) {
	db := new(MockDatabase)
	snippetID := model.ID("1")
	db.On("DeleteSnippet", snippetID).Return(nil) // Mock successful delete

	req := NewRequestBuilder().Delete(snippetID).Build()

	// Test successful Delete
	_, retType, err := req.Execute(db)
	assert.Nil(t, err)
	assert.Equal(t, ReturnBoolean, retType)

	db.AssertExpectations(t)
}

func TestRequestExecute_InsertTeam(t *testing.T) {
	db := new(MockDatabase)
	team := model.Team{Name: "newTeam", DisplayName: "New Team", PasswordHash: "passhash", AdminHash: "adminhash"}
	db.On("InsertTeam", team.Name, team.DisplayName, team.PasswordHash, team.AdminHash).Return(nil) // Mock successful insert

	req := NewRequestBuilder().NewTeam(team).Build()

	// Test successful InsertTeam
	_, retType, err := req.Execute(db)
	assert.Nil(t, err)
	assert.Equal(t, ReturnBoolean, retType)

	db.AssertExpectations(t)
}

func TestRequestExecute_UpdateTeam(t *testing.T) {
	db := new(MockDatabase)
	team := model.Team{Name: "existingTeam", DisplayName: "Updated Team", PasswordHash: "passhash", AdminHash: "adminhash"}
	db.On("UpdateTeam", team).Return(nil) // Mock successful update

	req := NewRequestBuilder().UpdateTeam(team).Build()

	// Test successful UpdateTeam
	_, retType, err := req.Execute(db)
	assert.Nil(t, err)
	assert.Equal(t, ReturnBoolean, retType)

	db.AssertExpectations(t)
}

func TestRequestExecute_DeleteTeam(t *testing.T) {
	db := new(MockDatabase)
	teamID := "team1"
	db.On("DeleteTeam", teamID).Return(nil) // Mock successful delete

	req := NewRequestBuilder().DeleteTeam(teamID).Build()

	// Test successful DeleteTeam
	_, retType, err := req.Execute(db)
	assert.Nil(t, err)
	assert.Equal(t, ReturnNone, retType)

	db.AssertExpectations(t)
}

func TestRequestExecute_Get_Negative(t *testing.T) {
	db := new(MockDatabase)
	invalidID := model.ID("2")

	// Test case for failed Get due to non-existent ID
	db.On("GetByID", invalidID).Return(model.Snippet{}, errors.New("snippet not found"))
	req := NewRequestBuilder().Get(invalidID).Build()
	_, _, err := req.Execute(db)
	assert.NotNil(t, err)

	// Test case for failed Get due to incorrect data type (e.g., passing a non-ID type)
	req.Data = "not an ID"
	_, _, err = req.Execute(db)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Request.Data for Get operation needs to be a string")

	db.AssertExpectations(t)
}

func TestRequestExecute_Insert_Negative(t *testing.T) {
	db := new(MockDatabase)
	snippet := model.Snippet{ID: "1", Content: "Sample"}
	db.On("InsertSnippet", snippet).Return(model.Snippet{}, errors.New("insert error"))
	req := NewRequestBuilder().Insert(snippet).Build()
	_, _, err := req.Execute(db)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestRequestExecute_Update_Negative(t *testing.T) {
	db := new(MockDatabase)
	snippet := model.Snippet{ID: "1", Content: "Updated Sample"}
	db.On("UpdateSnippet", snippet).Return(errors.New("update error"))
	req := NewRequestBuilder().Update(snippet).Build()
	_, _, err := req.Execute(db)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestRequestExecute_Delete_Negative(t *testing.T) {
	db := new(MockDatabase)
	snippetID := model.ID("1")
	db.On("DeleteSnippet", snippetID).Return(errors.New("delete error"))
	req := NewRequestBuilder().Delete(snippetID).Build()
	_, _, err := req.Execute(db)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestRequestExecute_InsertTeam_Negative(t *testing.T) {
	db := new(MockDatabase)
	team := model.Team{Name: "newTeam", DisplayName: "New Team", PasswordHash: "passhash", AdminHash: "adminhash"}
	db.On("InsertTeam", team.Name, team.DisplayName, team.PasswordHash, team.AdminHash).Return(errors.New("insert team error"))
	req := NewRequestBuilder().NewTeam(team).Build()
	_, _, err := req.Execute(db)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestRequestExecute_UpdateTeam_Negative(t *testing.T) {
	db := new(MockDatabase)
	team := model.Team{Name: "existingTeam", DisplayName: "Updated Team", PasswordHash: "passhash", AdminHash: "adminhash"}
	db.On("UpdateTeam", team).Return(errors.New("update team error"))
	req := NewRequestBuilder().UpdateTeam(team).Build()
	_, _, err := req.Execute(db)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestRequestExecute_DeleteTeam_Negative(t *testing.T) {
	db := new(MockDatabase)
	teamID := "team1"
	db.On("DeleteTeam", teamID).Return(errors.New("delete team error"))
	req := NewRequestBuilder().DeleteTeam(teamID).Build()
	_, _, err := req.Execute(db)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}
