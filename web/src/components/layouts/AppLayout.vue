<template>
  <div class="layout">
    <Menubar :model="menuItems">
      <template #start>
        <span class="app-title">Merch Admin</span>
      </template>
      <template #end>
        <div class="flex align-items-center gap-3">
          <span v-if="auth.userEmail" class="user-email">{{ auth.userEmail }}</span>
          <Button label="Выход" severity="secondary" text @click="logout" />
        </div>
      </template>
    </Menubar>
    <main class="p-4">
      <RouterView />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { RouterView, useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';

const auth = useAuthStore();
const router = useRouter();

const menuItems = computed(() => [
  { label: 'Пользователи', icon: 'pi pi-users', command: () => router.push({ name: 'users' }) },
  { label: 'Компании', icon: 'pi pi-building', command: () => router.push({ name: 'companies' }) },
  { label: 'Торговые точки', icon: 'pi pi-map-marker', command: () => router.push({ name: 'retail-points' }) },
  { label: 'Бренды', icon: 'pi pi-tag', command: () => router.push({ name: 'brands' }) },
  { label: 'Категории', icon: 'pi pi-sitemap', command: () => router.push({ name: 'categories' }) },
  { label: 'Продукты', icon: 'pi pi-box', command: () => router.push({ name: 'products' }) },
  { label: 'Визиты', icon: 'pi pi-calendar', command: () => router.push({ name: 'visits' }) },
  { label: 'Позиции визита', icon: 'pi pi-list', command: () => router.push({ name: 'visit-items' }) },
  { label: 'Отчёты', icon: 'pi pi-chart-bar', command: () => router.push({ name: 'reports-companies' }) },
]);

const logout = () => {
  auth.logout();
  router.push({ name: 'login' });
};
</script>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-title {
  font-weight: 600;
  font-size: 1.1rem;
}

main {
  flex: 1;
  background: var(--surface-ground);
}

.user-email {
  font-weight: 500;
}
</style>
