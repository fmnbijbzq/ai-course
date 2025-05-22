import request from '@/utils/request'
import type { 
  ClassAddRequest, 
  ClassEditRequest,
  Class,
  PaginationData
} from '@/types/class'

export const addClass = (data: ClassAddRequest) => {
  return request.post<Class>('/class/add', data)
}

export const editClass = (id: number, data: ClassEditRequest) => {
  return request.put<Class>(`/class/${id}`, data)
}

export const deleteClass = (id: number) => {
  return request.delete<void>(`/class/${id}`)
}

export const getClassList = (params: { page: number; page_size: number }) => {
  return request.get<PaginationData>('/class/list', { params })
}