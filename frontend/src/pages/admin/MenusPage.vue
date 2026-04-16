<template>
  <AdminPageShell
    eyebrow="Navigation"
    title="菜单管理"
    description="维护后台导航结构、层级关系与启停状态。筛选区与应用管理同步为同一套纵向工作台风格。"
  >
    <template #meta>
      <div class="admin-pill">
        <span>节点总数</span>
        <strong>{{ totalMenus }}</strong>
      </div>
    </template>

    <section class="admin-workspace">
      <div class="admin-workspace__top">
        <div class="admin-workspace__summary">
          <article class="admin-simple-metric">
            <p class="admin-simple-metric__eyebrow">全部菜单</p>
            <p class="admin-simple-metric__value">{{ totalMenus }}</p>
            <p class="admin-simple-metric__hint">当前所有导航节点数量</p>
          </article>
          <article class="admin-simple-metric">
            <p class="admin-simple-metric__eyebrow">顶级节点</p>
            <p class="admin-simple-metric__value">{{ rootCount }}</p>
            <p class="admin-simple-metric__hint">不包含父级的一级菜单</p>
          </article>
          <article class="admin-simple-metric">
            <p class="admin-simple-metric__eyebrow">启用中</p>
            <p class="admin-simple-metric__value">{{ activeCount }}</p>
            <p class="admin-simple-metric__hint">当前仍在显示的导航节点</p>
          </article>
        </div>

        <div class="admin-workspace__filters">
          <div class="admin-workspace__filters-main">
            <div class="admin-workspace__filters-copy">
              <h3 class="admin-workspace__filters-title">筛选菜单</h3>
              <p class="admin-workspace__filters-description">
                状态与标题搜索拆成清晰两层，既能快速缩小范围，也不会挤压后续操作区域。
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
                <label for="menu-title">菜单标题</label>
                <el-input id="menu-title" v-model="filters.title" placeholder="搜索菜单标题" @keydown.enter="handleSearch" />
              </div>
            </div>

            <div class="admin-filter-actions">
              <el-button class="admin-neutral-button" @click="handleReset">重置筛选</el-button>
              <el-button type="primary" class="admin-soft-button" @click="handleSearch">查询</el-button>
            </div>
          </div>

          <aside class="admin-workspace__filters-side">
            <div>
              <p class="admin-workspace__filters-side-label">Visible Rows</p>
              <p class="admin-workspace__filters-side-value">{{ formattedMenus.length }}</p>
            </div>
            <p class="admin-workspace__filters-side-hint">
              当前树形列表展示 {{ formattedMenus.length }} 个节点，可直接在表格中展开或收起层级。
            </p>
          </aside>
        </div>
      </div>

      <div class="admin-workspace__actions">
        <div class="admin-workspace__actions-copy">
          菜单列表保留层级展开逻辑，但把主要操作集中到右侧按钮组，阅读路径会更顺畅。
        </div>
        <div class="admin-workspace__actions-buttons">
          <el-button type="primary" class="admin-soft-button" @click="openForm('add')">新增菜单</el-button>
          <el-button type="danger" class="admin-danger-button" :disabled="multipleSelection.length === 0" @click="handleBatchDelete">
            批量删除
          </el-button>
        </div>
      </div>

      <div class="admin-workspace__table">
        <el-table
          v-if="formattedMenus.length > 0"
          :data="formattedMenus"
          border
          stripe
          :max-height="520"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="菜单" min-width="280">
            <template #default="{ row }">
              <div class="menu-item" :style="{ paddingLeft: `${row.level * 18}px` }">
                <button v-if="row.hasChildren" type="button" class="menu-item__toggle" @click="toggleMenu(row.id)">
                  <icon :icon="row.expanded ? 'solar:alt-arrow-down-outline' : 'solar:alt-arrow-right-outline'" width="16" height="16" />
                </button>
                <span v-else class="menu-item__toggle menu-item__toggle--empty"></span>
                <div class="menu-item__icon">
                  <img v-if="isUrl(row.icon)" :src="row.icon" alt="icon" width="22" height="22" />
                  <icon v-else :icon="row.icon || 'solar:sidebar-code-outline'" :color="row.icon_color || '#3f8f94'" width="18" height="18" />
                </div>
                <div>
                  <strong>{{ row.title }}</strong>
                  <p>ID #{{ row.id }} · {{ row.parent_id ? '子级菜单' : '顶级菜单' }}</p>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="order_id" label="排序" width="90" align="center" />
          <el-table-column label="状态" width="120" align="center">
            <template #default="{ row }">
              <span class="status-pill" :class="row.status === 'active' ? 'is-success' : 'is-warning'">
                {{ row.status === 'active' ? '启用中' : '停用中' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="190">
            <template #default="{ row }">
              {{ new Date(row.created_at).toLocaleString() }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" align="center" fixed="right">
            <template #default="{ row }">
              <div class="row-actions">
                <el-button link class="admin-link-button admin-link-button--edit" @click="openForm('edit', row)">
                  <icon icon="solar:pen-new-square-outline" width="16" height="16" />
                </el-button>
                <el-button link class="admin-link-button admin-link-button--delete" @click="deleteMenu(row.id)">
                  <icon icon="solar:trash-bin-trash-outline" width="16" height="16" />
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <el-empty v-else description="暂无菜单数据" />

        <div class="admin-table-footer">
          <el-pagination
            background
            layout="total, sizes, prev, pager, next, jumper"
            :total="totalMenus"
            :page-size="pageSize"
            :current-page="currentPage"
            :page-sizes="pageSizeOptions"
            @current-change="handlePageChange"
            @size-change="handlePageSizeChange"
          />
        </div>
      </div>
    </section>

    <el-dialog :model-value="showForm" :title="isEditing ? '编辑菜单' : '新增菜单'" class="admin-dialog" width="720px" @close="cancelForm">
      <el-form :model="menu" label-width="96px">
        <div class="admin-form-grid">
          <el-form-item label="标题" class="full-span">
            <el-input v-model="menu.title" placeholder="请输入菜单标题" />
          </el-form-item>
          <el-form-item label="图标" class="full-span">
            <div class="icon-editor">
              <el-input v-model="menu.icon" placeholder="请输入 Iconify 图标或 URL" />
              <div class="icon-editor__preview">
                <img v-if="isUrl(menu.icon)" :src="menu.icon" width="26" height="26" alt="icon" />
                <icon v-else :icon="menu.icon || 'solar:sidebar-code-outline'" :color="menu.icon_color || '#3f8f94'" width="20" height="20" />
              </div>
              <el-color-picker v-model="menu.icon_color" />
            </div>
            <a href="https://icon-sets.iconify.design/?category=General" target="_blank" class="admin-icon-link" rel="noreferrer">
              前往 Iconify 选择图标
            </a>
          </el-form-item>
          <el-form-item label="父级菜单">
            <el-select v-model="menu.parent_id" placeholder="请选择父级菜单" clearable>
              <el-option v-for="parent in parentMenus" :key="parent.id" :label="parent.title" :value="parent.id" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-switch v-model="menu.status" active-value="active" inactive-value="inactive" active-text="启用" inactive-text="停用" />
          </el-form-item>
          <el-form-item label="排序">
            <el-input-number v-model="menu.order_id" :min="-100" :max="100" />
          </el-form-item>
        </div>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button class="admin-neutral-button" @click="resetForm">重置</el-button>
          <el-button class="admin-neutral-button" @click="cancelForm">取消</el-button>
          <el-button type="primary" @click="submitMenu">保存</el-button>
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
      originalMenu: {},
      menu: {
        title: '',
        icon: '',
        icon_color: '#3f8f94',
        status: 'active',
        order_id: 0,
        parent_id: null,
      },
      menus: [],
      formattedMenus: [],
      parentMenus: [],
      expandedIds: new Set(),
      currentPage: 1,
      pageSize: 10,
      totalMenus: 0,
      pageSizeOptions: [10, 20, 30, 40],
      multipleSelection: [],
      filters: {
        title: '',
        status: '',
      },
    };
  },
  computed: {
    rootCount() {
      return this.menus.filter((item) => item.parent_id === null).length;
    },
    activeCount() {
      return this.menus.filter((item) => item.status === 'active').length;
    },
  },
  created() {
    this.fetchMenus();
  },
  methods: {
    isUrl(string) {
      return !!string && (string.startsWith('http://') || string.startsWith('https://'));
    },
    fetchMenus() {
      axios.get(`${process.env.VUE_APP_API_URL}/menus`, {
        params: {
          title: this.filters.title,
          status: this.filters.status,
          pageSize: 1000,
        },
        headers: { Authorization: `Bearer ${this.$store.state.token}` },
      }).then((response) => {
        this.menus = response.data.menus.sort((a, b) => a.order_id - b.order_id);
        this.totalMenus = this.menus.length;
        this.formatMenus();
        this.populateParentMenus();
      }).catch((error) => {
        console.error('Error fetching menus:', error);
      });
    },
    formatMenus() {
      const buildTree = (menus, parentId, level) => menus
        .filter((item) => item.parent_id === parentId)
        .map((item) => {
          const children = buildTree(menus, item.id, level + 1);
          return {
            ...item,
            level,
            children,
            hasChildren: children.length > 0,
            expanded: this.expandedIds.has(item.id),
          };
        });

      const flattenTree = (items) => items.reduce((acc, item) => {
        acc.push(item);
        if (item.expanded && item.children.length > 0) {
          acc.push(...flattenTree(item.children));
        }
        return acc;
      }, []);

      this.formattedMenus = flattenTree(buildTree(this.menus, null, 0));
    },
    populateParentMenus() {
      this.parentMenus = this.menus.filter((item) => item.parent_id === null);
    },
    toggleMenu(id) {
      if (this.expandedIds.has(id)) {
        this.expandedIds.delete(id);
      } else {
        this.expandedIds.add(id);
      }
      this.formatMenus();
    },
    handleSearch() {
      this.currentPage = 1;
      this.fetchMenus();
    },
    handleReset() {
      this.filters = { title: '', status: '' };
      this.currentPage = 1;
      this.fetchMenus();
    },
    handleSelectionChange(val) {
      this.multipleSelection = val;
    },
    handleBatchDelete() {
      if (!this.multipleSelection.length) return;
      this.$confirm('此操作将永久删除所选菜单，是否继续？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(() => {
        const ids = this.multipleSelection.map((item) => item.id);
        axios.post(`${process.env.VUE_APP_API_URL}/menus/batch_delete`, { ids }, {
          headers: { Authorization: `Bearer ${this.$store.state.token}` },
        }).then(() => {
          this.fetchMenus();
          this.multipleSelection = [];
          this.$message.success('菜单已删除');
        });
      }).catch(() => {
        this.$message.info('已取消删除');
      });
    },
    handlePageChange(page) {
      this.currentPage = page;
      this.fetchMenus();
    },
    handlePageSizeChange(size) {
      this.pageSize = size;
      this.currentPage = 1;
      this.fetchMenus();
    },
    openForm(action, menu = null) {
      this.isEditing = action === 'edit';
      if (this.isEditing && menu) {
        this.menu = { ...menu };
        this.originalMenu = { ...menu };
      } else {
        this.resetForm();
      }
      this.showForm = true;
    },
    hasChanges() {
      return Object.keys(this.menu).some((key) => this.menu[key] !== this.originalMenu[key]);
    },
    submitMenu() {
      if (this.isEditing && this.hasChanges()) {
        axios.put(`${process.env.VUE_APP_API_URL}/menus/${this.menu.id}`, this.menu, {
          headers: { Authorization: `Bearer ${this.$store.state.token}` },
        }).then(() => {
          this.fetchMenus();
          this.showForm = false;
          this.resetForm();
          this.$message.success('菜单已更新');
        }).catch((error) => {
          console.error('Error updating menu:', error);
          this.$message.error('更新菜单时出错');
        });
      } else if (!this.isEditing) {
        axios.post(`${process.env.VUE_APP_API_URL}/menus`, this.menu, {
          headers: { Authorization: `Bearer ${this.$store.state.token}` },
        }).then(() => {
          this.fetchMenus();
          this.showForm = false;
          this.resetForm();
          this.$message.success('菜单已添加');
        }).catch((error) => {
          console.error('Error adding menu:', error);
          this.$message.error('添加菜单时出错');
        });
      } else {
        this.$message.info('未检测到任何修改');
      }
    },
    deleteMenu(id) {
      this.$confirm('此操作将永久删除该菜单，是否继续？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(() => {
        axios.delete(`${process.env.VUE_APP_API_URL}/menus/${id}`, {
          headers: { Authorization: `Bearer ${this.$store.state.token}` },
        }).then(() => {
          this.fetchMenus();
          this.$message.success('菜单已删除');
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
      this.menu = {
        title: '',
        icon: '',
        icon_color: '#3f8f94',
        status: 'active',
        order_id: 0,
        parent_id: null,
      };
      this.originalMenu = {};
      this.isEditing = false;
    },
  },
};
</script>

<style scoped>
.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.menu-item__toggle {
  width: 28px;
  height: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: #64748b;
  cursor: pointer;
}

.menu-item__toggle:hover {
  background: #eef8f8;
}

.menu-item__toggle--empty {
  cursor: default;
}

.menu-item__icon {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: #eef8f8;
  overflow: hidden;
}

.menu-item strong {
  display: block;
  font-size: 14px;
  color: #111827;
}

.menu-item p {
  margin: 4px 0 0;
  font-size: 12px;
  color: #94a3b8;
}

.status-pill {
  min-height: 28px;
  padding: 0 10px;
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
}

.status-pill.is-success {
  color: #15803d;
  background: #ecfdf3;
}

.status-pill.is-warning {
  color: #b45309;
  background: #fff7ed;
}

.row-actions {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.full-span {
  grid-column: 1 / -1;
}

.icon-editor {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto auto;
  gap: 10px;
  align-items: center;
  width: 100%;
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
