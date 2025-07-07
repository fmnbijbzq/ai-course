import type { ApiResponse } from './class'

export interface Teacher {
  id: number
  name: string
  title: string
  email: string
}

export interface TeacherListResponse extends ApiResponse<Teacher[]> {} 