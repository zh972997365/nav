<template>
  <div class="login-screen">
    <div class="login-layout">
      <section class="login-stage">
        <div class="login-stage__frame">
          <div class="login-stage__glow login-stage__glow--top"></div>
          <div class="login-stage__glow login-stage__glow--bottom"></div>

          <header class="login-stage__header">
            <div class="login-stage__brand">
              <div class="login-stage__logo">
                <img v-if="isUrl(displayIcon)" :src="displayIcon" alt="logo" class="login-stage__logo-image" />
                <icon v-else :icon="displayIcon" :color="displayIconColor" width="28" height="28" />
              </div>
              <div class="login-stage__brand-copy">
                <p class="login-stage__eyebrow">Console Access</p>
                <strong>{{ siteTitle || 'ZNav' }}</strong>
              </div>
            </div>

            <div class="login-stage__status">
              <span class="login-stage__status-dot"></span>
              <span>Secure Entry</span>
            </div>
          </header>

          <div class="login-stage__body">
            <section class="login-copy">
              <div class="login-copy__hero">
                <p class="login-copy__eyebrow">管理入口</p>
                <div class="login-copy__divider"></div>
              </div>

              <div class="login-copy__overview">
                <div class="login-copy__overview-main">
                  <span class="login-copy__overview-label">后台工作台</span>
                  <strong>统一入口</strong>
                  <p>维护站点信息，让后台操作保持清晰、稳定和高效。</p>
                </div>
                <div class="login-copy__overview-side">
                  <span>安全访问</span>
                  <span>统一入口</span>
                </div>
              </div>

              <div class="login-copy__feature-list">
                <article class="login-copy__feature">
                  <span class="login-copy__feature-index">01</span>
                  <div>
                    <strong>应用管理</strong>
                    <p>维护应用信息、入口链接与展示内容。</p>
                  </div>
                </article>
                <article class="login-copy__feature">
                  <span class="login-copy__feature-index">02</span>
                  <div>
                    <strong>菜单管理</strong>
                    <p>调整导航层级、分组与页面跳转关系。</p>
                  </div>
                </article>
                <article class="login-copy__feature">
                  <span class="login-copy__feature-index">03</span>
                  <div>
                    <strong>账号与设置</strong>
                    <p>管理后台账号，并维护系统基础配置。</p>
                  </div>
                </article>
              </div>
            </section>

            <section class="login-panel">
              <div class="login-panel__surface">
                <div class="login-panel__head">
                  <div class="login-panel__title">
                    <div class="login-panel__title-logo">
                      <img v-if="isUrl(displayIcon)" :src="displayIcon" alt="logo" class="login-panel__title-logo-image" />
                      <icon v-else :icon="displayIcon" :color="displayIconColor" width="28" height="28" />
                    </div>
                    <div class="login-panel__title-copy">
                      <p class="login-panel__eyebrow">Welcome Back</p>
                      <h2>登录管理后台</h2>
                    </div>
                  </div>
                </div>

                <div class="login-panel__meta">
                  <div class="login-panel__meta-item">
                    <span class="login-panel__meta-label">状态</span>
                    <strong>安全登录</strong>
                  </div>
                  <div class="login-panel__meta-item">
                    <span class="login-panel__meta-label">默认账号</span>
                    <strong><code>admin</code></strong>
                  </div>
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
                      <el-input
                        v-model="loginForm.password"
                        type="password"
                        show-password
                        placeholder="请输入密码"
                      />
                    </div>
                  </el-form-item>

                  <el-button type="primary" native-type="submit" class="login-submit">
                    <icon icon="solar:login-2-outline" width="16" height="16" />
                    <span>进入后台</span>
                  </el-button>
                </el-form>

                <footer class="login-panel__footer">
                  <span>登录后进入后台控制台</span>
                  <span>请使用具备管理权限的账号操作</span>
                </footer>
              </div>
            </section>
          </div>
        </div>
      </section>
    </div>
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
  overflow: hidden;
  box-sizing: border-box;
  padding: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background:
    radial-gradient(circle at top left, rgba(125, 211, 252, 0.14), transparent 24%),
    radial-gradient(circle at bottom right, rgba(94, 234, 212, 0.12), transparent 26%),
    linear-gradient(180deg, #f4f9fb, #fbfefe 52%, #f2f8f9);
}

.login-layout {
  width: min(1280px, 100%);
  margin: 0 auto;
}

.login-stage {
  width: 100%;
}

.login-stage__frame {
  position: relative;
  min-height: min(840px, calc(100vh - 48px));
  overflow: hidden;
  border-radius: 36px;
  border: 1px solid rgba(255, 255, 255, 0.82);
  background:
    linear-gradient(145deg, rgba(16, 32, 51, 0.96), rgba(46, 87, 108, 0.95) 46%, rgba(244, 249, 251, 0.92) 46.1%, rgba(255, 255, 255, 0.9) 100%);
  box-shadow: 0 32px 90px rgba(15, 23, 42, 0.1);
  backdrop-filter: blur(18px);
}

.login-stage__glow {
  position: absolute;
  border-radius: 999px;
  filter: blur(4px);
  pointer-events: none;
}

