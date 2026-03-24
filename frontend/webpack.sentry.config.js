const SentryWebpackPlugin = require("@sentry/webpack-plugin");

module.exports = {
  // ... other webpack config
  devtool: "source-map",
  plugins: [
    new SentryWebpackPlugin({
      org: process.env.SENTRY_ORG || "phenotype",
      project: process.env.SENTRY_PROJECT || "phenotype-web",
      authToken: process.env.SENTRY_AUTH_TOKEN,
      release: process.env.SENTRY_RELEASE || "production",
      dist: "production",
      include: "./dist",
      ignore: ["node_modules", "test"],
      sourceMapReference: true,
      stripPrefix: ["webpack://"],
      urlPrefix: "~/",
    }),
  ],
};

// CI/CD Integration Example (.gitlab-ci.yml):
/*
stages:
  - deploy

deploy:
  stage: deploy
  script:
    - npm install
    - npm run build
    - node scripts/upload-sourcemaps.js
  only:
    - main
    - production
  environment:
    name: production
*/

// Environment variables needed:
// SENTRY_ORG=phenotype
// SENTRY_PROJECT=phenotype-web
// SENTRY_AUTH_TOKEN=<your-auth-token>
// SENTRY_RELEASE=$CI_COMMIT_SHA
