//go:build integration

package integrationtests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"minimal_sns_app/configs"
	"minimal_sns_app/domain/models"
	"minimal_sns_app/handlers"
	"minimal_sns_app/repository"
	"minimal_sns_app/testhelpers"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func TestGetFriendListIntegration(t *testing.T) {
	// テスト用の設定をロード
	conf := configs.Get()

	// テスト用のデータベース接続をセットアップ
	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	cleanupFunc, err := setupTestData(db)
	if err != nil {
		t.Fatalf("failed to setup test data: %v", err)
	}
	defer cleanupFunc()

	// Echo インスタンスの作成
	e := echo.New()

	// レポジトリの作成
	friendRepo := repository.NewFriendRepository(db)

	// ハンドラの作成
	friendHandler := handlers.NewFriendHandler(friendRepo)

	// ハンドラにルートを登録
	friendHandler.RegisterRoutes(e)

	// テスト用のサーバーを設定
	ts := httptest.NewServer(e)
	defer ts.Close()

	// テスト用のクライアントを作成
	client := &http.Client{}

	// 正しい user_id で GET リクエストを実行
	resp, err := client.Get(fmt.Sprintf("%s/get_friend_list?id=1", ts.URL))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	defer resp.Body.Close()

	// ステータスコードの検証
	testhelpers.AssertEqual(t, http.StatusOK, resp.StatusCode)

	// レスポンスボディの内容を読み取り
	body, err := io.ReadAll(resp.Body)
	testhelpers.AssertNoError(t, err)

	// レスポンスボディを期待する構造体にデコード
	var gotFriends []models.Friend
	err = json.Unmarshal(body, &gotFriends)
	testhelpers.AssertNoError(t, err)

	// 期待するフレンドリスト
	expectedFriends := []models.Friend{
		{
			ID:   2,
			Name: "Bob",
		},
	}

	// 取得したフレンドリストと期待するリストを比較
	testhelpers.AssertDeepEqual(t, expectedFriends, gotFriends)
}

// setupTestData inserts necessary test data into the database and returns
// a function to cleanup that data.
// It utilizes transactions to rollback changes after test completion.
func setupTestData(db *sql.DB) (cleanupFunc func(), err error) {
	// トランザクションを開始
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	// Alice を挿入
	result, err := tx.Exec("INSERT INTO users (name) VALUES ('Alice')")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	aliceID, err := result.LastInsertId() // Alice の ID を取得
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Bob を挿入
	result, err = tx.Exec("INSERT INTO users (name) VALUES ('Bob')")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	bobID, err := result.LastInsertId() // Bob の ID を取得
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Alice から Bob への友達リンクを挿入
	_, err = tx.Exec("INSERT INTO friend_link (user1_id, user2_id) VALUES (?, ?)", aliceID, bobID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	// テスト後にデータを削除するための関数を返す
	cleanupFunc = func() {
		// トランザクションを開始
		tx, err := db.Begin()
		if err != nil {
			log.Printf("failed to begin transaction: %v", err)
			return
		}

		// テストデータを削除
		if _, err := tx.Exec("DELETE FROM friend_link"); err != nil {
			log.Printf("failed to delete from friend_link: %v", err)
			tx.Rollback()
			return
		}

		if _, err := tx.Exec("DELETE FROM users"); err != nil {
			log.Printf("failed to delete from users: %v", err)
			tx.Rollback()
			return
		}

		// トランザクションをコミット
		if err := tx.Commit(); err != nil {
			log.Printf("failed to commit transaction: %v", err)
			tx.Rollback()
			return
		}
	}

	return cleanupFunc, nil
}
