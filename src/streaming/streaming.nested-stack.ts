import { NestedStack } from "aws-cdk-lib";
import { Construct } from "constructs";
import { NestedStackProps } from "aws-cdk-lib/core/lib/nested-stack";
import { Vpc } from "aws-cdk-lib/aws-ec2";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import { join } from "path";

export interface StreamingNestedStackProps extends NestedStackProps {
  vpc: Vpc;
}

export class StreamingNestedStack extends NestedStack {
  updateRenditionLambda = new GoFunction(this, "UpdateRendition", {
    entry: join(__dirname, "update-rendition.go"),
    vpc: this.props.vpc,
  });

  constructor(
    scope: Construct,
    id: string,
    private readonly props: StreamingNestedStackProps
  ) {
    super(scope, id, props);
  }
}
