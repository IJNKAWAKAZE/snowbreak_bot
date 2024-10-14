package utils

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
)

type Character struct {
	Name     string `json:"name"`     // 名字
	ThumbURL string `json:"thumbURL"` // 立绘
}

type Weapon struct {
	Name     string `json:"name"`     // 名字
	ThumbURL string `json:"thumbURL"` // 立绘
	Url      string `json:"url"`      // 网址
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

func GetCharacterByName(name string) Character {
	var char Character
	path := "./assets/strategy"
	d, _ := os.Open(path)
	fs, _ := d.Readdir(-1)
	for _, f := range fs {
		n := f.Name()[:len(f.Name())-4]
		if strings.Contains(n, name) {
			char.Name = name
			char.ThumbURL = path + "/" + f.Name()
			break
		}
	}
	return char
}

func GetCharactersByName(name string) []Character {
	var characterList []Character
	var characters []Character
	charactersJson := RedisGet("characterList")
	json.Unmarshal([]byte(charactersJson), &characters)
	for _, char := range characters {
		if strings.Contains(char.Name, name) {
			characterList = append(characterList, char)
		}
	}
	sort.Slice(characterList, func(i, j int) bool {
		return characterList[i].Name > characterList[j].Name
	})
	return characterList
}

func GetWeaponsByName(name string) []Weapon {
	var weaponList []Weapon
	var weapons []Weapon
	weaponsJson := RedisGet("weaponList")
	json.Unmarshal([]byte(weaponsJson), &weapons)
	for _, weapon := range weapons {
		if strings.Contains(weapon.Name, name) {
			weaponList = append(weaponList, weapon)
		}
	}
	sort.Slice(weaponList, func(i, j int) bool {
		return weaponList[i].Name > weaponList[j].Name
	})
	return weaponList
}

func GetWeaponByName(name string) Weapon {
	var weapons []Weapon
	weaponsJson := RedisGet("weaponList")
	json.Unmarshal([]byte(weaponsJson), &weapons)
	for _, weapon := range weapons {
		if weapon.Name == name {
			return weapon
		}
	}
	return Weapon{}
}
