// Package bluearchive 插件主体部分
package bluearchive

import (
	"github.com/KomeiDiSanXian/BlueArchive/bluearchive/utils"
	"github.com/KomeiDiSanXian/BlueArchive/bluearchive/wiki"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var engine = control.Register("碧蓝档案", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault: false,
	Brief:            "ba相关信息查询",
	Help: "bluearchive\n" +
		"- .ba活动\t查询活动信息" +
		"- .ba公告\t查看公告",
	PrivateDataFolder: "bluearchive",
})

func init() {
	// 完全匹配触发
	// 使用合并消息转发
	engine.OnFullMatch(".ba活动").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			w := wiki.NewWikiData()
			if err := w.Request(); err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERROR: 请求错误"))
				return
			}
			layout := "1月02日 15:04"
			msg := make(message.Message, 0, len(*w.Events))
			eventmsg := w.Events.PrintEvent(layout)
			for _, sendmsg := range eventmsg {
				msg = append(msg, ctxext.FakeSenderForwardNode(ctx, utils.Txt2Img(ctx, sendmsg)))
			}
			if id := ctx.Send(msg).ID(); id == 0 {
				ctx.SendChain(message.Text("ERROR: 可能被风控了"))
			}
		})
	engine.OnFullMatch(".ba公告").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			w := wiki.NewWikiData()
			if err := w.Request(); err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERROR: 请求错误"))
				return
			}
			msg := make(message.Message, 0, len(*w.Announcements))
			announce := w.Announcements.PrintAnnouncements()
			for _, sendmsg := range announce {
				msg = append(msg, ctxext.FakeSenderForwardNode(ctx, utils.Txt2Img(ctx, sendmsg)))
			}
			if id := ctx.Send(msg).ID(); id == 0 {
				ctx.SendChain(message.Text("ERROR: 可能被风控了"))
			}
		})
}
