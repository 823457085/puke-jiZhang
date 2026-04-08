<template>
  <view id="app">
    <slot />
  </view>
</template>

<script>
export default {
  onLaunch() {
    // 检查登录状态
    const userId = uni.getStorageSync('userId')
    if (!userId) {
      // 未登录，跳转到登录页
      uni.reLaunch({ url: '/pages/index/login' })
    }
  },
  onShow() {},
  onHide() {}
}
</script>

<style>
/* 全局样式 */
page {
  background-color: #f8fafc;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  font-size: 15px;
  color: #1e293b;
}

view {
  box-sizing: border-box;
}

/* 金额颜色 */
.text-green { color: #16a34a; }
.text-red { color: #dc2626; }
.text-blue { color: #2563eb; }
.text-gray { color: #64748b; }

/* 金额数字 */
.amount-positive { color: #16a34a; font-size: 32rpx; font-weight: 700; }
.amount-negative { color: #dc2626; font-size: 32rpx; font-weight: 700; }

/* 卡片 */
.card {
  background: #ffffff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.06);
}

/* 主按钮 */
.btn-primary {
  background: #2563eb;
  color: #ffffff;
  border-radius: 8rpx;
  height: 96rpx;
  line-height: 96rpx;
  text-align: center;
  font-size: 32rpx;
  font-weight: 500;
  border: none;
}

.btn-primary:active {
  background: #1d4ed8;
}

.btn-primary[disabled] {
  background: #94a3b8;
}

/* 次按钮 */
.btn-secondary {
  background: #f1f5f9;
  color: #1e293b;
  border-radius: 8rpx;
  height: 96rpx;
  line-height: 96rpx;
  text-align: center;
  font-size: 32rpx;
  font-weight: 500;
}

/* 输入框 */
.input {
  background: #f8fafc;
  border: 2rpx solid #e2e8f0;
  border-radius: 8rpx;
  height: 96rpx;
  padding: 0 24rpx;
  font-size: 32rpx;
}

.input:focus {
  border-color: #2563eb;
}
</style>
