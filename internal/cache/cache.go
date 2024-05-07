package cache

import (
	"fmt"
	"tg_reader_bot/internal/app"
)

func (cm *ChannelsManager) LoadUsersData() {
	db := app.GetDatabase()

	groupsRows, err := db.Query("SELECT `id`, `userid`, `groupid`, `name`, `title` FROM `groups`")
	if err != nil {
		panic(err)
	}
	defer groupsRows.Close()

	channelKeyWords := make(map[int64]*ChannelInfo)

	for groupsRows.Next() {
		var group rowGroups
		if err := groupsRows.Scan(&group.ID, &group.UserID, &group.TelegramID, &group.Name, &group.Title); err != nil {
			panic(err)
		}

		channel, _ := cm.AddChannelToUser(group.UserID, group.ID, group.TelegramID, group.Name, group.Title, false)
		channelKeyWords[group.ID] = channel
	}

	keywordsRows, err := db.Query("SELECT `id`, `group_fk`, keyword FROM `keywords`")
	if err != nil {
		panic(err)
	}

	defer keywordsRows.Close()

	for keywordsRows.Next() {
		var keyword rowKeywords
		if err := keywordsRows.Scan(&keyword.DatabaseID, &keyword.GroupFK, &keyword.Keyword); err != nil {
			panic(err)
		}

		if channel, ok := channelKeyWords[keyword.GroupFK]; ok {
			channel.AddKeyword(123, keyword.DatabaseID, keyword.Keyword, false)
		}
	}
}

func (cm *ChannelsManager) AddChannelToUser(userID int64, databaseID int64, channelID int64, name, title string, addToDB bool) (*ChannelInfo, error) {
	if addToDB {
		db := app.GetDatabase()

		result, err := db.Exec("INSERT INTO `groups` (`id`, `userid`, `groupid`, `name`, `title`) VALUES (NULL, ?, ?, ?, ?)",
			userID, channelID, name, title)

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
			TelegramID: userID,
			Channels:   make(map[int64]*ChannelInfo),
		}
		cm.Users[userID] = user
	}

	channel, ok := cm.Channels[channelID]
	if !ok {
		channel = &ChannelInfo{
			UsersKeyWords: make(map[int64]*ChannelKeyWords),
			DatabaseID:    databaseID,
			TelegramID:    channelID,
			Name:          name,
			Title:         title,
			LastParseTime: 0,
		}
		cm.Channels[channelID] = channel
	}

	user.Channels[channelID] = channel

	return channel, nil
}

func (cm *ChannelsManager) RemoveChannelFromUser(userID int64, channelID int64, addToDB bool) error {
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

	if channel, ok := cm.Channels[channelID]; ok {
		delete(channel.UsersKeyWords, userID)
	}

	return nil
}

func (channel *ChannelInfo) AddKeyword(userID int64, keywordID int64, keyword string, addToDB bool) error {
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
		usersKeywords = &ChannelKeyWords{Keywords: make(map[int64]string)}
		usersKeywords.Keywords[keywordID] = keyword
		channel.UsersKeyWords[userID] = usersKeywords
	}

	return nil
}

func (channel *ChannelInfo) RemoveKeyword(userID int64, keywordID int64, addToDB bool) error {
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

func (cm *ChannelsManager) GetUserData(userID int64, create bool) *UserData {
	if user, ok := cm.Users[userID]; ok {
		return user
	} else if create {
		user := &UserData{TelegramID: userID, Channels: map[int64]*ChannelInfo{}}
		cm.Users[userID] = user
		return user
	}
	return nil
}

func (cm *ChannelInfo) GetUserKeyWords(userID int64) *map[int64]string {
	if user, ok := cm.UsersKeyWords[userID]; ok {
		return &user.Keywords
	}
	return nil
}

func (cm *ChannelInfo) GetUserKeyWordsCount(userID int64) int {
	keyWords := cm.GetUserKeyWords(userID)
	if keyWords != nil {
		return len(*keyWords)
	}

	return 0
}

func (user *UserData) HasChannelByID(telegramID int64) bool {
	_, ok := user.Channels[telegramID]
	return ok
}

func (user *UserData) GetActiveChannel() *ChannelInfo {
	return user.Channels[user.ActiveChannelID]
}
