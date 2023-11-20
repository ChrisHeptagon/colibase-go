import { defineConfig } from 'astro/config';
import node from '@astrojs/node';

import react from "@astrojs/react";

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
  })]
});