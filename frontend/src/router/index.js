import { createRouter, createWebHistory } from 'vue-router';
import { routes } from './routes';
import { registerAuthGuard } from './guards';

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

registerAuthGuard(router);

export default router;
