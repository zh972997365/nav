<template>
  <AdminPageShell
    eyebrow="Overview"
    title="仪表盘"
    description="后台首页聚焦最关键的概览信息、规模占比和最近新增内容，帮助你快速进入状态。"
  >
    <template #meta>
      <div class="admin-pill">
        <span>今日总览</span>
        <strong>{{ totalApplications + totalMenus + totalUsers }}</strong>
      </div>
    </template>

    <div class="dashboard-layout">
      <section class="dashboard-hero-card">
        <div class="dashboard-hero-card__copy">
          <p class="dashboard-hero-card__eyebrow">System summary</p>
          <h3>在一个页面里，快速掌握后台当前的整体状态。</h3>
          <p>这里保留最值得先看的概览数据、结构规模和最近新增内容，帮助你进入后台后马上找到重点。</p>
        </div>
        <div class="dashboard-hero-card__stats">
          <div class="dashboard-mini-stat">
            <span>应用</span>
            <strong>{{ totalApplications }}</strong>
          </div>
          <div class="dashboard-mini-stat">
            <span>菜单</span>
            <strong>{{ totalMenus }}</strong>
          </div>
          <div class="dashboard-mini-stat">
            <span>用户</span>
            <strong>{{ totalUsers }}</strong>
          </div>
        </div>
      </section>

      <section class="dashboard-panel">
        <div class="dashboard-panel__head">
          <div>
            <h3>应用占比</h3>
            <p>用于快速感知当前应用规模的整体分布</p>
          </div>
          <span class="dashboard-panel__value">{{ applicationsPercentage }}%</span>
        </div>
        <div ref="applicationsChart" class="dashboard-chart"></div>
      </section>

      <section class="dashboard-panel">
        <div class="dashboard-panel__head">
          <div>
            <h3>菜单占比</h3>
            <p>帮助判断导航结构目前的复杂度与层级体量</p>
          </div>
          <span class="dashboard-panel__value">{{ menusPercentage }}%</span>
        </div>
        <div ref="menusChart" class="dashboard-chart"></div>
      </section>

      <section class="dashboard-panel">
        <div class="dashboard-panel__head">
          <div>
            <h3>用户占比</h3>
            <p>查看后台账户数量变化带来的管理压力</p>
          </div>
          <span class="dashboard-panel__value">{{ usersPercentage }}%</span>
        </div>
        <div ref="usersChart" class="dashboard-chart"></div>
      </section>
    </div>

    <section class="admin-table-card dashboard-table">
      <div class="admin-card__header">
        <div>
          <h3 class="admin-card__title">最近新增应用</h3>
          <p class="admin-card__description">按创建时间倒序展示，方便你快速确认最近录入或同步进后台的内容。</p>
        </div>
      </div>

      <el-table :data="recentApplications" border stripe :max-height="420">
        <el-table-column label="应用" min-width="240">
          <template #default="{ row }">
            <div class="dashboard-item">
              <div class="dashboard-item__icon">
                <img v-if="isUrl(row.icon)" :src="row.icon" alt="icon" width="22" height="22" />
                <icon v-else :icon="row.icon" :color="row.icon_color || '#3f8f94'" width="18" height="18" />
              </div>
              <div>
                <strong>{{ row.title }}</strong>
                <el-tooltip :content="row.description || '暂无描述'" placement="top-start" :show-after="250">
                  <p class="dashboard-item__description">{{ row.description || '暂无描述' }}</p>
                </el-tooltip>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="link" label="链接" min-width="220" show-overflow-tooltip />
        <el-table-column label="菜单" min-width="180">
          <template #default="{ row }">
            <div class="dashboard-menu">
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
            <span class="dashboard-status" :class="row.status === 'active' ? 'is-active' : 'is-inactive'">
              {{ row.status === 'active' ? '启用中' : '停用中' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="190">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
      </el-table>
    </section>
  </AdminPageShell>
</template>

<script>
import axios from 'axios';
import * as echarts from 'echarts';
import { Icon } from '@iconify/vue';
import AdminPageShell from '../../components/admin/AdminPageShell.vue';

export default {
  components: {
    Icon,
    AdminPageShell,
  },
  data() {
    return {
      totalUsers: 0,
      totalApplications: 0,
      totalMenus: 0,
      recentApplications: [],
      menus: [],
      applicationsPercentage: 0,
      menusPercentage: 0,
      usersPercentage: 0,
      charts: {},
    };
  },
  mounted() {
    this.fetchDashboardData();
  },
  beforeUnmount() {
    Object.values(this.charts).forEach((chart) => {
      if (chart) chart.dispose();
    });
  },
  methods: {
    createRingChart(el, value, color, key) {
      if (!el) return;
      if (this.charts[key]) {
        this.charts[key].dispose();
      }

      const chart = echarts.init(el);
      chart.setOption({
        tooltip: {
          trigger: 'item',
          formatter: '{c}%',
        },
        series: [
          {
            type: 'pie',
            radius: ['70%', '84%'],
            label: { show: false },
            silent: true,
            data: [
              { value, itemStyle: { color } },
              { value: 100 - value, itemStyle: { color: '#e8f2f4' } },
            ],
          },
        ],
      });
      this.charts[key] = chart;
    },
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
    calculatePercentage(value) {
      const maxValue = 500;
      return parseFloat(Math.min((value / maxValue) * 100, 100).toFixed(1));
    },
    fetchDashboardData() {
      const token = this.$store.state.token;
      axios.all([
        axios.get(`${process.env.VUE_APP_API_URL}/users/total`, { headers: { Authorization: `Bearer ${token}` } }),
        axios.get(`${process.env.VUE_APP_API_URL}/applications/total`, { headers: { Authorization: `Bearer ${token}` } }),
        axios.get(`${process.env.VUE_APP_API_URL}/menus/total`, { headers: { Authorization: `Bearer ${token}` } }),
        axios.get(`${process.env.VUE_APP_API_URL}/applications/recent`, { headers: { Authorization: `Bearer ${token}` } }),
        axios.get(`${process.env.VUE_APP_API_URL}/menus`, {
          params: { pageSize: 1000 },
          headers: { Authorization: `Bearer ${token}` },
        }),
      ]).then(axios.spread((usersRes, appsRes, menusRes, recentRes, allMenusRes) => {
        this.totalUsers = usersRes.data.total;
        this.totalApplications = appsRes.data.total;
        this.totalMenus = menusRes.data.total;
        this.recentApplications = recentRes.data.applications;
        this.menus = allMenusRes.data.menus;

        this.applicationsPercentage = this.calculatePercentage(this.totalApplications);
        this.menusPercentage = this.calculatePercentage(this.totalMenus);
        this.usersPercentage = this.calculatePercentage(this.totalUsers);

        this.createRingChart(this.$refs.applicationsChart, this.applicationsPercentage, '#3f8f94', 'apps');
        this.createRingChart(this.$refs.menusChart, this.menusPercentage, '#0f172a', 'menus');
        this.createRingChart(this.$refs.usersChart, this.usersPercentage, '#15803d', 'users');
      })).catch((error) => {
        console.error('Error fetching dashboard data:', error);
      });
    },
  },
};
</script>

<style scoped>
.dashboard-layout {
  display: grid;
  grid-template-columns: 1.25fr repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.dashboard-hero-card,
.dashboard-panel {
  border: 1px solid rgba(255, 255, 255, 0.78);
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.04);
}

.dashboard-hero-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 20px;
  background: linear-gradient(135deg, rgba(89, 171, 176, 0.94), rgba(30, 41, 59, 0.94));
  color: #f4feff;
}

.dashboard-hero-card__eyebrow {
  margin: 0 0 10px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: rgba(222, 248, 246, 0.82);
}

.dashboard-hero-card h3 {
  margin: 0;
  font-size: 24px;
  line-height: 1.25;
}

.dashboard-hero-card__copy p:last-child {
  margin: 12px 0 0;
  font-size: 13px;
  line-height: 1.75;
  color: rgba(240, 253, 252, 0.78);
}

.dashboard-hero-card__stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  margin-top: 18px;
}

