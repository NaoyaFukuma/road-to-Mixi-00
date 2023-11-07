-- usersテーブルへのテストデータ挿入
INSERT INTO users (user_id, name) VALUES (1, 'Alice');
INSERT INTO users (user_id, name) VALUES (2, 'Bob');
INSERT INTO users (user_id, name) VALUES (3, 'Charlie');

-- friend_linkテーブルへのテストデータ挿入
INSERT INTO friend_link (user1_id, user2_id) VALUES (1, 2);
INSERT INTO friend_link (user1_id, user2_id) VALUES (1, 3);
INSERT INTO friend_link (user1_id, user2_id) VALUES (2, 1);

-- block_listテーブルへのテストデータ挿入
INSERT INTO block_list (user1_id, user2_id) VALUES (1, 3);
