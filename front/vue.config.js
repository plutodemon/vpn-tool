const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8181,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        pathRewrite:{'^/api':''},
        ws: true,
        changeOrigin: true
      },
    }
  }
})
