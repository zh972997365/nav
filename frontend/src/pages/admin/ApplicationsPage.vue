<template>
  <AdminPageShell
    eyebrow="Applications"
    title="应用管理"
    description="统一维护站点应用的标题、图标、菜单归属与状态。"
  >
    <template #meta>
      <div class="admin-pill">
        <span>应用总数</span>
        <strong>{{ totalApplications }}</strong>
      </div>
    </template>

    <section class="admin-workspace apps-page">
      <div class="admin-workspace__top">
        <div class="admin-workspace__summary">
          <article class="admin-simple-metric">
            <p class="admin-simple-metric__eyebrow">全部应用</p>
            <p class="admin-simple-metric__value">{{ totalApplications }}</p>
            <p class="admin-simple-metric__hint">当前已收录到后台的全部应用数量</p>
          </article>
          <article class="admin-simple-metric">
            <p class="admin-simple-metric__eyebrow">启用中</p>
            <p class="admin-simple-metric__value">{{ activeCount }}</p>
            <p class="admin-simple-metric__hint">正在展示或可用的应用</p>
          </article>
          <article class="admin-simple-metric">
            <p class="admin-simple-metric__eyebrow">已选择</p>
            <p class="admin-simple-metric__value">{{ multipleSelection.length }}</p>
            <p class="admin-simple-metric__hint">当前可批量操作的选中项</p>
          </article>
        </div>

        <div class="admin-workspace__filters">
          <div class="admin-workspace__filters-main">
            <div class="admin-workspace__filters-copy">
              <h3 class="admin-workspace__filters-title">筛选应用</h3>
              <p class="admin-workspace__filters-description">
                先按状态快速缩小范围，再结合标题和菜单定位目标应用。
              </p>
            </div>

            <div class="admin-filter-tabs">
              <button type="button" class="admin-filter-tab" :class="{ 'is-active': !filters.status }" @click="filters.status = ''">
                <span>全部</span>
              </button>
              <button type="button" class="admin-filter-tab" :class="{ 'is-active': filters.status === 'active' }" @click="filters.status = 'active'">
                <span>启用</span>
              </button>
              <button type="button" class="admin-filter-tab" :class="{ 'is-active': filters.status === 'inactive' }" @click="filters.status = 'inactive'">
                <span>停用</span>
              </button>
            </div>

            <div class="admin-filter-grid">
              <div class="admin-filter-field admin-filter-field--wide">
                <label for="app-title">应用标题</label>
                <el-input id="app-title" v-model="filters.title" placeholder="搜索应用标题" @keydown.enter="handleSearch" />
              </div>

              <div class="admin-filter-field">
                <label for="app-menu">所属菜单</label>
                <el-cascader
                  id="app-menu"
                  v-model="filters.menu_id"
                  :options="menuOptions"
                  placeholder="选择菜单"
                  clearable
                  :props="{ expandTrigger: 'hover', emitPath: false, checkStrictly: true }"
                >
                  <template #default="{ data }">
                    <span class="apps-page__menu-option">
                      <img v-if="isUrl(data.icon)" :src="data.icon" width="18" height="18" alt="icon" />
                      <icon v-else :icon="data.icon" width="16" height="16" :style="{ color: data.icon_color || '#3f8f94' }" />
                      <span>{{ data.label }}</span>
                    </span>
                  </template>
                </el-cascader>
              </div>
            </div>

            <div class="admin-filter-actions">
              <el-button type="primary" class="admin-soft-button" @click="handleSearch">查询</el-button>
              <el-button class="admin-neutral-button" @click="handleReset">重置筛选</el-button>
            </div>
          </div>

          <aside class="admin-workspace__filters-side apps-filters-side">
            <div>
              <p class="admin-workspace__filters-side-label">Active Filters</p>
              <p class="admin-workspace__filters-side-value">{{ activeFilterCount }}</p>
            </div>
            <p class="admin-workspace__filters-side-hint">
              当前已启用 {{ activeFilterCount }} 个筛选条件，列表会随之联动更新。
            </p>
          </aside>
        </div>
      </div>

      <div class="admin-workspace__actions">
        <div class="admin-workspace__actions-copy">
          应用列表保留批量操作，筛选与操作分区展示，页面会更清爽也更容易定位。
        </div>
        <div class="admin-workspace__actions-buttons">
          <el-button type="primary" class="admin-soft-button" @click="openForm('add')">新增应用</el-button>
          <el-button type="danger" class="admin-danger-button" :disabled="multipleSelection.length === 0" @click="handleBatchDelete">
            批量删除
          </el-button>
        </div>
      </div>

      <div class="admin-workspace__table">
        <el-table
          v-if="applications.length > 0"
          :data="applications"
          border
          stripe
          :max-height="560"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="应用" min-width="300">
            <template #default="{ row }">
              <div class="apps-row">
                <div class="apps-row__icon">
                  <img v-if="isUrl(row.icon)" :src="row.icon" width="24" height="24" alt="icon" />
                  <icon v-else :icon="row.icon" :color="row.icon_color || '#3f8f94'" width="18" height="18" />
                </div>
                <div class="apps-row__content">
                  <strong>{{ row.title }}</strong>
                  <el-tooltip :content="row.description || '暂无描述'" placement="top-start" :show-after="250">
                    <p class="apps-row__description">{{ row.description || '暂无描述' }}</p>
                  </el-tooltip>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="link" label="链接" min-width="220" show-overflow-tooltip />
          <el-table-column label="菜单" min-width="180">
            <template #default="{ row }">
              <div class="apps-page__menu-option">
                <img v-if="isUrl(getMenuIcon(row.menu_id))" :src="getMenuIcon(row.menu_id)" width="18" height="18" alt="menu" />
                <icon
                  v-else
                  :icon="getMenuIcon(row.menu_id) || 'solar:sidebar-code-outline'"
                  :color="getMenuIconColor(row.menu_id) || '#94a3b8'"
                  width="16"
                  height="16"
                />
                <span>{{ getMenuTitle(row.menu_id) }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="120" align="center">
            <template #default="{ row }">
              <span class="apps-status" :class="row.status === 'active' ? 'is-active' : 'is-inactive'">
                {{ row.status === 'active' ? '启用中' : '停用中' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="order_id" label="排序" width="90" align="center" />
          <el-table-column label="创建时间" width="190">
            <template #default="{ row }">
              {{ new Date(row.created_at).toLocaleString() }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" align="center" fixed="right">
            <template #default="{ row }">
              <div class="apps-row__actions">
                <el-button link class="admin-link-button admin-link-button--edit" @click="openForm('edit', row)">
                  <icon icon="solar:pen-new-square-outline" width="16" height="16" />
                </el-button>
                <el-button link class="admin-link-button admin-link-button--delete" @click="deleteApp(row.id)">
                  <icon icon="solar:trash-bin-trash-outline" width="16" height="16" />
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <el-empty v-else description="暂无应用数据" />

        <div class="admin-table-footer">
          <el-pagination
            background
            layout="total, sizes, prev, pager, next, jumper"
            :total="totalApplications"
            :page-size="pageSize"
            :current-page="currentPage"
            :page-sizes="pageSizeOptions"
            @current-change="handlePageChange"
            @size-change="handlePageSizeChange"
          />
        </div>
      </div>
    </section>

    <el-dialog
      :model-value="showForm"
      :title="isEditing ? '编辑应用' : '新增应用'"
      class="admin-dialog"
      width="760px"
      @close="cancelForm"
    >
      <el-form :model="app" label-width="90px">
        <div class="admin-form-grid">
          <el-form-item label="链接" class="full-span">
            <el-input v-model="app.link" placeholder="请输入应用链接" />
          </el-form-item>
          <el-form-item label="采集" class="full-span">
            <el-button type="primary" @click="fetchWebsiteData" :loading="isFetching">采集站点信息</el-button>
          </el-form-item>
          <el-form-item label="标题" class="full-span">
            <el-input v-model="app.title" placeholder="请输入应用标题" />
          </el-form-item>
          <el-form-item label="图标" class="full-span">
            <div class="icon-editor">
              <el-input v-model="app.icon" placeholder="请输入 Iconify 图标或 URL" />
              <el-button text class="icon-editor__action" @click="syncIconToImageHost" :disabled="!isUrl(app.icon) || isSyncingIcon">
                <icon :icon="isSyncingIcon ? 'eos-icons:loading' : 'solar:restart-outline'" :spin="isSyncingIcon" width="18" height="18" />
              </el-button>
              <div class="icon-editor__preview">
                <img v-if="isUrl(app.icon)" :src="app.icon" width="26" height="26" alt="icon" />
                <icon v-else :icon="app.icon || 'solar:widget-5-outline'" :color="app.icon_color || '#3f8f94'" width="20" height="20" />
              </div>
              <el-color-picker v-model="app.icon_color" />
            </div>
            <a href="https://icon-sets.iconify.design/?category=General" target="_blank" class="admin-icon-link" rel="noreferrer">
              前往 Iconify 选择图标
            </a>
          </el-form-item>
          <el-form-item label="描述" class="full-span">
            <el-input v-model="app.description" type="textarea" placeholder="请输入应用描述" />
          </el-form-item>
          <el-form-item label="状态">
            <el-switch v-model="app.status" active-value="active" inactive-value="inactive" active-text="启用" inactive-text="停用" />
          </el-form-item>
          <el-form-item label="菜单">
            <el-cascader
              v-model="app.menu_id"
              :options="menuOptions"
              placeholder="请选择菜单"
              clearable
              :props="{ expandTrigger: 'hover', emitPath: false, checkStrictly: true }"
            >
              <template #default="{ data }">
                <span class="apps-page__menu-option">
                  <img v-if="isUrl(data.icon)" :src="data.icon" width="18" height="18" alt="icon" />
                  <icon v-else :icon="data.icon" width="16" height="16" :style="{ color: data.icon_color || '#3f8f94' }" />
                  <span>{{ data.label }}</span>
                </span>
              </template>
            </el-cascader>
          </el-form-item>
          <el-form-item label="排序">
            <el-input-number v-model="app.order_id" :min="-100" :max="100" />
          </el-form-item>
        </div>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button class="admin-neutral-button" @click="resetForm">重置</el-button>
          <el-button class="admin-neutral-button" @click="cancelForm">取消</el-button>
          <el-button type="primary" @click="submitApp">保存</el-button>
        </div>
      </template>
    </el-dialog>
  </AdminPageShell>
</template>

<script>
import { Icon } from '@iconify/vue';
import axios from 'axios';
import AdminPageShell from '../../components/admin/AdminPageShell.vue';

export default {
  components: {
    Icon,
    AdminPageShell,
  },
  data() {
    return {
      showForm: false,
      isEditing: false,
      originalApp: {},
      isFetching: false,
      isSyncingIcon: false,
      app: {
        title: '',
        icon: '',
        icon_color: '#3f8f94',
        link: '',
        description: '',
        status: 'active',
        menu_id: '',
        order_id: 0,
      },
      applications: [],
      menus: [],
      menuOptions: [],
      currentPage: 1,
      pageSize: 10,
      totalApplications: 0,
      pageSizeOptions: [10, 20, 30, 40],
      multipleSelection: [],
      filters: {
        title: '',
        status: '',
        menu_id: '',
      },
    };
  },
  computed: {
    settings: {
      get() {
        return this.$store.getters.siteSettings;
      },
      set(value) {
        this.$store.commit('setSiteSettings', value);
      },
    },
    activeCount() {
      return this.applications.filter((item) => item.status === 'active').length;
    },
    activeFilterCount() {
      return Object.values(this.filters).filter(Boolean).length;
    },
  },
  created() {
    this.fetchMenus();
  },
  methods: {
    isUrl(string) {
      return !!string && (string.startsWith('http://') || string.startsWith('https://'));
    },
    getMenuIcon(menuId) {
      const menu = this.menus.find((item) => item.id === menuId);
      return menu ? menu.icon : '';
    },
    getMenuIconColor(menuId) {
      const menu = this.menus.find((item) => item.id === menuId);
      return menu ? menu.icon_color : '';
    },
    getMenuTitle(menuId) {
      const menu = this.menus.find((item) => item.id === menuId);
      return menu ? menu.title : '未指定';
    },
    fetchApplications() {
      axios.get(`${process.env.VUE_APP_API_URL}/applications`, {
        params: {
          page: this.currentPage,
          pageSize: this.pageSize,
          title: this.filters.title,
          status: this.filters.status,
          menuId: this.filters.menu_id,
        },
        headers: { Authorization: `Bearer ${this.$store.state.token}` },
      }).then((response) => {
        const uniqueApplications = Array.from(new Set(response.data.applications.map((item) => item.id)))
          .map((id) => response.data.applications.find((item) => item.id === id));
        this.applications = uniqueApplications.sort((a, b) => a.menu_id - b.menu_id);
        this.totalApplications = response.data.total;
      }).catch((error) => {
        console.error('Error fetching applications:', error);
      });
    },
    fetchMenus() {
      axios.get(`${process.env.VUE_APP_API_URL}/menus`, {
        params: { pageSize: 1000 },
        headers: { Authorization: `Bearer ${this.$store.state.token}` },
      }).then((response) => {
        this.menus = response.data.menus;
        this.menuOptions = this.menus
          .filter((menu) => menu.parent_id === null)
          .map((menu) => ({
            value: menu.id,
            label: menu.title,
            icon: menu.icon,
            icon_color: menu.icon_color,
            children: this.menus
              .filter((submenu) => submenu.parent_id === menu.id)
              .map((submenu) => ({
                value: submenu.id,
                label: submenu.title,
                icon: submenu.icon,
                icon_color: submenu.icon_color,
              })),
          }));
        this.fetchApplications();
      }).catch((error) => {
        console.error('Error fetching menus:', error);
      });
    },
    handleSearch() {
      this.currentPage = 1;
      this.fetchApplications();
    },
    handleReset() {
      this.filters = { title: '', status: '', menu_id: '' };
      this.currentPage = 1;
      this.fetchApplications();
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
    handleBatchDelete() {
      if (!this.multipleSelection.length) return;
      this.$confirm('此操作将永久删除所选应用，是否继续？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(() => {
        const ids = this.multipleSelection.map((item) => item.id);
        axios.post(`${process.env.VUE_APP_API_URL}/applications/batch_delete`, { ids }, {
          headers: { Authorization: `Bearer ${this.$store.state.token}` },
        }).then(() => {
          this.fetchApplications();
          this.multipleSelection = [];
          this.$message.success('应用已删除');
        });
      }).catch(() => {
        this.$message.info('已取消删除');
      });
    },
    handlePageChange(page) {
      this.currentPage = page;
      this.fetchApplications();
    },
    handlePageSizeChange(size) {
      this.pageSize = size;
      this.currentPage = 1;
      this.fetchApplications();
    },
    hasChanges() {
      return Object.keys(this.app).some((key) => this.app[key] !== this.originalApp[key]);
    },
    submitApp() {
      if (!this.app.menu_id) {
        this.$message.error('请选择一个菜单');
        return;
      }

      if (this.isEditing && this.hasChanges()) {
        axios.put(`${process.env.VUE_APP_API_URL}/applications/${this.app.id}`, this.app, {
          headers: { Authorization: `Bearer ${this.$store.state.token}` },
        }).then(() => {
          this.fetchApplications();
          this.$message.success('应用已更新');
          this.resetForm();
          this.showForm = false;
        }).catch((error) => {
          console.error('Error submitting application:', error);
        });
      } else if (!this.isEditing) {
        axios.post(`${process.env.VUE_APP_API_URL}/applications`, this.app, {
          headers: { Authorization: `Bearer ${this.$store.state.token}` },
        }).then(() => {
          this.fetchApplications();
          this.$message.success('应用已添加');
          this.resetForm();
          this.showForm = false;
        }).catch((error) => {
          console.error('Error submitting application:', error);
        });
      } else {
        this.$message.info('未检测到任何修改');
      }
    },
    editApp(application) {
      this.app = { ...application };
      this.originalApp = { ...application };
      this.isEditing = true;
      this.showForm = true;
    },
    deleteApp(id) {
      this.$confirm('此操作将永久删除该应用，是否继续？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(() => {
        axios.delete(`${process.env.VUE_APP_API_URL}/applications/${id}`, {
          headers: { Authorization: `Bearer ${this.$store.state.token}` },
        }).then(() => {
          this.fetchApplications();
          this.$message.success('应用已删除');
        });
      }).catch(() => {
        this.$message.info('已取消删除');
      });
    },
    cancelForm() {
      this.resetForm();
      this.showForm = false;
    },
    resetForm() {
      this.app = {
        title: '',
        icon: '',
        icon_color: '#3f8f94',
        link: '',
        description: '',
        status: 'active',
        menu_id: '',
        order_id: 0,
      };
      this.originalApp = {};
      this.isEditing = false;
    },
    openForm(action, app = null) {
      if (action === 'edit' && app) {
        this.editApp(app);
      } else {
        this.resetForm();
        this.showForm = true;
      }
    },
    async fetchWebsiteData() {
      if (!this.app.link) {
        this.$message.error('请先输入应用链接');
        return;
      }
      if (!this.isUrl(this.app.link)) {
        this.$message.error('请输入有效的 URL，例如 https://example.com');
        return;
      }

      this.isFetching = true;
      try {
        const response = await axios.get(`${process.env.VUE_APP_API_URL}/scrape-website`, {
          params: { url: this.app.link },
          headers: { Authorization: `Bearer ${this.$store.state.token}` },
        });
        const data = response.data;
        this.app.title = data.title || '';
        this.app.icon = data.logo_link || '';
        this.app.description = data.description || '';
        this.$message.success('站点信息已填入');
      } catch (error) {
        console.error('Error fetching website data:', error);
        this.$message.error('采集失败，请稍后再试');
      } finally {
        this.isFetching = false;
      }
    },
    async syncIconToImageHost() {
      if (!this.isUrl(this.app.icon)) {
        this.$message.error('当前图标不是有效 URL');
        return;
      }
      if (!this.settings || !this.settings.image_host_token) {
        this.$message.error('请先在站点设置中保存图床 Token');
        return;
      }

      this.isSyncingIcon = true;
      try {
        const response = await axios.post(`${process.env.VUE_APP_API_URL}/sync-icon`, {
          icon_url: this.app.icon,
          token: this.settings.image_host_token,
        });
        if (response.data.icon_url) {
          this.app.icon = response.data.icon_url;
          this.$message.success('图标已同步到图床');
        } else {
          this.$message.error('同步图标失败');
        }
      } catch (error) {
        console.error('Error syncing icon to image host:', error);
        this.$message.error('同步图标失败');
      } finally {
        this.isSyncingIcon = false;
      }
    },
  },
};
</script>

<style scoped>
.apps-page {
  border: 1px solid #d9e9eb;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(246, 251, 251, 0.98));
}

.apps-page :deep(.admin-simple-metric) {
  border-color: #dcebed;
  background: #ffffff;
}

.apps-page :deep(.admin-simple-metric__eyebrow) {
  color: #7f98a3;
}

.apps-page :deep(.admin-simple-metric__hint),
.apps-page :deep(.admin-workspace__actions-copy) {
  color: #6e8791;
}

.apps-page :deep(.admin-workspace__filters) {
  border-color: #dcebed;
  background: linear-gradient(180deg, #ffffff, #f4fbfb);
}

.apps-page :deep(.admin-filter-tabs) {
  background: #eef8f8;
  border-color: #dcebed;
}

.apps-page :deep(.admin-workspace__actions) {
  border-top-color: #e4eff1;
  border-bottom-color: #e4eff1;
  background: rgba(247, 252, 252, 0.86);
}

.apps-page :deep(.admin-soft-button.el-button),
.apps-page :deep(.el-button--primary:not(.is-text):not(.is-link)) {
  background: linear-gradient(135deg, #63b4b1, #478f95) !important;
  box-shadow: 0 14px 28px rgba(71, 143, 149, 0.16) !important;
}

.apps-page :deep(.admin-soft-button.el-button:hover),
.apps-page :deep(.admin-soft-button.el-button:focus),
.apps-page :deep(.el-button--primary:not(.is-text):not(.is-link):hover),
.apps-page :deep(.el-button--primary:not(.is-text):not(.is-link):focus) {
  background: linear-gradient(135deg, #55a4a1, #3f8388) !important;
}

.apps-page :deep(.admin-neutral-button.el-button) {
  background: #ffffff !important;
  border-color: #dcebed !important;
}

.apps-page :deep(.el-input__wrapper),
.apps-page :deep(.el-cascader .el-input__wrapper) {
  box-shadow: 0 0 0 1px #dcebed inset !important;
  background: #ffffff !important;
}

.apps-page :deep(.el-input__wrapper.is-focus),
.apps-page :deep(.el-cascader .el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px rgba(71, 143, 149, 0.3), 0 0 0 4px rgba(71, 143, 149, 0.08) !important;
}

.apps-filters-side {
  background: linear-gradient(180deg, #6cbab6, #4e9196) !important;
  color: #f4feff !important;
}

.apps-page__menu-option {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.apps-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.apps-row__icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: #eef8f8;
  overflow: hidden;
}

.apps-row__content {
  min-width: 0;
}

.apps-row__content strong {
  display: block;
  font-size: 14px;
  color: #16313a;
}

.apps-row__description {
  display: -webkit-box;
  margin: 4px 0 0;
  overflow: hidden;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
  text-overflow: ellipsis;
  max-width: 260px;
  font-size: 12px;
  color: #84a0ad;
  cursor: help;
}

.apps-status {
  min-height: 28px;
  padding: 0 10px;
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
}

.apps-status.is-active {
  color: #1f7a67;
  background: #eaf8f2;
}

.apps-status.is-inactive {
  color: #8f6a44;
  background: #fff5e8;
}

.apps-row__actions {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.full-span {
  grid-column: 1 / -1;
}

.icon-editor {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto auto auto;
  gap: 10px;
  align-items: center;
  width: 100%;
}

.icon-editor__action {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  border: 1px solid #dcebed;
  background: #fff;
}

.icon-editor__preview {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: #eef8f8;
}

@media (max-width: 768px) {
  .icon-editor {
    grid-template-columns: 1fr;
  }
}
</style>
