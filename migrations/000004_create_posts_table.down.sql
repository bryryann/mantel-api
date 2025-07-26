DROP INDEX IF EXISTS idx_posts_user_id;
DROP INDEX IF EXISTS idx_posts_created_id;

DROP TRIGGER IF EXISTS trigger_update_post_version ON posts;

DROP FUNCTION IF EXISTS update_post_version;

DROP TABLE IF EXISTS posts;
