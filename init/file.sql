\c voblako_db;

CREATE TABLE IF NOT EXISTS public.file_metadata (
    id UUID PRIMARY KEY UNIQUE NOT NULL,
    owner_id INT NOT NULL,
    filename TEXT NOT NULL
        CHECK (filename <> '')
        CONSTRAINT max_len_email CHECK(LENGTH(filename) <= 50),
    content_type TEXT NOT NULL
        CHECK (content_type <> ''),
    size BIGINT NOT NULL,
    upload_time TIMESTAMP DEFAULT NOW() NOT NULL,
    update_time TIMESTAMP DEFAULT NOW() NOT NULL
        CONSTRAINT updated_time_after_created_time CHECK (update_time >= upload_time),
    storage_key TEXT UNIQUE NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_time TIMESTAMP DEFAULT NULL
        CONSTRAINT deleted_time_after_created_time CHECK (deleted_time >= upload_time)
);

CREATE OR REPLACE FUNCTION change_metadata_update_time()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.update_time := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_metadata_trigger
    BEFORE UPDATE ON public.file_metadata
    FOR EACH ROW
EXECUTE PROCEDURE change_metadata_update_time();

CREATE OR REPLACE FUNCTION set_deleted_time()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_deleted = TRUE AND OLD.is_deleted = FALSE THEN
        NEW.deleted_time := NOW();
    END IF;

    IF NEW.is_deleted = FALSE AND OLD.is_deleted = TRUE THEN
        NEW.deleted_time := NULL;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_deleted_time_trigger
    BEFORE UPDATE ON public.file_metadata
    FOR EACH ROW
EXECUTE FUNCTION set_deleted_time();

