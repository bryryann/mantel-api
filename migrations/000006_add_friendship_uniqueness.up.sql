ALTER TABLE friendships
ADD CONSTRAINT friendships_sender_receiver_unique
UNIQUE (sender_id, receiver_id);

CREATE INDEX IF NOT EXISTS friendships_sender_id_idx
ON friendships (sender_id);

CREATE INDEX IF NOT EXISTS friendships_receiver_id_idx
ON friendships (receiver_id);
