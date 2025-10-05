import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const AppLayout = () => import('../components/layouts/AppLayout.vue');
const LoginView = () => import('../views/LoginView.vue');
const UsersView = () => import('../views/UsersView.vue');
const CompaniesView = () => import('../views/CompaniesView.vue');
const RetailPointsView = () => import('../views/RetailPointsView.vue');
const BrandsView = () => import('../views/BrandsView.vue');
const CategoriesView = () => import('../views/CategoriesView.vue');
const ProductsView = () => import('../views/ProductsView.vue');
const VisitsView = () => import('../views/VisitsView.vue');
const VisitItemsView = () => import('../views/VisitItemsView.vue');
const ReportsView = () => import('../views/ReportsView.vue');

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { public: true },
    },
    {
      path: '/',
      component: AppLayout,
      children: [
        { path: '', redirect: '/users' },
        { path: 'users', name: 'users', component: UsersView },
        { path: 'companies', name: 'companies', component: CompaniesView },
        { path: 'retail-points', name: 'retail-points', component: RetailPointsView },
        { path: 'brands', name: 'brands', component: BrandsView },
        { path: 'categories', name: 'categories', component: CategoriesView },
        { path: 'products', name: 'products', component: ProductsView },
        { path: 'visits', name: 'visits', component: VisitsView },
        { path: 'visit-items', name: 'visit-items', component: VisitItemsView },
        { path: 'reports/companies', name: 'reports-companies', component: ReportsView },
      ],
    },
  ],
});

router.beforeEach((to, from, next) => {
  const auth = useAuthStore();
  if (!to.meta.public && !auth.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } });
    return;
  }
  if (to.name === 'login' && auth.isAuthenticated) {
    next({ name: 'users' });
    return;
  }
  next();
});

export default router;
