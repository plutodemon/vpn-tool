<template>
  <div class="box">
    <div class="form">
      <el-form :model="ruleForm" status-icon ref="ruleForm" label-position="top" class="form1" :disabled="visible">
        <el-form-item label="用户名" prop="userName" class="form-item" :rules="rules.name">
          <el-input v-model="ruleForm.userName" class="form-input"></el-input>
        </el-form-item>
        <el-form-item label="邮 箱" prop="userEmail" class="form-item" :rules="rules.email">
          <el-input v-model="ruleForm.userEmail" class="form-input"></el-input>
        </el-form-item>
        <el-form-item label="密 码" prop="userPass" class="form-item" :rules="rules.name">
          <el-input type="password" v-model="ruleForm.userPass" autocomplete="off" class="form-input"></el-input>
        </el-form-item>
        <el-form-item label="网络编号" prop="netId" class="form-item" :rules="rules.name">
          <el-input v-model.number="ruleForm.netId" class="form-input"></el-input>
        </el-form-item>
        <el-form-item label="IP地址" prop="netIp" class="form-item" :rules="rules.ip">
          <el-input v-model="ruleForm.netIp" class="form-input"></el-input>
        </el-form-item>
        <el-form-item label="用户等级" prop="userLevel" class="form-item">
          <el-input v-model.number="ruleForm.userLevel" class="form-input"></el-input>
        </el-form-item>

        <el-form-item>
          <el-button @click="submitForm('ruleForm')" class="form-button">提交</el-button>
          <el-button @click="resetForm('ruleForm')" class="form-button">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <el-card class="card" shadow="always" v-show="visible">
      <img
          :src="url"
          class="image"
          alt="获取失败"/>
      <div style="padding: 14px">
        <span>密钥：</span>
        <div class="bottom">
          {{ key }}
        </div>
        <el-button text class="button" @click="handleClose">Operating</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import {defineComponent, markRaw} from "vue";
import {User} from "@element-plus/icons-vue";
import {ElButton, ElDrawer} from 'element-plus'
import {ElMessage, ElMessageBox} from "element-plus";
import axios from "axios";

export default defineComponent({
  name: "AddUser",
  data() {
    return {
      ruleForm: {
        userName: '',
        userPass: '',
        userEmail: '',
        netId: '',
        netIp: '',
        userLevel: ''
      },
      url: '',
      key: '',
      visible: false,
      rules: {
        name: [
          {required: true, message: '此项不能为空', trigger: 'change',},
        ],
        email: [
          {required: true, message: '请输入邮箱', trigger: 'blur'},
          {
            validator: function (rule, value, callback) {
              if (/^([a-zA-Z0-9]+[-_\.]?)+@[a-zA-Z0-9]+\.[a-z]+$/.test(value) === false) {
                callback(new Error("请输入正确邮箱"));
              } else {
                //校验通过
                callback();
              }
            }, trigger: 'blur'
          },
        ],
        ip: [
          {required: true, message: '请输入地址', trigger: 'blur'},
          {
            validator: function (rule, value, callback) {
              if (/^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/.test(value) === false) {
                callback(new Error("请输入正确地址"));
              } else {
                //校验通过
                callback();
              }
            }, trigger: 'blur'
          },
        ]
      }
    };
  },
  methods: {
    submitForm() {
      //进行表单验证
      this.$refs.ruleForm.validate(valid => {
        //表单验证失败
        if (!valid) {
          //提示语
          this.$message("请正确填写");
          return false;
        }
        ElMessageBox.confirm(
            '即将创建新用户，确认继续？',
            '提示',
            {
              confirmButtonText: '确认',
              cancelButtonText: '返回',
              type: 'success',
              icon: markRaw(User),
            }
        ).then(() => {
          console.log(this.ruleForm)
          axios.post("/api/ui/createUser", this.ruleForm).then(res => {
            console.log(res)
            this.url = res.data.url
            this.key = res.data.key
            this.visible = true
            if (res.status === 200) {
              ElMessage({
                type: 'success',
                message: '创建成功',
              })
            } else {
              ElMessage({
                type: 'danger',
                message: '创建失败',
              })
            }
          })
        })
            .catch(() => {
              ElMessage({
                type: 'info',
                message: '未创建',
              })
            })
      })
    },
    handleClose() {
      ElMessageBox.confirm('确定关闭嘛？')
          .then(() => {
            this.visible = false
            this.$refs.ruleForm.resetFields();
          })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    }
  }
})
</script>

<style scoped>
.box {
  display: flex;
  flex-direction: row;
  justify-content: space-around;
}

.form {
  display: flex;
  flex-direction: column;
  justify-content: center;
  width: 700px;
  height: 600px;
  margin-top: 20px;
  margin-left: 35px;
  margin-right: 35px;

  border-radius: 10px;
  background: rgba(238, 233, 233, 0.25);
  box-shadow: 0 8px 28px 0 rgba(137, 144, 232, 0.27);
  backdrop-filter: blur(7.5px);
  -webkit-backdrop-filter: blur(7.5px);
  border: 1px solid rgba(255, 255, 255, 0.18);
}

.form1 {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.form-item {
  width: 400px;
  margin-top: 10px;
  margin-bottom: 10px;
}

.form-input {
  height: 40px;
  background: rgba(238, 233, 233, 0.25);
  box-shadow: 0 8px 28px 0 rgba(137, 144, 232, 0.27);
  backdrop-filter: blur(7.5px);
  -webkit-backdrop-filter: blur(7.5px);
}

.form-button {
  width: 150px;
  margin-left: 20px;
  margin-right: 20px;
  box-shadow: 0 7px 7px 0 rgba(105, 112, 224, 0.37);
  backdrop-filter: blur(15.5px);
  -webkit-backdrop-filter: blur(15.5px);
  border-radius: 10px;
}

.card {
  border-radius: 10px;
  background: rgba(238, 233, 233, 0.25);
  box-shadow: 0 8px 28px 0 rgba(137, 144, 232, 0.27);
  backdrop-filter: blur(7.5px);
  -webkit-backdrop-filter: blur(7.5px);
  border: 1px solid rgba(255, 255, 255, 0.18);
  width: 1000px;
  height: 600px;
  display: flex;
  justify-content: center;
  flex-direction: column;
  margin-top: 20px;
  margin-left: 35px;
  margin-right: 15px;
  align-items: center;
}
</style>