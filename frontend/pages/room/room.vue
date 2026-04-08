<template>
  <view class="page">
    <!-- 顶部：房间码 -->
    <view class="room-code-card card">
      <view class="code-label">房间码</view>
      <view class="code-value" @click="copyCode">{{ roomInfo.room_code }}</view>
      <view class="share-btn" @click="shareRoom">分享邀请好友</view>
    </view>

    <!-- 成员列表 -->
    <view class="card members-card">
      <view class="card-title">成员 {{ members.length }}人</view>
      <view class="member-item" v-for="m in members" :key="m.id">
        <view class="member-avatar">{{ m.user.nickname[0] }}</view>
        <view class="member-name">{{ m.user.nickname }}</view>
        <view
          class="member-balance"
          :class="getBalanceClass(m.user.id)"
        >
          {{ formatBalance(m.user.id) }}
        </view>
      </view>
    </view>

    <!-- 账单流水 -->
    <view class="section">
      <view class="section-title">账单流水</view>
      <view class="bill-item card" v-for="b in bills" :key="b.id">
        <view class="bill-payer">{{ b.payer.nickname }}</view>
        <view class="bill-arrow">→</view>
        <view class="bill-receiver">{{ b.receiver.nickname }}</view>
        <view class="bill-amount">+¥{{ (b.amount / 100).toFixed(2) }}</view>
        <view class="bill-note" v-if="b.note">{{ b.note }}</view>
        <view class="bill-time">{{ formatTime(b.created_at) }}</view>
      </view>

      <view class="empty-bills" v-if="bills.length === 0">
        <text class="text-gray">还没有账单，记录第一笔吧</text>
      </view>
    </view>

    <!-- 浮动按钮 -->
    <view class="fab" @click="showAddBill">
      <text>+ 记一笔</text>
    </view>

    <!-- 底部结算栏 -->
    <view class="bottom-bar" v-if="roomInfo.status === 'active'">
      <view class="settle-btn" @click="goSettle">
        <text>去结算</text>
      </view>
    </view>

    <!-- 记一笔弹窗 -->
    <view class="modal-mask" v-if="showBillModal" @click="closeBillModal">
      <view class="modal-content" @click.stop>
        <view class="modal-header">
          <text class="modal-title">记一笔</text>
          <text class="modal-close" @click="closeBillModal">✕</text>
        </view>

        <view class="modal-body">
          <view class="select-label">谁付钱？</view>
          <view class="member-select">
            <view
              class="member-chip"
              :class="{ active: selectedPayer === m.user.id }"
              v-for="m in members"
              :key="m.user.id"
              @click="selectedPayer = m.user.id"
            >
              {{ m.user.nickname }}
            </view>
          </view>

          <view class="select-label">谁收钱？</view>
          <view class="member-select">
            <view
              class="member-chip"
              :class="{ active: selectedReceiver === m.user.id }"
              v-for="m in members"
              :key="m.user.id"
              @click="selectedReceiver = m.user.id"
            >
              {{ m.user.nickname }}
            </view>
          </view>

          <view class="select-label">金额</view>
          <view class="amount-input-wrap">
            <text class="yuan-sign">¥</text>
            <input
              class="amount-input"
              type="digit"
              v-model="billAmount"
              placeholder="0.00"
              placeholder-class="amount-placeholder"
            />
          </view>

          <view class="select-label">备注（选填）</view>
          <input
            class="input note-input"
            v-model="billNote"
            placeholder="第3局斗地主..."
            placeholder-class="input-placeholder"
          />

          <button
            class="btn-primary"
            @click="submitBill"
            :disabled="!canSubmitBill"
          >
            确认
          </button>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import { room, bill } from '../../utils/api.js'

