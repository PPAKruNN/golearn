package inmemory

type AuthRepository struct {
	Auths map[string]int
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		Auths: map[string]int{},
	}
}

func (r AuthRepository) RegisterToken(token string, accountId int) {
	r.Auths[token] = accountId
}

func (r AuthRepository) DecodeToken(token string) (int, error) {
	accountId := r.Auths[token]

	return accountId, nil
}

func (r AuthRepository) Reset() error {
	r.Auths = map[string]int{}
	return nil
}
