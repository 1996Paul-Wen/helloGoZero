import request from './request'
import type { ApiResponse, SavePwdReq, PasswordItem } from '@/types'

/** 保存一条密码记录 */
export function savePassword(data: SavePwdReq) {
  return request.post<any, ApiResponse<null>>('/pwdManage/saveOne', data)
}

/** 搜索密码记录 */
export function queryPasswords(queryStr: string) {
  return request.post<any, ApiResponse<PasswordItem[]>>('/pwdManage/query', { query: queryStr })
}
