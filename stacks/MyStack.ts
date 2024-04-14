import { StackContext, Api } from "sst/constructs";

export function API({ stack }: StackContext) {
  stack.setDefaultFunctionProps({
    runtime: "go",
  })

  const api = new Api(stack, "api", {
    routes: {
      "GET /": "packages/functions/src/lambda.go",
    },
  })

  stack.addOutputs({
    ApiEndpoint: api.url,
  })
}
