# AHNU-Retrieval
> 检索AHNU图书馆空闲座位脚本，该脚本仅供学习
### 一、 安装使用
#### 1 clone该项目
``` bash
git clone https://github.com/lthyxsl/AHNU-Retrieval.git
```
#### 2. 安装依赖
``` bash
go mod download
```

#### 3. 修改配置
> 3.1 复制 conf/config.json.template 为conf/config.json
> 3.2 修改 conf/config.json 文件中的 *tbUserName* 、*tbPassWord*、 *date* 和 *option*等字段值
#### 4. 运行程序
``` bash
go run main.go
```

### 二、 配置解析
```
{
  "tbUserName": "", // 登录用户名，即学号
  "tbPassWord": "", // 登录密码
  "date": "",   // 预约时间，默认值为当天， 格式 "YYYY-MM-DD" 如 "2022-05-20"
  "option": 0,  // 预约的选项，默认值为0，即 花津校区图书馆2楼电子阅览室  其值 与 *urls*数组中的*index*值对应
  "urls": [   // 预约地址配置
    {
      "index": 0,
      "title": "花津校区图书馆2楼电子阅览室",
      "url": "https://libzwxt.ahnu.edu.cn/seatwx/Room.aspx?rid=18&fid=1"
    },
    {
      "index": 1,
      "title": "花津校区图书馆2楼南",
      "url": "https://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=1&fid=1"
    },
    {
      "index": 2,
      "title": "花津校区图书馆3楼南自然科学",
      "url": "https://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=6&fid=3"
    },
    {
      "index": 3,
      "title": "花津校区图书馆3楼北社科一",
      "url": "https://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=5&fid=4"
    },
    {
      "index": 4,
      "title": "花津校区图书馆4楼北社科三",
      "url": "https://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=3&fid=5"
    },
    {
      "index": 5,
      "title": "花津校区图书馆4楼南社科二",
      "url": "https://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=4&fid=6"
    },
    {
      "index": 6,
      "title": "花津校区图书馆3楼公共区域东",
      "url": "https://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=13&fid=9"
    },
    {
      "index": 7,
      "title": "花津校区图书馆3楼公共区域西",
      "url": "https://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=14&fid=9"
    },
    {
      "index": 8,
      "title": "花津校区图书馆4楼公共区域东",
      "url": "https://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=15&fid=10"
    },
    {
      "index": 9,
      "title": "花津校区图书馆4楼公共区域西",
      "url": "https://libzwxt.ahnu.edu.cn/SeatWx/Room.aspx?rid=16&fid=10"
    }
  ]
}

```
