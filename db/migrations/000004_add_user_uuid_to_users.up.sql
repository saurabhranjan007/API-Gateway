ALTER TABLE users ADD COLUMN user_uuid VARCHAR(36) NOT NULL;
UPDATE users SET user_uuid = gen_random_uuid();
ALTER TABLE users ALTER COLUMN user_uuid SET DEFAULT gen_random_uuid();

-- ALTER TABLE users ADD COLUMN user_uuid VARCHAR(36) NOT NULL UNIQUE;
