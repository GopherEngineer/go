const path = require('path')
const TerserPlugin = require('terser-webpack-plugin')

module.exports = {
	resolve: {
		extensions: ['.tsx', '.ts', '.js'],
	},
	module: {
		rules: [
			{
				test: /\.(ts|tsx)$/,
				include: path.resolve(__dirname, 'src'),
				use: 'ts-loader'
			},
			{
				test: /\.css$/i,
				include: path.resolve(__dirname, 'src'),
				use: ['style-loader', 'css-loader', 'postcss-loader'],
			}
		]
	},
	devServer: {
		static: {
			directory: path.join(__dirname, 'public'),
		},
		compress: false,
		port: 4000,
	},
	output: {
		filename: 'index.js',
		path: path.resolve(__dirname, 'public'),
	},
	optimization: {
		minimize: true,
		minimizer: [
			new TerserPlugin({
				extractComments: false,
			}),
		],
	},
}