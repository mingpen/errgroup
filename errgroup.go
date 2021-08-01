package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

// 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

func main() {
	httpServer := &http.Server{Addr: ":" + os.Getenv("PORT"),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})}
	g, ctx := errgroup.WithContext(context.Background())

	// 启动http服务
	g.Go(func() error {
		err := httpServer.ListenAndServe()
		log.Println("httpserver routine ,over:", err)
		return err
	})

	// 监听 signal事件
	onInterrupt := func() {
		err := httpServer.Shutdown(ctx)
		log.Println("httpServer Shutdown return:", err)
	}
	//
	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c)
		select {
		case s := <-c:
			if os.Interrupt == s {
				onInterrupt() // 响应
				return errors.New("Interrupt")
			}
			log.Println(s.String())
		case <-ctx.Done():
			break
		}
		log.Println("in signal.Notify,Other Error:", ctx.Err())
		return nil
	})
	err := g.Wait()
	log.Println("Exit with ", err)
}
