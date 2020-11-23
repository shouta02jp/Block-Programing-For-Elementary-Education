module.exports = {
  devServer: {
    host: 'localhost',
    port: 8888,
    disableHostCheck: true,
    proxy: {
      "/api/": {
        target: "http://localhost:8080/",
      }
    }
  },
  transpileDependencies: [
    'vuetify'
  ]
}