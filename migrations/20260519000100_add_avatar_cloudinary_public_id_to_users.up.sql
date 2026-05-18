ALTER TABLE users
ADD COLUMN IF NOT EXISTS avatar_cloudinary_public_id TEXT;
