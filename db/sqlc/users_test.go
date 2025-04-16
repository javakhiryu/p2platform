package db

import (
	"context"
	"p2platform/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		TelegramID: util.RandomInt(10000, 999999999),
		TgUsername: util.RandomTgUsername(),
		FirstName:  util.ToPgText(util.RandomString(7)),
		LastName:   util.ToPgText(util.RandomString(8)),
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.TelegramID, user.TelegramID)
	require.Equal(t, arg.TgUsername, user.TgUsername)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, arg.LastName)
	return user
}

func CreateUserWithExactId(t *testing.T, telegramId int64) User {
	arg := CreateUserParams{
		TelegramID: telegramId,
		TgUsername: util.RandomTgUsername(),
		FirstName:  util.ToPgText(util.RandomString(7)),
		LastName:   util.ToPgText(util.RandomString(8)),
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.TelegramID, user.TelegramID)
	require.Equal(t, arg.TgUsername, user.TgUsername)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, arg.LastName)
	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.TelegramID)
	require.NoError(t, err)
	require.Equal(t, user1.TelegramID, user2.TelegramID)
	require.Equal(t, user1.TgUsername, user2.TgUsername)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, user1.LastName, user2.LastName)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
}

func TestUpdateUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	newTgUsername :=util.ToPgText(util.RandomTgUsername())
	newLastName := util.ToPgText(util.RandomString(8))

	arg := UpdateUserParams{
		TelegramID: user1.TelegramID,
		TgUsername: newTgUsername,
		LastName: newLastName,
	}
	user2, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, user1.TelegramID, user2.TelegramID)
	require.Equal(t, newTgUsername.String, user2.TgUsername)
	require.Equal(t, user1.FirstName, user2.FirstName)
	require.Equal(t, newLastName, user2.LastName)
	require.Equal(t, user1.CreatedAt, user2.CreatedAt)
	require.WithinDuration(t, time.Now(), user2.UpdatedAt, 100*time.Millisecond)
}

func TestDeleteUser(t *testing.T){
	user := CreateRandomUser(t)
	err := testStore.DeleteUser(context.Background(), user.TelegramID)
	require.NoError(t, err)
}