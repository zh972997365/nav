<template>
  <div ref="topAnchor" class="home-shell">
    <header class="home-topbar">
      <button type="button" class="brand" @click="goToHome">
        <span class="brand__logo">
          <img v-if="isUrl(siteSettings.icon)" :src="siteSettings.icon" alt="logo" />
          <icon
            v-else
            :icon="siteSettings.icon || 'solar:planet-outline'"
            :color="siteSettings.icon_color || '#3f8f94'"
            width="22"
            height="22"
          />
        </span>
        <span class="brand__text">
          <strong>{{ siteSettings.title || 'ZNav' }}</strong>
          <small>个人导航工作台</small>
        </span>
      </button>

      <div class="topbar-actions">
        <div
          v-if="isLoggedIn"
          class="user-dropdown"
          @mouseenter="openUserMenu"
          @mouseleave="closeUserMenu"
        >
          <button type="button" class="user-trigger" @click="toggleUserMenu">
            <span class="user-trigger__avatar">{{ displayUsername.slice(0, 1).toUpperCase() }}</span>
            <span class="user-trigger__meta">
              <strong class="user-trigger__name">{{ displayUsername }}</strong>
              <small class="user-trigger__role">已登录用户</small>
            </span>
            <span class="user-trigger__arrow" :class="{ 'is-open': userMenuVisible }">
              <icon icon="solar:alt-arrow-down-outline" width="16" height="16" />
            </span>
          </button>

          <div v-if="userMenuVisible" class="user-menu">
            <button type="button" class="user-menu__item" @click="goToDashboard">后台管理</button>
            <button type="button" class="user-menu__item" @click="logout">退出登录</button>
          </div>
        </div>
        <button v-else type="button" class="soft-btn" @click="goToLogin">登录后台</button>
      </div>
    </header>

    <main class="home-main">
      <section class="toolbar-card">
        <div class="search-box">
          <icon icon="solar:magnifer-outline" width="18" height="18" />
          <input v-model.trim="searchQuery" type="text" placeholder="搜索应用名称、描述或分组" />
          <button v-if="searchQuery" type="button" @click="searchQuery = ''">清空</button>
        </div>

      </section>

      <section class="content-layout">
        <aside class="catalog">
          <div class="catalog__inner">
            <p class="catalog__label">目录导航</p>
            <button
              v-for="section in sectionSummaries"
              :key="`nav-${section.id}`"
              type="button"
              class="catalog__item"
              :class="{ 'is-active': activeMenu === String(section.id) }"
              @click="scrollToMenu(section.id)"
            >
              <span>{{ section.title }}</span>
              <strong>{{ section.count }}</strong>
            </button>
          </div>
        </aside>

        <div class="sections">
          <section
            v-for="section in displaySections"
            :id="`menu-${section.id}`"
            :key="section.id"
            class="section-card"
          >
            <div class="section-card__header">
              <div>
                <h2>{{ section.title }}</h2>
                <p>{{ section.description }}</p>
              </div>
              <div class="section-card__badge">
                <span>入口数量</span>
                <strong>{{ section.totalCount }}</strong>
              </div>
            </div>

            <div v-if="section.directApps.length" class="app-grid">
              <article
                v-for="app in section.directApps"
                :key="`app-${app.id}`"
                class="app-card"
                @click="handleCardClick(app)"
              >
                <div class="app-card__icon">
                  <img v-if="isUrl(app.icon)" :src="app.icon" alt="" />
                  <icon
                    v-else
                    :icon="app.icon || 'solar:planet-outline'"
                    :color="app.icon_color || '#3f8f94'"
                    width="24"
                    height="24"
                  />
                </div>
                <div class="app-card__content">
                  <h3>{{ app.title }}</h3>
                  <p>{{ app.description || '暂无补充说明' }}</p>
                </div>
              </article>
            </div>

            <div
              v-for="group in section.children"
              :id="`menu-${group.id}`"
              :key="`group-${group.id}`"
              class="subsection"
            >
              <div class="subsection__header">
                <strong>{{ group.title }}</strong>
                <span>{{ group.apps.length }} 项</span>
              </div>
              <div class="app-grid">
                <article
                  v-for="app in group.apps"
                  :key="`subapp-${app.id}`"
                  class="app-card"
                  @click="handleCardClick(app)"
                >
                  <div class="app-card__icon">
                    <img v-if="isUrl(app.icon)" :src="app.icon" alt="" />
                    <icon
                      v-else
                      :icon="app.icon || 'solar:planet-outline'"
                      :color="app.icon_color || '#3f8f94'"
                      width="24"
                      height="24"
                    />
                  </div>
                  <div class="app-card__content">
                    <h3>{{ app.title }}</h3>
                    <p>{{ app.description || '暂无补充说明' }}</p>
                  </div>
                </article>
              </div>
            </div>
          </section>

          <el-empty v-if="!displaySections.length" description="当前搜索没有匹配结果" />
        </div>
      </section>
    </main>

    <button type="button" class="backtop" @click="scrollToTop">返回顶部</button>

    <footer class="home-footer">
      <p>{{ siteSettings.footer || '个人导航首页' }}</p>
      <span>{{ siteSettings.icp || '未配置备案信息' }}</span>
    </footer>
  </div>
