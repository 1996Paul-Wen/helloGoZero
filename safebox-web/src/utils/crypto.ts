import CryptoJS from 'crypto-js'

const PBKDF2_ITERATIONS = 10000
const KEY_SIZE = 256 / 32 // 256 bits = 8 words (32-bit each)
const SALT = 'SafeBox2024SaltKey'

/**
 * 计算描述+用户名的 SHA256 哈希，取后 4 位字符
 */
function hashSuffix(description: string, username: string): string {
  const combined = description + username
  const hash = CryptoJS.SHA256(combined)
  return hash.toString(CryptoJS.enc.Hex).slice(-4)
}

/**
 * 从用户输入的字符串派生 AES-256 密钥
 * 派生规则：用户密钥 + hash(描述+用户名)后4位，再 PBKDF2
 */
function deriveKey(passphrase: string, description: string, username: string): CryptoJS.lib.WordArray {
  const finalKey = passphrase + hashSuffix(description, username)
  return CryptoJS.PBKDF2(finalKey, SALT, {
    keySize: KEY_SIZE,
    iterations: PBKDF2_ITERATIONS,
  })
}

/**
 * 使用 AES-256-CBC 加密明文密码
 * @param plainText 明文密码
 * @param passphrase 用户输入的加密密钥字符串
 * @param description 密码记录描述（参与密钥派生）
 * @param username 密码记录用户名（参与密钥派生）
 * @returns Base64 编码的密文（IV + ciphertext）
 */
export function encrypt(
  plainText: string,
  passphrase: string,
  description: string,
  username: string,
): string {
  const key = deriveKey(passphrase, description, username)
  const iv = CryptoJS.lib.WordArray.random(128 / 8) // 16 bytes random IV

  const encrypted = CryptoJS.AES.encrypt(plainText, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  })

  // 将 IV 拼接到密文前面，Base64 编码后传输
  const combined = iv.concat(encrypted.ciphertext)
  return combined.toString(CryptoJS.enc.Base64)
}

/**
 * 使用 AES-256-CBC 解密密文
 * @param cipherText Base64 编码的密文（IV + ciphertext）
 * @param passphrase 用户输入的解密密钥字符串
 * @param description 密码记录描述（参与密钥派生）
 * @param username 密码记录用户名（参与密钥派生）
 * @returns 解密后的明文密码
 */
export function decrypt(
  cipherText: string,
  passphrase: string,
  description: string,
  username: string,
): string {
  try {
    const key = deriveKey(passphrase, description, username)
    const combined = CryptoJS.enc.Base64.parse(cipherText)

    // 提取 IV（前16字节）和密文
    const iv = CryptoJS.lib.WordArray.create(combined.words.slice(0, 4)) // 16 bytes = 4 words
    const ciphertext = CryptoJS.lib.WordArray.create(combined.words.slice(4))

    const decrypted = CryptoJS.AES.decrypt(
      { ciphertext } as CryptoJS.lib.CipherParams,
      key,
      {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7,
      }
    )

    const result = decrypted.toString(CryptoJS.enc.Utf8)
    if (!result) throw new Error('解密失败：密钥不匹配或数据损坏')
    return result
  } catch (error) {
    console.error('Decrypt error:', error)
    throw new Error('解密失败，请检查加密密钥是否正确')
  }
}
