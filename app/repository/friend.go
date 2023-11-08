// repository/friend_repository.go

package repository

import (
	"database/sql"
	"minimal_sns_app/domain/models"
	"minimal_sns_app/logutils"
)

// FriendRepository defines the interface for friend data access.
type FriendRepository interface {
	GetFriends(userID int64) ([]models.Friend, error)
	Get2hopsFriends(userID int64) ([]models.Friend, error)
	DeleteFriend(userID int64, friendID int64) error
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
	query := `SELECT u.id, u.name FROM users AS u
			JOIN friend_link AS fl ON u.id = fl.user2_id
			WHERE fl.user1_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		logutils.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var friend models.Friend
		if err := rows.Scan(&friend.ID, &friend.Name); err != nil {
			logutils.Error(err.Error())
			return nil, err
		}
		friends = append(friends, friend)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		logutils.Error(err.Error())
		return nil, err
	}

	return friends, nil
}

// GetFrineds retrieves a list of two hops friends for a given user ID.
func (r *friendRepository) Get2hopsFriends(userID int64) ([]models.Friend, error) {
	var friends []models.Friend
	query := `SELECT DISTINCT u2.id, u2.name FROM users AS u1
						JOIN friend_link AS fl1 ON u1.id = fl1.user1_id
						JOIN friend_link AS fl2 ON fl1.user2_id = fl2.user1_id
						JOIN users AS u2 ON fl2.user2_id = u2.id
						LEFT JOIN block_list AS bl ON bl.user1_id = u1.id AND bl.user2_id = u2.id
						WHERE u1.id = ? AND u2.id != u1.id AND u2.id NOT IN (
							SELECT user2_id FROM friend_link WHERE user1_id = u1.id
						) AND bl.user1_id IS NULL`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		logutils.Error("Failed to get 2 hops friends")
		logutils.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var friend models.Friend
		if err := rows.Scan(&friend.ID, &friend.Name); err != nil {
			logutils.Error("Failed to scan 2 hops friends")
			logutils.Error(err.Error())
			return nil, err
		}
		// 追加する前に、ブロックされていない友達だけをリストに追加する。
		friends = append(friends, friend)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		logutils.Error("Failed to iterate over rows")
		logutils.Error(err.Error())
		return nil, err
	}
	return friends, nil
}

// DeleteFriend deletes a friend for a given user ID and friend ID.
func (r *friendRepository) DeleteFriend(userID int64, friendID int64) error {
	query := `DELETE FROM friend_link WHERE user1_id = ? AND user2_id = ?`
	_, err := r.db.Exec(query, userID, friendID)
	if err != nil {
		logutils.Error(err.Error())
		return err
	}
	return nil
}
