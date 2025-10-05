<template>
  <div class="login-page flex align-items-center justify-content-center">
    <Card class="login-card">
      <template #title>
        <div class="flex align-items-center gap-2">
          <i class="pi pi-lock" />
          <span>Вход в систему</span>
        </div>
      </template>
      <template #content>
        <form class="flex flex-column gap-3" @submit.prevent="onSubmit">
          <span class="p-float-label">
            <InputText id="email" v-model="email" type="email" required autofocus />
            <label for="email">Email</label>
          </span>
          <span class="p-float-label">
            <Password id="password" v-model="password" :feedback="false" toggleMask required />
            <label for="password">Пароль</label>
          </span>
          <Button type="submit" label="Войти" :loading="loading" />
        </form>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useToast } from 'primevue/usetoast';
import api from '../services/api';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const route = useRoute();
const toast = useToast();
const auth = useAuthStore();

const email = ref(auth.userEmail || '');
const password = ref('');
const loading = ref(false);

const onSubmit = async () => {
  loading.value = true;
  try {
    const { data } = await api.post('/auth/login', {
      email: email.value,
      password: password.value,
    });
    auth.setToken(data.token, email.value);
    const redirect = route.query.redirect || '/users';
    router.push(redirect);
  } catch (error) {
    console.error(error);
    toast.add({ severity: 'error', summary: 'Ошибка авторизации', detail: 'Проверьте email и пароль', life: 3000 });
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: var(--surface-ground);
}

.login-card {
  width: min(28rem, 100%);
}
</style>
