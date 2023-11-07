# プロジェクト名を定義
NAME := minimal_sns_app

# 依存関係として$(NAME)を指定
.PHONY: $(NAME)
$(NAME): all

# コンテナ起動 バックグラウンド
.PHONY: all
all:
	docker compose -p $(NAME) up -d

# コンテナ起動 フォアグラウンド(ログを見たいとき)
.PHONY: up-logs
up-logs:
	docker compose -p $(NAME) up

# コンテナを停止
.PHONY: stop
stop:
	docker compose -p $(NAME) stop

# プロジェクト名を定義
NAME := minimal_sns_app

# プログラムに関連するコンテナだけを停止し削除
.PHONY: clean-containers
clean-containers:
	docker compose -p $(NAME) down --remove-orphans

# プロジェクトに関連するコンテナとイメージを削除
.PHONY: clean-images
clean-images:
	docker compose -p $(NAME) down --rmi local

# プロジェクトに関連するコンテナとボリュームを削除
.PHONY: clean-volumes
clean-volumes:
	docker compose -p $(NAME) down -v --remove-orphans

# プロジェクトに関連するネットワークだけを削除
.PHONY: clean-networks
clean-networks:
	docker network prune -f --filter label=com.docker.compose.project=$(NAME)

# コンテナ、イメージ、ボリューム、ネットワークを削除
.PHONY: fclean
fclean: clean-containers clean-images clean-volumes clean-networks

# コンテナを再構築
.PHONY: re
re: fclean all

# モックを使ったテスト
.PHONY: test-mock
test-mock:
	docker compose -p $(NAME) exec -T app go test -v -cover -tags=mock ./...

# コンテナ起動が必要な統合テスト
.PHONY: test-integration
test-integration: all
	- docker compose -p $(NAME) exec app go test -v -cover ./...
	$(MAKE) clean-volumes

# イメージを構築
.PHONY: build
build:
	docker compose -p $(NAME) build

# キャッシュを使わずにイメージを構築
.PHONY: ncbuild
ncbuild:
	docker compose -p $(NAME) build --no-cache

# イメージを構築（詳細なプログレス情報を表示）
.PHONY: build-verbose
build-verbose:
	docker compose -p $(NAME) build --progress=plain

# キャッシュを使わずにイメージを構築（詳細なプログレス情報を表示）
.PHONY: ncbuild-verbose
ncbuild-verbose:
	docker compose -p $(NAME) build --no-cache --progress=plain

# ホストのDocker環境をクリーンにする
.PHONY: all-clean-docker
all-clean-docker:
	docker ps -q | xargs -r docker stop
	docker ps -aq | xargs -r docker rm
	docker images -q | xargs -r docker rmi -f
	docker volume ls -q | xargs -r docker volume rm
	docker network ls --filter type=custom -q | xargs -r docker network rm
