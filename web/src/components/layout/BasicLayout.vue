<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div class="logo">学生管理系统</div>
      <div class="user-info" v-if="userStore.currentUser">
        <el-dropdown @command="handleCommand">
          <span class="el-dropdown-link">
            {{ userStore.currentUser.name }}
            <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人信息</el-dropdown-item>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>

    <el-container>
      <el-aside width="200px" class="aside">
        <el-menu
          :default-active="activeMenu"
          class="el-menu-vertical"
          :router="true"
        >
          <el-menu-item index="/app/dashboard">
            <el-icon><menu /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/app/class">
            <el-icon><reading /></el-icon>
            <span>班级管理</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-main class="main">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ArrowDown, Menu, Reading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const handleCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      userStore.logout()
      ElMessage.success('退出登录成功')
      router.push('/login')
      break
  }
}
</script>

<style scoped lang="scss">
.layout-container {
  height: 100vh;

  .header {
    background-color: #fff;
    border-bottom: 1px solid #dcdfe6;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 20px;

    .logo {
      font-size: 20px;
      font-weight: bold;
      color: #409eff;
    }

    .user-info {
      cursor: pointer;

      .el-dropdown-link {
        display: flex;
        align-items: center;
        color: #606266;
      }
    }
  }

  .aside {
    background-color: #fff;
    border-right: 1px solid #dcdfe6;

    .el-menu {
      border-right: none;
    }
  }

  .main {
    background-color: #f5f7fa;
    padding: 20px;
  }
}
</style>