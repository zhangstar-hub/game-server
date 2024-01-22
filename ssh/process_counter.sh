#!/bin/bash
counter_file="pids/pid_count.txt"

# 检查计数器文件是否存在，不存在则创建
if [ ! -e "$counter_file" ]; then
    echo "1" > "$counter_file"
fi

exec 200<> "$counter_file"

# 获取文件锁
if flock -x 200; then
    # 读取当前计数
    current_count=$(cat "$counter_file")

    # 更新计数
    ((current_count++))

    # 如果超过1000，重置为1
    if [ "$current_count" -gt 1000 ]; then
        current_count=1
    fi

    # 将新计数写入文件
    echo "$current_count" > "$counter_file"

    # 释放文件锁
    flock -u 200
    echo $current_count
else
    echo 0
fi

# 关闭文件描述符
exec 200>&-
