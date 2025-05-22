import axios from 'axios'
import type { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'

interface ApiResponse<T = any> {
  code?: number
  message: string
  data: T
}

// 创建自定义的 request 实例，并添加类型
const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
  withCredentials: true
}) as AxiosInstance & {
  get<T = any>(url: string, config?: any): Promise<T>
  post<T = any>(url: string, data?: any, config?: any): Promise<T>
  put<T = any>(url: string, data?: any, config?: any): Promise<T>
  delete<T = any>(url: string, config?: any): Promise<T>
}

// 确保 service 有 post 方法
console.log('Axios instance methods:', {
  post: typeof service.post,
  get: typeof service.get,
  put: typeof service.put,
  delete: typeof service.delete
})

service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 确保 headers 对象存在
    if (!config.headers) {
      // 如果 headers 不存在，我们跳过设置而不是创建新的
      return config
    }

    // 设置 Content-Type
    if (config.method?.toUpperCase() === 'POST') {
      config.headers['Content-Type'] = 'application/json'
    }

    // 添加认证 token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    return config
  },
  (error) => {
    console.error('Request Error:', error)
    return Promise.reject(error)
  }
)

service.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const data = response.data
    
    // 检查响应数据结构
    if (!data) {
      console.error('响应数据为空:', data)
      return Promise.reject(new Error('响应数据为空'))
    }

    // 检查业务状态码
    if (data.code !== undefined && data.code !== 0) {
      console.error('业务状态码错误:', data)
      return Promise.reject(new Error(data.message || '请求失败'))
    }

    // 如果是登录或注册接口
    if (response.config.url?.includes('/user/login') || response.config.url?.includes('/user/register')) {
      const responseData = data.data
      
      // 确保返回的数据包含必要的字段
      if (!responseData) {
        console.error('响应数据为空:', data)
        return Promise.reject(new Error('响应数据结构不正确: data 为空'))
      }
      if (!responseData.token) {
        console.error('缺少 token:', responseData)
        return Promise.reject(new Error('响应数据结构不正确: 缺少 token'))
      }
      if (!responseData.user) {
        console.error('缺少 user:', responseData)
        return Promise.reject(new Error('响应数据结构不正确: 缺少 user'))
      }
    }

    // 对于修改操作，显示成功消息
    const method = response.config.method?.toUpperCase()
    if (method && ['POST', 'PUT', 'DELETE'].includes(method)) {
      ElMessage({
        type: 'success',
        message: data.message || '操作成功'
      })
    }

    // 只返回响应数据中的 data 字段
    return data.data
  },
  (error) => {
    console.error('Response Error:', error)

    let errorMessage = '';

    // 处理401错误（用户账户密码错误）
    if (error.response?.status === 401) {
      errorMessage = '用户名或密码错误，请重新输入'

      // 清除token并延迟跳转到登录页，确保用户能看到错误信息
      localStorage.removeItem('token')
      setTimeout(() => {
        window.location.href = '/login'
      }, 1500)
    } else {
      // 处理其他错误
      const response = error.response?.data
      errorMessage = response?.message || error.message || '请求失败'
    }

    // 显示错误消息（只显示一次）
    ElMessage.error(errorMessage)

    return Promise.reject(new Error(errorMessage))
  }
)

export default service
