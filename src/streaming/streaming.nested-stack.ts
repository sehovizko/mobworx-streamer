import { join } from "path";
import { HttpApi, HttpMethod } from "@aws-cdk/aws-apigatewayv2-alpha";
import { HttpLambdaIntegration } from "@aws-cdk/aws-apigatewayv2-integrations-alpha";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import { NestedStack } from "aws-cdk-lib";
import { Vpc } from "aws-cdk-lib/aws-ec2";
import { NestedStackProps } from "aws-cdk-lib/core/lib/nested-stack";
import { Construct } from "constructs";

export interface StreamingNestedStackProps extends NestedStackProps {
  vpc: Vpc;
  api: HttpApi;
}

export class StreamingNestedStack extends NestedStack {
  updatePartLambda = new GoFunction(this, "UpdatePart", {
    entry: join(__dirname, "update-part.go"),
    vpc: this.props.vpc,
  });

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
        path: "/live/update/part",
        methods: [HttpMethod.POST],
        integration: new HttpLambdaIntegration(
          "updateRenditionHttp",
          this.updatePartLambda
        ),
      },
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
