package user

type User struct {
	ID   int64
	Name Name
}

func New(id int64, rawName string) (User, error) {
	name, err := NewName(rawName)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:   id,
		Name: name,
	}, nil
}
