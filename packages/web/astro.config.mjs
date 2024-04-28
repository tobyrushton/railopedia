import { defineConfig } from 'astro/config'
import tailwind from '@astrojs/tailwind'
import svelte from '@astrojs/svelte'
import aws from 'astro-sst'

// https://astro.build/config
export default defineConfig({
    integrations: [tailwind(), svelte()],
    output: 'hybrid',
    adapter: aws(),
})
