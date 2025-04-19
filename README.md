# Dify 扩展

## 调试

### `ping`

```bash
curl -X POST 'http://127.0.0.1:44444/dify' -H 'Content-Type: application/json' -d '{
  "point": "ping"
}'
```

### `foods`

```bash
curl -X POST 'http://127.0.0.1:44444/dify' -H 'Content-Type: application/json' -d '{
  "point": "app.external_data_tool.query",
  "params": {
    "app_id": "93d5cc7a-5f61-4734-ba89-84706bbc8c97",
    "tool_variable": "food_list",
    "inputs": {
      "choice": "foods"
    },
    "query": "\u4eca\u5929\u6709\u4ec0\u4e48\u597d\u5403\u7684\uff1f\"
  }
}'
```

### `drinks`

```bash
curl -X POST 'http://127.0.0.1:44444/dify' -H 'Content-Type: application/json' -d '{
  "point": "app.external_data_tool.query",
  "params": {
    "app_id": "93d5cc7a-5f61-4734-ba89-84706bbc8c97",
    "tool_variable": "food_list",
    "inputs": {
      "choice": "drinks"
    },
    "query": "\u4eca\u5929\u6709\u4ec0\u4e48\u597d\u559d\u7684\uff1f"
  }
}'
```

## 参考

### 源码

- [Dify NodeJs Extension Template](https://github.com/crazywoola/dify-extension)

### 视频

- [Dify 外部数据工具 Function Call => External data tool](https://www.bilibili.com/video/BV1UN4y1b7s1/)