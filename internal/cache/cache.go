package cache

import (
	"sync"
	"tg_reader_bot/internal/app"
)

const (
	StateNone = iota
	WaitingChannelName
	WaitingKeyWord
)

type ChannelCache struct {
	Name     string
	KeyWords map[int64]string
}

type UserCache struct {
	ID           int64
	State        uint32
	ActiveMenuID int
	Channels     map[int64]ChannelCache
	Mutex        sync.RWMutex
}

func CreateUser(userID int64) *UserCache {
	cache := UserCache{Channels: make(map[int64]ChannelCache)}
	cache.fillUserChannels()
	return &cache
}

func (uc *UserCache) AddGroup(groupID int64) {
	db := app.GetDatabase()

	_, err := db.Exec("INSERT INTO `groups` (`id`, `userid`, `groupid`) VALUES (NULL, ?, ?)", uc.ID, groupID)
	if err != nil {
		return
	}

	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()

	uc.Channels[groupID] = ChannelCache{}
}

func (uc *UserCache) AddKeyword(groupID int64, keyword string) {
	db := app.GetDatabase()

	result, err := db.Exec("INSERT INTO `keywords` (`id`, `groupid`, `keyword`) VALUES (NULL, ?, ?)", groupID, keyword)
	if err != nil {
		return
	}

	keywordID, err := result.LastInsertId()
	if err != nil {
		return
	}

	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()

	if channel, ok := uc.Channels[groupID]; ok {
		channel.KeyWords[keywordID] = keyword
		uc.Channels[groupID] = channel
	} else {
		uc.Channels[groupID] = ChannelCache{
			KeyWords: map[int64]string{
				keywordID: keyword,
			},
		}
	}
}

func (uc *UserCache) RemoveKeyword(groupID, keywordID int64) {
	db := app.GetDatabase()

	_, err := db.Exec("DELETE FROM `keywords` WHERE `groupid` = ? AND `id` = ?", groupID, keywordID)
	if err != nil {
		return
	}

	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()

	if channel, ok := uc.Channels[groupID]; ok {
		delete(channel.KeyWords, keywordID)
		uc.Channels[groupID] = channel
	}
}

func (uc *UserCache) RemoveGroup(groupID int64) {
	db := app.GetDatabase()

	_, err := db.Exec("DELETE FROM `groups` WHERE `userid` = ? AND `groupid` = ?", uc.ID, groupID)
	if err != nil {
		return
	}

	uc.Mutex.Lock()
	defer uc.Mutex.Unlock()

	delete(uc.Channels, groupID)
}

func (u *UserCache) fillUserChannels() {
	db := app.GetDatabase()

	rows, err := db.Query("SELECT g.groupid, k.id, k.keyword FROM `groups` g LEFT JOIN `keywords` k ON g.groupid = k.groupid WHERE g.userid = ?", u.ID)
	if err != nil {
		return
	}
	defer rows.Close()

	groups := make(map[int64]ChannelCache)
	for rows.Next() {
		var groupid, keywordID int64
		var keyword string
		err := rows.Scan(&groupid, &keywordID, &keyword)
		if err != nil {
			return
		}

		channel, ok := groups[groupid]
		if !ok {
			channel = ChannelCache{}
		}
		channel.KeyWords[keywordID] = keyword
		groups[groupid] = channel
	}

	u.Mutex.Lock()
	defer u.Mutex.Unlock()

	for id, channel := range groups {
		u.Channels[id] = channel
	}
}

func (u *UserCache) HasChannelByID(channelID int64) bool {
	u.Mutex.RLock()
	defer u.Mutex.RUnlock()
	_, ok := u.Channels[channelID]
	return ok
}
