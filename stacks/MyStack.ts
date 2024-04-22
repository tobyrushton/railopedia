import { StackContext, Api } from "sst/constructs";

export function API({ stack }: StackContext) {
  stack.setDefaultFunctionProps({
    runtime: "go",
  })

  const api = new Api(stack, "api", {
    routes: {
      "GET /search": "packages/functions/src/endpoints/search.go",
    },
  })

  stack.addOutputs({
    ApiEndpoint: api.url,
  })
}
