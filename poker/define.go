package poker

import (
	"strings"
)

//一些共用的扑克定义抽象

/*card就是普通的int32，从0到53
定义点数(num)从A-K对应整型变量0-12
定义花色(color)为0=方块，1=梅花，2=红桃，3=黑桃
则：card=color*13+num，color=card/13, num=card%13
定义牌值(rank)为牌的实际等级，在斗地主中2>A>K，rank算法见下面代码
*/

const (
	BlackJoker = 52
	RedJoker   = 53
	ColorCards = 13
)

var NumberString = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "X", "J", "Q", "K"}
var ColorString = []string{"♠️", "♥️", "♣️", "♦️"}

func GetCard(num, color int32) int32 {
	return color*ColorCards + num
}

func GetNum(card int32) int32 {
	if card >= BlackJoker {
		return -1
	}
	return card % ColorCards
}

func GetColor(card int32) int32 {
	if card >= BlackJoker {
		return -1
	}
	return card / ColorCards
}

func GetPrintColor(card int32) string {
	color := GetColor(card)
	if color >= 0 {
		return ColorString[color]
	}
	return "-"
}

func GetPrintNum(card int32) string {
	if card == BlackJoker {
		return "Y"
	}
	if card == RedJoker {
		return "Z"
	}
	num := card % ColorCards
	return NumberString[num]
}

func Cards2String(cards []int32) string {
	var sb strings.Builder
	for _, card := range cards {
		sb.WriteString(GetPrintNum(card))
		sb.WriteString(" ")
	}
	return sb.String()
}

func Cards2ColorString(cards []int32) string {
	var sb strings.Builder
	for _, card := range cards {
		if color := GetPrintColor(card); color != "-" {
			sb.WriteString(GetPrintColor(card))
		}
		sb.WriteString(GetPrintNum(card))
		sb.WriteString(" ")
	}
	return sb.String()
}