.login-stage__glow--top {
  top: -120px;
  left: -80px;
  width: 280px;
  height: 280px;
  background: radial-gradient(circle, rgba(110, 231, 231, 0.24), transparent 70%);
}

.login-stage__glow--bottom {
  right: -120px;
  bottom: -140px;
  width: 360px;
  height: 360px;
  background: radial-gradient(circle, rgba(63, 143, 148, 0.12), transparent 72%);
}

.login-stage__header {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 24px 28px 0;
}

.login-stage__brand {
  display: inline-flex;
  align-items: center;
  gap: 14px;
}

.login-stage__logo {
  width: 58px;
  height: 58px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.12);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.08);
  overflow: hidden;
}

.login-stage__logo-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.login-stage__brand-copy {
  display: flex;
  flex-direction: column;
  gap: 5px;
  color: #ffffff;
}

.login-stage__eyebrow,
.login-copy__eyebrow,
.login-panel__eyebrow {
  margin: 0;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.login-stage__eyebrow {
  color: rgba(216, 249, 248, 0.78);
}

.login-stage__brand-copy strong {
  font-size: 20px;
  line-height: 1;
}

.login-stage__status {
  min-height: 38px;
  padding: 0 14px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.12);
  color: #7ce7dd;
  font-size: 12px;
}

.login-stage__status-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: #7ce7dd;
  box-shadow: 0 0 0 4px rgba(124, 231, 221, 0.18);
}

.login-stage__body {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: minmax(0, 1.06fr) minmax(380px, 440px);
  align-items: center;
  gap: 28px;
  padding: 24px 28px 28px;
  min-height: 0;
}

.login-copy {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 22px 16px 12px 8px;
  color: #f4feff;
  align-self: center;
}

.login-copy__eyebrow {
  color: rgba(216, 249, 248, 0.72);
}

.login-copy__hero h1 {
  margin: 14px 0 0;
  max-width: 520px;
  font-size: clamp(36px, 3.5vw, 52px);
  line-height: 1.04;
}

.login-copy__divider {
  width: 88px;
  height: 4px;
  margin-top: 22px;
  border-radius: 999px;
  background: linear-gradient(90deg, rgba(124, 231, 221, 0.9), rgba(124, 231, 221, 0.12));
}

.login-copy__overview {
  max-width: 560px;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 18px;
  margin-top: 30px;
  padding: 20px 22px;
  border-radius: 26px;
  background:
    radial-gradient(circle at top right, rgba(124, 231, 221, 0.16), transparent 30%),
    linear-gradient(180deg, rgba(104, 191, 187, 0.22), rgba(18, 35, 53, 0.3));
  border: 1px solid rgba(210, 242, 244, 0.12);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.04);
  backdrop-filter: blur(12px);
  transition: transform 0.25s ease, box-shadow 0.25s ease, border-color 0.25s ease, background 0.25s ease;
}

.login-copy__overview:hover {
  transform: translateY(-3px);
  border-color: rgba(210, 242, 244, 0.2);
  box-shadow: 0 18px 36px rgba(15, 23, 42, 0.14), inset 0 1px 0 rgba(255, 255, 255, 0.05);
  background:
    radial-gradient(circle at top right, rgba(124, 231, 221, 0.2), transparent 30%),
    linear-gradient(180deg, rgba(104, 191, 187, 0.26), rgba(18, 35, 53, 0.34));
}

.login-copy__overview-label {
  display: inline-flex;
  align-items: center;
  width: fit-content;
  min-height: 28px;
  padding: 0 12px;
  margin-bottom: 14px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.1);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #dffafa;
}

.login-copy__overview-main strong {
  display: block;
  max-width: 420px;
  font-size: 24px;
  line-height: 1.3;
  color: #ffffff;
}

.login-copy__overview-main p {
  max-width: 430px;
  margin: 10px 0 0;
  font-size: 13px;
  line-height: 1.75;
  color: rgba(226, 232, 240, 0.8);
}

.login-copy__overview-side {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  gap: 10px;
}

.login-copy__overview-side span {
  min-height: 34px;
  padding: 0 14px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.1);
  color: rgba(232, 250, 250, 0.9);
  font-size: 12px;
}

.login-copy__feature-list {
  max-width: 560px;
  display: grid;
  gap: 14px;
  margin-top: 18px;
}

.login-copy__feature {
  display: grid;
  grid-template-columns: 56px minmax(0, 1fr);
  gap: 16px;
  align-items: center;
  padding: 16px 18px;
  border-radius: 22px;
  background: rgba(18, 35, 53, 0.22);
  border: 1px solid rgba(210, 242, 244, 0.08);
  backdrop-filter: blur(10px);
  transition: transform 0.22s ease, background 0.22s ease, border-color 0.22s ease, box-shadow 0.22s ease;
}

.login-copy__feature:hover {
  transform: translateY(-2px);
  background: rgba(23, 43, 64, 0.34);
  border-color: rgba(210, 242, 244, 0.14);
  box-shadow: 0 14px 28px rgba(15, 23, 42, 0.12);
}

