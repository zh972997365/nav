<template>
  <div class="console-shell">
    <aside class="console-sidebar">
      <div class="console-sidebar__top">
        <div class="console-brand" @click="goToHome">
          <div class="console-brand__logo">
            <img v-if="isUrl(displayIcon)" :src="displayIcon" alt="brand" class="console-brand__logo-image" />
            <icon v-else :icon="displayIcon" width="20" height="20" color="white" />
          </div>
          <div class="console-brand__copy">
            <strong>{{ siteTitle }}</strong>
            <span>管理后台</span>
          </div>
        </div>
      </div>

      <div class="console-sidebar__body">
        <p class="console-sidebar__label">Navigation</p>
        <nav class="console-nav">
          <button
            v-for="item in navItems"
            :key="item.path"
            type="button"
            class="console-nav__item"
            :class="{ 'is-active': $route.path === item.path }"
            @click="goTo(item.path)"
          >
            <span class="console-nav__accent"></span>
            <span class="console-nav__icon">
              <icon :icon="item.icon" width="16" height="16" />
            </span>
            <span class="console-nav__text">{{ item.label }}</span>
          </button>
        </nav>
      </div>

      <div class="console-sidebar__bottom">
        <div class="console-account">
          <div class="console-account__avatar">{{ accountInitial }}</div>
          <div class="console-account__copy">
            <strong>{{ username }}</strong>
            <span>{{ userRole }}</span>
          </div>
        </div>

        <button type="button" class="console-logout" @click="logout">
          <icon icon="mdi:logout-variant" width="16" height="16" />
          <span>退出登录</span>
        </button>
      </div>
    </aside>

    <main class="console-main">
      <header class="console-topbar">
        <div>
          <p class="console-topbar__eyebrow">{{ currentPageInfo.eyebrow }}</p>
          <h1 class="console-topbar__title">{{ currentPageInfo.title }}</h1>
        </div>
        <div class="console-topbar__meta">
          <div class="console-chip">
            <icon icon="mdi:calendar-blank-outline" width="14" height="14" />
            <span>{{ formattedDate }}</span>
          </div>
        </div>
      </header>

      <section class="console-main__body">
        <router-view />
      </section>
    </main>
  </div>
</template>

<script>
import { Icon } from '@iconify/vue';
import { mapGetters } from 'vuex';
import { ADMIN_NAV_ITEMS } from '../constants/navigation';

