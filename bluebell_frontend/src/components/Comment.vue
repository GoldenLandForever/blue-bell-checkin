<template>
  <div class="comment-container">
    <h3 class="comment-title">评论区</h3>
    
    <!-- 调试信息 -->
    <div class="debug-info" v-if="showDebug">
      <p>登录状态: {{ isLogin }}</p>
      <p>sourceId: {{ sourceId }}</p>
      <p>评论数量: {{ comments.length }}</p>
    </div>

    <!-- 评论输入框 -->
    <div class="comment-input-section" v-if="isLogin">
      <textarea
        v-model="commentContent"
        placeholder="写下你的评论..."
        class="comment-textarea"
      ></textarea>
      <button @click="submitComment" class="submit-btn">提交评论</button>
    </div>
    <div class="not-login" v-else>
      <p>请先登录后再评论</p>
    </div>

    <!-- 评论列表 -->
    <div class="comment-list">
      <div v-for="comment in comments" :key="comment.commentID" class="comment-item">
        <div class="comment-header">
          <span class="comment-author">{{ comment.author_id || '匿名用户' }}</span>
          <span class="comment-time">{{ formatDate(comment.create_time) }}</span>
        </div>
        <div class="comment-content">{{ comment.content }}</div>
      </div>
      <div v-if="comments.length === 0" class="no-comment">
        <p>暂无评论</p>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  name: 'Comment',
  props: {
    sourceId: {
      type: String,
      required: true,
    },
    // 允许父组件直接传递登录状态
    loginStatus: {
      type: Boolean,
      required: false,
      default: null
    },
  },
  data() {
    return {
      comments: [],
      commentContent: '',
      showDebug: true, // 调试模式开启
    };
  },
  computed: {
      isLogin() {
        try {
          // 优先使用父组件传递的登录状态
          if (this.loginStatus !== null) {
            return this.loginStatus;
          }
          
          // 使用Vuex getter获取登录状态(推荐方式)
          if (this.$store && this.$store.getters) {
            return this.$store.getters.isLogin;
          }
          
          // 备选方案: 直接访问state中的isLogin字段
          if (this.$store && this.$store.state) {
            return this.$store.state.isLogin;
          }
          
          console.error('无法获取登录状态');
          return false;
        } catch (error) {
          console.error('获取登录状态失败:', error);
          return false;
        }
      }
    },
  mounted() {
    // console.log('评论组件挂载成功');
    // console.log('登录状态:', this.isLogin);
    // console.log('sourceId:', this.sourceId);
    this.getCommentList();
  },
  methods: {
    // 获取评论列表
    getCommentList() {
      this.$axios({
        method: 'get',
        url: '/comment',
        params: {
          post_id: this.sourceId, // 确保是字符串类型
          page: 1,  // 页码，默认为1
          size: 10, // 每页数量，默认为10
          order: 'DESC' // 排序方式，降序
        }
      })
        .then(response => {
          console.log('获取评论列表成功:', response.data);
          // 增加对response和response.data的检查
          if (response.code === 1000) {
            this.comments = response.data;
          } else {
            console.error('获取评论列表失败:', response.data.msg || '未知错误');
            this.comments = [];
          }
        })
        .catch(error => {
          console.error('获取评论列表失败:', error);
          alert('获取评论列表失败，请稍后再试');
          this.comments = [];
        });
    },

    // 提交评论
    submitComment() {
      if (!this.commentContent.trim()) {
        alert('评论内容不能为空');
        return;
      }

      this.$axios({
        method: 'post',
        url: '/comment',
        data: {
          post_id: String(this.sourceId), // 确保是字符串类型
          content: this.commentContent
        }
      })
        .then(response => {
          console.log(response.data);
          // 增加对response和response.data的检查
          if (response.code == 1000) {
            this.commentContent = '';
            this.getCommentList(); // 重新获取评论列表
          } else {
            console.log(response.data.msg);
            alert('评论提交失败: ' + (response.data.msg || '未知错误'));
          }
        })
        .catch(error => {
          console.error('提交评论失败:', error);
          alert('评论提交失败，请稍后再试');
        });
    },

    // 格式化日期
    formatDate(dateStr) {
      const date = new Date(dateStr);
      return date.toLocaleString();
    },
  },
};
</script>

<style scoped>
.comment-container {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.comment-title {
  font-size: 18px;
  margin-bottom: 20px;
  color: #333;
}

.comment-input-section {
  margin-bottom: 30px;
}

.comment-textarea {
  width: 100%;
  height: 120px;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  resize: none;
  font-size: 14px;
  margin-bottom: 10px;
}

.submit-btn {
  background-color: #409eff;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  float: right;
}

.submit-btn:hover {
  background-color: #3a8ee6;
}

.not-login {
  padding: 20px;
  text-align: center;
  color: #999;
  background-color: #f9f9f9;
  border-radius: 4px;
  margin-bottom: 30px;
}

.comment-list {
  clear: both;
}

.comment-item {
  padding: 16px 0;
  border-bottom: 1px solid #eee;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.comment-author {
  font-weight: bold;
  color: #333;
}

.comment-time {
  color: #999;
  font-size: 12px;
}

.comment-content {
  color: #333;
  line-height: 1.6;
}

.debug-info {
  background-color: #f0f8ff;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 15px;
  font-size: 12px;
  color: #666;
}

.no-comment {
  padding: 30px 0;
  text-align: center;
  color: #999;
}
</style>