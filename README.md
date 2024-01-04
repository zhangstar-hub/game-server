# Socket游戏服务器

## 基础功能

1. 平滑重启
2. TCP长连接，以处理粘包问题
3. 心跳检测，长时间未发请求自动断开连接
4. 请求限流 令牌桶+滑动窗口
5. 中间件
   1. 已定义login中间件，拒绝未登录请求
6. Mysql自动建表
7. 路由映射
8. 玩家上下文管理