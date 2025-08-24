CREATE TABLE IF NOT EXISTS follows (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    follower_id integer NOT NULL,
    followee_id integer NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_follower FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_followee FOREIGN KEY (followee_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_follow UNIQUE (follower_id, followee_id),
    CONSTRAINT no_self_follow CHECK (follower_id <> followee_id)
);

CREATE INDEX IF NOT EXISTS idx_follower_id ON follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_followee_id ON follows(followee_id);