export default {
  data() {
    return {
      roomId: null,
      roomInfo: {},
      members: [],
      bills: [],
      balances: [],
      showBillModal: false,
      selectedPayer: '',
      selectedReceiver: '',
      billAmount: '',
      billNote: ''
    }
  },
  computed: {
    canSubmitBill() {
      return this.selectedPayer && this.selectedReceiver &&
             this.selectedPayer !== this.selectedReceiver &&
             this.billAmount > 0
    }
  },
  onLoad(options) {
    this.roomId = options.id
    this.loadRoom()
  },
  methods: {
    async loadRoom() {
      uni.showLoading({ title: '加载中...' })
      try {
        const [err, res] = await room.getRoom(this.roomId)
        if (err) throw err

        this.roomInfo = res.room || {}
        this.members = res.members || []
        this.bills = res.bills || []
        this.balances = res.balances || []

        // 默认选中当前用户为付款人
        const userId = uni.getStorageSync('userId')
        if (this.members.length > 0) {
          this.selectedPayer = this.members[0].user.id
        }
      } catch (e) {
        uni.showToast({ title: '加载失败', icon: 'none' })
      } finally {
        uni.hideLoading()
      }
    },
    getBalanceClass(userId) {
      const bal = this.balances.find(b => b.user_id === userId)
      if (!bal) return ''
      if (bal.balance > 0) return 'balance-positive'
      if (bal.balance < 0) return 'balance-negative'
      return ''
    },
    formatBalance(userId) {
      const bal = this.balances.find(b => b.user_id === userId)
      if (!bal) return '¥0.00'
      const amount = bal.balance / 100
      if (amount > 0) return `+¥${amount.toFixed(2)}`
      if (amount < 0) return `-¥${(-amount).toFixed(2)}`
      return '¥0.00'
    },
    formatTime(dateStr) {
      if (!dateStr) return ''
      const d = new Date(dateStr)
      return `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
    },
    copyCode() {
      uni.setClipboardData({
        data: this.roomInfo.room_code,
        success: () => uni.showToast({ title: '房间码已复制', icon: 'none' })
      })
    },
    shareRoom() {
      uni.showModal({
        title: '邀请好友',
        content: `房间码：${this.roomInfo.room_code}\n让好友在"加入房间"中输入此码`,
        showCancel: false
      })
    },
    showAddBill() {
      this.showBillModal = true
    },
    closeBillModal() {
      this.showBillModal = false
      this.selectedPayer = ''
      this.selectedReceiver = ''
      this.billAmount = ''
      this.billNote = ''
    },
    async submitBill() {
      if (!this.canSubmitBill) return

      const amountFen = Math.round(parseFloat(this.billAmount) * 100)

      try {
        const [err] = await bill.create(
          this.roomId,
          this.selectedPayer,
          this.selectedReceiver,
          amountFen,
          this.billNote
        )
        if (err) throw err

        uni.showToast({ title: '记账成功' })
        this.closeBillModal()
        this.loadRoom()
      } catch (e) {
        uni.showToast({ title: e.message || '记账失败', icon: 'none' })
      }
    },
    goSettle() {
      uni.navigateTo({ url: `/pages/settle/settle?id=${this.roomId}` })
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

.room-code-card {
  text-align: center;
  margin-bottom: 24rpx;
}

.code-label {
  font-size: 26rpx;
  color: #64748b;
  margin-bottom: 8rpx;
}

.code-value {
  font-size: 56rpx;
  font-weight: 700;
  color: #2563eb;
  letter-spacing: 8rpx;
  margin-bottom: 16rpx;
}

.share-btn {
  font-size: 26rpx;
  color: #2563eb;
}

.card-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 20rpx;
}

.member-item {
  display: flex;
  align-items: center;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f1f5f9;
}

.member-item:last-child {
  border-bottom: none;
}

.member-avatar {
  width: 64rpx;
  height: 64rpx;
  background: #2563eb;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  font-weight: 600;
  margin-right: 16rpx;
}

.member-name {
  flex: 1;
  font-size: 30rpx;
  color: #1e293b;
}

.member-balance {
  font-size: 32rpx;
  font-weight: 700;
  color: #1e293b;
}

.balance-positive { color: #16a34a; }
.balance-negative { color: #dc2626; }

.section-title {
  font-size: 28rpx;
  color: #64748b;
  margin-bottom: 20rpx;
  font-weight: 500;
}

.bill-item {
  position: relative;
  padding: 20rpx 24rpx;
}

.bill-payer, .bill-receiver {
  display: inline;
  font-size: 30rpx;
  font-weight: 500;
}

.bill-arrow {
  display: inline;
  color: #64748b;
  margin: 0 12rpx;
}

.bill-amount {
  position: absolute;
  right: 24rpx;
  top: 20rpx;
  font-size: 32rpx;
  font-weight: 700;
  color: #16a34a;
}

.bill-note {
  display: block;
  font-size: 24rpx;
  color: #64748b;
  margin-top: 8rpx;
}

.bill-time {
  position: absolute;
  right: 24rpx;
  bottom: 20rpx;
  font-size: 24rpx;
  color: #94a3b8;
}

.empty-bills {
  text-align: center;
  padding: 60rpx 0;
}

.fab {
  position: fixed;
  right: 32rpx;
  bottom: 180rpx;
  width: 120rpx;
  height: 120rpx;
  background: #2563eb;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  font-weight: 500;
  box-shadow: 0 4rpx 20rpx rgba(37, 99, 235, 0.4);
  z-index: 100;
}

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 24rpx 32rpx;
  background: #fff;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.06);
}

.settle-btn {
  height: 96rpx;
  line-height: 96rpx;
  text-align: center;
  background: #16a34a;
  color: #fff;
  border-radius: 8rpx;
  font-size: 34rpx;
  font-weight: 600;
}

/* 弹窗 */
.modal-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 1000;
  display: flex;
  align-items: flex-end;
}

.modal-content {
  background: #fff;
  width: 100%;
  border-radius: 32rpx 32rpx 0 0;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 32rpx 32rpx 24rpx;
  border-bottom: 1rpx solid #f1f5f9;
}

.modal-title {
  font-size: 34rpx;
  font-weight: 600;
  color: #1e293b;
}

.modal-close {
  font-size: 40rpx;
  color: #94a3b8;
  padding: 8rpx;
}

.modal-body {
  padding: 32rpx;
}

.select-label {
  font-size: 28rpx;
  color: #64748b;
  margin-bottom: 16rpx;
  margin-top: 24rpx;
}

.select-label:first-child {
  margin-top: 0;
}

.member-select {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.member-chip {
  padding: 12rpx 24rpx;
  background: #f1f5f9;
  color: #1e293b;
  border-radius: 8rpx;
  font-size: 28rpx;
  border: 2rpx solid transparent;
}

.member-chip.active {
  background: #eff6ff;
  color: #2563eb;
  border-color: #2563eb;
}

.amount-input-wrap {
  display: flex;
  align-items: center;
  background: #f8fafc;
  border: 2rpx solid #e2e8f0;
  border-radius: 8rpx;
  padding: 0 24rpx;
  height: 96rpx;
}

.yuan-sign {
  font-size: 40rpx;
  font-weight: 700;
  color: #1e293b;
  margin-right: 8rpx;
}

.amount-input {
  flex: 1;
  font-size: 40rpx;
  font-weight: 700;
  background: transparent;
}

.note-input {
  font-size: 30rpx;
}

.modal-body .btn-primary {
  margin-top: 40rpx;
}
</style>
