<template>
  <AdminPageShell
    eyebrow="Access"
    title="用户管理"
    description="维护后台账户、权限角色与启停状态，并支持管理员编辑用户信息或普通成员修改自己的密码。"
  >
    <template #meta>
      <div class="admin-pill">
        <span>账户总数</span>
        <strong>{{ users.length }}</strong>
      </div>
    </template>

    <section class="admin-workspace">
      <div class="admin-workspace__top">
        <div class="admin-workspace__summary">
          <article class="admin-simple-metric">
            <p class="admin-simple-metric__eyebrow">全部账户</p>
            <p class="admin-simple-metric__value">{{ users.length }}</p>
            <p class="admin-simple-metric__hint">系统中的后台账号总数</p>
          </article>
          <article class="admin-simple-metric">
            <p class="admin-simple-metric__eyebrow">管理员</p>
            <p class="admin-simple-metric__value">{{ adminCount }}</p>
            <p class="admin-simple-metric__hint">拥有更高权限的账号数量</p>
          </article>
          <article class="admin-simple-metric">
            <p class="admin-simple-metric__eyebrow">启用中</p>
            <p class="admin-simple-metric__value">{{ enabledCount }}</p>
            <p class="admin-simple-metric__hint">当前仍可登录后台的账号</p>
          </article>
        </div>
      </div>

      <div class="admin-workspace__actions">
        <div class="admin-workspace__actions-copy">
          用户页保留表格为核心，角色与状态可以直接在行内切换；如果你不是管理员，只会看到自己可操作的部分。
        </div>
        <div class="admin-workspace__actions-buttons">
          <el-button v-if="isAdmin" type="primary" class="admin-soft-button" @click="openUserForm('add')">新增用户</el-button>
        </div>
      </div>

      <div class="admin-workspace__table">
        <el-table v-if="users.length > 0" :data="users" border stripe style="width: 100%" :max-height="520">
          <el-table-column label="用户" min-width="220">
            <template #default="{ row }">
              <div class="user-item">
                <div class="user-item__avatar">{{ row.username ? row.username.slice(0, 1).toUpperCase() : 'U' }}</div>
                <div>
                  <strong>{{ row.username }}</strong>
                  <p>ID #{{ row.id }}</p>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="角色" width="130" align="center">
            <template #default="{ row }">
              <button
                type="button"
                class="toggle-pill"
                :class="row.is_admin ? 'is-admin' : 'is-member'"
                @click="toggleUserAdmin(row)"
              >
                {{ row.is_admin ? '管理员' : '成员' }}
              </button>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="130" align="center">
            <template #default="{ row }">
              <button
                type="button"
                class="toggle-pill"
                :class="row.status === 'enabled' ? 'is-enabled' : 'is-disabled'"
                @click="toggleUserStatus(row)"
              >
                {{ row.status === 'enabled' ? '启用中' : '已停用' }}
              </button>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="190">
            <template #default="{ row }">
              {{ new Date(row.created_at).toLocaleString() }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" align="center" fixed="right">
            <template #default="{ row }">
              <div v-if="isAdmin || row.username === loggedInUser" class="row-actions">
                <el-button link class="admin-link-button admin-link-button--edit" @click="openUserForm('edit', row)">
                  <icon icon="solar:pen-new-square-outline" width="16" height="16" />
                </el-button>
                <el-button
                  v-if="isAdmin && row.username !== loggedInUser"
                  link
                  class="admin-link-button admin-link-button--delete"
                  @click="handleDelete(row)"
                >
                  <icon icon="solar:trash-bin-trash-outline" width="16" height="16" />
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <el-empty v-else description="暂无用户数据" />
      </div>
    </section>

    <el-dialog :title="dialogTitle" v-model="showForm" width="640px" class="admin-dialog">
      <template v-if="isEditing">
        <el-form :model="formModel" :rules="rules" ref="form" label-width="96px">
          <template v-if="isAdmin">
            <el-form-item label="用户名">
              <el-input v-model="currentUser.username" disabled />
            </el-form-item>
            <el-form-item label="新密码" prop="password">
              <el-input v-model="currentUser.password" type="password" show-password placeholder="如需修改密码，请输入新密码" />
            </el-form-item>
            <el-form-item label="管理员">
              <el-switch v-model="currentUser.is_admin" :active-value="true" :inactive-value="false" />
            </el-form-item>
            <el-form-item label="状态">
              <el-switch v-model="currentUser.status" :active-value="'enabled'" :inactive-value="'disabled'" />
            </el-form-item>
          </template>

          <template v-else>
            <el-form-item label="原密码" prop="currentPassword">
              <el-input v-model="passwordForm.currentPassword" type="password" show-password placeholder="请输入原密码" />
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
              <el-input v-model="passwordForm.newPassword" type="password" show-password placeholder="请输入新密码" />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input v-model="passwordForm.confirmPassword" type="password" show-password placeholder="请再次输入新密码" />
            </el-form-item>
          </template>
        </el-form>
      </template>

      <template v-else>
        <el-form :model="currentUser" :rules="rules" ref="form" label-width="96px">
          <el-form-item label="用户名" prop="username">
            <el-input v-model="currentUser.username" placeholder="请输入用户名" />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input v-model="currentUser.password" type="password" show-password placeholder="请输入密码" />
          </el-form-item>
          <el-form-item v-if="isAdmin" label="管理员">
            <el-switch v-model="currentUser.is_admin" :active-value="true" :inactive-value="false" />
          </el-form-item>
          <el-form-item v-if="isAdmin" label="状态">
            <el-switch v-model="currentUser.status" :active-value="'enabled'" :inactive-value="'disabled'" />
          </el-form-item>
        </el-form>
      </template>

      <template #footer>
        <div class="dialog-footer">
          <el-button class="admin-neutral-button" @click="showForm = false">取消</el-button>
          <el-button type="primary" @click="submitUser">保存</el-button>
        </div>
      </template>
    </el-dialog>
  </AdminPageShell>
</template>

<script>
import axios from 'axios';
import { Icon } from '@iconify/vue';
import AdminPageShell from '../../components/admin/AdminPageShell.vue';

export default {
  components: {
    Icon,
    AdminPageShell,
  },
  data() {
    return {
      users: [],
      currentUser: {
        id: null,
        username: '',
        password: '',
        status: 'enabled',
        is_admin: false,
      },
      passwordForm: {
        currentPassword: '',
        newPassword: '',
        confirmPassword: '',
      },
      isEditing: false,
      showForm: false,
      originalUserData: null,
      loggedInUser: this.$store.state.username,
      isAdmin: false,
      rules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { pattern: /^[a-zA-Z0-9]+$/, message: '用户名只能包含英文字母和数字', trigger: 'blur' },
        ],
        password: [
          { min: 5, message: '密码至少需要 5 个字符', trigger: 'blur' },
        ],
        currentPassword: [
          { required: true, message: '请输入原密码', trigger: 'blur' },
        ],
        newPassword: [
          { required: true, message: '请输入新密码', trigger: 'blur' },
          { min: 5, message: '新密码至少为 5 个字符', trigger: 'blur' },
        ],
        confirmPassword: [
          { required: true, message: '请确认新密码', trigger: 'blur' },
          { validator: this.validateConfirmPassword, trigger: 'blur' },
        ],
      },
    };
  },
  computed: {
    adminCount() {
      return this.users.filter((item) => item.is_admin).length;
    },
    enabledCount() {
      return this.users.filter((item) => item.status === 'enabled').length;
    },
    dialogTitle() {
      if (this.isEditing) {
        return this.isAdmin ? '编辑用户' : '修改密码';
      }
      return '新增用户';
    },
    formModel() {
      return this.isEditing ? (this.isAdmin ? this.currentUser : this.passwordForm) : this.currentUser;
    },
  },
  created() {
    this.getLoggedInUser();
  },
  methods: {
    getLoggedInUser() {
      axios.get(`${process.env.VUE_APP_API_URL}/user`, {
        headers: { Authorization: `Bearer ${this.$store.state.token}` },
      }).then((response) => {
        this.loggedInUser = response.data.user.username;
        this.isAdmin = response.data.user.is_admin;
        this.fetchUsers();
      }).catch((error) => {
        console.error('Failed to fetch logged in user:', error);
      });
    },
    fetchUsers() {
      axios.get(`${process.env.VUE_APP_API_URL}/users`, {
        headers: { Authorization: `Bearer ${this.$store.state.token}` },
      }).then((response) => {
        this.users = response.data.users;
      }).catch((error) => {
        console.error('Failed to fetch users:', error);
      });
    },
    openUserForm(action, user) {
      this.isEditing = action === 'edit';
      if (this.isEditing) {
        this.currentUser = { ...user, password: '' };
        this.originalUserData = { ...user };
        this.passwordForm = {
          currentPassword: '',
          newPassword: '',
          confirmPassword: '',
        };
      } else {
        this.currentUser = {
          id: null,
          username: '',
          password: '',
          status: 'enabled',
          is_admin: false,
        };
        this.originalUserData = null;
      }
      this.showForm = true;
    },
    validateConfirmPassword(rule, value, callback) {
      if (value !== this.passwordForm.newPassword) {
        callback(new Error('两次输入的新密码不匹配'));
      } else {
        callback();
      }
    },
    submitUser() {
      this.$refs.form.validate((valid) => {
        if (!valid) {
          this.$message.error('请正确填写表单');
          return;
        }

        let url;
        let method;
        let actionSuccess;
        let actionFail;
        let data;

        if (this.isEditing) {
          if (this.isAdmin) {
            url = `${process.env.VUE_APP_API_URL}/users`;
            method = 'put';
            actionSuccess = '已修改';
            actionFail = '修改失败';
            data = { ...this.currentUser };
            if (!data.password) delete data.password;
          } else {
            url = `${process.env.VUE_APP_API_URL}/users/password`;
            method = 'put';
            actionSuccess = '密码已修改';
            actionFail = '密码修改失败';
            data = {
              id: this.currentUser.id,
              current_password: this.passwordForm.currentPassword,
              new_password: this.passwordForm.newPassword,
            };
          }
        } else {
          url = `${process.env.VUE_APP_API_URL}/users`;
          method = 'post';
          actionSuccess = '已添加';
          actionFail = '添加失败';
          data = { ...this.currentUser };
        }

        axios[method](url, data, {
          headers: { Authorization: `Bearer ${this.$store.state.token}` },
        }).then(() => {
          this.fetchUsers();
          this.showForm = false;

          if (this.isAdmin && this.isEditing) {
            const messages = [];
            if (this.currentUser.password) messages.push(`${this.currentUser.username} 密码已修改`);
            if (this.currentUser.is_admin !== this.originalUserData.is_admin) {
              messages.push(`${this.currentUser.username} ${this.currentUser.is_admin ? '已授予管理员权限' : '已取消管理员权限'}`);
            }
            if (this.currentUser.status !== this.originalUserData.status) {
              messages.push(`${this.currentUser.username} ${this.currentUser.status === 'enabled' ? '已启用' : '已停用'}`);
            }
            if (!messages.length) messages.push(`${this.currentUser.username} 信息未修改`);
            messages.forEach((msg) => this.$message.success(msg));
          } else {
            this.$message.success(`${this.currentUser.username} ${actionSuccess}`);
          }
        }).catch((error) => {
          const errorMsg = error && error.response && error.response.data && error.response.data.error;
          if (errorMsg === '原密码不正确') {
            this.$message.error(`${this.currentUser.username} 原密码不匹配`);
          } else if (errorMsg) {
            this.$message.error(`${this.currentUser.username}: ${errorMsg}`);
          } else {
            this.$message.error(`${this.currentUser.username} 用户${actionFail}`);
          }
        });
      });
    },
    toggleUserStatus(user) {
      if (!this.isAdmin) {
        this.$message.error('您没有权限修改用户状态');
        return;
      }
      if (this.isActionDisabled(user)) {
        this.$message.error('不允许修改自己的状态');
        return;
      }

      user.status = user.status === 'enabled' ? 'disabled' : 'enabled';
      axios.put(`${process.env.VUE_APP_API_URL}/users/status`, user, {
        headers: { Authorization: `Bearer ${this.$store.state.token}` },
      }).then(() => {
        this.$message.success(`${user.username} ${user.status === 'enabled' ? '已启用' : '已停用'}`);
      }).catch((error) => {
        const errorMsg = error && error.response && error.response.data && error.response.data.error;
        this.$message.error(errorMsg ? `${user.username}: ${errorMsg}` : `${user.username} 更新用户状态失败`);
      });
    },
    toggleUserAdmin(user) {
      if (!this.isAdmin) {
        this.$message.error('您没有权限修改管理员权限');
        return;
      }
      if (this.isActionDisabled(user)) {
        this.$message.error('不允许修改自己的管理员权限');
        return;
      }

      const newIsAdmin = !user.is_admin;
      axios.put(`${process.env.VUE_APP_API_URL}/users/admin`, {
        id: user.id,
        is_admin: newIsAdmin,
      }, {
        headers: { Authorization: `Bearer ${this.$store.state.token}` },
      }).then(() => {
        user.is_admin = newIsAdmin;
        this.$message.success(`${user.username} ${newIsAdmin ? '已授予管理员权限' : '已取消管理员权限'}`);
      }).catch((error) => {
        const errorMsg = error && error.response && error.response.data && error.response.data.error;
        this.$message.error(errorMsg ? `${user.username}: ${errorMsg}` : `${user.username} 更新管理员权限失败`);
      });
    },
    handleDelete(user) {
      if (!this.isAdmin) {
        this.$message.error('您没有权限删除用户');
        return;
      }
      if (this.isActionDisabled(user)) {
        this.$message.error('不允许删除自己的账户');
        return;
      }

      this.$confirm(`此操作将永久删除用户 ${user.username}，是否继续？`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(() => {
        axios.delete(`${process.env.VUE_APP_API_URL}/users/${user.id}`, {
          headers: { Authorization: `Bearer ${this.$store.state.token}` },
        }).then(() => {
          this.fetchUsers();
          this.$message.success(`${user.username} 用户已删除`);
        });
      });
    },
    isActionDisabled(user) {
      return user.username === this.loggedInUser;
    },
  },
};
</script>

<style scoped>
.user-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-item__avatar {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: #eaf6f7;
  color: #3f8f94;
  font-size: 13px;
  font-weight: 700;
}

.user-item strong {
  display: block;
  font-size: 14px;
  color: #111827;
}

.user-item p {
  margin: 4px 0 0;
  font-size: 12px;
  color: #94a3b8;
}

.toggle-pill {
  min-height: 30px;
  padding: 0 10px;
  display: inline-flex;
  align-items: center;
  border: none;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}

.toggle-pill.is-admin {
  color: #3f8f94;
  background: #eaf6f7;
}

.toggle-pill.is-member {
  color: #64748b;
  background: #f8fafc;
}

.toggle-pill.is-enabled {
  color: #15803d;
  background: #ecfdf3;
}

.toggle-pill.is-disabled {
  color: #b45309;
  background: #fff7ed;
}

.row-actions {
  display: flex;
  justify-content: center;
  gap: 8px;
}
</style>
