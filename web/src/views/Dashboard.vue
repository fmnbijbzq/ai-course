<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="welcome-card">
          <template #header>
            <div class="card-header">
              <span>欢迎回来</span>
            </div>
          </template>
          <div class="welcome-content">
            <h3>{{ userStore.currentUser?.name }}</h3>
            <p>学号：{{ userStore.currentUser?.student_id }}</p>
          </div>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card class="stat-card">
          <template #header>
            <div class="card-header">
              <span>班级信息</span>
            </div>
          </template>
          <div class="stat-content">
            <el-statistic title="班级总数" :value="classCount">
              <template #suffix>
                <el-icon><School /></el-icon>
              </template>
            </el-statistic>
          </div>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card class="quick-actions">
          <template #header>
            <div class="card-header">
              <span>快捷操作</span>
            </div>
          </template>
          <div class="action-buttons">
            <el-button type="primary" @click="router.push('/app/class')">
              班级管理
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getClassList } from '@/api/class'
import { School } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const classCount = ref(0)

onMounted(async () => {
  try {
    const res = await getClassList({ page: 1, page_size: 1 })
    classCount.value = res.data.total
  } catch (error) {
    console.error('Failed to fetch class count:', error)
  }
})
</script>

<style scoped lang="scss">
.dashboard {
  padding: 20px;

  .el-card {
    margin-bottom: 20px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }

  .welcome-card {
    .welcome-content {
      text-align: center;
      padding: 20px 0;

      h3 {
        margin: 0 0 10px;
        font-size: 24px;
        color: #303133;
      }

      p {
        margin: 0;
        color: #606266;
      }
    }
  }

  .stat-card {
    .stat-content {
      padding: 20px 0;
      text-align: center;
    }
  }

  .quick-actions {
    .action-buttons {
      display: flex;
      justify-content: center;
      padding: 20px 0;

      .el-button {
        min-width: 120px;
      }
    }
  }
}
</style>