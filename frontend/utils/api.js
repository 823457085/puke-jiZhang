// API 配置
const BASE_URL = process.env.NODE_ENV === 'development'
  ? 'http://localhost:8080'
  : 'https://api.pukeapp.com' // TODO: 替换为实际域名

// 获取本地存储的 userID
function getUserId() {
  return uni.getStorageSync('userId') || ''
}

// 通用的请求封装
async function request(url, method = 'GET', data = {}) {
  const userId = getUserId()

  const header = {
    'Content-Type': 'application/json'
  }

  // MVP阶段：简单通过header传userID
  if (userId) {
    header['X-User-ID'] = userId
  }

  const token = uni.getStorageSync('token')
  if (token) {
    header['Authorization'] = `Bearer ${token}`
  }

  const [err, res] = await uni.request({
    url: BASE_URL + url,
    method,
    data,
    header
  })

  if (err) {
    return [err, null]
  }

  if (res.statusCode >= 400) {
    return [new Error(res.data?.error || '请求失败'), null]
  }

  return [null, res.data]
}

// ========== API 接口 ==========

// 认证
export const auth = {
  // 登录（MVP：直接用openid注册，生产需调用微信code换openid）
  login(openid, nickname, avatarUrl) {
    return request('/api/v1/auth/login', 'POST', {
      openid,
      nickname,
      avatar_url: avatarUrl
    })
  }
}

// 用户
export const user = {
  getMe() {
    return request('/api/v1/user/me')
  }
}

// 房间
export const room = {
  // 创建房间
  create(name, gameType = 'generic') {
    return request('/api/v1/rooms', 'POST', { name, game_type: gameType })
  },

  // 获取我的房间列表
  getMyRooms() {
    return request('/api/v1/rooms')
  },

  // 获取房间详情
  getRoom(roomId) {
    return request(`/api/v1/rooms/${roomId}`)
  },

  // 加入房间
  join(code) {
    return request('/api/v1/rooms/join', 'POST', { code })
  },

  // 关闭房间
  close(roomId) {
    return request(`/api/v1/rooms/${roomId}/close`, 'POST')
  },

  // 获取欠款余额
  getBalance(roomId) {
    return request(`/api/v1/rooms/${roomId}/balance`)
  },

  // 结算
  settle(roomId) {
    return request(`/api/v1/rooms/${roomId}/settle`, 'POST')
  }
}

// 账单
export const bill = {
  // 记一笔
  create(roomId, payerId, receiverId, amount, note = '') {
    return request(`/api/v1/rooms/${roomId}/bills`, 'POST', {
      payer_id: payerId,
      receiver_id: receiverId,
      amount,
      note
    })
  },

  // 删除账单
  delete(roomId, billId) {
    return request(`/api/v1/rooms/${roomId}/bills/${billId}`, 'DELETE')
  }
}

export default {
  auth,
  user,
  room,
  bill
}
