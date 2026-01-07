package main

import (
	"log"
	"sync"
	"time"

	"github.com/rei0721/ghhook"
)

// Debouncer 防抖器
type Debouncer struct {
	mu    sync.Mutex
	timer *time.Timer
	wait  time.Duration
}

func NewDebouncer(wait time.Duration) *Debouncer {
	return &Debouncer{wait: wait}
}

func (d *Debouncer) Do(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.wait, f)
}

func main() {
	// 创建新的监听器
	hook, err := ghhook.New(":9901", "qwq")
	if err != nil {
		log.Fatal(err)
	}

	// 创建防抖器，防止短时间内多次触发构建 (例如 2 秒内的重复推送)
	pushDebouncer := NewDebouncer(2 * time.Second)

	// 注册 push 事件钩子
	hook.On("push", func(ctx *ghhook.Context) error {
		// 捕获当前的上下文信息
		repoName := ctx.Repo.FullName
		message := ctx.Push.HeadCommit.Message

		log.Printf("📥 收到推送信号 -> %s", repoName)

		// 防抖处理
		pushDebouncer.Do(func() {
			log.Printf("📦 [Debounced] 开始处理仓库 %s 的推送: %s", repoName, message)
			// 这里可以添加实际的业务逻辑，例如调用部署脚本等

			// 执行 Shell 命令
			result, err := ctx.Exec("./run.sh", ctx.Repo.Name)
			if err != nil {
				log.Printf("部署失败: %v", err)
				return
			}
			log.Printf("部署输出: %s", result.Stdout)

			// 执行 HTTP 请求
			resp, err := ctx.HTTP("POST", "http://38.14.250.76:9901/notify",
				ghhook.WithHeaders(map[string]string{
					"Content-Type": "application/json",
				}),
				ghhook.WithBody(`{"event": "push", "repo": "`+ctx.Repo.FullName+`"}`),
			)
			if err != nil {
				log.Printf("通知失败: %v", err)
				return
			}
			log.Printf("通知已发送: %d", resp.StatusCode)
		})

		return nil
	})

	// 注册 issue 事件钩子
	hook.On("issues", func(ctx *ghhook.Context) error {
		log.Printf("📝 Rei 新 Issue: %s", ctx.Issue.Title)
		return nil
	})

	// 启动监听器
	hook.Run()
}
