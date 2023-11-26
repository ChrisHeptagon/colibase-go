import { defineConfig } from 'astro/config';
import node from '@astrojs/node';
import react from "@astrojs/react";

import svelte from "@astrojs/svelte";

// https://astro.build/config
export default defineConfig({
  output: "server",
  vite: {
    ssr: {
      external: ['node:http']
    }
  },
  adapter: node({
    mode: "standalone"
  }),
  integrations: [react({
    experimentalReactChildren: true
  }), svelte()]
});