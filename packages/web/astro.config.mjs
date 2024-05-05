import { defineConfig } from 'astro/config'
import tailwind from '@astrojs/tailwind'
import svelte from '@astrojs/svelte'
import aws from 'astro-sst'

import simpleStackStream from 'simple-stack-stream'

// https://astro.build/config
export default defineConfig({
    integrations: [tailwind(), svelte(), simpleStackStream()],
    output: 'hybrid',
    adapter: aws(),
})
