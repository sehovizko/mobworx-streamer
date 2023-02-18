import { App, Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import { EndpointNestedStack } from "./endpoint/endpoint.nested-stack";
import { VpcNestedStack } from "./vpc.nested-stack";

export class StreamerStack extends Stack {
  constructor(scope: Construct, id: string, props: StackProps = {}) {
    super(scope, id, props);
    const vpcNestedStack = new VpcNestedStack(this, "VPC");
    new EndpointNestedStack(this, "EndpointNestedStack", {
      vpc: vpcNestedStack.vpc,
    });
  }
}

// for development, use account/region from cdk cli
const devEnv = {
  account: process.env.CDK_DEFAULT_ACCOUNT,
  region: process.env.CDK_DEFAULT_REGION,
};

const app = new App();

new StreamerStack(app, "mobworx-streamer-dev", { env: devEnv });
// new MyStack(app, 'mobworx-streamer-prod', { env: prodEnv });

app.synth();
