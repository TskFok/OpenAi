# 钉钉机器人回调接口

使用gpt3.5接口

``````
/chat 提供给钉钉机器人回调地址的接口
/chat-web 流式请求页面 只在firefox浏览器能使用流式传输,废弃
/chat-web-sse 新的流式请求页面 使用sse
/chat-web-ws 新的流式请求页面 使用websocket

非守护进程运行 go run main.go
守护进程运行 go run main.go bg
日志在chat.log
关闭进程 kill -2 pid
``````
