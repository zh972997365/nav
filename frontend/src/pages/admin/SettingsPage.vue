<template>
  <AdminPageShell
    eyebrow="Configuration"
    title="站点设置"
    description="集中维护站点名称、图标、页脚信息与图床 Token，让后台和前台使用同一份基础品牌配置。"
  >
    <template #meta>
      <div class="admin-pill">
        <span>配置状态</span>
        <strong>{{ settings && settings.title ? '已完成' : '待完善' }}</strong>
      </div>
    </template>

    <div class="settings-layout">
      <section class="admin-card settings-preview">
        <div class="admin-card__header">
          <div>
            <h3 class="admin-card__title">配置预览</h3>
            <p class="admin-card__description">展示当前站点标识与已保存的基础信息。</p>
          </div>
        </div>

        <div class="settings-preview__identity">
          <div class="settings-preview__icon">
            <img v-if="isUrl(settings.icon)" :src="settings.icon" alt="icon" class="settings-preview__image" />
            <icon v-else :icon="settings.icon || 'solar:planet-outline'" width="22" height="22" :color="settings.icon_color || '#3f8f94'" />
          </div>
          <div>
            <strong>{{ settings.title || '未设置站点名称' }}</strong>
            <p>{{ settings.footer || '暂无页脚信息' }}</p>
          </div>
        </div>

        <div class="settings-preview__items">
          <div class="settings-preview__item">
            <span>备案号</span>
            <strong>{{ settings.icp || '未填写' }}</strong>
          </div>
          <div class="settings-preview__item">
            <span>图床 Token</span>
            <strong>{{ settings.image_host_token ? '已保存' : '未配置' }}</strong>
          </div>
        </div>
      </section>

      <section class="admin-card">
        <div class="admin-card__header">
          <div>
            <h3 class="admin-card__title">编辑配置</h3>
            <p class="admin-card__description">保存后会同步刷新全局站点信息，并影响登录页与后台品牌展示。</p>
          </div>
        </div>

        <el-form :model="settings" label-width="96px">
          <div class="admin-form-grid">
            <el-form-item label="站点名称" class="full-span">
              <el-input v-model="settings.title" placeholder="请输入站点名称" />
            </el-form-item>

            <el-form-item label="图床 Token" class="full-span">
              <div class="settings-inline">
                <el-input v-model="settings.image_host_token" placeholder="请输入图床 Token" />
                <el-button class="admin-neutral-button" @click="redirectToTokenPage">获取 Token</el-button>
              </div>
            </el-form-item>

            <el-form-item label="图标" class="full-span">
              <div class="settings-inline settings-inline--icon">
                <el-input v-model="settings.icon" placeholder="请输入 Iconify 图标或 URL" />
                <div class="settings-inline__preview">
                  <img v-if="isUrl(settings.icon)" :src="settings.icon" width="24" height="24" alt="icon" />
                  <icon v-else :icon="settings.icon || 'solar:planet-outline'" width="18" height="18" :color="settings.icon_color || '#3f8f94'" />
                </div>
                <el-color-picker v-model="settings.icon_color" />
              </div>
              <a href="https://icon-sets.iconify.design/?category=General" target="_blank" class="admin-icon-link" rel="noreferrer">
                前往 Iconify 选择图标
              </a>
            </el-form-item>

            <el-form-item label="备案号">
              <el-input v-model="settings.icp" placeholder="请输入备案号" />
            </el-form-item>

            <el-form-item label="页脚信息">
              <el-input v-model="settings.footer" type="textarea" placeholder="请输入页脚信息" />
            </el-form-item>
          </div>

          <el-form-item>
            <el-button class="admin-soft-button" type="primary" @click="saveSettings">保存设置</el-button>
          </el-form-item>
        </el-form>
      </section>
    </div>
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
  computed: {
    settings: {
      get() {
        return this.$store.getters.siteSettings;
      },
      set(value) {
        this.$store.commit('setSiteSettings', value);
      },
    },
  },
  mounted() {
    this.$store.dispatch('fetchSiteSettings');
  },
  methods: {
    redirectToTokenPage() {
      window.open('https://img.ink/user/settings.html', '_blank');
    },
    isUrl(string) {
      return !!string && (string.startsWith('http://') || string.startsWith('https://'));
    },
    saveSettings() {
      const settingsToSave = { ...this.settings };

      if (this.settings.image_host_token) {
        axios.post(`${process.env.VUE_APP_API_URL}/save-token`, {
          token: this.settings.image_host_token,
        }).then(() => {
          this.$message.success('Token 已保存');
        }).catch((error) => {
          this.$message.error(`保存 Token 失败：${error.message}`);
        });
      }

      this.$store.dispatch('updateSiteSettings', settingsToSave)
        .then(() => {
          this.$message.success('设置已保存');
        })
        .catch((error) => {
          this.$message.error(`保存设置失败：${error.message}`);
        });
    },
  },
};
</script>

<style scoped>
.settings-layout {
  display: grid;
  grid-template-columns: 360px minmax(0, 1fr);
  gap: 16px;
}

.settings-preview__identity {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 0;
}

.settings-preview__icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: #eef8f8;
  overflow: hidden;
}

.settings-preview__image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.settings-preview__identity strong {
  display: block;
  font-size: 16px;
  color: #111827;
}

.settings-preview__identity p {
  margin: 6px 0 0;
  font-size: 13px;
  color: #64748b;
  line-height: 1.7;
}

.settings-preview__items {
  display: grid;
  gap: 12px;
}

.settings-preview__item {
  padding: 14px;
  border-radius: 14px;
  background: #f8fbfc;
  border: 1px solid #eef4f7;
}

.settings-preview__item span {
  display: block;
  font-size: 12px;
  color: #94a3b8;
}

.settings-preview__item strong {
  display: block;
  margin-top: 6px;
  font-size: 14px;
  color: #111827;
}

.full-span {
  grid-column: 1 / -1;
}

.settings-inline {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 10px;
  align-items: center;
  width: 100%;
}

.settings-inline--icon {
  grid-template-columns: minmax(0, 1fr) auto auto;
}

.settings-inline__preview {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: #eef8f8;
}

@media (max-width: 1080px) {
  .settings-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .settings-inline,
  .settings-inline--icon {
    grid-template-columns: 1fr;
  }
}
</style>
