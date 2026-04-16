import http from '../../services/http/client';

const state = () => ({
  siteSettings: {
    id: null,
    title: '',
    icon: '',
    icon_color: '',
    footerInfo: '',
    recordNumber: '',
    image_host_token: '',
  },
});

const mutations = {
  setSiteSettings(currentState, settings) {
    currentState.siteSettings = settings;
  },
  updateSiteSetting(currentState, { key, value }) {
    currentState.siteSettings[key] = value;
  },
};

const actions = {
  fetchSiteSettings({ commit }) {
    return http.get('/settings')
      .then((response) => {
        commit('setSiteSettings', response.data);
        return response.data;
      })
      .catch((error) => {
        console.error('Failed to fetch site settings:', error);
        throw error;
      });
  },
  updateSiteSettings({ commit }, settings) {
    return http.put('/settings', settings)
      .then((response) => {
        commit('setSiteSettings', response.data);
        return response.data;
      })
      .catch((error) => {
        console.error('Failed to update site settings:', error);
        throw error;
      });
  },
};

const getters = {
  siteSettings: (currentState) => currentState.siteSettings,
};

export default {
  namespaced: false,
  state,
  mutations,
  actions,
  getters,
};
