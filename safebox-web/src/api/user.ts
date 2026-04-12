import request from './request'
import type { ApiResponse, CreateUserReq, LoginData, UserInfo } from '@/types'

/** 创建用户 */
export function createUser(data: CreateUserReq) {
  return request.post<any, ApiResponse<{ UserID: number }>>('/user/create', data)
}

/** 用户登录 */
export function loginUser(data: CreateUserReq) {
  return request.post<any, ApiResponse<LoginData>>('/user/login', data)
}

/** 获取当前用户信息 */
export function describeUser() {
  return request.post<any, ApiResponse<UserInfo>>('/user/describe', {})
}
