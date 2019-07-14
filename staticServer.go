package main

import (
    "net/http"
)

// 启动前端
func main() {
    http.ListenAndServe(":8080",http.FileServer(http.Dir("static")))
}
