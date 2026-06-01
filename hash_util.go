package manor

import (
    "crypto/md5"
    "encoding/hex"

    "github.com/OneOfOne/xxhash"
)

// xxhash 64 计算, 返回uint64的整数
func XXHashSum64(data []byte) uint64 {
    xh := xxhash.New64()
    xh.Write(data)
    return xh.Sum64()
}

// md5摘要计算, 返回为md5后的二进制, 长度为16
func MD5Digest(data []byte) []byte {
    h := md5.New()
    h.Write(data)
    return h.Sum(nil)
}

// md5摘要计算, 返回md5后的16进制字符串, 长度为32
func MD5HexDigest(data []byte) string {
    h := md5.New()
    h.Write(data)
    return hex.EncodeToString(h.Sum(nil))
}

