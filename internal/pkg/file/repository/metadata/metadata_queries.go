package repository

const (
	GetFilesListQuery = `
		SELECT id, owner_id, filename, content_type, "size", upload_time, update_time, storage_key, 
		       is_deleted, deleted_time
		FROM public.file_metadata
		WHERE owner_id = $1 AND is_deleted = $2
		LIMIT $3 OFFSET $4;
	`

	GetMetadataQuery = `
		SELECT id, owner_id, filename, content_type, "size", upload_time, update_time, storage_key, 
		       is_deleted, deleted_time
		FROM public.file_metadata
		WHERE id = $1 AND NOT(is_deleted);
	`

	UploadMetadataQuery = `
		INSERT INTO public.file_metadata (id, owner_id, filename, content_type, size, storage_key)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, owner_id, filename, content_type, "size", upload_time, update_time, storage_key, 
		       is_deleted, deleted_time;
	`

	UpdateFilenameQuery = `
		UPDATE public.file_metadata
		SET filename = $2
		WHERE id = $1;
	`

	UpdateSizeQuery = `
		UPDATE public.file_metadata
		SET size = $2
		WHERE id = $1;
	`

	DeleteFileQuery = `
		UPDATE public.file_metadata
		SET is_deleted = TRUE
		WHERE id = $1;
	`
)
