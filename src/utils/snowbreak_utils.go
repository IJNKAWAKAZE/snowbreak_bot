package utils

import (
	"encoding/json"
)

type Character struct {
	Name     string `json:"name"`     // 名字
	ThumbURL string `json:"thumbURL"` // 立绘
}

func GetCharacters() []Character {
	var characters []Character
	charactersJson := RedisGet("characterList")
	json.Unmarshal([]byte(charactersJson), &characters)
	return characters
}
