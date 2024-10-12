package utils

import (
	"encoding/json"
	"os"
	"strings"
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

func GetLocalCharacters() []Character {
	var characters []Character
	path := "./assets/images"
	d, _ := os.Open(path)
	fs, _ := d.Readdir(-1)
	for _, f := range fs {
		var char Character
		char.Name = f.Name()[:len(f.Name())-4]
		char.ThumbURL = path + "/" + f.Name()
		characters = append(characters, char)
	}
	return characters
}

func GetCharacter(name string) Character {
	var char Character
	path := "./assets/strategy"
	d, _ := os.Open(path)
	fs, _ := d.Readdir(-1)
	for _, f := range fs {
		n := f.Name()[:len(f.Name())-4]
		if n == name {
			char.Name = name
			char.ThumbURL = path + "/" + f.Name()
			break
		}
	}
	return char
}

func GetCharactersByName(name string) []Character {
	var characterList []Character
	for _, char := range GetLocalCharacters() {
		if strings.Contains(char.Name, name) {
			characterList = append(characterList, char)
		}
	}
	return characterList
}
