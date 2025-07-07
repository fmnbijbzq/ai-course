<template>
  <div class="profile-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>个人信息</h1>
    </div>

    <!-- 个人信息头部 -->
    <div class="profile-header">
      <div class="avatar-section">
        <el-avatar :size="120" :src="userInfo.avatar || undefined" class="profile-avatar">
          <i class="el-icon-user-solid"></i>
        </el-avatar>
        <el-button 
          size="small" 
          type="primary" 
          class="upload-btn"
          @click="handleAvatarUpload"
        >
          <i class="el-icon-camera"></i>
        </el-button>
      </div>
      <div class="user-info">
        <h2>{{ userInfo.name }}</h2>
        <p class="user-subtitle">{{ userInfo.student_id }} · {{ userInfo.department || '计算机学院' }}</p>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="profile-content">
      <!-- 基本信息卡片 -->
      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
            <el-button 
              type="text" 
              @click="editMode = !editMode"
              :icon="editMode ? 'el-icon-check' : 'el-icon-edit'"
            >
              {{ editMode ? '保存' : '编辑' }}
            </el-button>
          </div>
        </template>
        
        <el-form :model="userInfo" label-width="100px" class="profile-form">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="学号">
                <el-input 
                  v-model="userInfo.student_id" 
                  :disabled="true"
                  class="readonly-input"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="姓名">
                <el-input 
                  v-model="userInfo.name" 
                  :disabled="!editMode"
                />
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="性别">
                <el-select 
                  v-model="userInfo.gender" 
                  :disabled="!editMode"
                  placeholder="请选择性别"
                >
                  <el-option label="男" value="male"></el-option>
                  <el-option label="女" value="female"></el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="出生日期">
                <el-date-picker
                  v-model="userInfo.birth_date"
                  type="date"
                  placeholder="选择日期"
                  :disabled="!editMode"
                />
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="学院">
                <el-input 
                  v-model="userInfo.department" 
                  :disabled="!editMode"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="专业">
                <el-input 
                  v-model="userInfo.major" 
                  :disabled="!editMode"
                />
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="年级">
                <el-input 
                  v-model="userInfo.grade" 
                  :disabled="!editMode"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="班级">
                <el-input 
                  v-model="userInfo.class_name" 
                  :disabled="!editMode"
                />
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="邮箱">
                <el-input 
                  v-model="userInfo.email" 
                  :disabled="!editMode"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="手机">
                <el-input 
                  v-model="userInfo.phone" 
                  :disabled="!editMode"
                />
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
      </el-card>

      <!-- 修改密码卡片 -->
      <el-card class="password-card">
        <template #header>
          <span>修改密码</span>
        </template>
        
        <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
          <el-form-item label="当前密码" prop="currentPassword">
            <el-input 
              v-model="passwordForm.currentPassword" 
              type="password" 
              placeholder="请输入当前密码"
              show-password
            />
          </el-form-item>
          <el-form-item label="新密码" prop="newPassword">
            <el-input 
              v-model="passwordForm.newPassword" 
              type="password" 
              placeholder="请输入新密码"
              show-password
            />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input 
              v-model="passwordForm.confirmPassword" 
              type="password" 
              placeholder="请确认新密码"
              show-password
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handlePasswordChange">
              修改密码
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 账户安全卡片 -->
      <el-card class="security-card">
        <template #header>
          <span>账户安全</span>
        </template>
        
        <div class="security-items">
          <div class="security-item" @click="handleSecurityAction('phone')">
            <div class="security-info">
              <h4>绑定手机</h4>
              <p>{{ userInfo.phone ? `已绑定：${hidePhone(userInfo.phone)}` : '未绑定' }}</p>
            </div>
            <i class="el-icon-arrow-right"></i>
          </div>
          
          <div class="security-item" @click="handleSecurityAction('email')">
            <div class="security-info">
              <h4>绑定邮箱</h4>
              <p>{{ userInfo.email ? `已绑定：${hideEmail(userInfo.email)}` : '未绑定' }}</p>
            </div>
            <i class="el-icon-arrow-right"></i>
          </div>
          
          <div class="security-item" @click="handleSecurityAction('loginRecord')">
            <div class="security-info">
              <h4>登录记录</h4>
              <p>查看近期登录记录</p>
            </div>
            <i class="el-icon-arrow-right"></i>
          </div>
        </div>
      </el-card>

      <!-- 系统设置卡片 -->
      <el-card class="settings-card">
        <template #header>
          <span>系统设置</span>
        </template>
        
        <div class="settings-items">
          <div class="settings-item" @click="handleSettingsAction('notification')">
            <div class="settings-info">
              <h4>消息通知</h4>
              <p>设置系统消息、作业提醒等通知方式</p>
            </div>
            <i class="el-icon-arrow-right"></i>
          </div>
          
          <div class="settings-item" @click="handleSettingsAction('privacy')">
            <div class="settings-info">
              <h4>隐私设置</h4>
              <p>控制个人信息的可见范围</p>
            </div>
            <i class="el-icon-arrow-right"></i>
          </div>
          
          <div class="settings-item" @click="handleSettingsAction('about')">
            <div class="settings-info">
              <h4>关于系统</h4>
              <p>版本信息、使用帮助、联系我们</p>
            </div>
            <i class="el-icon-arrow-right"></i>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/store/user'