</template>

<script>
import axios from 'axios';
import { mapGetters } from 'vuex';
import { Icon } from '@iconify/vue';

const PAGE_SIZE = 100;

export default {
  components: { Icon },
  data() {
    return {
      menus: [],
      applications: [],
      activeMenu: '',
      searchQuery: '',
      userMenuVisible: false,
      userMenuCloseTimer: null,
    };
  },
  computed: {
    ...mapGetters(['siteSettings', 'username']),
    isLoggedIn() {
      return !!this.$store.state.token;
    },
    displayUsername() {
      return this.username || 'admin';
    },
    activeApplications() {
      return this.applications
        .filter((app) => String(app.status || '').toLowerCase() === 'active')
        .sort((a, b) => this.normalizeOrder(a.order_id) - this.normalizeOrder(b.order_id));
    },
    topLevelMenus() {
      return this.buildTree();
    },
    sectionSummaries() {
      return this.topLevelMenus
        .map((menu) => ({ ...menu, count: this.countSectionApps(menu) }))
        .filter((menu) => menu.count > 0);
    },
    displaySections() {
      const keyword = this.searchQuery.toLowerCase();
      return this.topLevelMenus
        .map((menu) => {
          const directApps = this.filterAppsByMenu(menu.id, keyword);
          const children = (menu.children || [])
            .map((child) => ({ ...child, apps: this.filterAppsByMenu(child.id, keyword) }))
            .filter((child) => child.apps.length > 0 || this.matchesKeyword(child.title, keyword));
          const totalCount = directApps.length + children.reduce((sum, child) => sum + child.apps.length, 0);

          if (!keyword || this.matchesKeyword(menu.title, keyword) || totalCount > 0) {
            return {
              ...menu,
              description: menu.children.length
                ? `当前分组下包含 ${menu.children.length} 个子类，已展示 ${totalCount} 个可用入口。`
                : `聚合这一类常用入口，当前共展示 ${totalCount} 个可用应用。`,
              directApps,
              children,
              totalCount,
            };
          }

          return null;
        })
        .filter((section) => section && section.totalCount > 0);
    },
    visibleApplicationCount() {
      return this.displaySections.reduce((sum, section) => sum + section.totalCount, 0);
    },
  },
  created() {
    this.fetchMenus();
    this.fetchAllApplications();
    this.$store.dispatch('fetchSiteSettings');
    if (this.isLoggedIn) {
      this.$store.dispatch('fetchUserData').catch(() => {});
    }
  },
  methods: {
    openUserMenu() {
      if (this.userMenuCloseTimer) {
        clearTimeout(this.userMenuCloseTimer);
        this.userMenuCloseTimer = null;
      }
      this.userMenuVisible = true;
    },
    closeUserMenu() {
      this.userMenuCloseTimer = window.setTimeout(() => {
        this.userMenuVisible = false;
        this.userMenuCloseTimer = null;
      }, 180);
    },
    toggleUserMenu() {
      if (this.userMenuCloseTimer) {
        clearTimeout(this.userMenuCloseTimer);
        this.userMenuCloseTimer = null;
      }
      this.userMenuVisible = !this.userMenuVisible;
    },
    isUrl(value) {
      return /^(http|https):\/\//.test(value || '');
    },
    normalizeId(value) {
      if (value === undefined || value === null || value === '') {
        return null;
      }

      const normalized = Number(value);
      return Number.isNaN(normalized) ? String(value) : normalized;
    },
    normalizeOrder(value) {
      const normalized = Number(value);
      return Number.isNaN(normalized) ? 0 : normalized;
    },
    normalizeParentId(parentId) {
      const normalized = this.normalizeId(parentId);
      return normalized === 0 ? null : normalized;
    },
    matchesKeyword(value, keyword) {
      if (!keyword) return true;
      return String(value || '').toLowerCase().includes(keyword);
    },
    buildTree(parentId = null) {
      return this.menus
        .filter((menu) => this.normalizeParentId(menu.parent_id) === this.normalizeId(parentId))
        .sort((a, b) => this.normalizeOrder(a.order_id) - this.normalizeOrder(b.order_id))
        .map((menu) => ({
          ...menu,
          children: this.buildTree(menu.id),
        }));
    },
    countSectionApps(menu) {
      const direct = this.activeApplications.filter((app) => this.normalizeId(app.menu_id) === this.normalizeId(menu.id)).length;
      const children = (menu.children || []).reduce((sum, child) => sum + this.countSectionApps(child), 0);
      return direct + children;
    },
    filterAppsByMenu(menuId, keyword = '') {
      return this.activeApplications.filter((app) => {
        if (this.normalizeId(app.menu_id) !== this.normalizeId(menuId)) return false;
        if (!keyword) return true;
        return [app.title, app.description].some((field) => this.matchesKeyword(field, keyword));
      });
    },
    async fetchPaged(endpoint, listKey) {
      let page = 1;
      let total = 0;
      const items = [];

      do {
        const response = await axios.get(`${process.env.VUE_APP_API_URL}${endpoint}`, {
          params: {
            page,
            pageSize: PAGE_SIZE,
          },
        });

        const batch = Array.isArray(response.data[listKey]) ? response.data[listKey] : [];
        total = Number(response.data.total || 0);
        items.push(...batch);
        page += 1;

        if (batch.length === 0) {
          break;
        }
      } while (items.length < total);

      return items;
    },
    async fetchMenus() {
      try {
        const menus = await this.fetchPaged('/menus', 'menus');
        this.menus = menus.filter((menu) => String(menu.status || '').toLowerCase() === 'active');
      } catch (error) {
        console.error('Failed to fetch menus:', error);
      }
    },
    async fetchAllApplications() {
      try {
        this.applications = await this.fetchPaged('/applications', 'applications');
      } catch (error) {
        console.error('Failed to fetch applications:', error);
      }
    },
    handleCardClick(app) {
      if (app.link) {
        window.open(app.link, '_blank', 'noopener');
      }
    },
    scrollToMenu(menuId) {
      this.activeMenu = String(menuId);
      this.$nextTick(() => {
        const target = document.getElementById(`menu-${menuId}`);
        if (target) {
          target.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
      });
    },
    scrollToTop() {
      const target = this.$refs.topAnchor;
      if (target && typeof target.scrollIntoView === 'function') {
        target.scrollIntoView({ behavior: 'smooth', block: 'start' });
        return;
      }

      window.scrollTo({ top: 0, behavior: 'smooth' });
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
    },
    goToLogin() {
      this.$router.push('/login');
    },
    goToDashboard() {
      this.userMenuVisible = false;
      this.$router.push('/dashboard');
    },
    logout() {
      this.userMenuVisible = false;
      this.$store.commit('clearToken');
      this.$router.push('/');
    },
    goToHome() {
      this.scrollToTop();
    },
  },
};
</script>

<style scoped>
.home-shell { min-height: 100vh; padding: 22px; color: var(--admin-text); }
.home-topbar, .toolbar-card, .catalog__inner, .section-card {
  border: 1px solid rgba(255,255,255,.78);
  background: rgba(255,255,255,.84);
  box-shadow: var(--admin-shadow-sm);
  backdrop-filter: blur(16px);
}
.home-topbar {
  position: sticky;
  top: 16px;
  z-index: 20;
  width: min(1480px,100%);
  margin: 0 auto 18px;
  min-height: 78px;
  padding: 14px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  border-radius: 28px;
}
.brand {
  padding: 0;
  display: inline-flex;
  align-items: center;
  gap: 14px;
  border: none;
  background: transparent;
  cursor: pointer;
  color: inherit;
}
.brand__logo {
  width: 52px;
  height: 52px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 18px;
  background: linear-gradient(135deg, rgba(99,179,181,.12), rgba(125,211,252,.1));
  overflow: hidden;
}
.brand__logo img, .app-card__icon img { width: 100%; height: 100%; object-fit: cover; }
.brand__text { display: flex; flex-direction: column; align-items: flex-start; gap: 2px; }
.brand__text strong { font-size: 18px; }
.brand__text small, .metric-card p, .section-card__header p, .app-card__content p, .home-footer, .section-card__badge span {
  color: var(--admin-text-soft);
}
.topbar-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
}
.soft-btn, .ghost-btn, .logout-btn, .backtop, .user-trigger, .user-menu__item {
  min-height: 40px;
  padding: 0 16px;
  border-radius: 999px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--admin-transition);
}
.soft-btn, .backtop {
  border: none;
  color: #fff;
  background: linear-gradient(135deg, #63b4b1, #478f95);
  box-shadow: 0 14px 28px rgba(71,143,149,.16);
}
.ghost-btn, .logout-btn, .user-trigger, .user-menu__item {
  border: 1px solid var(--admin-border);
  background: rgba(255,255,255,.78);
  color: var(--admin-text);
}
.logout-btn {
  color: #9f1239;
  border-color: #fbcfe8;
  background: #fff7fb;
}
.soft-btn:hover, .ghost-btn:hover, .logout-btn:hover, .catalog__item:hover, .app-card:hover, .user-trigger:hover, .user-menu__item:hover {
  transform: translateY(-1px);
}
.user-dropdown {
  position: relative;
  display: flex;
  align-items: center;
  width: fit-content;
}
.user-trigger {
  min-width: 20px;
  padding: 8px 12px 8px 10px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  justify-content: flex-start;
  border-radius: 20px;
  background: linear-gradient(180deg, rgba(255,255,255,.92), rgba(243,250,251,.96));
  box-shadow: inset 0 0 0 1px rgba(220, 235, 237, 0.9);
}
.user-trigger__avatar {
  width: 36px;
  height: 36px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  background: linear-gradient(135deg, #63b4b1, #478f95);
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  flex-shrink: 0;
}
.user-trigger__meta {
  min-width: 0;
  display: flex;
  flex: 1;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
}
.user-trigger__name {
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  line-height: 1.2;
}
.user-trigger__role {
  color: var(--admin-text-faint);
  font-size: 11px;
  line-height: 1.2;
}
.user-trigger__arrow {
  width: 26px;
  height: 26px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: rgba(234, 246, 247, 0.9);
  color: var(--admin-primary);
  flex-shrink: 0;
  transition: transform var(--admin-transition), background var(--admin-transition);
}
.user-trigger__arrow.is-open {
  transform: rotate(180deg);
}
.user-menu {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  width: 100%;
  min-width: 100%;
  box-sizing: border-box;
  padding: 8px;
  display: grid;
  gap: 6px;
  border: 1px solid var(--admin-border);
  border-radius: 20px;
  background: rgba(255,255,255,.97);
  box-shadow: 0 18px 32px rgba(15, 23, 42, 0.08);
  backdrop-filter: blur(18px);
}
.user-menu__item {
  width: 100%;
  min-height: 42px;
  justify-content: center;
  border-radius: 14px;
  background: rgba(255,255,255,.85);
}
.home-main { width: min(1480px,100%); margin: 0 auto; display: flex; flex-direction: column; gap: 18px; }
.toolbar-card {
  padding: 20px;
  border-radius: 30px;
}
.search-box {
  min-height: 56px;
  padding: 0 16px;
  display: grid;
  grid-template-columns: auto minmax(0,1fr) auto;
  align-items: center;
  gap: 12px;
  border-radius: 20px;
  border: 1px solid var(--admin-border);
  background: rgba(255,255,255,.92);
}
.search-box input, .search-box button {
  border: none;
  background: transparent;
  outline: none;
  font-size: 14px;
}
.search-box button { cursor: pointer; color: var(--admin-text-soft); }
.metric-grid { margin-top: 18px; display: grid; grid-template-columns: repeat(3, minmax(0,1fr)); gap: 12px; }
.metric-card, .catalog__item, .app-card, .section-card__badge {
  border: 1px solid var(--admin-border);
  background: rgba(255,255,255,.82);
}
.metric-card {
  padding: 18px;
  border-radius: 22px;
}
.metric-card span, .subsection__header span {
  font-size: 12px;
  color: var(--admin-text-faint);
}
.metric-card strong, .section-card__badge strong {
  display: block;
  margin-top: 10px;
  font-size: 32px;
  line-height: 1;
}
.content-layout { display: grid; grid-template-columns: 260px minmax(0,1fr); gap: 18px; align-items: start; }
.catalog { position: sticky; top: 24px; }
.catalog__inner {
  padding: 20px;
  border-radius: 30px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.catalog__label {
  margin: 0 0 10px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: .16em;
  text-transform: uppercase;
  color: var(--admin-primary);
}
.catalog__item {
  width: 100%;
  min-height: 46px;
  padding: 10px 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border-radius: 16px;
  cursor: pointer;
  transition: all var(--admin-transition);
  text-align: left;
}
.catalog__item.is-active {
  background: var(--admin-primary-soft);
  color: var(--admin-primary);
  border-color: rgba(63,143,148,.16);
}
.sections { display: grid; gap: 18px; }
.section-card {
  padding: 20px;
  border-radius: 30px;
}
.section-card__header, .subsection__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}
.section-card__header h2, .app-card__content h3, .subsection__header strong { margin: 0; }
.section-card__header h2 {
  font-size: 22px;
  line-height: 1.15;
}
.section-card__header p {
  margin: 6px 0 0;
  font-size: 13px;
  line-height: 1.8;
}
.section-card__badge {
  min-width: 110px;
  padding: 14px 16px;
  border-radius: 18px;
  text-align: right;
}
.app-grid {
  margin-top: 18px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px,1fr));
  gap: 14px;
}
.app-card {
  padding: 16px;
  display: grid;
  grid-template-columns: auto minmax(0,1fr);
  align-items: center;
  gap: 14px;
  border-radius: 22px;
  cursor: pointer;
  transition: all var(--admin-transition);
}
.app-card__icon {
  width: 48px;
  height: 48px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: #eef8f8;
  overflow: hidden;
}
.app-card__content p {
  margin-top: 8px;
  font-size: 12px;
  line-height: 1.7;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.app-card__content h3 {
  font-size: 14px;
  line-height: 1.3;
}
.subsection__header strong {
  font-size: 15px;
}
.subsection {
  margin-top: 18px;
  padding-top: 18px;
  border-top: 1px solid #eef4f7;
}
.backtop {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 15;
}
.home-footer {
  width: min(1480px,100%);
  margin: 18px auto 0;
  padding: 18px 4px 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-size: 13px;
}
.home-footer p { margin: 0; }
@media (max-width: 1240px) {
  .content-layout { grid-template-columns: 1fr; }
  .catalog { position: static; }
}
@media (max-width: 900px) {
  .home-shell { padding: 14px; }
  .home-topbar { flex-direction: column; align-items: stretch; }
  .topbar-actions { justify-content: flex-start; }
  .user-trigger,
  .user-menu { min-width: 100%; }
  .metric-grid, .app-grid { grid-template-columns: 1fr; }
  .section-card__header, .subsection__header, .home-footer { flex-direction: column; align-items: flex-start; }
  .section-card__badge { width: 100%; text-align: left; }
}
@media (max-width: 640px) {
  .toolbar-card, .catalog__inner, .section-card { padding: 18px; border-radius: 24px; }
}
</style>