.dashboard-mini-stat {
  padding: 12px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.08);
}

.dashboard-mini-stat span {
  display: block;
  font-size: 12px;
  color: rgba(204, 251, 241, 0.72);
}

.dashboard-mini-stat strong {
  display: block;
  margin-top: 6px;
  font-size: 22px;
  color: #ffffff;
}

.dashboard-panel {
  padding: 16px;
}

.dashboard-panel__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 8px;
}

.dashboard-panel__head h3 {
  margin: 0;
  font-size: 15px;
}

.dashboard-panel__head p {
  margin: 6px 0 0;
  font-size: 12px;
  color: #64748b;
}

.dashboard-panel__value {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
}

.dashboard-chart {
  width: 100%;
  height: 160px;
}

.dashboard-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.dashboard-item__icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: #eef8f8;
  overflow: hidden;
}

.dashboard-item strong {
  display: block;
  font-size: 14px;
  color: #111827;
}

.dashboard-item__description {
  margin: 4px 0 0;
  font-size: 12px;
  color: #94a3b8;
  display: -webkit-box;
  overflow: hidden;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
  text-overflow: ellipsis;
  cursor: help;
  max-width: 260px;
}

.dashboard-menu {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.dashboard-status {
  min-height: 28px;
  padding: 0 10px;
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
}

.dashboard-status.is-active {
  color: #15803d;
  background: #ecfdf3;
}

.dashboard-status.is-inactive {
  color: #b45309;
  background: #fff7ed;
}

@media (max-width: 1180px) {
  .dashboard-layout {
    grid-template-columns: 1fr 1fr;
  }

  .dashboard-hero-card {
    grid-column: 1 / -1;
  }
}

@media (max-width: 768px) {
  .dashboard-layout,
  .dashboard-hero-card__stats {
    grid-template-columns: 1fr;
  }
}
</style>
