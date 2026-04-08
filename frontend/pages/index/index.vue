<template>
  <view class="page">
    <!-- 顶部：标题 -->
    <view class="header">
      <text class="page-title">记牌官</text>
      <view class="header-right" @click="goCreateRoom">
        <text>+ 创建房间</text>
      </view>
    </view>

    <!-- 进行中的房间 -->
    <view class="section" v-if="activeRooms.length > 0">
      <view class="section-title">进行中的房间</view>
      <view
        class="room-card card"
        v-for="room in activeRooms"
        :key="room.id"
        @click="goRoom(room)"
      >
        <view class="room-header">
          <view class="room-name">{{ room.name }}</view>
          <view class="room-tag tag-blue">进行中</view>
        </view>
        <view class="room-info">
          <text class="text-gray">{{ room.game_type_text }}</text>
          <text class="text-gray m-l">·</text>
          <text class="text-gray m-l">{{ room.member_count }}人</text>
        </view>
        <view class="room-balance" v-if="room.my_balance > 0">
          <text class="text-green">应收 ¥{{ (room.my_balance / 100).toFixed(2) }}</text>
        </view>
        <view class="room-balance" v-else-if="room.my_balance < 0">
          <text class="text-red">应付 ¥{{ (-room.my_balance / 100).toFixed(2) }}</text>
        </view>
      </view>
    </view>

    <!-- 历史房间 -->
    <view class="section" v-if="historyRooms.length > 0">
      <view class="section-title">历史房间</view>
      <view
        class="room-card card"
        v-for="room in historyRooms"
        :key="room.id"
        @click="goRoom(room)"
      >
        <view class="room-header">
          <view class="room-name">{{ room.name }}</view>
          <view class="room-tag tag-gray">已结算</view>
        </view>
        <view class="room-info">
          <text class="text-gray">{{ room.game_type_text }}</text>
          <text class="text-gray m-l">·</text>
          <text class="text-gray m-l">{{ room.created_at }}</text>
        </view>
      </view>
    </view>

    <!-- 空状态 -->
    <view class="empty-state" v-if="!loading && activeRooms.length === 0 && historyRooms.length === 0">
      <view class="empty-icon">🃏</view>
      <view class="empty-text">还没有房间</view>
      <view class="empty-sub">创建或加入一个房间开始记账吧</view>
      <button class="btn-primary" @click="goCreateRoom" style="width: 400rpx; margin-top: 40rpx;">
        创建房间
      </button>
    </view>

    <!-- 底部操作栏 -->
    <view class="bottom-bar">
      <view class="bottom-btn" @click="goJoinRoom">
        <text>加入房间</text>
      </view>
      <view class="bottom-btn btn-blue" @click="goCreateRoom">
        <text>+ 创建房间</text>
      </view>
    </view>
  </view>
</template>

<script>
import { room } from '../../utils/api.js'

export default {
  data() {
    return {
      loading: false,
      activeRooms: [],
      historyRooms: []
    }
  },
  onShow() {
    this.loadRooms()
  },
  onPullDownRefresh() {
    this.loadRooms().finally(() => uni.stopPullDownRefresh())
  },
  methods: {
    async loadRooms() {
      this.loading = true
      try {
        const [err, res] = await room.getMyRooms()
        if (err) throw err

        // 分离进行中和历史房间
        this.activeRooms = (res.rooms || []).filter(r => r.status === 'active').map(r => ({
          ...r,
          game_type_text: this.gameTypeText(r.game_type),
          member_count: r.member_count || 0,
          created_at: this.formatDate(r.created_at)
        }))

        this.historyRooms = (res.rooms || []).filter(r => r.status === 'closed').map(r => ({
          ...r,
          game_type_text: this.gameTypeText(r.game_type),
          created_at: this.formatDate(r.created_at)
        }))
      } catch (e) {
        uni.showToast({ title: '加载失败', icon: 'none' })
      } finally {
        this.loading = false
      }
    },
    goRoom(r) {
      uni.navigateTo({ url: `/pages/room/room?id=${r.id}` })
    },
    goCreateRoom() {
      uni.navigateTo({ url: '/pages/create-room/create-room' })
    },
    goJoinRoom() {
      uni.navigateTo({ url: '/pages/join-room/join-room' })
    },
    gameTypeText(type_) {
      const map = { poker: '扑克', mahjong: '麻将', generic: '通用' }
      return map[type_] || '通用'
    },
    formatDate(dateStr) {
      if (!dateStr) return ''
      const d = new Date(dateStr)
      const now = new Date()
      const diff = Math.floor((now - d) / (1000 * 60 * 60 * 24))
      if (diff === 0) return '今天'
      if (diff === 1) return '昨天'
      if (diff < 7) return `${diff}天前`
      return `${d.getMonth() + 1}月${d.getDate()}日`
    }
  }
}
</script>

<style scoped>
.page {
  padding: 0 32rpx 160rpx;
  min-height: 100vh;
  background: #f8fafc;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 0;
  position: sticky;
  top: 0;
  background: #f8fafc;
  z-index: 10;
}

.page-title {
  font-size: 40rpx;
  font-weight: 700;
  color: #1e293b;
}

.header-right {
  font-size: 28rpx;
  color: #2563eb;
}

.section {
  margin-bottom: 40rpx;
}

.section-title {
  font-size: 28rpx;
  color: #64748b;
  margin-bottom: 20rpx;
  font-weight: 500;
}

.room-card {
  position: relative;
}

.room-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12rpx;
}

.room-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #1e293b;
}

.room-tag {
  font-size: 22rpx;
  padding: 4rpx 12rpx;
  border-radius: 4rpx;
}

.tag-blue {
  background: #eff6ff;
  color: #2563eb;
}

.tag-gray {
  background: #f1f5f9;
  color: #64748b;
}

.room-info {
  font-size: 26rpx;
}

.m-l {
  margin-left: 12rpx;
}

.room-balance {
  position: absolute;
  right: 24rpx;
  bottom: 24rpx;
  font-size: 28rpx;
  font-weight: 700;
}

.empty-state {
  text-align: center;
  padding: 120rpx 0;
}

.empty-icon {
  font-size: 120rpx;
  margin-bottom: 24rpx;
}

.empty-text {
  font-size: 36rpx;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 12rpx;
}

.empty-sub {
  font-size: 28rpx;
  color: #64748b;
}

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  gap: 24rpx;
  padding: 24rpx 32rpx;
  background: #ffffff;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.06);
}

.bottom-btn {
  flex: 1;
  height: 96rpx;
  line-height: 96rpx;
  text-align: center;
  background: #f1f5f9;
  color: #1e293b;
  border-radius: 8rpx;
  font-size: 32rpx;
  font-weight: 500;
}

.btn-blue {
  background: #2563eb;
  color: #ffffff;
}
</style>
