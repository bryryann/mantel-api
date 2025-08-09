CREATE TABLE IF NOT EXISTS friendships (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sender_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status TEXT NOT NULL CHECK (status IN ('pending', 'accepted', 'blocked')) DEFAULT 'pending',
    version INTEGER NOT NULL DEFAULT 1,

    CHECK (sender_id <> receiver_id)
    UNIQUE (LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id))
);

CREATE OR REPLACE FUNCTION update_friendship_version()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    NEW.version := OLD.version + 1;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_friendship_version
BEFORE UPDATE ON friendships
FOR EACH ROW
EXECUTE FUNCTION update_friendship_version();
