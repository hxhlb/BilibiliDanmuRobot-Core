package handler

import (
	"github.com/Akegarasu/blivedm-go/message"
	"github.com/xbclub/BilibiliDanmuRobot-Core/utiles/danmuProcess"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
)

// 弹幕处理
func (w *wsHandler) receiveDanmu() {
	//弹幕处理的功能类接口
	danmuProcessFuncList := createDanmuProcessFuncList(w)
	w.client.OnDanmaku(func(danmaku *message.Danmaku) {
		sender := *danmaku.Sender
		strContent := danmaku.Content
		if sender.Uid != w.userId {
			//移除表情包内容，[]形式
			re := regexp.MustCompile("\\[(.*?)\\]")
			strContent = re.ReplaceAllString(strContent, "")
			if len(strContent) > 0 {
				for _, danmuProcessFunc := range danmuProcessFuncList {
					danmuProcessFunc.SetDanmu(strContent, sender)
					danmuProcessFunc.DoDanmuProcess()
				}
			}
		}
		// 实时输出弹幕消息
		logx.Infof("%d %s:%s", danmaku.Sender.Uid, danmaku.Sender.Uname, danmaku.Content)
	})
}

func createDanmuProcessFuncList(w *wsHandler) []danmuProcess.DanmuProcess {
	//弹幕处理的功能类的集合
	danmuProcessClass := new(danmuProcess.DanmuProcessClass)
	var danmuProcessFuncList []danmuProcess.DanmuProcess
	//判断启用机器人聊天
	//if w.svc.Config.FuzzyMatchCmd {
	gptClass := danmuProcessClass.GptClass.Create()
	gptClass.SetConfig(w.svc)
	danmuProcessFuncList = append(danmuProcessFuncList, gptClass)
	//}
	//判断启用机器人抽签
	if w.svc.Config.DrawByLot {
		drawByLotClass := danmuProcessClass.DrawByLotClass.Create()
		drawByLotClass.SetConfig(w.svc)
		danmuProcessFuncList = append(danmuProcessFuncList, drawByLotClass)
	}
	// 签到
	if w.svc.Config.SignInEnable {
		signInClass := danmuProcessClass.SignInClass.Create()
		signInClass.SetConfig(w.svc)
		danmuProcessFuncList = append(danmuProcessFuncList, signInClass)
	}
	return danmuProcessFuncList
}
