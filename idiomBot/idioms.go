package main

import (
	"bufio"
	_ "container/list"
	"fmt"
	"github.com/cao-guang/pinyin"
	"github.com/deckarep/golang-set"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var multipleIsOpen = false
var mulInfo = NewMulInfo()
var idiomsMap = make(map[string]IdiomAndFirstPinYin)

var allFourLetterWordsSet = mapset.NewSet()
var userInfoMap = make(map[string]userInfo)

type IdiomAndFirstPinYin struct {
	idiom       string
	firstPinYin string
}

func NewIdiomAndFirstPinYin(idiom string, firstPinYin string) IdiomAndFirstPinYin {
	return IdiomAndFirstPinYin{
		idiom:       idiom,
		firstPinYin: firstPinYin,
	}
}

type userInfo struct {
	userId    string
	lastWord  string
	userName  string
	openIdiom bool
	path      string
	value     int
	lifeValue int
}

type result struct {
	content string
	error   string
}

type mulGameInfo struct {
	channelId   string
	lastWord    string
	path        string
	userInfoMap map[string]int
}

func NewMulInfo() mulGameInfo {
	return mulGameInfo{}
}
func NewMulInfoExample(channelId string, lastWord string, path string, userInfoMap map[string]int) mulGameInfo {
	return mulGameInfo{
		channelId:   channelId,
		lastWord:    lastWord,
		path:        path,
		userInfoMap: userInfoMap,
	}
}
func NewUserInfoExample(userId string, lastWord string, userName string, openIdiom bool, path string, value int, lifeValue int) userInfo {
	return userInfo{
		userId:    userId,
		lastWord:  lastWord,
		userName:  userId,
		openIdiom: openIdiom,
		path:      path,
		value:     value,
		lifeValue: lifeValue,
	}
}

func IdiomDataInit() bool {
	if GetIdiomsMap() && GetAllFourLetterWords() {
		return true
	}
	return false
}

func GetIdiomsMap() bool {
	if len(idiomsMap) > 1 {
		return true
	}
	file, err := os.Open("selectedidiom.txt")
	if err != nil {
		log.Println("读取selectedidiom.txt失败")
		return false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	pinyin.LoadingPYFileName("pinyin.txt")
	for scanner.Scan() {
		lineText := scanner.Text()
		wordPinYin, err := pinyin.To_Py(lineText, " ", "")
		wordPinYinArray := strings.Split(wordPinYin, " ")
		if err != nil {
			fmt.Println(err)
		}
		idiomsMap[wordPinYin] = NewIdiomAndFirstPinYin(lineText, wordPinYinArray[0])
	}
	return true
}

func GetAllFourLetterWords() bool {
	if allFourLetterWordsSet.Contains("谢格洛夫") {
		return true
	}
	file, err := os.Open("fourletterwords.txt")
	if err != nil {
		log.Println("读取fourletterwords.txt失败")
		return false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		allFourLetterWordsSet.Add(lineText)
	}
	return true
}
func FirstStartGame(userInfo *userInfo) string {
	idiom := GenerateRandomIdiom()
	userInfo.lastWord = idiom
	userInfo.path += idiom
	userInfo.path += "->"
	return idiom
}

func MulFirstStartGame() string {
	idiom := GenerateRandomIdiom()
	mulInfo.lastWord = idiom
	mulInfo.path += idiom
	mulInfo.path += "->"
	return idiom
}
func EndGame(userInfo *userInfo) string {
	userInfo.lastWord = ""
	userInfo.openIdiom = false
	var res = new(strings.Builder)
	res.WriteString("单人模式结束，想再玩时可以使用@单人模式\n")
	res.WriteString("游戏中接龙词语为：")
	rs := []rune(userInfo.path)
	end := len(rs) - 2
	res.WriteString(string(rs[0:end]))
	res.WriteString("\n最终得分为")
	res.WriteString(strconv.Itoa(userInfo.value))
	return res.String()
}

func EndMulGame() string {
	var res = new(strings.Builder)
	res.WriteString("游戏结束，想再玩时可以使用@多人模式\n")
	res.WriteString("游戏中接龙词语为：")
	rs := []rune(mulInfo.path)
	end := len(rs) - 2
	res.WriteString(string(rs[0:end]))
	winnerName := ""
	winnerValue := 0
	for k := range mulInfo.userInfoMap {
		res.WriteString("\n")
		res.WriteString(k)
		res.WriteString("的最终得分为")
		nowValue := mulInfo.userInfoMap[k]
		if winnerValue < nowValue {
			winnerName = k
			winnerValue = nowValue
		}
		res.WriteString(strconv.Itoa(mulInfo.userInfoMap[k]))
	}
	if winnerName != "" {
		s := "\n恭喜" + winnerName + "夺得本次接龙第一名"
		res.WriteString(s)
	}

	mulInfo = NewMulInfo()
	multipleIsOpen = false
	return res.String()
}

func GenerateRandomIdiom() string {
	j := 0
	keys := make([]string, len(idiomsMap))
	for k := range idiomsMap {
		keys[j] = k
		j++
	}
	rand.Seed(time.Now().UnixNano())
	return idiomsMap[keys[rand.Intn(len(idiomsMap))]].idiom
}

func isIdiom(chinese string) bool {
	if len(chinese) == 12 && allFourLetterWordsSet.Contains(chinese) {
		return true
	}
	return false
}

func ContinueConcatenateDragon(word string, userInfo *userInfo) result {
	lastWordPinYin, err := pinyin.To_Py(userInfo.lastWord, " ", "")
	if err != nil {
		fmt.Println(err)
	}
	wordPinYin, err := pinyin.To_Py(word, " ", "")
	if err != nil {
		fmt.Println(err)
	}
	if !isIdiom(word) {
		return result{
			content: "",
			error:   "您输入的不是成语",
		}
	}
	if !CheckIsMatch(wordPinYin, lastWordPinYin) {
		return result{
			content: "",
			error:   "你接错了，请再来一次。",
		}
	}
	userInfo.path += word
	userInfo.path += "->"
	userInfo.value++
	return result{
		content: GetNextMatchIdiom(wordPinYin, userInfo),
		error:   "",
	}

}

func MulContinueConcatenateDragon(word string, userName string) result {
	lastWordPinYin, err := pinyin.To_Py(mulInfo.lastWord, " ", "")
	if err != nil {
		fmt.Println(err)
	}
	wordPinYin, err := pinyin.To_Py(word, " ", "")
	if err != nil {
		fmt.Println(err)
	}
	if !isIdiom(word) {
		return result{
			content: "",
			error:   "您输入的不是成语",
		}
	}
	if !CheckIsMatch(wordPinYin, lastWordPinYin) {
		return result{
			content: "",
			error:   "你接错了，请再来一次。",
		}
	}
	mulInfo.path += word
	mulInfo.path += "->"
	mulInfo.userInfoMap[userName]++
	return result{
		content: MulGetNextMatchIdiom(wordPinYin),
		error:   "",
	}

}

func CheckIsMatch(inputPinYin string, lastWordPinYin string) bool {
	inputPinYinArray := strings.Split(inputPinYin, " ")
	lastWordArray := strings.Split(lastWordPinYin, " ")
	if inputPinYinArray[0] != lastWordArray[3] {
		return false
	}
	return true
}

func GetNextMatchIdiom(word string, userInfo *userInfo) string {
	wordArray := strings.Split(word, " ")
	wordLastPinYin := wordArray[3]
	for k := range idiomsMap {
		if idiomsMap[k].firstPinYin == wordLastPinYin {
			userInfo.path += idiomsMap[k].idiom
			userInfo.path += "->"
			return idiomsMap[k].idiom
		}
	}
	return ""
}

func MulGetNextMatchIdiom(word string) string {
	wordArray := strings.Split(word, " ")
	wordLastPinYin := wordArray[3]
	for k := range idiomsMap {
		if idiomsMap[k].firstPinYin == wordLastPinYin {
			mulInfo.path += idiomsMap[k].idiom
			mulInfo.path += "->"
			return idiomsMap[k].idiom
		}
	}
	return ""
}

func SaveUserData(info userInfo) {
	userInfoMap[info.userId] = info
}
