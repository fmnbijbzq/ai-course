import request from '@/utils/request'
import type { LoginForm, RegisterForm, LoginResponse, RegisterResponse } from '@/types/user'

export const login = (data: LoginForm) => {
  console.log('Logging in with data:', data)
  // 使用 axios 直接发送 POST 请求，确保方法正确
  return request({
    method: 'POST',
    url: '/user/login',
    data
  }) as Promise<LoginResponse>
}

export const register = (data: RegisterForm) => {
  console.log('Registering with data:', data)
  // 使用 axios 直接发送 POST 请求，确保方法正确
  return request({
    method: 'POST',
    url: '/user/register',
    data
  }) as Promise<RegisterResponse>
}