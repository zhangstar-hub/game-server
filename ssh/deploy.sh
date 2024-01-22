#!/bin/bash
cd ..

cmd=$1
rebuild=$2

echo $cmd $rebuild

port=8080
pid_file=pids/app.pid

if [ ! -d "logs" ]; then
  mkdir -p logs
fi

if [ ! -d "pids" ]; then
  mkdir -p pids
fi

start() {
  if nc -z -w 1 127.0.0.1 "$port"; then
      echo "进程已经启动了"
      exit 0
  fi
  returned_count=$(ssh/process_counter.sh)
  echo $returned_count
  if [ ! -f "main" ] || [ "$rebuild" == "rebuild" ]; then
    go build main.go
  fi
  nohup ./main -pid=$returned_count > logs/app.log 2>&1 &
  pid=$!
  echo $pid > $pid_file
  echo $pid > "pids/app_$returned_count.pid"
  echo "项目启动: 进程ID为 $(cat $pid_file)"
}

stop() {
  if [ ! -f "$pid_file" ]; then
    echo "进程文件app.pid不存在, 请手动关闭进程"
    exit 0
  fi
  kill $(cat $pid_file)
  rm $pid_file
}

shutdown() {
  if [ ! -f "$pid_file" ]; then
    echo "进程文件app.pid不存在, 请手动关闭进程"
    exit 0
  fi
  kill $(cat $pid_file)
  kill $(cat $pid_file)
  rm $pid_file
}

restart() {
  if nc -z -w 1 127.0.0.1 "$port"; then
    stop
  fi
  start
}

case "$cmd" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    shutdown)
        shutdown
        ;;
    restart)
        restart
        ;;
    *)
        echo "Usage: $0 {start|stop|shutdown|restart} [rebuild]"
        echo "start 启动服务"
        echo "stop 停止服务, 但不关闭旧连接"
        echo "shutdown 停止服务, 关闭所有连接"
        echo "restart 启动新服务, stop旧服务"
        echo "* rebuild 如果存在则重新构建项目"
        exit 1
esac

exit 0