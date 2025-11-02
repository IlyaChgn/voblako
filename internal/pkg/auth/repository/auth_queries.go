package repository

const (
	GetUserByEmailQuery = `
		SELECT u.id, u.email, u.password_hash
		FROM public.user u
		WHERE u.email = $1;
	`

	CreateUserQuery = `
		INSERT
		INTO public.user (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email;
	`
)
