
-- Adds a foreign key constraint to the posts table
-- Ensures user_id in posts must match an existing id in users table
-- Prevents orphaned posts and maintains referential integrity
-- Links posts.user_id to users.id
-- Ensures every post belongs to a valid user
ALTER TABLE
  posts
ADD
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id);
