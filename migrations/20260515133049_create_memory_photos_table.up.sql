CREATE TABLE IF NOT EXISTS memory_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    memory_id UUID NOT NULL REFERENCES memories(id) ON DELETE CASCADE,
    uploaded_by UUID NOT NULL REFERENCES users(id),
    media_type VARCHAR(20) NOT NULL DEFAULT 'image',
    photo_url TEXT NOT NULL,
    cloudinary_public_id TEXT,
    caption TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);