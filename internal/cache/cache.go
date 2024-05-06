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
	DatabaseID int
	TelegramID int64
	Name       string
	Title      string
	KeyWords   map[int32]string
}

type UserCache struct {
	TelegramID      int64
	State           uint32
	Channels        map[int64]*ChannelCache
	DataLoaded      bool
	ActiveMenuID    int
	ActiveChannelID int64
	Mutex           sync.RWMutex
}

func CreateUser(user *tg.User) (*UserCache, error) {
	cache := UserCache{TelegramID: user.ID, State: StateNone, Channels: make(map[int64]*ChannelCache)}
	err := cache.loadUserData(user)
	return &cache, err
}

func (uc *UserCache) AddGroup(channel *tg.Channel) (*ChannelCache, error) {
	db := app.GetDatabase()

	if cache, exists := uc.Channels[channel.ID]; exists {
		return cache, nil
	}

	result, err := db.Exec("INSERT INTO `groups` (`id`, `userid`, `telegram_id`, `name`, `title`) VALUES (NULL, ?, ?, ?, ?)",
		uc.TelegramID, channel.ID, channel.Username, channel.Title)
	if err != nil {
		return nil, fmt.Errorf("error adding group to database: %v", err)
	}

	channelDBID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error get insert channel id: %v", err)
	}

	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()

	cache := ChannelCache{
		DatabaseID: int(channelDBID),
		TelegramID: channel.ID,
		Name:       channel.Username,
		Title:      channel.Title,
		KeyWords:   make(map[int32]string),
	}

	uc.Channels[channel.ID] = &cache
	return &cache, nil
}

func (uc *UserCache) RemoveGroup(id int64) error {
	db := app.GetDatabase()

	_, err := db.Exec("DELETE FROM `groups` WHERE `userid` = ? AND `telegram_id` = ?", uc.TelegramID, id)
	if err != nil {
		return fmt.Errorf("error remove group: %v", err)
	}

	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()

	delete(uc.Channels, id)
	return nil
}

type groupData struct {
	groupID    int
	telegramID int64
	groupName  string
	groupTitle string
	keywordID  sql.NullInt32
	keyword    sql.NullString
}

func (u *UserCache) loadUserData(user *tg.User) error {
	db := app.GetDatabase()

	rows, err := db.Query("SELECT g.id, g.telegram_id, g.name, g.title, k.id, k.keyword FROM `groups` g LEFT JOIN `keywords` k ON g.id = k.groupid WHERE g.userid = ?", u.TelegramID)
	if err != nil {
		return fmt.Errorf("error select user data: %v", err)
	}

	defer rows.Close()

	groups := make(map[int64]ChannelCache)
	for rows.Next() {
		var groupData groupData

		err := rows.Scan(&groupData.groupID, &groupData.telegramID, &groupData.groupName, &groupData.groupTitle, &groupData.keywordID, &groupData.keyword)
		if err != nil {
			return fmt.Errorf("error scan user data: %v", err)
		}

		group, ok := groups[groupData.telegramID]
		if !ok {
			group = ChannelCache{KeyWords: make(map[int32]string)}
			group.DatabaseID = groupData.groupID
			group.Title = groupData.groupTitle
			group.Name = groupData.groupName
			groups[groupData.telegramID] = group
		}

		if groupData.keywordID.Valid {
			group.KeyWords[groupData.keywordID.Int32] = groupData.keyword.String
		}
	}

	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	u.DataLoaded = true
	for id, channel := range groups {
		u.Channels[id] = &channel
	}

	return nil
}

func (channel *ChannelCache) AddKeyword(keyword string) error {
	db := app.GetDatabase()

	result, err := db.Exec("INSERT INTO `keywords` (`id`, `groupid`, `keyword`) VALUES (NULL, ?, ?)", channel.DatabaseID, keyword)
	if err != nil {
		return fmt.Errorf("error adding keyword: %v", err)
	}

	keywordID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	channel.KeyWords[int32(keywordID)] = keyword
	return nil
}

func (channel *ChannelCache) RemoveKeyword(keywordID int32) error {
	db := app.GetDatabase()

	_, err := db.Exec("DELETE FROM `keywords` WHERE `id` = ?", keywordID)
	if err != nil {
		return fmt.Errorf("error remove keyword: %v", err)
	}

	delete(channel.KeyWords, keywordID)
	return nil
}

func (uc *UserCache) HasChannelByID(channelID int64) bool {
	uc.Mutex.RLock()
	defer uc.Mutex.RUnlock()
	_, ok := uc.Channels[channelID]
	return ok
}

func (uc *UserCache) SetState(state uint32) {
	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()
	uc.State = state
}

func (uc *UserCache) GetState() uint32 {
	uc.Mutex.RLock()
	defer uc.Mutex.RUnlock()
	return uc.State
}
