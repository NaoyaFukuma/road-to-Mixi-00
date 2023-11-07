#!/bin/bash

# すべてのコンテナを停止する
docker stop $(docker ps -aq)

# すべてのコンテナを削除する
docker rm $(docker ps -aq)

# すべてのDockerイメージを削除する
docker rmi $(docker images -q)

# 未使用のボリュームを削除する
docker volume rm $(docker volume ls -qf dangling=true)

# 未使用のネットワークを削除する
docker network rm $(docker network ls -q)

# 未使用のリソースを一掃する
docker system prune -a --volumes
