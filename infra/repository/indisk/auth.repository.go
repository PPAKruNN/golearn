package indisk

import (
	"encoding/json"
	"os"
	"path"
)

const (
	AUTH_DATA_FILENAME = "auth.json"
)

type AuthRepository struct {
	Auths      map[string]int
	pathToFile string
}

func NewAuthRepository(dir string) *AuthRepository {
	repo := &AuthRepository{
		Auths:      map[string]int{},
		pathToFile: path.Join(dir, AUTH_DATA_FILENAME),
	}

	repo.loadIntoMemory()

	return repo
}

func (r *AuthRepository) RegisterToken(token string, accountId int) {

	r.loadIntoMemory()

	r.Auths[token] = accountId

	r.save()
}

func (r *AuthRepository) DecodeToken(token string) (int, error) {
	r.loadIntoMemory()

	accountId := r.Auths[token]

	return accountId, nil
}

func (r *AuthRepository) openHandle() *os.File {
	handle := openOrCreateFile(r.pathToFile)
	return handle
}

func (r *AuthRepository) Reset() error {
	err := resetFile(r.pathToFile)
	return err
}

func (r *AuthRepository) save() {
	marshal, _ := json.MarshalIndent(r.Auths, "", "  ")
	saveInFile(r.pathToFile, marshal)
}

func (r *AuthRepository) loadIntoMemory() {
	// Load accounts from disk
	handle := r.openHandle()
	defer handle.Close()

	var auths map[string]int = map[string]int{}
	json.NewDecoder(handle).Decode(&auths)

	r.Auths = auths
}
