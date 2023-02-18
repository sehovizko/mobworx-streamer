import { join } from "path";
import { HttpApi, HttpMethod } from "@aws-cdk/aws-apigatewayv2-alpha";
import { AddRoutesOptions } from "@aws-cdk/aws-apigatewayv2-alpha/lib/http/api";
import { HttpLambdaIntegration } from "@aws-cdk/aws-apigatewayv2-integrations-alpha";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import { NestedStack } from "aws-cdk-lib";
import { IVpc } from "aws-cdk-lib/aws-ec2";
import { NestedStackProps } from "aws-cdk-lib/core/lib/nested-stack";
import { Construct } from "constructs";

export interface EndpointNestedStackProps extends NestedStackProps {
  vpc: IVpc;
}

export class EndpointNestedStack extends NestedStack {
  queryAdminEventsLambda = new GoFunction(this, "QueryAdminEventsLambda", {
    entry: join(__dirname, "adminevents", "query-admin-events.go"),
    vpc: this.props?.vpc,
  });

  queryAdminRoomsLambda = new GoFunction(this, "QueryAdminRoomsLambda", {
    entry: join(__dirname, "adminevents", "query-admin-rooms.go"),
    vpc: this.props?.vpc,
  });

  queryMunitStatsLambda = new GoFunction(this, "MStatsLambda", {
    entry: join(__dirname, "mstats", "query-munit-stats.go"),
    vpc: this.props?.vpc,
  });

  queryRoomParticipantsLambda = new GoFunction(this, "roomLambda", {
    entry: join(__dirname, "room", "query-participants.go"),
    vpc: this.props?.vpc,
  });

  api = new HttpApi(this, "EndpointApi");

  constructor(
    scope: Construct,
    id: string,
    private readonly props?: EndpointNestedStackProps
  ) {
    super(scope, id, props);

    [
      {
        path: "/v1/adminEvents/{eventType}",
        methods: [HttpMethod.GET],
        integration: new HttpLambdaIntegration(
          "queryAdminEvents",
          this.queryAdminEventsLambda
        ),
      },
      {
        path: "/v1/adminRooms",
        methods: [HttpMethod.GET],
        integration: new HttpLambdaIntegration(
          "queryAdminRooms",
          this.queryAdminRoomsLambda
        ),
      },
      {
        path: "/v1/mststats",
        methods: [HttpMethod.GET],
        integration: new HttpLambdaIntegration(
          "queryMunitStats",
          this.queryMunitStatsLambda
        ),
      },
      {
        path: "/v1/participants/{roomId}",
        methods: [HttpMethod.GET],
        integration: new HttpLambdaIntegration(
          "queryRoomParticipants",
          this.queryRoomParticipantsLambda
        ),
      },
    ].forEach((route: AddRoutesOptions) => this.api.addRoutes(route));
  }
}
