package repository

import "github.com/acoshift/web5/entity"

// FindUserByEmail finds user by email
func FindUserByEmail(q Queryer, email string) (*entity.User, error) {
	var x entity.User
	err := q.QueryRow(`
		select
			id, email, password
		from users
		where email = $1
	`, email).Scan(
		&x.ID, &x.Email, &x.Password,
	)
	if err != nil {
		return nil, err
	}
	return &x, nil
}
