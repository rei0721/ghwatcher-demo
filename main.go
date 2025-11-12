package main

import (
	"log"

	"github.com/kawaiirei0/ghwatcher"
)

func main() {
	// 创建新的监听器
	w, err := ghwatcher.New(":6001", "qwq")
	if err != nil {
		log.Fatal(err)
	}

	// 注册 push 事件钩子
	w.On("push", func(ctx *ghwatcher.Context) error {
		log.Printf("📦 Rei 仓库 %s 收到推送: %s",
			ctx.Repo.FullName,
			ctx.Push.HeadCommit.Message)
		return nil
	})

	// 注册 issue 事件钩子
	w.On("issues", func(ctx *ghwatcher.Context) error {
		log.Printf("📝 Rei 新 Issue: %s", ctx.Issue.Title)
		return nil
	})

	// 启动监听器
	w.Run()
}