import type { User } from '@/types/user'

const userStore = useUserStore()
const editMode = ref(false)
const passwordFormRef = ref()

// 用户信息
const userInfo = reactive<User & {
  gender?: string
  birth_date?: Date
  department?: string
  major?: string
  grade?: string
  class_name?: string
  email?: string
  phone?: string
  avatar?: string
}>({
  id: 0,
  student_id: '',
  name: '',
  gender: '',
  birth_date: undefined,
  department: '',
  major: '',
  grade: '',
  class_name: '',
  email: '',
  phone: '',
  avatar: ''
})

// 密码修改表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码修改规则
const passwordRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { 
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ]
}

// 初始化用户信息
onMounted(() => {
  if (userStore.currentUser) {
    Object.assign(userInfo, userStore.currentUser)
  }
})

// 处理头像上传
const handleAvatarUpload = () => {
  ElMessage.info('头像上传功能暂未实现')
}

// 处理密码修改
const handlePasswordChange = () => {
  if (!passwordFormRef.value) return
  
  passwordFormRef.value.validate((valid: boolean) => {
    if (valid) {
      ElMessage.success('密码修改成功')
      // 重置表单
      Object.assign(passwordForm, {
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      })
    }
  })
}

// 处理安全设置点击
const handleSecurityAction = (action: string) => {
  switch (action) {
    case 'phone':
      ElMessage.info('手机绑定功能暂未实现')
      break
    case 'email':
      ElMessage.info('邮箱绑定功能暂未实现')
      break
    case 'loginRecord':
      ElMessage.info('登录记录功能暂未实现')
      break
  }
}

// 处理系统设置点击
const handleSettingsAction = (action: string) => {
  switch (action) {
    case 'notification':
      ElMessage.info('消息通知设置功能暂未实现')
      break
    case 'privacy':
      ElMessage.info('隐私设置功能暂未实现')
      break
    case 'about':
      ElMessage.info('关于系统功能暂未实现')
      break
  }
}

// 隐藏手机号
const hidePhone = (phone: string) => {
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// 隐藏邮箱
const hideEmail = (email: string) => {
  const [name, domain] = email.split('@')
  const hiddenName = name.substring(0, 2) + '***' + name.substring(name.length - 1)
  return hiddenName + '@' + domain
}
</script>

<style scoped lang="scss">
.profile-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
  
  h1 {
    font-size: 24px;
    font-weight: bold;
    color: #303133;
  }
}

.profile-header {
  background: linear-gradient(135deg, #409EFF, #36A3F7);
  color: white;
  padding: 30px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 20px;
  
  .avatar-section {
    position: relative;
    
    .profile-avatar {
      border: 3px solid white;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    
    .upload-btn {
      position: absolute;
      bottom: 0;
      right: 0;
      width: 32px;
      height: 32px;
      border-radius: 50%;
      padding: 0;
      border: 2px solid white;
    }
  }
  
  .user-info {
    h2 {
      margin: 0 0 8px 0;
      font-size: 28px;
      font-weight: bold;
    }
    
    .user-subtitle {
      margin: 0;
      font-size: 16px;
      opacity: 0.9;
    }
  }
}

.profile-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
  
  .info-card {
    grid-column: 1 / -1;
  }
  
  .password-card {
    grid-column: 1 / 2;
  }
  
  .security-card,
  .settings-card {
    grid-column: 2 / 3;
  }
}

.card-header {
  display: flex;
  justify-content: between;
  align-items: center;
  
  span {
    font-size: 18px;
    font-weight: 600;
  }
}

.profile-form {
  .readonly-input {
    :deep(.el-input__inner) {
      background-color: #f5f7fa;
      color: #909399;
    }
  }
}

.security-items,
.settings-items {
  .security-item,
  .settings-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px 0;
    border-bottom: 1px solid #EBEEF5;
    cursor: pointer;
    transition: background-color 0.3s;
    
    &:hover {
      background-color: #f5f7fa;
      border-radius: 4px;
      margin: 0 -10px;
      padding: 15px 10px;
    }
    
    &:last-child {
      border-bottom: none;
    }
    
    .security-info,
    .settings-info {
      h4 {
        margin: 0 0 4px 0;
        font-size: 14px;
        font-weight: 500;
        color: #303133;
      }
      
      p {
        margin: 0;
        font-size: 12px;
        color: #909399;
      }
    }
    
    .el-icon-arrow-right {
      color: #C0C4CC;
    }
  }
}

@media (max-width: 768px) {
  .profile-content {
    grid-template-columns: 1fr;
    
    .info-card,
    .password-card,
    .security-card,
    .settings-card {
      grid-column: 1 / -1;
    }
  }
  
  .profile-header {
    flex-direction: column;
    text-align: center;
    
    .user-info {
      margin-top: 10px;
    }
  }
}
</style>