package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/TskFok/OpenAi/app/global"
	"github.com/TskFok/OpenAi/app/process"
	"github.com/TskFok/OpenAi/bootstrap"
	"github.com/TskFok/OpenAi/router"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed public/*
var Fs embed.FS

func main() {
	//守护进程
	args := os.Args

	if len(args) != 1 && args[1] == "bg" {
		process.InitProcess()
	}

	//加载配置
	bootstrap.Init()

	//加载router
	router.InitRouter()

	//静态资源加载
	router.Handle.Any("public/*filepath", func(context *gin.Context) {
		staticServer := http.FileServer(http.FS(Fs))
		staticServer.ServeHTTP(context.Writer, context.Request)
	})

	//浏览器图标
	router.Handle.StaticFileFS("/favicon.ico", "./public/static/favicon.ico", http.FS(Fs))

	addr := fmt.Sprintf(":%d", 443)
	if global.AppMode == gin.DebugMode {
		addr = fmt.Sprintf(":%d", 9988)
	}

	s := &http.Server{
		Addr:           addr,
		Handler:        router.Handle,
		ReadTimeout:    time.Duration(20) * time.Second,
		WriteTimeout:   time.Duration(20) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if global.AppMode == gin.DebugMode {
			//不使用https
			if err := s.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
				log.Printf("listen: %s\n", err)
			}
		} else {
			//使用https
			if err := s.ListenAndServeTLS(global.TlsCert, global.TlsKey); err != nil && errors.Is(err, http.ErrServerClosed) {
				log.Printf("listen: %s\n", err)
			}
		}
	}()

	//接收信号关闭
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
