import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import tsconfigPaths from 'vite-tsconfig-paths';
import fs from 'fs';

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());

  return {
    plugins: [react(), tsconfigPaths()],
    server: {
      proxy: {
        '/api': {
          target: 'http://localhost:4001',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ''),
        },
      },
      https: {
        key: fs.readFileSync('./certs/localhost-key.pem'),
        cert: fs.readFileSync('./certs/localhost.pem'),
      },
      host: 'localhost',
    },
    css: {
      modules: {
        scopeBehaviour: 'local',
      },
    },
  }
});
