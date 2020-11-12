module.exports = {
    outputDir: "../public",
    assetsDir: "static",
    pages: {
        index: {
            entry: "src/main.js",
            template: "public/index.html",
            filename: "index.html",
            title: "首页",
        }
    },
    chainWebpack: config => {
        config.performance.hints = false
        config.optimization.delete('splitChunks')
    },
    devServer: {
        disableHostCheck: true
    }
};