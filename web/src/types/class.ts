export interface Class {
  id: number
  code: string
  name: string
  description: string
  teacher_id: number
}

export interface ClassAddRequest {
  code: string
  name: string
  description?: string
  teacher_id: number
}

export interface ClassEditRequest {
  id: number
  code: string
  name: string
  description?: string
  teacher_id: number
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface PaginationData {
  list: Class[]
  total: number
}

export interface ClassListResponse {
  code: number
  message: string
  data: PaginationData
}

export interface ClassResponse {
  code: number
  message: string
  data: Class
}