.login-copy__feature-index {
  width: 56px;
  height: 56px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.1);
  color: #dffafa;
  font-size: 13px;
  font-weight: 700;
}

.login-copy__feature strong {
  display: block;
  font-size: 16px;
  color: #ffffff;
}

.login-copy__feature p {
  margin: 6px 0 0;
  font-size: 13px;
  line-height: 1.7;
  color: rgba(226, 232, 240, 0.8);
}

.login-panel {
  display: flex;
  align-items: center;
  align-self: center;
}

.login-panel__surface {
  width: 100%;
  padding: 34px 30px 26px;
  border-radius: 30px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.94), rgba(247, 251, 252, 0.9));
  box-shadow: 0 26px 64px rgba(15, 23, 42, 0.1);
  backdrop-filter: blur(20px);
  transition: transform 0.25s ease, box-shadow 0.25s ease;
}

.login-panel__surface:hover {
  transform: translateY(-2px);
  box-shadow: 0 30px 72px rgba(15, 23, 42, 0.12);
}

.login-panel__head {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.login-panel__title {
  display: flex;
  align-items: center;
  gap: 16px;
}

.login-panel__title-logo {
  width: 58px;
  height: 58px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 18px;
  background: linear-gradient(180deg, rgba(234, 246, 247, 0.98), rgba(242, 250, 250, 0.92));
  box-shadow: inset 0 0 0 1px rgba(63, 143, 148, 0.08);
  overflow: hidden;
}

.login-panel__title-logo-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.login-panel__title-copy {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.login-panel__eyebrow {
  color: #3f8f94;
}

.login-panel__title-copy h2 {
  margin: 0;
  font-size: 34px;
  line-height: 1.08;
  color: #111827;
}

.login-panel__meta {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 24px;
}

.login-panel__meta-item {
  padding: 14px 16px;
  border-radius: 18px;
  background: linear-gradient(180deg, rgba(234, 246, 247, 0.96), rgba(244, 250, 251, 0.96));
  border: 1px solid #dcebed;
  transition: transform 0.22s ease, border-color 0.22s ease, box-shadow 0.22s ease, background 0.22s ease;
}

.login-panel__meta-item:hover {
  transform: translateY(-2px);
  border-color: #c6dde1;
  background: linear-gradient(180deg, rgba(236, 248, 249, 0.98), rgba(248, 252, 252, 0.96));
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.06);
}

.login-panel__meta-label {
  display: block;
  margin-bottom: 8px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #7a8b9c;
}

.login-panel__meta-item strong {
  font-size: 14px;
  color: #102033;
}

.login-panel__meta-item code {
  padding: 2px 6px;
  border-radius: 999px;
  background: rgba(63, 143, 148, 0.12);
  color: #3f8f94;
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
  min-height: 48px;
  margin-top: 12px;
  border: none !important;
  border-radius: 18px;
  background: linear-gradient(135deg, #63b4b1, #478f95) !important;
  box-shadow: 0 16px 32px rgba(71, 143, 149, 0.2);
  transition: transform 0.22s ease, box-shadow 0.22s ease, filter 0.22s ease;
}

.login-submit:hover,
.login-submit:focus {
  transform: translateY(-1px);
  box-shadow: 0 20px 36px rgba(71, 143, 149, 0.24);
  filter: brightness(1.02);
}

.login-panel__footer {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
  margin-top: 18px;
  padding-top: 18px;
  border-top: 1px solid #e7eef1;
  font-size: 12px;
  color: #8b9aaa;
}

@media (max-width: 1120px) {
  .login-screen {
    height: auto;
    min-height: 100vh;
    overflow: visible;
  }

  .login-stage__frame {
    min-height: auto;
  }

  .login-stage__body {
    grid-template-columns: 1fr;
  }

  .login-copy {
    padding: 8px 8px 0;
  }
}

@media (max-width: 768px) {
  .login-screen {
    padding: 14px;
  }

  .login-stage__frame {
    border-radius: 28px;
    background:
      linear-gradient(180deg, rgba(16, 32, 51, 0.96), rgba(46, 87, 108, 0.94) 34%, rgba(244, 249, 251, 0.92) 34.1%, rgba(255, 255, 255, 0.9) 100%);
  }

  .login-stage__header,
  .login-stage__body {
    padding-left: 16px;
    padding-right: 16px;
  }

  .login-stage__header {
    flex-direction: column;
    align-items: flex-start;
    padding-top: 16px;
  }

  .login-stage__body {
    gap: 18px;
    padding-top: 18px;
    padding-bottom: 16px;
  }

  .login-copy {
    padding: 12px 0 0;
  }

  .login-copy__hero h1 {
    font-size: 34px;
  }

  .login-copy__overview {
    grid-template-columns: 1fr;
  }

  .login-panel__surface {
    padding: 24px 20px 20px;
    border-radius: 24px;
  }

  .login-panel__title {
    align-items: flex-start;
  }

  .login-panel__title-copy h2 {
    font-size: 28px;
  }

  .login-panel__meta {
    grid-template-columns: 1fr;
  }
}
</style>
