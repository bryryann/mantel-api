DROP INDEX IF EXISTS friendships_receiver_id_idx;
DROP INDEX IF EXISTS friendships_sender_id_idx;

ALTER TABLE friendships
DROP CONSTRAINT IF EXISTS friendships_sender_receiver_unique;
