package manor

import (
    "time"
    "math"
    "errors"
    "encoding/base64"
    "encoding/hex"
    "encoding/binary"
)

// 数据结构常用操作的封装

// 判断整数切片中是否包含元素
func IntArrayContain(arr []int, v int) bool {
    for _, x := range arr {
        if x == v {
            return true
        }
    }
    return false
}

// unix时间戳转时间字符串
// timestamp: 时间戳, 单位: s
// timeFormat: 格式化方式, 例如: "2006-01-02 15:04:05"
func TimestampToString(timestamp int64, timeFormat string) string {
    unixTime := time.Unix(timestamp, 0)
    return unixTime.Format(timeFormat)
}

// 字符串格式的时间转时间戳, 
// 注意: 按照本地时区进行转换
// 返回: 转换成功error为nil, 返回时间戳, 单位为s; 转换失败error不为nil
func StringToTimestamp(timeAsStr, timeFormat string) (int64, error) {
    // UTC
    // ts, err := time.Parse(timeFormat, timeAsStr)
    loc, err := time.LoadLocation("Local")
    if err != nil {
        return 0, err
    }
    ts, err := time.ParseInLocation(timeFormat, timeAsStr, loc)
    if err != nil {
        return 0, err
    }
    return ts.Unix(), nil
}

// 按照指定的时区转换时间戳, 比如: UTC - UTC时间, Asia/Shanghai - 东八区时间
func TimeZoneStringToTimestamp(timeAsStr, timeFormat, timeZone string) (int64, error) {
    loc, err := time.LoadLocation(timeZone)
    if err != nil {
        return 0, err
    }
    ts, err := time.ParseInLocation(timeFormat, timeAsStr, loc)
    if err != nil {
        return 0, err
    }
    return ts.Unix(), nil
}


// 两个字节转float16的数值表示
// 返回值: float64 对应float16表示的浮点数
func ByteToFloat16(high, low byte) float64 {
    s := 1.0
    if high >> 7 == 1 {
        s = -1.0
    }

    exp := int(high & 0x7c) >> 2
    E := exp - 15

    man := int(high & 0x03) << 8 + int(low)
    M := 1.0
    for i := 0; i < 10; i++ {
        v := man & (0x0200 >> uint(i))
        if v != 0 {
            M += 1 / math.Pow(2, float64(i + 1))
        }
    }

    return s * M * math.Pow(2, float64(E))
}

// base64编码
func Base64Encode(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}

// base64解码
func Base64Decode(encodedData string) ([]byte, error) {
    return base64.StdEncoding.DecodeString(encodedData)
}

// 校验int是否在uint8范围内
// 校验通过返回value, 否则返回defaultValue默认值
func Uint8Transform(value, defaultValue int) int {
    if value >= 0 && value <= 255 {
        return value
    }
    return defaultValue
}

// 校验int是否在uint16范围内
// 校验通过返回value, 否则返回defaultValue默认值
func Uint16Transform(value, defaultValue int) int {
    if value >= 0 && value <= 65535 {
        return value
    }
    return defaultValue
}

// 编码二进制为hex格式的字符串
func EncodeToHex(src []byte) string {
    return hex.EncodeToString(src)
}

// 解码hex为二进制数组
func HexDecode(s string) ([]byte, error) {
    return hex.DecodeString(s)
}

// uint64数字小端方式打包字节流
func Uint64PackLittleEndian(num uint64) []byte {
    b := make([]byte, 8)
    b[0] = byte(num)
    b[1] = byte(num >> 8)
    b[2] = byte(num >> 16)
    b[3] = byte(num >> 24)
    b[4] = byte(num >> 32)
    b[5] = byte(num >> 40)
    b[6] = byte(num >> 48)
    b[7] = byte(num >> 56)
    return b
}

