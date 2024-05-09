package cache

import (
	"fmt"
	"tg_reader_bot/internal/app"
)

type rowGroups struct {
	id        int64
	userID    int64
	channelID int64
	lastMsgID int
	name      string
	title     string
}

type rowKeywords struct {
	databaseID int64
	groupFK    int64
	keyword    string
}

type userKeyWords struct {
	userID  int64
	channel *PeerInfo
}

func (cm *PeersManager) LoadUsersData() {
	db := app.GetDatabase()

	groupsRows, err := db.Query("SELECT `id`, `userid`, `groupid`, `lastmsgid`, `name`, `title` FROM `groups`")
	if err != nil {
		panic(err)
	}
	defer groupsRows.Close()

	channelKeyWords := make(map[int64]userKeyWords)

	for groupsRows.Next() {
		var group rowGroups
		if err := groupsRows.Scan(&group.id, &group.userID, &group.channelID, &group.lastMsgID, &group.name, &group.title); err != nil {
			panic(err)
		}

		channel, _ := cm.AddChannelToUser(group.userID, group.id, group.channelID, group.lastMsgID, group.name, group.title, false)
		channelKeyWords[group.id] = userKeyWords{userID: group.userID, channel: channel}
	}

	keywordsRows, err := db.Query("SELECT `id`, `group_fk`, keyword FROM `keywords`")
	if err != nil {
		panic(err)
	}

	defer keywordsRows.Close()

	for keywordsRows.Next() {
		var keyword rowKeywords
		if err := keywordsRows.Scan(&keyword.databaseID, &keyword.groupFK, &keyword.keyword); err != nil {
			panic(err)
		}

		if data, ok := channelKeyWords[keyword.groupFK]; ok {
			data.channel.AddKeyword(data.userID, keyword.databaseID, keyword.keyword, false)
		}
	}
}

func (cm *PeersManager) AddChannelToUser(userID int64, databaseID int64, channelID int64, lastmsgid int, name, title string, addToDB bool) (*PeerInfo, error) {
	if addToDB {
		db := app.GetDatabase()

		result, err := db.Exec("INSERT INTO `groups` (`id`, `userid`, `groupid`, `lastmsgid`, `name`, `title`) VALUES (NULL, ?, ?, ?, ?, ?)",
			userID, channelID, lastmsgid, name, title)

		if err != nil {
			return nil, fmt.Errorf("error adding group to database: %v", err)
		}

		databaseID, err = result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("error get insert channel id: %v", err)
		}

		cm.Mutex.Lock()
		defer cm.Mutex.Unlock()
	}

	user, ok := cm.Users[userID]
	if !ok {
		user = &UserData{
			// TelegramID: userID, TODO
			Channels: make(map[int64]*PeerInfo),
		}
		cm.Users[userID] = user
	}

	channel, ok := cm.Peers[channelID]
	if !ok {
		channel = &PeerInfo{
			UsersKeyWords: make(map[int64]*PeerKeyWords),
			DatabaseID:    databaseID,
			// TelegramID:    channelID, TODO
			Name:      name,
			Title:     title,
			LastMsgID: lastmsgid,
		}
		cm.Peers[channelID] = channel
	}

	user.Channels[channelID] = channel

	return channel, nil
}

func (cm *PeersManager) RemoveChannelFromUser(userID int64, channelID int64, addToDB bool) error {
	if addToDB {
		db := app.GetDatabase()
		_, err := db.Exec("DELETE FROM `groups` WHERE `userid` = ? AND `groupid` = ?", userID, channelID)
		if err != nil {
			return fmt.Errorf("error remove group: %v", err)
		}

		cm.Mutex.Lock()
		defer cm.Mutex.Unlock()
	}

	if user, ok := cm.Users[userID]; ok {
		delete(user.Channels, channelID)
	}

	if channel, ok := cm.Peers[channelID]; ok {
		delete(channel.UsersKeyWords, userID)
	}

	return nil
}

func (channel *PeerInfo) AddKeyword(userID int64, keywordID int64, keyword string, addToDB bool) error {
	if addToDB {
		db := app.GetDatabase()

		result, err := db.Exec("INSERT INTO `keywords` (`id`, `group_fk`, `keyword`) VALUES (NULL, ?, ?)", channel.DatabaseID, keyword)
		if err != nil {
			return fmt.Errorf("error adding keyword: %v", err)
		}

		keywordID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	}

	usersKeywords, ok := channel.UsersKeyWords[userID]
	if ok {
		usersKeywords.Keywords[keywordID] = keyword
	} else {
		usersKeywords = &PeerKeyWords{Keywords: make(map[int64]string)}
		usersKeywords.Keywords[keywordID] = keyword
		channel.UsersKeyWords[userID] = usersKeywords
	}

	return nil
}

func (channel *PeerInfo) RemoveKeyword(userID int64, keywordID int64, addToDB bool) error {
	if addToDB {
		db := app.GetDatabase()

		_, err := db.Exec("DELETE FROM `keywords` WHERE `id` = ?", keywordID)
		if err != nil {
			return fmt.Errorf("error remove keyword: %v", err)
		}
	}

	if usersKeywords, ok := channel.UsersKeyWords[userID]; ok {
		delete(usersKeywords.Keywords, keywordID)
	}

	return nil
}

func (cm *PeersManager) GetUserData(userID int64, create bool) *UserData {
	if user, ok := cm.Users[userID]; ok {
		return user
	} else if create {
		user := &UserData{Channels: map[int64]*PeerInfo{}} // TODO TelegramID: userID,
		cm.Users[userID] = user
		return user
	}
	return nil
}

func (channel *PeerInfo) GetUserKeyWords(userID int64) *map[int64]string {
	if user, ok := channel.UsersKeyWords[userID]; ok {
		return &user.Keywords
	}
	return nil
}

func (channel *PeerInfo) GetUserKeyWordsCount(userID int64) int {
	keyWords := channel.GetUserKeyWords(userID)
	if keyWords != nil {
		return len(*keyWords)
	}

	return 0
}

func (user *UserData) HasChannelByID(ID int64) bool {
	_, ok := user.Channels[ID]
	return ok
}

func (user *UserData) GetActivePeerID() int64 {
	return user.ActivePeerID
}

func (user *UserData) GetActiveChannel() *PeerInfo {
	return user.Channels[user.GetActivePeerID()]
}

func (user *UserData) GetID() int64 {
	return user.Peer.UserID
}
