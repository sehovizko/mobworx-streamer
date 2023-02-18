import { NestedStack } from "aws-cdk-lib";
import { Vpc } from "aws-cdk-lib/aws-ec2";

export class VpcNestedStack extends NestedStack {
  vpc = new Vpc(this, "VPC");
}
