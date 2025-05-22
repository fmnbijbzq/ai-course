import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types/user'

export const useUserStore = defineStore('user', () => {
  // 尝试从localStorage获取用户信息
  const userJson = localStorage.getItem('user')
  let savedUser: User | null = null

  try {
    if (userJson) {
      savedUser = JSON.parse(userJson)
    }
  } catch (e) {
    console.error('Failed to parse user from localStorage', e)
  }

  // 状态
  const currentUser = ref<User | null>(savedUser)
  const token = ref<string | null>(localStorage.getItem('token'))

  // Actions
  function setUser(user: User) {
    currentUser.value = user
    // 将用户信息保存到localStorage
    localStorage.setItem('user', JSON.stringify(user))
  }

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function logout() {
    currentUser.value = null
    token.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  // 检查是否已登录
  function isLoggedIn(): boolean {
    return !!token.value // 只检查token是否存在
  }

  return {
    currentUser,
    token,
    setUser,
    setToken,
    logout,
    isLoggedIn
  }
})