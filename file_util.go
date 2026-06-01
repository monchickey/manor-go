package manor

import (
    "os"
    "encoding/gob"
)

// 判断文件或目录是不是存在
func PathIsExist(path string) bool {
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}

// 判断path是不是文件
func PathIsFile(path string) bool {
    fi, err := os.Stat(path)
    if err != nil {
        return false
    }
    return !fi.IsDir()
}

// 判断path是不是目录
func PathIsDir(path string) bool {
    fi, err := os.Stat(path)
    if err != nil {
        return false
    }
    return fi.IsDir()
}

// 获得文件的修改时间, 返回字符串示例: 2019-01-01 11:11:00
// 当文件不存在或出错时返回空字符串
func GetFileModifyTime(filename string) string {
    fileInfo, err := os.Stat(filename)
    if err == nil {
        modTime := fileInfo.ModTime()
        return modTime.Format("2006-01-02 15:04:05")
    }
    return ""
}

// 将指定内容以gob方式序列化到文件
// 特殊数据类型需要提前注册
// filePath: 文件绝对路径
func GobSerialize(filePath string, content interface{}) error {
    fp, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return err
    }
    defer fp.Close()
    encoder := gob.NewEncoder(fp)
    err = encoder.Encode(content)
    if err != nil {
        return err
    }
    return nil
}

// 将gob内容反序列化为数据
// filePath: 文件绝对路径
// content: 接收返回内容的指针
func GobDeserialize(filePath string, content interface{}) error {
    fp, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer fp.Close()
    decoder := gob.NewDecoder(fp)
    err = decoder.Decode(content)
    if err != nil {
        return err
    }
    return nil
}
