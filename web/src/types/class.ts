export interface Class {
  id: number
  class_name: string
  created_at: string
  updated_at: string
}

export interface ClassAddRequest {
  class_name: string
}

export interface ClassEditRequest {
  class_name: string
}

export interface ClassResponse {
  message: string
  data: Class
}

export interface ClassListResponse {
  message: string
  data: {
    total: number
    list: Class[]
  }
}