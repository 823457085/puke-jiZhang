<template>
  <view class="page">
    <view class="tab-bar">
      <view
        class="tab"
        :class="{ active: tab === 'create' }"
        @click="tab = 'create'"
      >创建房间</view>
      <view
        class="tab"
        :class="{ active: tab === 'join' }"
        @click="tab = 'join'"
      >加入房间</view>
    </view>

    <!-- 创建房间 -->
    <view v-if="tab === 'create'" class="form">
      <view class="form-item">
        <view class="form-label">房间名称</view>
        <input
          class="input"
          v-model="roomName"
          placeholder="给房间起个名字"
          placeholder-class="input-placeholder"
        />
      </view>

      <view class="form-item">
        <view class="form-label">游戏类型</view>
        <view class="type-list">
          <view
            class="type-chip"
            :class="{ active: gameType === 'poker' }"
            @click="gameType = 'poker'"
          >扑克</view>
          <view
            class="type-chip"
            :class="{ active: gameType === 'mahjong' }"
            @click="gameType = 'mahjong'"
          >麻将</view>
          <view
            class="type-chip"
            :class="{ active: gameType === 'generic' }"
            @click="gameType = 'generic'"
          >通用</view>
        </view>
      </view>

      <button
        class="btn-primary"
        @click="handleCreate"
        :loading="loading"
        :disabled="!roomName"
      >
        创建房间
      </button>
    </view>

    <!-- 加入房间 -->
    <view v-if="tab === 'join'" class="form">
      <view class="form-item">
        <view class="form-label">房间码</view>
        <input
          class="input"
          v-model="roomCode"
          placeholder="请输入6位房间码"
          placeholder-class="input-placeholder"
          maxlength="6"
          type="number"
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
  </view>
</template>

<script>
import { room } from '../../utils/api.js'

export default {
  data() {
    return {
      tab: 'create',
      roomName: '',
      gameType: 'generic',
      roomCode: '',
      loading: false
    }
  },
  methods: {
    async handleCreate() {
      if (!this.roomName) return
      this.loading = true
      try {
        const [err, res] = await room.create(this.roomName, this.gameType)
        if (err) throw err
        uni.navigateTo({ url: `/pages/room/room?id=${res.room.id}` })
      } catch (e) {
        uni.showToast({ title: e.message || '创建失败', icon: 'none' })
      } finally {
        this.loading = false
      }
    },
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
  padding: 32rpx;
}

.tab-bar {
  display: flex;
  background: #f1f5f9;
  border-radius: 8rpx;
  padding: 6rpx;
  margin-bottom: 48rpx;
}

.tab {
  flex: 1;
  height: 72rpx;
  line-height: 72rpx;
  text-align: center;
  font-size: 30rpx;
  color: #64748b;
  border-radius: 6rpx;
}

.tab.active {
  background: #fff;
  color: #2563eb;
  font-weight: 600;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.08);
}

.form-item {
  margin-bottom: 40rpx;
}

.form-label {
  font-size: 28rpx;
  color: #64748b;
  margin-bottom: 16rpx;
}

.type-list {
  display: flex;
  gap: 16rpx;
}

.type-chip {
  flex: 1;
  height: 80rpx;
  line-height: 80rpx;
  text-align: center;
  background: #f1f5f9;
  color: #1e293b;
  border-radius: 8rpx;
  font-size: 30rpx;
  border: 2rpx solid transparent;
}

.type-chip.active {
  background: #eff6ff;
  color: #2563eb;
  border-color: #2563eb;
}
</style>
