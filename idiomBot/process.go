package main

import (
	"context"
	"fmt"
	"github.com/cao-guang/pinyin"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
	"strconv"
	"time"
)

// Processor is a struct to process message
type Processor struct {
	api openapi.OpenAPI
}

// ProcessMessage is a function to process message
func (p Processor) ProcessMessage(input string, data *dto.WSATMessageData) error {
	ctx := context.Background()
	cmd := message.ParseCommand(input)
	userId := data.Author.ID
	userInfo := userInfoMap[userId]
	toCreate := &dto.MessageToCreate{
		Content: "默认回复" + message.Emoji(307),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}
	if userInfo.openIdiom {
		if cmd.Cmd == "退出" {
			if !userInfo.openIdiom {
				return nil
			}
			toCreate.Content = EndGame(&userInfo)
			delete(userInfoMap, userInfo.userId)
			p.sendReply(ctx, data.ChannelID, toCreate)
			return nil
		}

	}
	switch cmd.Cmd {
	case "单人模式":
		if userInfo.openIdiom {
			toCreate.Content = "正在统计本局数据，即将重新开始单人接龙\n"
			toCreate.Content += "	" + EndGame(&userInfo) + "\n" + "新的一局开始:"
			userInfoMap[data.Author.ID] = NewUserInfoExample(data.Author.ID, "", data.Author.Username, true, "", 0, 3)
			userInfo = userInfoMap[data.Author.ID]
			toCreate.Content += FirstStartGame(&userInfo)
		} else {
			userInfo.openIdiom = true
			toCreate.Content = "单人成语接龙开始：\n" + FirstStartGame(&userInfo)
		}
		SaveUserData(userInfo)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "hi":
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "time":
		toCreate.Content = genReplyContent(data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "提示":
		if userInfo.lifeValue < 1 {
			toCreate.Content = "本局提示次数已经用完，不可提示"
			p.sendReply(ctx, data.ChannelID, toCreate)
			return nil
		}
		userInfoLastWordPinYin, err := pinyin.To_Py(userInfo.lastWord, " ", "")
		if err != nil {
			fmt.Println(err)
		}
		idiom := GetNextMatchIdiom(userInfoLastWordPinYin, &userInfo)
		userInfo.lifeValue--
		userInfo.lastWord = idiom
		toCreate.Content = idiom + "\n本局剩余提示次数为：" + strconv.Itoa(userInfo.lifeValue)
		SaveUserData(userInfo)
		p.sendReply(ctx, data.ChannelID, toCreate)
	default:
		if !userInfo.openIdiom {
			return nil
		}
		result := ContinueConcatenateDragon(cmd.Cmd, &userInfo)
		if result.error == "" {
			userInfo.lastWord = result.content
			toCreate.Content = result.content
			SaveUserData(userInfo)
		} else {
			toCreate.Content = result.error
		}
		p.sendReply(ctx, data.ChannelID, toCreate)
	}

	return nil
}

func genReplyContent(data *dto.WSATMessageData) string {
	var tpl = `你好：%s
在子频道 %s 收到消息。
收到的消息发送时时间为：%s
当前本地时间为：%s

消息来自：%s
`

	msgTime, _ := data.Timestamp.Time()
	return fmt.Sprintf(
		tpl,
		message.MentionUser(data.Author.ID),
		message.MentionChannel(data.ChannelID),
		msgTime, time.Now().Format(time.RFC3339),
	)
}
