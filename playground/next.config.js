/** @type {import('next').NextConfig} */

module.exports = {
    reactStrictMode: true,
    webpack: function (config, options) {
        config.experiments = { layers: true, asyncWebAssembly: true };
        return config;
    }
}