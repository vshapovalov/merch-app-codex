<template>
  <div class="layout" :class="{ 'drawer-visible': drawerVisible }">
    <Drawer
      v-model:visible="drawerVisible"
      class="app-drawer"
      position="left"
      :modal="false"
      :dismissable="false"
      :show-close-icon="false"
      :block-scroll="false"
      append-to="self"
    >
      <template #container>
        <div class="drawer-surface">
          <div class="drawer-user" v-if="userDisplayName">
            <span class="drawer-user-name">{{ userDisplayName }}</span>
          </div>
          <nav class="drawer-menu">
            <RouterLink
              v-for="item in menuItems"
              :key="item.routeName"
              :to="{ name: item.routeName }"
              class="drawer-menu-item"
              :class="{ active: route.name === item.routeName }"
              @click="handleMenuClick"
            >
              <i :class="['pi', item.icon, 'drawer-menu-icon']" aria-hidden="true" />
              <span class="drawer-menu-label">{{ item.label }}</span>
            </RouterLink>
          </nav>
          <div class="drawer-footer">
            <Button
              label="Выход"
              icon="pi pi-sign-out"
              severity="secondary"
              text
              class="drawer-logout"
              @click="handleLogout"
            />
          </div>
        </div>
      </template>
    </Drawer>

    <div class="content-area">
      <Toolbar class="topbar">
        <template #start>
          <div class="topbar-left">
            <Button
              icon="pi pi-bars"
              text
              severity="secondary"
              class="toggle-button"
              @click="toggleDrawer"
            />
            <span class="page-title">{{ currentPageTitle }}</span>
          </div>
        </template>
        <template #end>
          <div class="topbar-right">
            <span v-if="auth.userEmail" class="user-email">{{ auth.userEmail }}</span>
            <Button label="Выход" severity="secondary" text @click="handleLogout" />
          </div>
        </template>
      </Toolbar>
      <main class="page-content">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();
const drawerVisible = ref(true);

const menuItems = [
  { label: 'Пользователи', icon: 'pi-users', routeName: 'users' },
  { label: 'Компании', icon: 'pi-building', routeName: 'companies' },
  { label: 'Торговые точки', icon: 'pi-map-marker', routeName: 'retail-points' },
  { label: 'Бренды', icon: 'pi-tag', routeName: 'brands' },
  { label: 'Категории', icon: 'pi-sitemap', routeName: 'categories' },
  { label: 'Продукты', icon: 'pi-box', routeName: 'products' },
  { label: 'Визиты', icon: 'pi-calendar', routeName: 'visits' },
  { label: 'Позиции визита', icon: 'pi-list', routeName: 'visit-items' },
  { label: 'Отчёты', icon: 'pi-chart-bar', routeName: 'reports-companies' },
];

const currentPageTitle = computed(() => {
  const activeItem = menuItems.find((item) => item.routeName === route.name);
  return activeItem ? activeItem.label : '';
});

const userDisplayName = computed(() => auth.userEmail ?? '');

const isSmallViewport = () => typeof window !== 'undefined' && window.innerWidth < 992;

onMounted(() => {
  if (isSmallViewport()) {
    drawerVisible.value = false;
  }
});

const toggleDrawer = () => {
  drawerVisible.value = !drawerVisible.value;
};

const handleMenuClick = () => {
  if (isSmallViewport()) {
    drawerVisible.value = false;
  }
};

const handleLogout = () => {
  drawerVisible.value = false;
  auth.logout();
  router.push({ name: 'login' });
};
</script>

<style scoped>
.layout {
  position: relative;
  min-height: 100vh;
  background: var(--surface-ground);
  display: flex;
  flex-direction: column;
}

.layout.drawer-visible .content-area {
  margin-left: 18rem;
}

@media (max-width: 991px) {
  .layout.drawer-visible .content-area {
    margin-left: 0;
  }
}

:deep(.app-drawer.p-drawer) {
  width: 18rem;
}

:deep(.app-drawer.p-drawer .p-drawer-content) {
  padding: 0;
  height: 100%;
}

.drawer-surface {
  background: var(--surface-card);
  border-right: 1px solid var(--surface-border);
  height: 100%;
  display: flex;
  flex-direction: column;
}

.drawer-user {
  padding: 1.5rem 1.5rem 1rem;
  border-bottom: 1px solid var(--surface-border);
}

.drawer-user-name {
  font-weight: 600;
  font-size: 1.05rem;
  color: var(--text-color);
}

.drawer-menu {
  display: flex;
  flex-direction: column;
  padding: 1rem 0.75rem;
  gap: 0.25rem;
  flex: 1;
}

.drawer-menu-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 0.85rem;
  border-radius: 0.75rem;
  text-decoration: none;
  color: var(--text-color);
  transition: background-color 0.2s ease, color 0.2s ease;
}

.drawer-menu-item:hover {
  background: color-mix(in srgb, var(--primary-color) 12%, transparent);
  color: var(--primary-color);
}

.drawer-menu-item.active {
  background: color-mix(in srgb, var(--primary-color) 18%, transparent);
  color: var(--primary-color);
  font-weight: 600;
}

.drawer-menu-icon {
  font-size: 1.1rem;
}

.drawer-menu-label {
  white-space: nowrap;
}

.drawer-footer {
  padding: 1rem 0.75rem 1.5rem;
  border-top: 1px solid var(--surface-border);
}

.drawer-logout {
  width: 100%;
  justify-content: flex-start;
}

.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  transition: margin-left 0.2s ease;
}

.topbar {
  border-radius: 0;
  border: none;
  background: var(--surface-card);
  border-bottom: 1px solid var(--surface-border);
  padding: 0 1.5rem;
}

.topbar-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.toggle-button {
  min-width: auto;
}

.page-title {
  font-weight: 600;
  font-size: 1.15rem;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.user-email {
  font-weight: 500;
}

.page-content {
  flex: 1;
  padding: 1.5rem;
}
</style>
