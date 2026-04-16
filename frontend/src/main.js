import { createApp } from 'vue';
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import 'element-plus/theme-chalk/index.css';

import App from './App.vue';
import router from './router';
import store from './store';
import { applySiteMeta } from './services/site-meta';
import './styles/index.css';

const app = createApp(App);

app.use(store);
app.use(router);
app.use(ElementPlus);

app.mount('#app');

store.dispatch('fetchSiteSettings')
  .then(() => {
    applySiteMeta(store.getters.siteSettings);
    store.watch(
      (state) => state.siteSettings,
      (newSettings) => {
        applySiteMeta(newSettings);
      }
    );
  })
  .catch(() => {
    applySiteMeta({
      title: 'ZNav',
      icon: 'https://img1.baidu.com/it/u=1217061905,2277984247&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=500',
    });
  });
