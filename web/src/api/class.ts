import request from '@/utils/request'
import type { 
  ClassAddRequest, 
  ClassEditRequest, 
  ClassResponse, 
  ClassListResponse 
} from '@/types/class'

export const addClass = (data: ClassAddRequest) => {
  return request.post<ClassResponse>('/class/add', data)
}

export const editClass = (id: number, data: ClassEditRequest) => {
  return request.put<ClassResponse>(`/class/${id}`, data)
}

export const deleteClass = (id: number) => {
  return request.delete<{ message: string }>(`/class/${id}`)
}

export const getClassList = (params: { page: number; page_size: number }) => {
  return request.get<ClassListResponse>('/class/list', { params })
}