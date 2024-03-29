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
		log.Println("in httpServer routine ,before ListenAndServe")
		err := httpServer.ListenAndServe()
		log.Println("httpserver routine ,over:", err)
		return err
	})

	// 监听 signal事件
	g.Go(func() error {
		log.Println("in signal.Notify routine ")
		c := make(chan os.Signal, 1)
		signal.Notify(c)
		defer func() {
			signal.Stop(c)
		}()

		for stop := false; !stop; {
			select {
			case s := <-c:
				if s != nil {
					log.Println("one signal:", s.String())
				}
				if os.Interrupt == s {
					// onInterrupt() // 响应退出 (为了解耦提到外面)
					return errors.New("Interrupt")
				}
			case <-ctx.Done():
				stop = true
			}
		}
		log.Println("in signal.Notify,Other Error:", ctx.Err())
		return nil
	})

	// 由于http server 里面直接封装了，没法接入其他的一个goroutine监听group的异常
	onInterrupt := func() { // 响应其他goroutine failed
		err := httpServer.Shutdown(context.Background())
		log.Println("httpServer Shutdown return:", err)
	}
	g.Go(func() error {
		log.Println("onInterrupt start")
		<-ctx.Done()
		onInterrupt()
		log.Println("onInterrupt over:", ctx.Err())
		return ctx.Err()
	})

	// 等待启动的routine全部结束
	err := g.Wait()
	log.Println(">>Exit with ", err)
}
