package handlers

import (
	"minimal_sns_app/domain/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFriendRepositoryはFriendRepositoryのモックです。
type MockFriendRepository struct {
	mock.Mock
}

// GetFriendsはモック関数です。
func (_m *MockFriendRepository) GetFriends(userID int64) ([]models.Friend, error) {
	ret := _m.Called(userID)
	return ret.Get(0).([]models.Friend), ret.Error(1)
}

func TestGetFriendList(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/?id=1", nil)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	mockRepo := new(MockFriendRepository)
	handler := NewFriendHandler(mockRepo)

	expectedFriends := []models.Friend{{ID: 1, Name: "Alice"}}
	mockRepo.On("GetFriends", int64(1)).Return(expectedFriends, nil)

	// テスト対象のメソッドを実行します
	if assert.NoError(t, handler.GetFriendList(context)) {
		assert.Equal(t, http.StatusOK, recorder.Code)
		mockRepo.AssertExpectations(t) // モックに期待した呼び出しがあったかを確認
	}
}

// Get2hopsFriendsはモック関数です。
func (_m *MockFriendRepository) Get2hopsFriends(userID int64) ([]models.Friend, error) {
	ret := _m.Called(userID)
	return ret.Get(0).([]models.Friend), ret.Error(1)
}

func TestGet2hopsFriends(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/?id=1", nil)
	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	mockRepo := new(MockFriendRepository)
	handler := NewFriendHandler(mockRepo)

	expected2hopsFriends := []models.Friend{
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
	mockRepo.On("Get2hopsFriends", int64(1)).Return(expected2hopsFriends, nil)

	// テスト対象のメソッドを実行します
	if assert.NoError(t, handler.Get2hopsFriends(context)) {
		assert.Equal(t, http.StatusOK, recorder.Code)
		// レスポンスボディの内容が期待するJSON文字列になっているかを検証します
		assert.JSONEq(t, `[
			{"id":2,"name":"Bob"},
			{"id":3,"name":"Charlie"}
		]`, recorder.Body.String())
		mockRepo.AssertExpectations(t) // モックに期待した呼び出しがあったかを確認
	}
}

// DeleteFriendはモック関数です。
func (_m *MockFriendRepository) DeleteFriend(userID int64, friendID int64) error {
	ret := _m.Called(userID, friendID)
	return ret.Error(0)
}
