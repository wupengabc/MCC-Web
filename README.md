
### 项目说明

该项目使用 Beego 框架 + bootstrap V5 搭建，进程守护端暂不开放源码，以下是使用步骤：

1. **克隆项目**
   ```bash
   git clone https://git.wupeng1.top/2508932142/MCC-Web.git
   cd MCC-Web
   ```

2. **安装依赖**
   使用 `go mod` 来安装项目所需的依赖包。
   ```bash
   go mod tidy
   ```

3. **配置端口**
   在 `conf/app.conf` 文件中配置端口。

4. **运行项目**
   在项目根目录下执行以下命令来启动项目：
   ```bash
   bee run
   ```

5. **访问项目**
   项目启动后，可以在浏览器中通过 `http://localhost:端口` 访问。

6. **反馈问题**
   如发现问题或者有什么建议的功能请点击 [问题反馈](https://gitee.com/sg250/MCC-Web/issues)

### TODO
- Bot启动自动执行某些指令,
- 完善Bot分享
- 完善权限升级
- 后台管理界面的设计


### V1.3 更新日志
- 完善账户信息界面
- 支持修改用户名, 但修改前需手动删除所有Bot
- 支持修改密码
- 修复部分Bug, 提高稳定性

### V1.2 更新日志

- 完善搜索功能, 现可使用搜索功能搜索列表内的机器人
- 添加更新日志弹窗和打开按钮, 方便用户获取更新内容
- 完善快捷指令功能, 现可增删修改快捷指令
- 由于设计缺陷, 同一用户只能创建一个相同id的bot

### V1.1 更新日志

- 完善Bot配置修改, 可在配置修改里删除Bot
- 目前发现1.20.4 1.19.4版本不稳定 暂时不建议使用这两个版本 建议在跨版本服务器使用1.18.2
- 现在支持SRV解析的域名, 使用服主提供的不带端口的链接时可不输入端口
- 管理面板支持音乐播放功能
- 修复部分bug, 提高稳定性

### V1.0 更新日志

- 支持添加Bot,但是还暂未支持Bot的配置修改和Bot的删除
- 快捷指令UI设计完成,实际功能还未实现
- 添加Bot搜索框,但暂不支持使用
- 暂不支持账户修改以及权限升级