export default {
  name: 'AdminLayout',
  components: {
    Icon,
  },
  computed: {
    ...mapGetters(['siteSettings']),
    navItems() {
      return ADMIN_NAV_ITEMS;
    },
    currentPageInfo() {
      return ADMIN_NAV_ITEMS.find((item) => item.path === this.$route.path) || {
        eyebrow: 'Admin',
        title: '管理后台',
      };
    },
    displayIcon() {
      return this.siteSettings.icon || 'solar:planet-outline';
    },
    siteTitle() {
      return this.siteSettings.title || 'ZNav';
    },
    username() {
      return (this.$store.state.user && this.$store.state.user.username) || '管理员';
    },
    userRole() {
      return this.$store.state.user && this.$store.state.user.is_admin ? '超级管理员' : '普通成员';
    },
    accountInitial() {
      return this.username ? this.username.slice(0, 1).toUpperCase() : 'A';
    },
    token() {
      return (this.$store.state && this.$store.state.token) || null;
    },
    formattedDate() {
      return new Intl.DateTimeFormat('zh-CN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      }).format(new Date());
    },
  },
  watch: {
    '$route.path'() {
      if (!this.token) {
        this.$router.push('/login');
        return;
      }
      document.title = `${this.currentPageInfo.title} | ${this.siteTitle}`;
    },
  },
  created() {
    this.$store.dispatch('fetchSiteSettings');
    if (this.token) {
      this.$store.dispatch('fetchUserData').catch(() => {});
    } else {
      this.$router.push('/login');
    }
  },
  methods: {
    isUrl(string) {
      return /^(http|https):\/\//.test(string || '');
    },
    goTo(path) {
      if (this.$route.path !== path) {
        this.$router.push(path);
      }
    },
    goToHome() {
      this.$router.push('/');
    },
    logout() {
      this.$store.dispatch('logout');
      this.$router.push('/login');
    },
  },
};
</script>

<style scoped>
.console-shell {
  box-sizing: border-box;
  display: grid;
  grid-template-columns: 288px minmax(0, 1fr);
  height: 100vh;
  overflow: hidden;
  background:
    radial-gradient(circle at top left, rgba(147, 197, 253, 0.08), transparent 24%),
    linear-gradient(180deg, #f4f9fb, #fafdfd 45%, #f3f8f9);
}

.console-sidebar {
  box-sizing: border-box;
  display: grid;
  grid-template-rows: auto minmax(0, 1fr) auto;
  height: 100vh;
  min-height: 0;
  padding: 18px;
  background:
    radial-gradient(circle at top left, rgba(125, 211, 252, 0.16), transparent 28%),
    linear-gradient(180deg, #1e293b, #172033 48%, #111827);
  color: #e2e8f0;
  overflow: hidden;
}

.console-sidebar__top {
  min-height: 0;
}

.console-brand {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 10px 10px 14px;
  border-radius: 22px;
  cursor: pointer;
}

.console-brand__logo {
  width: 46px;
  height: 46px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: linear-gradient(135deg, #7dd3fc, #38bdf8);
  box-shadow: 0 18px 32px rgba(56, 189, 248, 0.22);
  overflow: hidden;
}

.console-brand__logo-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.console-brand__copy {
  min-width: 0;
}

.console-brand__copy strong {
  display: block;
  font-size: 16px;
  color: #f8fafc;
}

.console-brand__copy span {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: rgba(226, 232, 240, 0.68);
}

.console-sidebar__body {
  min-height: 0;
  display: flex;
  flex-direction: column;
  margin-top: 8px;
  overflow: hidden;
}

.console-sidebar__label {
  margin: 0 8px 12px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: rgba(191, 219, 254, 0.56);
}

.console-nav {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 0;
  overflow: auto;
  padding-right: 4px;
}

.console-nav__item {
  position: relative;
  min-height: 48px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 0 12px;
  border: 1px solid transparent;
  border-radius: 18px;
  background: transparent;
  color: rgba(226, 232, 240, 0.78);
  cursor: pointer;
  transition: all 0.2s ease;
}

.console-nav__item:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #ffffff;
}

.console-nav__item.is-active {
  background: linear-gradient(135deg, rgba(125, 211, 252, 0.16), rgba(56, 189, 248, 0.16));
  border-color: rgba(125, 211, 252, 0.16);
  color: #ffffff;
}

.console-nav__accent {
  width: 3px;
  height: 20px;
  border-radius: 999px;
  background: transparent;
  transition: all 0.2s ease;
}

.console-nav__item.is-active .console-nav__accent {
  background: #7dd3fc;
  box-shadow: 0 0 0 6px rgba(125, 211, 252, 0.08);
}

.console-nav__icon {
  width: 30px;
  height: 30px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.06);
}

.console-nav__text {
  min-width: 0;
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

.console-sidebar__bottom {
  min-height: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-top: 14px;
  margin-top: 12px;
  border-top: 1px solid rgba(148, 163, 184, 0.14);
  background: linear-gradient(180deg, rgba(15, 23, 42, 0), rgba(15, 23, 42, 0.92) 18%);
}

.console-account {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(148, 163, 184, 0.14);
}

.console-account__avatar {
  width: 40px;
  height: 40px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: rgba(125, 211, 252, 0.16);
  color: #e0f2fe;
  font-size: 13px;
  font-weight: 700;
}

.console-account__copy {
  min-width: 0;
  flex: 1;
}

.console-account__copy strong {
  display: block;
  font-size: 13px;
  color: #f8fafc;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.console-account__copy span {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: rgba(226, 232, 240, 0.68);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.console-logout {
  width: 100%;
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 0 14px;
  border: 1px solid rgba(125, 211, 252, 0.16);
  border-radius: 16px;
  background: rgba(125, 211, 252, 0.1);
  color: #f0f9ff;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s ease;
}

.console-logout:hover {
  background: rgba(125, 211, 252, 0.18);
  transform: translateY(-1px);
}

.console-main {
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  min-width: 0;
  height: 100vh;
  padding: 22px;
  overflow: hidden;
}

.console-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
}

.console-topbar__eyebrow {
  margin: 0 0 8px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #0f766e;
}

.console-topbar__title {
  margin: 0;
  font-size: 30px;
  line-height: 1.08;
  color: #0f172a;
}

.console-chip {
  min-height: 38px;
  padding: 0 14px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 999px;
  background: rgba(250, 254, 255, 0.9);
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.05);
  font-size: 12px;
  color: #6b7280;
  backdrop-filter: blur(14px);
}

.console-main__body {
  min-width: 0;
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding-right: 4px;
}

@media (max-width: 1080px) {
  .console-shell {
    grid-template-columns: 1fr;
    height: auto;
  }

  .console-sidebar {
    height: auto;
    grid-template-rows: auto auto auto;
  }

  .console-sidebar__body {
    overflow: visible;
  }

  .console-main {
    height: auto;
  }
}

@media (max-width: 768px) {
  .console-main {
    padding: 16px;
  }

  .console-topbar {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
