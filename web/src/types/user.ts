export interface LoginForm {
  student_id: string
  password: string
}

export interface RegisterForm {
  student_id: string
  name: string
  password: string
  confirmPassword?: string // 添加确认密码字段，但在API请求时不会发送
}

export interface User {
  id: number
  student_id: string
  name: string
}

export interface LoginResponse {
  code: number
  message: string
  data: {
    user: User
    token: string
  }
}

export interface RegisterResponse {
  code: number
  message: string
  data: {
    user: User
    token: string
  }
}