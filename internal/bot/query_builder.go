package bot

import "encoding/json"

func (b *Bot) buildQuery(action uint32, data string) []byte {
	res, _ := json.Marshal(QueryHeader{Action: action, Time: b.startTime, Data: data})
	return res
}
