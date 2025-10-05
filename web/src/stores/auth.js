import { defineStore } from 'pinia';

const TOKEN_KEY = 'merch_app_token';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem(TOKEN_KEY) || '',
    userEmail: localStorage.getItem('merch_app_email') || '',
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token),
  },
  actions: {
    setToken(token, email) {
      this.token = token;
      localStorage.setItem(TOKEN_KEY, token);
      if (email) {
        this.userEmail = email;
        localStorage.setItem('merch_app_email', email);
      }
    },
    logout() {
      this.token = '';
      this.userEmail = '';
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem('merch_app_email');
    },
  },
});
