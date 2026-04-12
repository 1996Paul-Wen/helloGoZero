/** 后端统一响应格式 */
export interface ApiResponse<T = any> {
  Code: number
  Msg: string
  Data: T
  TraceID: string
}

/** 用户信息 */
export interface UserInfo {
  Id: number
  Username: string
  HashPassword: string
  Creator: string
  Updator: string
  CreateTime: string
  UpdateTime: string
}

/** 登录响应 */
export interface LoginData {
  token: string
}

/** 创建/登录用户请求 */
export interface CreateUserReq {
  name: string
  password: string
}

/** 密码管理条目 */
export interface PasswordItem {
  Id: number
  UserId: number
  Description: string
  Username: string
  Password: string
  Creator: string
  Updator: string
  CreateTime: string
  UpdateTime: string
}

/** 保存密码请求 */
export interface SavePwdReq {
  description: string
  username: string
  password: string
}

/** 查询密码请求 */
export interface QueryPwdReq {
  query: string
}
