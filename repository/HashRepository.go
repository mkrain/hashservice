package repository

import (
	"crypto/sha512"
)

//HashRepository this is the keeper and generator of hashes
type HashRepository struct {
	hashes map[string][]byte
}

//NewHashRepository repsrents a backing store for password hashes
func NewHashRepository() *HashRepository {
	hash := HashRepository{
		hashes: make(map[string][]byte),
	}

	return &hash
}

//GetHash returns the hash for the given value
func (repo *HashRepository) GetHash(value string) []byte {
	return repo.hashes[value]
}

//ContainsKey return true if the key is in the repository, false otherwise
func (repo *HashRepository) ContainsKey(value string) bool {
	if _, ok := repo.hashes[value]; ok {
		return true
	}

	return false
}

//CreateHash creates a SHA-512 hash from the given string
func (repo *HashRepository) CreateHash(value string) []byte {
	if hash, ok := repo.hashes[value]; ok {
		return hash
	}

	sha512 := sha512.New512_256()
	bytes := []byte(value)
	sha512.Write(bytes)
	hash := sha512.Sum(nil)

	repo.hashes[value] = hash

	return hash
}