// 小端方式解包字节数组为uint64数字
func Uint64UnpackLittleEndian(b []byte) uint64 {
    _ = b[7]
    return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | 
        uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

// uint64数字大端方式打包字节流
func Uint64PackBigEndian(v uint64) []byte {
    b := make([]byte, 8)
    b[0] = byte(v >> 56)
    b[1] = byte(v >> 48)
    b[2] = byte(v >> 40)
    b[3] = byte(v >> 32)
    b[4] = byte(v >> 24)
    b[5] = byte(v >> 16)
    b[6] = byte(v >> 8)
    b[7] = byte(v)
    return b
}

// 从字节数组大端方式解包uint64
func Uint64UnpackBigEndian(b []byte) uint64 {
    _ = b[7]
    return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 | 
        uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56

}

// uint32小端打包
func Uint32PackLittleEndian(v uint32) []byte {
    b := make([]byte, 4)
    b[0] = byte(v)
    b[1] = byte(v >> 8)
    b[2] = byte(v >> 16)
    b[3] = byte(v >> 24)
    return b
}

// uint32小端解包
func Uint32UnpackLittleEndian(b []byte) uint32 {
    _ = b[3]
    return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

// uint32大端打包
func Uint32PackBigEndian(v uint32) []byte {
    b := make([]byte, 4)
    b[0] = byte(v >> 24)
    b[1] = byte(v >> 16)
    b[2] = byte(v >> 8)
    b[3] = byte(v)
    return b
}

// uint32大端解包
func Uint32UnpackBigEndian(b []byte) uint32 {
    _ = b[3]
    return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

// uint16小端打包
func Uint16PackLittleEndian(v uint16) []byte {
    b := make([]byte, 2)
    b[0] = byte(v)
    b[1] = byte(v >> 8)
    return b
}

// uint16小端解包
func Uint16UnpackLittleEndian(b []byte) uint16 {
    _ = b[1]
    return uint16(b[0]) | uint16(b[1])<<8
}

// uint16大端打包
func Uint16PackBigEndian(v uint16) []byte {
    b := make([]byte, 2)
    b[0] = byte(v >> 8)
    b[1] = byte(v)
    return b
}

// uint16大端解包
func Uint16UnpackBigEndian(b []byte) uint16 {
    _ = b[1]
    return uint16(b[1]) | uint16(b[0])<<8
}

// 打包uint8序列为二进制序列
func Uint8ToBytes(seq []uint8) []byte {
    dstBytes := make([]byte, len(seq))
    for i, value := range seq {
        dstBytes[i] = value
    }
    return dstBytes
}

// 打包uint16序列为二进制序列[大端]
func Uint16ToBytesBigEndian(seq []uint16) []byte {
    dstBytes := make([]byte, len(seq) * 2)
    for i, value := range seq {
        high := uint8(value >> 8)
        low := uint8(value)
        dstBytes[2 * i] = high
        dstBytes[2 * i + 1] = low
    }
    return dstBytes
}

// 打包uint16序列为二进制序列[小端]
func Uint16ToBytesLittleEndian(seq []uint16) []byte {
    dstBytes := make([]byte, len(seq) * 2)
    for i , value := range seq {
        high := uint8(value >> 8)
        low := uint8(value)
        dstBytes[2 * i] = low
        dstBytes[2 * i + 1] = high
    }
    return dstBytes
}

// 解包二进制数组为uint16切片[大端]
func BytesToUint16BigEndian(bs []byte) ([]uint16, error) {
    if len(bs) % 2 != 0 {
        return nil, errors.New("Byte length must be a multiple of 2!")
    }
    src := make([]uint16, len(bs) / 2)
    for i := 0; i < len(bs); i += 2 {
        src[i / 2] = uint16(bs[i]) << 8 | uint16(bs[i + 1])
    }
    return src, nil
}

// 解包二进制数组为uint16切片[小端]
func BytesToUint16LittleEndian(bs []byte) ([]uint16, error) {
    if len(bs) % 2 != 0 {
        return nil, errors.New("Byte length must be a multiple of 2!")
    }
    src := make([]uint16, len(bs) / 2)
    for i := 0; i < len(bs); i += 2 {
        src[i / 2] = uint16(bs[i + 1]) << 8 | uint16(bs[i])
    }
    return src, nil
}

// float32小端方式打包为字节数组
func Float32PackLittleEndian(f float32) []byte {
    // 转为uint32
    num := math.Float32bits(f)
    value := make([]byte, 4)
    binary.LittleEndian.PutUint32(value, num)
    return value
}

// float32小端方式解包
func Float32UnpackLittleEndian(v []byte) float32 {
    _ = v[3]
    num := binary.LittleEndian.Uint32(v)
    return math.Float32frombits(num)
}

// float64小端方式打包为字节数组
func Float64PackLittleEndian(f float64) []byte {
    // 转为uint64
    num := math.Float64bits(f)
    value := make([]byte, 8)
    binary.LittleEndian.PutUint64(value, num)
    return value
}

// float64小端方式解包
func Float64UnpackLittleEndian(v []byte) float64 {
    _ = v[7]
    num := binary.LittleEndian.Uint64(v)
    return math.Float64frombits(num)
}

// 设置uint32指定位的值
// index是从高位到低位, 范围: 1 - 32
func SetUint32Bit(num *uint32, index, value uint8) (err error) {
    if index > 32 || index == 0 {
        err = errors.New("Index range is 1 to 32.")
        return
    }
    err = nil
    switch value {
    case 0:
        *num &= ^(1 << (32 - index))
    case 1:
        *num |= 1 << (32 - index)
    default:
        err = errors.New("Value can only be 0 or 1.")
    }
    return
}

// 获取uint32指定位的值, index从高位到低位: 1 - 32
func GetUint32Bit(num uint32, index uint8) (value uint8, err error) {
    if index > 32 || index == 0 {
        err = errors.New("Index range is 1 to 32.")
        return
    }
    err = nil
    movBit := 32 - index
    value = uint8((num & (1 << movBit)) >> movBit)
    return
}

// 设置uint64指定位的值
// index是从高位到低位, index依次为: 1 - 64
func SetUint64Bit(num *uint64, index, value uint8) (err error) {
    if index > 64 || index == 0 {
        err = errors.New("Index range is 1 to 64.")
        return
    }
    err = nil
    switch value {
    case 0:
        *num &= ^(1 << (64 - index))
    case 1:
        *num |= 1 << (64 - index)
    default:
        err = errors.New("Value can only be 0 or 1.")
    }
    return
}

// 获取uint64指定位的值
// index从高位到低位依次为: 1 - 64
func GetUint64Bit(num uint64, index uint8) (value uint8, err error) {
    if index > 64 || index == 0 {
        err = errors.New("Index range is 1 to 64.")
        return
    }
    err = nil
    movBit := 64 - index
    value = uint8((num & (1 << movBit)) >> movBit)
    return
}


// 计算数字的绝对值, int64类型
func Int64Abs(v int64) int64 {
    n := v >> 63
    return (v ^ n) - n
}


// 计算区间[x1, y1]和[x2, y2]的交集区间
// 比如: [-1, 6]与[2, 8]返回: [2, 6], error为nil
// 如果重合在1点, 比如[3,5]与[5,7], 则返回: [5,5]
// 如果区间没有交集, 则error不为nil
func IntervalIntersection(x1, y1, x2, y2 int64) (left, right int64, err error) {
    if x2 <= y1 && x1 <= y2 {
        err = nil
        left, right = x1, y1
        if x2 > x1 {
            left = x2
        }
        if y2 < y1 {
            right = y2
        }
        return
    }
    err = errors.New("No intersection of the two groups.")
    return
}
