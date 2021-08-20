# 创建项目

```shell script
# Create a go module for CLI.
go mod init democtl

# Get Cobra library
go get -u github.com/spf13/cobra/cobra

# Create a bare minimum skeleton
cobra init --pkg-name democtl

# 创建create命令
cobra add create

# 创建card命令
cobra add card
```
