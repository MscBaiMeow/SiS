package whitelist

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/miaoscraft/SiS/data"
	"regexp"
	"strconv"
)

var expQQ = regexp.MustCompile(`^(?:([0-9]+)|\[CQ:at,qq=([0-9]+)])$`) // 匹配一个QQ或At
var expName = regexp.MustCompile(`^([0-9A-Za-z_]{3,16})$`)            // 匹配一个玩家名

func Info(args []string, fromQQ int64, ret func(string)) bool {
	// 找出当前想查询的人的QQ
	switch len(args) {
	case 1:
		qqInfo(fromQQ, ret)
		return true
	case 2:
		if sms := expQQ.FindStringSubmatch(args[1]); len(sms) == 3 { // 匹配一个QQ或At
			for _, sm := range sms[1:3] { // [3]sms中后两项有一项为空，另一项为QQ
				qq, err := strconv.ParseInt(sm, 10, 64)
				if err != nil {
					continue
				}
				qqInfo(qq, ret)
				return true
			}
		} else if sm := expName.FindStringSubmatch(args[1]); len(sm) == 2 { // 匹配一个玩家名
			nameInfo(sm[1], ret)
			return true
		}
		return false
	default:
		return false
	}
}

func qqInfo(targetQQ int64, ret func(string)) {
	// 查询本人的绑定
	ID, err := data.GetWhitelistByQQ(targetQQ)
	if err != nil {
		Logger.Errorf("读取玩家绑定的ID出错: %v", err)
		ret("数据库查询失败惹(つД`)ノ")
		return
	}
	if ID == uuid.Nil {
		ret("这个还没有绑定白名单呢")
		return
	}

	// 根据UUID找到名字
	name, err := getName(ID)
	if err != nil {
		ret("游戏名查询失败惹(つД`)ノ")
		return
	}
	ret(name)
}

func nameInfo(targetName string, ret func(string)) {
	name, id, err := GetUUID(targetName)
	if err != nil {
		Logger.Errorf("查询UUID失败: %v", err)
		ret("查无此人")
		return
	}

	qq, err := data.GetWhitelistByUUID(id)
	if err != nil {
		Logger.Errorf("数据库查询QQ失败: %v", err)
		ret("数据库出问题了(つД`)ノ")
		return
	}

	if qq == 0 {
		ret(fmt.Sprintf("没人绑定%s哟~", name))
	} else {
		ret(fmt.Sprintf("啊呐占用%s的是%d哟", name, qq))
	}
}
