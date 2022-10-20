import { defineConfig } from "astro/config";
import node from "@astrojs/node";
import tailwind from "@astrojs/tailwind";

// https://astro.build/config
export default defineConfig({
  output: "server",
  integrations: [tailwind()],
  adapter: node({
    mode: "standalone",
  }),
  server: {
    port: 3000,
    host: "0.0.0.0",
  },
});
