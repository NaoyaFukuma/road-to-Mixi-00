#!/usr/bin/env python3

import random

# ユーザーIDのリスト
user_ids = list(range(1, 101))

# 友達関係をランダムに生成
friend_links = []
block_list = []
for user1_id in user_ids:
    # 各ユーザーに対してランダムな数の友達関係を生成
    friends = random.sample([u for u in user_ids if u != user1_id], random.randint(1, 5))
    for user2_id in friends:
        friend_links.append((user1_id, user2_id))

    # 同様にブロックリストを生成
    blocks = random.sample([u for u in user_ids if u != user1_id and u not in friends], random.randint(1, 3))
    for user2_id in blocks:
        block_list.append((user1_id, user2_id))

# SQLステートメントの生成
friend_links_sql = "INSERT INTO friend_link (user1_id, user2_id) VALUES\n" + ",\n".join([f"({f[0]}, {f[1]})" for f in friend_links]) + ";"
block_list_sql = "INSERT INTO block_list (user1_id, user2_id) VALUES\n" + ",\n".join([f"({b[0]}, {b[1]})" for b in block_list]) + ";"

# SQLステートメントを出力
print(friend_links_sql)
print()
print(block_list_sql)
