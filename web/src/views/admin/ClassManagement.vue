<template>
  <div class="class-management">
    <div class="header">
      <h2>班级管理</h2>
      <el-button type="primary" @click="showAddDialog">添加班级</el-button>
    </div>

    <el-table :data="classList" border style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="class_name" label="班级名称" />
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column prop="updated_at" label="更新时间" width="180" />
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
        <el-form-item label="班级名称" prop="class_name">
          <el-input v-model="form.class_name" placeholder="请输入班级名称" />
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
  ClassEditRequest
} from '@/types/class'

const classList = ref<Class[]>([])
const total = ref(0)
const pagination = reactive({
  page: 1,
  page_size: 10
})
const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const currentId = ref<number | null>(null)

const form = reactive<ClassAddRequest | ClassEditRequest>({
  class_name: ''
})

const rules = {
  class_name: [
    { required: true, message: '请输入班级名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ]
}

onMounted(() => {
  fetchClassList()
})

const fetchClassList = async () => {
  try {
    const res = await getClassList(pagination)
    classList.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    ElMessage.error('获取班级列表失败')
  }
}

const handleSizeChange = () => {
  // 当每页显示数量变化时，重置页码为1并重新获取数据
  pagination.page = 1
  fetchClassList()
}

const showAddDialog = () => {
  dialogTitle.value = '添加班级'
  currentId.value = null
  form.class_name = ''
  dialogVisible.value = true
}

const handleEdit = (row: Class) => {
  dialogTitle.value = '编辑班级'
  currentId.value = row.id
  form.class_name = row.class_name
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
      ElMessage.success('删除成功')
      fetchClassList()
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

const submitForm = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (currentId.value) {
          await editClass(currentId.value, form as ClassEditRequest)
          ElMessage.success('更新成功')
        } else {
          await addClass(form as ClassAddRequest)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        fetchClassList()
      } catch (error) {
        ElMessage.error('操作失败')
      }
    }
  })
}
</script>

<style scoped lang="scss">
.class-management {
  padding: 20px;

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .el-pagination {
    margin-top: 20px;
    justify-content: flex-end;
  }
}
</style>