package postgres

import (
	domainUser "go_pg_http/internal/domain/user"
)

func toUserRow(usr domainUser.User) UserRow {
	return UserRow{
		ID:   usr.ID,
		Name: usr.Name.String(),
	}
}

func toDomainUser(row UserRow) (domainUser.User, error) {
	return domainUser.New(row.ID, row.Name)
}
