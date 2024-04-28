import { AstroSite, StackContext, use } from 'sst/constructs'
import { API } from './MyStack'

export function Astro({ stack }: StackContext) {
    const api = use(API)

    const site = new AstroSite(stack, 'AstroSite', {
        path: 'packages/web',
        environment: {
            API_URL: api.url
        }
    })

    stack.addOutputs({
        SiteUrl: site.url,
    })
}
