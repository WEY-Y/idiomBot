package main

import (
	"context"
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

// 消息处理器，持有 openapi 对象
var processor Processor

func main() {
	IdiomDataInit()
	ctx := context.Background()

	botToken := token.New(token.TypeBot) //创建空token
	if err := botToken.LoadFromConfig(getConfigPath("config.yaml")); err != nil {
		log.Fatalln(err)
	}

	api := botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)

	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}

	processor = Processor{api: api}

	intent := websocket.RegisterHandlers(
		ATMessageEventHandler(),
		ReadyHandler(),
		ErrorNotifyHandler(),
	)
	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}

// ReadyHandler 自定义 ReadyHandler 感知连接成功事件
func ReadyHandler() event.ReadyHandler {
	return func(event *dto.WSPayload, data *dto.WSReadyData) {
		log.Println("ready event receive: ", data)
	}
}

func ErrorNotifyHandler() event.ErrorNotifyHandler {
	return func(err error) {
		log.Println("error notify receive: ", err)
	}
}

// ATMessageEventHandler 实现处理 at 消息的回调
func ATMessageEventHandler() event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		if !IdiomDataInit() {
			fmt.Println("初始化配置失败")
			return nil
		}
		input := strings.ToLower(message.ETLInput(data.Content))
		log.Println("data.content", data.Content)

		if input == "多人模式" {
			if !multipleIsOpen {
				multipleIsOpen = true
				mulInfo = NewMulInfoExample(data.ChannelID, "", "", make(map[string]int))
			} else {
				toCreate := &dto.MessageToCreate{
					Content: "多人模式不可重复开启",
					MessageReference: &dto.MessageReference{
						// 引用这条消息
						MessageID:             data.ID,
						IgnoreGetMessageError: true,
					},
				}
				processor.sendReply(context.Background(), data.ChannelID, toCreate)
				return nil
			}
		}
		if multipleIsOpen {
			_, ok := mulInfo.userInfoMap[data.Author.Username]
			if !ok {
				mulInfo.userInfoMap[data.Author.Username] = 0
			}
			return processor.MulProcessMessage(input, data)
		} else {
			userInfo, exist := userInfoMap[data.Author.ID]
			fmt.Println(userInfo.userName)
			if exist == false {
				userInfoMap[data.Author.ID] = NewUserInfoExample(data.Author.ID, "", data.Author.Username, false, "", 0, 3)
			}
			return processor.ProcessMessage(input, data)
		}

	}
}

func getConfigPath(name string) string {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("%s/%s", path.Dir(filename), name)
	}
	return ""
}
