package main

import (
	"context"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
)

// Processor is a struct to process message

func (p Processor) MulProcessMessage(input string, data *dto.WSATMessageData) error {
	ctx := context.Background()
	cmd := message.ParseCommand(input)
	//userId := data.Author.ID
	userName := data.Author.Username
	//userInfo := NewUserInfoExample(userId,"",data.Author.Username,false,"",0);
	toCreate := &dto.MessageToCreate{
		Content: "默认回复" + message.Emoji(307),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}
	if cmd.Cmd == "单人模式" {
		toCreate.Content = "指令无效，请退出多人模式后，再尝试"
		p.sendReply(ctx, data.ChannelID, toCreate)
		return nil
	}
	if cmd.Cmd == "退出" {
		if !multipleIsOpen {
			return nil
		}
		toCreate.Content = EndMulGame()
		p.sendReply(ctx, data.ChannelID, toCreate)
		return nil
	} else if cmd.Cmd == "多人模式" {
		//if multipleIsOpen {
		//	toCreate.Content = "多人模式不可重复开启"
		//	p.sendReply(ctx, data.ChannelID, toCreate)
		//	return nil
		//}
		toCreate.Content = "多人成语接龙开启:\n" + MulFirstStartGame()
		p.sendReply(ctx, data.ChannelID, toCreate)
	} else {
		if !multipleIsOpen {
			return nil
		}
		result := MulContinueConcatenateDragon(cmd.Cmd, userName)
		if result.error == "" {
			mulInfo.lastWord = result.content
			toCreate.Content = result.content
		} else {
			toCreate.Content = result.error
		}
		p.sendReply(ctx, data.ChannelID, toCreate)
	}
	return nil
}
