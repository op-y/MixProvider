# MixProvider
MixProvider 是按照Open Falcon监控系统短信报警接口规范开发的Provider程序。

## 主要功能：
* HTTP Server 用于接收请求
* 转发告警邮件到其他HTTP API
* 发送告警短信到其他微信HTTP API
* 发送告警短信到微信接口
* 发送告警短信到容联云通讯接口
* 发送告警短信到搜狐sendcloud接口

## 文件说明:
* cfg.json      -- 配置文件
* control       -- 控制脚本
* filter.go     -- 短信过滤代码
* provider.go   -- 入口HTTP Server  代码
* relayer.go    -- 转发其他HTTP API 代码
* sendcloud.go  -- 搜狐sendcloud接口发送短信代码
* wechat.go     -- 微信接口发送短信代码
* yuntongxun.go -- 容联云通讯接口发送短信代码

## 使用方法
1. git clone https://github.com/op-y/MixProvider.git
2. cd MixProvider
3. go build -o mix-provider
4. 根据实际情况修改cfg.json配置
5. ./control start

## 修改历史
* 2017/9/27 根据业务需要，将将短信和微信接口分离，去掉了邮件接口
