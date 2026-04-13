import request from './request'
import type { ApiResponse, SavePwdReq, UpdatePwdReq, PasswordItem } from '@/types'

/** 保存一条密码记录 */
export function savePassword(data: SavePwdReq) {
  return request.post<any, ApiResponse<null>>('/pwdManage/saveOne', data)
}

/** 搜索密码记录 */
export function queryPasswords(queryStr: string) {
  return request.post<any, ApiResponse<PasswordItem[]>>('/pwdManage/query', { query: queryStr })
}

/** 更新一条密码记录 */
export function updatePassword(data: UpdatePwdReq) {
  return request.post<any, ApiResponse<number>>('/pwdManage/updateOne', data)
}

/** 删除一条密码记录 */
export function deletePassword(id: number) {
  return request.post<any, ApiResponse<null>>('/pwdManage/deleteOne', { id })
}
