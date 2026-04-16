<template>
  <div class="login-screen">
    <section class="login-showcase">
      <div class="login-showcase__backdrop"></div>
      <div class="login-showcase__content">
        <p class="login-showcase__eyebrow">Console Access</p>
        <h1>{{ siteTitle || 'ZNav' }}</h1>
        <p class="login-showcase__description">
          一个更轻盈、更稳定的后台入口。把复杂操作收进清晰结构里，让你登录后能立刻进入工作状态。
        </p>

        <div class="login-showcase__highlights">
          <article class="login-showcase__highlight">
            <span class="login-showcase__highlight-index">01</span>
            <div>
              <strong>统一的管理体验</strong>
              <p>应用、菜单、用户与配置维护保持同一套后台语言，减少切换成本。</p>
            </div>
          </article>

          <article class="login-showcase__highlight">
            <span class="login-showcase__highlight-index">02</span>
            <div>
              <strong>更清晰的工作动线</strong>
              <p>筛选、编辑、批量操作和预览都回到更自然的阅读节奏里。</p>
            </div>
          </article>

          <article class="login-showcase__highlight">
            <span class="login-showcase__highlight-index">03</span>
            <div>
              <strong>稳定且可持续扩展</strong>
              <p>界面先保证日常使用舒服，再为后续继续打磨页面留足空间。</p>
            </div>
          </article>
        </div>

        <div class="login-showcase__footer">
          <div class="login-showcase__footer-item">
            <span>聚焦内容</span>
          </div>
          <div class="login-showcase__footer-item">
            <span>统一布局</span>
          </div>
          <div class="login-showcase__footer-item">
            <span>体验更自然</span>
          </div>
        </div>
      </div>
    </section>

    <section class="login-panel">
      <div class="login-panel__card">
        <div class="login-panel__brand">
          <div class="login-panel__logo">
            <img v-if="isUrl(displayIcon)" :src="displayIcon" alt="logo" class="login-panel__logo-image" />
            <icon v-else :icon="displayIcon" :color="displayIconColor" width="26" height="26" />
          </div>
          <div>
            <p class="login-panel__eyebrow">Welcome Back</p>
            <h2>登录管理后台</h2>
          </div>
        </div>

        <div class="login-panel__intro">
          <p class="login-panel__description">请输入账号和密码，进入控制台继续管理站点内容与导航结构。</p>
          <div class="login-panel__badge">安全登录</div>
        </div>

        <el-form ref="loginForm" :model="loginForm" class="login-form" @submit.prevent="login">
          <el-form-item prop="username">
            <div class="login-field">
              <label>用户名</label>
              <el-input v-model="loginForm.username" placeholder="请输入用户名" />
            </div>
          </el-form-item>

          <el-form-item prop="password">
            <div class="login-field">
              <label>密码</label>
              <el-input v-model="loginForm.password" type="password" show-password placeholder="请输入密码" />
            </div>
          </el-form-item>

          <el-button type="primary" native-type="submit" class="login-submit">
            <icon icon="solar:login-2-outline" width="16" height="16" />
            <span>进入后台</span>
          </el-button>
        </el-form>

        <div class="login-panel__footer">
          默认管理员账号通常为 <code>admin</code>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import axios from 'axios';
import { mapGetters } from 'vuex';
import { Icon } from '@iconify/vue';

export default {
  components: {
    Icon,
  },
  data() {
    return {
      loginForm: {
        username: '',
        password: '',
      },
    };
  },
  computed: {
    ...mapGetters(['siteSettings']),
    displayIcon() {
      return this.siteSettings.icon || 'solar:shield-user-outline';
    },
    displayIconColor() {
      return this.siteSettings.icon_color || '#3f8f94';
    },
    siteTitle() {
      return this.siteSettings.title || 'ZNav';
    },
  },
  mounted() {
    if (!this.siteSettings || !this.siteSettings.title) {
      this.$store.dispatch('fetchSiteSettings');
    }
  },
  methods: {
    isUrl(string) {
      return /^(http|https):\/\//.test(string || '');
    },
    login() {
      this.$refs.loginForm.validate((valid) => {
        if (!valid) {
          this.$message.error('请完整填写登录信息');
          return;
        }

        axios.post(`${process.env.VUE_APP_API_URL}/login`, {
          username: this.loginForm.username,
          password: this.loginForm.password,
        }).then((response) => {
          this.$store.commit('setToken', response.data.token);
          this.$router.push('/dashboard');
        }).catch((error) => {
          const errorMsg = error && error.response && error.response.data && error.response.data.error;
          this.$message.error(errorMsg || '登录失败，请稍后重试');
        });
      });
    },
  },
};
</script>

