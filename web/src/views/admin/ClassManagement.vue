<template>
  <div class="class-management">
    <div class="header">
      <h2>班级管理</h2>
      <el-button type="primary" @click="showAddDialog">添加班级</el-button>
    </div>

    <el-table :data="classList" border style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="班级名称" />
      <el-table-column prop="description" label="班级描述" />
      <el-table-column label="操作" width="180">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.page_size"
      :total="total"
      @current-change="fetchClassList"
      @size-change="handleSizeChange"
      layout="total, sizes, prev, pager, next"
    />

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="30%"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="班级名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入班级名称" />
        </el-form-item>
        <el-form-item label="班级描述" prop="description">
          <el-input v-model="form.description" type="textarea" placeholder="请输入班级描述" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import {
  addClass,
  editClass,
  deleteClass,
  getClassList
} from '@/api/class'
import type {
  Class,
  ClassAddRequest,
  ClassEditRequest,
  ClassListResponse,
  PaginationData
} from '@/types/class'

const classList = ref<Class[]>([])
const total = ref<number>(0)
const pagination = reactive({
  page: 1,
  page_size: 10
})
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const currentId = ref<number | null>(null)

const form = reactive<ClassAddRequest>({
  name: '',
  description: '',
  teacher_id: 0
})

const rules = {
  name: [
    { required: true, message: '请输入班级名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ]
}

onMounted(() => {
  fetchClassList()
})

const fetchClassList = async () => {
  try {
    const response = await getClassList(pagination) as any
    if (response && response.list) {
      classList.value = response.list
      total.value = response.total
    } else {
      ElMessage.error('获取班级列表失败：数据格式错误')
    }
  } catch (error: any) {
    console.error('获取班级列表失败:', error)
    ElMessage.error(error.message || '获取班级列表失败')
  }
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchClassList()
}

const showAddDialog = () => {
  dialogTitle.value = '添加班级'
  currentId.value = null
  form.name = ''
  form.description = ''
  form.teacher_id = 0
  dialogVisible.value = true
}

const handleEdit = (row: Class) => {
  dialogTitle.value = '编辑班级'
  currentId.value = row.id
  form.name = row.name
  form.description = row.description
  form.teacher_id = row.teacher_id
  dialogVisible.value = true
}

const handleDelete = (row: Class) => {
  ElMessageBox.confirm('确定要删除这个班级吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteClass(row.id)
      fetchClassList()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const submitForm = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    
    if (currentId.value === null) {
      // 添加班级
      await addClass(form)
    } else {
      // 编辑班级
      await editClass(currentId.value, {
        id: currentId.value,
        name: form.name,
        description: form.description,
        teacher_id: form.teacher_id
      })
    }

    dialogVisible.value = false
    fetchClassList()
  } catch (error) {
    ElMessage.error('表单验证失败')
  }
}
</script>

<style scoped>
.class-management {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.el-pagination {
  margin-top: 20px;
  text-align: right;
}
</style>