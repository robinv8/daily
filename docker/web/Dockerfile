FROM node:16-alpine as build

WORKDIR /app

COPY ui/package.json ui/pnpm-lock.yaml ./

RUN corepack enable && corepack prepare pnpm@latest --activate && pnpm install

COPY ui .

RUN pnpm run build


FROM node:16-alpine

WORKDIR /app

COPY --from=build /app/dist .

RUN npm init -y && npm install astro @astrojs/node @astrojs/react react react-dom

EXPOSE 3000
CMD [ "node", "server/entry.mjs" ]