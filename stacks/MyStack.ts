import { StackContext, Api } from "sst/constructs";

export function API({ stack }: StackContext) {
  const api = new Api(stack, "api", {
    routes: {
      "GET /search": {
        function: {
          runtime: "go",
          timeout: 20,
          memorySize: 256,
          handler: "packages/functions/src/endpoints/search.go",
        }
      },
    },
  })

  stack.addOutputs({
    ApiEndpoint: api.url,
  })

  return api
}
