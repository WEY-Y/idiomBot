## 1.业务逻辑

下面是成语接龙机器人的游戏模式流程。

![image](https://github.com/WEY-Y/idiomBot/blob/main/idiomBot/pic/logic1.png)

下面是成语接龙机器人的接龙判定逻辑。

![image](https://github.com/WEY-Y/idiomBot/blob/main/idiomBot/pic/logic2.png)

## 2.数据结构

```
成语字典map  key是拼音 value是一个结构体（结构体内存放了 成语的汉字 首位汉字的拼音）
var idiomsMap = make(map[string]IdiomAndFirstPinYin)
所有常用四字词语 第一遍过筛使用
var allFourLetterWordsSet = mapset.NewSet()
游戏中的用户信息
var userInfoMap = make(map[string]userInfo)
```



## 3.函数功能

见函数名 较为规范易懂 如：CheckIsMatch()、GenerateRandomIdiom()。



## 4.文件功能

main.go 主要包含主体连接逻辑

process.go 主要进行关键词匹配对应功能点

idioms.go 主要是跟成语接龙相关的逻辑实现

