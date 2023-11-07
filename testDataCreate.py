#!/usr/bin/env python3

from random import randint

# Pythonスクリプトで大量のテストデータを生成
user_count = 100  # 生成するユーザ数
friends_per_user = 5  # 1人のユーザにつき、何人の友達を生成するか
blocks_per_user = 2  # 1人のユーザにつき、何人をブロックするか

# usersテーブル用のデータ
with open('./mysql/2_insert_users.sql', 'w') as f:
    f.write('INSERT INTO users (user_id, name) VALUES\n')
    for i in range(1, user_count + 1):
        name = f"'User{i}'"
        if i < user_count:
            f.write(f"({i}, {name}),\n")
        else:
            f.write(f"({i}, {name});\n")

# friend_linkテーブル用のデータ
with open('./mysql/3_insert_friend_links.sql', 'w') as f:
    f.write('INSERT INTO friend_link (user1_id, user2_id) VALUES\n')
    for i in range(1, user_count + 1):
        friends = set()
        while len(friends) < friends_per_user:
            friend_id = randint(1, user_count)
            if friend_id != i and friend_id not in friends:
                friends.add(friend_id)
                f.write(f"({i}, {friend_id}),\n")

# 最後の行のカンマを取り除くには、後でファイルを編集するか、ロジックを追加してください。

# block_listテーブル用のデータ
# 同様に実装できます。
