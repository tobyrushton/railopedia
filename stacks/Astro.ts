import { AstroSite, StackContext, use } from 'sst/constructs'
import { API } from './MyStack'

export function Astro({ stack }: StackContext) {
    const api = use(API)

    const site = new AstroSite(stack, 'AstroSite', {
        path: 'packages/web',
        customDomain: stack.stage === 'prod' ? 'railopedia.com' : undefined, 
        environment: {
            API_URL: api.url
        }
    })

    stack.addOutputs({
        SiteUrl: site.customDomainUrl || site.url,
    })
}
