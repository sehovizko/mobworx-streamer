import { join } from "path";
import { HttpApi, HttpMethod } from "@aws-cdk/aws-apigatewayv2-alpha";
import { AddRoutesOptions } from "@aws-cdk/aws-apigatewayv2-alpha/lib/http/api";
import { HttpLambdaIntegration } from "@aws-cdk/aws-apigatewayv2-integrations-alpha";
import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import { NestedStack, NestedStackProps } from "aws-cdk-lib";
import { Vpc } from "aws-cdk-lib/aws-ec2";
import { CfnCacheCluster, CfnSubnetGroup } from "aws-cdk-lib/aws-elasticache";
import { Construct } from "constructs";

export interface EndpointNestedStackProps extends NestedStackProps {
  vpc: Vpc;
}

export class EndpointNestedStack extends NestedStack {
  redisSubnetGroup = new CfnSubnetGroup(this, "RedisSubnetGroup", {
    subnetIds: this.props.vpc.privateSubnets.map((subnet) => subnet.subnetId),
    description: "Redis subnet group",
  });

  redisCluster = new CfnCacheCluster(this, "RedisCluster", {
    engine: "redis",
    cacheNodeType: "cache.t3.micro",
    numCacheNodes: 1,
    vpcSecurityGroupIds: [this.props.vpc.vpcDefaultSecurityGroup.toString()],
    cacheSubnetGroupName: this.redisSubnetGroup.cacheSubnetGroupName,
  });

  queryAdminEventsLambda = new GoFunction(this, "QueryAdminEventsLambda", {
    entry: join(__dirname, "adminevents", "query-admin-events.go"),
    vpc: this.props.vpc,
    environment: {},
  });

  queryAdminRoomsLambda = new GoFunction(this, "QueryAdminRoomsLambda", {
    entry: join(__dirname, "adminevents", "query-admin-rooms.go"),
    vpc: this.props.vpc,
  });

  queryMunitStatsLambda = new GoFunction(this, "MStatsLambda", {
    entry: join(__dirname, "mstats", "query-munit-stats.go"),
    vpc: this.props.vpc,
    environment: {
      REDIS_ADDRESS: this.redisCluster.attrRedisEndpointAddress,
    },
  });

  queryRoomParticipantsLambda = new GoFunction(this, "roomLambda", {
    entry: join(__dirname, "room", "query-participants.go"),
    vpc: this.props?.vpc,
  });

  api = new HttpApi(this, "EndpointApi");

  constructor(
    scope: Construct,
    id: string,
    private readonly props: EndpointNestedStackProps
  ) {
    super(scope, id, props);

    this.redisCluster.addDependency(this.redisSubnetGroup);

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
