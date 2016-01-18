package servercfg

import (
    "os"
    "encoding/json"
    "fmt"
)

type ServerCfg struct {
    ClientPath string `json:"clientPath"`
    Error404 string `json:"error404"`
    Port string `json:"port"`
    DefaultExt string `json:"defaultExt"`
}

func Load(cfgPath string) (ServerCfg) {
    file, err := os.Open(cfgPath)
    if err != nil {
        fmt.Println(cfgPath + " not found. Using default server config.")
    }
    
    fBytes := make([]byte, 65536)
    _, err = file.Read(fBytes)
    
    defer file.Close()
    
    cfg := ServerCfg{
        ClientPath: "client",
        Error404: "404.html",
        Port: "8080",
        DefaultExt: ".html",
    }
    
    json.Unmarshal(fBytes, &cfg)
    
    return cfg
}
