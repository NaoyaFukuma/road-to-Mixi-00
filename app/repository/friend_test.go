//go:build mock

package repository

import (
	"minimal_sns_app/domain/models"
	"minimal_sns_app/testhelpers"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetFriends(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewFriendRepository(db)

	tests := []struct {
		name        string
		prepareMock func()
		userID      int64
		expected    []models.Friend
		expectError bool
	}{
		{
			name: "With valid user, returns friends list",
			prepareMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Alice").
					AddRow(2, "Bob")
				mock.ExpectQuery("^SELECT f.id, f.name FROM users AS f JOIN friend_link AS fl ON f.user_id = fl.user2_id WHERE fl.user1_id = \\?$").
					WithArgs(1).
					WillReturnRows(rows)
			},
			userID: 1,
			expected: []models.Friend{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			expectError: false,
		},
		// ... 他のテストケースをここに追加 ...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMock()
			friends, err := r.GetFriends(tt.userID)
			if tt.expectError {
				testhelpers.AssertError(t, err)
			} else {
				testhelpers.AssertNoError(t, err)
				testhelpers.AssertEqual(t, len(tt.expected), len(friends))
				for i, friend := range friends {
					testhelpers.AssertEqual(t, tt.expected[i], friend)
				}
			}
		})
	}
}
