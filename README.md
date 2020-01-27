# Backend

## Running a Server
1. `cp config.example.yml config.yml`
2. 相应位置填入石墨文档 Cookie 以及文档 ID
3. 使用 `go run .` 即可运行项目

## Usage
- 住宿信息列表 ([https://shimo.im/sheets/6c6GKvX83hRCVdG8](https://shimo.im/sheets/6c6GKvX83hRCVdG8))
    - `GET /accommodations/json` 返回 JSON 形式的住宿信息列表
    - `GET /accommodations/csv` 返回 CSV 形式的住宿信息列表
- 诊断平台信息列表 ([https://shimo.im/sheets/kDQJ6vWgWWwq8r8H](https://shimo.im/sheets/kDQJ6vWgWWwq8r8H))
    - `GET /platforms/json` 返回 JSON 形式的诊断平台信息列表
    - `GET /platforms/csv` 返回 CSV 形式的诊断平台信息列表
