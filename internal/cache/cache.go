package cache

import (
	"database/sql"
	"fmt"
	"sync"
	"tg_reader_bot/internal/app"

	"github.com/gotd/td/tg"
)

const (
	StateNone = iota
	WaitingChannelName
	WaitingKeyWord
)

type ChannelCache struct {
	ID         int
	TelegramID int64
	Name       string
	KeyWords   map[int32]string
}

type UserCache struct {
	ID           int64
	State        uint32
	ActiveMenuID int
	Channels     map[string]*ChannelCache
	DataLoaded   bool
	Mutex        sync.RWMutex
}

func CreateUser(user *tg.User) (*UserCache, error) {
	cache := UserCache{ID: user.ID, State: StateNone, Channels: make(map[string]*ChannelCache)}
	err := cache.loadUserData(user)
	return &cache, err
}

func (uc *UserCache) AddGroup(channel *tg.Channel) (*ChannelCache, error) {
	db := app.GetDatabase()

	if cache, exists := uc.Channels[channel.Username]; exists {
		return cache, nil
	}

	_, err := db.Exec("INSERT INTO `groups` (`id`, `userid`, `groupname`) VALUES (NULL, ?, ?)", uc.ID, channel.Username)
	if err != nil {
		return nil, fmt.Errorf("error adding group to database: %v", err)
	}

	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()

	cache := ChannelCache{KeyWords: make(map[int32]string)}
	uc.Channels[channel.Username] = &cache
	return &cache, nil
}

func (uc *UserCache) RemoveGroup(groupName string) error {
	db := app.GetDatabase()

	_, err := db.Exec("DELETE FROM `groups` WHERE `userid` = ? AND `groupname` = ?", uc.ID, groupName)
	if err != nil {
		return fmt.Errorf("error remove group: %v", err)
	}

	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()

	delete(uc.Channels, groupName)
	return nil
}

type groupData struct {
	groupID   int
	groupName string
	keywordID sql.NullInt32
	keyword   sql.NullString
}

func (u *UserCache) loadUserData(user *tg.User) error {
	db := app.GetDatabase()

	_, err := db.Exec("INSERT INTO users (id, userid, username) VALUES (NULL, ?, ?) ON DUPLICATE KEY UPDATE username = ?", u.ID, user.Username, user.Username)
	if err != nil {
		return fmt.Errorf("error creating/updating user in database: %v", err)
	}

	rows, err := db.Query("SELECT g.id, g.groupname, k.id, k.keyword FROM `groups` g LEFT JOIN `keywords` k ON g.id = k.groupid WHERE g.userid = ?", u.ID)
	if err != nil {
		return fmt.Errorf("error select user data: %v", err)
	}

	defer rows.Close()

	groups := make(map[string]ChannelCache)
	for rows.Next() {
		var data groupData

		err := rows.Scan(&data.groupID, &data.groupName, &data.keywordID, &data.keyword)
		if err != nil {
			return fmt.Errorf("error scan user data: %v", err)
		}

		group, ok := groups[data.groupName]
		if !ok {
			group = ChannelCache{KeyWords: make(map[int32]string)}
			group.ID = data.groupID
			groups[data.groupName] = group
		}

		if data.keywordID.Valid {
			group.KeyWords[data.keywordID.Int32] = data.keyword.String
		}
	}
	fmt.Println("user groups ", groups)

	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.DataLoaded = true
	for name, channel := range groups {
		u.Channels[name] = &channel
	}

	return nil
}

func (u *UserCache) HasChannelByID(channelID int64) bool {
	for _, channel := range u.Channels {
		if channel.TelegramID == channelID {
			return true
		}
	}
	return false
}

func (uc *UserCache) SetState(state uint32) {
	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()
	uc.State = state
}

func (cache *ChannelCache) AddKeyword(keyword string) error {
	db := app.GetDatabase()

	result, err := db.Exec("INSERT INTO `keywords` (`id`, `groupid`, `keyword`) VALUES (NULL, ?, ?)", cache.ID, keyword)
	if err != nil {
		return fmt.Errorf("error adding keyword: %v", err)
	}

	keywordID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	cache.KeyWords[int32(keywordID)] = keyword
	return nil
}

func (uc *ChannelCache) RemoveKeyword(keywordID int32) error {
	db := app.GetDatabase()

	_, err := db.Exec("DELETE FROM `keywords` WHERE `id` = ?", keywordID)
	if err != nil {
		return fmt.Errorf("error remove keyword: %v", err)
	}

	delete(uc.KeyWords, keywordID)
	return nil
}

func (uc *UserCache) GetState() uint32 {
	uc.Mutex.RLock()
	defer uc.Mutex.RUnlock()
	return uc.State
}

func (uc *UserCache) SetActiveMenuID(menuID int) {
	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()
	uc.ActiveMenuID = menuID
}

func (uc *UserCache) GetActiveMenuID() int {
	uc.Mutex.RLock()
	defer uc.Mutex.RUnlock()
	return uc.ActiveMenuID
}
