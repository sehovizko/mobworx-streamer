import { NestedStack } from "aws-cdk-lib";
import { Construct } from "constructs";
import { NestedStackProps } from "aws-cdk-lib/core/lib/nested-stack";
import { Vpc } from "aws-cdk-lib/aws-ec2";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import { join } from "path";
import { HttpApi, HttpMethod } from "@aws-cdk/aws-apigatewayv2-alpha";
import { HttpLambdaIntegration } from "@aws-cdk/aws-apigatewayv2-integrations-alpha";

export interface StreamingNestedStackProps extends NestedStackProps {
  vpc: Vpc;
  api: HttpApi;
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

    [
      {
        path: "/live/update/rendition",
        methods: [HttpMethod.POST],
        integration: new HttpLambdaIntegration(
          "updateRenditionHttp",
          this.updateRenditionLambda
        ),
      },
    ].forEach((route) => this.props.api.addRoutes(route));
  }
}
