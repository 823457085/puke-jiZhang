<template>
  <view class="page">
    <view class="user-card card">
      <view class="user-avatar">{{ userInfo.nickname[0] }}</view>
      <view class="user-name">{{ userInfo.nickname }}</view>
    </view>

    <view class="menu-list">
      <view class="menu-item card">
        <text class="menu-text">清除登录记录</text>
        <text class="menu-arrow">›</text>
      </view>
    </view>

    <view class="app-info">
      <view class="app-name">记牌官</view>
      <view class="app-version">版本 1.0.0</view>
      <view class="app-tag">无广告 · 极简记账</view>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      userInfo: {}
    }
  },
  onShow() {
    const info = uni.getStorageSync('userInfo') || {}
    this.userInfo = info
  },
  methods: {
    logout() {
      uni.showModal({
        title: '确认退出',
        content: '确定要清除登录记录吗？',
        success: (res) => {
          if (res.confirm) {
            uni.clearStorageSync()
            uni.reLaunch({ url: '/pages/index/login' })
          }
        }
      })
    }
  }
}
</script>

<style scoped>
.page {
  padding: 24rpx 32rpx;
  min-height: 100vh;
  background: #f8fafc;
}

.user-card {
  display: flex;
  align-items: center;
  padding: 32rpx;
}

.user-avatar {
  width: 100rpx;
  height: 100rpx;
  background: #2563eb;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48rpx;
  font-weight: 600;
  margin-right: 24rpx;
}

.user-name {
  font-size: 36rpx;
  font-weight: 600;
  color: #1e293b;
}

.menu-list {
  margin-bottom: 40rpx;
}

.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 32rpx 24rpx;
}

.menu-text {
  font-size: 30rpx;
  color: #1e293b;
}

.menu-arrow {
  font-size: 36rpx;
  color: #cbd5e1;
}

.app-info {
  text-align: center;
  padding: 60rpx 0;
}

.app-name {
  font-size: 36rpx;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 8rpx;
}

.app-version {
  font-size: 26rpx;
  color: #94a3b8;
  margin-bottom: 12rpx;
}

.app-tag {
  font-size: 24rpx;
  color: #2563eb;
  background: #eff6ff;
  padding: 6rpx 20rpx;
  border-radius: 20rpx;
  display: inline-block;
}
</style>
