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

func TestSQLUserRepository_GetUserByEmailWithProfile(t *testing.T) {
	defer cleanupTestData(t)

	repo := NewSQLUserRepository(testDB)
	userID, err := repo.CreateUser("John Doe", "john@doe.com", "testpassword", "avatar")
	assert.Nil(t, err)
	assert.Greater(t, userID, 0)

	userData, err := repo.GetUserByEmailWithProfile("john@doe.com")
	assert.NoError(t, err)
	assert.NotNil(t, userData)
	assert.NotZero(t, userData.Profile.UserID)
}

func TestSQLUserRepository_GetUsers(t *testing.T) {
	defer cleanupTestData(t)

	repo := NewSQLUserRepository(testDB)

	userID1, err := repo.CreateUser("John Doe", "john@doe.com", "testpassword", "avatar1")
	assert.NoError(t, err)
	assert.Greater(t, userID1, 0)

	userID2, err := repo.CreateUser("Jane Doe", "jane@doe.com", "testpassword", "avatar2")
	assert.NoError(t, err)
	assert.Greater(t, userID2, 0)

	users, err := repo.GetUsers()

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)

	assert.NotZero(t, users[0].ID)
	assert.NotEmpty(t, users[0].Email)
	assert.NotZero(t, users[0].Profile.UserID)

	assert.NotZero(t, users[1].ID)
	assert.NotEmpty(t, users[1].Email)
	assert.NotZero(t, users[1].Profile.UserID)
}

func generateString(n int) string {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = 'a'
	}
	return string(buf)
}
