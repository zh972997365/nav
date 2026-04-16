export function registerAuthGuard(router) {
  router.beforeEach((to, from, next) => {
    const loggedIn = !!localStorage.getItem('token');
    if (to.matched.some((record) => record.meta.requiresAuth) && !loggedIn) {
      next('/login');
      return;
    }

    setTimeout(() => {
      next();
    }, 200);
  });
}
