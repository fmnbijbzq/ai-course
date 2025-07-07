<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <h2 class="card-title">登录</h2>
      </template>

      <el-form
        ref="formRef"
        :model="loginForm"
        :rules="rules"
        label-width="80px"
        @submit.prevent="handleSubmit"
      >
        <el-form-item label="学号" prop="student_id">
          <el-input
            v-model="loginForm.student_id"
            placeholder="请输入学号"
          />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading" class="submit-btn">
            登录
          </el-button>
          <el-link type="primary" :underline="false" @click="router.push('/register')" class="register-link">
            没有账号？立即注册
          </el-link>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { login } from '@/api/user'
import { useUserStore } from '@/store/user'
import type { LoginForm } from '@/types/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const loginForm = reactive<LoginForm>({
  student_id: '',
  password: ''
})

const rules = {
  student_id: [
    { required: true, message: '请输入学号', trigger: 'blur' },
    { min: 5, message: '学号长度不能小于5位', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        loading.value = true
        const res = await login(loginForm)
        console.log('Login response:', res)

        // 后端返回的结构是 { code, message, data: { user, token } }
        if (!res.data || !res.data.token || !res.data.user) {
          throw new Error('登录响应数据不完整')
        }

        // 先设置用户信息和令牌
        userStore.setUser(res.data.user)
        userStore.setToken(res.data.token)

        // 显示登录成功消息
        ElMessage.success('登录成功')

        // 获取重定向地址
        const redirect = (router.currentRoute.value.query.redirect as string) || '/app/dashboard'
        
        // 使用 nextTick 确保状态更新后再跳转
        nextTick(() => {
          router.push(redirect).catch((err) => {
            console.error('Navigation error:', err)
            // 如果跳转失败，默认跳转到仪表盘
        router.push('/app/dashboard')
          })
        })
      } catch (error: any) {
        console.error('Login error:', error)
        ElMessage.error(error.message || '登录失败，请重试')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped lang="scss">
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;

  .login-card {
    width: 400px;

    .card-title {
      text-align: center;
      margin: 0;
      color: #303133;
    }

    .submit-btn {
      width: 100%;
      margin-bottom: 15px;
    }

    .register-link {
      display: block;
      text-align: center;
    }
  }
}
</style>