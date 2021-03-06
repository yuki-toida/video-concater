const path = require('path');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const { VueLoaderPlugin } = require('vue-loader');

module.exports = (_, argv) => ({
  entry: './src/index.js',
  output: {
    filename: 'app.js',
    path: argv.mode == 'development'
      ? path.resolve(__dirname, '../static/js')
      : path.resolve(__dirname, './_build')
  },
  devtool: 'source-map',
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: 'babel-loader',
        options: {
          presets: [['env', { modules: false }]]
        }
      },
      {
        // enforce: 'pre'を指定することによって
        // enforce: 'pre'がついていないローダーより早く処理が実行される
        // 今回はbabel-loaderで変換する前にコードを検証したいため、指定が必要
        enforce: 'pre',
        test: /\.js$/,
        exclude: /node_modules/,
        loader: 'eslint-loader',
        options: {
          fix: true,
          failOnError: true,
        },
      },
      {
        test: /\.vue$/,
        loader: 'vue-loader'
      },
    ]
  },
  resolve: {
    alias: {
      'vue$': 'vue/dist/vue.esm.js'
    }
  },
  plugins: [
    new VueLoaderPlugin(),
    new CopyWebpackPlugin([
      { from: 'static', to: path.resolve(__dirname, '../static') }
    ])
  ],  
});