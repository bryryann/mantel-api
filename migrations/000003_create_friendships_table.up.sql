CREATE TABLE IF NOT EXISTS friendships (
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status TEXT NOT NULL CHECK (status IN ('pending', 'accepted', 'blocked')) DEFAULT 'pending',
    version INTEGER NOT NULL DEFAULT 1,

    PRIMARY KEY (user_id, friend_id),
    CHECK (user_id <> friend_id)
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
