package manor

import (
    "os"

    "github.com/goccy/go-yaml"
)

// 加载yaml配置文件
// configFile: 配置文件路径; configs: 结构体指针, 注意是传地址
func GetYamlConfig(configFile string, configs interface{}) error {
    configContent, err := os.ReadFile(configFile)
    if err != nil {
        return err
    }
    err = yaml.Unmarshal(configContent, configs)
    if err != nil {
        return err
    }
    return nil
}
