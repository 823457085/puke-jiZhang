<template>
  <view class="page">
    <!-- 结算结果 -->
    <view class="card">
      <view class="card-title">结算结果</view>
      <view class="balance-item" v-for="b in balances" :key="b.user_id">
        <view class="balance-name">{{ b.user.nickname }}</view>
        <view
          class="balance-amount"
          :class="b.balance > 0 ? 'text-green' : b.balance < 0 ? 'text-red' : ''"
        >
          {{ b.balance > 0 ? '+' : '' }}{{ (b.balance / 100).toFixed(2) }}元
        </view>
      </view>
    </view>

    <!-- 最优还款方案 -->
    <view class="card" v-if="settlements.length > 0">
      <view class="card-title">最优还款方案</view>
      <view class="settlement-item" v-for="(s, i) in settlements" :key="i">
        <view class="settlement-arrow">
          <text class="from-name">{{ s.from_user_name }}</text>
          <text class="arrow">→</text>
          <text class="to-name">{{ s.to_user_name }}</text>
        </view>
        <view
          class="settlement-amount"
          :class="s.amount > 0 ? 'text-green' : ''"
        >
          {{ (s.amount / 100).toFixed(2) }}元
        </view>
      </view>
    </view>

    <!-- 无需结算 -->
    <view class="empty-state card" v-if="balances.length > 0 && settlements.length === 0">
      <text class="text-gray">所有账目已结清，无需转账</text>
    </view>

    <!-- 确认结算 -->
    <button class="btn-primary confirm-btn" @click="confirmSettle">
      确认结算并分享
    </button>
  </view>
</template>

<script>
import { room } from '../../utils/api.js'

export default {
  data() {
    return {
      roomId: null,
      balances: [],
      settlements: [],
      loading: false
    }
  },
  async onLoad(options) {
    this.roomId = options.id
    await this.loadSettlement()
  },
  methods: {
    async loadSettlement() {
      uni.showLoading({ title: '计算中...' })
      try {
        // 获取余额
        const [err1, res1] = await room.getBalance(this.roomId)
        if (err1) throw err1
        this.balances = res1.balances || []

        // 获取还款方案
        const [err2, res2] = await room.settle(this.roomId)
        if (err2) throw err2
        this.settlements = res2.settlements || []
      } catch (e) {
        uni.showToast({ title: '加载失败', icon: 'none' })
      } finally {
        uni.hideLoading()
      }
    },
    async confirmSettle() {
      uni.showModal({
        title: '确认结算',
        content: '结算后房间将关闭，确定要结算吗？',
        success: async (res) => {
          if (res.confirm) {
            try {
              await room.close(this.roomId)
              uni.showToast({ title: '结算成功' })
              setTimeout(() => {
                uni.switchTab({ url: '/pages/index/index' })
              }, 1500)
            } catch (e) {
              uni.showToast({ title: e.message || '结算失败', icon: 'none' })
            }
          }
        }
      })
    }
  }
}
</script>

<style scoped>
.page {
  padding: 24rpx 32rpx 160rpx;
  min-height: 100vh;
  background: #f8fafc;
}

.card-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 24rpx;
}

.balance-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f1f5f9;
}

.balance-item:last-child {
  border-bottom: none;
}

.balance-name {
  font-size: 32rpx;
  color: #1e293b;
}

.balance-amount {
  font-size: 36rpx;
  font-weight: 700;
  color: #1e293b;
}

.settlement-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f1f5f9;
}

.settlement-item:last-child {
  border-bottom: none;
}

.settlement-arrow {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.from-name, .to-name {
  font-size: 30rpx;
  color: #1e293b;
  font-weight: 500;
}

.arrow {
  color: #64748b;
  font-size: 28rpx;
}

.settlement-amount {
  font-size: 32rpx;
  font-weight: 700;
  color: #1e293b;
}

.empty-state {
  text-align: center;
  padding: 60rpx 0;
}

.confirm-btn {
  margin-top: 40rpx;
}
</style>
