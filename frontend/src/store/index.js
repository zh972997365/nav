import { createStore } from 'vuex';
import auth from './modules/auth';
import site from './modules/site';

export default createStore({
  state: {
    token: auth.state().token,
    user: auth.state().user,
    siteSettings: site.state().siteSettings,
  },
  mutations: {
    setToken(state, token) {
      auth.mutations.setToken(state, token);
    },
    setUser(state, user) {
      auth.mutations.setUser(state, user);
    },
    clearToken(state) {
      auth.mutations.clearToken(state);
    },
    setSiteSettings(state, settings) {
      site.mutations.setSiteSettings(state, settings);
    },
    updateSiteSetting(state, payload) {
      site.mutations.updateSiteSetting(state, payload);
    },
  },
  actions: {
    fetchUserData(context) {
      return auth.actions.fetchUserData(context);
    },
    fetchSiteSettings(context) {
      return site.actions.fetchSiteSettings(context);
    },
    updateSiteSettings(context, payload) {
      return site.actions.updateSiteSettings(context, payload);
    },
    logout(context) {
      return auth.actions.logout(context);
    },
  },
  getters: {
    isAuthenticated: (state) => auth.getters.isAuthenticated(state),
    userData: (state) => auth.getters.userData(state),
    siteSettings: (state) => site.getters.siteSettings(state),
    username: (state) => auth.getters.username(state),
  },
});
