import request from '@/utils/request'
import type { 
  ClassAddRequest, 
  ClassEditRequest,
  Class,
  PaginationData,
  ClassResponse,
  ClassListResponse
} from '@/types/class'

export const addClass = (data: ClassAddRequest) => {
  return request.post<ClassResponse>('/api/class/add', data)
}

export const editClass = (id: number, data: ClassEditRequest) => {
  return request.put<ClassResponse>(`/api/class/${id}`, data)
}

export const deleteClass = (id: number) => {
  return request.delete<void>(`/api/class/${id}`)
}

export const getClassList = (params: { page: number; page_size: number }) => {
  return request.get<ClassListResponse>('/api/class/list', { params })
}