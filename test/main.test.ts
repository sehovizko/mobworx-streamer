import { App } from "aws-cdk-lib";
import { Template } from "aws-cdk-lib/assertions";
import { StreamerStack } from "../src/main";

test("Snapshot", () => {
  const app = new App();
  const stack = new StreamerStack(app, "test");

  const template = Template.fromStack(stack);
  expect(template.toJSON()).toMatchSnapshot();
});
