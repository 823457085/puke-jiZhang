<template>
  <view class="page">
    <view class="form-item">
      <view class="form-label">房间码</view>
      <input
        class="input code-input"
        v-model="roomCode"
        placeholder="请输入6位房间码"
        placeholder-class="input-placeholder"
        maxlength="6"
        type="number"
        focus
      />
    </view>

    <button
      class="btn-primary"
      @click="handleJoin"
      :loading="loading"
      :disabled="roomCode.length !== 6"
    >
      加入房间
    </button>
  </view>
</template>

<script>
import { room } from '../../utils/api.js'

export default {
  data() {
    return {
      roomCode: '',
      loading: false
    }
  },
  methods: {
    async handleJoin() {
      if (this.roomCode.length !== 6) return
      this.loading = true
      try {
        const [err, res] = await room.join(this.roomCode)
        if (err) throw err
        uni.navigateTo({ url: `/pages/room/room?id=${res.room.id}` })
      } catch (e) {
        uni.showToast({ title: '房间不存在或已关闭', icon: 'none' })
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.page {
  padding: 48rpx 32rpx;
}

.form-item {
  margin-bottom: 48rpx;
}

.form-label {
  font-size: 28rpx;
  color: #64748b;
  margin-bottom: 16rpx;
}

.code-input {
  font-size: 40rpx;
  letter-spacing: 8rpx;
  text-align: center;
  font-weight: 700;
}
</style>
