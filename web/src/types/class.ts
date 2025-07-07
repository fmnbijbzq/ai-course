export interface Class {
  id: number
  name: string
  description: string
  teacher_id: number
}

export interface ClassAddRequest {
  name: string
  description?: string
  teacher_id: number
}

export interface ClassEditRequest {
  id: number
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
  page: number
  page_size: number
}

export interface ClassListResponse {
  list: Class[]
  total: number
  page: number
  page_size: number
}

export interface ClassResponse extends ApiResponse<Class> {}