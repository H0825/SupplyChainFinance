import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('scf_user') || 'null')
  }),
  getters: {
    role: (s) => (s.user?.role || 'guest').toLowerCase(),
    isLoggedIn: (s) => !!s.user
  },
  actions: {
    setUser(user) {
      this.user = user
      localStorage.setItem('scf_user', JSON.stringify(user))
    },
    logout() {
      this.user = null
      localStorage.removeItem('scf_user')
    }
  }
})
