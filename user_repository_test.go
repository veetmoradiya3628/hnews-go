package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLUserRepository_CreateUser(t *testing.T) {
	defer cleanupTestData(t)

	repo := NewSQLUserRepository(testDB)

	userID, err := repo.CreateUser("John Doe", "john@doe.com", "testpassword", "avatar")
	assert.Nil(t, err)
	assert.Greater(t, userID, 0)

	user, err := repo.GetUserByEmail("john@doe.com")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)

	// userID, err = repo.CreateUser("John Doe", "john2@doe.com", generateString(73), "avatar")
	// assert.Error(t, err)
	// assert.Equal(t, userID, 0)
}

func TestSQLUserRepository_DuplicateEmail(t *testing.T) {
	defer cleanupTestData(t)

	repo := NewSQLUserRepository(testDB)

	userID, err := repo.CreateUser("John Doe", "john@doe.com", "testpassword", "avatar")
	assert.Nil(t, err)
	assert.Greater(t, userID, 0)

	_, err = repo.CreateUser("John Doe", "john@doe.com", "testpassword", "avatar")
	assert.Error(t, err)
}

func TestSQLUserRepository_Authentication_Success(t *testing.T) {
	defer cleanupTestData(t)

	repo := NewSQLUserRepository(testDB)

	currUserID, err := repo.CreateUser("John Doe", "john@doe.com", "testpassword", "avatar")
	assert.Nil(t, err)
	assert.Greater(t, currUserID, 0)

	authUserID, err := repo.Authenticate("john@doe.com", "testpassword")
	assert.NoError(t, err)
	assert.Equal(t, currUserID, authUserID)
}

func TestSQLUserRepository_Authentication_WrongPassword(t *testing.T) {
	defer cleanupTestData(t)

	repo := NewSQLUserRepository(testDB)

	currUserID, err := repo.CreateUser("John Doe", "john@doe.com", "testpassword", "avatar")
	assert.Nil(t, err)
	assert.Greater(t, currUserID, 0)

	_, err = repo.Authenticate("john@doe.com", "testpassword1")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidCredential, err)
}

func generateString(n int) string {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = 'a'
	}
	return string(buf)
}
