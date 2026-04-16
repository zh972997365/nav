import http from '../../services/http/client';

const state = () => ({
  token: localStorage.getItem('token') || null,
  user: {
    username: '',
  },
});

const mutations = {
  setToken(currentState, token) {
    currentState.token = token;
    localStorage.setItem('token', token);
  },
  setUser(currentState, user) {
    currentState.user = user;
  },
  clearToken(currentState) {
    currentState.token = null;
    currentState.user = {};
    localStorage.removeItem('token');
  },
};

const actions = {
  fetchUserData({ commit, state, dispatch }) {
    return new Promise((resolve, reject) => {
      if (!state.token) {
        reject(new Error('No token available'));
        return;
      }

      http.get('/user')
        .then((response) => {
          commit('setUser', response.data.user);
          resolve(response.data.user);
        })
        .catch((error) => {
          if (error.response && error.response.status === 401) {
            dispatch('logout');
          }
          reject(error);
        });
    });
  },
  logout({ commit }) {
    commit('clearToken');
  },
};

const getters = {
  isAuthenticated: (currentState) => !!currentState.token,
  userData: (currentState) => currentState.user,
  username: (currentState) => currentState.user.username || '管理员',
};

export default {
  namespaced: false,
  state,
  mutations,
  actions,
  getters,
};
