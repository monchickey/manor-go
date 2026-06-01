package manor

import (
    "fmt"
    "errors"
    "strings"
)

// 经纬度范围
const (
    MAX_LON = 180
    MIN_LON = -180
    MAX_LAT = 90
    MIN_LAT = -90
)

var geoBase32Chars = []string{
    "0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
    "b", "c", "d", "e", "f", "g", "h", "j", "k", "m",
    "n", "p", "q", "r", "s", "t", "u", "v", "w", "x",
    "y", "z",
}

var geoCodesByChar = map[byte]uint8{
    '0': 0,
    '1': 1,
    '2': 2,
    '3': 3,
    '4': 4,
    '5': 5,
    '6': 6,
    '7': 7,
    '8': 8,
    '9': 9,
    'b': 10,
    'c': 11,
    'd': 12,
    'e': 13,
    'f': 14,
    'g': 15,
    'h': 16,
    'j': 17,
    'k': 18,
    'm': 19,
    'n': 20,
    'p': 21,
    'q': 22,
    'r': 23,
    's': 24,
    't': 25,
    'u': 26,
    'v': 27,
    'w': 28,
    'x': 29,
    'y': 30,
    'z': 31,
}

// 对单个坐标值进行geohash编码, 位数: digits
// 注意是使用低digits位编码
func coordinateHash(coordinate, min, max float64, digits int) uint32 {
    hashValue := uint32(0)
    for i := 0; i < digits; i++ {
        mid := (min + max) / 2
        if(coordinate >= mid) {
            hashValue |= 1 << (digits - i - 1)
            min = mid
        } else {
            max = mid
        }
    }
    return hashValue
}


// 对单个坐标值做反向hash, 位数: digits
// 注意是使用uint32的高digits位
func coordinateReverseHash(hashValue uint32, min, max float64, digits int) float64 {
    mid := (min + max) / 2
    for i := 1; i <= digits; i++ {
        bitValue, _ := GetUint32Bit(hashValue, uint8(i))
        if bitValue == 1 {
            min = mid
        } else {
            max = mid
        }
        mid = (min + max) / 2
    }
    return mid
}


func base32Encode(src []byte) (string, error) {
    dst := make([]string, len(src))
    for i, c := range src {
        if c > 31 {
            err := errors.New("The source Numbers range from 0 to 31.")
            return "", err
        }
        dst[i] = geoBase32Chars[c]
    }
    return strings.Join(dst, ""), nil
}

func base32Decode(encoded string) ([]byte, error) {
    encodedBytes := []byte(encoded)
    decodedBytes := make([]byte,len(encodedBytes))
    for i, enc := range encodedBytes {
        if dec, ok := geoCodesByChar[enc]; ok {
            decodedBytes[i] = dec
        } else {
            err := errors.New(fmt.Sprintf("Not other characters! %c", enc))
            return decodedBytes, err
        }
    }
    return decodedBytes, nil
}

// geohash编码函数
// longitude: 经度
// latitude: 纬度
// precision: 精度, 范围: 1-12
// returns: geohash编码后的字符串, 错误信息
func GeohashEncode(longitude, latitude float64, precision int) (string, error) {
    if precision < 1 || precision > 12 {
        err := errors.New("Precision range from 1 to 12.")
        return "", err
    }

    // 总位数
    totalDigits := precision * 5
    // 单个坐标值hash位数
    lonDigits := totalDigits / 2
    latDigits := lonDigits
    if totalDigits % 2 != 0 {
        lonDigits++
    }

    lonHash := coordinateHash(longitude, MIN_LON, MAX_LON, lonDigits)
    latHash := coordinateHash(latitude, MIN_LAT, MAX_LAT, latDigits)
    lonHash = lonHash << (32 - lonDigits)
    latHash = latHash << (32 - latDigits)

    mergeHash := uint64(0)
    var i int
    for i = 1; i <= lonDigits; i++ {
        lonbit, _ := GetUint32Bit(lonHash, uint8(i))
        latbit, _ := GetUint32Bit(latHash, uint8(i))
        SetUint64Bit(&mergeHash, uint8(i * 2 - 1), lonbit)
        SetUint64Bit(&mergeHash, uint8(i * 2), latbit)
    }

    if totalDigits % 2 != 0 {
        lonbit, _ := GetUint32Bit(lonHash, uint8(i))
        SetUint64Bit(&mergeHash, uint8(i * 2 - 1), lonbit)
    }

    // 分组转换字节
    hashBytes := make([]byte, precision)
    for i = 1; i <= precision; i++ {
        hashBytes[i - 1] = uint8((mergeHash >> (64 - i * 5)) & 31)
    }

    geoDst, err := base32Encode(hashBytes)

    return geoDst, err
}


// 分割hash为两个坐标的子hash
func splitHash(value uint64, digits int) (uint32, uint32) {
    lonHash := uint32(0)
    latHash := uint32(0)

    var i = 1
    for ; i <= digits - 1; i+=2 {
        lonbit, _ := GetUint64Bit(value, uint8(i))
        latbit, _ := GetUint64Bit(value, uint8(i + 1))
        SetUint32Bit(&lonHash, uint8(i / 2 + 1), lonbit)
        SetUint32Bit(&latHash, uint8(i / 2 + 1), latbit)
    }

    if digits & 0x1 == 1 {
        lonbit, _ := GetUint64Bit(value, uint8(i))
        SetUint32Bit(&lonHash, uint8(i / 2 + 1), lonbit)
    }

    return lonHash, latHash
}

// geohash解码函数
// encoded: hash编码后的字符串, 长度: 1-12
// returns: 经度, 纬度, 错误信息
func GeohashDecode(encoded string) (longitude, latitude float64, err error) {
    precision := len(encoded)
    if precision > 12 || precision == 0 {
        err = errors.New("Encoded length can only be 1 to 12.")
        return
    }

    decoded, err := base32Decode(encoded)
    if err != nil {
        return
    }

    decodedHash := uint64(0)
    for i, dec := range decoded {
        movBit := 64 - (i + 1) * 5
        decodedHash |= (uint64(dec) << movBit)
    }

    lonHash, latHash := splitHash(decodedHash, precision * 5)
    lonDigits := precision * 5 / 2
    latDigits := lonDigits
    if precision & 0x1 == 1 {
        lonDigits++
    }

    longitude = coordinateReverseHash(lonHash, MIN_LON, MAX_LON, lonDigits)
    latitude = coordinateReverseHash(latHash, MIN_LAT, MAX_LAT, latDigits)

    return
}

