module.exports = {
  apps: [{
    name: 'pichub-api',
    script: './pichub-api',  // 直接使用当前目录下的二进制文件
    cwd: './dist',  // 设置工作目录为 dist
    exec_mode: 'fork',  // Go 程序是二进制的，使用 fork 模式
    instances: 1,  // 实例数量
    autorestart: true,  // 自动重启
    watch: false,  // 不需要监视文件变化
    max_memory_restart: '1G',  // 超过内存限制时重启
  }]
};
