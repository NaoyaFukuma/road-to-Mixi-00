// repository/friend_repository.go

package repository

import (
	"database/sql"
	"minimal_sns_app/domain/models"
)

// FriendRepository defines the interface for friend data access.
type FriendRepository interface {
	GetFriends(userID int64) ([]models.Friend, error)
}

type friendRepository struct {
	db *sql.DB
}

// NewFriendRepository creates a new instance of a FriendRepository.
func NewFriendRepository(db *sql.DB) FriendRepository {
	return &friendRepository{db: db}
}

// GetFriends retrieves a list of friends for a given user ID.
func (r *friendRepository) GetFriends(userID int64) ([]models.Friend, error) {
	var friends []models.Friend
	query := `SELECT f.id, f.name FROM users AS f
			JOIN friend_link AS fl ON f.user_id = fl.user2_id
			WHERE fl.user1_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var friend models.Friend
		if err := rows.Scan(&friend.ID, &friend.Name); err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}
