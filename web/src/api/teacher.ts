import request from '@/utils/request'
import type { TeacherListResponse } from '@/types/teacher'

export const getTeacherList = () => {
  return request.get<TeacherListResponse>('/api/teacher/list')
} 