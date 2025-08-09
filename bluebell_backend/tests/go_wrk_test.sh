#!/bin/bash

# Bluebell论坛go-wrk压测脚本
# 使用方法: ./go_wrk_test.sh [场景]

# 设置变量
BASE_URL="http://localhost:8080"
CONCURRENT=100
DURATION=30s
TIMEOUT=5s

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Bluebell论坛go-wrk压测开始${NC}"

# 检查go-wrk是否安装
if ! command -v go-wrk &> /dev/null; then
    echo -e "${YELLOW}安装go-wrk...${NC}"
    go install github.com/tsliwowicz/go-wrk@latest
fi

# 压测场景函数
test_post_list() {
    echo -e "${GREEN}📋 测试: 帖子列表查询${NC}"
    go-wrk -c $CONCURRENT -d $DURATION -T $TIMEOUT \
        -H "Accept: application/json" \
        "$BASE_URL/api/v1/posts?page=1&size=10"
}

test_post_detail() {
    echo -e "${GREEN}📖 测试: 帖子详情查看${NC}"
    # 假设有帖子ID 1-1000
    go-wrk -c $CONCURRENT -d $DURATION -T $TIMEOUT \
        -H "Accept: application/json" \
        "$BASE_URL/api/v1/post/1"
}

test_hot_posts() {
    echo -e "${GREEN}🔥 测试: 热门帖子排行${NC}"
    go-wrk -c $CONCURRENT -d $DURATION -T $TIMEOUT \
        -H "Accept: application/json" \
        "$BASE_URL/api/v1/posts/hot?days=7"
}

test_community_posts() {
    echo -e "${GREEN}🏘️ 测试: 社区帖子列表${NC}"
    # 假设社区ID为1
    go-wrk -c $CONCURRENT -d $DURATION -T $TIMEOUT \
        -H "Accept: application/json" \
        "$BASE_URL/api/v1/community/1/posts?page=1&size=10"
}

test_vote_post() {
    echo -e "${GREEN}👍 测试: 投票操作${NC}"
    # 创建测试数据
    curl -s -X POST "$BASE_URL/api/v1/vote" \
        -H "Content-Type: application/json" \
        -d '{"post_id": 1, "direction": 1}' \
        > /dev/null 2>&1 || true
    
    go-wrk -c 50 -d $DURATION -T $TIMEOUT \
        -M POST \
        -H "Content-Type: application/json" \
        -B '{"post_id": 1, "direction": 1}' \
        "$BASE_URL/api/v1/vote"
}

test_create_post() {
    echo -e "${GREEN}✍️ 测试: 创建帖子${NC}"
    go-wrk -c 20 -d $DURATION -T $TIMEOUT \
        -M POST \
        -H "Content-Type: application/json" \
        -B '{"title": "压测帖子", "content": "这是压测内容", "community_id": 1}' \
        "$BASE_URL/api/v1/post"
}

# 完整压测流程
run_full_test() {
    echo -e "${GREEN}🎯 开始完整压测流程${NC}"
    
    echo -e "${YELLOW}1. 预热缓存...${NC}"
    curl -s "$BASE_URL/api/v1/posts?page=1&size=10" > /dev/null 2>&1
    
    echo -e "${YELLOW}2. 帖子列表压测${NC}"
    test_post_list
    
    echo -e "${YELLOW}3. 帖子详情压测${NC}"
    test_post_detail
    
    echo -e "${YELLOW}4. 热门排行压测${NC}"
    test_hot_posts
    
    echo -e "${YELLOW}5. 社区帖子压测${NC}"
    test_community_posts
    
    echo -e "${YELLOW}6. 投票操作压测${NC}"
    test_vote_post
    
    echo -e "${YELLOW}7. 发帖操作压测${NC}"
    test_create_post
    
    echo -e "${GREEN}✅ 完整压测完成${NC}"
}

# 渐进式压测
run_progressive_test() {
    echo -e "${GREEN}📈 渐进式压测（逐步增加并发）${NC}"
    
    for concurrent in 10 50 100 200 500 1000; do
        echo -e "${YELLOW}并发数: $concurrent${NC}"
        CONCURRENT=$concurrent
        test_post_list
        echo ""
        sleep 2
    done
}

# 性能基准测试
run_baseline_test() {
    echo -e "${GREEN}📊 性能基准测试${NC}"
    
    # 单机性能测试
    echo -e "${YELLOW}单机100并发测试${NC}"
    CONCURRENT=100 DURATION=60s
    test_post_list
    
    # 缓存命中率测试
    echo -e "${YELLOW}缓存命中率测试${NC}"
    for i in {1..100}; do
        curl -s "$BASE_URL/api/v1/posts?page=1&size=10" > /dev/null 2>&1
    done
    
    # 查看Redis命中率
    redis.Client.Info("stats").Result()
}

# 主程序
case "${1:-full}" in
    "list")
        test_post_list
        ;;
    "detail")
        test_post_detail
        ;;
    "hot")
        test_hot_posts
        ;;
    "community")
        test_community_posts
        ;;
    "vote")
        test_vote_post
        ;;
    "post")
        test_create_post
        ;;
    "progressive")
        run_progressive_test
        ;;
    "baseline")
        run_baseline_test
        ;;
    "full"|*)
        run_full_test
        ;;
esac

echo -e "${GREEN}🎉 压测脚本执行完成！${NC}"
echo -e "${YELLOW}查看完整报告请查看: tests/load_test_report.txt${NC}"