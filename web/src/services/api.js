import axios from 'axios';
import router from '../router';
import { useAuthStore } from '../stores/auth';

const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL,
});

api.interceptors.request.use((config) => {
  const auth = useAuthStore();
  if (auth.token) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${auth.token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status;
    if (status === 401) {
      const auth = useAuthStore();
      auth.logout();
      if (router.currentRoute.value.name !== 'login') {
        router.push({ name: 'login', query: { redirect: router.currentRoute.value.fullPath } });
      }
    }
    return Promise.reject(error);
  }
);

export default api;