<style scoped>
.login-screen {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1.18fr) minmax(400px, 460px);
  background:
    radial-gradient(circle at top left, rgba(125, 211, 252, 0.12), transparent 24%),
    radial-gradient(circle at bottom right, rgba(94, 234, 212, 0.1), transparent 22%),
    linear-gradient(180deg, #f4f9fb, #fbfefe 52%, #f2f8f9);
}

.login-showcase,
.login-panel {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 36px;
}

.login-showcase__backdrop {
  position: absolute;
  inset: 36px 0 36px 36px;
  border-radius: 36px;
  background:
    radial-gradient(circle at top right, rgba(191, 219, 254, 0.16), transparent 28%),
    linear-gradient(135deg, rgba(89, 171, 176, 0.96), rgba(30, 41, 59, 0.96)),
    #0f172a;
  box-shadow: 0 24px 80px rgba(15, 23, 42, 0.16);
}

.login-showcase__content {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 620px;
  color: #f4feff;
}

.login-showcase__eyebrow,
.login-panel__eyebrow {
  margin: 0 0 10px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.login-showcase__eyebrow {
  color: rgba(216, 249, 248, 0.88);
}

.login-showcase h1 {
  margin: 0;
  max-width: 420px;
  font-size: 56px;
  line-height: 0.95;
}

.login-showcase__description {
  max-width: 560px;
  margin: 20px 0 0;
  font-size: 15px;
  line-height: 1.9;
  color: rgba(226, 232, 240, 0.82);
}

.login-showcase__highlights {
  display: grid;
  gap: 14px;
  margin-top: 32px;
}

.login-showcase__highlight {
  display: grid;
  grid-template-columns: 52px minmax(0, 1fr);
  gap: 14px;
  align-items: start;
  padding: 18px 20px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(10px);
}

.login-showcase__highlight-index {
  width: 52px;
  height: 52px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.1);
  font-size: 13px;
  font-weight: 700;
  color: #dffafa;
}

.login-showcase__highlight strong {
  display: block;
  font-size: 15px;
}

.login-showcase__highlight p {
  margin: 8px 0 0;
  font-size: 13px;
  line-height: 1.75;
  color: rgba(226, 232, 240, 0.8);
}

.login-showcase__footer {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-top: 24px;
}

.login-showcase__footer-item {
  padding: 14px 16px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.06);
}

.login-showcase__footer-item span {
  display: block;
  font-size: 11px;
  color: rgba(216, 249, 248, 0.72);
}

.login-showcase__footer-item strong {
  display: block;
  margin-top: 8px;
  font-size: 15px;
  color: #ffffff;
}

.login-panel__card {
  width: 100%;
  padding: 32px;
  border-radius: 30px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  background: rgba(255, 255, 255, 0.86);
  box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08);
  backdrop-filter: blur(18px);
}

.login-panel__brand {
  display: flex;
  align-items: center;
  gap: 14px;
}

.login-panel__logo {
  width: 58px;
  height: 58px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 18px;
  background: linear-gradient(135deg, rgba(99, 179, 181, 0.1), rgba(125, 211, 252, 0.08));
  overflow: hidden;
}

.login-panel__logo-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.login-panel__eyebrow {
  color: #3f8f94;
}

.login-panel__brand h2 {
  margin: 0;
  font-size: 30px;
  color: #111827;
}

.login-panel__intro {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
  margin-top: 20px;
}

.login-panel__description {
  margin: 0;
  font-size: 14px;
  line-height: 1.8;
  color: #64748b;
}

.login-panel__badge {
  flex-shrink: 0;
  min-height: 34px;
  padding: 0 14px;
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  background: #eaf6f7;
  color: #3f8f94;
  font-size: 12px;
  font-weight: 600;
}

.login-form {
  margin-top: 26px;
}

.login-field {
  width: 100%;
}

.login-field label {
  display: block;
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 600;
  color: #475569;
}

.login-submit {
  width: 100%;
  min-height: 46px;
  margin-top: 10px;
  border: none !important;
  border-radius: 16px;
  background: linear-gradient(135deg, #63b4b1, #478f95) !important;
  box-shadow: 0 16px 32px rgba(71, 143, 149, 0.2);
}

.login-panel__footer {
  margin-top: 20px;
  font-size: 12px;
  color: #94a3b8;
}

.login-panel__footer code {
  padding: 2px 6px;
  border-radius: 999px;
  background: #eaf6f7;
  color: #3f8f94;
}

@media (max-width: 1080px) {
  .login-screen {
    grid-template-columns: 1fr;
  }

  .login-showcase__backdrop {
    inset: 20px;
  }

  .login-showcase,
  .login-panel {
    padding: 20px;
  }
}

@media (max-width: 768px) {
  .login-showcase h1 {
    font-size: 42px;
  }

  .login-showcase__footer {
    grid-template-columns: 1fr;
  }

  .login-showcase__highlight {
    grid-template-columns: 1fr;
  }

  .login-panel__intro {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
