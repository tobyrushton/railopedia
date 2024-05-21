import { StackContext, Api } from "sst/constructs";

export function API({ stack }: StackContext) {
  const api = new Api(stack, "api", {
    customDomain: stack.stage === "prod" ? "api.railopedia.com": undefined,
    routes: {
      "GET /search": {
        function: {
          runtime: "container",
          timeout: 60,
          memorySize: '1 GB',
          handler: "packages/functions",
        }
      },
    },
  })

  stack.addOutputs({
    ApiEndpoint: api.customDomainUrl || api.url,
  })

  return api
}
