const { awscdk, TextFile } = require("projen");

const project = new awscdk.AwsCdkTypeScriptApp({
  name: "mobworx-streamer",
  repository: "https://github.com/sehovizko/mobwrox-streamer.git",
  keywords: ["hls", "streaming", "golang"],
  deps: [
    "@aws-cdk/aws-lambda-go-alpha",
    "@aws-cdk/aws-apigatewayv2-alpha",
    "@aws-cdk/aws-apigatewayv2-integrations-alpha",
  ],

  prettier: true,
  release: true,
  releaseBranches: {
    dev: {},
  },
  defaultReleaseBranch: "main",
  workflowNodeVersion: "18",
  cdkVersion: "2.65.0",

  renovatebot: true,
  autoApproveUpgrades: true,
  autoApproveOptions: {
    allowedUsernames: ["sehovizko"],
  },
  license: false,

  lambdaOptions: {
    runtime: awscdk.LambdaRuntime.NODEJS_18_X,
  },
});

new TextFile(project, ".nvmrc", {
  marker: false,
  lines: ["v18"],
});

new TextFile(project, ".editorconfig", {
  marker: false,
  lines: `root = true
[*]
charset = utf-8
insert_final_newline = true
trim_trailing_whitespace = true
indent_size = 2
tab_width = 2
`.split("\n"),
});

project.synth();
