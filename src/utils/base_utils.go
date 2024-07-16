package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	bot "snowbreak_bot/config"
	"time"
)

var ctx = context.Background()

type GroupInvite struct {
	Id           string    `json:"id" gorm:"primaryKey"`
	GroupName    string    `json:"groupName"`
	GroupNumber  int64     `json:"groupNumber"`
	UserName     string    `json:"userName"`
	UserNumber   int64     `json:"userNumber"`
	MemberName   string    `json:"memberName"`
	MemberNumber int64     `json:"memberNumber"`
	CreateTime   time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime   time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark       string    `json:"remark"`
}

type GroupJoined struct {
	Id          string    `json:"id" gorm:"primaryKey"`
	GroupName   string    `json:"groupName"`
	GroupNumber int64     `json:"groupNumber"`
	News        int64     `json:"news"`
	CreateTime  time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime  time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark      string    `json:"remark"`
}

// SaveInvite 保存邀请记录
func SaveInvite(message *tgbotapi.Message, member *tgbotapi.User) {
	id, _ := gonanoid.New(32)
	groupInvite := GroupInvite{
		Id:           id,
		GroupName:    message.Chat.Title,
		GroupNumber:  message.Chat.ID,
		UserName:     message.From.FullName(),
		UserNumber:   message.From.ID,
		MemberName:   member.FullName(),
		MemberNumber: member.ID,
	}

	bot.DBEngine.Table("group_invite").Create(&groupInvite)
}

// SaveJoined 保存入群记录
func SaveJoined(message *tgbotapi.Message) {
	id, _ := gonanoid.New(32)
	groupJoined := GroupJoined{
		Id:          id,
		GroupName:   message.Chat.Title,
		GroupNumber: message.Chat.ID,
		News:        0,
	}

	bot.DBEngine.Table("group_joined").Create(&groupJoined)
}

// GetJoinedGroups 获取加入的群组
func GetJoinedGroups() []int64 {
	var groups []int64
	bot.DBEngine.Raw("select group_number from group_joined where news = 1 group by group_number").Scan(&groups)
	return groups
}

// GetJoinedByChatId 查询入群记录
func GetJoinedByChatId(chatId int64) *gorm.DB {
	return bot.DBEngine.Raw("select * from group_joined where group_number = ? limit 1", chatId)
}

func GetImg(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("获取图片失败", err)
		return nil
	}
	pic, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	return pic
}

// RedisSet redis存值
func RedisSet(key string, val interface{}, expiration time.Duration) {
	err := bot.GoRedis.Set(ctx, key, val, expiration).Err()
	if err != nil {
		log.Println(err)
	}
}

// RedisGet redis取值
func RedisGet(key string) string {
	val, err := bot.GoRedis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ""
		}
		log.Println(err)
	}
	return val
}

// RedisIsExists 判断redis值是否存在
func RedisIsExists(key string) bool {
	val := RedisGet(key)
	if val == "" {
		return false
	}
	return true
}

// RedisDel redis根据key删除
func RedisDel(key string) {
	err := bot.GoRedis.Del(ctx, key).Err()
	if err != nil {
		log.Println(err)
	}
}

// RedisScanKeys 扫描匹配keys
func RedisScanKeys(match string) (*redis.ScanIterator, context.Context) {
	return bot.GoRedis.Scan(ctx, 0, match, 0).Iterator(), ctx
}

// RedisSetList redis添加链表元素
func RedisSetList(key string, val interface{}) {
	err := bot.GoRedis.RPush(ctx, key, val).Err()
	if err != nil {
		log.Println(err)
	}
}

// RedisGetList redis获取所有链表元素
func RedisGetList(key string) []string {
	val, err := bot.GoRedis.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		log.Println(err)
	}
	return val
}

// RedisDelListItem redis移除链表元素
func RedisDelListItem(key string, val string) {
	err := bot.GoRedis.LRem(ctx, key, 0, val).Err()
	if err != nil {
		log.Println(err)
	}
}

// RedisAddSet redis集合添加元素
func RedisAddSet(key string, val string) {
	err := bot.GoRedis.SAdd(ctx, key, val).Err()
	if err != nil {
		log.Println(err)
	}
}

// RedisSetIsExists redis集合是否包含元素
func RedisSetIsExists(key string, val string) bool {
	exists, err := bot.GoRedis.SIsMember(ctx, key, val).Result()
	if err != nil {
		log.Println(err)
	}
	return exists
}

// RedisDelSetItem redis移除集合元素
func RedisDelSetItem(key string, val string) {
	err := bot.GoRedis.SRem(ctx, key, val).Err()
	if err != nil {
		log.Println(err)
	}
}
