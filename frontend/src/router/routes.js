import AdminLayout from '../layouts/AdminLayout.vue';
import LoginPage from '../pages/auth/LoginPage.vue';
import DashboardPage from '../pages/admin/DashboardPage.vue';
import ApplicationsPage from '../pages/admin/ApplicationsPage.vue';
import MenusPage from '../pages/admin/MenusPage.vue';
import UsersPage from '../pages/admin/UsersPage.vue';
import SettingsPage from '../pages/admin/SettingsPage.vue';
import HomePage from '../pages/public/HomePage.vue';

export const routes = [
  {
    path: '/',
    component: HomePage,
  },
  {
    path: '/login',
    component: LoginPage,
  },
  {
    path: '/',
    component: AdminLayout,
    children: [
      { path: '/dashboard', component: DashboardPage, meta: { requiresAuth: true } },
      { path: '/apps', component: ApplicationsPage, meta: { requiresAuth: true } },
      { path: '/menus', component: MenusPage, meta: { requiresAuth: true } },
      { path: '/admin', component: UsersPage, meta: { requiresAuth: true } },
      { path: '/settings', component: SettingsPage, meta: { requiresAuth: true } },
    ],
  },
];
