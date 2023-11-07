package main

import (
	"database/sql"
	"minimal_sns_app/configs"
	"minimal_sns_app/handlers"
	"minimal_sns_app/repository"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	// アプリケーション設定のロード
	conf := configs.Get()

	// データベース接続のセットアップ
	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Echo インスタンスの作成
	e := echo.New()

	// レポジトリとハンドラの作成
	friendRepo := repository.NewFriendRepository(db)
	friendHandler := handlers.NewFriendHandler(friendRepo)

	// ハンドラにルートを登録させる
	friendHandler.RegisterRoutes(e)

	// サーバーの起動
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
